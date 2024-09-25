### SG attached to test-vpc1/vsi0-subnet10
resource "ibm_is_security_group" "test-vpc1--vsi0-subnet10" {
  name           = "sg-test-vpc1--vsi0-subnet10"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc1_id
}
# Internal. required-connections[2]: (instance test-vpc1/vsi0-subnet10)->(instance test-vpc1/vsi0-subnet11); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc1--vsi0-subnet10-0" {
  group     = ibm_is_security_group.test-vpc1--vsi0-subnet10.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc1--vsi0-subnet11.id
  tcp {
  }
}

### SG attached to test-vpc1/vsi0-subnet11
resource "ibm_is_security_group" "test-vpc1--vsi0-subnet11" {
  name           = "sg-test-vpc1--vsi0-subnet11"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc1_id
}
# Internal. required-connections[2]: (instance test-vpc1/vsi0-subnet10)->(instance test-vpc1/vsi0-subnet11); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc1--vsi0-subnet11-0" {
  group     = ibm_is_security_group.test-vpc1--vsi0-subnet11.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc1--vsi0-subnet10.id
  tcp {
  }
}
