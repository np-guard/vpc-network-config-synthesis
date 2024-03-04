package ir

import (
	"fmt"
	"log"
	"math"
)

type ICMPCodeType struct {
	// ICMP type allowed.
	Type int

	// ICMP code allowed. If omitted, any code is allowed
	Code *int
}

type ICMP struct {
	*ICMPCodeType
}

func (t ICMP) InverseDirection() Protocol {
	if t.ICMPCodeType == nil {
		return nil
	}

	if invType := inverseICMPType(t.Type); invType != undefinedICMP {
		return ICMP{ICMPCodeType: &ICMPCodeType{Type: invType, Code: t.Code}}
	}
	return nil
}

// Based on https://datatracker.ietf.org/doc/html/rfc792

const (
	echoReply              = 0
	destinationUnreachable = 3
	sourceQuench           = 4
	redirect               = 5
	echo                   = 8
	timeExceeded           = 11
	parameterProblem       = 12
	timestamp              = 13
	timestampReply         = 14
	informationRequest     = 15
	informationReply       = 16

	undefinedICMP = -2
)

// inverseICMPType returns the reply type for request type and vice versa.
// When there is no inverse, returns undefinedICMP
func inverseICMPType(t int) int {
	switch t {
	case echo:
		return echoReply
	case echoReply:
		return echo

	case timestamp:
		return timestampReply
	case timestampReply:
		return timestamp

	case informationRequest:
		return informationReply
	case informationReply:
		return informationRequest

	case destinationUnreachable, sourceQuench, redirect, timeExceeded, parameterProblem:
		return undefinedICMP
	default:
		log.Panicf("Impossible ICMP type: %v", t)
	}
	return undefinedICMP
}

//nolint:revive // magic numbers are fine here
func ValidateICMP(t, c int) error {
	possibleCodes := map[int]int{
		echoReply:              0,
		destinationUnreachable: 5,
		sourceQuench:           0,
		redirect:               3,
		echo:                   0,
		timeExceeded:           1,
		parameterProblem:       0,
		timestamp:              0,
		timestampReply:         0,
		informationRequest:     0,
		informationReply:       0,
	}
	max, ok := possibleCodes[t]
	if !ok {
		return fmt.Errorf("invalid ICMP type %v", t)
	}
	if c > max {
		return fmt.Errorf("ICMP code %v is invalid for ICMP type %v", c, t)
	}
	return nil
}

const (
	newDestinationUnreachable = 0
	newRedirect               = 6
	newTimeExceeded           = 10
	newEcho                   = 17
	newEchoReply              = 18
	newSourceQuench           = 19
)

func mapToNew(t, code int) int {
	switch t {
	case destinationUnreachable:
		return newDestinationUnreachable + code
	case redirect:
		return newRedirect + code
	case timeExceeded:
		return newTimeExceeded + code
	case echo:
		return newEcho
	case echoReply:
		return newEchoReply
	case sourceQuench:
		return newSourceQuench
	default:
		return t
	}
}

func mapToOld(newCode int) (t int, code int) {
	switch {
	case newCode < newRedirect:
		t = newDestinationUnreachable
	case newCode < newTimeExceeded:
		t = newRedirect
	case newCode < parameterProblem:
		t = newTimeExceeded
	case newCode == newEcho:
		t = echo
	case newCode == newEchoReply:
		t = echoReply
	case newCode == newSourceQuench:
		t = sourceQuench
	default:
		t = newCode
	}
	code = newCode - t
	return
}

type ICMPSet uint32

func (s ICMPSet) IsSubset(other ICMPSet) bool {
	return s|other == other
}

func (s ICMPSet) Union(other ICMPSet) ICMPSet {
	return s | other
}

const (
	allDestinationUnreachable = 0b00000000000000111111
	allRedirect               = 0b00000000001111000000
	allTimeExceeded           = 0b00000000110000000000
	allOther                  = 0b11111111000000000000
)

func FromICMP(t ICMP) ICMPSet {
	if t.ICMPCodeType == nil {
		return allDestinationUnreachable | allRedirect | allTimeExceeded | allOther
	}
	if t.Code == nil {
		return math.MaxUint32
	}
	return 1 << mapToNew(t.Type, *t.Code)
}
