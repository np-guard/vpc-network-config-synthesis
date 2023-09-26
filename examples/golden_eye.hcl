
module "segments" {
  edge = [
    ibm_is_subnet.ky-testenv-edge-subnet-1.id,
    ibm_is_subnet.ky-testenv-edge-subnet-2.id,
    ibm_is_subnet.ky-testenv-edge-subnet-3.id
  ]
  private = [
    ibm_is_subnet.ky-testenv-private-subnet-1.id,
    ibm_is_subnet.ky-testenv-private-subnet-2.id,
    ibm_is_subnet.ky-testenv-private-subnet-3.id
  ]
  transit = [
    ibm_is_subnet.ky-testenv-transit-subnet-1.id,
    ibm_is_subnet.ky-testenv-transit-subnet-2.id,
    ibm_is_subnet.ky-testenv-transit-subnet-3.id
  ]
}

module "externals" {
  public_internet = "0.0.0.0/0"
}

#generate:acl
resource "npguard_connection" "edge_external" {
  src = segments.edge
  dst = externals.public_internet

  tcp {
  }
}

#generate:acl
resource "npguard_connection" "edge_fc" {
  src = segments.edge
  dst = segments.edge

  tcp {
  }
}

#generate:acl
resource "npguard_connection" "private_edge" {
  src = segments.private
  dst = segments.edge

  tcp {
  }
}

#generate:acl
resource "npguard_connection" "private_fc" {
  src = segments.private
  dst = segments.private

  tcp {
  }
}

#generate:acl
resource "npguard_connection" "transit_private" {
  src = segments.transit
  dst = segments.private

  tcp {
  }
}

#generate:acl
resource "npguard_connection" "transit_fc" {
  src = segments.transit
  dst = segments.transit

  tcp {
  }
}
