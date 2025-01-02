 | Acl | Subnet | Direction | Rule priority | Allow or deny | Source | Destination | Protocol | Value | Description | 
 |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  | 
 | testacl5-vpc--sub1-1 | sub1-1 | Outbound | 1 | Allow | 10.240.1.0/24, any port | 1.1.1.0/31, any port | TCP | - |  | 
 | testacl5-vpc--sub1-2 | sub1-2 | Outbound | 1 | Allow | 10.240.2.0/24, any port | 1.1.1.1, ports 1-20 | UDP | - |  | 
 | testacl5-vpc--sub1-3 | sub1-3 | Outbound | 1 | Allow | 10.240.3.0/24 | 10.240.64.0/24 | ALL | - |  | 
 | testacl5-vpc--sub2-1 | sub2-1 | Inbound | 1 | Allow | 10.240.3.0/24 | 10.240.64.0/24 | ALL | - |  | 
 | testacl5-vpc--sub2-2 | sub2-2 | Outbound | 1 | Allow | 10.240.65.0/24, ports 11-20 | 10.240.128.0/24, any port | TCP | - |  | 
 | testacl5-vpc--sub3-1 | sub3-1 | Inbound | 1 | Allow | 10.240.65.0/24, ports 11-20 | 10.240.128.0/24, any port | TCP | - |  | 
