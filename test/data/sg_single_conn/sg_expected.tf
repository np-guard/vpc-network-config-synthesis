### SG attached to test-vpc/appdata-endpoint-gateway
resource "ibm_is_security_group" "test-vpc--appdata-endpoint-gateway" {
  name           = "sg-test-vpc--appdata-endpoint-gateway"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc_id
}

### SG attached to test-vpc/be
resource "ibm_is_security_group" "test-vpc--be" {
  name           = "sg-test-vpc--be"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc_id
}
# Internal. required-connections[0]: (instance test-vpc/fe)->(instance test-vpc/be); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--be-0" {
  group     = ibm_is_security_group.test-vpc--be.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc--fe.id
  tcp {
  }
}

### SG attached to test-vpc/fe
resource "ibm_is_security_group" "test-vpc--fe" {
  name           = "sg-test-vpc--fe"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc_id
}
# Internal. required-connections[0]: (instance test-vpc/fe)->(instance test-vpc/be); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--fe-0" {
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
  vpc            = local.sg_synth_test-vpc_id
}

### SG attached to test-vpc/policydb-endpoint-gateway
resource "ibm_is_security_group" "test-vpc--policydb-endpoint-gateway" {
  name           = "sg-test-vpc--policydb-endpoint-gateway"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc_id
}

### SG attached to test-vpc/proxy
resource "ibm_is_security_group" "test-vpc--proxy" {
  name           = "sg-test-vpc--proxy"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc_id
}
