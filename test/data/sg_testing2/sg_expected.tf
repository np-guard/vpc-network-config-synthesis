# Internal. required-connections[2]: (instance fe-ky)->(instance be-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-be-ky-0" {
  group     = be-ky.id
  direction = "inbound"
  remote    = fe-ky.id
  tcp {
  }
}
# Internal. required-connections[3]: (instance be-ky)->(instance opa-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-be-ky-1" {
  group     = be-ky.id
  direction = "outbound"
  remote    = opa-ky.id
  tcp {
  }
}
# Internal. required-connections[1]: (instance proxy-ky)->(instance fe-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-fe-ky-0" {
  group     = fe-ky.id
  direction = "inbound"
  remote    = proxy-ky.id
  tcp {
  }
}
# Internal. required-connections[2]: (instance fe-ky)->(instance be-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-fe-ky-1" {
  group     = fe-ky.id
  direction = "outbound"
  remote    = be-ky.id
  tcp {
  }
}
# Internal. required-connections[3]: (instance be-ky)->(instance opa-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-opa-ky-0" {
  group     = opa-ky.id
  direction = "inbound"
  remote    = be-ky.id
  tcp {
  }
}
# External. required-connections[0]: (external public internet)->(instance proxy-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-proxy-ky-0" {
  group     = proxy-ky.id
  direction = "inbound"
  remote    = "0.0.0.0/0"
}
# Internal. required-connections[1]: (instance proxy-ky)->(instance fe-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-proxy-ky-1" {
  group     = proxy-ky.id
  direction = "outbound"
  remote    = fe-ky.id
  tcp {
  }
}
