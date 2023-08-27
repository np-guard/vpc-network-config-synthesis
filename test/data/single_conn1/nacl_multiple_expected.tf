# subnet1-ky [10.240.10.0/24]
resource "ibm_is_network_acl" "acl_subnet1-ky" {
  name           = "acl_subnet1-ky"
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
}

# subnet3-ky [10.240.30.0/24]
resource "ibm_is_network_acl" "acl_subnet3-ky" {
  name           = "acl_subnet3-ky"
  resource_group = var.resource_group_id
  vpc            = var.vpc_id
  # Internal. required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
  rules {
    name        = "# rule0"
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
    name        = "# rule1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.30.0/24"
    destination = "10.240.10.0/24"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
}
