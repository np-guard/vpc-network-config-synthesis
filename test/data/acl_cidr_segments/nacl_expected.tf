# test-vpc-ky/sub1-1-ky [10.240.1.0/24]
resource "ibm_is_network_acl" "acl-test-vpc-ky/sub1-1-ky" {
  name           = "acl-test-vpc-ky/sub1-1-ky"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_vpc_id
  # Internal. required-connections[0]: (segment need-dns)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "10.240.64.0/24"
  }
  # Internal. response to required-connections[0]: (segment need-dns)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/24"
    destination = "10.240.1.0/24"
  }
  # Internal. required-connections[2]: (segment cidrSegment1)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.2.0/23"
    destination = "10.240.1.0/24"
  }
  # Internal. response to required-connections[2]: (segment cidrSegment1)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule3"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "10.240.2.0/23"
  }
  # Deny other internal communication; see rfc1918#3; item 0,0
  rules {
    name        = "rule4"
    action      = "deny"
    direction   = "outbound"
    source      = "10.0.0.0/8"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 0,0
  rules {
    name        = "rule5"
    action      = "deny"
    direction   = "inbound"
    source      = "10.0.0.0/8"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 0,1
  rules {
    name        = "rule6"
    action      = "deny"
    direction   = "outbound"
    source      = "10.0.0.0/8"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 0,1
  rules {
    name        = "rule7"
    action      = "deny"
    direction   = "inbound"
    source      = "172.16.0.0/12"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 0,2
  rules {
    name        = "rule8"
    action      = "deny"
    direction   = "outbound"
    source      = "10.0.0.0/8"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 0,2
  rules {
    name        = "rule9"
    action      = "deny"
    direction   = "inbound"
    source      = "192.168.0.0/16"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 1,0
  rules {
    name        = "rule10"
    action      = "deny"
    direction   = "outbound"
    source      = "172.16.0.0/12"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 1,0
  rules {
    name        = "rule11"
    action      = "deny"
    direction   = "inbound"
    source      = "10.0.0.0/8"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 1,1
  rules {
    name        = "rule12"
    action      = "deny"
    direction   = "outbound"
    source      = "172.16.0.0/12"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 1,1
  rules {
    name        = "rule13"
    action      = "deny"
    direction   = "inbound"
    source      = "172.16.0.0/12"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 1,2
  rules {
    name        = "rule14"
    action      = "deny"
    direction   = "outbound"
    source      = "172.16.0.0/12"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 1,2
  rules {
    name        = "rule15"
    action      = "deny"
    direction   = "inbound"
    source      = "192.168.0.0/16"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 2,0
  rules {
    name        = "rule16"
    action      = "deny"
    direction   = "outbound"
    source      = "192.168.0.0/16"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 2,0
  rules {
    name        = "rule17"
    action      = "deny"
    direction   = "inbound"
    source      = "10.0.0.0/8"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 2,1
  rules {
    name        = "rule18"
    action      = "deny"
    direction   = "outbound"
    source      = "192.168.0.0/16"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 2,1
  rules {
    name        = "rule19"
    action      = "deny"
    direction   = "inbound"
    source      = "172.16.0.0/12"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 2,2
  rules {
    name        = "rule20"
    action      = "deny"
    direction   = "outbound"
    source      = "192.168.0.0/16"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 2,2
  rules {
    name        = "rule21"
    action      = "deny"
    direction   = "inbound"
    source      = "192.168.0.0/16"
    destination = "192.168.0.0/16"
  }
  # External. required-connections[1]: (segment need-dns)->(external dns); allowed-protocols[0]
  rules {
    name        = "rule22"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "8.8.8.8"
    udp {
      port_min = 53
      port_max = 53
    }
  }
}

# test-vpc-ky/sub1-2-ky [10.240.2.0/24]
resource "ibm_is_network_acl" "acl-test-vpc-ky/sub1-2-ky" {
  name           = "acl-test-vpc-ky/sub1-2-ky"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_vpc_id
  # Internal. required-connections[2]: (segment cidrSegment1)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.2.0/24"
    destination = "10.240.1.0/24"
  }
  # Internal. response to required-connections[2]: (segment cidrSegment1)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.1.0/24"
    destination = "10.240.2.0/24"
  }
  # Internal. required-connections[2]: (segment cidrSegment1)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.2.0/24"
    destination = "10.240.64.0/24"
  }
  # Internal. response to required-connections[2]: (segment cidrSegment1)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule3"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/24"
    destination = "10.240.2.0/24"
  }
}

# test-vpc-ky/sub1-3-ky [10.240.3.0/24]
resource "ibm_is_network_acl" "acl-test-vpc-ky/sub1-3-ky" {
  name           = "acl-test-vpc-ky/sub1-3-ky"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_vpc_id
  # Internal. required-connections[2]: (segment cidrSegment1)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.3.0/24"
    destination = "10.240.1.0/24"
  }
  # Internal. response to required-connections[2]: (segment cidrSegment1)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.1.0/24"
    destination = "10.240.3.0/24"
  }
  # Internal. required-connections[2]: (segment cidrSegment1)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.3.0/24"
    destination = "10.240.64.0/24"
  }
  # Internal. response to required-connections[2]: (segment cidrSegment1)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule3"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/24"
    destination = "10.240.3.0/24"
  }
}

# test-vpc-ky/sub2-1-ky [10.240.64.0/24]
resource "ibm_is_network_acl" "acl-test-vpc-ky/sub2-1-ky" {
  name           = "acl-test-vpc-ky/sub2-1-ky"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_vpc_id
  # Internal. required-connections[0]: (segment need-dns)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.1.0/24"
    destination = "10.240.64.0/24"
  }
  # Internal. response to required-connections[0]: (segment need-dns)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "10.240.1.0/24"
  }
  # Internal. required-connections[2]: (segment cidrSegment1)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.2.0/23"
    destination = "10.240.64.0/24"
  }
  # Internal. response to required-connections[2]: (segment cidrSegment1)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule3"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "10.240.2.0/23"
  }
  # Internal. required-connections[3]: (segment cidrSegment2)->(subnet test-vpc-ky/sub3-1-ky); allowed-protocols[0]
  rules {
    name        = "rule4"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "10.240.128.0/24"
    tcp {
    }
  }
  # Internal. response to required-connections[3]: (segment cidrSegment2)->(subnet test-vpc-ky/sub3-1-ky); allowed-protocols[0]
  rules {
    name        = "rule5"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.128.0/24"
    destination = "10.240.64.0/24"
    tcp {
    }
  }
  # Deny other internal communication; see rfc1918#3; item 0,0
  rules {
    name        = "rule6"
    action      = "deny"
    direction   = "outbound"
    source      = "10.0.0.0/8"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 0,0
  rules {
    name        = "rule7"
    action      = "deny"
    direction   = "inbound"
    source      = "10.0.0.0/8"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 0,1
  rules {
    name        = "rule8"
    action      = "deny"
    direction   = "outbound"
    source      = "10.0.0.0/8"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 0,1
  rules {
    name        = "rule9"
    action      = "deny"
    direction   = "inbound"
    source      = "172.16.0.0/12"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 0,2
  rules {
    name        = "rule10"
    action      = "deny"
    direction   = "outbound"
    source      = "10.0.0.0/8"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 0,2
  rules {
    name        = "rule11"
    action      = "deny"
    direction   = "inbound"
    source      = "192.168.0.0/16"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 1,0
  rules {
    name        = "rule12"
    action      = "deny"
    direction   = "outbound"
    source      = "172.16.0.0/12"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 1,0
  rules {
    name        = "rule13"
    action      = "deny"
    direction   = "inbound"
    source      = "10.0.0.0/8"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 1,1
  rules {
    name        = "rule14"
    action      = "deny"
    direction   = "outbound"
    source      = "172.16.0.0/12"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 1,1
  rules {
    name        = "rule15"
    action      = "deny"
    direction   = "inbound"
    source      = "172.16.0.0/12"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 1,2
  rules {
    name        = "rule16"
    action      = "deny"
    direction   = "outbound"
    source      = "172.16.0.0/12"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 1,2
  rules {
    name        = "rule17"
    action      = "deny"
    direction   = "inbound"
    source      = "192.168.0.0/16"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 2,0
  rules {
    name        = "rule18"
    action      = "deny"
    direction   = "outbound"
    source      = "192.168.0.0/16"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 2,0
  rules {
    name        = "rule19"
    action      = "deny"
    direction   = "inbound"
    source      = "10.0.0.0/8"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 2,1
  rules {
    name        = "rule20"
    action      = "deny"
    direction   = "outbound"
    source      = "192.168.0.0/16"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 2,1
  rules {
    name        = "rule21"
    action      = "deny"
    direction   = "inbound"
    source      = "172.16.0.0/12"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 2,2
  rules {
    name        = "rule22"
    action      = "deny"
    direction   = "outbound"
    source      = "192.168.0.0/16"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 2,2
  rules {
    name        = "rule23"
    action      = "deny"
    direction   = "inbound"
    source      = "192.168.0.0/16"
    destination = "192.168.0.0/16"
  }
  # External. required-connections[1]: (segment need-dns)->(external dns); allowed-protocols[0]
  rules {
    name        = "rule24"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "8.8.8.8"
    udp {
      port_min = 53
      port_max = 53
    }
  }
}

# test-vpc-ky/sub2-2-ky [10.240.65.0/24]
resource "ibm_is_network_acl" "acl-test-vpc-ky/sub2-2-ky" {
  name           = "acl-test-vpc-ky/sub2-2-ky"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_vpc_id
  # Internal. required-connections[3]: (segment cidrSegment2)->(subnet test-vpc-ky/sub3-1-ky); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.65.0/24"
    destination = "10.240.128.0/24"
    tcp {
    }
  }
  # Internal. response to required-connections[3]: (segment cidrSegment2)->(subnet test-vpc-ky/sub3-1-ky); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.128.0/24"
    destination = "10.240.65.0/24"
    tcp {
    }
  }
}

# test-vpc-ky/sub3-1-ky [10.240.128.0/24]
resource "ibm_is_network_acl" "acl-test-vpc-ky/sub3-1-ky" {
  name           = "acl-test-vpc-ky/sub3-1-ky"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_vpc_id
  # Internal. required-connections[3]: (segment cidrSegment2)->(subnet test-vpc-ky/sub3-1-ky); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/18"
    destination = "10.240.128.0/24"
    tcp {
    }
  }
  # Internal. response to required-connections[3]: (segment cidrSegment2)->(subnet test-vpc-ky/sub3-1-ky); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.128.0/24"
    destination = "10.240.64.0/18"
    tcp {
    }
  }
}
