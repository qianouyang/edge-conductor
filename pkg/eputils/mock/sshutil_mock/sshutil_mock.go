//
//   Copyright (c) 2022 Intel Corporation.
//
//   SPDX-License-Identifier: Apache-2.0
//
//
//

// Code generated by MockGen. DO NOT EDIT.
// Source: ep/pkg/eputils (interfaces: SSHApiInterface,SSHDialInterface)

// Package mock is a generated GoMock package.
package mock

import (
	plugins "ep/pkg/api/plugins"
	fs "io/fs"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	ssh "golang.org/x/crypto/ssh"
)

// MockSSHApiInterface is a mock of SSHApiInterface interface.
type MockSSHApiInterface struct {
	ctrl     *gomock.Controller
	recorder *MockSSHApiInterfaceMockRecorder
}

// MockSSHApiInterfaceMockRecorder is the mock recorder for MockSSHApiInterface.
type MockSSHApiInterfaceMockRecorder struct {
	mock *MockSSHApiInterface
}

// NewMockSSHApiInterface creates a new mock instance.
func NewMockSSHApiInterface(ctrl *gomock.Controller) *MockSSHApiInterface {
	mock := &MockSSHApiInterface{ctrl: ctrl}
	mock.recorder = &MockSSHApiInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSSHApiInterface) EXPECT() *MockSSHApiInterfaceMockRecorder {
	return m.recorder
}

