/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package confio

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	configModel "github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"

	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

const EndpointVPE string = "endpoint_gateway"

func ReadDefs(filename string) (*ir.ConfigDefs, error) {
	config, err := ReadModel(filename)
	if err != nil {
		return nil, err
	}

	subnets, err := parseSubnets(config)
	if err != nil {
		return nil, err
	}
	instances, nifs, err := parseInstancesNifs(config)
	if err != nil {
		return nil, err
	}
	vpes, vpeEndpoints, err := parseVPEs(config)
	if err != nil {
		return nil, err
	}
	vpcs, err := ParseVPCs(config)
	if err != nil {
		return nil, err
	}

	return &ir.ConfigDefs{
		VPCs:           vpcs,
		Subnets:        subnets,
		NIFs:           nifs,
		Instances:      instances,
		VPEReservedIPs: vpeEndpoints,
		VPEs:           vpes,
	}, nil
}

func ParseVPCs(config *configModel.ResourcesContainerModel) (map[ir.ID]*ir.VPCDetails, error) {
	VPCs := make(map[ir.ID]*ir.VPCDetails, len(config.VpcList))
	for _, vpc := range config.VpcList {
		addressPrefixes := netset.NewIPBlock()
		for _, addressPrefix := range vpc.AddressPrefixes {
			address, err := netset.IPBlockFromCidr(*addressPrefix.CIDR)
			if err != nil {
				return nil, err
			}
			addressPrefixes = addressPrefixes.Union(address)
		}
		VPCs[*vpc.Name] = &ir.VPCDetails{AddressPrefixes: addressPrefixes}
	}

	return VPCs, validateVpcs(VPCs)
}

func parseSubnets(config *configModel.ResourcesContainerModel) (map[ir.ID]*ir.SubnetDetails, error) {
	subnets := make(map[ir.ID]*ir.SubnetDetails, len(config.SubnetList))
	for _, subnet := range config.SubnetList {
		uniqueName := ScopingString(*subnet.VPC.Name, *subnet.Name)
		cidr, err := netset.IPBlockFromCidr(*subnet.Ipv4CIDRBlock)
		if err != nil {
			return nil, err
		}
		subnetDetails := ir.SubnetDetails{
			NamedEntity: ir.NamedEntity(*subnet.Name),
			VPC:         *subnet.VPC.Name,
			CIDR:        cidr,
		}
		subnets[uniqueName] = &subnetDetails
	}
	return subnets, nil
}

func parseInstancesNifs(config *configModel.ResourcesContainerModel) (instances map[ir.ID]*ir.InstanceDetails,
	nifs map[ir.ID]*ir.NifDetails, err error) {
	instances = make(map[ir.ID]*ir.InstanceDetails, len(config.InstanceList))
	nifs = make(map[ir.ID]*ir.NifDetails)
	for _, instance := range config.InstanceList {
		instanceUniqueName := ScopingString(*instance.VPC.Name, *instance.Name)
		instanceNifs := make([]ir.ID, len(instance.NetworkInterfaces))
		for i := range instance.NetworkInterfaces {
			nifUniqueName := ScopingString(instanceUniqueName, *instance.NetworkInterfaces[i].Name)
			nifIP, err := netset.IPBlockFromIPAddress(*instance.NetworkInterfaces[i].PrimaryIP.Address)
			if err != nil {
				return nil, nil, err
			}
			nifDetails := ir.NifDetails{
				NamedEntity: ir.NamedEntity(*instance.NetworkInterfaces[i].Name),
				Instance:    ScopingString(*instance.VPC.Name, *instance.Name),
				VPC:         *instance.VPC.Name,
				IP:          nifIP,
				Subnet:      ScopingString(*instance.VPC.Name, *instance.NetworkInterfaces[i].Subnet.Name),
			}
			nifs[nifUniqueName] = &nifDetails
			instanceNifs[i] = nifUniqueName
		}
		instanceDetails := ir.InstanceDetails{
			NamedEntity: ir.NamedEntity(*instance.Name),
			VPC:         *instance.VPC.Name,
			Nifs:        instanceNifs,
		}
		instances[instanceUniqueName] = &instanceDetails
	}
	return instances, nifs, nil
}

func parseVPEs(config *configModel.ResourcesContainerModel) (vpes map[ir.ID]*ir.VPEDetails,
	vpeReservedIPs map[ir.ID]*ir.VPEReservedIPsDetails, err error) {
	vpes = make(map[ir.ID]*ir.VPEDetails)
	vpeReservedIPs = make(map[ir.ID]*ir.VPEReservedIPsDetails)

	for _, vpe := range config.EndpointGWList {
		if *vpe.ResourceType != EndpointVPE {
			continue
		}
		uniqueVpeName := ScopingString(*vpe.VPC.Name, *vpe.Name)
		vpeDetails := ir.VPEDetails{
			NamedEntity:    ir.NamedEntity(*vpe.Name),
			VPEReservedIPs: []ir.ID{},
			VPC:            *vpe.VPC.Name,
		}
		vpes[uniqueVpeName] = &vpeDetails
	}
	for _, subnet := range config.SubnetList {
		for _, r := range subnet.ReservedIps {
			t, ok := r.Target.(*vpcv1.ReservedIPTarget)
			if !ok || t == nil || r.Address == nil || t.ResourceType == nil || *t.ResourceType != EndpointVPE || t.Name == nil {
				continue
			}
			VPEName := ScopingString(*subnet.VPC.Name, *t.Name)
			subnetName := ScopingString(*subnet.VPC.Name, *subnet.Name)
			uniqueVpeReservedIPName := ScopingString(VPEName, *r.Name)
			vpeIP, err := netset.IPBlockFromIPAddress(*r.Address)
			if err != nil {
				return nil, nil, err
			}
			vpeReservedIPDetails := ir.VPEReservedIPsDetails{
				NamedEntity: ir.NamedEntity(*r.Name),
				VPEName:     VPEName,
				Subnet:      subnetName,
				IP:          vpeIP,
				VPC:         vpes[VPEName].VPC,
			}
			vpeReservedIPs[uniqueVpeReservedIPName] = &vpeReservedIPDetails
			vpe := vpes[VPEName]
			vpe.VPEReservedIPs = append(vpe.VPEReservedIPs, uniqueVpeReservedIPName)
			vpes[VPEName] = vpe
		}
	}
	return vpes, vpeReservedIPs, nil
}

func validateVpcs(vpcs map[ir.ID]*ir.VPCDetails) error {
	for vpcName1, vpcDetails1 := range vpcs {
		for vpcName2, vpcDetails2 := range vpcs {
			if vpcName1 >= vpcName2 {
				continue
			}
			if ir.Overlap(vpcDetails1.AddressPrefixes, vpcDetails2.AddressPrefixes) {
				return fmt.Errorf("vpcs %s and %s have overlapping IP address spaces", vpcName1, vpcName2)
			}
		}
	}
	return nil
}

func ScopingString(s1, s2 string) string {
	return s1 + "/" + s2
}
