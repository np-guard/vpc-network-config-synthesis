# Attached subnets: test-vpc3/subnet30
resource "ibm_is_network_acl" "test-vpc3--subnet30" {
  name           = "test-vpc3--subnet30"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc3_id
  # Deny all communication; subnet test-vpc3/subnet30[10.240.192.0/24] does not have required connections
  rules {
    name        = "rule0"
    action      = "deny"
    direction   = "inbound"
    source      = "0.0.0.0/0"
    destination = "10.240.192.0/24"
  }
  # Deny all communication; subnet test-vpc3/subnet30[10.240.192.0/24] does not have required connections
  rules {
    name        = "rule1"
    action      = "deny"
    direction   = "outbound"
    source      = "10.240.192.0/24"
    destination = "0.0.0.0/0"
  }
}
