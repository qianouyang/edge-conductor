/*
* Copyright (c) 2021 Intel Corporation.
*
* SPDX-License-Identifier: Apache-2.0
*
 */

package conductorutils

import (
	papi "ep/pkg/api/plugins"
	log "github.com/sirupsen/logrus"
)

func GetResourceValueFromCustomcfg(cfg *papi.Customconfig, name string) (string, error) {
	resources := cfg.Resources
	for _, r := range resources {
		if r.Name == name {
			return r.Value, nil
		}
	}
	log.Errorf("Resource %s not found.", name)
	return "", errResource
}
