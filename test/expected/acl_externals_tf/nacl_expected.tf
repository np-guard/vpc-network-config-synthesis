# Attached subnets: test-vpc1/subnet1
resource "ibm_is_network_acl" "acl-test-vpc1--subnet1" {
  name           = "acl-test-vpc1--subnet1"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc1_id
  # Deny other internal communication; see rfc1918#3; item 0,0
  rules {
    name        = "rule0"
    action      = "deny"
    direction   = "outbound"
    source      = "10.0.0.0/8"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 0,0
  rules {
    name        = "rule1"
    action      = "deny"
    direction   = "inbound"
    source      = "10.0.0.0/8"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 0,1
  rules {
    name        = "rule2"
    action      = "deny"
    direction   = "outbound"
    source      = "10.0.0.0/8"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 0,1
  rules {
    name        = "rule3"
    action      = "deny"
    direction   = "inbound"
    source      = "172.16.0.0/12"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 0,2
  rules {
    name        = "rule4"
    action      = "deny"
    direction   = "outbound"
    source      = "10.0.0.0/8"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 0,2
  rules {
    name        = "rule5"
    action      = "deny"
    direction   = "inbound"
    source      = "192.168.0.0/16"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 1,0
  rules {
    name        = "rule6"
    action      = "deny"
    direction   = "outbound"
    source      = "172.16.0.0/12"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 1,0
  rules {
    name        = "rule7"
    action      = "deny"
    direction   = "inbound"
    source      = "10.0.0.0/8"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 1,1
  rules {
    name        = "rule8"
    action      = "deny"
    direction   = "outbound"
    source      = "172.16.0.0/12"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 1,1
  rules {
    name        = "rule9"
    action      = "deny"
    direction   = "inbound"
    source      = "172.16.0.0/12"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 1,2
  rules {
    name        = "rule10"
    action      = "deny"
    direction   = "outbound"
    source      = "172.16.0.0/12"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 1,2
  rules {
    name        = "rule11"
    action      = "deny"
    direction   = "inbound"
    source      = "192.168.0.0/16"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 2,0
  rules {
    name        = "rule12"
    action      = "deny"
    direction   = "outbound"
    source      = "192.168.0.0/16"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 2,0
  rules {
    name        = "rule13"
    action      = "deny"
    direction   = "inbound"
    source      = "10.0.0.0/8"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 2,1
  rules {
    name        = "rule14"
    action      = "deny"
    direction   = "outbound"
    source      = "192.168.0.0/16"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 2,1
  rules {
    name        = "rule15"
    action      = "deny"
    direction   = "inbound"
    source      = "172.16.0.0/12"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 2,2
  rules {
    name        = "rule16"
    action      = "deny"
    direction   = "outbound"
    source      = "192.168.0.0/16"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 2,2
  rules {
    name        = "rule17"
    action      = "deny"
    direction   = "inbound"
    source      = "192.168.0.0/16"
    destination = "192.168.0.0/16"
  }
  # External. required-connections[0]: (external dns)->(subnet test-vpc1/subnet1); allowed-protocols[0]
  rules {
    name        = "rule18"
    action      = "allow"
    direction   = "inbound"
    source      = "8.8.8.8"
    destination = "10.240.10.0/24"
  }
  # External. response to required-connections[0]: (external dns)->(subnet test-vpc1/subnet1); allowed-protocols[0]
  rules {
    name        = "rule19"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.10.0/24"
    destination = "8.8.8.8"
  }
}

# Attached subnets: test-vpc1/subnet2
resource "ibm_is_network_acl" "acl-test-vpc1--subnet2" {
  name           = "acl-test-vpc1--subnet2"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc1_id
  # Deny other internal communication; see rfc1918#3; item 0,0
  rules {
    name        = "rule0"
    action      = "deny"
    direction   = "outbound"
    source      = "10.0.0.0/8"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 0,0
  rules {
    name        = "rule1"
    action      = "deny"
    direction   = "inbound"
    source      = "10.0.0.0/8"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 0,1
  rules {
    name        = "rule2"
    action      = "deny"
    direction   = "outbound"
    source      = "10.0.0.0/8"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 0,1
  rules {
    name        = "rule3"
    action      = "deny"
    direction   = "inbound"
    source      = "172.16.0.0/12"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 0,2
  rules {
    name        = "rule4"
    action      = "deny"
    direction   = "outbound"
    source      = "10.0.0.0/8"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 0,2
  rules {
    name        = "rule5"
    action      = "deny"
    direction   = "inbound"
    source      = "192.168.0.0/16"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 1,0
  rules {
    name        = "rule6"
    action      = "deny"
    direction   = "outbound"
    source      = "172.16.0.0/12"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 1,0
  rules {
    name        = "rule7"
    action      = "deny"
    direction   = "inbound"
    source      = "10.0.0.0/8"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 1,1
  rules {
    name        = "rule8"
    action      = "deny"
    direction   = "outbound"
    source      = "172.16.0.0/12"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 1,1
  rules {
    name        = "rule9"
    action      = "deny"
    direction   = "inbound"
    source      = "172.16.0.0/12"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 1,2
  rules {
    name        = "rule10"
    action      = "deny"
    direction   = "outbound"
    source      = "172.16.0.0/12"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 1,2
  rules {
    name        = "rule11"
    action      = "deny"
    direction   = "inbound"
    source      = "192.168.0.0/16"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 2,0
  rules {
    name        = "rule12"
    action      = "deny"
    direction   = "outbound"
    source      = "192.168.0.0/16"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 2,0
  rules {
    name        = "rule13"
    action      = "deny"
    direction   = "inbound"
    source      = "10.0.0.0/8"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 2,1
  rules {
    name        = "rule14"
    action      = "deny"
    direction   = "outbound"
    source      = "192.168.0.0/16"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 2,1
  rules {
    name        = "rule15"
    action      = "deny"
    direction   = "inbound"
    source      = "172.16.0.0/12"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 2,2
  rules {
    name        = "rule16"
    action      = "deny"
    direction   = "outbound"
    source      = "192.168.0.0/16"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 2,2
  rules {
    name        = "rule17"
    action      = "deny"
    direction   = "inbound"
    source      = "192.168.0.0/16"
    destination = "192.168.0.0/16"
  }
  # External. required-connections[1]: (subnet test-vpc1/subnet2)->(external public internet); allowed-protocols[0]
  rules {
    name        = "rule18"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.20.0/24"
    destination = "0.0.0.0/0"
  }
  # External. response to required-connections[1]: (subnet test-vpc1/subnet2)->(external public internet); allowed-protocols[0]
  rules {
    name        = "rule19"
    action      = "allow"
    direction   = "inbound"
    source      = "0.0.0.0/0"
    destination = "10.240.20.0/24"
  }
}

# Attached subnets: test-vpc1/subnet3
resource "ibm_is_network_acl" "acl-test-vpc1--subnet3" {
  name           = "acl-test-vpc1--subnet3"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_test-vpc1_id
  # Deny all communication; subnet test-vpc1/subnet3[10.240.30.0/24] does not have required connections
  rules {
    name        = "rule0"
    action      = "deny"
    direction   = "inbound"
    source      = "0.0.0.0/0"
    destination = "10.240.30.0/24"
  }
  # Deny all communication; subnet test-vpc1/subnet3[10.240.30.0/24] does not have required connections
  rules {
    name        = "rule1"
    action      = "deny"
    direction   = "outbound"
    source      = "10.240.30.0/24"
    destination = "0.0.0.0/0"
  }
}
