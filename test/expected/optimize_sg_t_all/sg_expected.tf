### SG sg1 is not attached to anything
resource "ibm_is_security_group" "sg1" {
  name           = "sg-sg1"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc1_id
}
resource "ibm_is_security_group_rule" "sg1-0" {
  group     = ibm_is_security_group.sg1.id
  direction = "inbound"
  local     = "0.0.0.0/0"
  remote    = "0.0.0.0/0"
}
resource "ibm_is_security_group_rule" "sg1-1" {
  group     = ibm_is_security_group.sg1.id
  direction = "outbound"
  local     = "0.0.0.0/0"
  remote    = "0.0.0.0/0"
}

### SG test-vpc1--vsi1 is attached to ni1
resource "ibm_is_security_group" "test-vpc1--vsi1" {
  name           = "sg-test-vpc1--vsi1"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc1_id
}
resource "ibm_is_security_group_rule" "test-vpc1--vsi1-0" {
  group     = ibm_is_security_group.test-vpc1--vsi1.id
  direction = "outbound"
  local     = "0.0.0.0/0"
  remote    = "0.0.0.0/29"
  tcp {
    port_max = 10
  }
}
resource "ibm_is_security_group_rule" "test-vpc1--vsi1-1" {
  group     = ibm_is_security_group.test-vpc1--vsi1.id
  direction = "outbound"
  local     = "0.0.0.0/0"
  remote    = "0.0.0.2/31"
}

### SG test-vpc1--vsi2 is attached to ni2
resource "ibm_is_security_group" "test-vpc1--vsi2" {
  name           = "sg-test-vpc1--vsi2"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc1_id
}

### SG test-vpc1--vsi3a is attached to ni3a
resource "ibm_is_security_group" "test-vpc1--vsi3a" {
  name           = "sg-test-vpc1--vsi3a"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc1_id
}

### SG test-vpc1--vsi3b is attached to ni3b
resource "ibm_is_security_group" "test-vpc1--vsi3b" {
  name           = "sg-test-vpc1--vsi3b"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc1_id
}

### SG wombat-hesitate-scorn-subprime is not attached to anything
resource "ibm_is_security_group" "wombat-hesitate-scorn-subprime" {
  name           = "sg-wombat-hesitate-scorn-subprime"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc1_id
}
resource "ibm_is_security_group_rule" "wombat-hesitate-scorn-subprime-0" {
  group     = ibm_is_security_group.wombat-hesitate-scorn-subprime.id
  direction = "inbound"
  local     = "0.0.0.0/0"
  remote    = ibm_is_security_group.wombat-hesitate-scorn-subprime.id
}
resource "ibm_is_security_group_rule" "wombat-hesitate-scorn-subprime-1" {
  group     = ibm_is_security_group.wombat-hesitate-scorn-subprime.id
  direction = "outbound"
  local     = "0.0.0.0/0"
  remote    = "0.0.0.0/0"
}
