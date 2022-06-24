/*
* Copyright (c) 2022 Intel Corporation.
*
* SPDX-License-Identifier: Apache-2.0
*
 */
package main

import (
	plugin "ep/pkg/plugin"
	"flag"
	"fmt"
	"os"
	"reflect"
	"testing"

	mpatch "github.com/undefinedlabs/go-mpatch"
)

var (
	errTest = fmt.Errorf("test error")
)

func unpatch(t *testing.T, m *mpatch.Patch) {
	err := m.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
}

func unpatchAll(t *testing.T, pList []*mpatch.Patch) {
	for _, p := range pList {
		if p != nil {
			if err := p.Unpatch(); err != nil {
				t.Errorf("unpatch error: %v", err)
			}
		}
	}
}

func patchOsExit(t *testing.T) {
	var exitGuard *mpatch.Patch
	exitGuard, _ = mpatch.PatchMethod(os.Exit, func(code int) {
		unpatch(t, exitGuard)
		panic(code)
	})
}

func patchFlagArgs(t *testing.T, args []string) *mpatch.Patch {
	patch, patchErr := mpatch.PatchMethod(flag.Args, func() []string {
		return args
	})
	if patchErr != nil {
		t.Errorf("patch error: %v", patchErr)
		return nil
	}
	return patch
}

func patchEnablePluginRemoteLog(t *testing.T, err error) *mpatch.Patch {
	patch, patchErr := mpatch.PatchMethod(plugin.EnablePluginRemoteLog, func(name string) error {
		return err
	})
	if patchErr != nil {
		t.Errorf("patch error: %v", patchErr)
		return nil
	}
	return patch
}

func patchStartPlugin(t *testing.T, err error) *mpatch.Patch {
	patch, patchErr := mpatch.PatchMethod(plugin.StartPlugin, func(name string, errch chan error) error {
		return err
	})
	if patchErr != nil {
		t.Errorf("patch error: %v", patchErr)
		return nil
	}
	return patch
}

func patchWaitPluginFinished(t *testing.T, err error) *mpatch.Patch {
	patch, patchErr := mpatch.PatchMethod(plugin.WaitPluginFinished, func(name string) error {
		return err
	})
	if patchErr != nil {
		t.Errorf("patch error: %v", patchErr)
		return nil
	}
	return patch
}

func TestMainFunction(t *testing.T) {
	cases := []struct {
		funcBeforeTest func() []*mpatch.Patch
		wantExitCode   interface{}
	}{
		{
			funcBeforeTest: func() []*mpatch.Patch {
				patchOsExit(t)
				patch := patchFlagArgs(t, make([]string, 1000))
				return []*mpatch.Patch{patch}
			},
			wantExitCode: 1,
		},
		{
			funcBeforeTest: func() []*mpatch.Patch {
				patch := patchFlagArgs(t, []string{})
				return []*mpatch.Patch{patch}
			},
			wantExitCode: nil,
		},
		{
			funcBeforeTest: func() []*mpatch.Patch {
				patchOsExit(t)
				patchFlagArgs := patchFlagArgs(t, []string{"aaa"})
				patchEnablePluginRemoteLog := patchEnablePluginRemoteLog(t, errTest)
				return []*mpatch.Patch{patchFlagArgs, patchEnablePluginRemoteLog}
			},
			wantExitCode: 1,
		},
		{
			funcBeforeTest: func() []*mpatch.Patch {
				patchOsExit(t)
				patchFlagArgs := patchFlagArgs(t, []string{"aaa"})
				patchEnablePluginRemoteLog := patchEnablePluginRemoteLog(t, nil)
				patchStartPlugin := patchStartPlugin(t, errTest)
				return []*mpatch.Patch{patchFlagArgs, patchEnablePluginRemoteLog, patchStartPlugin}
			},
			wantExitCode: 1,
		},
		{
			funcBeforeTest: func() []*mpatch.Patch {
				patchOsExit(t)
				patchFlagArgs := patchFlagArgs(t, []string{"aaa"})
				patchEnablePluginRemoteLog := patchEnablePluginRemoteLog(t, nil)
				patchStartPlugin := patchStartPlugin(t, nil)
				patchWaitPluginFinished := patchWaitPluginFinished(t, errTest)
				return []*mpatch.Patch{patchFlagArgs, patchEnablePluginRemoteLog, patchStartPlugin, patchWaitPluginFinished}
			},
			wantExitCode: 1,
		},
		{
			funcBeforeTest: func() []*mpatch.Patch {
				patchFlagArgs := patchFlagArgs(t, []string{"aaa"})
				patchEnablePluginRemoteLog := patchEnablePluginRemoteLog(t, nil)
				patchStartPlugin := patchStartPlugin(t, nil)
				patchWaitPluginFinished := patchWaitPluginFinished(t, nil)
				return []*mpatch.Patch{patchFlagArgs, patchEnablePluginRemoteLog, patchStartPlugin, patchWaitPluginFinished}
			},
			wantExitCode: nil,
		},
	}

	for n, testCase := range cases {
		t.Logf("TestMainFunction case %d start", n)
		func() {
			if testCase.funcBeforeTest != nil {
				pList := testCase.funcBeforeTest()
				defer unpatchAll(t, pList)
			}
			defer func() {
				exitCode := recover()
				if !reflect.DeepEqual(exitCode, testCase.wantExitCode) {
					t.Errorf("Unexpected exit code: %v", exitCode)
				}
			}()
			main()
		}()
		t.Logf("TestMainFunction case %d End", n)
	}
	t.Log("Done")
}
