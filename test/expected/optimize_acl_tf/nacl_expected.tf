# Attached subnets: subnet1
resource "ibm_is_network_acl" "acl1" {
  name           = "acl1"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc1_id
  

  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "172.217.22.46"
    destination = "10.240.10.0/24"
  }
  

  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.20.0/24"
    destination = "10.240.10.0/24"
  }
  

  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.30.0/24"
    destination = "10.240.10.0/24"
    tcp {
    }
  }
  

  rules {
    name        = "rule3"
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
    name        = "rule4"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.10.0/24"
    destination = "172.217.22.46"
  }
  

  rules {
    name        = "rule5"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.10.0/24"
    destination = "10.240.20.0/24"
  }
  

  rules {
    name        = "rule6"
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
    name        = "rule7"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.10.0/24"
    destination = "10.240.30.0/24"
    tcp {
    }
  }
}

# Attached subnets: subnet2
resource "ibm_is_network_acl" "acl2" {
  name           = "acl2"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc1_id
  

  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "0.0.0.0/0"
    destination = "10.240.20.0/24"
  }
  

  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.10.0/24"
    destination = "10.240.20.0/24"
  }
  

  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.20.0/24"
    destination = "0.0.0.0/0"
  }
  

  rules {
    name        = "rule3"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.20.0/24"
    destination = "10.240.10.0/24"
  }
}

# Attached subnets: subnet3
resource "ibm_is_network_acl" "acl3" {
  name           = "acl3"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc1_id
  

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
  

  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.10.0/24"
    destination = "10.240.30.0/24"
    tcp {
    }
  }
  

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
  

  rules {
    name        = "rule3"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.30.0/24"
    destination = "10.240.10.0/24"
    tcp {
    }
  }
}

# No attached subnets
resource "ibm_is_network_acl" "capitol-siren-chirpy-doornail" {
  name           = "capitol-siren-chirpy-doornail"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc1_id
  

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
