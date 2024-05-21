# testacl5-vpc/sub1-2 [10.240.2.0/24]
resource "ibm_is_network_acl" "acl-testacl5-vpc--sub1-2" {
  name           = "acl-testacl5-vpc--sub1-2"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_vpc_id
  # Internal. required-connections[0]: (segment cidrSegment)->(subnet testacl5-vpc/sub3-1); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.2.0/24"
    destination = "10.240.128.0/24"
  }
  # Internal. response to required-connections[0]: (segment cidrSegment)->(subnet testacl5-vpc/sub3-1); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.128.0/24"
    destination = "10.240.2.0/24"
  }
  # Internal. required-connections[1]: (subnet testacl5-vpc/sub2-1)->(segment cidrSegment); allowed-protocols[0]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/24"
    destination = "10.240.2.0/24"
  }
  # Internal. response to required-connections[1]: (subnet testacl5-vpc/sub2-1)->(segment cidrSegment); allowed-protocols[0]
  rules {
    name        = "rule3"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.2.0/24"
    destination = "10.240.64.0/24"
  }
}

# testacl5-vpc/sub1-3 [10.240.3.0/24]
resource "ibm_is_network_acl" "acl-testacl5-vpc--sub1-3" {
  name           = "acl-testacl5-vpc--sub1-3"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_vpc_id
  # Internal. required-connections[0]: (segment cidrSegment)->(subnet testacl5-vpc/sub3-1); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.3.0/24"
    destination = "10.240.128.0/24"
  }
  # Internal. response to required-connections[0]: (segment cidrSegment)->(subnet testacl5-vpc/sub3-1); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.128.0/24"
    destination = "10.240.3.0/24"
  }
  # Internal. required-connections[1]: (subnet testacl5-vpc/sub2-1)->(segment cidrSegment); allowed-protocols[0]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/24"
    destination = "10.240.3.0/24"
  }
  # Internal. response to required-connections[1]: (subnet testacl5-vpc/sub2-1)->(segment cidrSegment); allowed-protocols[0]
  rules {
    name        = "rule3"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.3.0/24"
    destination = "10.240.64.0/24"
  }
}

# testacl5-vpc/sub2-1 [10.240.64.0/24]
resource "ibm_is_network_acl" "acl-testacl5-vpc--sub2-1" {
  name           = "acl-testacl5-vpc--sub2-1"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_vpc_id
  # Internal. required-connections[1]: (subnet testacl5-vpc/sub2-1)->(segment cidrSegment); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "10.240.3.0/23"
  }
  # Internal. response to required-connections[1]: (subnet testacl5-vpc/sub2-1)->(segment cidrSegment); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.3.0/23"
    destination = "10.240.64.0/24"
  }
}

# testacl5-vpc/sub3-1 [10.240.128.0/24]
resource "ibm_is_network_acl" "acl-testacl5-vpc--sub3-1" {
  name           = "acl-testacl5-vpc--sub3-1"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_vpc_id
  # Internal. required-connections[0]: (segment cidrSegment)->(subnet testacl5-vpc/sub3-1); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.3.0/23"
    destination = "10.240.128.0/24"
  }
  # Internal. response to required-connections[0]: (segment cidrSegment)->(subnet testacl5-vpc/sub3-1); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.128.0/24"
    destination = "10.240.3.0/23"
  }
}
