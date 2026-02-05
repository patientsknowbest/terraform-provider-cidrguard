# Terraform CIDR Guard Provider

The `cidr-guard` provider offers a simple yet powerful data source to manage and validate a central registry of CIDR blocks. Its primary purpose is to prevent overlapping IP ranges in complex environments where multiple teams or modules might be allocating network space.

This provider is particularly useful for:
-   **Centralized IP Address Management (IPAM):** Define all your VPC or network CIDRs in one place as a single source of truth.
-   **Overlap Prevention:** The provider will fail `terraform plan` if any two CIDR blocks in the registry overlap, preventing costly network configuration errors before they happen.
-   **Data Enrichment:** It automatically calculates and exports useful information about each CIDR block, such as the first and last IP, the address count, and more.

## Example Usage

The following example defines a registry with two distinct network blocks. The provider validates that they do not overlap and exports their details.

```terraform
terraform {
  required_providers {
    cidr-guard = {
      source  = "adacasolutions/cidr-guard"
      version = "1.0.0" # Or your desired version
    }
  }
}

provider "cidr-guard" {
  # No provider configuration is required.
}

data "cidr_guard_registry" "main" {
  networks = [
    {
      name        = "vpc-main"
      cidr        = "10.0.0.0/16"
      description = "Main VPC for core services."
    },
    {
      name        = "vpc-analytics"
      cidr        = "10.1.0.0/16"
      description = "VPC for the analytics platform."
    },
  ]
}

# --- Using the Outputs ---

# Output the entire network map, keyed by the network name.
output "all_networks" {
  description = "The full map of all network details."
  value       = data.cidr_guard_registry.main.network
}

# Access details for a specific network directly by its name.
output "analytics_vpc_first_ip" {
  description = "The first usable IP address of the analytics VPC."
  value       = data.cidr_guard_registry.main.network["vpc-analytics"].first_ip
}

# The provider will return an error if configurations overlap.
# To test this, uncomment the following data source.
/*
data "cidr_guard_registry" "test_overlap" {
  networks = [
    {
      name = "network-a"
      cidr = "192.168.0.0/16"
    },
    {
      name = "network-b-overlaps"
      cidr = "192.168.128.0/17"
    }
  ]
}
*/
```

## Schema Reference

### `cidr_guard_registry` Data Source

#### Arguments

-   `networks` (Required, List of Objects): A list of network blocks to register and validate. Each object in the list has the following attributes:
    -   `name` (Required, String): A unique name for the network block. This name will be used as the key in the output map.
    -   `cidr` (Required, String): The CIDR notation for the network block (e.g., `10.0.0.0/16`).
    -   `description` (Optional, String): A description of the network's purpose.

#### Attributes

-   `network` (Computed, Map of Objects): A map containing the details of each validated network block, keyed by the `name` provided in the input. Each object in the map has the following attributes:
    -   `cidr` (String): The original CIDR block string.
    -   `description` (String): The description of the network block.
    -   `first_ip` (String): The first IP address in the CIDR range.
    -   `last_ip` (String): The last IP address in the CIDR range.
    -   `base_ip` (String): The network address (the first IP).
    -   `prefix` (String): The CIDR prefix in the format `A.B.C.D/L`.
    -   `length` (Number): The length of the CIDR prefix (e.g., `16`).
    -   `count` (Number): The total number of IP addresses in the CIDR block.
