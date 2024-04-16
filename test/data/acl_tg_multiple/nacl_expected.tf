# test-vpc0_subnet0 [10.240.0.0/24]
resource "ibm_is_network_acl" "acl-test-vpc0_subnet0" {
  name           = "acl-test-vpc0_subnet0"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_vpc_id
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
  # Internal. required-connections[1]: (segment segment1)->(subnet test-vpc0_subnet3); allowed-protocols[0]
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

# test-vpc0_subnet2 [10.240.4.0/24]
resource "ibm_is_network_acl" "acl-test-vpc0_subnet2" {
  name           = "acl-test-vpc0_subnet2"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_vpc_id
  # Internal. required-connections[0]: (segment segment1)->(segment segment1); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.0.0/24"
    destination = "10.240.4.0/24"
  }
  # Internal. response to required-connections[0]: (segment segment1)->(segment segment1); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.4.0/24"
    destination = "10.240.0.0/24"
  }
  # Internal. required-connections[1]: (segment segment1)->(subnet test-vpc0_subnet3); allowed-protocols[0]
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

# test-vpc0_subnet3 [10.240.5.0/24]
resource "ibm_is_network_acl" "acl-test-vpc0_subnet3" {
  name           = "acl-test-vpc0_subnet3"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_vpc_id
  # Internal. required-connections[1]: (segment segment1)->(subnet test-vpc0_subnet3); allowed-protocols[0]
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
  # Internal. required-connections[1]: (segment segment1)->(subnet test-vpc0_subnet3); allowed-protocols[0]
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

# test-vpc0_subnet4 [10.240.8.0/24]
resource "ibm_is_network_acl" "acl-test-vpc0_subnet4" {
  name           = "acl-test-vpc0_subnet4"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_vpc_id
  # Internal. required-connections[2]: (subnet test-vpc0_subnet4)->(subnet test-vpc0_subnet5); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.8.0/24"
    destination = "10.240.9.0/24"
    icmp {
      type = 0
      code = 0
    }
  }
  # Internal. response to required-connections[2]: (subnet test-vpc0_subnet4)->(subnet test-vpc0_subnet5); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.9.0/24"
    destination = "10.240.8.0/24"
    icmp {
      type = 8
      code = 0
    }
  }
}

# test-vpc0_subnet5 [10.240.9.0/24]
resource "ibm_is_network_acl" "acl-test-vpc0_subnet5" {
  name           = "acl-test-vpc0_subnet5"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_vpc_id
  # Internal. required-connections[2]: (subnet test-vpc0_subnet4)->(subnet test-vpc0_subnet5); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.8.0/24"
    destination = "10.240.9.0/24"
    icmp {
      type = 0
      code = 0
    }
  }
  # Internal. response to required-connections[2]: (subnet test-vpc0_subnet4)->(subnet test-vpc0_subnet5); allowed-protocols[0]
  rules {
    name        = "rule1"
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
