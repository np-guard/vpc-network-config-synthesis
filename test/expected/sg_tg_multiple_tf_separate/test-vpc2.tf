### SG attached to test-vpc2/vsi0-subnet20
resource "ibm_is_security_group" "test-vpc2--vsi0-subnet20" {
  name           = "sg-test-vpc2--vsi0-subnet20"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc2--vsi0-subnet20_id
}
# Internal. required-connections[1]: (instance test-vpc2/vsi0-subnet20)->(instance test-vpc2/vsi2-subnet20); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc2--vsi0-subnet20-0" {
  group     = ibm_is_security_group.test-vpc2--vsi0-subnet20.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc2--vsi2-subnet20.id
  tcp {
    port_min = 53
    port_max = 53
  }
}

### SG attached to test-vpc2/vsi1-subnet20
resource "ibm_is_security_group" "test-vpc2--vsi1-subnet20" {
  name           = "sg-test-vpc2--vsi1-subnet20"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc2--vsi1-subnet20_id
}

### SG attached to test-vpc2/vsi2-subnet20
resource "ibm_is_security_group" "test-vpc2--vsi2-subnet20" {
  name           = "sg-test-vpc2--vsi2-subnet20"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc2--vsi2-subnet20_id
}
# Internal. required-connections[1]: (instance test-vpc2/vsi0-subnet20)->(instance test-vpc2/vsi2-subnet20); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc2--vsi2-subnet20-0" {
  group     = ibm_is_security_group.test-vpc2--vsi2-subnet20.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc2--vsi0-subnet20.id
  tcp {
    port_min = 53
    port_max = 53
  }
}
