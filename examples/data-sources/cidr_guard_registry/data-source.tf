# data "cidr_guard_registry" "main" {
#   networks = [
#     {
#       name        = "vpc-main"
#       cidr        = "10.0.0.0/16"
#       description = "Main VPC"
#     },
#     {
#       name        = "vpc-secondary"
#       cidr        = "10.1.0.0/16"
#       description = "Secondary VPC"
#     },
#   ]
# }
#
# output "all_networks" {
#   description = "The full map of all network details, keyed by network name."
#   value       = data.cidr_guard_registry.main.network
# }
#
# output "main_vpc_first_ip" {
#   description = "The first IP of the main VPC, accessed directly by its name."
#   value       = data.cidr_guard_registry.main.network["vpc-main"].first_ip
# }
