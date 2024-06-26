# test-vpc1/subnet10 [10.240.64.0/24]
resource "ibm_is_network_acl" "acl-test-vpc1--subnet10" {
  name           = "acl-test-vpc1--subnet10"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.name_test-vpc1_id
  # Internal. required-connections[3]: (subnet test-vpc1/subnet10)->(subnet test-vpc1/subnet11); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "10.240.80.0/24"
    icmp {
      type = 0
    }
  }
  # Internal. response to required-connections[3]: (subnet test-vpc1/subnet10)->(subnet test-vpc1/subnet11); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.80.0/24"
    destination = "10.240.64.0/24"
    icmp {
      type = 8
    }
  }
}

# test-vpc1/subnet11 [10.240.80.0/24]
resource "ibm_is_network_acl" "acl-test-vpc1--subnet11" {
  name           = "acl-test-vpc1--subnet11"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.name_test-vpc1_id
  # Internal. required-connections[3]: (subnet test-vpc1/subnet10)->(subnet test-vpc1/subnet11); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/24"
    destination = "10.240.80.0/24"
    icmp {
      type = 0
    }
  }
  # Internal. response to required-connections[3]: (subnet test-vpc1/subnet10)->(subnet test-vpc1/subnet11); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.80.0/24"
    destination = "10.240.64.0/24"
    icmp {
      type = 8
    }
  }
}
