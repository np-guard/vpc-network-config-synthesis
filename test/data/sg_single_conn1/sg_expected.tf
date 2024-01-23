### SG attached to A
resource "ibm_is_security_group" "A" {
  name           = "sg-A"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_vpc_id
}
# Internal. required-connections[0]: (nif ni3b)->(nif ni2); allowed-protocols[0]
resource "ibm_is_security_group_rule" "A-0" {
  group     = ibm_is_security_group.A.id
  direction = "outbound"
  remote    = ibm_is_security_group.B.id
  tcp {
    port_min = 443
    port_max = 443
  }
}

### SG attached to B
resource "ibm_is_security_group" "B" {
  name           = "sg-B"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_vpc_id
}
# Internal. required-connections[0]: (nif ni3b)->(nif ni2); allowed-protocols[0]
resource "ibm_is_security_group_rule" "B-0" {
  group     = ibm_is_security_group.B.id
  direction = "inbound"
  remote    = ibm_is_security_group.A.id
  tcp {
    port_min = 443
    port_max = 443
  }
}
