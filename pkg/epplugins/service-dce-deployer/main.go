/*
* Copyright (c) 2022 Intel Corporation.
*
* SPDX-License-Identifier: Apache-2.0
*
 */

// Template auto-generated once, maintained by plugin owner.

package servicedcedeployer

import (
	eputils "ep/pkg/eputils"
	"ep/pkg/executor"
)

func PluginMain(in eputils.SchemaMapData, outp *eputils.SchemaMapData) error {
	input_ep_params := input_ep_params(in)
	input_serviceconfig := input_serviceconfig(in)

	for _, service := range input_serviceconfig.Components {
		if service.Type == "dce" && service.Executor != nil {
			if service.Executor.Deploy != "" {
				err := executor.Run(service.Executor.Deploy, input_ep_params, service)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
