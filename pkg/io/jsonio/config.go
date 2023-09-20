package jsonio

import (
	"encoding/json"
	"os"

	configmodel "github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

func ReadDefs(filename string) (*ir.ConfigDefs, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := configmodel.ResourcesContainerModel{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	subnetMap := make(map[string]ir.IP)
	for _, subnet := range config.SubnetList {
		subnetMap[*subnet.Name] = ir.IPFromString(*subnet.Ipv4CIDRBlock)
	}

	nifToIP := make(map[string]ir.IP)
	instanceToNif := make(map[string][]string)
	for _, instance := range config.InstanceList {
		nifs := make([]string, len(instance.NetworkInterfaces))
		for i := range instance.NetworkInterfaces {
			nif := &instance.NetworkInterfaces[i]
			nifs[i] = *nif.Name
			nifToIP[*nif.Name] = ir.IPFromString(*nif.PrimaryIP.Address)
		}
		instanceToNif[*instance.Name] = nifs
	}
	return &ir.ConfigDefs{
		Subnets:        subnetMap,
		NifToIP:        nifToIP,
		InstanceToNifs: instanceToNif,
	}, nil
}
