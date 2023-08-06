package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"vpc-network-config-synthesis/pkg/synth"
)

func main() {
	filename := os.Args[1]
	spec, err := synth.UnmarshalSpec(filename)
	if err != nil {
		log.Fatalf("Could not parse %s: %s", filename, err)
	}
	dataJSON, err := json.MarshalIndent(spec, "", "    ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(string(dataJSON))
}
