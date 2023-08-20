package jsonio

import (
	"encoding/json"
	"os"
)

func ReadSubnetMap(filename string) (map[string]string, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := map[string][]map[string]interface{}{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	subnetMap := make(map[string]string)
	for _, subnet := range config["subnets"] {
		subnetMap[subnet["name"].(string)] = subnet["ipv4_cidr_block"].(string)
	}
	return subnetMap, nil
}
