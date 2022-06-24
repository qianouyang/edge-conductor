/*
* Copyright (c) 2022 Intel Corporation.
*
* SPDX-License-Identifier: Apache-2.0
*
 */
package capihostprovision

import (
	pluginapi "ep/pkg/api/plugins"
	certmgr "ep/pkg/certmgr"
	eputils "ep/pkg/eputils"
	capiutils "ep/pkg/eputils/capiutils"
	kubeutils "ep/pkg/eputils/kubeutils"
	"ep/pkg/executor"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	REGSERVERCERTFILE = "cert/pki/registry/registry.pem"
)

var (
	errDeploymentLaunchFail = errors.New("Deployment launch fail")
	errImgpkgBinMissing     = errors.New("imgpkg bin missing")
	errNodeNotReady         = errors.New("Node not ready")
)

func crioReleaseDownload(ep_params *pluginapi.EpParams, workFolder string, capiSetting *pluginapi.CapiSetting) error {
	if capiSetting.CRI.Name != capiutils.CONFIG_RUNTIME_CRIO {
		return nil
	}

	var err error
	imgpkgBin := filepath.Join(ep_params.Runtimedir, "bin", "imgpkg")
	if !eputils.FileExists(imgpkgBin) {
		log.Errorf("User delete Imgpkg bin in %s", ep_params.Runtimedir)
		return errImgpkgBinMissing
	}

	crioRelease := "crio.tar.gz"
	crioDownloadFolder := filepath.Join(workFolder, "crio")
	registry := fmt.Sprintf("%s:%s/library/", ep_params.Ekconfig.Parameters.GlobalSettings.ProviderIP, ep_params.Ekconfig.Parameters.GlobalSettings.RegistryPort)
	cmd := exec.Command(imgpkgBin, "pull", "-i", registry+crioRelease, "-o", crioDownloadFolder, "--registry-username", ep_params.Ekconfig.Parameters.Customconfig.Registry.User, "--registry-password", ep_params.Ekconfig.Parameters.Customconfig.Registry.Password, "--registry-ca-cert-path", certmgr.REGSERVERCERTFILE)
	_, err = eputils.RunCMD(cmd)
	if err != nil {
		log.Errorf("Failed to pull crio release. %v", err)
		return err
	}

	return nil
}

func DeploymentReady(management_kubeconfig, namespace, deploymentName string) error {
	count := 0

	MgrDeployment, err := kubeutils.NewDeployment(namespace, deploymentName, "", management_kubeconfig)
	if err != nil {
		log.Errorf("Failed to new deployment object %v", err)
		return err
	}

	for count < TIMEOUT {
		err = MgrDeployment.Get()
		if err != nil {
			log.Errorf("Failed to get deployment %s, %v", deploymentName, err)
			return err

		}

		status := MgrDeployment.GetStatus()
		if status.ReadyReplicas == status.Replicas {
			log.Infof("Deployment %s: is ready", deploymentName)
			break
		}

		log.Infof("Deployment %s is not ready, waiting", deploymentName)
		time.Sleep(WAIT_10_SEC * time.Second)
		count++
	}

	if count >= TIMEOUT {
		log.Errorf("Deployment %s: launch fail", deploymentName)
		return errDeploymentLaunchFail
	}

	return nil
}

func checkByoHosts(ep_params *pluginapi.EpParams, workFolder, management_kubeconfig string, clusterConfig *pluginapi.CapiClusterConfig, tmpl *capiutils.CapiTemplate) error {
	ready := false
	count := 0

	for count < TIMEOUT {
		cmd := exec.Command(ep_params.Workspace+"/kubectl", "get", "byohosts", "-n", clusterConfig.WorkloadCluster.Namespace, "--kubeconfig", management_kubeconfig)
		outputStr, err := eputils.RunCMD(cmd)
		if err != nil {
			log.Errorf("Failed to get workload config. %v", err)
			return err
		}

		lines := strings.Count(outputStr, "\n")
		if lines < 2 {
			log.Infof("sleep %d sec", WAIT_10_SEC)
			time.Sleep(WAIT_10_SEC * time.Second)
			count++
		} else {
			ready = true
			break
		}
	}

	if !ready {
		log.Errorf("Node is not ready, please check")
		return errNodeNotReady
	}

	return nil
}

func byohHostProvision(ep_params *pluginapi.EpParams, workFolder, management_kubeconfig string, clusterConfig *pluginapi.CapiClusterConfig, tmpl *capiutils.CapiTemplate) error {
	var err error

	err = crioReleaseDownload(ep_params, workFolder, &tmpl.CapiSetting)
	if err != nil {
		log.Errorf("Crio release pull fail, %v", err)
		return err
	}

	err = DeploymentReady(management_kubeconfig, "byoh-system", "byoh-controller-manager")
	if err != nil {
		log.Errorf("ByohCtlMgr deployment launch fail, %v", err)
		return err
	}

	err = executor.Run(clusterConfig.ByohAgent.InitScript, ep_params, tmpl.CapiSetting)
	if err != nil {
		log.Errorf("ByohAgent pre-provision failed, %v", err)
		return err
	}

	err = checkByoHosts(ep_params, workFolder, management_kubeconfig, clusterConfig, tmpl)
	if err != nil {
		log.Errorf("Failed to get available byohost, %v", err)
		return err
	}

	return nil
}
