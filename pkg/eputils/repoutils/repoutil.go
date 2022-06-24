/*
* Copyright (c) 2022 Intel Corporation.
*
* SPDX-License-Identifier: Apache-2.0
*
 */

//go:generate mockgen -destination=./mock/repoutils_mock.go -package=mock -copyright_file=../../../api/schemas/license-header.txt ep/pkg/eputils/repoutils RepoUtilsInterface

package repoutils

import (
	orasutils "ep/pkg/eputils/orasutils"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/url"
)

type (
	RepoUtilsInterface interface {
		PushFileToRepo(filepath, subRef, rev string) (string, error)
		PullFileFromRepo(filepath string, targeturl string) error
	}
)

var (
	errNoPushClient = errors.New("File push client not found!")
	errNoPullClient = errors.New("File pull client not found!")
)

func PushFileToRepo(filepath, subRef, rev string) (string, error) {
	var ref string
	var err error
	if orasutils.OrasCli != nil {
		ref, err = orasutils.OrasCli.OrasPushFile(filepath, subRef, rev)
		if err != nil {
			log.Errorln("Failed to push file", filepath, err)
			return "", err
		}
	} else {
		return "", errNoPushClient
	}
	return ref, nil
}

func PullFileFromRepo(filepath string, targeturl string) error {
	u, err := url.Parse(targeturl)
	if err != nil {
		log.Errorln("Failed to pull file", filepath, err)
		return err
	}
	if u.Scheme == "oci" {
		if orasutils.OrasCli != nil {
			err := orasutils.OrasCli.OrasPullFile(filepath, targeturl)
			if err != nil {
				log.Errorln("Failed to pull file", filepath, err)
				return err
			}
		} else {
			return errNoPullClient
		}
	}
	return nil
}
