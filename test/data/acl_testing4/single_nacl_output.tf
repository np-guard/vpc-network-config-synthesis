# Permissive network ACLs
resource "ibm_is_network_acl" "acl-for-test-vpc1-ky" {
  name = "acl-for-test-vpc1-ky"
  resource_group = var.resource_group_id
  vpc  = var.vpc_id
  rules {
    name        = "subnet1-out-1"
    action      = "allow"
    source      = "10.240.10.0/24"
    destination = "172.217.22.46"
    direction   = "outbound"
  }
  rules {
    name        = "subnet1-out-2"
    action      = "allow"
    source      = "10.240.10.0/24"
    destination = "10.240.20.0/24"
    direction   = "outbound"
  }
  rules {
    name        = "subnet1-out-3"
    action      = "allow"
    source      = "10.240.10.0/24"
    destination = "10.240.30.0/24"
    direction   = "outbound"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  rules {
    name        = "subnet1-out-4"
    action      = "allow"
    source      = "10.240.10.0/24"
    destination = "10.240.30.0/24"
    direction   = "outbound"
    tcp {
      source_port_max = 443
      source_port_min = 443
    }
  }
  rules {
    name        = "subnet1-in-1"
    action      = "allow"
    source      = "172.217.22.46"
    destination = "10.240.10.0/24"
    direction   = "inbound"
  }
  rules {
    name        = "subnet1-in-2"
    action      = "allow"
    source      = "10.240.20.0/24"
    destination = "10.240.10.0/24"
    direction   = "inbound"
  }
  rules {
    name        = "subnet1-in-3"
    action      = "allow"
    source      = "10.240.30.0/24"
    destination = "10.240.10.0/24"
    direction   = "inbound"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
  rules {
    name        = "subnet1-in-4"
    action      = "allow"
    source      = "10.240.30.0/24"
    destination = "10.240.10.0/24"
    direction   = "inbound"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  rules {
    name        = "subnet2-out-1"
    action      = "allow"
    source      = "10.240.20.0/24"
    destination = "0.0.0.0/0"
    direction   = "outbound"
  }
  rules {
    name        = "subnet2-out-2"
    action      = "allow"
    source      = "10.240.20.0/24"
    destination = "10.240.10.0/24"
    direction   = "outbound"
  }
  rules {
    name        = "subnet2-in-1"
    action      = "allow"
    source      = "0.0.0.0/0"
    destination = "10.240.20.0/24"
    direction   = "inbound"
  }
  rules {
    name        = "subnet2-in-2"
    action      = "allow"
    source      = "10.240.10.0/24"
    destination = "10.240.20.0/24"
    direction   = "inbound"
  }
  rules {
    name        = "subnet3-out-1"
    action      = "allow"
    source      = "10.240.30.0/24"
    destination = "10.240.10.0/24"
    direction   = "outbound"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  rules {
    name        = "subnet3-out-2"
    action      = "allow"
    source      = "10.240.30.0/24"
    destination = "10.240.10.0/24"
    direction   = "outbound"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
  rules {
    name        = "subnet3-in-1"
    action      = "allow"
    source      = "10.240.10.0/24"
    destination = "10.240.30.0/24"
    direction   = "inbound"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  rules {
    name        = "subnet3-in-2"
    action      = "allow"
    source      = "10.240.10.0/24"
    destination = "10.240.30.0/24"
    direction   = "inbound"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
}
