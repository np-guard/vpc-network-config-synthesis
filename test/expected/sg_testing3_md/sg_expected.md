 | SG | Direction | Remote type | Remote | Protocol | Protocol params | Description | 
 |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  | 
 | test-vpc/be | Inbound | Security group | test-vpc/fe | TCP | any port | Internal. required-connections[2]: (instance test-vpc/fe)->(instance test-vpc/be); allowed-protocols[0] | 
 | test-vpc/be | Outbound | Security group | test-vpc/opa | ALL |  | Internal. required-connections[3]: (instance test-vpc/be)->(instance test-vpc/opa); allowed-protocols[0] | 
 | test-vpc/be | Outbound | Security group | test-vpc/policydb-endpoint-gateway | ALL |  | Internal. required-connections[4]: (instance test-vpc/be)->(vpe test-vpc/policydb-endpoint-gateway); allowed-protocols[0] | 
 | test-vpc/fe | Inbound | Security group | test-vpc/proxy | TCP | ports 9000-9000 | Internal. required-connections[1]: (instance test-vpc/proxy)->(instance test-vpc/fe); allowed-protocols[0] | 
 | test-vpc/fe | Outbound | Security group | test-vpc/be | TCP | any port | Internal. required-connections[2]: (instance test-vpc/fe)->(instance test-vpc/be); allowed-protocols[0] | 
 | test-vpc/opa | Inbound | Security group | test-vpc/be | ALL |  | Internal. required-connections[3]: (instance test-vpc/be)->(instance test-vpc/opa); allowed-protocols[0] | 
 | test-vpc/opa | Outbound | Security group | test-vpc/policydb-endpoint-gateway | ALL |  | Internal. required-connections[5]: (instance test-vpc/opa)->(vpe test-vpc/policydb-endpoint-gateway); allowed-protocols[0] | 
 | test-vpc/policydb-endpoint-gateway | Inbound | Security group | test-vpc/be | ALL |  | Internal. required-connections[4]: (instance test-vpc/be)->(vpe test-vpc/policydb-endpoint-gateway); allowed-protocols[0] | 
 | test-vpc/policydb-endpoint-gateway | Inbound | Security group | test-vpc/opa | ALL |  | Internal. required-connections[5]: (instance test-vpc/opa)->(vpe test-vpc/policydb-endpoint-gateway); allowed-protocols[0] | 
 | test-vpc/proxy | Inbound | CIDR block | Any IP | ALL |  | External. required-connections[0]: (external public internet)->(instance test-vpc/proxy); allowed-protocols[0] | 
 | test-vpc/proxy | Outbound | Security group | test-vpc/fe | TCP | ports 9000-9000 | Internal. required-connections[1]: (instance test-vpc/proxy)->(instance test-vpc/fe); allowed-protocols[0] | 
