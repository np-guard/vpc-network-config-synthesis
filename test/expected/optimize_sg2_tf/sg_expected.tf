### SG test-vpc1--vsi1 is attached to ni1
resource "ibm_is_security_group" "test-vpc1--vsi1" {
  name           = "sg-test-vpc1--vsi1"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc1_id
}
resource "ibm_is_security_group_rule" "test-vpc1--vsi1-0" {
  group     = ibm_is_security_group.test-vpc1--vsi1.id
  direction = "outbound"
  local     = "0.0.0.0/0"
  remote    = "1.1.1.0/30"
  tcp {
  }
}
resource "ibm_is_security_group_rule" "test-vpc1--vsi1-1" {
  group     = ibm_is_security_group.test-vpc1--vsi1.id
  direction = "outbound"
  local     = "0.0.0.0/0"
  remote    = "1.1.1.1"
}
resource "ibm_is_security_group_rule" "test-vpc1--vsi1-2" {
  group     = ibm_is_security_group.test-vpc1--vsi1.id
  direction = "outbound"
  local     = "0.0.0.0/0"
  remote    = "1.1.1.3"
}
