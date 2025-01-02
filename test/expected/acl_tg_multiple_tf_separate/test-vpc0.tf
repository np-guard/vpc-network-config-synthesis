# Attached subnets: test-vpc0/subnet0
resource "ibm_is_network_acl" "test-vpc0--subnet0" {
  name           = "test-vpc0--subnet0"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc0_id
  # Internal. required-connections[0]: (segment segment1)->(segment segment1); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.0.0/24"
    destination = "10.240.4.0/24"
  }
  # Internal. response to required-connections[0]: (segment segment1)->(segment segment1); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.4.0/24"
    destination = "10.240.0.0/24"
  }
  # Internal. required-connections[1]: (segment segment1)->(subnet test-vpc0/subnet3); allowed-protocols[0]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.0.0/24"
    destination = "10.240.5.0/24"
    udp {
      port_min = 53
      port_max = 53
    }
  }
}

# Attached subnets: test-vpc0/subnet1
resource "ibm_is_network_acl" "test-vpc0--subnet1" {
  name           = "test-vpc0--subnet1"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc0_id
  # Deny all communication; subnet test-vpc0/subnet1[10.240.1.0/24] does not have required connections
  rules {
    name        = "rule0"
    action      = "deny"
    direction   = "inbound"
    source      = "0.0.0.0/0"
    destination = "10.240.1.0/24"
  }
  # Deny all communication; subnet test-vpc0/subnet1[10.240.1.0/24] does not have required connections
  rules {
    name        = "rule1"
    action      = "deny"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "0.0.0.0/0"
  }
}

# Attached subnets: test-vpc0/subnet2
resource "ibm_is_network_acl" "test-vpc0--subnet2" {
  name           = "test-vpc0--subnet2"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc0_id
  # Internal. required-connections[0]: (segment segment1)->(segment segment1); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.4.0/24"
    destination = "10.240.0.0/24"
  }
  # Internal. response to required-connections[0]: (segment segment1)->(segment segment1); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.0.0/24"
    destination = "10.240.4.0/24"
  }
  # Internal. required-connections[1]: (segment segment1)->(subnet test-vpc0/subnet3); allowed-protocols[0]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.4.0/24"
    destination = "10.240.5.0/24"
    udp {
      port_min = 53
      port_max = 53
    }
  }
}

# Attached subnets: test-vpc0/subnet3
resource "ibm_is_network_acl" "test-vpc0--subnet3" {
  name           = "test-vpc0--subnet3"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc0_id
  # Internal. required-connections[1]: (segment segment1)->(subnet test-vpc0/subnet3); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.0.0/24"
    destination = "10.240.5.0/24"
    udp {
      port_min = 53
      port_max = 53
    }
  }
  # Internal. required-connections[1]: (segment segment1)->(subnet test-vpc0/subnet3); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.4.0/24"
    destination = "10.240.5.0/24"
    udp {
      port_min = 53
      port_max = 53
    }
  }
}

# Attached subnets: test-vpc0/subnet4
resource "ibm_is_network_acl" "test-vpc0--subnet4" {
  name           = "test-vpc0--subnet4"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc0_id
  # Internal. required-connections[2]: (subnet test-vpc0/subnet4)->(subnet test-vpc0/subnet5); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.8.0/24"
    destination = "10.240.9.0/24"
    icmp {
      type = 4
    }
  }
}

# Attached subnets: test-vpc0/subnet5
resource "ibm_is_network_acl" "test-vpc0--subnet5" {
  name           = "test-vpc0--subnet5"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc0_id
  # Internal. required-connections[2]: (subnet test-vpc0/subnet4)->(subnet test-vpc0/subnet5); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.8.0/24"
    destination = "10.240.9.0/24"
    icmp {
      type = 4
    }
  }
}
