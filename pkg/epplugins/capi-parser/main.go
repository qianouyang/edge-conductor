/*
* Copyright (c) 2021 Intel Corporation.
*
* SPDX-License-Identifier: Apache-2.0
*
 */

// Template auto-generated once, maintained by plugin owner.

package capiparser

import (
	pluginapi "ep/pkg/api/plugins"
	eputils "ep/pkg/eputils"
	capiutils "ep/pkg/eputils/capiutils"
	cutils "ep/pkg/eputils/conductorutils"
	"errors"
	log "github.com/sirupsen/logrus"
	"path/filepath"
)

var (
	errAppendFile   = errors.New("Fail to append files")
	errCertMgrCfg   = errors.New("Cert manager config is missing in manifest.")
	errParam        = errors.New("Incorrect parameter")
	errCAPIProvider = errors.New("Failed to get CAPI infra provider config!")
	errCAPIManifest = errors.New("CAPI manifest Lost!")
	errClusterProv  = errors.New("Failed to get Cluster providers yaml list.")
	errCAPIKindLost = errors.New("CAPI manifest kind info Lost!")
	errMgmtCluster  = errors.New("Failed to get management cluster binary list")
)

func getManagementClusterBinaryList(kindprovider *pluginapi.ClustermanifestClusterProvidersItems0, images *pluginapi.Images, files *pluginapi.Files) error {
	kindImage, err := cutils.GetImageFromProvider(kindprovider, "img_node")
	if err != nil {
		log.Errorln("Failed to find container image for KIND node.")
		return err
	}
	kindHAProxy, err := cutils.GetImageFromProvider(kindprovider, "img_haproxy")
	if err != nil {
		log.Errorln("Failed to find container image for KIND haproxy.")
		return err
	}
	kindBin, kindSHA256, err := cutils.GetBinaryFromProvider(kindprovider, "kindtool")
	if err != nil {
		log.Errorln("Failed to find binary for KIND.")
		return err
	}

	images.Images = append(images.Images, &pluginapi.ImagesItems0{Name: "kind", URL: kindImage})
	images.Images = append(images.Images, &pluginapi.ImagesItems0{Name: "kindhaproxy", URL: kindHAProxy})

	files.Files = append(files.Files, &pluginapi.FilesItems0{
		URL:      kindBin,
		Hash:     kindSHA256,
		Hashtype: "sha256",
		Urlreplacement: &pluginapi.FilesItems0Urlreplacement{
			New:    "capi/kind",
			Origin: eputils.GetBaseUrl(kindBin),
		},
	})

	if files.Files == nil {
		return errAppendFile
	}

	return nil
}

func getDockerImagesList(manifest *pluginapi.ClustermanifestClusterapi, infra_provider capiutils.CapiInfraProvider, images *pluginapi.Images) {
	config_name := capiutils.GetManifestConfigNameByCapiInfraProvider(infra_provider)

	images.Images = []*pluginapi.ImagesItems0{}
	for _, config := range manifest.Configs {
		if config.Name == config_name {
			for _, image := range config.Images {
				images.Images = append(images.Images, &pluginapi.ImagesItems0{Name: "", URL: image})
			}
		}
	}
}

func generateFileItemsByURL(url string, subpath string) *pluginapi.FilesItems0 {
	dir := eputils.GetBaseUrl(url)
	subRef := filepath.Join("capi", subpath)
	return &pluginapi.FilesItems0{
		URL: url,
		Urlreplacement: &pluginapi.FilesItems0Urlreplacement{
			New:    subRef,
			Origin: dir,
		},
	}
}

func getProviderYamlList(manifest *pluginapi.ClustermanifestClusterapi, infra_provider capiutils.CapiInfraProvider, files *pluginapi.Files) error {
	config_name := capiutils.GetManifestConfigNameByCapiInfraProvider(infra_provider)

	files.Files = []*pluginapi.FilesItems0{}
	for _, config := range manifest.Configs {
		if config == nil {
			log.Warnln("Invalid parameter: nil configs in manifest.")
			continue
		}

		if config.Name == config_name {
			files.Files = append(files.Files, generateFileItemsByURL(config.BinURL, "bin"))
			for _, provider := range config.Providers {
				subpath := filepath.Join(provider.Parameters.ProviderLabel, provider.Parameters.Version)
				files.Files = append(files.Files, generateFileItemsByURL(provider.URL, subpath))
				files.Files = append(files.Files, generateFileItemsByURL(provider.Parameters.Metadata, subpath))
			}
		}
	}
	if manifest.CertManager == nil {
		return errCertMgrCfg
	}
	files.Files = append(files.Files, generateFileItemsByURL(manifest.CertManager.URL, "cert-manager"))

	return nil
}

func PluginMain(in eputils.SchemaMapData, outp *eputils.SchemaMapData) error {
	input_ep_params := input_ep_params(in)
	input_cluster_manifest := input_cluster_manifest(in)

	output_docker_images := output_docker_images(outp)
	output_files := output_files(outp)

	if input_ep_params == nil || input_ep_params.Ekconfig == nil {
		log.Errorln("Failed to find Ekconfigs for ClusterAPI cluster.")
		return errParam
	}

	infra_provider, err := capiutils.GetInfraProvider(input_ep_params.Ekconfig)
	if err != nil {
		log.Errorln(err)
		return errCAPIProvider
	}

	if input_cluster_manifest == nil || input_cluster_manifest.Clusterapi == nil {
		log.Errorln("Failed to find manifest for ClusterAPI cluster.")
		return errCAPIManifest
	}

	getDockerImagesList(input_cluster_manifest.Clusterapi, infra_provider, output_docker_images)
	if err = getProviderYamlList(input_cluster_manifest.Clusterapi, infra_provider, output_files); err != nil {
		log.Errorln(err)
		return errClusterProv
	}

	kindprovider, err := cutils.GetClusterManifest(input_cluster_manifest, "kind")
	if err != nil {
		log.Errorln("Failed to find kind info for ClusterAPI cluster.")
		return errCAPIKindLost
	}

	if err = getManagementClusterBinaryList(kindprovider, output_docker_images, output_files); err != nil {
		log.Errorln(err)
		return errMgmtCluster
	}

	log.Debugf("%v", output_docker_images)
	log.Debugf("%v", output_files)

	return nil
}
