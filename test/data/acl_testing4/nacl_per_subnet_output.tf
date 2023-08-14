resource "ibm_is_network_acl" "subnet1-acl" {
  name           = "subnet1-acl"
  resource_group = var.resource_group_id
  vpc            = var.vpc_id

  rules {
    name        = "out-1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.10.0/24"
    destination = "172.217.22.46"
  }
  rules {
    name        = "out-2"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.10.0/24"
    destination = "10.240.20.0/24"
  }
  rules {
    name        = "out-3"
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
    name        = "out-4"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.10.0/24"
    destination = "10.240.30.0/24"
    tcp {
      source_port_max = 443
      source_port_min = 443
    }
  }
  rules {
    name        = "in-1"
    action      = "allow"
    direction   = "inbound"
    source      = "172.217.22.46"
    destination = "10.240.10.0/24"
  }
  rules {
    name        = "in-2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.20.0/24"
    destination = "10.240.10.0/24"
  }
  rules {
    name        = "in-3"
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
    name        = "in-4"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.30.0/24"
    destination = "10.240.10.0/24"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
}

resource "ibm_is_network_acl" "subnet2-acl" {
  name           = "subnet2-acl"
  resource_group = var.resource_group_id
  vpc            = var.vpc_id

  rules {
    name        = "out-1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.20.0/24"
    destination = "0.0.0.0/0"
  }
  rules {
    name        = "out-2"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.20.0/24"
    destination = "10.240.10.0/24"
  }
  rules {
    name        = "in-1"
    action      = "allow"
    direction   = "inbound"
    source      = "0.0.0.0/0"
    destination = "10.240.20.0/24"
  }
  rules {
    name        = "in-2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.10.0/24"
    destination = "10.240.20.0/24"
  }
}

resource "ibm_is_network_acl" "subnet3-acl" {
  name           = "subnet3-acl"
  resource_group = var.resource_group_id
  vpc            = var.vpc_id

  rules {
    name        = "out-1"
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
    name        = "out-2"
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
    name        = "in-1"
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
    name        = "in-2"
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
