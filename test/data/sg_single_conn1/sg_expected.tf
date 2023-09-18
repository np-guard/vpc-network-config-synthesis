# Internal. required-connections[0]: (nif ni3b)->(nif ni2); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-A-0" {
  group     = ibm_is_security_group.A.id
  direction = "outbound"
  remote    = ibm_is_security_group.B.id
  tcp {
    port_min = 443
    port_max = 443
  }
}
# Internal. required-connections[0]: (nif ni3b)->(nif ni2); allowed-protocols[0]
resource "ibm_is_security_group_rule" "sgrule-B-0" {
  group     = ibm_is_security_group.B.id
  direction = "inbound"
  remote    = ibm_is_security_group.A.id
  tcp {
    source_port_min = 443
    source_port_max = 443
  }
}
