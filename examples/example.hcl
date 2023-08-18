define {
  segment {
    private = [
      "10.10.10.0/24",
      "10.10.10.0/24",
    ]
  }

  external {
    google = "10.240.30.0/24"
  }
}

connection {

  src cidr "10.10.10.0/24" { }
  dst cidr "10.10.10.0/24" { }

  tcp {
    min_port = 80
    max_port = 100
  }

  udp {
    min_port = 80
    max_port = 100
  }

  icmp {
    type = 2
    code = 5
  }
}
