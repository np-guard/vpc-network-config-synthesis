package confio

import (
	"encoding/json"
	"os"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	configModel "github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

func readModel(filename string) (*configModel.ResourcesContainerModel, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	model := configModel.ResourcesContainerModel{}
	err = json.Unmarshal(bytes, &model)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func ReadDefs(filename string) (*ir.ConfigDefs, error) {
	config, err := readModel(filename)
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

	addressPrefixes := make([]ir.CIDR, 0)
	for _, vpc := range config.VpcList {
		for _, addressPrefix := range vpc.AddressPrefixes {
			addressPrefixes = append(addressPrefixes, ir.CidrFromString(*addressPrefix.CIDR))
		}
	}

	return &ir.ConfigDefs{
		Subnets:         subnetMap,
		NIFToIP:         nifToIP,
		InstanceToNIFs:  instanceToNIF,
		VPEToIP:         vpeToIP,
		AddressPrefixes: addressPrefixes,
	}, nil
}
