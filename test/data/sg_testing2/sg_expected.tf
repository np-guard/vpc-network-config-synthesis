# Internal. required-connections[2]: (instance fe-ky)->(instance be-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "be-ky-rule" {
  group     = be-ky.id
  direction = "inbound"
  remote    = "10.240.2.5"
  tcp {
  }
}
# Internal. required-connections[3]: (instance be-ky)->(instance opa-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "be-ky-rule" {
  group     = be-ky.id
  direction = "outbound"
  remote    = "10.240.2.4"
  tcp {
  }
}
# Internal. required-connections[1]: (instance proxy-ky)->(instance fe-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "fe-ky-rule" {
  group     = fe-ky.id
  direction = "inbound"
  remote    = "10.240.1.4"
  tcp {
  }
}
# Internal. required-connections[2]: (instance fe-ky)->(instance be-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "fe-ky-rule" {
  group     = fe-ky.id
  direction = "outbound"
  remote    = "10.240.2.6"
  tcp {
  }
}
# Internal. required-connections[3]: (instance be-ky)->(instance opa-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "opa-ky-rule" {
  group     = opa-ky.id
  direction = "inbound"
  remote    = "10.240.2.6"
  tcp {
  }
}
# External. required-connections[0]: (external public internet)->(instance proxy-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "proxy-ky-rule" {
  group     = proxy-ky.id
  direction = "inbound"
  remote    = "0.0.0.0/0"
}
# Internal. required-connections[1]: (instance proxy-ky)->(instance fe-ky); allowed-protocols[0]
resource "ibm_is_security_group_rule" "proxy-ky-rule" {
  group     = proxy-ky.id
  direction = "outbound"
  remote    = "10.240.2.5"
  tcp {
  }
}
