/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package confio

import (
	"errors"
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	configModel "github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"

	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

const EndpointVPE string = "endpoint_gateway"

func ReadDefs(filename string) (*ir.ConfigDefs, error) {
	config, err := readModel(filename)
	if err != nil {
		return nil, err
	}

	subnets, err1 := parseSubnets(config)
	instances, nifs, err2 := parseInstancesNifs(config)
	vpes, vpeEndpoints, err3 := parseVPEs(config)
	vpcs, err4 := parseVPCs(config)
	err5 := validateVpcs(vpcs)
	if err := errors.Join(err1, err2, err3, err4, err5); err != nil {
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

func parseVPCs(config *configModel.ResourcesContainerModel) (map[ir.ID]*ir.VPCDetails, error) {
	res := make(map[ir.ID]*ir.VPCDetails, len(config.VpcList))
	for _, vpc := range config.VpcList {
		addressPrefixes := netset.NewIPBlock()
		for _, addressPrefix := range vpc.AddressPrefixes {
			address, err := netset.IPBlockFromCidr(*addressPrefix.CIDR)
			if err != nil {
				return nil, err
			}
			addressPrefixes = addressPrefixes.Union(address)
		}
		res[*vpc.Name] = &ir.VPCDetails{AddressPrefixes: addressPrefixes}
	}
	return res, nil
}

func parseSubnets(config *configModel.ResourcesContainerModel) (map[ir.ID]*ir.SubnetDetails, error) {
	subnets := make(map[ir.ID]*ir.SubnetDetails, len(config.SubnetList))
	for _, subnet := range config.SubnetList {
		cidr, err := netset.IPBlockFromCidr(*subnet.Ipv4CIDRBlock)
		if err != nil {
			return nil, err
		}
		subnetDetails := ir.SubnetDetails{
			CIDR: cidr,
		}
		subnets[ScopingString(*subnet.VPC.Name, *subnet.Name)] = &subnetDetails
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
			nifIP, err := netset.IPBlockFromIPAddress(*instance.NetworkInterfaces[i].PrimaryIP.Address)
			if err != nil {
				return nil, nil, err
			}
			nifDetails := ir.NifDetails{
				Instance: ScopingString(*instance.VPC.Name, *instance.Name),
				IP:       nifIP,
				Subnet:   ScopingString(*instance.VPC.Name, *instance.NetworkInterfaces[i].Subnet.Name),
			}
			nifUniqueName := ScopingString(instanceUniqueName, *instance.NetworkInterfaces[i].Name)
			nifs[nifUniqueName] = &nifDetails
			instanceNifs[i] = nifUniqueName
		}
		instanceDetails := ir.InstanceDetails{
			Nifs: instanceNifs,
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
		vpeDetails := ir.VPEDetails{
			VPEReservedIPs: []ir.ID{},
		}
		vpes[ScopingString(*vpe.VPC.Name, *vpe.Name)] = &vpeDetails
	}
	for _, subnet := range config.SubnetList {
		for _, r := range subnet.ReservedIps {
			t, ok := r.Target.(*vpcv1.ReservedIPTarget)
			if !ok || t == nil || r.Address == nil || t.ResourceType == nil || *t.ResourceType != EndpointVPE || t.Name == nil {
				continue
			}
			VPEName := ScopingString(*subnet.VPC.Name, *t.Name)
			subnetName := ScopingString(*subnet.VPC.Name, *subnet.Name)

			vpeIP, err := netset.IPBlockFromIPAddress(*r.Address)
			if err != nil {
				return nil, nil, err
			}
			vpeReservedIPDetails := ir.VPEReservedIPsDetails{
				VPEName: VPEName,
				Subnet:  subnetName,
				IP:      vpeIP,
			}
			uniqueVpeReservedIPName := ScopingString(VPEName, *r.Name)
			vpeReservedIPs[uniqueVpeReservedIPName] = &vpeReservedIPDetails
			vpe := vpes[VPEName]
			vpe.VPEReservedIPs = append(vpe.VPEReservedIPs, uniqueVpeReservedIPName)
			vpes[VPEName] = vpe
		}
	}
	return vpes, vpeReservedIPs, nil
}

func validateVpcs(vpcs map[ir.ID]*ir.VPCDetails) error {
	if vpcs == nil {
		return nil
	}
	for vpcName1, vpcDetails1 := range vpcs {
		for vpcName2, vpcDetails2 := range vpcs {
			if vpcName1 >= vpcName2 {
				continue
			}
			if vpcDetails1.AddressPrefixes.Overlap(vpcDetails2.AddressPrefixes) {
				return fmt.Errorf("vpcs %s and %s have overlapping IP address spaces", vpcName1, vpcName2)
			}
		}
	}
	return nil
}

func ScopingString(s1, s2 string) string {
	return s1 + "/" + s2
}
