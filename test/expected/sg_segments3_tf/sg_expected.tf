### SG attached to test-vpc0/vsi0-subnet0
resource "ibm_is_security_group" "test-vpc0--vsi0-subnet0" {
  name           = "sg-test-vpc0--vsi0-subnet0"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0_id
}

### SG attached to test-vpc0/vsi0-subnet1
resource "ibm_is_security_group" "test-vpc0--vsi0-subnet1" {
  name           = "sg-test-vpc0--vsi0-subnet1"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0_id
}

### SG attached to test-vpc0/vsi0-subnet2
resource "ibm_is_security_group" "test-vpc0--vsi0-subnet2" {
  name           = "sg-test-vpc0--vsi0-subnet2"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0_id
}

### SG attached to test-vpc0/vsi0-subnet3
resource "ibm_is_security_group" "test-vpc0--vsi0-subnet3" {
  name           = "sg-test-vpc0--vsi0-subnet3"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0_id
}

### SG attached to test-vpc0/vsi0-subnet4
resource "ibm_is_security_group" "test-vpc0--vsi0-subnet4" {
  name           = "sg-test-vpc0--vsi0-subnet4"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0_id
}
# Internal. required-connections[0]: (segment subnetSegment)->(segment nifSegment); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc0--vsi0-subnet4-0" {
  group     = ibm_is_security_group.test-vpc0--vsi0-subnet4.id
  direction = "inbound"
  remote    = "10.240.64.0/24"
}
# Internal. required-connections[0]: (segment subnetSegment)->(segment nifSegment); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc0--vsi0-subnet4-1" {
  group     = ibm_is_security_group.test-vpc0--vsi0-subnet4.id
  direction = "inbound"
  remote    = "10.240.128.0/24"
}

### SG attached to test-vpc0/vsi0-subnet5
resource "ibm_is_security_group" "test-vpc0--vsi0-subnet5" {
  name           = "sg-test-vpc0--vsi0-subnet5"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0_id
}

### SG attached to test-vpc0/vsi1-subnet0
resource "ibm_is_security_group" "test-vpc0--vsi1-subnet0" {
  name           = "sg-test-vpc0--vsi1-subnet0"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0_id
}

### SG attached to test-vpc0/vsi1-subnet1
resource "ibm_is_security_group" "test-vpc0--vsi1-subnet1" {
  name           = "sg-test-vpc0--vsi1-subnet1"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0_id
}

### SG attached to test-vpc0/vsi1-subnet2
resource "ibm_is_security_group" "test-vpc0--vsi1-subnet2" {
  name           = "sg-test-vpc0--vsi1-subnet2"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0_id
}

### SG attached to test-vpc0/vsi1-subnet3
resource "ibm_is_security_group" "test-vpc0--vsi1-subnet3" {
  name           = "sg-test-vpc0--vsi1-subnet3"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0_id
}

### SG attached to test-vpc0/vsi1-subnet4
resource "ibm_is_security_group" "test-vpc0--vsi1-subnet4" {
  name           = "sg-test-vpc0--vsi1-subnet4"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0_id
}

### SG attached to test-vpc0/vsi1-subnet5
resource "ibm_is_security_group" "test-vpc0--vsi1-subnet5" {
  name           = "sg-test-vpc0--vsi1-subnet5"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc0_id
}
# Internal. required-connections[0]: (segment subnetSegment)->(segment nifSegment); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc0--vsi1-subnet5-0" {
  group     = ibm_is_security_group.test-vpc0--vsi1-subnet5.id
  direction = "inbound"
  remote    = "10.240.64.0/24"
}
# Internal. required-connections[0]: (segment subnetSegment)->(segment nifSegment); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc0--vsi1-subnet5-1" {
  group     = ibm_is_security_group.test-vpc0--vsi1-subnet5.id
  direction = "inbound"
  remote    = "10.240.128.0/24"
}

### SG attached to test-vpc1/vsi0-subnet10
resource "ibm_is_security_group" "test-vpc1--vsi0-subnet10" {
  name           = "sg-test-vpc1--vsi0-subnet10"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc1_id
}
# Internal. required-connections[0]: (segment subnetSegment)->(segment nifSegment); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc1--vsi0-subnet10-0" {
  group     = ibm_is_security_group.test-vpc1--vsi0-subnet10.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc0--vsi0-subnet4.id
}
# Internal. required-connections[0]: (segment subnetSegment)->(segment nifSegment); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc1--vsi0-subnet10-1" {
  group     = ibm_is_security_group.test-vpc1--vsi0-subnet10.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc0--vsi1-subnet5.id
}

### SG attached to test-vpc1/vsi0-subnet11
resource "ibm_is_security_group" "test-vpc1--vsi0-subnet11" {
  name           = "sg-test-vpc1--vsi0-subnet11"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc1_id
}

### SG attached to test-vpc2/vsi0-subnet20
resource "ibm_is_security_group" "test-vpc2--vsi0-subnet20" {
  name           = "sg-test-vpc2--vsi0-subnet20"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc2_id
}
# Internal. required-connections[0]: (segment subnetSegment)->(segment nifSegment); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc2--vsi0-subnet20-0" {
  group     = ibm_is_security_group.test-vpc2--vsi0-subnet20.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc0--vsi0-subnet4.id
}
# Internal. required-connections[0]: (segment subnetSegment)->(segment nifSegment); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc2--vsi0-subnet20-1" {
  group     = ibm_is_security_group.test-vpc2--vsi0-subnet20.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc0--vsi1-subnet5.id
}

### SG attached to test-vpc2/vsi1-subnet20
resource "ibm_is_security_group" "test-vpc2--vsi1-subnet20" {
  name           = "sg-test-vpc2--vsi1-subnet20"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc2_id
}
# Internal. required-connections[0]: (segment subnetSegment)->(segment nifSegment); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc2--vsi1-subnet20-0" {
  group     = ibm_is_security_group.test-vpc2--vsi1-subnet20.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc0--vsi0-subnet4.id
}
# Internal. required-connections[0]: (segment subnetSegment)->(segment nifSegment); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc2--vsi1-subnet20-1" {
  group     = ibm_is_security_group.test-vpc2--vsi1-subnet20.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc0--vsi1-subnet5.id
}

### SG attached to test-vpc2/vsi2-subnet20
resource "ibm_is_security_group" "test-vpc2--vsi2-subnet20" {
  name           = "sg-test-vpc2--vsi2-subnet20"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc2_id
}
# Internal. required-connections[0]: (segment subnetSegment)->(segment nifSegment); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc2--vsi2-subnet20-0" {
  group     = ibm_is_security_group.test-vpc2--vsi2-subnet20.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc0--vsi0-subnet4.id
}
# Internal. required-connections[0]: (segment subnetSegment)->(segment nifSegment); allowed-protocols[0]
resource "ibm_is_security_group_rule" "test-vpc2--vsi2-subnet20-1" {
  group     = ibm_is_security_group.test-vpc2--vsi2-subnet20.id
  direction = "outbound"
  remote    = ibm_is_security_group.test-vpc0--vsi1-subnet5.id
}

### SG attached to test-vpc3/vsi0-subnet30
resource "ibm_is_security_group" "test-vpc3--vsi0-subnet30" {
  name           = "sg-test-vpc3--vsi0-subnet30"
  resource_group = local.sg_synth_resource_group_id
  vpc            = local.sg_synth_test-vpc3_id
}
