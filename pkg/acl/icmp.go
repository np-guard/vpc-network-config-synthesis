package acl

import "log"

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
