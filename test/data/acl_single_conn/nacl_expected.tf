# test-vpc1/subnet1 [10.240.10.0/24]
resource "ibm_is_network_acl" "acl-test-vpc1--subnet1" {
  name           = "acl-test-vpc1--subnet1"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc1_id
  # Internal. required-connections[0]: (subnet test-vpc1/subnet1)->(subnet test-vpc1/subnet3); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.10.0/24"
    destination = "10.240.30.0/24"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  # Internal. response to required-connections[0]: (subnet test-vpc1/subnet1)->(subnet test-vpc1/subnet3); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.30.0/24"
    destination = "10.240.10.0/24"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
  # Internal. inverse of required-connections[0]: (subnet test-vpc1/subnet1)->(subnet test-vpc1/subnet3); allowed-protocols[0]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.30.0/24"
    destination = "10.240.10.0/24"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  # Internal. response to inverse of required-connections[0]: (subnet test-vpc1/subnet1)->(subnet test-vpc1/subnet3); allowed-protocols[0]
  rules {
    name        = "rule3"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.10.0/24"
    destination = "10.240.30.0/24"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
}

# test-vpc1/subnet2 [10.240.20.0/24]
resource "ibm_is_network_acl" "acl-test-vpc1--subnet2" {
  name           = "acl-test-vpc1--subnet2"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc1_id
  # Deny all communication; subnet test-vpc1/subnet2[10.240.20.0/24] does not have required connections
  rules {
    name        = "rule0"
    action      = "deny"
    direction   = "inbound"
    source      = "0.0.0.0/0"
    destination = "10.240.20.0/24"
  }
  # Deny all communication; subnet test-vpc1/subnet2[10.240.20.0/24] does not have required connections
  rules {
    name        = "rule1"
    action      = "deny"
    direction   = "outbound"
    source      = "10.240.20.0/24"
    destination = "0.0.0.0/0"
  }
}

# test-vpc1/subnet3 [10.240.30.0/24]
resource "ibm_is_network_acl" "acl-test-vpc1--subnet3" {
  name           = "acl-test-vpc1--subnet3"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc1_id
  # Internal. required-connections[0]: (subnet test-vpc1/subnet1)->(subnet test-vpc1/subnet3); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.10.0/24"
    destination = "10.240.30.0/24"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  # Internal. response to required-connections[0]: (subnet test-vpc1/subnet1)->(subnet test-vpc1/subnet3); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.30.0/24"
    destination = "10.240.10.0/24"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
  # Internal. inverse of required-connections[0]: (subnet test-vpc1/subnet1)->(subnet test-vpc1/subnet3); allowed-protocols[0]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.30.0/24"
    destination = "10.240.10.0/24"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  # Internal. response to inverse of required-connections[0]: (subnet test-vpc1/subnet1)->(subnet test-vpc1/subnet3); allowed-protocols[0]
  rules {
    name        = "rule3"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.10.0/24"
    destination = "10.240.30.0/24"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
}
