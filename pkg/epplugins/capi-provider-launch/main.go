/*
* Copyright (c) 2021 Intel Corporation.
*
* SPDX-License-Identifier: Apache-2.0
*
 */

// Template auto-generated once, maintained by plugin owner.

package capiproviderlaunch

import (
	"bytes"
	pluginapi "ep/pkg/api/plugins"
	eputils "ep/pkg/eputils"
	capiutils "ep/pkg/eputils/capiutils"
	cutils "ep/pkg/eputils/conductorutils"
	repoutils "ep/pkg/eputils/repoutils"
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

type ProviderConfig struct {
	Name, Version, Label string
}

type CertManagerConfig struct {
	Version string
}

type ClusterCtlConfig struct {
	CoreProvider, BootstrapProvider, ControlPlaneProvider, InfrastructureProvider *ProviderConfig
	CertManager                                                                   *CertManagerConfig
	RuntimeDir                                                                    string
}

const template_config_yaml = `
cert-manager:
  version: "{{ .CertManager.Version }}"
  url: "{{ .RuntimeDir }}/cert-manager/{{ .CertManager.Version }}/cert-manager.yaml"

providers:
- name: "{{ .CoreProvider.Name }}"
  type: "CoreProvider"
  url: "{{ .RuntimeDir }}/{{ .CoreProvider.Label}}/{{ .CoreProvider.Version}}/core-components.yaml"
- name: "{{ .BootstrapProvider.Name }}"
  type: "BootstrapProvider"
  url: "{{ .RuntimeDir }}/{{ .BootstrapProvider.Label}}/{{ .BootstrapProvider.Version}}/bootstrap-components.yaml"
- name: "{{ .ControlPlaneProvider.Name }}"
  type: "ControlPlaneProvider"
  url: "{{ .RuntimeDir }}/{{ .ControlPlaneProvider.Label}}/{{ .ControlPlaneProvider.Version}}/control-plane-components.yaml"
- name: "{{ .InfrastructureProvider.Name }}"
  type: "InfrastructureProvider"
  url: "{{ .RuntimeDir }}/{{ .InfrastructureProvider.Label}}/{{ .InfrastructureProvider.Version}}/infrastructure-components.yaml"
`

const kind_config_yaml = `
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  apiServerAddress: "{{ .Ekconfig.Parameters.GlobalSettings.ProviderIP }}"
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 8088
    protocol: TCP
  - containerPort: 443
    hostPort: 9443
    protocol: TCP
  extraMounts:
    - containerPath: /etc/containerd/certs.d/
      hostPath: {{ .Runtimedir }}/data/cert/
containerdConfigPatches:
  - |-
      [plugins."io.containerd.grpc.v1.cri".registry.mirrors."k8s.gcr.io"]
        endpoint = ["https://{{ .Ekconfig.Parameters.GlobalSettings.ProviderIP }}:{{ .Ekconfig.Parameters.GlobalSettings.RegistryPort }}/v2/k8s.gcr.io", "https://k8s.gcr.io"]
      [plugins."io.containerd.grpc.v1.cri".registry.mirrors."gcr.io"]
        endpoint = ["https://{{ .Ekconfig.Parameters.GlobalSettings.ProviderIP }}:{{ .Ekconfig.Parameters.GlobalSettings.RegistryPort }}/v2/gcr.io", "https://gcr.io"]
      [plugins."io.containerd.grpc.v1.cri".registry.mirrors."quay.io"]
        endpoint = ["https://{{ .Ekconfig.Parameters.GlobalSettings.ProviderIP }}:{{ .Ekconfig.Parameters.GlobalSettings.RegistryPort }}/v2/quay.io", "https://quay.io"]
      [plugins."io.containerd.grpc.v1.cri".registry.mirrors."projects.registry.vmware.com"]
        endpoint = ["https://{{ .Ekconfig.Parameters.GlobalSettings.ProviderIP }}:{{ .Ekconfig.Parameters.GlobalSettings.RegistryPort }}/v2/projects.registry.vmware.com/", "https://projects.registry.vmware.com/"]
      [plugins."io.containerd.grpc.v1.cri".registry.configs."{{ .Ekconfig.Parameters.GlobalSettings.ProviderIP }}:{{ .Ekconfig.Parameters.GlobalSettings.RegistryPort }}".tls]
        ca_file = "/etc/containerd/certs.d/{{ .Ekconfig.Parameters.GlobalSettings.ProviderIP }}:{{ .Ekconfig.Parameters.GlobalSettings.RegistryPort }}/ca.crt"
      [plugins."io.containerd.grpc.v1.cri".registry.configs."{{ .Ekconfig.Parameters.GlobalSettings.ProviderIP }}:{{ .Ekconfig.Parameters.GlobalSettings.RegistryPort }}".auth]
        username = "{{ .Ekconfig.Parameters.Customconfig.Registry.User }}"
        password = "{{ .Ekconfig.Parameters.Customconfig.Registry.Password }}"
`

var (
	errCertMgrCfg       = errors.New("Cert manager config is missing in manifest.")
	errProvConfig       = errors.New("provider config lost!")
	errProviderLost     = errors.New("provider lost!")
	errCAPIManifest     = errors.New("CAPI manifest Lost!")
	errCAPIProvider     = errors.New("CAPI provider config invalidate")
	errEKConfigParm     = errors.New("ekconfig parameters missing!")
	errGenCfgClusterctl = errors.New("Failed to generate config file for clusterctl")
	errGenProvRepoCctl  = errors.New("Failed to generate local provider repo for clusterctl!")
	errGetCAPIInfraProv = errors.New("Failed to get CAPI infra provider config!")
	errInitClusterctl   = errors.New("Failed to init clusterctl!")
	errLaunchMgmtClster = errors.New("Failed to launch management cluster")
	errPullFile         = errors.New("Failed to pull file")
	errRunClusterctlCmd = errors.New("Failed to run clusterctl command!")
	errIncorrectParam   = errors.New("Incorrect parameter")
	errNoKindBinInReg   = errors.New("No kin bin on registry")
	errPullingFile      = errors.New("Pulling file failure!")
)

func launchManagementCluster(ep_params *pluginapi.EpParams, clusterManifest *pluginapi.Clustermanifest, files *pluginapi.Files) error {
	mgr_cluster_kubeconfig := ""
	var err error

	for _, ext := range ep_params.Extensions {
		if ext.Name == capiutils.CAPI_BYOH || ext.Name == capiutils.CAPI_METAL3 {
			for _, ext_section := range ext.Extension.Extension {
				if ext_section.Name == capiutils.EXTENSION_INFRA_PROVIDER {
					for _, config := range ext_section.Config {
						if config.Name == capiutils.CONFIG_MANAGEMENT_CLUSTER_KUBECONFIG {
							mgr_cluster_kubeconfig = config.Value
						}
					}
				}
			}
		}
	}

	if mgr_cluster_kubeconfig != "" {
		log.Infof("User provides management cluster")
		m_kubeconfig := filepath.Join(ep_params.Runtimedir, capiutils.MANAGEMENT_KUBECONFIG)
		_, err = eputils.CopyFile(m_kubeconfig, mgr_cluster_kubeconfig)
		if err != nil {
			log.Errorf("Failed to copy %s", mgr_cluster_kubeconfig)
			return err
		}
		return nil
	}

	kindURL := ""
	for _, file := range files.Files {
		if strings.Contains(file.Mirrorurl, "capi/kind") {
			kindURL = file.Mirrorurl
		}
	}

	if kindURL == "" {
		return errNoKindBinInReg
	}

	kindBin := filepath.Join(ep_params.Runtimebin, "kind")
	err = repoutils.PullFileFromRepo(kindBin, kindURL)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	err = os.Chmod(kindBin, 0700)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	kindTemplateConfigPath := filepath.Join(ep_params.Runtimedir, "kindTemplateconfig.yaml")
	err = eputils.WriteStringToFile(kind_config_yaml, kindTemplateConfigPath)
	if err != nil {
		log.Errorf("Write kind config template file fail, %v", err)
		return err
	}

	kindConfigPath := filepath.Join(ep_params.Runtimedir, "kindconfig.yaml")
	err = eputils.FileTemplateConvert(kindTemplateConfigPath, kindConfigPath)
	if err != nil {
		log.Errorf("Gent kind config fail, %v", err)
		return err
	}

	kubeconfigPath := filepath.Join(ep_params.Runtimedir, capiutils.MANAGEMENT_KUBECONFIG)
	if ep_params.Ekconfig == nil || ep_params.Ekconfig.Parameters == nil ||
		ep_params.Ekconfig.Parameters.GlobalSettings == nil {
		return errEKConfigParm
	}
	kindprovider, err := cutils.GetClusterManifest(clusterManifest, "kind")
	if err != nil {
		return errEKConfigParm
	}
	kindImageNode, err := cutils.GetImageFromProvider(kindprovider, "img_node")
	if err != nil {
		return errEKConfigParm
	}

	imageUrl := ep_params.Ekconfig.Parameters.GlobalSettings.ProviderIP +
		":" + ep_params.Ekconfig.Parameters.GlobalSettings.RegistryPort +
		"/docker.io/" + kindImageNode

	// If mngr cluster exist, it will be deleled first.
	cmd := exec.Command(kindBin, "delete", "cluster", "--name", capiutils.MANAGEMENT_CLUSTER_NAME)
	_, err = eputils.RunCMD(cmd)
	if err != nil {
		log.Errorf("Failed to delete existed management cluster, %v", err)
		return errRunClusterctlCmd
	}

	log.Infof("Starting management cluster ...")
	cmd = exec.Command(kindBin, "create", "cluster", "--name", capiutils.MANAGEMENT_CLUSTER_NAME, "--config", kindConfigPath, "--kubeconfig", kubeconfigPath, "--image", imageUrl)
	_, err = eputils.RunCMD(cmd)

	if err != nil {
		log.Errorf("Create management cluster fail, %v", err)
		return errRunClusterctlCmd
	}

	os.Remove(kindTemplateConfigPath)
	os.Remove(kindConfigPath)

	return nil
}

func generateClusterCtlConfig(manifest *pluginapi.ClustermanifestClusterapi, infra_provider capiutils.CapiInfraProvider, ep_params *pluginapi.EpParams) (config ClusterCtlConfig, err error) {
	config.RuntimeDir = filepath.Join(ep_params.Runtimedir, "clusterapi")

	if manifest.CertManager == nil {
		err = errCertMgrCfg
		return
	}
	config.CertManager = &CertManagerConfig{Version: manifest.CertManager.Version}

	config_name := capiutils.GetManifestConfigNameByCapiInfraProvider(infra_provider)
	for _, mconfig := range manifest.Configs {
		if mconfig == nil {
			log.Warnln("Invalid parameter: nil configs in manifest.")
			continue
		}

		if mconfig.Name == config_name {
			for _, provider := range mconfig.Providers {
				if provider.Parameters == nil {
					log.Errorf("provider %s miss config", provider.Name)
					err = errProvConfig
					return
				}

				pconfig := ProviderConfig{Name: provider.Name, Version: provider.Parameters.Version, Label: provider.Parameters.ProviderLabel}
				switch provider.ProviderType {
				case "CoreProvider":
					config.CoreProvider = &pconfig
				case "BootstrapProvider":
					config.BootstrapProvider = &pconfig
				case "ControlPlaneProvider":
					config.ControlPlaneProvider = &pconfig
				case "InfrastructureProvider":
					config.InfrastructureProvider = &pconfig
				}
			}
		}
	}

	if config.CoreProvider == nil || config.BootstrapProvider == nil || config.ControlPlaneProvider == nil || config.InfrastructureProvider == nil {
		err = errProviderLost
		return
	}

	return
}

func generateLocalPRPath(file *pluginapi.FilesItems0, config *ClusterCtlConfig) (target string) {
	file_name := path.Base(file.URL)

	if strings.Contains(file.Mirrorurl, config.CoreProvider.Label) {
		target = filepath.Join(config.RuntimeDir, config.CoreProvider.Label, config.CoreProvider.Version, file_name)
	} else if strings.Contains(file.Mirrorurl, config.BootstrapProvider.Label) {
		target = filepath.Join(config.RuntimeDir, config.BootstrapProvider.Label, config.BootstrapProvider.Version, file_name)
	} else if strings.Contains(file.Mirrorurl, config.ControlPlaneProvider.Label) {
		target = filepath.Join(config.RuntimeDir, config.ControlPlaneProvider.Label, config.ControlPlaneProvider.Version, file_name)
	} else if strings.Contains(file.Mirrorurl, config.InfrastructureProvider.Label) {
		target = filepath.Join(config.RuntimeDir, config.InfrastructureProvider.Label, config.InfrastructureProvider.Version, file_name)
	} else if strings.Contains(file.Mirrorurl, "cert-manager") {
		target = filepath.Join(config.RuntimeDir, "cert-manager", config.CertManager.Version, file_name)
	}
	return
}

func generateLocalProviderRepo(files *pluginapi.Files, clusterctl_config *ClusterCtlConfig) error {
	for _, file := range files.Files {
		target := generateLocalPRPath(file, clusterctl_config)
		if target == "" {
			continue
		}

		err := repoutils.PullFileFromRepo(target, file.Mirrorurl)
		if err != nil {
			log.Errorf("%s, %s", err, target)
			return errPullFile
		}
	}

	return nil
}

func generateClusterctlConfig(config *ClusterCtlConfig) (targetPath string, err error) {
	tpconfig := template.Must(template.New("pconfigyaml").Parse(template_config_yaml))

	if err = eputils.CreateFolderIfNotExist(config.RuntimeDir); err != nil {
		return
	}
	targetPath = filepath.Join(config.RuntimeDir, "config.yaml")

	var content bytes.Buffer
	if err = tpconfig.Execute(&content, *config); err != nil {
		return
	}

	if err = eputils.WriteStringToFile(content.String(), targetPath); err != nil {
		return
	}

	return
}

func launchCapiProvider(ep_params *pluginapi.EpParams, config *ClusterCtlConfig, config_path string, input_files *pluginapi.Files) error {
	bin := filepath.Join(ep_params.Runtimebin, "clusterctl")
	err := repoutils.PullFileFromRepo(bin, input_files.Files[0].Mirrorurl)
	if err != nil {
		log.Errorf("%s", err)
		return errPullingFile
	}
	err = os.Chmod(bin, 0700)
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	mClusterConfig := capiutils.GetManagementClusterKubeconfig(ep_params)
	if mClusterConfig != "" {
		cmd = exec.Command(bin, "init", "--core", config.CoreProvider.Name+":"+config.CoreProvider.Version, "--bootstrap", config.BootstrapProvider.Name+":"+config.BootstrapProvider.Version, "--control-plane", config.ControlPlaneProvider.Name+":"+config.ControlPlaneProvider.Version, "--infrastructure", config.InfrastructureProvider.Name+":"+config.InfrastructureProvider.Version, "--config", config_path, "--kubeconfig", mClusterConfig)
	} else {
		cmd = exec.Command(bin, "init", "--core", config.CoreProvider.Name+":"+config.CoreProvider.Version, "--bootstrap", config.BootstrapProvider.Name+":"+config.BootstrapProvider.Version, "--control-plane", config.ControlPlaneProvider.Name+":"+config.ControlPlaneProvider.Version, "--infrastructure", config.InfrastructureProvider.Name+":"+config.InfrastructureProvider.Version, "--config", config_path)
	}

	_, err = eputils.RunCMD(cmd)
	if err != nil {
		log.Errorf("%s", err)
		return errRunClusterctlCmd
	}

	return nil
}

func PluginMain(in eputils.SchemaMapData, outp *eputils.SchemaMapData) error {
	input_ep_params := input_ep_params(in)
	input_cluster_manifest := input_cluster_manifest(in)
	input_files := input_files(in)

	log.Infof("Plugin: capi-provider-launch")

	if input_ep_params == nil || input_ep_params.Ekconfig == nil {
		log.Errorln("Failed to find Ekconfigs for ClusterAPI cluster.")
		return errIncorrectParam
	}

	infra_provider, err := capiutils.GetInfraProvider(input_ep_params.Ekconfig)
	if err != nil {
		log.Errorln(err)
		return errGetCAPIInfraProv
	}

	if input_cluster_manifest == nil || input_cluster_manifest.Clusterapi == nil {
		log.Errorln("Failed to find manifest for ClusterAPI cluster.")
		return errCAPIManifest
	}

	clusterCtlConfig, err := generateClusterCtlConfig(input_cluster_manifest.Clusterapi, infra_provider, input_ep_params)
	if err != nil {
		log.Errorln(err)
		return errCAPIProvider
	}

	if err := generateLocalProviderRepo(input_files, &clusterCtlConfig); err != nil {
		log.Errorln(err)
		return errGenProvRepoCctl
	}

	configFilePath, err := generateClusterctlConfig(&clusterCtlConfig)
	if err != nil {
		log.Errorln(err)
		return errGenCfgClusterctl
	}

	if err := launchManagementCluster(input_ep_params, input_cluster_manifest, input_files); err != nil {
		log.Errorln(err)
		return errLaunchMgmtClster
	}

	if err := launchCapiProvider(input_ep_params, &clusterCtlConfig, configFilePath, input_files); err != nil {
		log.Errorf("%s", err)
		return errInitClusterctl
	}

	return nil
}
