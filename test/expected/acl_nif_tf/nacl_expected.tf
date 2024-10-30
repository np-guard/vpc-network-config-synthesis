# Attached subnets: test-vpc0/subnet0
resource "ibm_is_network_acl" "acl-test-vpc0--subnet0" {
  name           = "acl-test-vpc0--subnet0"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc0_id
  # Internal. required-connections[0]: (subnet test-vpc0/subnet4)->(instance test-vpc0/vsi0-subnet0); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.8.0/24"
    destination = "10.240.0.0/24"
    icmp {
      type = 0
    }
  }
  # Internal. response to required-connections[0]: (subnet test-vpc0/subnet4)->(instance test-vpc0/vsi0-subnet0); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.0.0/24"
    destination = "10.240.8.0/24"
    icmp {
      type = 8
    }
  }
}

# Attached subnets: test-vpc0/subnet1
resource "ibm_is_network_acl" "acl-test-vpc0--subnet1" {
  name           = "acl-test-vpc0--subnet1"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc0_id
  # Internal. required-connections[1]: (nif test-vpc0/vsi0-subnet1/scale-clambake-endearing-abridged)->(subnet test-vpc0/subnet4); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "10.240.8.0/24"
    tcp {
    }
  }
  # Internal. response to required-connections[1]: (nif test-vpc0/vsi0-subnet1/scale-clambake-endearing-abridged)->(subnet test-vpc0/subnet4); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.8.0/24"
    destination = "10.240.1.0/24"
    tcp {
    }
  }
}

# Attached subnets: test-vpc0/subnet2
resource "ibm_is_network_acl" "acl-test-vpc0--subnet2" {
  name           = "acl-test-vpc0--subnet2"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc0_id
  # Deny all communication; subnet test-vpc0/subnet2[10.240.4.0/24] does not have required connections
  rules {
    name        = "rule0"
    action      = "deny"
    direction   = "inbound"
    source      = "0.0.0.0/0"
    destination = "10.240.4.0/24"
  }
  # Deny all communication; subnet test-vpc0/subnet2[10.240.4.0/24] does not have required connections
  rules {
    name        = "rule1"
    action      = "deny"
    direction   = "outbound"
    source      = "10.240.4.0/24"
    destination = "0.0.0.0/0"
  }
}

# Attached subnets: test-vpc0/subnet3
resource "ibm_is_network_acl" "acl-test-vpc0--subnet3" {
  name           = "acl-test-vpc0--subnet3"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc0_id
  # Deny all communication; subnet test-vpc0/subnet3[10.240.5.0/24] does not have required connections
  rules {
    name        = "rule0"
    action      = "deny"
    direction   = "inbound"
    source      = "0.0.0.0/0"
    destination = "10.240.5.0/24"
  }
  # Deny all communication; subnet test-vpc0/subnet3[10.240.5.0/24] does not have required connections
  rules {
    name        = "rule1"
    action      = "deny"
    direction   = "outbound"
    source      = "10.240.5.0/24"
    destination = "0.0.0.0/0"
  }
}

# Attached subnets: test-vpc0/subnet4
resource "ibm_is_network_acl" "acl-test-vpc0--subnet4" {
  name           = "acl-test-vpc0--subnet4"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc0_id
  # Internal. required-connections[0]: (subnet test-vpc0/subnet4)->(instance test-vpc0/vsi0-subnet0); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.8.0/24"
    destination = "10.240.0.0/24"
    icmp {
      type = 0
    }
  }
  # Internal. response to required-connections[0]: (subnet test-vpc0/subnet4)->(instance test-vpc0/vsi0-subnet0); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.0.0/24"
    destination = "10.240.8.0/24"
    icmp {
      type = 8
    }
  }
  # Internal. required-connections[1]: (nif test-vpc0/vsi0-subnet1/scale-clambake-endearing-abridged)->(subnet test-vpc0/subnet4); allowed-protocols[0]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.1.0/24"
    destination = "10.240.8.0/24"
    tcp {
    }
  }
  # Internal. response to required-connections[1]: (nif test-vpc0/vsi0-subnet1/scale-clambake-endearing-abridged)->(subnet test-vpc0/subnet4); allowed-protocols[0]
  rules {
    name        = "rule3"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.8.0/24"
    destination = "10.240.1.0/24"
    tcp {
    }
  }
}

