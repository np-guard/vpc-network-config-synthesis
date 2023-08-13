resource "ibm_is_network_acl" "acl-for-test-vpc-ky" {
  name = "acl-for-test-vpc-ky"
  resource_group = var.resource_group_id
  vpc  = var.vpc_id

  rules {
    name        = "sub1-1-out-1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "8.8.8.8"
    udp {
      port_min = 53
      port_max = 53
    }
  }
  rules {
    name        = "sub1-1-out-2"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "10.240.2.0/23"
    tcp {}
  }
  rules {
    name        = "sub1-1-out-3"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "10.240.128.0/24"
    icmp {
      code = 0
      type = 0
    }
  }
  rules {
    name        = "sub1-1-in-1"
    action      = "allow"
    direction   = "inbound"
    source      = "8.8.8.8"
    destination = "10.240.1.0/24"
    udp {
      source_port_min = 53
      source_port_max = 53
    }
  }
  rules {
    name        = "sub1-1-in-2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.2.0/23"
    destination = "10.240.1.0/24"
    tcp {}
  }
  rules {
    name        = "sub1-1-in-3"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.128.0/24"
    destination = "10.240.1.0/24"
    icmp {
      code = 0
      type = 0
    }
  }
  rules {
    name        = "sub1-2-out-1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.2.0/23"
    destination = "10.240.1.0/24"
    tcp {}
  }
  rules {
    name        = "sub1-2-out-2"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.2.0/23"
    destination = "10.240.2.0/23"
    tcp {}
  }
  rules {
    name        = "sub1-2-in-1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.1.0/24"
    destination = "10.240.2.0/23"
    tcp {}
  }
  rules {
    name        = "sub1-2-in-2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.2.0/23"
    destination = "10.240.2.0/23"
    tcp {}
  }
  rules {
    name        = "sub2-1-out-1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "8.8.8.8"
    udp {
      port_min = 53
      port_max = 53
    }
  }
  rules {
    name        = "sub2-1-out-2"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "10.240.65.0/24"
  }
  rules {
    name        = "sub2-1-out-3"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "10.240.128.0/24"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
  rules {
    name        = "sub2-1-out-4"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "10.240.128.0/24"
    icmp {
      code = 0
      type = 0
    }
  }
  rules {
    name        = "sub2-1-in-1"
    action      = "allow"
    direction   = "inbound"
    source      = "8.8.8.8"
    destination = "10.240.64.0/24"
    udp {
      source_port_min = 53
      source_port_max = 53
    }
  }
  rules {
    name        = "sub2-1-in-2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.65.0/24"
    destination = "10.240.64.0/24"
  }
  rules {
    name        = "sub2-1-in-3"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.128.0/24"
    destination = "10.240.64.0/24"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  rules {
    name        = "sub2-1-in-4"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.128.0/24"
    destination = "10.240.64.0/24"
    icmp {
      code = 0
      type = 0
    }
  }
  rules {
    name        = "sub2-2-out-1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.65.0/24"
    destination = "10.240.64.0/24"
  }
  rules {
    name        = "sub2-2-in-1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/24"
    destination = "10.240.65.0/24"
  }
  rules {
    name        = "sub3-1-out-1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.128.0/24"
    destination = "10.240.64.0/24"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  rules {
    name        = "sub3-1-out-2"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.128.0/24"
    destination = "10.240.64.0/24"
    icmp {
      code = 0
      type = 0
    }
  }
  rules {
    name        = "sub3-1-out-3"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.128.0/24"
    destination = "10.240.1.0/24"
    icmp {
      code = 0
      type = 0
    }
  }
  rules {
    name        = "sub3-1-in-1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/24"
    destination = "10.240.128.0/24"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
  rules {
    name        = "sub3-1-in-2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/24"
    destination = "10.240.128.0/24"
    icmp {
      code = 0
      type = 0
    }
  }
  rules {
    name        = "sub3-1-in-3"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.1.0/24"
    destination = "10.240.128.0/24"
    icmp {
      code = 0
      type = 0
    }
  }
}