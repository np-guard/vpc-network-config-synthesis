/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package confio

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	configModel "github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"

	"github.com/np-guard/models/pkg/ipblock"
)

const EndpointVPE string = "endpoint_gateway"

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

	instances, nifs := parseInstancesNifs(config)
	vpes, vpeEndpoints := parseVPEs(config)
	vpcs := parseVPCs(config)
	err = validateVpcs(vpcs)
	if err != nil {
		return nil, err
	}

	return &ir.ConfigDefs{
		VPCs:         vpcs,
		Subnets:      parseSubnets(config),
		NIFs:         nifs,
		Instances:    instances,
		VPEEndpoints: vpeEndpoints,
		VPEs:         vpes,
	}, nil
}

func parseVPCs(config *configModel.ResourcesContainerModel) map[ir.ID]ir.VPCDetails {
	VPCs := make(map[ir.ID]ir.VPCDetails)
	for _, vpc := range config.VpcList {
		addressPrefixes := make([]ir.CIDR, 0)
		for _, addressPrefix := range vpc.AddressPrefixes {
			addressPrefixes = append(addressPrefixes, ir.CidrFromString(*addressPrefix.CIDR))
		}
		VPCs[ir.ID(*vpc.Name)] = ir.VPCDetails{AddressPrefixes: addressPrefixes}
	}
	return VPCs
}

func parseSubnets(config *configModel.ResourcesContainerModel) map[ir.ID]ir.SubnetDetails {
	subnets := make(map[ir.ID]ir.SubnetDetails)
	for _, subnet := range config.SubnetList {
		uniqueName := ir.ID(scopingString(*subnet.VPC.Name, *subnet.Name))
		subnetDetails := ir.SubnetDetails{
			NamedEntity: ir.NamedEntity(*subnet.Name),
			VPC:         ir.ID(*subnet.VPC.Name),
			CIDR:        ir.IPFromString(*subnet.Ipv4CIDRBlock),
		}
		subnets[uniqueName] = subnetDetails
	}
	return subnets
}

func parseInstancesNifs(config *configModel.ResourcesContainerModel) (instances map[ir.ID]ir.InstanceDetails,
	nifs map[ir.ID]ir.NifDetails) {
	instances = make(map[ir.ID]ir.InstanceDetails)
	nifs = make(map[ir.ID]ir.NifDetails)
	for _, instance := range config.InstanceList {
		instanceUniqueName := scopingString(*instance.VPC.Name, *instance.Name)
		instanceNifs := make([]ir.ID, len(instance.NetworkInterfaces))
		for i := range instance.NetworkInterfaces {
			nifUniqueName := scopingString(instanceUniqueName, *instance.NetworkInterfaces[i].Name)
			nifDetails := ir.NifDetails{
				NamedEntity: ir.NamedEntity(*instance.NetworkInterfaces[i].Name),
				Instance:    ir.ID(scopingString(*instance.VPC.Name, *instance.Name)),
				IP:          ir.IPFromString(*instance.NetworkInterfaces[i].PrimaryIP.Address),
			}
			nifs[ir.ID(nifUniqueName)] = nifDetails
			instanceNifs[i] = ir.ID(nifUniqueName)
		}
		instanceDetails := ir.InstanceDetails{
			NamedEntity: ir.NamedEntity(*instance.Name),
			VPC:         ir.ID(*instance.VPC.Name),
			Nifs:        instanceNifs,
		}
		instances[ir.ID(instanceUniqueName)] = instanceDetails
	}
	return instances, nifs
}

func parseVPEs(config *configModel.ResourcesContainerModel) (vpes map[ir.ID]ir.VPEDetails, vpeEndpoints map[ir.ID]ir.VPEEndpointDetails) {
	vpes = make(map[ir.ID]ir.VPEDetails)
	vpeEndpoints = make(map[ir.ID]ir.VPEEndpointDetails)

	for _, vpe := range config.EndpointGWList {
		if *vpe.ResourceType == EndpointVPE {
			uniqueVpeName := scopingString(*vpe.VPC.Name, *vpe.Name)
			vpeDetails := ir.VPEDetails{
				NamedEntity: ir.NamedEntity(*vpe.Name),
				VPEEndpoint: []ir.ID{},
				VPC:         ir.ID(*vpe.VPC.Name),
			}
			vpes[ir.ID(uniqueVpeName)] = vpeDetails
		}
	}

	for _, subnet := range config.SubnetList {
		for _, r := range subnet.ReservedIps {
			if t, ok := r.Target.(*vpcv1.ReservedIPTarget); ok && t != nil && r.Address != nil {
				if r.ResourceType != nil && *t.ResourceType == EndpointVPE && t.Name != nil {
					VPEName := ir.ID(scopingString(*subnet.VPC.Name, *t.Name))
					subnetName := ir.ID(scopingString(*subnet.VPC.Name, *subnet.Name))
					uniqueVpeEndpointName := scopingString(string(VPEName), *r.Name)
					vpeEndpointDetails := ir.VPEEndpointDetails{
						NamedEntity: ir.NamedEntity(*r.Name),
						VPEName:     VPEName,
						Subnet:      subnetName,
						IP:          ir.IPFromString(*r.Address),
					}
					vpeEndpoints[ir.ID(uniqueVpeEndpointName)] = vpeEndpointDetails
					vpe := vpes[VPEName]
					vpe.VPEEndpoint = append(vpe.VPEEndpoint, ir.ID(uniqueVpeEndpointName))
					vpes[VPEName] = vpe
				}
			}
		}
	}

	return vpes, vpeEndpoints
}

func validateVpcs(vpcs map[ir.ID]ir.VPCDetails) error {
	for vpcName1, vpcDetails1 := range vpcs {
		for vpcName2, vpcDetails2 := range vpcs {
			if vpcName1 == vpcName2 {
				continue
			}
			for _, addressPrefix1 := range vpcDetails1.AddressPrefixes {
				for _, addressPrefix2 := range vpcDetails2.AddressPrefixes {
					address1, err := ipblock.FromCidr(addressPrefix1.String())
					if err != nil {
						return err
					}
					address2, err := ipblock.FromCidr(addressPrefix2.String())
					if err != nil {
						return err
					}
					if !address1.Intersect(address2).IsEmpty() {
						return fmt.Errorf("vpcs %s and %s are overlapping", string(vpcName1), string(vpcName2))
					}
				}
			}
		}
	}
	return nil
}

func scopingString(s1, s2 string) string {
	return s1 + "/" + s2
}
