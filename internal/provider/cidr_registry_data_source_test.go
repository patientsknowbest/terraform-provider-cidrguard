package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// connect.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"cidr-guard": providerserver.NewProtocol6WithError(New("test")()),
}

func TestCidrGuardRegistryDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "cidr_guard_registry" "test" {
  networks = [
    {
      name        = "vpc-main"
      cidr        = "10.0.0.0/16"
      description = "Main VPC"
    },
    {
      name        = "vpc-secondary"
      cidr        = "10.1.0.0/16"
      description = "Secondary VPC"
    },
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cidr_guard_registry.test", "network.vpc-main.cidr", "10.0.0.0/16"),
					resource.TestCheckResourceAttr("data.cidr_guard_registry.test", "network.vpc-main.description", "Main VPC"),
					resource.TestCheckResourceAttr("data.cidr_guard_registry.test", "network.vpc-main.first_ip", "10.0.0.0"),
					resource.TestCheckResourceAttr("data.cidr_guard_registry.test", "network.vpc-main.last_ip", "10.0.255.255"),
					resource.TestCheckResourceAttr("data.cidr_guard_registry.test", "network.vpc-main.prefix", "10.0.0.0/16"),
					resource.TestCheckResourceAttr("data.cidr_guard_registry.test", "network.vpc-main.length", "16"),
					resource.TestCheckResourceAttr("data.cidr_guard_registry.test", "network.vpc-main.base_ip", "10.0.0.0"),
					resource.TestCheckResourceAttr("data.cidr_guard_registry.test", "network.vpc-main.count", "65536"),

					resource.TestCheckResourceAttr("data.cidr_guard_registry.test", "network.vpc-secondary.cidr", "10.1.0.0/16"),
					resource.TestCheckResourceAttr("data.cidr_guard_registry.test", "network.vpc-secondary.description", "Secondary VPC"),
					resource.TestCheckResourceAttr("data.cidr_guard_registry.test", "network.vpc-secondary.first_ip", "10.1.0.0"),
					resource.TestCheckResourceAttr("data.cidr_guard_registry.test", "network.vpc-secondary.last_ip", "10.1.255.255"),
					resource.TestCheckResourceAttr("data.cidr_guard_registry.test", "network.vpc-secondary.prefix", "10.1.0.0/16"),
					resource.TestCheckResourceAttr("data.cidr_guard_registry.test", "network.vpc-secondary.length", "16"),
					resource.TestCheckResourceAttr("data.cidr_guard_registry.test", "network.vpc-secondary.base_ip", "10.1.0.0"),
					resource.TestCheckResourceAttr("data.cidr_guard_registry.test", "network.vpc-secondary.count", "65536"),
				),
			},
		},
	})
}

func TestCidrGuardRegistryDataSource_duplicateNames(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "cidr_guard_registry" "test" {
  networks = [
    {
      name        = "duplicate-name"
      cidr        = "10.0.0.0/16"
      description = "First instance"
    },
    {
      name        = "duplicate-name"
      cidr        = "10.1.0.0/16"
      description = "Second instance"
    },
  ]
}
`,
				ExpectError: regexp.MustCompile(`(?s)The following network names are used more than once: 'duplicate-name'`),
			},
		},
	})
}

func TestCidrGuardRegistryDataSource_overlap(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "cidr_guard_registry" "test" {
  networks = [
    {
      name        = "network-a"
      cidr        = "10.0.0.0/16"
      description = "Network A"
    },
    {
      name        = "network-b"
      cidr        = "10.0.128.0/17"
      description = "Network B (overlaps with A)"
    },
  ]
}
`,
				ExpectError: regexp.MustCompile(`(?s)CIDR blocks for networks 'network-a'.*and 'network-b'.*overlap`),
			},
		},
	})
}
