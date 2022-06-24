/*
* Copyright (c) 2021 Intel Corporation.
*
* SPDX-License-Identifier: Apache-2.0
*
 */

package app

import (
	epapiplugins "ep/pkg/api/plugins"
	"ep/pkg/eputils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"regexp"
	"sigs.k8s.io/yaml"
	"strings"
)

var (
	errYml           = fmt.Errorf("Not a valid yaml file.")
	errSelector      = fmt.Errorf("Cannot find name in selector")
	errCompSelect    = fmt.Errorf("Component selector error, file")
	errNameOverride  = fmt.Errorf("override err")
	errMergeOverride = fmt.Errorf("override merge err")
	errUnmarshalOver = fmt.Errorf("override unmarshal to component error")
	errYmlOverride   = fmt.Errorf("override yaml unmarshal error")
	errSSHPath       = fmt.Errorf("Choose ESP for OS provider. There is no SSH path for ESP provision")
)

type EKBaseConfig struct {
	Use        []string `yaml:"Use"`
	ParamsMap  map[string]interface{}
	File       string
	Parameters interface{}
}

type EKParams struct {
	Ekconfig  *EKBaseConfig
	Workspace string
}

func get_item_yaml(name string, yml string) (string, error) {
	lines := strings.Split(yml, "\n")
	if len(lines) <= 0 {
		return "", errYml
	}
	matched := false
	y := ""
	rs, _ := regexp.Compile(fmt.Sprintf("^%s:", name))
	re, _ := regexp.Compile("^[#\\- \t]")
	for _, l := range lines {
		if match := re.MatchString(l); len(l) > 0 && (!match) && matched {
			break
		}
		if match := rs.MatchString(l); match {
			matched = true
		}
		if matched {
			y = fmt.Sprintf("%s\n%s", y, l)
		}
	}
	return y, nil
}

func load_ek_bcfg(file string) (*EKBaseConfig, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	bcfg := EKBaseConfig{}
	useyml, err := get_item_yaml("Use", string(data))
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(useyml), &bcfg)
	if err != nil {
		return nil, err
	}
	paramsyml, err := get_item_yaml("Parameters", string(data))
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(paramsyml), &bcfg.ParamsMap)
	if err != nil {
		return nil, err
	}
	bcfg.File = file
	return &bcfg, nil
}

func load_ek_bcfg_recursively(file string) ([]*EKBaseConfig, error) {
	bcfgs := []*EKBaseConfig{}
	bcfg, err := load_ek_bcfg(file)
	if err != nil {
		return nil, err
	}
	for _, uf := range bcfg.Use {
		bcs, err := load_ek_bcfg_recursively(uf)
		if err != nil {
			return nil, err
		}
		bcfgs = append(bcfgs, bcs...)
	}
	bcfgs = append(bcfgs, bcfg)
	return bcfgs, nil
}

func merge_component(ekconfig *epapiplugins.Ekconfig, compMap map[string]interface{}) error {
	name, ok := compMap["name"].(string)
	if !ok {
		return errSelector
	}
	overrideyaml := ""
	override, ok := compMap["override"].(map[string]interface{})
	if ok {
		data, err := yaml.Marshal(override)
		if err != nil {
			log.Warningf("%v override err: %v", name, err)
			return errNameOverride
		}
		overrideyaml = string(data)
	} else {
		log.Debugf("Cannot find override in %v", name)
	}
	found := false
	for _, v := range ekconfig.Components.Selector {
		if name == v.Name {
			found = true
			if len(v.OverrideYaml) > 0 {
				m := map[string]interface{}{}
				err := yaml.Unmarshal([]byte(v.OverrideYaml), &m)
				if err != nil {
					log.Warningf("%v override yaml unmarshal error: %v", name, err)
					return errYmlOverride
				}
				mm := eputils.MergeMaps(m, override)
				data, err := yaml.Marshal(mm)
				if err != nil {
					log.Warningf("%v override merge err: %v", name, err)
					return errMergeOverride
				}
				overrideyaml = string(data)
			}
			v.OverrideYaml = overrideyaml
			break
		}
	}
	if !found {
		ekconfig.Components.Selector = append(ekconfig.Components.Selector,
			&epapiplugins.EkconfigComponentsSelectorItems0{
				Name:         name,
				OverrideYaml: overrideyaml,
			})
	}
	comp := epapiplugins.Component{}
	err := yaml.Unmarshal([]byte(overrideyaml), &comp)
	if err != nil {
		log.Warningf("%v override unmarshal to component error: %v", name, err)
		return errUnmarshalOver
	}
	return nil
}

