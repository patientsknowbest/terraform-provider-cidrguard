terraform {
  required_providers {
    cidr-guard = {
      source = "adacasolutions/cidr-guard"
    }
  }
}

provider "cidr-guard" {
  # No provider configuration is required.
}
