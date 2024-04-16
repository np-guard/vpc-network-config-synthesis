### SG attached to test-vpc_be
resource "ibm_is_security_group" "test-vpc_be" {
  name           = "sg-test-vpc_be"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_vpc_id
}
# Internal. required-connections[2]: (instance test-vpc_fe)->(instance test-vpc_be); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc_be-0" {
  group     = ibm_is_security_group.test-vpc_be.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc_fe.id
  tcp {
  }
}
# Internal. required-connections[3]: (instance test-vpc_be)->(instance test-vpc_opa); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc_be-1" {
  group     = ibm_is_security_group.test-vpc_be.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc_opa.id
}

### SG attached to test-vpc_fe
resource "ibm_is_security_group" "test-vpc_fe" {
  name           = "sg-test-vpc_fe"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_vpc_id
}
# Internal. required-connections[1]: (instance test-vpc_proxy)->(instance test-vpc_fe); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc_fe-0" {
  group     = ibm_is_security_group.test-vpc_fe.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc_proxy.id
  tcp {
    port_min = 9000
    port_max = 9000
  }
}
# Internal. required-connections[2]: (instance test-vpc_fe)->(instance test-vpc_be); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc_fe-1" {
  group     = ibm_is_security_group.test-vpc_fe.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc_be.id
  tcp {
  }
}

### SG attached to test-vpc_opa
resource "ibm_is_security_group" "test-vpc_opa" {
  name           = "sg-test-vpc_opa"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_vpc_id
}
# Internal. required-connections[3]: (instance test-vpc_be)->(instance test-vpc_opa); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc_opa-0" {
  group     = ibm_is_security_group.test-vpc_opa.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc_be.id
}

### SG attached to test-vpc_proxy
resource "ibm_is_security_group" "test-vpc_proxy" {
  name           = "sg-test-vpc_proxy"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_vpc_id
}
# External. required-connections[0]: (external public internet)->(instance test-vpc_proxy); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc_proxy-0" {
  group     = ibm_is_security_group.test-vpc_proxy.id
  direction = "inbound"
  remote    = "0.0.0.0/0"
}
# Internal. required-connections[1]: (instance test-vpc_proxy)->(instance test-vpc_fe); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc_proxy-1" {
  group     = ibm_is_security_group.test-vpc_proxy.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc_fe.id
  tcp {
    port_min = 9000
    port_max = 9000
  }
}
