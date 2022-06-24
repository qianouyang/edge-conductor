/*
* Copyright (c) 2021 Intel Corporation.
*
* SPDX-License-Identifier: Apache-2.0
*
 */

// Template auto-generated once, maintained by plugin owner.

package capiprovisionbinarydownload

import (
	pluginapi "ep/pkg/api/plugins"
	certmgr "ep/pkg/certmgr"
	eputils "ep/pkg/eputils"
	capiutils "ep/pkg/eputils/capiutils"
	docker "ep/pkg/eputils/docker"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const (
	REGSERVERCERTFILE = "cert/pki/registry/registry.pem"
)

var (
	errProvider = errors.New("Please select one provider")
)

func launchIpaDownload(ep_params *pluginapi.EpParams, workFolder string, clusterConfig *pluginapi.CapiClusterConfig, tmpl *capiutils.CapiTemplate) error {
	var err error

	dstFile := filepath.Join(workFolder, "ironic-containers.yaml")
	err = capiutils.TmplFileRendering(tmpl, workFolder, clusterConfig.BaremetelOperator.IronicContainers, dstFile)
	if err != nil {
		log.Errorf("Failed to render %s, %v", clusterConfig.BaremetelOperator.IronicContainers, err)
		return err
	}

	var ironicContainers pluginapi.Containers
	err = eputils.LoadSchemaStructFromYamlFile(&ironicContainers, dstFile)
	if err != nil {
		log.Errorf("Load capi cluster config failed, %v", err)
		return err
	}

	for _, c := range ironicContainers.Containers {
		if c.Name == "ipa-downloader" {
			err = docker.DockerRun(c)
			if err != nil {
				log.Errorf("Container %s run fail, %v", c.Name, err)
				return err
			}
		}
	}

	return nil
}

func copyIronicOsPrivisionImage(ep_params *pluginapi.EpParams, ironicHttpdFolder string, capiSetting *pluginapi.CapiSetting) error {
	var err error

	if capiSetting.IronicConfig.IronicOsImage != "" {
		osFileName := filepath.Base(capiSetting.IronicConfig.IronicOsImage)
		provisionOsImage := filepath.Join(ironicHttpdFolder, osFileName)

		osImagePath := filepath.Join(ep_params.Workspace, osFileName)

		if !eputils.FileExists(osImagePath) {
			log.Warnf("No OS image in Workspace, %s", osImagePath)
			return nil
		}

		if !eputils.FileExists(provisionOsImage) {
			log.Infof("Copy OS image for provision, %s", osImagePath)
			log.Infof("Please wait for a while")

			_, err = eputils.CopyFile(provisionOsImage, osImagePath)
			if err != nil {
				return err
			}

			sha256sumFilePath := provisionOsImage + ".shasum"
			sha256sum, _ := eputils.GenFileSHA256(provisionOsImage)
			err = eputils.WriteStringToFile(sha256sum, sha256sumFilePath)
			if err != nil {
				log.Errorf("Failed to write data %s, reason: %v", sha256sumFilePath, err)
				return err
			}
		} else {
			log.Infof("Os image is already ready")
		}
	}

	return nil
}

func downloadByohResource(ep_params *pluginapi.EpParams, workFolder string, capiSetting *pluginapi.CapiSetting) error {
	var err error
	imgpkgUrl := capiSetting.ByohConfig.DownloadBinURL
	bundleRegistry := capiSetting.ByohConfig.BundleRegistry
	bundleImage := capiSetting.ByohConfig.BundleImage
	byohAgentUrl := capiSetting.ByohConfig.HostAgentBinURL
	byohBundleFile := filepath.Join(workFolder, "byoh-bundle")
	imgpkgBin := filepath.Join(ep_params.Runtimedir, "bin", "imgpkg")
	byohAgentBin := filepath.Join(ep_params.Runtimedir, "bin", "byohHostAgent")
	crioReleaseUrl := capiSetting.CRI.BinURL
	crioRelease := "crio.tar.gz"
	crioDownloadPath := filepath.Join(workFolder, crioRelease)

	if !eputils.FileExists(imgpkgBin) {
		err = eputils.DownloadFile(imgpkgBin, imgpkgUrl)
		if err != nil {
			log.Errorf("Imgpkg bin download fail")
			return err
		}
	}

	if !eputils.FileExists(byohAgentBin) {
		err = eputils.DownloadFile(byohAgentBin, byohAgentUrl)
		if err != nil {
			log.Errorf("ByohAgent bin download fail")
			return err
		}
	}

	err = os.Chmod(imgpkgBin, 0700)
	if err != nil {
		log.Errorf("fail to change %s access right", imgpkgBin)
		return err
	}

	cmd := exec.Command(imgpkgBin, "pull", "--recursive", "-i", bundleRegistry+"/"+bundleImage, "-o", byohBundleFile)
	_, err = eputils.RunCMD(cmd)
	if err != nil {
		log.Errorf("Failed to pull %s. %v", byohBundleFile, err)
		return err
	}

	registry := fmt.Sprintf("%s:%s/library/", ep_params.Ekconfig.Parameters.GlobalSettings.ProviderIP, ep_params.Ekconfig.Parameters.GlobalSettings.RegistryPort)
	cmd = exec.Command(imgpkgBin, "push", "-i", registry+bundleImage, "-f", byohBundleFile, "--registry-username", ep_params.Ekconfig.Parameters.Customconfig.Registry.User, "--registry-password", ep_params.Ekconfig.Parameters.Customconfig.Registry.Password, "--registry-ca-cert-path", certmgr.REGSERVERCERTFILE)
	_, err = eputils.RunCMD(cmd)
	if err != nil {
		log.Errorf("Failed to push %s to local registry. %v", byohBundleFile, err)
		return err
	}

	if eputils.FileExists(crioRelease) {
		os.RemoveAll(crioRelease)
	}
	err = eputils.DownloadFile(crioDownloadPath, crioReleaseUrl)
	if err != nil {
		log.Errorf("Crio release download fail")
		return err
	}

	cmd = exec.Command(imgpkgBin, "push", "-i", registry+crioRelease, "-f", crioDownloadPath, "--registry-username", ep_params.Ekconfig.Parameters.Customconfig.Registry.User, "--registry-password", ep_params.Ekconfig.Parameters.Customconfig.Registry.Password, "--registry-ca-cert-path", certmgr.REGSERVERCERTFILE)
	_, err = eputils.RunCMD(cmd)
	if err != nil {
		log.Errorf("Failed to push %s to local registry. %v", crioDownloadPath, err)
		return err
	}

	return nil
}

func PluginMain(in eputils.SchemaMapData, outp *eputils.SchemaMapData) error {
	input_ep_params := input_ep_params(in)
	input_cluster_manifest := input_cluster_manifest(in)

	var err error
	var provider string
	providers := make([]string, 0)
	for _, p := range input_ep_params.Ekconfig.Parameters.Extensions {
		for _, i := range capiutils.InfraProviderList {
			if p == i {
				providers = append(providers, p)
			}
		}
	}

	if len(providers) != 1 {
		return errProvider
	} else {
		provider = providers[0]
	}

	workFolder := filepath.Join(input_ep_params.Runtimedir, provider)
	if err = eputils.CreateFolderIfNotExist(workFolder); err != nil {
		return err
	}

	var clusterConfig pluginapi.CapiClusterConfig
	clusterConfig.WorkloadCluster = new(pluginapi.CapiClusterConfigWorkloadCluster)
	clusterConfig.BaremetelOperator = new(pluginapi.CapiClusterConfigBaremetelOperator)
	err = eputils.LoadSchemaStructFromYamlFile(&clusterConfig, input_ep_params.Ekconfig.Cluster.Config)
	if err != nil {
		log.Errorf("Load capi cluster config failed, %v", err)
		return err
	}

	var capiSetting pluginapi.CapiSetting
	capiSetting.Provider = provider
	capiSetting.InfraProvider = new(pluginapi.CapiSettingInfraProvider)
	capiSetting.ByohConfig = new(pluginapi.CapiSettingByohConfig)
	capiSetting.IronicConfig = new(pluginapi.CapiSettingIronicConfig)

	var tmpl capiutils.CapiTemplate
	capiutils.GetCapiSetting(input_ep_params, input_cluster_manifest, &clusterConfig, &capiSetting)
	err = capiutils.GetCapiTemplate(input_ep_params, capiSetting, &tmpl)
	if err != nil {
		log.Errorf("CapiHostProvision, get CapiTemplate failed, %v", err)
		return err
	}

	if provider == capiutils.CAPI_METAL3 {
		ironic_data_dir := filepath.Join(workFolder, "ironic")
		ironic_image_dir := filepath.Join(ironic_data_dir, "html", "images")
		err = os.MkdirAll(ironic_data_dir, 0755)
		if err != nil {
			return err
		}

		err = os.MkdirAll(ironic_image_dir, 0755)
		if err != nil {
			return err
		}

		err = launchIpaDownload(input_ep_params, workFolder, &clusterConfig, &tmpl)
		if err != nil {
			log.Errorf("Ironic provision agent download failed, %v", err)
			return err
		}

		err = copyIronicOsPrivisionImage(input_ep_params, ironic_image_dir, &capiSetting)
		if err != nil {
			log.Errorf("Copy ironic os provision image fail, %v", err)
			return err
		}
	} else if provider == capiutils.CAPI_BYOH {
		err = downloadByohResource(input_ep_params, workFolder, &capiSetting)
		if err != nil {
			log.Errorf("Byoh resource download failed, %v", err)
			return err
		}
	}

	return nil
}
