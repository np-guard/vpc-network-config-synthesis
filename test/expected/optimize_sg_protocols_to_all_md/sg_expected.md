 | SG | Direction | Remote type | Remote | Protocol | Protocol params | Description | 
 |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  | 
 | test-vpc1--vsi1 | Outbound | Security group | test-vpc1--vsi2 | ALL |  |  | 
 | test-vpc1--vsi1 | Outbound | CIDR block | 0.0.0.0/31 | TCP | ports 1-1 |  | 
 | test-vpc1--vsi1 | Outbound | CIDR block | 0.0.0.0/31 | UDP | ports 1-1 |  | 
 | test-vpc1--vsi1 | Outbound | CIDR block | 0.0.0.0/31 | ICMP | Type: Any, Code: Any |  | 
 | test-vpc1--vsi2 | Inbound | Security group | test-vpc1--vsi1 | ALL |  |  | 
 | wombat-hesitate-scorn-subprime | Inbound | Security group | wombat-hesitate-scorn-subprime | ALL |  |  | 