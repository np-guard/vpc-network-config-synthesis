### SG attached to test-vpc0/vsi0-subnet0
resource "ibm_is_security_group" "test-vpc0--vsi0-subnet0" {
  name           = "sg-test-vpc0/vsi0-subnet0"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.name_test-vpc0_id
}
# Internal. required-connections[0]: (instance test-vpc0/vsi0-subnet0)->(instance test-vpc0/vsi1-subnet4); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc0--vsi0-subnet0-0" {
  group     = ibm_is_security_group.test-vpc0--vsi0-subnet0.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc0--vsi1-subnet4.id
}

### SG attached to test-vpc0/vsi1-subnet4
resource "ibm_is_security_group" "test-vpc0--vsi1-subnet4" {
  name           = "sg-test-vpc0/vsi1-subnet4"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.name_test-vpc0_id
}
# Internal. required-connections[0]: (instance test-vpc0/vsi0-subnet0)->(instance test-vpc0/vsi1-subnet4); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc0--vsi1-subnet4-0" {
  group     = ibm_is_security_group.test-vpc0--vsi1-subnet4.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc0--vsi0-subnet0.id
}

### SG attached to test-vpc2/vsi0-subnet20
resource "ibm_is_security_group" "test-vpc2--vsi0-subnet20" {
  name           = "sg-test-vpc2/vsi0-subnet20"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.name_test-vpc2_id
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

### SG attached to test-vpc2/vsi2-subnet20
resource "ibm_is_security_group" "test-vpc2--vsi2-subnet20" {
  name           = "sg-test-vpc2/vsi2-subnet20"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.name_test-vpc2_id
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
