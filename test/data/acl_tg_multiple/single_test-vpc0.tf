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
    }
  }
  # Deny all communication; subnet test-vpc0/subnet1[10.240.1.0/24] does not have required connections
  rules {
    name        = "rule12"
    action      = "deny"
    direction   = "inbound"
    source      = "0.0.0.0/0"
    destination = "10.240.1.0/24"
  }
  # Deny all communication; subnet test-vpc0/subnet1[10.240.1.0/24] does not have required connections
  rules {
    name        = "rule13"
    action      = "deny"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "0.0.0.0/0"
  }
}
