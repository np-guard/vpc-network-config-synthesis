package spec

import (
	"fmt"
	"log"
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

func (t ICMP) Name() string { return "ICMP" }

func (t ICMP) InverseDirection() Protocol {
	if t.ICMPCodeType == nil {
		return nil
	}

	if invType := inverseICMPType(t.Type); invType != undefinedICMP {
		return &ICMP{ICMPCodeType: &ICMPCodeType{Type: invType, Code: t.Code}}
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
// When there is no inverse, returns `undefinedICMP`
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
		log.Fatalf("Impossible ICMP type: %v", t)
	}
	return undefinedICMP
}

// FIX: use go21 slices.Contains
func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//nolint:revive // magic numbers are fine here
func ValidateICMP(t, c int) error {
	possibleCodes := map[int][]int{
		echoReply:              {0},
		destinationUnreachable: {0, 1, 2, 3, 4, 5},
		sourceQuench:           {0},
		redirect:               {0, 1, 2, 3},
		echo:                   {0},
		timeExceeded:           {0, 1},
		parameterProblem:       {0},
		timestamp:              {13, 14},
		timestampReply:         {13, 14},
		informationRequest:     {15, 16},
		informationReply:       {15, 16},
	}
	options, ok := possibleCodes[t]
	if !ok {
		return fmt.Errorf("invalid ICMP type %v", t)
	}
	if !contains(options, c) {
		return fmt.Errorf("ICMP code %v is invalid for ICMP type %v", c, t)
	}
	return nil
}
