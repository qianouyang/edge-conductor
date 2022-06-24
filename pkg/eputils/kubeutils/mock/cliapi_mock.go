//
//   Copyright (c) 2022 Intel Corporation.
//
//   SPDX-License-Identifier: Apache-2.0
//
//
//

// Code generated by MockGen. DO NOT EDIT.
// Source: ep/pkg/eputils/kubeutils (interfaces: KubeClientWrapper)

// Package mock is a generated GoMock package.
package mock

import (
	plugins "ep/pkg/api/plugins"
	kubeutils "ep/pkg/eputils/kubeutils"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "k8s.io/api/core/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
)

// MockKubeClientWrapper is a mock of KubeClientWrapper interface.
type MockKubeClientWrapper struct {
	ctrl     *gomock.Controller
	recorder *MockKubeClientWrapperMockRecorder
}

// MockKubeClientWrapperMockRecorder is the mock recorder for MockKubeClientWrapper.
type MockKubeClientWrapperMockRecorder struct {
	mock *MockKubeClientWrapper
}

// NewMockKubeClientWrapper creates a new mock instance.
func NewMockKubeClientWrapper(ctrl *gomock.Controller) *MockKubeClientWrapper {
	mock := &MockKubeClientWrapper{ctrl: ctrl}
	mock.recorder = &MockKubeClientWrapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKubeClientWrapper) EXPECT() *MockKubeClientWrapperMockRecorder {
	return m.recorder
}

// ClientFromEPKubeConfig mocks base method.
func (m *MockKubeClientWrapper) ClientFromEPKubeConfig(arg0 *plugins.Filecontent) (*kubernetes.Clientset, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClientFromEPKubeConfig", arg0)
	ret0, _ := ret[0].(*kubernetes.Clientset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClientFromEPKubeConfig indicates an expected call of ClientFromEPKubeConfig.
func (mr *MockKubeClientWrapperMockRecorder) ClientFromEPKubeConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClientFromEPKubeConfig", reflect.TypeOf((*MockKubeClientWrapper)(nil).ClientFromEPKubeConfig), arg0)
}

// ClientFromKubeConfig mocks base method.
func (m *MockKubeClientWrapper) ClientFromKubeConfig(arg0 string) (*kubernetes.Clientset, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClientFromKubeConfig", arg0)
	ret0, _ := ret[0].(*kubernetes.Clientset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClientFromKubeConfig indicates an expected call of ClientFromKubeConfig.
func (mr *MockKubeClientWrapperMockRecorder) ClientFromKubeConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClientFromKubeConfig", reflect.TypeOf((*MockKubeClientWrapper)(nil).ClientFromKubeConfig), arg0)
}

// CreateNamespace mocks base method.
func (m *MockKubeClientWrapper) CreateNamespace(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNamespace", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateNamespace indicates an expected call of CreateNamespace.
func (mr *MockKubeClientWrapperMockRecorder) CreateNamespace(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNamespace", reflect.TypeOf((*MockKubeClientWrapper)(nil).CreateNamespace), arg0, arg1)
}

// GetNodeList mocks base method.
func (m *MockKubeClientWrapper) GetNodeList(arg0 *plugins.Filecontent, arg1 string) (*v1.NodeList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodeList", arg0, arg1)
	ret0, _ := ret[0].(*v1.NodeList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodeList indicates an expected call of GetNodeList.
func (mr *MockKubeClientWrapperMockRecorder) GetNodeList(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodeList", reflect.TypeOf((*MockKubeClientWrapper)(nil).GetNodeList), arg0, arg1)
}

// NewClient mocks base method.
func (m *MockKubeClientWrapper) NewClient(arg0 []byte) (*kubernetes.Clientset, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewClient", arg0)
	ret0, _ := ret[0].(*kubernetes.Clientset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewClient indicates an expected call of NewClient.
func (mr *MockKubeClientWrapperMockRecorder) NewClient(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewClient", reflect.TypeOf((*MockKubeClientWrapper)(nil).NewClient), arg0)
}

// NewConfigMap mocks base method.
func (m *MockKubeClientWrapper) NewConfigMap(arg0, arg1, arg2, arg3 string) (kubeutils.ConfigMapWrapper, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewConfigMap", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(kubeutils.ConfigMapWrapper)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewConfigMap indicates an expected call of NewConfigMap.
func (mr *MockKubeClientWrapperMockRecorder) NewConfigMap(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewConfigMap", reflect.TypeOf((*MockKubeClientWrapper)(nil).NewConfigMap), arg0, arg1, arg2, arg3)
}

// NewDeployment mocks base method.
func (m *MockKubeClientWrapper) NewDeployment(arg0, arg1, arg2, arg3 string) (kubeutils.DeploymentWrapper, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewDeployment", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(kubeutils.DeploymentWrapper)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewDeployment indicates an expected call of NewDeployment.
func (mr *MockKubeClientWrapperMockRecorder) NewDeployment(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewDeployment", reflect.TypeOf((*MockKubeClientWrapper)(nil).NewDeployment), arg0, arg1, arg2, arg3)
}

// NewRestClient mocks base method.
func (m *MockKubeClientWrapper) NewRestClient(arg0 []byte) (*rest.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRestClient", arg0)
	ret0, _ := ret[0].(*rest.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewRestClient indicates an expected call of NewRestClient.
func (mr *MockKubeClientWrapperMockRecorder) NewRestClient(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRestClient", reflect.TypeOf((*MockKubeClientWrapper)(nil).NewRestClient), arg0)
}

// NewSecret mocks base method.
func (m *MockKubeClientWrapper) NewSecret(arg0, arg1, arg2, arg3 string) (kubeutils.SecretWrapper, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewSecret", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(kubeutils.SecretWrapper)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewSecret indicates an expected call of NewSecret.
func (mr *MockKubeClientWrapperMockRecorder) NewSecret(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewSecret", reflect.TypeOf((*MockKubeClientWrapper)(nil).NewSecret), arg0, arg1, arg2, arg3)
}

// RestClientFromEPKubeConfig mocks base method.
func (m *MockKubeClientWrapper) RestClientFromEPKubeConfig(arg0 *plugins.Filecontent) (*rest.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RestClientFromEPKubeConfig", arg0)
	ret0, _ := ret[0].(*rest.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RestClientFromEPKubeConfig indicates an expected call of RestClientFromEPKubeConfig.
func (mr *MockKubeClientWrapperMockRecorder) RestClientFromEPKubeConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestClientFromEPKubeConfig", reflect.TypeOf((*MockKubeClientWrapper)(nil).RestClientFromEPKubeConfig), arg0)
}

// RestConfigFromKubeConfig mocks base method.
func (m *MockKubeClientWrapper) RestConfigFromKubeConfig(arg0 string) (*rest.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RestConfigFromKubeConfig", arg0)
	ret0, _ := ret[0].(*rest.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RestConfigFromKubeConfig indicates an expected call of RestConfigFromKubeConfig.
func (mr *MockKubeClientWrapperMockRecorder) RestConfigFromKubeConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestConfigFromKubeConfig", reflect.TypeOf((*MockKubeClientWrapper)(nil).RestConfigFromKubeConfig), arg0)
}
