package confio

import (
	"fmt"
	"log"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	configModel "github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

func makeACLRuleItem(rule *ir.ACLRule, current,
	next *vpcv1.NetworkACLRuleReference) vpcv1.NetworkACLRuleItemIntf {
	iPVersion := utils.Ptr(ipv4Const)
	direction := direction(rule.Direction)
	action := action(rule.Action)
	source := ip(rule.Source)
	destination := ip(rule.Destination)
	switch p := rule.Protocol.(type) {
	case ir.TCPUDP:
		data := tcpudp(p)
		result := &vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolTcpudp{
			Href:        current.Href,
			ID:          current.ID,
			Name:        current.Name,
			IPVersion:   iPVersion,
			Action:      action,
			Direction:   direction,
			Source:      source,
			Destination: destination,
			Before:      next,

			Protocol:           data.Protocol,
			SourcePortMin:      data.SourcePortMin,
			SourcePortMax:      data.SourcePortMax,
			DestinationPortMin: data.DestinationPortMin,
			DestinationPortMax: data.DestinationPortMax,
		}
		return result
	case ir.ICMP:
		data := icmp(p)
		result := &vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmp{
			Href:        current.Href,
			ID:          current.ID,
			Name:        current.Name,
			IPVersion:   iPVersion,
			Action:      action,
			Direction:   direction,
			Source:      source,
			Destination: destination,
			Before:      next,

			Protocol: data.Protocol,
			Type:     data.Type,
			Code:     data.Code,
		}
		return result
	case ir.AnyProtocol:
		data := all()
		result := &vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAll{
			Href:        current.Href,
			ID:          current.ID,
			Name:        current.Name,
			IPVersion:   iPVersion,
			Action:      action,
			Direction:   direction,
			Source:      source,
			Destination: destination,
			Before:      next,

			Protocol: data.Protocol,
		}
		return result
	default:
		log.Fatalf("Impossible protocol type for acl rule %v: %T", current.Name, p)
		return nil
	}
}

func makeACLItem(acl *ir.ACL, subnet *vpcv1.SubnetReference) *configModel.NetworkACL {
	ruleItems := make([]vpcv1.NetworkACLRuleItemIntf, len(acl.Rules()))
	rules := acl.Rules()

	var next *vpcv1.NetworkACLRuleReference
	for i := len(ruleItems) - 1; i >= 0; i-- {
		name := utils.Ptr(fmt.Sprintf("rule%v", i))
		ref := allocateRef()
		current := &vpcv1.NetworkACLRuleReference{
			Name: name,
			Href: ref.Href,
			ID:   ref.ID,
		}
		ruleItems[i] = makeACLRuleItem(&rules[i], current, next)
		next = current
	}

	ref := allocateRef()
	result := configModel.NewNetworkACL(&vpcv1.NetworkACL{
		ID:      ref.ID,
		CRN:     ref.CRN,
		Href:    ref.Href,
		Name:    utils.Ptr(acl.Name()),
		Subnets: []vpcv1.SubnetReference{*subnet},
		Rules:   ruleItems,
	})
	result.Tags = []string{}
	return result
}

func findSubnet(model *configModel.ResourcesContainerModel, name string) int {
	for i, subnet := range model.SubnetList {
		if subnet.Name != nil && *subnet.Name == name {
			return i
		}
	}
	return -1
}

func updateACL(model *configModel.ResourcesContainerModel, collection *ir.ACLCollection) error {
	vpc := model.SubnetList[0].VPC
	resourceGroup := model.SubnetList[0].ResourceGroup
	for _, subnetName := range collection.SortedACLSubnets() {
		acl := collection.ACLs[subnetName]
		if acl == nil {
			continue
		}
		subnetIndex := findSubnet(model, subnetName)
		subnet := model.SubnetList[subnetIndex]
		subnetRef := &vpcv1.SubnetReference{
			Name:         subnet.Name,
			CRN:          subnet.CRN,
			Href:         subnet.Href,
			ID:           subnet.ID,
			ResourceType: subnet.ResourceType,
		}
		aclItem := makeACLItem(acl, subnetRef)
		aclItem.ResourceGroup = resourceGroup
		aclItem.VPC = vpc
		model.NetworkACLList = append(model.NetworkACLList, aclItem)
		model.SubnetList[subnetIndex].NetworkACL = &vpcv1.NetworkACLReference{
			ID:   aclItem.ID,
			CRN:  aclItem.CRN,
			Href: aclItem.Href,
			Name: aclItem.Name,
		}
	}
	return nil
}

func (w *Writer) WriteACL(collection *ir.ACLCollection) error {
	if err := updateACL(w.model, collection); err != nil {
		return err
	}
	return w.writeModel()
}
