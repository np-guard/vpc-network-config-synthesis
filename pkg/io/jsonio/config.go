package jsonio

import (
	"encoding/json"
	"os"

	"github.com/IBM/vpc-go-sdk/vpcv1"

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
	instanceToNIF := make(map[string][]string)
	for _, instance := range config.InstanceList {
		nifs := make([]string, len(instance.NetworkInterfaces))
		for i := range instance.NetworkInterfaces {
			nif := &instance.NetworkInterfaces[i]
			nifs[i] = *nif.Name
			nifToIP[*nif.Name] = ir.IPFromString(*nif.PrimaryIP.Address)
		}
		instanceToNIF[*instance.Name] = nifs
	}

	vpeToIP := make(map[string]ir.IP)
	for _, subnet := range config.SubnetList {
		for _, r := range subnet.ReservedIps {
			if t, ok := r.Target.(*vpcv1.ReservedIPTarget); ok && t != nil && r.Address != nil {
				if r.ResourceType != nil && *t.ResourceType == "endpoint_gateway" && t.Name != nil {
					vpeToIP[*t.Name] = ir.IPFromString(*r.Address)
				}
			}
		}
	}

	return &ir.ConfigDefs{
		Subnets:        subnetMap,
		NIFToIP:        nifToIP,
		InstanceToNIFs: instanceToNIF,
		VPEToIP:        vpeToIP,
	}, nil
}
