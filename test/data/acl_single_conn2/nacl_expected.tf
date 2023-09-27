# subnet1-ky [10.240.10.0/24]
resource "ibm_is_network_acl" "acl-subnet1-ky" {
  name           = "acl-subnet1-ky"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_vpc_id
  # Internal. required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
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
  # Internal. response to required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
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
  # Internal. inverse of required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
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
  # Internal. response to inverse of required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
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

# subnet3-ky [10.240.30.0/24]
resource "ibm_is_network_acl" "acl-subnet3-ky" {
  name           = "acl-subnet3-ky"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_vpc_id
  # Internal. required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
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
  # Internal. response to required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
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
  # Internal. inverse of required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
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
  # Internal. response to inverse of required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
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
