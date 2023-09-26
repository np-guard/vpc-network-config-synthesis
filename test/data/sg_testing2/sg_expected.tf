# Internal. required-connections[4]: (instance be-ky)->(vpe appdata-endpoint-gateway); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-appdata-endpoint-gateway-0" {
  group     = ibm_is_security_group.appdata-endpoint-gateway.id
  direction = "inbound"
  remote    = ibm_is_security_group.be-ky.id
  tcp {
  }
}
# Internal. required-connections[2]: (instance fe-ky)->(instance be-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-be-ky-0" {
  group     = ibm_is_security_group.be-ky.id
  direction = "inbound"
  remote    = ibm_is_security_group.fe-ky.id
  tcp {
  }
}
# Internal. required-connections[3]: (instance be-ky)->(instance opa-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-be-ky-1" {
  group     = ibm_is_security_group.be-ky.id
  direction = "outbound"
  remote    = ibm_is_security_group.opa-ky.id
  tcp {
    port_min = 8181
    port_max = 8181
  }
}
# Internal. required-connections[4]: (instance be-ky)->(vpe appdata-endpoint-gateway); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-be-ky-2" {
  group     = ibm_is_security_group.be-ky.id
  direction = "outbound"
  remote    = ibm_is_security_group.appdata-endpoint-gateway.id
  tcp {
  }
}
# Internal. required-connections[1]: (instance proxy-ky)->(instance fe-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-fe-ky-0" {
  group     = ibm_is_security_group.fe-ky.id
  direction = "inbound"
  remote    = ibm_is_security_group.proxy-ky.id
  tcp {
    port_min = 9000
    port_max = 9000
  }
}
# Internal. required-connections[2]: (instance fe-ky)->(instance be-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-fe-ky-1" {
  group     = ibm_is_security_group.fe-ky.id
  direction = "outbound"
  remote    = ibm_is_security_group.be-ky.id
  tcp {
  }
}
# Internal. required-connections[3]: (instance be-ky)->(instance opa-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-opa-ky-0" {
  group     = ibm_is_security_group.opa-ky.id
  direction = "inbound"
  remote    = ibm_is_security_group.be-ky.id
  tcp {
    port_min = 8181
    port_max = 8181
  }
}
# Internal. required-connections[5]: (instance opa-ky)->(vpe policydb-endpoint-gateway); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-opa-ky-1" {
  group     = ibm_is_security_group.opa-ky.id
  direction = "outbound"
  remote    = ibm_is_security_group.policydb-endpoint-gateway.id
  tcp {
  }
}
# Internal. required-connections[5]: (instance opa-ky)->(vpe policydb-endpoint-gateway); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-policydb-endpoint-gateway-0" {
  group     = ibm_is_security_group.policydb-endpoint-gateway.id
  direction = "inbound"
  remote    = ibm_is_security_group.opa-ky.id
  tcp {
  }
}
# External. required-connections[0]: (external public internet)->(instance proxy-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-proxy-ky-0" {
  group     = ibm_is_security_group.proxy-ky.id
  direction = "inbound"
  remote    = "0.0.0.0/0"
}
# Internal. required-connections[1]: (instance proxy-ky)->(instance fe-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-proxy-ky-1" {
  group     = ibm_is_security_group.proxy-ky.id
  direction = "outbound"
  remote    = ibm_is_security_group.fe-ky.id
  tcp {
    port_min = 9000
    port_max = 9000
  }
}
