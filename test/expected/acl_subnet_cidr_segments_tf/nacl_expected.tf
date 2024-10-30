# Attached subnets: testacl5-vpc/sub1-1
resource "ibm_is_network_acl" "acl-testacl5-vpc--sub1-1" {
  name           = "acl-testacl5-vpc--sub1-1"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_testacl5-vpc_id
  # Deny all communication; subnet testacl5-vpc/sub1-1[10.240.1.0/24] does not have required connections
  rules {
    name        = "rule0"
    action      = "deny"
    direction   = "inbound"
    source      = "0.0.0.0/0"
    destination = "10.240.1.0/24"
  }
  # Deny all communication; subnet testacl5-vpc/sub1-1[10.240.1.0/24] does not have required connections
  rules {
    name        = "rule1"
    action      = "deny"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "0.0.0.0/0"
  }
}

# Attached subnets: testacl5-vpc/sub1-2
resource "ibm_is_network_acl" "acl-testacl5-vpc--sub1-2" {
  name           = "acl-testacl5-vpc--sub1-2"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_testacl5-vpc_id
  # Internal. required-connections[0]: (segment cidrSegment)->(segment cidrSegment); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.2.0/24"
    destination = "10.240.2.0/23"
  }
  # Internal. response to required-connections[0]: (segment cidrSegment)->(segment cidrSegment); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.2.0/23"
    destination = "10.240.2.0/24"
  }
  # Internal. required-connections[1]: (segment cidrSegment)->(segment subnetSegment); allowed-protocols[0]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.2.0/24"
    destination = "10.240.64.0/24"
  }
  # Internal. response to required-connections[1]: (segment cidrSegment)->(segment subnetSegment); allowed-protocols[0]
  rules {
    name        = "rule3"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/24"
    destination = "10.240.2.0/24"
  }
  # Internal. required-connections[1]: (segment cidrSegment)->(segment subnetSegment); allowed-protocols[0]
  rules {
    name        = "rule4"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.2.0/24"
    destination = "10.240.65.0/24"
  }
  # Internal. response to required-connections[1]: (segment cidrSegment)->(segment subnetSegment); allowed-protocols[0]
  rules {
    name        = "rule5"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.65.0/24"
    destination = "10.240.2.0/24"
  }
}

# Attached subnets: testacl5-vpc/sub1-3
resource "ibm_is_network_acl" "acl-testacl5-vpc--sub1-3" {
  name           = "acl-testacl5-vpc--sub1-3"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_testacl5-vpc_id
  # Internal. required-connections[0]: (segment cidrSegment)->(segment cidrSegment); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.3.0/24"
    destination = "10.240.2.0/23"
  }
  # Internal. response to required-connections[0]: (segment cidrSegment)->(segment cidrSegment); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.2.0/23"
    destination = "10.240.3.0/24"
  }
  # Internal. required-connections[1]: (segment cidrSegment)->(segment subnetSegment); allowed-protocols[0]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.3.0/24"
    destination = "10.240.64.0/24"
  }
  # Internal. response to required-connections[1]: (segment cidrSegment)->(segment subnetSegment); allowed-protocols[0]
  rules {
    name        = "rule3"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/24"
    destination = "10.240.3.0/24"
  }
  # Internal. required-connections[1]: (segment cidrSegment)->(segment subnetSegment); allowed-protocols[0]
  rules {
    name        = "rule4"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.3.0/24"
    destination = "10.240.65.0/24"
  }
  # Internal. response to required-connections[1]: (segment cidrSegment)->(segment subnetSegment); allowed-protocols[0]
  rules {
    name        = "rule5"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.65.0/24"
    destination = "10.240.3.0/24"
  }
}

# Attached subnets: testacl5-vpc/sub2-1
resource "ibm_is_network_acl" "acl-testacl5-vpc--sub2-1" {
  name           = "acl-testacl5-vpc--sub2-1"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_testacl5-vpc_id
  # Internal. required-connections[1]: (segment cidrSegment)->(segment subnetSegment); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.2.0/23"
    destination = "10.240.64.0/24"
  }
  # Internal. response to required-connections[1]: (segment cidrSegment)->(segment subnetSegment); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "10.240.2.0/23"
  }
  # Internal. required-connections[2]: (segment subnetSegment)->(segment subnetSegment); allowed-protocols[0]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "10.240.65.0/24"
  }
  # Internal. response to required-connections[2]: (segment subnetSegment)->(segment subnetSegment); allowed-protocols[0]
  rules {
    name        = "rule3"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.65.0/24"
    destination = "10.240.64.0/24"
  }
}

# Attached subnets: testacl5-vpc/sub2-2
resource "ibm_is_network_acl" "acl-testacl5-vpc--sub2-2" {
  name           = "acl-testacl5-vpc--sub2-2"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_testacl5-vpc_id
  # Internal. required-connections[1]: (segment cidrSegment)->(segment subnetSegment); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.2.0/23"
    destination = "10.240.65.0/24"
  }
  # Internal. response to required-connections[1]: (segment cidrSegment)->(segment subnetSegment); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.65.0/24"
    destination = "10.240.2.0/23"
  }
  # Internal. required-connections[2]: (segment subnetSegment)->(segment subnetSegment); allowed-protocols[0]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.65.0/24"
    destination = "10.240.64.0/24"
  }
  # Internal. response to required-connections[2]: (segment subnetSegment)->(segment subnetSegment); allowed-protocols[0]
  rules {
    name        = "rule3"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/24"
    destination = "10.240.65.0/24"
  }
}

# Attached subnets: testacl5-vpc/sub3-1
resource "ibm_is_network_acl" "acl-testacl5-vpc--sub3-1" {
  name           = "acl-testacl5-vpc--sub3-1"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_testacl5-vpc_id
  # Deny all communication; subnet testacl5-vpc/sub3-1[10.240.128.0/24] does not have required connections
  rules {
    name        = "rule0"
    action      = "deny"
    direction   = "inbound"
    source      = "0.0.0.0/0"
    destination = "10.240.128.0/24"
  }
  # Deny all communication; subnet testacl5-vpc/sub3-1[10.240.128.0/24] does not have required connections
  rules {
    name        = "rule1"
    action      = "deny"
    direction   = "outbound"
    source      = "10.240.128.0/24"
    destination = "0.0.0.0/0"
  }
}
