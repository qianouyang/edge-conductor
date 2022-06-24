/*
* Copyright (c) 2022 Intel Corporation.
*
* SPDX-License-Identifier: Apache-2.0
*
 */
//go:generate mockgen -destination=./mock/goharborapi_mock.go -package=mock -copyright_file=../../../api/schemas/license-header.txt ep/pkg/eputils/restfulcli GoharborClientWrapper

package restfulcli

import (
	"encoding/base64"
	docker "ep/pkg/eputils/docker"
	"errors"
	"fmt"
	"github.com/go-openapi/validate"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

var (
	DayZeroCertFilePath = "cert/pki/ca.pem"
	errCertNull         = errors.New("Cert file is null")
	errHarborUrlNull    = errors.New("Harbor URL string is null")
	errPrjNull          = errors.New("Project string is null")
	errAuthNull         = errors.New("auth string is null")
	errHarborResponse   = errors.New("Harbor response error")
	errHarborAbnormal   = errors.New("Harbor response is abnormal")
	errHarborCertna     = errors.New("harbor certificate path does not exist")
	errHarborUrlEmpty   = errors.New("harbor URL is empty")
	errProjectName      = errors.New("project name is empty")
	errAuthEmpty        = errors.New("auth string is empty")
	errClientGet        = errors.New("client get error")
	errInternalServer   = errors.New("Internal server error")
	errHarborIPEmpty    = errors.New("input harbor IP is empty")
	errHarborPort       = errors.New("input harbor port is empty")
	errHarborUser       = errors.New("input harbor user is empty")
	errHarborPasswd     = errors.New("input harbor password is empty")
	errInputAuthSrv     = errors.New("input auth server address is empty")
	errInputPrjName     = errors.New("input project name is empty")
	errDay0CertFile     = errors.New("day-0 cert file path is empty")
)

func TlsBasicAuth(username, password string) string {
	var authType string = "basic "
	authUser := username + ":" + password
	return fmt.Sprintf("%s%s", authType, base64.StdEncoding.EncodeToString([]byte(authUser)))
}

func RegistryCreateProject(harborUrl, project, authStr, certFilePath string) error {
	_, err := os.Stat(certFilePath)
	if err != nil {
		return errCertNull
	}
	if len(strings.TrimSpace(harborUrl)) == 0 {
		return errHarborUrlNull
	}
	if len(strings.TrimSpace(project)) == 0 {
		return errPrjNull
	}
	if len(strings.TrimSpace(authStr)) == 0 {
		return errAuthNull
	}

	restApiUrl := fmt.Sprintf("https://%s/api/v2.0/projects", harborUrl)
	client := resty.New()
	client.SetRootCertificate(certFilePath)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", authStr).
		SetBody(map[string]interface{}{"project_name": project, "public": true}).
		Post(restApiUrl)
	if err != nil {
		log.Errorln("client post error:", err)
		return err
	}

	if resp != nil {
		StatusCode := fmt.Sprintf("%v", resp.StatusCode())
		if StatusCode == "201" || StatusCode == "409" {
			if StatusCode == "201" {
				log.Debugf("Create project %s successfully\n", project)
			}
			if StatusCode == "409" {
				log.Debugf("Project %s already exist\n", project)
			}
		} else {
			log.Errorf("Harbor response error: %s ", StatusCode)
			return errHarborResponse
		}
	} else {
		log.Infof("Harbor response is abnormal")
		return errHarborAbnormal
	}

	return nil
}

/*
RegistryProjectExists return 2 types variables, bool and error.
1. bool:
true: the project exists
false: the project doesn't exist
2. error:
error is not nil if error happened in RegistryProjectExists
error is nil if no error happened in RegistryProjectExists
*/
func RegistryProjectExists(harborUrl, project, authStr, certFilePath string) (bool, error) {
	_, err := os.Stat(certFilePath)
	if err != nil {
		log.Errorln("harbor certificate path is not exist")
		return false, errHarborCertna
	}
	if len(strings.TrimSpace(harborUrl)) == 0 {
		log.Errorln("harbor URL is empty")
		return false, errHarborUrlEmpty
	}
	if len(strings.TrimSpace(project)) == 0 {
		log.Errorln("project name is empty")
		return false, errProjectName
	}
	if len(strings.TrimSpace(authStr)) == 0 {
		log.Errorln("auth string is empty")
		return false, errAuthEmpty
	}

	restApiUrl := fmt.Sprintf("https://%s/api/v2.0/projects/%s", harborUrl, project)
	client := resty.New()
	client.SetRootCertificate(certFilePath)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", authStr).
		Get(restApiUrl)
	if err != nil {
		log.Errorln("client get error:", err)
		return false, errClientGet
	}

	if resp != nil {
		if resp.StatusCode() == http.StatusOK {
			log.Infof("Project %s exists\n", project)
		} else {
			log.Debugf("Harbor response code: %v ", resp.StatusCode())
			if "403" == fmt.Sprintf("%v", resp.StatusCode()) {
				log.Infof("Project %s not found", project)
				return false, nil
			} else {
				return false, errInternalServer
			}
		}
	} else {
		return false, errHarborAbnormal
	}

	return true, nil
}

func MapImageURLCreateHarborProject(harborIP, harborPort, harborUser, harborPass string, image []string) ([]string, error) {
	newImages := image

	if len(strings.TrimSpace(harborIP)) == 0 {
		err := errHarborIPEmpty
		return newImages, err
	}
	if len(strings.TrimSpace(harborPort)) == 0 {
		err := errHarborPort
		return newImages, err
	}
	if len(strings.TrimSpace(harborUser)) == 0 {
		err := errHarborUser
		return newImages, err
	}
	if len(strings.TrimSpace(harborPass)) == 0 {
		err := errHarborPasswd
		return newImages, err
	}
	if newImages, err := MapImageURLOnHarbor(image); err != nil {
		return newImages, err
	}

	auth, err := docker.GetAuthConf(harborIP, harborPort, harborUser, harborPass)
	if err != nil {
		return newImages, err
	}
	AuthStr := TlsBasicAuth(harborUser, harborPass)

	for _, mappedUrl := range newImages {
		log.Debugf("mapped URL is: %v", mappedUrl)

		// separate repository name in order to  create a project on day-0 harbor
		splitedimage := strings.SplitN(mappedUrl, "/", 2)

		if err := CreateHarborProject(auth.ServerAddress, splitedimage[0], AuthStr, DayZeroCertFilePath); err != nil {
			return newImages, err
		}
	}

	return newImages, nil
}

//MapImageURLOnHarbor
//=====================================================================================
// Possible inputs:
//
// IMAGE REGISTRY  |     docker.io            |    k8s.gcr.io               |
//--------------------------------------------------------------------------
// 1-layer-input   |       nginx              |         X                   |
//--------------------------------------------------------------------------
// 2-layer-input   |  docker.io/nginx         |  k8s.gcr.io/pause           |
//                 |  library/nginx
//                 |  portainerci/portainer   |                             |
//--------------------------------------------------------------------------
// 3-layer-input   |  docker.io/library/nginx |  k8s.gcr.io/coredns/cordns  |
//
// Step.1
// Filter docker.io and library

// Step.2
// Create project on harbor

// Step.3
// Tag and push to harbor on proper project

// Examples for docker.io
// =========================================================================================
// IMAGE LAYER     |     ORIGINAL IMG URL     |    MIRRORED IMG URL                       |
// -----------------------------------------------------------------------------------------
// 1-layer-input   |       nginx              |  harbor/docker.io/nginx                   |
// -----------------------------------------------------------------------------------------
// 2-layer-input   |  docker.io/nginx         |  harbor/docker.io/nginx                   |
//                 |  library/nginx           |  harbor/docker.io/library/nginx           |
//                 |  portainerci/portainer   |  harbor/docker.io/portainerci/portainer   |
// -----------------------------------------------------------------------------------------
// 3-layer-input   |  docker.io/library/nginx |  harbor/docker.io/library/nginx           |
// =========================================================================================

// Examples for k8s.gcr.io
// ====================================================================================
// IMAGE LAYER     |     ORIGINAL IMG URL       |    MIRRORED IMG URL                |
// ------------------------------------------------------------------------------------
// 2-layer-input   |  k8s.gcr.io/pause          |  harbor/k8s.gcr.io/pause           |
// ------------------------------------------------------------------------------------
// 3-layer-input   |  k8s.gcr.io/coredns/cordns |  harbor/k8s.gcr.io/coredns/cordns  |
// ====================================================================================

// =====================================================================================
func MapImageURLOnHarbor(image []string) ([]string, error) {
	log.Debugf("image list: %v", image)
	newImages := image

	for key, origin_url := range image {
		log.Debugf("Original URL is: %v", origin_url)
		//check first letter is not special letter
		if err := validate.Pattern("name", "body", origin_url, `^[a-zA-Z0-9]{1}`); err != nil {
			log.Infof("Image name %s must start with a letter or number", origin_url)
			return newImages, err
		}

		urlArray := strings.Split(origin_url, "/")

		var new_url string
		if len(urlArray) > 1 && strings.Contains(urlArray[0], ".") {
			new_url = origin_url
		} else {
			// add namespace "docker.io" or "docker.io/library"
			if strings.Contains(origin_url, "/") {
				new_url = fmt.Sprintf("%s%s", "docker.io/", origin_url)
			} else {
				new_url = fmt.Sprintf("%s%s", "docker.io/library/", origin_url)
			}
		}
		newImages[key] = new_url
	}

	return newImages, nil
}

func CreateHarborProject(authServerAddress, projectName, authStr, DayZeroCertFilePath string) error {
	if len(strings.TrimSpace(authServerAddress)) == 0 {
		err := errInputAuthSrv
		return err
	}
	if len(strings.TrimSpace(projectName)) == 0 {
		err := errInputPrjName
		return err
	}
	if len(strings.TrimSpace(authStr)) == 0 {
		err := errAuthEmpty
		return err
	}
	if len(strings.TrimSpace(DayZeroCertFilePath)) == 0 {
		err := errDay0CertFile
		return err
	}

	if existFlag, err := RegistryProjectExists(authServerAddress, projectName, authStr, DayZeroCertFilePath); err == nil {
		if !existFlag {
			if err := RegistryCreateProject(authServerAddress, projectName, authStr, DayZeroCertFilePath); err != nil {
				return err
			}
		}
	} else {
		return err
	}

	return nil
}
