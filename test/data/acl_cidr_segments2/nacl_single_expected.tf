resource "ibm_is_network_acl" "acl-testacl5-vpc--singleACL" {
  name           = "acl-testacl5-vpc--singleACL"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.name_testacl5-vpc_id
  # Internal. required-connections[0]: (segment cidrSegment1)->(segment cidrSegment2); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.2.0/24"
    destination = "10.240.64.0/23"
  }
  # Internal. response to required-connections[0]: (segment cidrSegment1)->(segment cidrSegment2); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/23"
    destination = "10.240.2.0/24"
  }
  # Internal. required-connections[0]: (segment cidrSegment1)->(segment cidrSegment2); allowed-protocols[0]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.3.0/24"
    destination = "10.240.64.0/23"
  }
  # Internal. response to required-connections[0]: (segment cidrSegment1)->(segment cidrSegment2); allowed-protocols[0]
  rules {
    name        = "rule3"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/23"
    destination = "10.240.3.0/24"
  }
  # Internal. required-connections[0]: (segment cidrSegment1)->(segment cidrSegment2); allowed-protocols[0]
  rules {
    name        = "rule4"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.2.0/23"
    destination = "10.240.64.0/24"
  }
  # Internal. response to required-connections[0]: (segment cidrSegment1)->(segment cidrSegment2); allowed-protocols[0]
  rules {
    name        = "rule5"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "10.240.2.0/23"
  }
  # Internal. required-connections[0]: (segment cidrSegment1)->(segment cidrSegment2); allowed-protocols[0]
  rules {
    name        = "rule6"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.2.0/23"
    destination = "10.240.65.0/24"
  }
  # Internal. response to required-connections[0]: (segment cidrSegment1)->(segment cidrSegment2); allowed-protocols[0]
  rules {
    name        = "rule7"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.65.0/24"
    destination = "10.240.2.0/23"
  }
}
