### SG attached to test-vpc0/vsi0-subnet0
resource "ibm_is_security_group" "test-vpc0--vsi0-subnet0" {
  name           = "sg-test-vpc0--vsi0-subnet0"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0--vsi0-subnet0_id
}
# Internal. required-connections[0]: (instance test-vpc0/vsi0-subnet0)->(instance test-vpc0/vsi1-subnet4); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc0--vsi0-subnet0-0" {
  group     = ibm_is_security_group.test-vpc0--vsi0-subnet0.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc0--vsi1-subnet4.id
}

### SG attached to test-vpc0/vsi0-subnet1
resource "ibm_is_security_group" "test-vpc0--vsi0-subnet1" {
  name           = "sg-test-vpc0--vsi0-subnet1"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0--vsi0-subnet1_id
}

### SG attached to test-vpc0/vsi0-subnet2
resource "ibm_is_security_group" "test-vpc0--vsi0-subnet2" {
  name           = "sg-test-vpc0--vsi0-subnet2"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0--vsi0-subnet2_id
}

### SG attached to test-vpc0/vsi0-subnet3
resource "ibm_is_security_group" "test-vpc0--vsi0-subnet3" {
  name           = "sg-test-vpc0--vsi0-subnet3"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0--vsi0-subnet3_id
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

### SG attached to test-vpc0/vsi1-subnet1
resource "ibm_is_security_group" "test-vpc0--vsi1-subnet1" {
  name           = "sg-test-vpc0--vsi1-subnet1"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0--vsi1-subnet1_id
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
# Internal. required-connections[0]: (instance test-vpc0/vsi0-subnet0)->(instance test-vpc0/vsi1-subnet4); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc0--vsi1-subnet4-0" {
  group     = ibm_is_security_group.test-vpc0--vsi1-subnet4.id
  direction = "inbound"
  remote    = ibm_is_security_group.test-vpc0--vsi0-subnet0.id
}

### SG attached to test-vpc0/vsi1-subnet5
resource "ibm_is_security_group" "test-vpc0--vsi1-subnet5" {
  name           = "sg-test-vpc0--vsi1-subnet5"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0--vsi1-subnet5_id
}
