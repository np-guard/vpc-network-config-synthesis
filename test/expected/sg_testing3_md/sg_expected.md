 | SG | Direction | Local | Remote type | Remote | Protocol | Protocol params | Description | 
 |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  | 
 | test-vpc/be | Inbound | 0.0.0.0/0 | Security group | test-vpc/fe | TCP | any port | Internal. required-connections[2]: (instance test-vpc/fe)->(instance test-vpc/be); allowed-protocols[0] | 
 | test-vpc/be | Outbound | 0.0.0.0/0 | Security group | test-vpc/opa | ALL |  | Internal. required-connections[3]: (instance test-vpc/be)->(instance test-vpc/opa); allowed-protocols[0] | 
 | test-vpc/be | Outbound | 0.0.0.0/0 | Security group | test-vpc/policydb-endpoint-gateway | ALL |  | Internal. required-connections[4]: (instance test-vpc/be)->(vpe test-vpc/policydb-endpoint-gateway); allowed-protocols[0] | 
 | test-vpc/fe | Inbound | 0.0.0.0/0 | Security group | test-vpc/proxy | TCP | ports 9000-9000 | Internal. required-connections[1]: (instance test-vpc/proxy)->(instance test-vpc/fe); allowed-protocols[0] | 
 | test-vpc/fe | Outbound | 0.0.0.0/0 | Security group | test-vpc/be | TCP | any port | Internal. required-connections[2]: (instance test-vpc/fe)->(instance test-vpc/be); allowed-protocols[0] | 
 | test-vpc/opa | Inbound | 0.0.0.0/0 | Security group | test-vpc/be | ALL |  | Internal. required-connections[3]: (instance test-vpc/be)->(instance test-vpc/opa); allowed-protocols[0] | 
 | test-vpc/opa | Outbound | 0.0.0.0/0 | Security group | test-vpc/policydb-endpoint-gateway | ALL |  | Internal. required-connections[5]: (instance test-vpc/opa)->(vpe test-vpc/policydb-endpoint-gateway); allowed-protocols[0] | 
 | test-vpc/policydb-endpoint-gateway | Inbound | 0.0.0.0/0 | Security group | test-vpc/be | ALL |  | Internal. required-connections[4]: (instance test-vpc/be)->(vpe test-vpc/policydb-endpoint-gateway); allowed-protocols[0] | 
 | test-vpc/policydb-endpoint-gateway | Inbound | 0.0.0.0/0 | Security group | test-vpc/opa | ALL |  | Internal. required-connections[5]: (instance test-vpc/opa)->(vpe test-vpc/policydb-endpoint-gateway); allowed-protocols[0] | 
 | test-vpc/proxy | Inbound | 0.0.0.0/0 | CIDR block | Any IP | ALL |  | External. required-connections[0]: (external public internet)->(instance test-vpc/proxy); allowed-protocols[0] | 
 | test-vpc/proxy | Outbound | 0.0.0.0/0 | Security group | test-vpc/fe | TCP | ports 9000-9000 | Internal. required-connections[1]: (instance test-vpc/proxy)->(instance test-vpc/fe); allowed-protocols[0] | 
