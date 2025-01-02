# Attached subnets: sub1-1
resource "ibm_is_network_acl" "testacl5-vpc--sub1-1" {
  name           = "testacl5-vpc--sub1-1"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_testacl5-vpc_id
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "1.1.1.0/31"
    tcp {
    }
  }
}

# Attached subnets: sub1-2
resource "ibm_is_network_acl" "testacl5-vpc--sub1-2" {
  name           = "testacl5-vpc--sub1-2"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_testacl5-vpc_id
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.2.0/24"
    destination = "2.2.2.2"
    udp {
      port_max = 20
    }
  }
}

# Attached subnets: sub1-3
resource "ibm_is_network_acl" "testacl5-vpc--sub1-3" {
  name           = "testacl5-vpc--sub1-3"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_testacl5-vpc_id
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.3.0/24"
    destination = "10.240.64.0/24"
  }
}

# Attached subnets: sub2-1
resource "ibm_is_network_acl" "testacl5-vpc--sub2-1" {
  name           = "testacl5-vpc--sub2-1"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_testacl5-vpc_id
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.3.0/24"
    destination = "10.240.64.0/24"
  }
}

# Attached subnets: sub2-2
resource "ibm_is_network_acl" "testacl5-vpc--sub2-2" {
  name           = "testacl5-vpc--sub2-2"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_testacl5-vpc_id
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.65.0/24"
    destination = "10.240.128.0/24"
    tcp {
      source_port_min = 11
      source_port_max = 20
    }
  }
}

# Attached subnets: sub3-1
resource "ibm_is_network_acl" "testacl5-vpc--sub3-1" {
  name           = "testacl5-vpc--sub3-1"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_testacl5-vpc_id
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.65.0/24"
    destination = "10.240.128.0/24"
    tcp {
      source_port_min = 11
      source_port_max = 20
    }
  }
}
