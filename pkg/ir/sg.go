package ir

import (
	"fmt"
	"reflect"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

type SGResource string

const (
	SGResourceNIF                  SGResource = "nif"
	SGResourceBareMetalNIF         SGResource = "bnif"
	SGResourceLoadBalancer         SGResource = "loadbalancer"
	SGResourceVPE                  SGResource = "vpe"
	SGResourceVPNServer            SGResource = "vpn"
	SGResourceFileShareMountTarget SGResource = "fsme"
)

type SecurityGroupName string

func (s SecurityGroupName) String() string {
	return string(s)
}

type RemoteType interface {
	fmt.Stringer
	IP | CIDR | SecurityGroupName
}

type SecurityGroupRule[T RemoteType] struct {
	Direction   Direction
	Remote      T
	Protocol    Protocol
	Explanation string
}

type SecurityGroup struct {
	Rules    []SecurityGroupRule[CIDR]
	Attached []SecurityGroupName
}

type SecurityGroupCollection struct {
	SGs map[SecurityGroupName]*SecurityGroup
}

type SGWriter interface {
	WriteSG(*SecurityGroupCollection) error
}

func (r *SecurityGroupRule[CIDR]) isRedundant(rules []SecurityGroupRule[CIDR]) bool {
	for i := range rules {
		if rules[i].mustSupersede(r) {
			return true
		}
	}
	return false
}

func (r *SecurityGroupRule[CIDR]) mustSupersede(other *SecurityGroupRule[CIDR]) bool {
	otherExplanation := other.Explanation
	other.Explanation = r.Explanation
	res := reflect.DeepEqual(r, other)
	other.Explanation = otherExplanation
	return res
}

func NewSecurityGroup() *SecurityGroup {
	return &SecurityGroup{Rules: []SecurityGroupRule[CIDR]{}, Attached: []SecurityGroupName{}}
}

func NewSGCollection() *SecurityGroupCollection {
	return &SecurityGroupCollection{SGs: map[SecurityGroupName]*SecurityGroup{}}
}

func (c *SecurityGroupCollection) LookupOrCreate(name SecurityGroupName) *SecurityGroup {
	acl, ok := c.SGs[name]
	if ok {
		return acl
	}
	newSG := NewSecurityGroup()
	c.SGs[name] = newSG
	return newSG
}

func (a *SecurityGroup) Add(rule *SecurityGroupRule[CIDR]) {
	if rule.isRedundant(a.Rules) {
		return
	}
	a.Rules = append(a.Rules, *rule)
}

func MergeSecurityGroupCollections(collections ...*SecurityGroupCollection) *SecurityGroupCollection {
	result := NewSGCollection()
	for _, c := range collections {
		for a := range c.SGs {
			sg := c.LookupOrCreate(a)
			for r := range sg.Rules {
				result.LookupOrCreate(a).Add(&sg.Rules[r])
			}
		}
	}
	return result
}

func (c *SecurityGroupCollection) SortedSGNames() []SecurityGroupName {
	return utils.SortedKeys(c.SGs)
}
