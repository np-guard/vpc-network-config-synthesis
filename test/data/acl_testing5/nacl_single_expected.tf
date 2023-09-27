resource "ibm_is_network_acl" "acl-1" {
  name           = "acl-1"
  resource_group = local.acl_synth_resource_group_id
  vpc            = local.acl_synth_vpc_id
  # Internal. required-connections[0]: (segment need-dns)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule0"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "10.240.64.0/24"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  # Internal. response to required-connections[0]: (segment need-dns)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule1"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/24"
    destination = "10.240.1.0/24"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
  # Internal. required-connections[0]: (segment need-dns)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule2"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.1.0/24"
    destination = "10.240.64.0/24"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  # Internal. response to required-connections[0]: (segment need-dns)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule3"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "10.240.1.0/24"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
  # Internal. required-connections[0]: (segment need-dns)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule4"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "10.240.1.0/24"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  # Internal. response to required-connections[0]: (segment need-dns)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule5"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.1.0/24"
    destination = "10.240.64.0/24"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
  # Internal. required-connections[0]: (segment need-dns)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule6"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/24"
    destination = "10.240.1.0/24"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  # Internal. response to required-connections[0]: (segment need-dns)->(segment need-dns); allowed-protocols[0]
  rules {
    name        = "rule7"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "10.240.64.0/24"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
  # Internal. required-connections[2]: (segment need-dns)->(subnet sub3-1-ky); allowed-protocols[0]
  rules {
    name        = "rule8"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "10.240.128.0/24"
    icmp {
      type = 8
      code = 0
    }
  }
  # Internal. response to required-connections[2]: (segment need-dns)->(subnet sub3-1-ky); allowed-protocols[0]
  rules {
    name        = "rule9"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.128.0/24"
    destination = "10.240.1.0/24"
    icmp {
      type = 0
      code = 0
    }
  }
  # Internal. required-connections[2]: (segment need-dns)->(subnet sub3-1-ky); allowed-protocols[0]
  rules {
    name        = "rule10"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.1.0/24"
    destination = "10.240.128.0/24"
    icmp {
      type = 8
      code = 0
    }
  }
  # Internal. response to required-connections[2]: (segment need-dns)->(subnet sub3-1-ky); allowed-protocols[0]
  rules {
    name        = "rule11"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.128.0/24"
    destination = "10.240.1.0/24"
    icmp {
      type = 0
      code = 0
    }
  }
  # Internal. required-connections[2]: (segment need-dns)->(subnet sub3-1-ky); allowed-protocols[0]
  rules {
    name        = "rule12"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "10.240.128.0/24"
    icmp {
      type = 8
      code = 0
    }
  }
  # Internal. response to required-connections[2]: (segment need-dns)->(subnet sub3-1-ky); allowed-protocols[0]
  rules {
    name        = "rule13"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.128.0/24"
    destination = "10.240.64.0/24"
    icmp {
      type = 0
      code = 0
    }
  }
  # Internal. required-connections[2]: (segment need-dns)->(subnet sub3-1-ky); allowed-protocols[0]
  rules {
    name        = "rule14"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/24"
    destination = "10.240.128.0/24"
    icmp {
      type = 8
      code = 0
    }
  }
  # Internal. response to required-connections[2]: (segment need-dns)->(subnet sub3-1-ky); allowed-protocols[0]
  rules {
    name        = "rule15"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.128.0/24"
    destination = "10.240.64.0/24"
    icmp {
      type = 0
      code = 0
    }
  }
  # Internal. required-connections[3]: (subnet sub1-1-ky)->(subnet sub1-2-ky); allowed-protocols[0]
  rules {
    name        = "rule16"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "10.240.2.0/24"
    tcp {
    }
  }
  # Internal. response to required-connections[3]: (subnet sub1-1-ky)->(subnet sub1-2-ky); allowed-protocols[0]
  rules {
    name        = "rule17"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.2.0/24"
    destination = "10.240.1.0/24"
    tcp {
    }
  }
  # Internal. required-connections[3]: (subnet sub1-1-ky)->(subnet sub1-2-ky); allowed-protocols[0]
  rules {
    name        = "rule18"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.1.0/24"
    destination = "10.240.2.0/24"
    tcp {
    }
  }
  # Internal. response to required-connections[3]: (subnet sub1-1-ky)->(subnet sub1-2-ky); allowed-protocols[0]
  rules {
    name        = "rule19"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.2.0/24"
    destination = "10.240.1.0/24"
    tcp {
    }
  }
  # Internal. required-connections[4]: (subnet sub1-1-ky)->(subnet sub1-3-ky); allowed-protocols[0]
  rules {
    name        = "rule20"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "10.240.3.0/24"
    tcp {
    }
  }
  # Internal. response to required-connections[4]: (subnet sub1-1-ky)->(subnet sub1-3-ky); allowed-protocols[0]
  rules {
    name        = "rule21"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.3.0/24"
    destination = "10.240.1.0/24"
    tcp {
    }
  }
  # Internal. required-connections[4]: (subnet sub1-1-ky)->(subnet sub1-3-ky); allowed-protocols[0]
  rules {
    name        = "rule22"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.1.0/24"
    destination = "10.240.3.0/24"
    tcp {
    }
  }
  # Internal. response to required-connections[4]: (subnet sub1-1-ky)->(subnet sub1-3-ky); allowed-protocols[0]
  rules {
    name        = "rule23"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.3.0/24"
    destination = "10.240.1.0/24"
    tcp {
    }
  }
  # Internal. required-connections[5]: (subnet sub1-2-ky)->(subnet sub1-3-ky); allowed-protocols[0]
  rules {
    name        = "rule24"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.2.0/24"
    destination = "10.240.3.0/24"
    tcp {
    }
  }
  # Internal. response to required-connections[5]: (subnet sub1-2-ky)->(subnet sub1-3-ky); allowed-protocols[0]
  rules {
    name        = "rule25"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.3.0/24"
    destination = "10.240.2.0/24"
    tcp {
    }
  }
  # Internal. required-connections[5]: (subnet sub1-2-ky)->(subnet sub1-3-ky); allowed-protocols[0]
  rules {
    name        = "rule26"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.2.0/24"
    destination = "10.240.3.0/24"
    tcp {
    }
  }
  # Internal. response to required-connections[5]: (subnet sub1-2-ky)->(subnet sub1-3-ky); allowed-protocols[0]
  rules {
    name        = "rule27"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.3.0/24"
    destination = "10.240.2.0/24"
    tcp {
    }
  }
  # Internal. required-connections[6]: (subnet sub2-1-ky)->(subnet sub2-2-ky); allowed-protocols[0]
  rules {
    name        = "rule28"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "10.240.65.0/24"
  }
  # Internal. response to required-connections[6]: (subnet sub2-1-ky)->(subnet sub2-2-ky); allowed-protocols[0]
  rules {
    name        = "rule29"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.65.0/24"
    destination = "10.240.64.0/24"
  }
  # Internal. required-connections[6]: (subnet sub2-1-ky)->(subnet sub2-2-ky); allowed-protocols[0]
  rules {
    name        = "rule30"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/24"
    destination = "10.240.65.0/24"
  }
  # Internal. response to required-connections[6]: (subnet sub2-1-ky)->(subnet sub2-2-ky); allowed-protocols[0]
  rules {
    name        = "rule31"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.65.0/24"
    destination = "10.240.64.0/24"
  }
  # Internal. required-connections[7]: (subnet sub3-1-ky)->(subnet sub2-1-ky); allowed-protocols[0]
  rules {
    name        = "rule32"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.128.0/24"
    destination = "10.240.64.0/24"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  # Internal. response to required-connections[7]: (subnet sub3-1-ky)->(subnet sub2-1-ky); allowed-protocols[0]
  rules {
    name        = "rule33"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.64.0/24"
    destination = "10.240.128.0/24"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
  # Internal. required-connections[7]: (subnet sub3-1-ky)->(subnet sub2-1-ky); allowed-protocols[0]
  rules {
    name        = "rule34"
    action      = "allow"
    direction   = "inbound"
    source      = "10.240.128.0/24"
    destination = "10.240.64.0/24"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  # Internal. response to required-connections[7]: (subnet sub3-1-ky)->(subnet sub2-1-ky); allowed-protocols[0]
  rules {
    name        = "rule35"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "10.240.128.0/24"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
  # Deny other internal communication; see rfc1918#3; item 0,0
  rules {
    name        = "rule36"
    action      = "deny"
    direction   = "outbound"
    source      = "10.0.0.0/8"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 0,0
  rules {
    name        = "rule37"
    action      = "deny"
    direction   = "inbound"
    source      = "10.0.0.0/8"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 0,1
  rules {
    name        = "rule38"
    action      = "deny"
    direction   = "outbound"
    source      = "10.0.0.0/8"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 0,1
  rules {
    name        = "rule39"
    action      = "deny"
    direction   = "inbound"
    source      = "172.16.0.0/12"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 0,2
  rules {
    name        = "rule40"
    action      = "deny"
    direction   = "outbound"
    source      = "10.0.0.0/8"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 0,2
  rules {
    name        = "rule41"
    action      = "deny"
    direction   = "inbound"
    source      = "192.168.0.0/16"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 1,0
  rules {
    name        = "rule42"
    action      = "deny"
    direction   = "outbound"
    source      = "172.16.0.0/12"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 1,0
  rules {
    name        = "rule43"
    action      = "deny"
    direction   = "inbound"
    source      = "10.0.0.0/8"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 1,1
  rules {
    name        = "rule44"
    action      = "deny"
    direction   = "outbound"
    source      = "172.16.0.0/12"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 1,1
  rules {
    name        = "rule45"
    action      = "deny"
    direction   = "inbound"
    source      = "172.16.0.0/12"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 1,2
  rules {
    name        = "rule46"
    action      = "deny"
    direction   = "outbound"
    source      = "172.16.0.0/12"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 1,2
  rules {
    name        = "rule47"
    action      = "deny"
    direction   = "inbound"
    source      = "192.168.0.0/16"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 2,0
  rules {
    name        = "rule48"
    action      = "deny"
    direction   = "outbound"
    source      = "192.168.0.0/16"
    destination = "10.0.0.0/8"
  }
  # Deny other internal communication; see rfc1918#3; item 2,0
  rules {
    name        = "rule49"
    action      = "deny"
    direction   = "inbound"
    source      = "10.0.0.0/8"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 2,1
  rules {
    name        = "rule50"
    action      = "deny"
    direction   = "outbound"
    source      = "192.168.0.0/16"
    destination = "172.16.0.0/12"
  }
  # Deny other internal communication; see rfc1918#3; item 2,1
  rules {
    name        = "rule51"
    action      = "deny"
    direction   = "inbound"
    source      = "172.16.0.0/12"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 2,2
  rules {
    name        = "rule52"
    action      = "deny"
    direction   = "outbound"
    source      = "192.168.0.0/16"
    destination = "192.168.0.0/16"
  }
  # Deny other internal communication; see rfc1918#3; item 2,2
  rules {
    name        = "rule53"
    action      = "deny"
    direction   = "inbound"
    source      = "192.168.0.0/16"
    destination = "192.168.0.0/16"
  }
  # External. required-connections[1]: (segment need-dns)->(external dns); allowed-protocols[0]
  rules {
    name        = "rule54"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.1.0/24"
    destination = "8.8.8.8"
  }
  # External. response to required-connections[1]: (segment need-dns)->(external dns); allowed-protocols[0]
  rules {
    name        = "rule55"
    action      = "allow"
    direction   = "inbound"
    source      = "8.8.8.8"
    destination = "10.240.1.0/24"
  }
  # External. required-connections[1]: (segment need-dns)->(external dns); allowed-protocols[0]
  rules {
    name        = "rule56"
    action      = "allow"
    direction   = "outbound"
    source      = "10.240.64.0/24"
    destination = "8.8.8.8"
  }
  # External. response to required-connections[1]: (segment need-dns)->(external dns); allowed-protocols[0]
  rules {
    name        = "rule57"
    action      = "allow"
    direction   = "inbound"
    source      = "8.8.8.8"
    destination = "10.240.64.0/24"
  }
}
