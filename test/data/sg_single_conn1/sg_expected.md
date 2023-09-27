 | SG | Direction | Protocol | Remote type | Remote | Protocol params | Description | 
 |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  |  :---  | 
 | A | Outbound | Security group | B | TCP | Ports 443-443 | Internal. required-connections[0]: (nif ni3b)->(nif ni2); allowed-protocols[0] | 
 | B | Inbound | Security group | A | TCP | Ports 443-443 | Internal. required-connections[0]: (nif ni3b)->(nif ni2); allowed-protocols[0] | 
