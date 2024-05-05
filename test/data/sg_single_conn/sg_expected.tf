### SG attached to test-vpc/be
resource "ibm_is_security_group" "test-vpc/be" {
  name           = "sg-test-vpc/be"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_vpc_id
}
# Internal. required-connections[0]: (instance test-vpc/fe)->(instance test-vpc/be); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc/be-0" {
  group     = ibm_is_security_group.test-vpc--be.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc--fe.id
  tcp {
  }
}

### SG attached to test-vpc/fe
resource "ibm_is_security_group" "test-vpc/fe" {
  name           = "sg-test-vpc/fe"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_vpc_id
}
# Internal. required-connections[0]: (instance test-vpc/fe)->(instance test-vpc/be); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc/fe-0" {
  group     = ibm_is_security_group.test-vpc--fe.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc--be.id
  tcp {
  }
}
