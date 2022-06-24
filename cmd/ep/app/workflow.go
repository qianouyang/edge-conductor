/*
* Copyright (c) 2022 Intel Corporation.
*
* SPDX-License-Identifier: Apache-2.0
*
 */
package app

import (
	epapiplugins "ep/pkg/api/plugins"
	eputils "ep/pkg/eputils"
	docker "ep/pkg/eputils/docker"
	orasutils "ep/pkg/eputils/orasutils"
	plugin "ep/pkg/plugin"
	wf "ep/pkg/workflow"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

var (
	errHost   = errors.New("Host is not found")
	errEKConf = errors.New("EK config file not provided.")
	errTopCfg = fmt.Errorf("Top Config Lost!")
)

func EpWfStart(epParams *epapiplugins.EpParams, name string) error {
	ekcfg := GetRuntimeTopConfig(epParams)
	if ekcfg == nil {
		return errTopCfg
	}
	address := fmt.Sprintf("%s:%s", ekcfg.Parameters.GlobalSettings.ProviderIP, ekcfg.Parameters.GlobalSettings.WorkflowPort)
	plugin.Address = address

	addonbin := "addon/bin/conductor-plugin"
	if eputils.FileExists(addonbin) {

		finished := make(chan bool)
		go func() {
			cmd := exec.Command(addonbin, address)
			_, err := eputils.RunCMDEx(cmd, true)
			if err != nil {
				log.Errorf("Failed to run addon %s on %s", addonbin, address)
			}
			finished <- true
		}()
		err := wf.Start(name, address, WfConfig)
		<-finished
		return err

	} else {
		return wf.Start(name, address, WfConfig)

	}
}

func setHostIptoNoProxy(input_ep_params *epapiplugins.EpParams) error {
	no_proxy := os.Getenv("no_proxy")
	if input_ep_params == nil {
		return errHost
	}

	epHost := input_ep_params.Ekconfig.Parameters.GlobalSettings.ProviderIP
	no_proxy = fmt.Sprintf("%s,%s", no_proxy, epHost)
	if err := os.Setenv("no_proxy", no_proxy); err != nil {
		return err
	}
	return nil
}

func EpUtilsInit(epParams *epapiplugins.EpParams) error {
	regcacert := epParams.Registrycert.Ca.Cert
	auth, err := docker.GetAuthConf(epParams.Ekconfig.Parameters.GlobalSettings.ProviderIP,
		epParams.Ekconfig.Parameters.GlobalSettings.RegistryPort,
		epParams.Ekconfig.Parameters.Customconfig.Registry.User,
		epParams.Ekconfig.Parameters.Customconfig.Registry.Password)
	if err != nil {
		return err
	}

	err = orasutils.OrasNewClient(auth, regcacert)
	if err != nil {
		log.Errorln("Failed to create an OrasClient", err)
		return err
	}
	eputils.SetTemplateParams(epParams)
	eputils.SetTemplateFuncs(funcs)
	return nil
}

func EpWfPreInit(epPms *epapiplugins.EpParams, p map[string]string) (*epapiplugins.EpParams, error) {
	var rfile, ekcfgPath string
	var err error
	epParams := epPms
	if epParams == nil {
		rfile, err = FileNameofRuntime(fnRuntimeInitParams)
		if err != nil {
			log.Errorln("Failed to get runtime file path:", err)
			return nil, err
		}
		if _, err := os.Stat(rfile); os.IsNotExist(err) {
			log.Errorln("Failed to open", rfile, err)
			return nil, err
		}
		epParams = new(epapiplugins.EpParams)
		if err := eputils.LoadSchemaStructFromYamlFile(epParams, rfile); err != nil {
			log.Error(err)
			return nil, err
		}
	}
	epParams.Cmdline = ""
	epParams.Kubeconfig = ""
	for k, v := range p {
		switch k {
		case Epcmdline:
			epParams.Cmdline = v
		case Epkubeconfig:
			epParams.Kubeconfig = v
		case EKConfigPath:
			ekcfgPath = v
		default:
			log.Warnf("%s is not defined in Epparams", k)
		}
	}
	if ekcfgPath == "" {
		if epParams.Ekconfigpath != "" {
			ekcfgPath = epParams.Ekconfigpath
		} else {
			return nil, errEKConf
		}
	}
	err = setupCustomConfig(ekcfgPath, epParams)
	if err != nil {
		log.Errorln("Failed to setup custom config:", err)
		return nil, err
	}
	if err = setHostIptoNoProxy(epParams); err != nil {
		log.Errorln("Failed to set HostIp to no proxy env", err)
		return nil, err
	}

	if err = EpUtilsInit(epParams); err != nil {
		return nil, err
	}

	return epParams, nil
}

func EpWfTearDown(epParams *epapiplugins.EpParams, rfile string) error {
	teardownCustomConfig(epParams)
	err := eputils.SaveSchemaStructToYamlFile(epParams, rfile)
	if err != nil {
		log.Errorln("Failed to save ep-params:", err)
		return err
	}
	return nil
}

func EpwfLoadServices(epParams *epapiplugins.EpParams) error {
	EkcfgComponentsSelector := &(epParams.Ekconfig.Components.Selector)
	index := len(*EkcfgComponentsSelector)
	*EkcfgComponentsSelector = append((*EkcfgComponentsSelector)[:0], (*EkcfgComponentsSelector)[index:]...)

	err := load_ek_services(epParams.Ekconfig, epParams.Ekconfigpath)
	if err != nil {
		return err
	}

	if EkcfgComponentsSelector == nil {
		log.Warnf("Components Selector not specified, use default value %s", DefaultComponentsSelector)
	}
	return nil
}
