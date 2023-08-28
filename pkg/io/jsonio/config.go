package jsonio

import (
	"encoding/json"
	"os"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

func ReadSubnetMap(filename string) (map[string]ir.IP, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := map[string][]map[string]interface{}{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	subnetMap := make(map[string]ir.IP)
	for _, subnet := range config["subnets"] {
		subnetMap[subnet["name"].(string)] = ir.IPFromString(subnet["ipv4_cidr_block"].(string))
	}
	return subnetMap, nil
}
