 | Acl | Subnet | Direction | Rule priority | Allow or deny | Protocol | Source | Destination | Value | Description | 
 |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  | 
 | acl-1 | 1 | Outbound | 1 | Allow | TCP | 10.240.10.0/24, any port | 10.240.30.0/24, ports 443-443 | - | Internal. required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0] | 
 | acl-1 | 1 | Inbound | 2 | Allow | TCP | 10.240.30.0/24, ports 443-443 | 10.240.10.0/24, any port | - | Internal. response to required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0] | 
 | acl-1 | 1 | Inbound | 3 | Allow | TCP | 10.240.10.0/24, any port | 10.240.30.0/24, ports 443-443 | - | Internal. required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0] | 
 | acl-1 | 1 | Outbound | 4 | Allow | TCP | 10.240.30.0/24, ports 443-443 | 10.240.10.0/24, any port | - | Internal. response to required-connections[0]: (subnet subnet1-ky)->(subnet subnet3-ky); allowed-protocols[0] | 
