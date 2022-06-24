/*
* Copyright (c) 2021 Intel Corporation.
*
* SPDX-License-Identifier: Apache-2.0
*
 */

package capiutils

import (
	"encoding/base64"
	pluginapi "ep/pkg/api/plugins"
	eputils "ep/pkg/eputils"
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	CAPI_METAL3 = "capi-metal3"
	CAPI_BYOH   = "capi-byoh"

	EXTENSION_INFRA_PROVIDER = "Infra-provider"
	EXTENSION_IRONIC_CONFIG  = "Ironic-config"

	CONFIG_MANAGEMENT_CLUSTER_KUBECONFIG          = "Management-cluster-kubeconfig"
	CONFIG_WORKLOAD_CLUSTER_NETWORK               = "Workload-cluster-network"
	CONFIG_WORKLOAD_CLUSTER_NETWORK_GATEWAY       = "Workload-cluster-network-gateway"
	CONFIG_WORKLOAD_CLUSTER_CONTROLPLANE_ENDPOINT = "Workload-cluster-controlplane-endpoint"
	CONFIG_WORKLOAD_CLUSTER_NODE_ADDRESS_START    = "Workload-cluster-node-address-start"
	CONFIG_WORKLOAD_CLUSTER_NODE_ADDRESS_END      = "Workload-cluster-node-address-end"
	CONFIG_WORKLOAD_CLUSTER_NODE_ADDRESS_PREFIX   = "Workload-cluster-node-address-prefix"
	CONFIG_WORKLOAD_CLUSTER_NODE_USERNAME         = "Workload-cluster-node-username"
	CONFIG_WORKLOAD_CLUSTER_NIC_NAME              = "Workload-cluster-nic-name"
	CONFIG_AUTHORIZED_SSH_PUBLIC_KEY              = "Authorized-ssh-public-key"
	CONFIG_IRONIC_PROVISION_NIC                   = "Ironic-provision-nic"
	CONFIG_IRONIC_PROVISION_IP                    = "Ironic-provision-ip"
	CONFIG_IRONIC_DHCP                            = "Ironic-dhcp-range"
	CONFIG_IRONIC_HTTP_PORT                       = "Ironic-http-port"
	CONFIG_IRONIC_OS_IMAGE                        = "Ironic-os-image"

	CONFIG_NAME_METAL3 = "metal3"
	CONFIG_NAME_BYOH   = "byoh"

	CONFIG_RUNTIME_CRIO       = "crio"
	CONFIG_RUNTIME_CONTAINERD = "containerd"

	CONFIG_WORKLOAD_CLUSTER_CONTROLPLANE = "controlplane"
	CONFIG_WORKLOAD_CLUSTER_WORKER       = "worker"

	MANAGEMENT_KUBECONFIG   = "m_kubeconfig"
	MANAGEMENT_CLUSTER_NAME = "capi-management"
)

var (
	InfraProviderList = []string{
		CAPI_METAL3,
		CAPI_BYOH,
	}
)

var (
	errInput       = errors.New("Incorrect input: parameter in ekconfig missing.")
	errProvider    = errors.New("Please select one provider")
	errCAPISetting = errors.New("Invalid CAPI Setting. InfraProvider info missing.")
	errNumberNodes = errors.New("Invalid CAPI Setting. Number of nodes in workload cluster invalid.")
)

type CapiTemplate struct {
	pluginapi.EpParams
	CapiSetting pluginapi.CapiSetting
}

func GetCapiTemplate(epparams *pluginapi.EpParams, setting pluginapi.CapiSetting, cp *CapiTemplate) error {
	err := eputils.ConvertSchemaStruct(epparams, cp)
	if err != nil {
		return err
	}
	cp.CapiSetting = setting
	return nil
}

