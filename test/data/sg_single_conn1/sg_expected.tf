# Internal. required-connections[0]: (nif ni3b)->(nif ni2); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-A-0" {
  group     = A-0.id
  direction = "outbound"
  remote    = "10.240.20.4"
  tcp {
    port_min = 443
    port_max = 443
  }
}
# Internal. required-connections[0]: (nif ni3b)->(nif ni2); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-B-0" {
  group     = B-0.id
  direction = "inbound"
  remote    = "10.240.30.4"
  tcp {
    source_port_min = 443
    source_port_max = 443
  }
}
