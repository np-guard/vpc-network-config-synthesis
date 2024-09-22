# test-vpc0/subnet0 [10.240.0.0/24]
resource "ibm_is_network_acl" "acl-test-vpc0--subnet0" {
  name           = "acl-test-vpc0--subnet0"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc0_id
  # Internal. required-connections[0]: (subnet test-vpc0/subnet0)->(subnet test-vpc0/subnet1); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
<<<<<<< HEAD:test/expected/acl_protocols_tf/output.tf
    direction   = "outbound"
    source      = "10.240.0.0/24"
    destination = "10.240.1.0/24"
    tcp {
=======
    direction   = "inbound"
    source      = "10.240.8.0/24"
    destination = "10.240.0.0/24"
    icmp {
      type = 0
>>>>>>> main:test/data/acl_nif/nacl_expected.tf
    }
  }
  # Internal. response to required-connections[0]: (subnet test-vpc0/subnet0)->(subnet test-vpc0/subnet1); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.1.0/24"
    destination = "10.240.0.0/24"
    tcp {
    }
  }
  # Internal. required-connections[0]: (subnet test-vpc0/subnet0)->(subnet test-vpc0/subnet1); allowed-protocols[1]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.0.0/24"
<<<<<<< HEAD:test/expected/acl_protocols_tf/output.tf
    destination = "10.240.1.0/24"
=======
    destination = "10.240.8.0/24"
    icmp {
      type = 8
    }
>>>>>>> main:test/data/acl_nif/nacl_expected.tf
  }
}

# test-vpc0/subnet1 [10.240.1.0/24]
resource "ibm_is_network_acl" "acl-test-vpc0--subnet1" {
  name           = "acl-test-vpc0--subnet1"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc0_id
  # Internal. required-connections[0]: (subnet test-vpc0/subnet0)->(subnet test-vpc0/subnet1); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.0.0/24"
    destination = "10.240.1.0/24"
    tcp {
    }
  }
  # Internal. response to required-connections[0]: (subnet test-vpc0/subnet0)->(subnet test-vpc0/subnet1); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "10.240.0.0/24"
    tcp {
    }
  }
  # Internal. required-connections[0]: (subnet test-vpc0/subnet0)->(subnet test-vpc0/subnet1); allowed-protocols[1]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.0.0/24"
    destination = "10.240.1.0/24"
  }
}

# test-vpc0/subnet2 [10.240.4.0/24]
resource "ibm_is_network_acl" "acl-test-vpc0--subnet2" {
  name           = "acl-test-vpc0--subnet2"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc0_id
  # Internal. required-connections[1]: (subnet test-vpc0/subnet2)->(subnet test-vpc0/subnet3); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.4.0/24"
    destination = "10.240.5.0/24"
    tcp {
      port_min = 8080
      port_max = 8080
    }
  }
  # Internal. response to required-connections[1]: (subnet test-vpc0/subnet2)->(subnet test-vpc0/subnet3); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.5.0/24"
    destination = "10.240.4.0/24"
    tcp {
      source_port_min = 8080
      source_port_max = 8080
    }
  }
  # Internal. required-connections[1]: (subnet test-vpc0/subnet2)->(subnet test-vpc0/subnet3); allowed-protocols[1]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.4.0/24"
    destination = "10.240.5.0/24"
    icmp {
      type = 3
      code = 2
    }
  }
}

# test-vpc0/subnet3 [10.240.5.0/24]
resource "ibm_is_network_acl" "acl-test-vpc0--subnet3" {
  name           = "acl-test-vpc0--subnet3"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc0_id
  # Internal. required-connections[1]: (subnet test-vpc0/subnet2)->(subnet test-vpc0/subnet3); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.4.0/24"
    destination = "10.240.5.0/24"
    tcp {
      port_min = 8080
      port_max = 8080
    }
  }
  # Internal. response to required-connections[1]: (subnet test-vpc0/subnet2)->(subnet test-vpc0/subnet3); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.5.0/24"
    destination = "10.240.4.0/24"
    tcp {
      source_port_min = 8080
      source_port_max = 8080
    }
  }
  # Internal. required-connections[1]: (subnet test-vpc0/subnet2)->(subnet test-vpc0/subnet3); allowed-protocols[1]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.4.0/24"
    destination = "10.240.5.0/24"
    icmp {
      type = 3
      code = 2
    }
  }
}

# test-vpc0/subnet4 [10.240.8.0/24]
resource "ibm_is_network_acl" "acl-test-vpc0--subnet4" {
  name           = "acl-test-vpc0--subnet4"
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
<<<<<<< HEAD:test/expected/acl_protocols_tf/output.tf
      type = 15
      code = 0
=======
      type = 0
>>>>>>> main:test/data/acl_nif/nacl_expected.tf
    }
  }
  # Internal. response to required-connections[2]: (subnet test-vpc0/subnet4)->(subnet test-vpc0/subnet5); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.9.0/24"
    destination = "10.240.8.0/24"
    icmp {
<<<<<<< HEAD:test/expected/acl_protocols_tf/output.tf
      type = 16
      code = 0
=======
      type = 8
>>>>>>> main:test/data/acl_nif/nacl_expected.tf
    }
  }
  # Internal. required-connections[2]: (subnet test-vpc0/subnet4)->(subnet test-vpc0/subnet5); allowed-protocols[1]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.8.0/24"
    destination = "10.240.9.0/24"
    udp {
    }
  }
}

# test-vpc0/subnet5 [10.240.9.0/24]
resource "ibm_is_network_acl" "acl-test-vpc0--subnet5" {
  name           = "acl-test-vpc0--subnet5"
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
      type = 15
      code = 0
    }
  }
  # Internal. response to required-connections[2]: (subnet test-vpc0/subnet4)->(subnet test-vpc0/subnet5); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.9.0/24"
    destination = "10.240.8.0/24"
    icmp {
      type = 16
      code = 0
    }
  }
  # Internal. required-connections[2]: (subnet test-vpc0/subnet4)->(subnet test-vpc0/subnet5); allowed-protocols[1]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.8.0/24"
    destination = "10.240.9.0/24"
    udp {
    }
  }
}

# test-vpc1/subnet10 [10.240.64.0/24]
resource "ibm_is_network_acl" "acl-test-vpc1--subnet10" {
  name           = "acl-test-vpc1--subnet10"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc1_id
  # Internal. required-connections[3]: (subnet test-vpc1/subnet10)->(subnet test-vpc1/subnet11); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "10.240.80.0/24"
    udp {
      port_min = 53
      port_max = 53
    }
  }
}

# test-vpc1/subnet11 [10.240.80.0/24]
resource "ibm_is_network_acl" "acl-test-vpc1--subnet11" {
  name           = "acl-test-vpc1--subnet11"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc1_id
  # Internal. required-connections[3]: (subnet test-vpc1/subnet10)->(subnet test-vpc1/subnet11); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/24"
    destination = "10.240.80.0/24"
    udp {
      port_min = 53
      port_max = 53
    }
  }
}

# test-vpc2/subnet20 [10.240.128.0/24]
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

# test-vpc3/subnet30 [10.240.192.0/24]
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