func GetCapiSetting(epparams *pluginapi.EpParams, clusterManifest *pluginapi.Clustermanifest, clusterConfig *pluginapi.CapiClusterConfig, setting *pluginapi.CapiSetting) {

	var extension *pluginapi.Extension
	cri := &pluginapi.CapiSettingCRI{}

	for _, ex := range epparams.Extensions {
		extension = ex.Extension
		if ex.Name == CAPI_BYOH {
			setting.ByohConfig.HostAgentBinURL = clusterManifest.Clusterapi.ByohConfig.HostAgentBinURL
			setting.ByohConfig.DownloadBinURL = clusterManifest.Clusterapi.ByohConfig.DownloadBinURL
			setting.ByohConfig.BundleRegistry = clusterManifest.Clusterapi.ByohConfig.BundleRegistry
			setting.ByohConfig.BundleImage = clusterManifest.Clusterapi.ByohConfig.BundleImage
			setting.ByohConfig.BundleTag = strings.Split(clusterManifest.Clusterapi.ByohConfig.BundleImage, ":")[1]

			if epparams.Ekconfig != nil && epparams.Ekconfig.Parameters != nil {
				setting.InfraProvider.WorkloadClusterControlPlaneNum, setting.InfraProvider.WorkloadClusterWorkerNodeNum = getWorkloadClusterNodesNum(epparams.Ekconfig.Parameters.Nodes)
			}
		}

		for _, item := range extension.Extension {
			for _, config := range item.Config {
				if config.Name == CONFIG_MANAGEMENT_CLUSTER_KUBECONFIG {
					setting.InfraProvider.ManagementClusterKubeconfig = config.Value
				} else if config.Name == CONFIG_WORKLOAD_CLUSTER_NETWORK {
					setting.InfraProvider.WorkloadClusterNetwork = config.Value
				} else if config.Name == CONFIG_WORKLOAD_CLUSTER_NETWORK_GATEWAY {
					setting.InfraProvider.WorkloadClusterNetworkGateway = config.Value
				} else if config.Name == CONFIG_WORKLOAD_CLUSTER_CONTROLPLANE_ENDPOINT {
					setting.InfraProvider.WorkloadClusterControlplaneEndpoint = config.Value
				} else if config.Name == CONFIG_WORKLOAD_CLUSTER_NODE_ADDRESS_START {
					setting.InfraProvider.WorkloadClusterNodeAddressStart = config.Value
				} else if config.Name == CONFIG_WORKLOAD_CLUSTER_NODE_ADDRESS_END {
					setting.InfraProvider.WorkloadClusterNodeAddressEnd = config.Value
				} else if config.Name == CONFIG_WORKLOAD_CLUSTER_NODE_ADDRESS_PREFIX {
					setting.InfraProvider.WorkloadClusterNodeAddressPrefix = config.Value
				} else if config.Name == CONFIG_WORKLOAD_CLUSTER_NODE_USERNAME {
					setting.InfraProvider.WorkloadClusterNodeUsername = config.Value
				} else if config.Name == CONFIG_WORKLOAD_CLUSTER_NIC_NAME {
					setting.InfraProvider.WorkloadClusterNicName = config.Value
				} else if config.Name == CONFIG_AUTHORIZED_SSH_PUBLIC_KEY {
					setting.InfraProvider.AuthorizedSSHPublicKey = config.Value
				} else if config.Name == CONFIG_IRONIC_PROVISION_NIC {
					setting.IronicConfig.IronicProvisionNic = config.Value
				} else if config.Name == CONFIG_IRONIC_PROVISION_IP {
					setting.IronicConfig.IronicProvisionIP = config.Value
				} else if config.Name == CONFIG_IRONIC_DHCP {
					setting.IronicConfig.IronicDhcpRange = config.Value
				} else if config.Name == CONFIG_IRONIC_HTTP_PORT {
					setting.IronicConfig.IronicHTTPPort = config.Value
				} else if config.Name == CONFIG_IRONIC_OS_IMAGE {
					setting.IronicConfig.IronicOsImage = config.Value
				}
			}
		}

		for _, config := range clusterManifest.Clusterapi.Configs {
			if "capi-"+config.Name == ex.Name {
				for _, criBin := range config.RuntimeBins {
					cri.Name = clusterManifest.Clusterapi.Runtime
					if cri.Name == criBin.Name {
						cri.Endpoint = "unix://" + filepath.Join("/var/run", cri.Name, cri.Name+".sock")
						cri.Version = criBin.Version
						cri.BinURL = criBin.URL
					}
				}
			}
		}
	}

	setting.InfraProvider.WorkloadClusterName = clusterConfig.WorkloadCluster.Name
	setting.InfraProvider.WorkloadClusterNamespace = clusterConfig.WorkloadCluster.Namespace
	setting.CRI = cri

	var AuthStr string
	if epparams.Ekconfig == nil ||
		epparams.Ekconfig.Parameters == nil ||
		epparams.Ekconfig.Parameters.Customconfig == nil ||
		epparams.Ekconfig.Parameters.Customconfig.Registry == nil {
		AuthStr = ""
	} else {
		AuthStr = base64.StdEncoding.EncodeToString([]byte(epparams.Ekconfig.Parameters.Customconfig.Registry.User + ":" + epparams.Ekconfig.Parameters.Customconfig.Registry.Password))
	}
	setting.Registry = &pluginapi.CapiSettingRegistry{Auth: AuthStr}
}

