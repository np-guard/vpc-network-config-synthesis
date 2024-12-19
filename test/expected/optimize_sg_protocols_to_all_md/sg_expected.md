 | SG | Direction | Local | Remote type | Remote | Protocol | Protocol params | Description | 
 |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  | 
 | sg1 | Inbound | 0.0.0.0/0 | CIDR block | Any IP | ALL |  |  | 
 | sg1 | Outbound | 0.0.0.0/0 | CIDR block | Any IP | ALL |  |  | 
 | test-vpc1--vsi1 | Outbound | 0.0.0.0/0 | Security group | test-vpc1--vsi2 | ALL |  |  | 
 | test-vpc1--vsi1 | Outbound | 0.0.0.0/0 | CIDR block | 0.0.0.0/30 | ICMP | Type: Any, Code: Any |  | 
 | test-vpc1--vsi1 | Outbound | 0.0.0.0/0 | CIDR block | 0.0.0.0/31 | ALL |  |  | 
 | test-vpc1--vsi1 | Outbound | 10.240.0.0/16 | Security group | test-vpc1--vsi3a | ALL |  |  | 
 | test-vpc1--vsi2 | Inbound | 0.0.0.0/0 | Security group | test-vpc1--vsi1 | ALL |  |  | 
 | test-vpc1--vsi3a | Inbound | 10.240.0.0/16 | Security group | test-vpc1--vsi1 | ALL |  |  | 
 | wombat-hesitate-scorn-subprime | Inbound | 0.0.0.0/0 | Security group | wombat-hesitate-scorn-subprime | ALL |  |  | 
 | wombat-hesitate-scorn-subprime | Outbound | 0.0.0.0/0 | CIDR block | Any IP | ALL |  |  | 
