# 1 [10.240.10.0/24]
resource "ibm_is_network_acl" "acl_1" {
  name           = "acl_1"
  resource_group = var.resource_group_id
  vpc            = var.vpc_id
  # Internal. required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
  rules {
    name        = "# rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.10.0/24"
    destination = "10.240.30.0/24"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  # Internal. required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
  rules {
    name        = "# rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.30.0/24"
    destination = "10.240.10.0/24"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
  # Internal. required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
  rules {
    name        = "# rule2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.10.0/24"
    destination = "10.240.30.0/24"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  # Internal. required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
  rules {
    name        = "# rule3"
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
    name        = "# rule4"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.30.0/24"
    destination = "10.240.10.0/24"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  # Internal. inverse of required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
  rules {
    name        = "# rule5"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.10.0/24"
    destination = "10.240.30.0/24"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
  # Internal. inverse of required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
  rules {
    name        = "# rule6"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.30.0/24"
    destination = "10.240.10.0/24"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  # Internal. inverse of required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
  rules {
    name        = "# rule7"
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
