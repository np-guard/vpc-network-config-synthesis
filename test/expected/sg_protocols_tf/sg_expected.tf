### SG attached to test-vpc0/vsi0-subnet0
resource "ibm_is_security_group" "test-vpc0--vsi0-subnet0" {
  name           = "sg-test-vpc0--vsi0-subnet0"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0--vsi0-subnet0_id
}
# Internal. required-connections[0]: (instance test-vpc0/vsi0-subnet0)->(instance test-vpc0/vsi0-subnet1); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc0--vsi0-subnet0-0" {
  group     = ibm_is_security_group.test-vpc0--vsi0-subnet0.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc0--vsi0-subnet1.id
  udp {
  }
}
# Internal. required-connections[0]: (instance test-vpc0/vsi0-subnet0)->(instance test-vpc0/vsi0-subnet1); allowed-protocols[1]
resource "ibm_is_security_group_rule" "test-vpc0--vsi0-subnet0-1" {
  group     = ibm_is_security_group.test-vpc0--vsi0-subnet0.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc0--vsi0-subnet1.id
  icmp {
  }
}

### SG attached to test-vpc0/vsi0-subnet1
resource "ibm_is_security_group" "test-vpc0--vsi0-subnet1" {
  name           = "sg-test-vpc0--vsi0-subnet1"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0--vsi0-subnet1_id
}
# Internal. required-connections[0]: (instance test-vpc0/vsi0-subnet0)->(instance test-vpc0/vsi0-subnet1); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc0--vsi0-subnet1-0" {
  group     = ibm_is_security_group.test-vpc0--vsi0-subnet1.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc0--vsi0-subnet0.id
  udp {
  }
}
# Internal. required-connections[0]: (instance test-vpc0/vsi0-subnet0)->(instance test-vpc0/vsi0-subnet1); allowed-protocols[1]
resource "ibm_is_security_group_rule" "test-vpc0--vsi0-subnet1-1" {
  group     = ibm_is_security_group.test-vpc0--vsi0-subnet1.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc0--vsi0-subnet0.id
  icmp {
  }
}

### SG attached to test-vpc0/vsi0-subnet2
resource "ibm_is_security_group" "test-vpc0--vsi0-subnet2" {
  name           = "sg-test-vpc0--vsi0-subnet2"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0--vsi0-subnet2_id
}
# Internal. required-connections[2]: (nif test-vpc0/vsi0-subnet2/graveyard-handmade-ransack-acquaint)->(nif test-vpc0/vsi0-subnet3/icky-balsamic-outgoing-leached); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc0--vsi0-subnet2-0" {
  group     = ibm_is_security_group.test-vpc0--vsi0-subnet2.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc0--vsi0-subnet3.id
  tcp {
  }
}
# Internal. required-connections[2]: (nif test-vpc0/vsi0-subnet2/graveyard-handmade-ransack-acquaint)->(nif test-vpc0/vsi0-subnet3/icky-balsamic-outgoing-leached); allowed-protocols[1]
resource "ibm_is_security_group_rule" "test-vpc0--vsi0-subnet2-1" {
  group     = ibm_is_security_group.test-vpc0--vsi0-subnet2.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc0--vsi0-subnet3.id
  icmp {
    type = 11
    code = 1
  }
}

### SG attached to test-vpc0/vsi0-subnet3
resource "ibm_is_security_group" "test-vpc0--vsi0-subnet3" {
  name           = "sg-test-vpc0--vsi0-subnet3"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0--vsi0-subnet3_id
}
# Internal. required-connections[2]: (nif test-vpc0/vsi0-subnet2/graveyard-handmade-ransack-acquaint)->(nif test-vpc0/vsi0-subnet3/icky-balsamic-outgoing-leached); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc0--vsi0-subnet3-0" {
  group     = ibm_is_security_group.test-vpc0--vsi0-subnet3.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc0--vsi0-subnet2.id
  tcp {
  }
}
# Internal. required-connections[2]: (nif test-vpc0/vsi0-subnet2/graveyard-handmade-ransack-acquaint)->(nif test-vpc0/vsi0-subnet3/icky-balsamic-outgoing-leached); allowed-protocols[1]
resource "ibm_is_security_group_rule" "test-vpc0--vsi0-subnet3-1" {
  group     = ibm_is_security_group.test-vpc0--vsi0-subnet3.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc0--vsi0-subnet2.id
  icmp {
    type = 11
    code = 1
  }
}

### SG attached to test-vpc0/vsi0-subnet4
resource "ibm_is_security_group" "test-vpc0--vsi0-subnet4" {
  name           = "sg-test-vpc0--vsi0-subnet4"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0--vsi0-subnet4_id
}

### SG attached to test-vpc0/vsi0-subnet5
resource "ibm_is_security_group" "test-vpc0--vsi0-subnet5" {
  name           = "sg-test-vpc0--vsi0-subnet5"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0--vsi0-subnet5_id
}

### SG attached to test-vpc0/vsi1-subnet0
resource "ibm_is_security_group" "test-vpc0--vsi1-subnet0" {
  name           = "sg-test-vpc0--vsi1-subnet0"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0--vsi1-subnet0_id
}
# Internal. required-connections[1]: (instance test-vpc0/vsi1-subnet0)->(instance test-vpc0/vsi1-subnet1); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc0--vsi1-subnet0-0" {
  group     = ibm_is_security_group.test-vpc0--vsi1-subnet0.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc0--vsi1-subnet1.id
  tcp {
    port_min = 8080
    port_max = 8080
  }
}
# Internal. required-connections[1]: (instance test-vpc0/vsi1-subnet0)->(instance test-vpc0/vsi1-subnet1); allowed-protocols[1]
resource "ibm_is_security_group_rule" "test-vpc0--vsi1-subnet0-1" {
  group     = ibm_is_security_group.test-vpc0--vsi1-subnet0.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc0--vsi1-subnet1.id
  udp {
    port_min = 53
    port_max = 53
  }
}
# Internal. required-connections[1]: (instance test-vpc0/vsi1-subnet0)->(instance test-vpc0/vsi1-subnet1); allowed-protocols[2]
resource "ibm_is_security_group_rule" "test-vpc0--vsi1-subnet0-2" {
  group     = ibm_is_security_group.test-vpc0--vsi1-subnet0.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc0--vsi1-subnet1.id
  icmp {
    type = 8
  }
}

