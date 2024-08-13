/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package confio

import (
	"encoding/json"
	"os"

	configModel "github.com/np-guard/cloud-resource-collector/pkg/ibm/datamodel"
)

func readModel(filename string) (*configModel.ResourcesContainerModel, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	model := configModel.ResourcesContainerModel{}
	err = json.Unmarshal(bytes, &model)
	if err != nil {
		return nil, err
	}
	return &model, nil
}
