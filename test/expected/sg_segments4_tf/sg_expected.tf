### SG test-vpc--appdata-endpoint-gateway is attached to test-vpc/appdata-endpoint-gateway
resource "ibm_is_security_group" "test-vpc--appdata-endpoint-gateway" {
  name           = "sg-test-vpc--appdata-endpoint-gateway"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc_id
}
# Internal. required-connections[0]: (segment vpeSegment)->(segment instanceSegment); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--appdata-endpoint-gateway-0" {
  group     = ibm_is_security_group.test-vpc--appdata-endpoint-gateway.id
  direction = "outbound"
  local     = "0.0.0.0/0"
  remote    = ibm_is_security_group.test-vpc--fe.id
}
# Internal. required-connections[0]: (segment vpeSegment)->(segment instanceSegment); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--appdata-endpoint-gateway-1" {
  group     = ibm_is_security_group.test-vpc--appdata-endpoint-gateway.id
  direction = "outbound"
  local     = "0.0.0.0/0"
  remote    = ibm_is_security_group.test-vpc--be.id
}

### SG test-vpc--be is attached to test-vpc/be
resource "ibm_is_security_group" "test-vpc--be" {
  name           = "sg-test-vpc--be"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc_id
}
# Internal. required-connections[0]: (segment vpeSegment)->(segment instanceSegment); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--be-0" {
  group     = ibm_is_security_group.test-vpc--be.id
  direction = "inbound"
  local     = "0.0.0.0/0"
  remote    = ibm_is_security_group.test-vpc--appdata-endpoint-gateway.id
}
# Internal. required-connections[0]: (segment vpeSegment)->(segment instanceSegment); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--be-1" {
  group     = ibm_is_security_group.test-vpc--be.id
  direction = "inbound"
  local     = "0.0.0.0/0"
  remote    = ibm_is_security_group.test-vpc--policydb-endpoint-gateway.id
}

### SG test-vpc--fe is attached to test-vpc/fe
resource "ibm_is_security_group" "test-vpc--fe" {
  name           = "sg-test-vpc--fe"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc_id
}
# Internal. required-connections[0]: (segment vpeSegment)->(segment instanceSegment); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--fe-0" {
  group     = ibm_is_security_group.test-vpc--fe.id
  direction = "inbound"
  local     = "0.0.0.0/0"
  remote    = ibm_is_security_group.test-vpc--appdata-endpoint-gateway.id
}
# Internal. required-connections[0]: (segment vpeSegment)->(segment instanceSegment); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--fe-1" {
  group     = ibm_is_security_group.test-vpc--fe.id
  direction = "inbound"
  local     = "0.0.0.0/0"
  remote    = ibm_is_security_group.test-vpc--policydb-endpoint-gateway.id
}

### SG test-vpc--opa is attached to test-vpc/opa
resource "ibm_is_security_group" "test-vpc--opa" {
  name           = "sg-test-vpc--opa"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc_id
}

### SG test-vpc--policydb-endpoint-gateway is attached to test-vpc/policydb-endpoint-gateway
resource "ibm_is_security_group" "test-vpc--policydb-endpoint-gateway" {
  name           = "sg-test-vpc--policydb-endpoint-gateway"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc_id
}
# Internal. required-connections[0]: (segment vpeSegment)->(segment instanceSegment); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--policydb-endpoint-gateway-0" {
  group     = ibm_is_security_group.test-vpc--policydb-endpoint-gateway.id
  direction = "outbound"
  local     = "0.0.0.0/0"
  remote    = ibm_is_security_group.test-vpc--fe.id
}
# Internal. required-connections[0]: (segment vpeSegment)->(segment instanceSegment); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc--policydb-endpoint-gateway-1" {
  group     = ibm_is_security_group.test-vpc--policydb-endpoint-gateway.id
  direction = "outbound"
  local     = "0.0.0.0/0"
  remote    = ibm_is_security_group.test-vpc--be.id
}

### SG test-vpc--proxy is attached to test-vpc/proxy
resource "ibm_is_security_group" "test-vpc--proxy" {
  name           = "sg-test-vpc--proxy"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc_id
}