// ContainerdCertificatePathCreateSudoNoPasswd mocks base method.
func (m *MockSSHApiInterface) ContainerdCertificatePathCreateSudoNoPasswd(arg0 string, arg1 *ssh.ClientConfig, arg2, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerdCertificatePathCreateSudoNoPasswd", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// ContainerdCertificatePathCreateSudoNoPasswd indicates an expected call of ContainerdCertificatePathCreateSudoNoPasswd.
func (mr *MockSSHApiInterfaceMockRecorder) ContainerdCertificatePathCreateSudoNoPasswd(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerdCertificatePathCreateSudoNoPasswd", reflect.TypeOf((*MockSSHApiInterface)(nil).ContainerdCertificatePathCreateSudoNoPasswd), arg0, arg1, arg2, arg3)
}

// CopyLocalFileToRemoteFile mocks base method.
func (m *MockSSHApiInterface) CopyLocalFileToRemoteFile(arg0 string, arg1 *ssh.ClientConfig, arg2, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CopyLocalFileToRemoteFile", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// CopyLocalFileToRemoteFile indicates an expected call of CopyLocalFileToRemoteFile.
func (mr *MockSSHApiInterfaceMockRecorder) CopyLocalFileToRemoteFile(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CopyLocalFileToRemoteFile", reflect.TypeOf((*MockSSHApiInterface)(nil).CopyLocalFileToRemoteFile), arg0, arg1, arg2, arg3)
}

// CopyLocalFileToRemoteRootFileSudoNoPasswd mocks base method.
func (m *MockSSHApiInterface) CopyLocalFileToRemoteRootFileSudoNoPasswd(arg0 string, arg1 *ssh.ClientConfig, arg2, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CopyLocalFileToRemoteRootFileSudoNoPasswd", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// CopyLocalFileToRemoteRootFileSudoNoPasswd indicates an expected call of CopyLocalFileToRemoteRootFileSudoNoPasswd.
func (mr *MockSSHApiInterfaceMockRecorder) CopyLocalFileToRemoteRootFileSudoNoPasswd(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CopyLocalFileToRemoteRootFileSudoNoPasswd", reflect.TypeOf((*MockSSHApiInterface)(nil).CopyLocalFileToRemoteRootFileSudoNoPasswd), arg0, arg1, arg2, arg3)
}

// CopyRemoteFileToLocalFile mocks base method.
func (m *MockSSHApiInterface) CopyRemoteFileToLocalFile(arg0 string, arg1 *ssh.ClientConfig, arg2, arg3 string, arg4 fs.FileMode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CopyRemoteFileToLocalFile", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// CopyRemoteFileToLocalFile indicates an expected call of CopyRemoteFileToLocalFile.
func (mr *MockSSHApiInterfaceMockRecorder) CopyRemoteFileToLocalFile(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CopyRemoteFileToLocalFile", reflect.TypeOf((*MockSSHApiInterface)(nil).CopyRemoteFileToLocalFile), arg0, arg1, arg2, arg3, arg4)
}

// CopyRemoteRootFileToLocalFileSudoNoPasswd mocks base method.
func (m *MockSSHApiInterface) CopyRemoteRootFileToLocalFileSudoNoPasswd(arg0 string, arg1 *ssh.ClientConfig, arg2, arg3 string, arg4 fs.FileMode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CopyRemoteRootFileToLocalFileSudoNoPasswd", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// CopyRemoteRootFileToLocalFileSudoNoPasswd indicates an expected call of CopyRemoteRootFileToLocalFileSudoNoPasswd.
func (mr *MockSSHApiInterfaceMockRecorder) CopyRemoteRootFileToLocalFileSudoNoPasswd(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CopyRemoteRootFileToLocalFileSudoNoPasswd", reflect.TypeOf((*MockSSHApiInterface)(nil).CopyRemoteRootFileToLocalFileSudoNoPasswd), arg0, arg1, arg2, arg3, arg4)
}

// GenSSHConfig mocks base method.
func (m *MockSSHApiInterface) GenSSHConfig(arg0 *plugins.Node) (*ssh.ClientConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenSSHConfig", arg0)
	ret0, _ := ret[0].(*ssh.ClientConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenSSHConfig indicates an expected call of GenSSHConfig.
func (mr *MockSSHApiInterfaceMockRecorder) GenSSHConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenSSHConfig", reflect.TypeOf((*MockSSHApiInterface)(nil).GenSSHConfig), arg0)
}

// RemoteFileExists mocks base method.
func (m *MockSSHApiInterface) RemoteFileExists(arg0 string, arg1 *ssh.ClientConfig, arg2 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoteFileExists", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoteFileExists indicates an expected call of RemoteFileExists.
func (mr *MockSSHApiInterfaceMockRecorder) RemoteFileExists(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoteFileExists", reflect.TypeOf((*MockSSHApiInterface)(nil).RemoteFileExists), arg0, arg1, arg2)
}

// RunRemoteCMD mocks base method.
func (m *MockSSHApiInterface) RunRemoteCMD(arg0 string, arg1 *ssh.ClientConfig, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunRemoteCMD", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunRemoteCMD indicates an expected call of RunRemoteCMD.
func (mr *MockSSHApiInterfaceMockRecorder) RunRemoteCMD(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunRemoteCMD", reflect.TypeOf((*MockSSHApiInterface)(nil).RunRemoteCMD), arg0, arg1, arg2)
}

// RunRemoteMultiCMD mocks base method.
func (m *MockSSHApiInterface) RunRemoteMultiCMD(arg0 string, arg1 *ssh.ClientConfig, arg2 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunRemoteMultiCMD", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunRemoteMultiCMD indicates an expected call of RunRemoteMultiCMD.
func (mr *MockSSHApiInterfaceMockRecorder) RunRemoteMultiCMD(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunRemoteMultiCMD", reflect.TypeOf((*MockSSHApiInterface)(nil).RunRemoteMultiCMD), arg0, arg1, arg2)
}

// RunRemoteNodeMultiCMD mocks base method.
func (m *MockSSHApiInterface) RunRemoteNodeMultiCMD(arg0 *plugins.Node, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunRemoteNodeMultiCMD", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunRemoteNodeMultiCMD indicates an expected call of RunRemoteNodeMultiCMD.
func (mr *MockSSHApiInterfaceMockRecorder) RunRemoteNodeMultiCMD(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunRemoteNodeMultiCMD", reflect.TypeOf((*MockSSHApiInterface)(nil).RunRemoteNodeMultiCMD), arg0, arg1)
}

// ServiceRestartSudoNoPasswd mocks base method.
func (m *MockSSHApiInterface) ServiceRestartSudoNoPasswd(arg0 string, arg1 *ssh.ClientConfig, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ServiceRestartSudoNoPasswd", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ServiceRestartSudoNoPasswd indicates an expected call of ServiceRestartSudoNoPasswd.
func (mr *MockSSHApiInterfaceMockRecorder) ServiceRestartSudoNoPasswd(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServiceRestartSudoNoPasswd", reflect.TypeOf((*MockSSHApiInterface)(nil).ServiceRestartSudoNoPasswd), arg0, arg1, arg2)
}

// WriteRemoteFile mocks base method.
func (m *MockSSHApiInterface) WriteRemoteFile(arg0 string, arg1 *ssh.ClientConfig, arg2, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteRemoteFile", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteRemoteFile indicates an expected call of WriteRemoteFile.
func (mr *MockSSHApiInterfaceMockRecorder) WriteRemoteFile(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteRemoteFile", reflect.TypeOf((*MockSSHApiInterface)(nil).WriteRemoteFile), arg0, arg1, arg2, arg3)
}

// MockSSHDialInterface is a mock of SSHDialInterface interface.
type MockSSHDialInterface struct {
	ctrl     *gomock.Controller
	recorder *MockSSHDialInterfaceMockRecorder
}

// MockSSHDialInterfaceMockRecorder is the mock recorder for MockSSHDialInterface.
type MockSSHDialInterfaceMockRecorder struct {
	mock *MockSSHDialInterface
}

// NewMockSSHDialInterface creates a new mock instance.
func NewMockSSHDialInterface(ctrl *gomock.Controller) *MockSSHDialInterface {
	mock := &MockSSHDialInterface{ctrl: ctrl}
	mock.recorder = &MockSSHDialInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSSHDialInterface) EXPECT() *MockSSHDialInterfaceMockRecorder {
	return m.recorder
}

// Dial mocks base method.
func (m *MockSSHDialInterface) Dial(arg0, arg1 string, arg2 *ssh.ClientConfig) (*ssh.Client, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dial", arg0, arg1, arg2)
	ret0, _ := ret[0].(*ssh.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Dial indicates an expected call of Dial.
func (mr *MockSSHDialInterfaceMockRecorder) Dial(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dial", reflect.TypeOf((*MockSSHDialInterface)(nil).Dial), arg0, arg1, arg2)
}