# Attached subnets: test-vpc0/subnet5
resource "ibm_is_network_acl" "acl-test-vpc0--subnet5" {
  name           = "acl-test-vpc0--subnet5"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc0_id
  # Deny all communication; subnet test-vpc0/subnet5[10.240.9.0/24] does not have required connections
  rules {
    name        = "rule0"
    action      = "deny"
    direction   = "inbound"
    source      = "0.0.0.0/0"
    destination = "10.240.9.0/24"
  }
  # Deny all communication; subnet test-vpc0/subnet5[10.240.9.0/24] does not have required connections
  rules {
    name        = "rule1"
    action      = "deny"
    direction   = "outbound"
    source      = "10.240.9.0/24"
    destination = "0.0.0.0/0"
  }
}

# Attached subnets: test-vpc1/subnet10
resource "ibm_is_network_acl" "acl-test-vpc1--subnet10" {
  name           = "acl-test-vpc1--subnet10"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc1_id
  # Deny all communication; subnet test-vpc1/subnet10[10.240.64.0/24] does not have required connections
  rules {
    name        = "rule0"
    action      = "deny"
    direction   = "inbound"
    source      = "0.0.0.0/0"
    destination = "10.240.64.0/24"
  }
  # Deny all communication; subnet test-vpc1/subnet10[10.240.64.0/24] does not have required connections
  rules {
    name        = "rule1"
    action      = "deny"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "0.0.0.0/0"
  }
}

# Attached subnets: test-vpc1/subnet11
resource "ibm_is_network_acl" "acl-test-vpc1--subnet11" {
  name           = "acl-test-vpc1--subnet11"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc1_id
  # Deny all communication; subnet test-vpc1/subnet11[10.240.80.0/24] does not have required connections
  rules {
    name        = "rule0"
    action      = "deny"
    direction   = "inbound"
    source      = "0.0.0.0/0"
    destination = "10.240.80.0/24"
  }
  # Deny all communication; subnet test-vpc1/subnet11[10.240.80.0/24] does not have required connections
  rules {
    name        = "rule1"
    action      = "deny"
    direction   = "outbound"
    source      = "10.240.80.0/24"
    destination = "0.0.0.0/0"
  }
}

# Attached subnets: test-vpc2/subnet20
resource "ibm_is_network_acl" "acl-test-vpc2--subnet20" {
  name           = "acl-test-vpc2--subnet20"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc2_id
  # Deny all communication; subnet test-vpc2/subnet20[10.240.128.0/24] does not have required connections
  rules {
    name        = "rule0"
    action      = "deny"
    direction   = "inbound"
    source      = "0.0.0.0/0"
    destination = "10.240.128.0/24"
  }
  # Deny all communication; subnet test-vpc2/subnet20[10.240.128.0/24] does not have required connections
  rules {
    name        = "rule1"
    action      = "deny"
    direction   = "outbound"
    source      = "10.240.128.0/24"
    destination = "0.0.0.0/0"
  }
}

# Attached subnets: test-vpc3/subnet30
resource "ibm_is_network_acl" "acl-test-vpc3--subnet30" {
  name           = "acl-test-vpc3--subnet30"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc3_id
  # Deny all communication; subnet test-vpc3/subnet30[10.240.192.0/24] does not have required connections
  rules {
    name        = "rule0"
    action      = "deny"
    direction   = "inbound"
    source      = "0.0.0.0/0"
    destination = "10.240.192.0/24"
  }
  # Deny all communication; subnet test-vpc3/subnet30[10.240.192.0/24] does not have required connections
  rules {
    name        = "rule1"
    action      = "deny"
    direction   = "outbound"
    source      = "10.240.192.0/24"
    destination = "0.0.0.0/0"
  }
}
