 | Acl | Subnet | Direction | Rule priority | Allow or deny | Source | Destination | Protocol | Value | Description | 
 |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  | 
 | acl1-1 | sub1-1 | Inbound | 1 | Allow | 8.8.8.8, ports 53-53 | 10.240.1.0/24, any port | UDP | - |  | 
 | acl1-1 | sub1-1 | Inbound | 2 | Allow | 10.240.2.0/23, any port | 10.240.1.0/24, any port | TCP | - |  | 
 | acl1-1 | sub1-1 | Inbound | 3 | Allow | 10.240.128.0/24 | 10.240.1.0/24 | ICMP | Type: 0, Code: 0 |  | 
 | acl1-1 | sub1-1 | Outbound | 4 | Allow | 10.240.1.0/24, any port | 8.8.8.8, ports 53-53 | UDP | - |  | 
 | acl1-1 | sub1-1 | Outbound | 5 | Allow | 10.240.1.0/24, any port | 10.240.2.0/23, any port | TCP | - |  | 
 | acl1-1 | sub1-1 | Outbound | 6 | Allow | 10.240.1.0/24 | 10.240.128.0/24 | ICMP | Type: 0, Code: 0 |  | 
 | acl1-2 | sub1-2, sub1-3 | Inbound | 1 | Allow | 10.240.1.0/24, any port | 10.240.2.0/23, any port | TCP | - |  | 
 | acl1-2 | sub1-2, sub1-3 | Inbound | 2 | Allow | 10.240.2.0/23, any port | 10.240.2.0/23, any port | TCP | - |  | 
 | acl1-2 | sub1-2, sub1-3 | Outbound | 3 | Allow | 10.240.2.0/23, any port | 10.240.1.0/24, any port | TCP | - |  | 
 | acl1-2 | sub1-2, sub1-3 | Outbound | 4 | Allow | 10.240.2.0/23, any port | 10.240.2.0/23, any port | TCP | - |  | 
 | acl2-1 | sub2-1 | Inbound | 1 | Allow | 8.8.8.8, ports 53-53 | 10.240.64.0/24, any port | UDP | - |  | 
 | acl2-1 | sub2-1 | Inbound | 2 | Allow | 10.240.65.0/24 | 10.240.64.0/24 | ALL | - |  | 
 | acl2-1 | sub2-1 | Inbound | 3 | Allow | 10.240.128.0/24, any port | 10.240.64.0/24, ports 443-443 | TCP | - |  | 
 | acl2-1 | sub2-1 | Inbound | 4 | Allow | 10.240.128.0/24 | 10.240.64.0/24 | ICMP | Type: 0, Code: 0 |  | 
 | acl2-1 | sub2-1 | Outbound | 5 | Allow | 10.240.64.0/24, any port | 8.8.8.8, ports 53-53 | UDP | - |  | 
 | acl2-1 | sub2-1 | Outbound | 6 | Allow | 10.240.64.0/24 | 10.240.65.0/24 | ALL | - |  | 
 | acl2-1 | sub2-1 | Outbound | 7 | Allow | 10.240.64.0/24, ports 443-443 | 10.240.128.0/24, any port | TCP | - |  | 
 | acl2-1 | sub2-1 | Outbound | 8 | Allow | 10.240.64.0/24 | 10.240.128.0/24 | ICMP | Type: 0, Code: 0 |  | 
 | acl2-2 | sub2-2 | Inbound | 1 | Allow | 10.240.64.0/24 | 10.240.65.0/24 | ALL | - |  | 
 | acl2-2 | sub2-2 | Outbound | 2 | Allow | 10.240.65.0/24 | 10.240.64.0/24 | ALL | - |  | 
 | acl3-1 | sub3-1 | Inbound | 1 | Allow | 10.240.64.0/24, ports 443-443 | 10.240.128.0/24, any port | TCP | - |  | 
 | acl3-1 | sub3-1 | Inbound | 2 | Allow | 10.240.64.0/24 | 10.240.128.0/24 | ICMP | Type: 0, Code: 0 |  | 
 | acl3-1 | sub3-1 | Inbound | 3 | Allow | 10.240.1.0/24 | 10.240.128.0/24 | ICMP | Type: 0, Code: 0 |  | 
 | acl3-1 | sub3-1 | Outbound | 4 | Allow | 10.240.128.0/24, any port | 10.240.64.0/24, ports 443-443 | TCP | - |  | 
 | acl3-1 | sub3-1 | Outbound | 5 | Allow | 10.240.128.0/24 | 10.240.64.0/24 | ICMP | Type: 0, Code: 0 |  | 
 | acl3-1 | sub3-1 | Outbound | 6 | Allow | 10.240.128.0/24 | 10.240.1.0/24 | ICMP | Type: 0, Code: 0 |  | 
 | disallow-laborious-compress-abiding |  | Inbound | 1 | Allow | Any IP | Any IP | ALL | - |  | 
 | disallow-laborious-compress-abiding |  | Outbound | 2 | Allow | Any IP | Any IP | ALL | - |  | 
