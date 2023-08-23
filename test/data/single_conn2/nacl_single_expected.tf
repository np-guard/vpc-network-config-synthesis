resource "ibm_is_network_acl" "acl1" {
  name = "acl1-${var.initials}"
  resource_group = var.resource_group_id
  vpc = var.vpc_id
  rules {
    name = "c0,p0,[subnet1-ky->subnet3-ky],src0,dst0,request,outbound,allow"
    action = "allow"
    direction = "outbound"
    source = "10.240.10.0/24"
    destination = "10.240.30.0/24"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  rules {
    name = "c0,p0,[subnet1-ky->subnet3-ky],src0,dst0,response,inbound,allow"
    action = "allow"
    direction = "inbound"
    source = "10.240.30.0/24"
    destination = "10.240.10.0/24"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
  rules {
    name = "c0,p0,[subnet1-ky->subnet3-ky],src0,dst0,request,inbound,allow"
    action = "allow"
    direction = "inbound"
    source = "10.240.10.0/24"
    destination = "10.240.30.0/24"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  rules {
    name = "c0,p0,[subnet1-ky->subnet3-ky],src0,dst0,response,outbound,allow"
    action = "allow"
    direction = "outbound"
    source = "10.240.30.0/24"
    destination = "10.240.10.0/24"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
  rules {
    name = "c1,p0,[subnet3-ky->subnet1-ky],src0,dst0,request,outbound,allow"
    action = "allow"
    direction = "outbound"
    source = "10.240.30.0/24"
    destination = "10.240.10.0/24"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  rules {
    name = "c1,p0,[subnet3-ky->subnet1-ky],src0,dst0,response,inbound,allow"
    action = "allow"
    direction = "inbound"
    source = "10.240.10.0/24"
    destination = "10.240.30.0/24"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
  rules {
    name = "c1,p0,[subnet3-ky->subnet1-ky],src0,dst0,request,inbound,allow"
    action = "allow"
    direction = "inbound"
    source = "10.240.30.0/24"
    destination = "10.240.10.0/24"
    tcp {
      port_min = 443
      port_max = 443
    }
  }
  rules {
    name = "c1,p0,[subnet3-ky->subnet1-ky],src0,dst0,response,outbound,allow"
    action = "allow"
    direction = "outbound"
    source = "10.240.10.0/24"
    destination = "10.240.30.0/24"
    tcp {
      source_port_min = 443
      source_port_max = 443
    }
  }
}
