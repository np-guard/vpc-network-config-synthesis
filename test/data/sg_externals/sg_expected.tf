### SG attached to test-vpc/be
resource "ibm_is_security_group" "test-vpc--be" {
  name           = "sg-test-vpc--be"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.name_test-vpc_id
}
# External. required-connections[8]: (instance test-vpc/be)->(external external5); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--be-0" {
  group     = ibm_is_security_group.test-vpc--be.id
  direction = "outbound"
  remote    = "0.0.0.0"
}

### SG attached to test-vpc/fe
resource "ibm_is_security_group" "test-vpc--fe" {
  name           = "sg-test-vpc--fe"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.name_test-vpc_id
}
# External. required-connections[4]: (instance test-vpc/fe)->(external external1); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--fe-0" {
  group     = ibm_is_security_group.test-vpc--fe.id
  direction = "outbound"
  remote    = "0.0.0.0/0"
}
# External. required-connections[5]: (instance test-vpc/fe)->(external external2); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--fe-1" {
  group     = ibm_is_security_group.test-vpc--fe.id
  direction = "outbound"
  remote    = "8.8.8.8"
}
# External. required-connections[6]: (instance test-vpc/fe)->(external external3); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--fe-2" {
  group     = ibm_is_security_group.test-vpc--fe.id
  direction = "outbound"
  remote    = "7.7.7.7"
}
# External. required-connections[7]: (instance test-vpc/fe)->(external external4); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--fe-3" {
  group     = ibm_is_security_group.test-vpc--fe.id
  direction = "outbound"
  remote    = "5.5.0.0/16"
}

### SG attached to test-vpc/opa
resource "ibm_is_security_group" "test-vpc--opa" {
  name           = "sg-test-vpc--opa"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.name_test-vpc_id
}
# External. required-connections[9]: (external external5)->(instance test-vpc/opa); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--opa-0" {
  group     = ibm_is_security_group.test-vpc--opa.id
  direction = "inbound"
  remote    = "0.0.0.0"
}

### SG attached to test-vpc/proxy
resource "ibm_is_security_group" "test-vpc--proxy" {
  name           = "sg-test-vpc--proxy"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.name_test-vpc_id
}
# External. required-connections[0]: (external external1)->(instance test-vpc/proxy); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--proxy-0" {
  group     = ibm_is_security_group.test-vpc--proxy.id
  direction = "inbound"
  remote    = "0.0.0.0/0"
}
# External. required-connections[1]: (external external2)->(instance test-vpc/proxy); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--proxy-1" {
  group     = ibm_is_security_group.test-vpc--proxy.id
  direction = "inbound"
  remote    = "8.8.8.8"
}
# External. required-connections[2]: (external external3)->(instance test-vpc/proxy); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--proxy-2" {
  group     = ibm_is_security_group.test-vpc--proxy.id
  direction = "inbound"
  remote    = "7.7.7.7"
}
# External. required-connections[3]: (external external4)->(instance test-vpc/proxy); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--proxy-3" {
  group     = ibm_is_security_group.test-vpc--proxy.id
  direction = "inbound"
  remote    = "5.5.0.0/16"
}
