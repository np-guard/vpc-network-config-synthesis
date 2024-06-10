 | SG | Direction | Remote type | Remote | Protocol | Protocol params | Description | 
 |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  | 
 | test-vpc/be | Outbound | IP address | 0.0.0.0 | ALL |  | External. required-connections[8]: (instance test-vpc/be)->(external external5); allowed-protocols[0] | 
 | test-vpc/fe | Outbound | CIDR block | Any IP | ALL |  | External. required-connections[4]: (instance test-vpc/fe)->(external external1); allowed-protocols[0] | 
 | test-vpc/fe | Outbound | IP address | 8.8.8.8 | ALL |  | External. required-connections[5]: (instance test-vpc/fe)->(external external2); allowed-protocols[0] | 
 | test-vpc/fe | Outbound | IP address | 7.7.7.7 | ALL |  | External. required-connections[6]: (instance test-vpc/fe)->(external external3); allowed-protocols[0] | 
 | test-vpc/fe | Outbound | CIDR block | 5.5.0.0/16 | ALL |  | External. required-connections[7]: (instance test-vpc/fe)->(external external4); allowed-protocols[0] | 
 | test-vpc/opa | Inbound | IP address | 0.0.0.0 | ALL |  | External. required-connections[9]: (external external5)->(instance test-vpc/opa); allowed-protocols[0] | 
 | test-vpc/proxy | Inbound | CIDR block | Any IP | ALL |  | External. required-connections[0]: (external external1)->(instance test-vpc/proxy); allowed-protocols[0] | 
 | test-vpc/proxy | Inbound | IP address | 8.8.8.8 | ALL |  | External. required-connections[1]: (external external2)->(instance test-vpc/proxy); allowed-protocols[0] | 
 | test-vpc/proxy | Inbound | IP address | 7.7.7.7 | ALL |  | External. required-connections[2]: (external external3)->(instance test-vpc/proxy); allowed-protocols[0] | 
 | test-vpc/proxy | Inbound | CIDR block | 5.5.0.0/16 | ALL |  | External. required-connections[3]: (external external4)->(instance test-vpc/proxy); allowed-protocols[0] | 