func load_ek_config(ekconfig *epapiplugins.Ekconfig, ekfile string) error {
	bcfgs, err := load_ek_bcfg_recursively(ekfile)
	if err != nil {
		return err
	}

	ekparams := &EKParams{
		Workspace: GetWorkspacePath(),
		Ekconfig:  &EKBaseConfig{},
	}

	params := map[string]interface{}{}
	for _, bcfg := range bcfgs {
		params = eputils.MergeMaps(params, bcfg.ParamsMap)
	}
	if err != nil {
		return err
	}
	ekparams.Ekconfig.Parameters = params["Parameters"]
	*ekconfig = epapiplugins.Ekconfig{
		Parameters: &epapiplugins.EkconfigParameters{
			GlobalSettings: &epapiplugins.EkconfigParametersGlobalSettings{},
		},
		OS:         &epapiplugins.EkconfigOS{},
		Cluster:    &epapiplugins.EkconfigCluster{},
		Components: &epapiplugins.EkconfigComponents{},
	}
	err = eputils.ConvertSchemaStruct(ekparams.Ekconfig.Parameters, ekconfig.Parameters)
	if err != nil {
		return err
	}
	if ekconfig.Parameters == nil {
		ekconfig.Parameters = &epapiplugins.EkconfigParameters{}
	}
	if ekconfig.Parameters.GlobalSettings == nil {
		ekconfig.Parameters.GlobalSettings = &epapiplugins.EkconfigParametersGlobalSettings{}
	}
	if ekconfig.Parameters.GlobalSettings.RegistryPort == "" {
		ekconfig.Parameters.GlobalSettings.RegistryPort = DefaultRegistryPort
	}
	if ekconfig.Parameters.GlobalSettings.ProviderIP == "" {
		ekconfig.Parameters.GlobalSettings.ProviderIP = GetHostDefaultIP()
	}
	if ekconfig.Parameters.GlobalSettings.WorkflowPort == "" {
		ekconfig.Parameters.GlobalSettings.WorkflowPort = DefaultWfPort
	}

	for _, bcfg := range bcfgs {
		log.Debugf("load bcfg from: %v", bcfg.File)
		data, err := ioutil.ReadFile(bcfg.File)
		if err != nil {
			return err
		}
		yml, err := eputils.StringTemplateConvertWithParams(string(data), ekparams)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}

		cfg := epapiplugins.Ekconfig{}
		json, err := yaml.YAMLToJSON([]byte(yml))
		if err != nil {
			return err
		}
		err = cfg.UnmarshalBinary(json)
		if err != nil {
			return err
		}
		if cfg.OS != nil && len(cfg.OS.Provider) > 0 {
			ekconfig.OS = cfg.OS
		}
		if cfg.Cluster != nil && len(cfg.Cluster.Provider) > 0 {
			ekconfig.Cluster = cfg.Cluster
		}

		if cfg.Components == nil {
			cfg.Components = &epapiplugins.EkconfigComponents{}
		}
		for _, m := range cfg.Components.Manifests {
			found := false
			for _, mm := range ekconfig.Components.Manifests {
				if m == mm {
					found = true
					break
				}
			}
			if !found {
				ekconfig.Components.Manifests = append(ekconfig.Components.Manifests, m)
			}
		}
		if err := load_ek_services(ekconfig, ekfile); err != nil {
			return err
		}
	}
	if ekconfig.Cluster == nil {
		ekconfig.Cluster = &epapiplugins.EkconfigCluster{}
	}
	if ekconfig.Cluster.Provider == "" {
		ekconfig.Cluster.Provider = DefaultClusterProvider
	}
	if ekconfig.Cluster.Provider != DefaultClusterProvider && ekconfig.Cluster.Config == "" {
		ekconfig.Cluster.Config = DefaultClusterConfig
	}
	if ekconfig.OS == nil {
		ekconfig.OS = &epapiplugins.EkconfigOS{}
	}
	if ekconfig.OS.Provider == "" {
		ekconfig.OS.Provider = DefaultOSProvider
	}
	if ekconfig.OS.Config == "" {
		ekconfig.OS.Config = DefaultOSConfig
	}
	if ekconfig.Parameters.DefaultSSHKeyPath == "" && ekconfig.OS.Provider == "esp" {
		return errSSHPath
	}
	eputils.DumpVar(ekconfig)
	if ekconfig.Validate(nil) != nil {
		log.Warningf("Verify ekconfig err: %v", err)
		return errEkconfig
	}

	return nil
}

func load_ek_services(ekconfig *epapiplugins.Ekconfig, ekfile string) error {
	bcfgs, err := load_ek_bcfg_recursively(ekfile)
	if err != nil {
		return err
	}
	ekparams := &EKParams{
		Workspace: GetWorkspacePath(),
		Ekconfig:  &EKBaseConfig{},
	}
	for _, bcfg := range bcfgs {
		data, err := ioutil.ReadFile(bcfg.File)
		if err != nil {
			return err
		}
		yml, err := eputils.StringTemplateConvertWithParams(string(data), ekparams)
		if err != nil {
			return err
		}
		componentYml, err := get_item_yaml("Components", yml)
		if err != nil {
			return err
		}
		componentMap := map[string]interface{}{}
		err = yaml.Unmarshal([]byte(componentYml), &componentMap)
		if err != nil {
			return err
		}

		if _, ok := componentMap["Components"]; ok {
			for k, v := range componentMap["Components"].(map[string]interface{}) {
				if k != "selector" {
					continue
				}
				if _, ok := v.([]interface{}); !ok {
					log.Warningf("Component selector error, file: %v", bcfg.File)
					return errCompSelect
				}
				for _, sv := range v.([]interface{}) {
					if _, ok := sv.(map[string]interface{}); !ok {
						log.Warningf("Component selector error, file: %v", bcfg.File)
						return errCompSelect
					}
					err := merge_component(ekconfig, sv.(map[string]interface{}))
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
