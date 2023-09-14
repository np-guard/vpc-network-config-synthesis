# Internal. required-connections[2]: (instance fe-ky)->(instance be-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-be-ky-0" {
  group     = be-ky-0.id
  direction = "inbound"
  remote    = "10.240.2.5"
  tcp {
  }
}
# Internal. required-connections[3]: (instance be-ky)->(instance opa-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-be-ky-1" {
  group     = be-ky-1.id
  direction = "outbound"
  remote    = "10.240.2.4"
  tcp {
  }
}
# Internal. required-connections[1]: (instance proxy-ky)->(instance fe-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-fe-ky-0" {
  group     = fe-ky-0.id
  direction = "inbound"
  remote    = "10.240.1.4"
  tcp {
  }
}
# Internal. required-connections[2]: (instance fe-ky)->(instance be-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-fe-ky-1" {
  group     = fe-ky-1.id
  direction = "outbound"
  remote    = "10.240.2.6"
  tcp {
  }
}
# Internal. required-connections[3]: (instance be-ky)->(instance opa-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-opa-ky-0" {
  group     = opa-ky-0.id
  direction = "inbound"
  remote    = "10.240.2.6"
  tcp {
  }
}
# External. required-connections[0]: (external public internet)->(instance proxy-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-proxy-ky-0" {
  group     = proxy-ky-0.id
  direction = "inbound"
  remote    = "0.0.0.0/0"
}
# Internal. required-connections[1]: (instance proxy-ky)->(instance fe-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-proxy-ky-1" {
  group     = proxy-ky-1.id
  direction = "outbound"
  remote    = "10.240.2.5"
  tcp {
  }
}
