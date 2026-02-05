terraform {
  required_providers {
    cidrguard = {
      source = "patientsknowbest/cidrguard"
    }
  }
}

provider "cidrguard" {
  # No provider configuration is required.
}
