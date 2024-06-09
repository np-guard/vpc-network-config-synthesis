### SG attached to test-vpc/appdata-endpoint-gateway
resource "ibm_is_security_group" "test-vpc--appdata-endpoint-gateway" {
  name           = "sg-test-vpc--appdata-endpoint-gateway"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.name_test-vpc_id
}

### SG attached to test-vpc/be
resource "ibm_is_security_group" "test-vpc--be" {
  name           = "sg-test-vpc--be"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.name_test-vpc_id
}
# Internal. required-connections[2]: (instance test-vpc/fe)->(instance test-vpc/be); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--be-0" {
  group     = ibm_is_security_group.test-vpc--be.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc--fe.id
  tcp {
  }
}
# Internal. required-connections[3]: (instance test-vpc/be)->(instance test-vpc/opa); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--be-1" {
  group     = ibm_is_security_group.test-vpc--be.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc--opa.id
}
# Internal. required-connections[4]: (instance test-vpc/be)->(vpe test-vpc/policydb-endpoint-gateway); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--be-2" {
  group     = ibm_is_security_group.test-vpc--be.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc--policydb-endpoint-gateway.id
}

### SG attached to test-vpc/fe
resource "ibm_is_security_group" "test-vpc--fe" {
  name           = "sg-test-vpc--fe"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.name_test-vpc_id
}
# Internal. required-connections[1]: (instance test-vpc/proxy)->(instance test-vpc/fe); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--fe-0" {
  group     = ibm_is_security_group.test-vpc--fe.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc--proxy.id
  tcp {
    port_min = 9000
    port_max = 9000
  }
}
# Internal. required-connections[2]: (instance test-vpc/fe)->(instance test-vpc/be); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--fe-1" {
  group     = ibm_is_security_group.test-vpc--fe.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc--be.id
  tcp {
  }
}

### SG attached to test-vpc/opa
resource "ibm_is_security_group" "test-vpc--opa" {
  name           = "sg-test-vpc--opa"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.name_test-vpc_id
}
# Internal. required-connections[3]: (instance test-vpc/be)->(instance test-vpc/opa); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--opa-0" {
  group     = ibm_is_security_group.test-vpc--opa.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc--be.id
}
# Internal. required-connections[5]: (instance test-vpc/opa)->(vpe test-vpc/policydb-endpoint-gateway); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--opa-1" {
  group     = ibm_is_security_group.test-vpc--opa.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc--policydb-endpoint-gateway.id
}

### SG attached to test-vpc/policydb-endpoint-gateway
resource "ibm_is_security_group" "test-vpc--policydb-endpoint-gateway" {
  name           = "sg-test-vpc--policydb-endpoint-gateway"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.name_test-vpc_id
}
# Internal. required-connections[4]: (instance test-vpc/be)->(vpe test-vpc/policydb-endpoint-gateway); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--policydb-endpoint-gateway-0" {
  group     = ibm_is_security_group.test-vpc--policydb-endpoint-gateway.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc--be.id
}
# Internal. required-connections[5]: (instance test-vpc/opa)->(vpe test-vpc/policydb-endpoint-gateway); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--policydb-endpoint-gateway-1" {
  group     = ibm_is_security_group.test-vpc--policydb-endpoint-gateway.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc--opa.id
}

### SG attached to test-vpc/proxy
resource "ibm_is_security_group" "test-vpc--proxy" {
  name           = "sg-test-vpc--proxy"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.name_test-vpc_id
}
# External. required-connections[0]: (external public internet)->(instance test-vpc/proxy); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--proxy-0" {
  group     = ibm_is_security_group.test-vpc--proxy.id
  direction = "inbound"
  remote    = "0.0.0.0/0"
}
# Internal. required-connections[1]: (instance test-vpc/proxy)->(instance test-vpc/fe); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--proxy-1" {
  group     = ibm_is_security_group.test-vpc--proxy.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc--fe.id
  tcp {
    port_min = 9000
    port_max = 9000
  }
}
