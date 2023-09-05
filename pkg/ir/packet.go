package ir

import (
	"fmt"
	"sort"
)

// Helpers for creation of ACLs

func (a *ACL) Rules() []Rule {
	rules := a.Internal
	if len(a.External) != 0 {
		rules = append(rules, makeDenyInternal()...)
		rules = append(rules, a.External...)
	}
	return rules
}

func (a *ACL) AppendInternal(rule *Rule) {
	if a.External == nil {
		panic("ACLs should be created with non-null Internal")
	}
	a.Internal = append(a.Internal, *rule)
}

func (a *ACL) Name() string {
	return fmt.Sprintf("acl-%v", a.Subnet)
}

func (a *ACL) AppendExternal(rule *Rule) {
	if a.External == nil {
		panic("ACLs should be created with non-null External")
	}
	a.External = append(a.External, *rule)
}

func NewCollection() *Collection {
	return &Collection{ACLs: map[string]*ACL{}}
}

func MergeCollections(collections ...*Collection) *Collection {
	result := NewCollection()
	for _, c := range collections {
		for a := range c.ACLs {
			acl := c.ACLs[a]
			for r := range acl.Internal {
				result.ACLs[a].AppendInternal(&acl.Internal[r])
			}
			for r := range acl.External {
				result.ACLs[a].AppendExternal(&acl.External[r])
			}
		}
	}
	return result
}

func (c *Collection) LookupOrCreate(name string) *ACL {
	acl, ok := c.ACLs[name]
	if ok {
		return acl
	}
	newACL := ACL{Subnet: name, Internal: []Rule{}, External: []Rule{}}
	c.ACLs[name] = &newACL
	return &newACL
}

func (c *Collection) SortedACLSubnets() []string {
	keys := make([]string, len(c.ACLs))
	i := 0
	for k := range c.ACLs {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

type Packet struct {
	Src, Dst    IP
	Protocol    Protocol
	Explanation string
}

func (r *Rule) Target() IP {
	if r.Direction == Inbound {
		return r.Destination
	}
	return r.Source
}

func AllowSend(packet Packet) *Rule {
	return packetRule(packet, Outbound, Allow)
}

func AllowReceive(packet Packet) *Rule {
	return packetRule(packet, Inbound, Allow)
}

func packetRule(packet Packet, direction Direction, action Action) *Rule {
	return &Rule{
		Action:      action,
		Source:      packet.Src,
		Destination: packet.Dst,
		Direction:   direction,
		Protocol:    packet.Protocol,
		Explanation: packet.Explanation,
	}
}

// makeDenyInternal prevents allowing external communications from accidentally allowing internal communications too
func makeDenyInternal() []Rule {
	localIPs := []IP{ // https://datatracker.ietf.org/doc/html/rfc1918#section-3
		{"10.0.0.0/8"},
		{"172.16.0.0/12"},
		{"192.168.0.0/16"},
	}
	var denyInternal []Rule
	for i, anyLocalIPSrc := range localIPs {
		for j, anyLocalIPDst := range localIPs {
			explanation := fmt.Sprintf("Deny other internal communication; see rfc1918#3; item %v,%v", i, j)
			denyInternal = append(denyInternal, []Rule{
				*packetRule(Packet{Src: anyLocalIPSrc, Dst: anyLocalIPDst, Protocol: AnyProtocol{}, Explanation: explanation}, Outbound, Deny),
				*packetRule(Packet{Src: anyLocalIPDst, Dst: anyLocalIPSrc, Protocol: AnyProtocol{}, Explanation: explanation}, Inbound, Deny),
			}...)
		}
	}
	return denyInternal
}
