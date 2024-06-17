### SG attached to test-vpc3/vsi0-subnet30
resource "ibm_is_security_group" "test-vpc3--vsi0-subnet30" {
  name           = "sg-test-vpc3--vsi0-subnet30"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.name_test-vpc3_id
}