func TmplFileRendering(tmpl *CapiTemplate, workFolder, url, dstFile string) error {
	var err error

	templatePath := filepath.Join(workFolder, "template.yaml")
	err = eputils.DownloadFile(templatePath, url)
	if err != nil {
		return err
	}

	templateData, err := ioutil.ReadFile(templatePath)
	if err != nil {
		log.Errorf("Ready cluster config %s failed, %v", url, err)
		return err
	}

	rawData, err := eputils.StringTemplateConvertWithParams(string(templateData), tmpl)
	if err != nil {
		log.Errorf("For metal3 provider, get bmo manifest failed, %v", err)
		return err
	}

	err = eputils.WriteStringToFile(rawData, dstFile)
	if err != nil {
		log.Errorf("Write string to file %s failed", dstFile)
		return err
	}
	os.RemoveAll(templatePath)

	return nil
}

// TODO refactor: using GetInfraProvider in capi plugins to get infra provider
// instead of implement in each plugins.

type CapiInfraProvider string

const (
	METAL3 = "capi-metal3"
	BYOH   = "capi-byoh"
)

var (
	SupportedInfraProvider = []CapiInfraProvider{
		METAL3,
		BYOH,
	}
)

func (provider CapiInfraProvider) IsSupported() bool {
	for _, p := range SupportedInfraProvider {
		if p == provider {
			return true
		}
	}

	return false
}

func GetInfraProvider(inputEkconfig *pluginapi.Ekconfig) (provider CapiInfraProvider, err error) {
	var providerNum int
	if inputEkconfig.Parameters == nil {
		err = errInput
		return
	}

	for _, p := range inputEkconfig.Parameters.Extensions {
		if CapiInfraProvider(p).IsSupported() {
			providerNum = providerNum + 1
			provider = CapiInfraProvider(p)
		}
	}

	if providerNum != 1 {
		err = errProvider
	}

	return
}

func GetManifestConfigNameByCapiInfraProvider(provider CapiInfraProvider) string {
	switch provider {
	case METAL3:
		return CONFIG_NAME_METAL3
	case BYOH:
		return CONFIG_NAME_BYOH
	}
	return ""
}

func GetManagementClusterKubeconfig(ep_params *pluginapi.EpParams) (configPath string) {
	mgr_cluster_kubeconfig := ""

	for _, ext := range ep_params.Extensions {
		if ext.Name == CAPI_BYOH || ext.Name == CAPI_METAL3 {
			for _, ext_section := range ext.Extension.Extension {
				if ext_section.Name == EXTENSION_INFRA_PROVIDER {
					for _, config := range ext_section.Config {
						if config.Name == CONFIG_MANAGEMENT_CLUSTER_KUBECONFIG {
							mgr_cluster_kubeconfig = config.Value
						}
					}
				}
			}
		}
	}

	if mgr_cluster_kubeconfig == "" {
		mgr_cluster_kubeconfig = filepath.Join(ep_params.Runtimedir, MANAGEMENT_KUBECONFIG)
	}

	return mgr_cluster_kubeconfig
}

func getWorkloadClusterNodesNum(nodes []*pluginapi.Node) (int64, int64) {
	nodeNum, controlPlaneNum, workerNum := int64(0), int64(0), int64(0)
	for _, node := range nodes {
		nodeNum = nodeNum + 1
		if node.Role == nil {
			return controlPlaneNum, workerNum
		}

		for _, role := range node.Role {
			switch role {
			case CONFIG_WORKLOAD_CLUSTER_CONTROLPLANE:
				controlPlaneNum = controlPlaneNum + 1
			case CONFIG_WORKLOAD_CLUSTER_WORKER:
				workerNum = workerNum + 1
			}
		}
	}

	if nodeNum != controlPlaneNum+workerNum {
		log.Errorf("Invalid EK config. Got %d workers and %d control plane in %d nodes. The number of Workload cluster nodes should be the sum of workers and control planes.", workerNum, controlPlaneNum, nodeNum)
		controlPlaneNum = 0
		workerNum = 0
	}

	if controlPlaneNum%2 == 0 {
		controlPlaneNum = controlPlaneNum - 1
		workerNum = nodeNum - controlPlaneNum
		log.Warnf("Number of control plane should be odd number. Revise workload cluster as %d control planes and %d workers.", controlPlaneNum, workerNum)
	}

	return controlPlaneNum, workerNum
}

func CheckCapiSetting(setting *pluginapi.CapiSetting) error {
	if setting.InfraProvider == nil {
		return errCAPISetting
	}

	if setting.ByohConfig != nil &&
		setting.ByohConfig.HostAgentBinURL != "" &&
		(setting.InfraProvider.WorkloadClusterControlPlaneNum == 0 ||
			setting.InfraProvider.WorkloadClusterWorkerNodeNum == 0) {
		return errNumberNodes
	}

	return nil
}