### SG attached to test-vpc0/vsi1-subnet1
resource "ibm_is_security_group" "test-vpc0--vsi1-subnet1" {
  name           = "sg-test-vpc0--vsi1-subnet1"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0--vsi1-subnet1_id
}
# Internal. required-connections[1]: (instance test-vpc0/vsi1-subnet0)->(instance test-vpc0/vsi1-subnet1); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc0--vsi1-subnet1-0" {
  group     = ibm_is_security_group.test-vpc0--vsi1-subnet1.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc0--vsi1-subnet0.id
  tcp {
    port_min = 8080
    port_max = 8080
  }
}
# Internal. required-connections[1]: (instance test-vpc0/vsi1-subnet0)->(instance test-vpc0/vsi1-subnet1); allowed-protocols[1]
resource "ibm_is_security_group_rule" "test-vpc0--vsi1-subnet1-1" {
  group     = ibm_is_security_group.test-vpc0--vsi1-subnet1.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc0--vsi1-subnet0.id
  udp {
    port_min = 53
    port_max = 53
  }
}
# Internal. required-connections[1]: (instance test-vpc0/vsi1-subnet0)->(instance test-vpc0/vsi1-subnet1); allowed-protocols[2]
resource "ibm_is_security_group_rule" "test-vpc0--vsi1-subnet1-2" {
  group     = ibm_is_security_group.test-vpc0--vsi1-subnet1.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc0--vsi1-subnet0.id
  icmp {
    type = 8
  }
}

### SG attached to test-vpc0/vsi1-subnet2
resource "ibm_is_security_group" "test-vpc0--vsi1-subnet2" {
  name           = "sg-test-vpc0--vsi1-subnet2"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0--vsi1-subnet2_id
}

### SG attached to test-vpc0/vsi1-subnet3
resource "ibm_is_security_group" "test-vpc0--vsi1-subnet3" {
  name           = "sg-test-vpc0--vsi1-subnet3"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0--vsi1-subnet3_id
}

### SG attached to test-vpc0/vsi1-subnet4
resource "ibm_is_security_group" "test-vpc0--vsi1-subnet4" {
  name           = "sg-test-vpc0--vsi1-subnet4"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0--vsi1-subnet4_id
}

### SG attached to test-vpc0/vsi1-subnet5
resource "ibm_is_security_group" "test-vpc0--vsi1-subnet5" {
  name           = "sg-test-vpc0--vsi1-subnet5"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0--vsi1-subnet5_id
}

### SG attached to test-vpc1/vsi0-subnet10
resource "ibm_is_security_group" "test-vpc1--vsi0-subnet10" {
  name           = "sg-test-vpc1--vsi0-subnet10"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc1--vsi0-subnet10_id
}
# External. required-connections[3]: (instance test-vpc1/vsi0-subnet10)->(external dns); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc1--vsi0-subnet10-0" {
  group     = ibm_is_security_group.test-vpc1--vsi0-subnet10.id
  direction = "outbound"
  remote    = "8.8.8.8"
  tcp {
  }
}

### SG attached to test-vpc1/vsi0-subnet11
resource "ibm_is_security_group" "test-vpc1--vsi0-subnet11" {
  name           = "sg-test-vpc1--vsi0-subnet11"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc1--vsi0-subnet11_id
}

### SG attached to test-vpc2/vsi0-subnet20
resource "ibm_is_security_group" "test-vpc2--vsi0-subnet20" {
  name           = "sg-test-vpc2--vsi0-subnet20"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc2--vsi0-subnet20_id
}
# Internal. required-connections[4]: (instance test-vpc2/vsi0-subnet20)->(instance test-vpc2/vsi2-subnet20); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc2--vsi0-subnet20-0" {
  group     = ibm_is_security_group.test-vpc2--vsi0-subnet20.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc2--vsi2-subnet20.id
}

### SG attached to test-vpc2/vsi1-subnet20
resource "ibm_is_security_group" "test-vpc2--vsi1-subnet20" {
  name           = "sg-test-vpc2--vsi1-subnet20"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc2--vsi1-subnet20_id
}
# External. required-connections[5]: (instance test-vpc2/vsi1-subnet20)->(external public internet); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc2--vsi1-subnet20-0" {
  group     = ibm_is_security_group.test-vpc2--vsi1-subnet20.id
  direction = "outbound"
  remote    = "0.0.0.0/0"
}

### SG attached to test-vpc2/vsi2-subnet20
resource "ibm_is_security_group" "test-vpc2--vsi2-subnet20" {
  name           = "sg-test-vpc2--vsi2-subnet20"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc2--vsi2-subnet20_id
}
# Internal. required-connections[4]: (instance test-vpc2/vsi0-subnet20)->(instance test-vpc2/vsi2-subnet20); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc2--vsi2-subnet20-0" {
  group     = ibm_is_security_group.test-vpc2--vsi2-subnet20.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc2--vsi0-subnet20.id
}

### SG attached to test-vpc3/vsi0-subnet30
resource "ibm_is_security_group" "test-vpc3--vsi0-subnet30" {
  name           = "sg-test-vpc3--vsi0-subnet30"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc3--vsi0-subnet30_id
}
