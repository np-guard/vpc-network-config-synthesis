resource "ibm_is_security_group" "sg1" {
  name           = "sg-sg1"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc1_id
}


resource "ibm_is_security_group" "test-vpc1--vsi1" {
  name           = "sg-test-vpc1--vsi1"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc1_id
}
resource "ibm_is_security_group_rule" "test-vpc1--vsi1-0" {
  group     = ibm_is_security_group.test-vpc1--vsi1.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc1--vsi2.id
  icmp {
    type = 5
  }
}
resource "ibm_is_security_group_rule" "test-vpc1--vsi1-1" {
  group     = ibm_is_security_group.test-vpc1--vsi1.id
  direction = "outbound"
  remote    = "0.0.0.0/31"
  icmp {
  }
}


resource "ibm_is_security_group" "test-vpc1--vsi2" {
  name           = "sg-test-vpc1--vsi2"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc1_id
}
resource "ibm_is_security_group_rule" "test-vpc1--vsi2-0" {
  group     = ibm_is_security_group.test-vpc1--vsi2.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc1--vsi1.id
  icmp {
    type = 5
  }
}


resource "ibm_is_security_group" "test-vpc1--vsi3a" {
  name           = "sg-test-vpc1--vsi3a"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc1_id
}


resource "ibm_is_security_group" "test-vpc1--vsi3b" {
  name           = "sg-test-vpc1--vsi3b"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc1_id
}


resource "ibm_is_security_group" "wombat-hesitate-scorn-subprime" {
  name           = "sg-wombat-hesitate-scorn-subprime"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc1_id
}
resource "ibm_is_security_group_rule" "wombat-hesitate-scorn-subprime-0" {
  group     = ibm_is_security_group.wombat-hesitate-scorn-subprime.id
  direction = "inbound"
  remote    = ibm_is_security_group.wombat-hesitate-scorn-subprime.id
}
