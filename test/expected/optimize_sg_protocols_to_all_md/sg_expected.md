 | SG | Direction | Remote type | Remote | Protocol | Protocol params | Description | 
 |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  | 
 | sg1 | Inbound | CIDR block | Any IP | ALL |  |  | 
 | sg1 | Outbound | CIDR block | Any IP | ALL |  |  | 
 | test-vpc1--vsi1 | Outbound | Security group | test-vpc1--vsi2 | ALL |  |  | 
 | test-vpc1--vsi1 | Outbound | CIDR block | 0.0.0.0/30 | ICMP | Type: Any, Code: Any |  | 
 | test-vpc1--vsi1 | Outbound | CIDR block | 0.0.0.0/31 | ALL |  |  | 
 | test-vpc1--vsi2 | Inbound | Security group | test-vpc1--vsi1 | ALL |  |  | 
 | wombat-hesitate-scorn-subprime | Inbound | Security group | wombat-hesitate-scorn-subprime | ALL |  |  | 
 | wombat-hesitate-scorn-subprime | Outbound | CIDR block | Any IP | ALL |  |  | 
