 | Acl | Subnet | Direction | Rule priority | Allow or deny | Protocol | Source | Destination | Value | Description | 
 |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  | 
 | acl1-1 | sub1-1 | Inbound | 1 | Allow | UDP | 8.8.8.8, ports 53-53 | 10.240.1.0/24, any port | - |  | 
 | acl1-1 | sub1-1 | Inbound | 2 | Allow | TCP | 10.240.2.0/23, any port | 10.240.1.0/24, any port | - |  | 
 | acl1-1 | sub1-1 | Inbound | 3 | Allow | ICMP | 10.240.128.0/24 | 10.240.1.0/24 | Type: 0, Code: 0 |  | 
 | acl1-1 | sub1-1 | Outbound | 4 | Allow | UDP | 10.240.1.0/24, any port | 8.8.8.8, ports 53-53 | - |  | 
 | acl1-1 | sub1-1 | Outbound | 5 | Allow | TCP | 10.240.1.0/24, any port | 10.240.2.0/23, any port | - |  | 
 | acl1-1 | sub1-1 | Outbound | 6 | Allow | ICMP | 10.240.1.0/24 | 10.240.128.0/24 | Type: 0, Code: 0 |  | 
 | acl1-2 | sub1-2, sub1-3 | Inbound | 1 | Allow | TCP | 10.240.1.0/24, any port | 10.240.2.0/23, any port | - |  | 
 | acl1-2 | sub1-2, sub1-3 | Inbound | 2 | Allow | TCP | 10.240.2.0/23, any port | 10.240.2.0/23, any port | - |  | 
 | acl1-2 | sub1-2, sub1-3 | Outbound | 3 | Allow | TCP | 10.240.2.0/23, any port | 10.240.1.0/24, any port | - |  | 
 | acl1-2 | sub1-2, sub1-3 | Outbound | 4 | Allow | TCP | 10.240.2.0/23, any port | 10.240.2.0/23, any port | - |  | 
 | acl2-1 | sub2-1 | Inbound | 1 | Allow | UDP | 8.8.8.8, ports 53-53 | 10.240.64.0/24, any port | - |  | 
 | acl2-1 | sub2-1 | Inbound | 2 | Allow | ALL | 10.240.65.0/24 | 10.240.64.0/24 | - |  | 
 | acl2-1 | sub2-1 | Inbound | 3 | Allow | TCP | 10.240.128.0/24, any port | 10.240.64.0/24, ports 443-443 | - |  | 
 | acl2-1 | sub2-1 | Inbound | 4 | Allow | ICMP | 10.240.128.0/24 | 10.240.64.0/24 | Type: 0, Code: 0 |  | 
 | acl2-1 | sub2-1 | Outbound | 5 | Allow | UDP | 10.240.64.0/24, any port | 8.8.8.8, ports 53-53 | - |  | 
 | acl2-1 | sub2-1 | Outbound | 6 | Allow | ALL | 10.240.64.0/24 | 10.240.65.0/24 | - |  | 
 | acl2-1 | sub2-1 | Outbound | 7 | Allow | TCP | 10.240.64.0/24, ports 443-443 | 10.240.128.0/24, any port | - |  | 
 | acl2-1 | sub2-1 | Outbound | 8 | Allow | ICMP | 10.240.64.0/24 | 10.240.128.0/24 | Type: 0, Code: 0 |  | 
 | acl2-2 | sub2-2 | Inbound | 1 | Allow | ALL | 10.240.64.0/24 | 10.240.65.0/24 | - |  | 
 | acl2-2 | sub2-2 | Outbound | 2 | Allow | ALL | 10.240.65.0/24 | 10.240.64.0/24 | - |  | 
 | acl3-1 | sub3-1 | Inbound | 1 | Allow | TCP | 10.240.64.0/24, ports 443-443 | 10.240.128.0/24, any port | - |  | 
 | acl3-1 | sub3-1 | Inbound | 2 | Allow | ICMP | 10.240.64.0/24 | 10.240.128.0/24 | Type: 0, Code: 0 |  | 
 | acl3-1 | sub3-1 | Inbound | 3 | Allow | ICMP | 10.240.1.0/24 | 10.240.128.0/24 | Type: 0, Code: 0 |  | 
 | acl3-1 | sub3-1 | Outbound | 4 | Allow | TCP | 10.240.128.0/24, any port | 10.240.64.0/24, ports 443-443 | - |  | 
 | acl3-1 | sub3-1 | Outbound | 5 | Allow | ICMP | 10.240.128.0/24 | 10.240.64.0/24 | Type: 0, Code: 0 |  | 
 | acl3-1 | sub3-1 | Outbound | 6 | Allow | ICMP | 10.240.128.0/24 | 10.240.1.0/24 | Type: 0, Code: 0 |  | 
 | disallow-laborious-compress-abiding |  | Inbound | 1 | Allow | ALL | Any IP | Any IP | - |  | 
 | disallow-laborious-compress-abiding |  | Outbound | 2 | Allow | ALL | Any IP | Any IP | - |  | 
