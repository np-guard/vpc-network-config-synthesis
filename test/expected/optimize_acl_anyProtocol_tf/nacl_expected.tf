# Attached subnets: sub1-1
resource "ibm_is_network_acl" "testacl5-vpc--sub1-1" {
  name           = "testacl5-vpc--sub1-1"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_testacl5-vpc_id
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "1.1.1.0/31"
  }
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "2.2.2.0/30"
  }
}
