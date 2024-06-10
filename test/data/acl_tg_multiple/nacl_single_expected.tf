# test-vpc0/singleACL [10.240.0.0/24]
resource "ibm_is_network_acl" "acl-test-vpc0--singleACL" {
  name           = "acl-test-vpc0--singleACL"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.name_test-vpc0_id
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
  # Internal. required-connections[0]: (segment segment1)->(segment segment1); allowed-protocols[0]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.0.0/24"
    destination = "10.240.4.0/24"
  }
  # Internal. response to required-connections[0]: (segment segment1)->(segment segment1); allowed-protocols[0]
  rules {
    name        = "rule3"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.4.0/24"
    destination = "10.240.0.0/24"
  }
  # Internal. required-connections[1]: (segment segment1)->(subnet test-vpc0/subnet3); allowed-protocols[0]
  rules {
    name        = "rule4"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.0.0/24"
    destination = "10.240.5.0/24"
    udp {
      port_min = 53
      port_max = 53
    }
  }
  # Internal. required-connections[1]: (segment segment1)->(subnet test-vpc0/subnet3); allowed-protocols[0]
  rules {
    name        = "rule5"
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
    name        = "rule6"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.4.0/24"
    destination = "10.240.5.0/24"
    udp {
      port_min = 53
      port_max = 53
    }
  }
  # Internal. required-connections[1]: (segment segment1)->(subnet test-vpc0/subnet3); allowed-protocols[0]
  rules {
    name        = "rule7"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.4.0/24"
    destination = "10.240.5.0/24"
    udp {
      port_min = 53
      port_max = 53
    }
  }
  # Internal. required-connections[2]: (subnet test-vpc0/subnet4)->(subnet test-vpc0/subnet5); allowed-protocols[0]
  rules {
    name        = "rule8"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.8.0/24"
    destination = "10.240.9.0/24"
    icmp {
      type = 0
      code = 0
    }
  }
  # Internal. response to required-connections[2]: (subnet test-vpc0/subnet4)->(subnet test-vpc0/subnet5); allowed-protocols[0]
  rules {
    name        = "rule9"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.9.0/24"
    destination = "10.240.8.0/24"
    icmp {
      type = 8
      code = 0
    }
  }
  # Internal. required-connections[2]: (subnet test-vpc0/subnet4)->(subnet test-vpc0/subnet5); allowed-protocols[0]
  rules {
    name        = "rule10"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.8.0/24"
    destination = "10.240.9.0/24"
    icmp {
      type = 0
      code = 0
    }
  }
  # Internal. response to required-connections[2]: (subnet test-vpc0/subnet4)->(subnet test-vpc0/subnet5); allowed-protocols[0]
  rules {
    name        = "rule11"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.9.0/24"
    destination = "10.240.8.0/24"
    icmp {
      type = 8
      code = 0
    }
  }
}

# test-vpc1/singleACL [10.240.64.0/24]
resource "ibm_is_network_acl" "acl-test-vpc1--singleACL" {
  name           = "acl-test-vpc1--singleACL"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.name_test-vpc1_id
  # Internal. required-connections[3]: (subnet test-vpc1/subnet10)->(subnet test-vpc1/subnet11); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "10.240.80.0/24"
    icmp {
      type = 0
      code = 0
    }
  }
  # Internal. response to required-connections[3]: (subnet test-vpc1/subnet10)->(subnet test-vpc1/subnet11); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.80.0/24"
    destination = "10.240.64.0/24"
    icmp {
      type = 8
      code = 0
    }
  }
  # Internal. required-connections[3]: (subnet test-vpc1/subnet10)->(subnet test-vpc1/subnet11); allowed-protocols[0]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/24"
    destination = "10.240.80.0/24"
    icmp {
      type = 0
      code = 0
    }
  }
  # Internal. response to required-connections[3]: (subnet test-vpc1/subnet10)->(subnet test-vpc1/subnet11); allowed-protocols[0]
  rules {
    name        = "rule3"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.80.0/24"
    destination = "10.240.64.0/24"
    icmp {
      type = 8
      code = 0
    }
  }
}
