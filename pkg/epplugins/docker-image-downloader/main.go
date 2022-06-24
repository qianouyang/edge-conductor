/*
* Copyright (c) 2022 Intel Corporation.
*
* SPDX-License-Identifier: Apache-2.0
*
 */
// Template auto-generated once, maintained by plugin owner.

package dockerimagedownloader

import (
	eputils "ep/pkg/eputils"
	docker "ep/pkg/eputils/docker"
	restfulcli "ep/pkg/eputils/restfulcli"
	"errors"
	log "github.com/sirupsen/logrus"
)

var (
	DayZeroCertFilePath = "cert/pki/ca.pem"
)

var (
	errEpparamsEkconfig = errors.New("epparams Ekconfig is not correct")
)

func PluginMain(in eputils.SchemaMapData, outp *eputils.SchemaMapData) error {
	input_ep_params := input_ep_params(in)
	input_docker_images := input_docker_images(in)

	if input_docker_images.Images == nil || len(input_docker_images.Images) == 0 {
		return nil
	}

	if input_ep_params.Ekconfig == nil || input_ep_params.Ekconfig.Parameters == nil || input_ep_params.Ekconfig.Parameters.GlobalSettings == nil || input_ep_params.Ekconfig.Parameters.Customconfig == nil {
		return errEpparamsEkconfig
	}

	auth, err := docker.GetAuthConf(input_ep_params.Ekconfig.Parameters.GlobalSettings.ProviderIP,
		input_ep_params.Ekconfig.Parameters.GlobalSettings.RegistryPort,
		input_ep_params.Ekconfig.Parameters.Customconfig.Registry.User,
		input_ep_params.Ekconfig.Parameters.Customconfig.Registry.Password)
	if err != nil {
		return err
	}

	var images []string
	var newImages []string

	for _, img := range input_docker_images.Images {
		url := img.URL
		images = append(images, url)
		log.Infof("Pull image %s", url)
		if err := docker.ImagePull(url, nil); err != nil {
			return err
		}
	}

	if newImages, err = restfulcli.MapImageURLCreateHarborProject(input_ep_params.Ekconfig.Parameters.GlobalSettings.ProviderIP,
		input_ep_params.Ekconfig.Parameters.GlobalSettings.RegistryPort,
		input_ep_params.Ekconfig.Parameters.Customconfig.Registry.User,
		input_ep_params.Ekconfig.Parameters.Customconfig.Registry.Password, images); err != nil {
		return err
	}

	for _, img := range newImages {
		prefixUrl := img

		newTag, err := docker.TagImageToLocal(prefixUrl, auth.ServerAddress)
		if err != nil {
			return err
		}
		log.Infof("Push %s to %s", prefixUrl, newTag)
		if err := docker.ImagePush(newTag, auth); err != nil {
			return err
		}

	}

	return nil
}
