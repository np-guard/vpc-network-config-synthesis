# Attached subnets: sub1-1
resource "ibm_is_network_acl" "acl1-1" {
  name           = "acl1-1"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_testacl5-vpc_id
  rules {
    name        = "rule0"
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
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.2.0/23"
    destination = "10.240.1.0/24"
    tcp {
    }
  }
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.128.0/24"
    destination = "10.240.1.0/24"
    icmp {
      type = 0
      code = 0
    }
  }
  rules {
    name        = "rule3"
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
    name        = "rule4"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "10.240.2.0/23"
    tcp {
    }
  }
  rules {
    name        = "rule5"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "10.240.128.0/24"
    icmp {
      type = 0
      code = 0
    }
  }
}

# Attached subnets: sub1-2, sub1-3
resource "ibm_is_network_acl" "acl1-2" {
  name           = "acl1-2"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_testacl5-vpc_id
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.1.0/24"
    destination = "10.240.2.0/23"
    tcp {
    }
  }
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.2.0/23"
    destination = "10.240.2.0/23"
    tcp {
    }
  }
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.2.0/23"
    destination = "10.240.1.0/24"
    tcp {
    }
  }
  rules {
    name        = "rule3"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.2.0/23"
    destination = "10.240.2.0/23"
    tcp {
    }
  }
}

# Attached subnets: sub2-1
resource "ibm_is_network_acl" "acl2-1" {
  name           = "acl2-1"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_testacl5-vpc_id
  rules {
    name        = "rule0"
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
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.65.0/24"
    destination = "10.240.64.0/24"
  }
  rules {
    name        = "rule2"
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
    name        = "rule3"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.128.0/24"
    destination = "10.240.64.0/24"
    icmp {
      type = 0
      code = 0
    }
  }
  rules {
    name        = "rule4"
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
    name        = "rule5"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "10.240.65.0/24"
  }
  rules {
    name        = "rule6"
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
    name        = "rule7"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "10.240.128.0/24"
    icmp {
      type = 0
      code = 0
    }
  }
}

# Attached subnets: sub2-2
resource "ibm_is_network_acl" "acl2-2" {
  name           = "acl2-2"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_testacl5-vpc_id
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/24"
    destination = "10.240.65.0/24"
  }
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.65.0/24"
    destination = "10.240.64.0/24"
  }
}

# Attached subnets: sub3-1
resource "ibm_is_network_acl" "acl3-1" {
  name           = "acl3-1"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_testacl5-vpc_id
  rules {
    name        = "rule0"
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
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/24"
    destination = "10.240.128.0/24"
    icmp {
      type = 0
      code = 0
    }
  }
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.1.0/24"
    destination = "10.240.128.0/24"
    icmp {
      type = 0
      code = 0
    }
  }
  rules {
    name        = "rule3"
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
    name        = "rule4"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.128.0/24"
    destination = "10.240.64.0/24"
    icmp {
      type = 0
      code = 0
    }
  }
  rules {
    name        = "rule5"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.128.0/24"
    destination = "10.240.1.0/24"
    icmp {
      type = 0
      code = 0
    }
  }
}

# No attached subnets
resource "ibm_is_network_acl" "disallow-laborious-compress-abiding" {
  name           = "disallow-laborious-compress-abiding"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_testacl5-vpc_id
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "0.0.0.0/0"
    destination = "0.0.0.0/0"
  }
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "outbound"
    source      = "0.0.0.0/0"
    destination = "0.0.0.0/0"
  }
}
