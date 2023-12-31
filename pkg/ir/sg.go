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
	SGResourceFileShareMountTarget SGResource = "fsmt"
)

type SGName string

func (s SGName) String() string {
	return string(s)
}

type RemoteType interface {
	fmt.Stringer
	// IP | CIDR | SGName
}

type SGRule struct {
	Direction   Direction
	Remote      RemoteType
	Protocol    Protocol
	Explanation string
}

type SG struct {
	Rules    []SGRule
	Attached []SGName
}

type SGCollection struct {
	SGs map[SGName]*SG
}

type SGWriter interface {
	WriteSG(*SGCollection) error
}

func (r *SGRule) isRedundant(rules []SGRule) bool {
	for i := range rules {
		if rules[i].mustSupersede(r) {
			return true
		}
	}
	return false
}

func (r *SGRule) mustSupersede(other *SGRule) bool {
	otherExplanation := other.Explanation
	other.Explanation = r.Explanation
	res := reflect.DeepEqual(r, other)
	other.Explanation = otherExplanation
	return res
}

func NewSG() *SG {
	return &SG{Rules: []SGRule{}, Attached: []SGName{}}
}

func NewSGCollection() *SGCollection {
	return &SGCollection{SGs: map[SGName]*SG{}}
}

func (c *SGCollection) LookupOrCreate(name SGName) *SG {
	sg, ok := c.SGs[name]
	if ok {
		return sg
	}
	newSG := NewSG()
	c.SGs[name] = newSG
	return newSG
}

func (a *SG) Add(rule *SGRule) {
	if rule.isRedundant(a.Rules) {
		return
	}
	a.Rules = append(a.Rules, *rule)
}

func MergeSGCollections(collections ...*SGCollection) *SGCollection {
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

func (c *SGCollection) Write(w Writer) error {
	return w.WriteSG(c)
}

func (c *SGCollection) SortedSGNames() []SGName {
	return utils.SortedKeys(c.SGs)
}
