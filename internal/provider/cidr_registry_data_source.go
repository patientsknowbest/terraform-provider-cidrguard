package provider

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/adacasolutions/terraform-provider-cidr-guard/internal/cidr"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &cidrGuardRegistryDataSource{}
)

// NewCidrGuardRegistryDataSource is a helper function to simplify the provider implementation.
func NewCidrGuardRegistryDataSource() datasource.DataSource {
	return &cidrGuardRegistryDataSource{}
}

// cidrGuardRegistryDataSource is the data source implementation.
type cidrGuardRegistryDataSource struct{}

// cidrGuardRegistryDataSourceModel maps the data source schema data.
type cidrGuardRegistryDataSourceModel struct {
	Networks []networkModel `tfsdk:"networks"`
	Network  types.Map      `tfsdk:"network"`
}

// networkModel maps the network block schema data.
type networkModel struct {
	Name        types.String `tfsdk:"name"`
	CIDR        types.String `tfsdk:"cidr"`
	Description types.String `tfsdk:"description"`
}

// networkDetailModel maps the output object schema data.
type networkDetailModel struct {
	CIDR        types.String `tfsdk:"cidr"`
	Description types.String `tfsdk:"description"`
	FirstIP     types.String `tfsdk:"first_ip"`
	LastIP      types.String `tfsdk:"last_ip"`
	Prefix      types.String `tfsdk:"prefix"`
	Length      types.Int64  `tfsdk:"length"`
	BaseIP      types.String `tfsdk:"base_ip"`
	Count       types.Number `tfsdk:"count"`
}

var networkDetailType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"cidr":        types.StringType,
		"description": types.StringType,
		"first_ip":    types.StringType,
		"last_ip":     types.StringType,
		"prefix":      types.StringType,
		"length":      types.Int64Type,
		"base_ip":     types.StringType,
		"count":       types.NumberType,
	},
}

// Metadata returns the data source type name.
func (d *cidrGuardRegistryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registry"
}

// Schema defines the schema for the data source.
func (d *cidrGuardRegistryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a way to manage and validate a registry of CIDR blocks, ensuring they do not overlap.",
		Attributes: map[string]schema.Attribute{
			"networks": schema.ListNestedAttribute{
				Description: "A list of network block definitions to be registered and validated.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "A unique name for this network block.",
							Required:    true,
						},
						"cidr": schema.StringAttribute{
							Description: "The CIDR block string (e.g., '10.0.0.0/16').",
							Required:    true,
						},
						"description": schema.StringAttribute{
							Description: "A description for this network block.",
							Optional:    true,
						},
					},
				},
			},
			"network": schema.MapAttribute{
				Description: "The details of each allocated network block, keyed by the network name.",
				Computed:    true,
				ElementType: networkDetailType,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *cidrGuardRegistryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state cidrGuardRegistryDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that network names are unique and collect them.
	nameCounts := make(map[string]int)
	for _, net := range state.Networks {
		nameCounts[net.Name.ValueString()]++
	}

	var duplicateNames []string
	for name, count := range nameCounts {
		if count > 1 {
			duplicateNames = append(duplicateNames, fmt.Sprintf("'%s'", name))
		}
	}

	if len(duplicateNames) > 0 {
		resp.Diagnostics.AddError(
			"Duplicate Network Names",
			fmt.Sprintf("The following network names are used more than once: %s. Network names must be unique.", strings.Join(duplicateNames, ", ")),
		)
		return
	}

	var cidrRanges []*cidr.Range
	for _, net := range state.Networks {
		cr, err := cidr.NewRange(
			net.Name.ValueString(),
			net.CIDR.ValueString(),
			net.Description.ValueString(),
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"CIDR Parsing Error",
				fmt.Sprintf("Could not parse CIDR for network '%s': %s", net.Name.ValueString(), err.Error()),
			)
			return
		}
		cidrRanges = append(cidrRanges, cr)
	}

	if err := cidr.ValidateNoOverlap(cidrRanges); err != nil {
		resp.Diagnostics.AddError("CIDR Overlap Error", err.Error())
		return
	}

	networkMap := make(map[string]attr.Value)
	for _, cr := range cidrRanges {
		networkDetail := networkDetailModel{
			CIDR:        types.StringValue(cr.CIDR),
			Description: types.StringValue(cr.Description),
			FirstIP:     types.StringValue(cr.FirstIP.String()),
			LastIP:      types.StringValue(cr.LastIP.String()),
			Prefix:      types.StringValue(cr.Prefix),
			Length:      types.Int64Value(int64(cr.Length)),
			BaseIP:      types.StringValue(cr.FirstIP.String()),
			Count:       types.NumberValue(new(big.Float).SetInt(cr.Count)),
		}

		obj, diags := types.ObjectValueFrom(ctx, networkDetailType.AttrTypes, &networkDetail)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		networkMap[cr.Name] = obj
	}

	networksMap, diags := types.MapValue(networkDetailType, networkMap)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	state.Network = networksMap

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
