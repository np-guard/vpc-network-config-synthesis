### SG test-vpc2--vsi0-subnet20 is attached to test-vpc2/vsi0-subnet20
resource "ibm_is_security_group" "test-vpc2--vsi0-subnet20" {
  name           = "sg-test-vpc2--vsi0-subnet20"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc2_id
}
# Internal. required-connections[1]: (instance test-vpc2/vsi0-subnet20)->(instance test-vpc2/vsi2-subnet20); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc2--vsi0-subnet20-0" {
  group     = ibm_is_security_group.test-vpc2--vsi0-subnet20.id
  direction = "outbound"
  local     = "0.0.0.0/0"
  remote    = ibm_is_security_group.test-vpc2--vsi2-subnet20.id
  tcp {
    port_min = 53
    port_max = 53
  }
}

### SG test-vpc2--vsi1-subnet20 is attached to test-vpc2/vsi1-subnet20
resource "ibm_is_security_group" "test-vpc2--vsi1-subnet20" {
  name           = "sg-test-vpc2--vsi1-subnet20"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc2_id
}

### SG test-vpc2--vsi2-subnet20 is attached to test-vpc2/vsi2-subnet20
resource "ibm_is_security_group" "test-vpc2--vsi2-subnet20" {
  name           = "sg-test-vpc2--vsi2-subnet20"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc2_id
}
# Internal. required-connections[1]: (instance test-vpc2/vsi0-subnet20)->(instance test-vpc2/vsi2-subnet20); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc2--vsi2-subnet20-0" {
  group     = ibm_is_security_group.test-vpc2--vsi2-subnet20.id
  direction = "inbound"
  local     = "0.0.0.0/0"
  remote    = ibm_is_security_group.test-vpc2--vsi0-subnet20.id
  tcp {
    port_min = 53
    port_max = 53
  }
}
