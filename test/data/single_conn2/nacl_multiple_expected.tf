resource "ibm_is_network_acl" "acl_10-240-10-0_24" {
  name           = "acl_10-240-10-0_24"
  resource_group = var.resource_group_id
  vpc            = var.vpc_id

  rules {
    # Internal. required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
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

  rules {
    # Internal. required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
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

  rules {
    # Internal. inverse of required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
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

  rules {
    # Internal. inverse of required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
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

resource "ibm_is_network_acl" "acl_10-240-30-0_24" {
  name           = "acl_10-240-30-0_24"
  resource_group = var.resource_group_id
  vpc            = var.vpc_id

  rules {
    # Internal. required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
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

  rules {
    # Internal. required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
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

  rules {
    # Internal. inverse of required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
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

  rules {
    # Internal. inverse of required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0]
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
