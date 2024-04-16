### SG attached to test-vpc_be
resource "ibm_is_security_group" "test-vpc_be" {
  name           = "sg-test-vpc_be"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_vpc_id
}
# Internal. required-connections[0]: (instance test-vpc_fe)->(instance test-vpc_be); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc_be-0" {
  group     = ibm_is_security_group.test-vpc_be.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc_fe.id
  tcp {
  }
}

### SG attached to test-vpc_fe
resource "ibm_is_security_group" "test-vpc_fe" {
  name           = "sg-test-vpc_fe"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_vpc_id
}
# Internal. required-connections[0]: (instance test-vpc_fe)->(instance test-vpc_be); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc_fe-0" {
  group     = ibm_is_security_group.test-vpc_fe.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc_be.id
  tcp {
  }
}
