 | Acl | Subnet | Direction | Rule priority | Allow or deny | Protocol | Source | Destination | Value | Description | 
 |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  | 
 | acl1 | subnet1 | Inbound | 1 | Allow | ALL | 172.217.22.46 | 10.240.10.0/24 | - |  | 
 | acl1 | subnet1 | Inbound | 2 | Allow | ALL | 10.240.20.0/24 | 10.240.10.0/24 | - |  | 
 | acl1 | subnet1 | Inbound | 3 | Allow | TCP | 10.240.30.0/24, any port | 10.240.10.0/24, any port | - |  | 
 | acl1 | subnet1 | Inbound | 4 | Allow | TCP | 10.240.30.0/24, any port | 10.240.10.0/24, ports 443-443 | - |  | 
 | acl1 | subnet1 | Outbound | 5 | Allow | ALL | 10.240.10.0/24 | 172.217.22.46 | - |  | 
 | acl1 | subnet1 | Outbound | 6 | Allow | ALL | 10.240.10.0/24 | 10.240.20.0/24 | - |  | 
 | acl1 | subnet1 | Outbound | 7 | Allow | TCP | 10.240.10.0/24, any port | 10.240.30.0/24, ports 443-443 | - |  | 
 | acl1 | subnet1 | Outbound | 8 | Allow | TCP | 10.240.10.0/24, any port | 10.240.30.0/24, any port | - |  | 
 | acl2 | subnet2 | Inbound | 1 | Allow | ALL | Any IP | 10.240.20.0/24 | - |  | 
 | acl2 | subnet2 | Inbound | 2 | Allow | ALL | 10.240.10.0/24 | 10.240.20.0/24 | - |  | 
 | acl2 | subnet2 | Outbound | 3 | Allow | ALL | 10.240.20.0/24 | Any IP | - |  | 
 | acl2 | subnet2 | Outbound | 4 | Allow | ALL | 10.240.20.0/24 | 10.240.10.0/24 | - |  | 
 | acl3 | subnet3 | Inbound | 1 | Allow | TCP | 10.240.10.0/24, any port | 10.240.30.0/24, ports 443-443 | - |  | 
 | acl3 | subnet3 | Inbound | 2 | Allow | TCP | 10.240.10.0/24, any port | 10.240.30.0/24, any port | - |  | 
 | acl3 | subnet3 | Outbound | 3 | Allow | TCP | 10.240.30.0/24, any port | 10.240.10.0/24, ports 443-443 | - |  | 
 | acl3 | subnet3 | Outbound | 4 | Allow | TCP | 10.240.30.0/24, any port | 10.240.10.0/24, any port | - |  | 
 | capitol-siren-chirpy-doornail |  | Inbound | 1 | Allow | ALL | Any IP | Any IP | - |  | 
 | capitol-siren-chirpy-doornail |  | Outbound | 2 | Allow | ALL | Any IP | Any IP | - |  | 
