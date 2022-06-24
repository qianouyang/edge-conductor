/*
* Copyright (c) 2021 Intel Corporation.
*
* SPDX-License-Identifier: Apache-2.0
*
 */

// Template auto-generated once, maintained by plugin owner.

package capiproviderlaunch

import (
	eputils "ep/pkg/eputils"
	repoutils "ep/pkg/eputils/repoutils"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	mpatch "github.com/undefinedlabs/go-mpatch"
)

var (
	errGeneral = errors.New("error")
)

func TestPluginMain(t *testing.T) {
	func_patch_os_file := func() []*mpatch.Patch {
		patch1, _ := mpatch.PatchMethod(os.Chmod, func(string, os.FileMode) error { return nil })
		patch2, _ := mpatch.PatchMethod(eputils.WriteStringToFile, func(string, string) error { return nil })
		patch3, _ := mpatch.PatchMethod(eputils.FileTemplateConvert, func(string, string) error { return nil })
		patch4, _ := mpatch.PatchMethod(eputils.RunCMD, func(*exec.Cmd) (string, error) { return "", nil })
		patch5, _ := mpatch.PatchMethod(eputils.CreateFolderIfNotExist, func(string) error { return nil })
		return []*mpatch.Patch{patch1, patch2, patch3, patch4, patch5}
	}
	func_patch_create_folder := func() []*mpatch.Patch {
		patch1, _ := mpatch.PatchMethod(eputils.CreateFolderIfNotExist, func(string) error { return nil })
		return []*mpatch.Patch{patch1}
	}

	func_failed_pull_clusterctl := func() []*mpatch.Patch {
		patches := func_patch_os_file()
		patch, _ := mpatch.PatchMethod(repoutils.PullFileFromRepo, func(file string, url string) error {
			if strings.Contains(file, "clusterctl") {
				return errGeneral
			} else {
				return nil
			}
		})
		patches = append(patches, patch)
		return patches
	}

	cases := []struct {
		name           string
		input          map[string][]byte
		expectError    bool
		expectErrorMsg string
		funcBeforeTest func() []*mpatch.Patch
	}{
		{
			name: "Ekconfigs_lost",
			input: map[string][]byte{
				"ep-params": nil,
			},
			funcBeforeTest: func_patch_create_folder,
			expectError:    true,
			expectErrorMsg: "Incorrect parameter",
		},
		{
			name: "Infra_provider_lost_in_EK_config",
			input: map[string][]byte{
				"ep-params": []byte(`{"ekconfig": {"Cluster": {"provider": "clusterapi"}}}`),
				"cluster-manifest": []byte(`{
					"clusterapi":{"configs":[{"name": "metal3", "bin_url": "aaa"}]}}`),
			},
			funcBeforeTest: func_patch_create_folder,
			expectError:    true,
			expectErrorMsg: "Failed to get CAPI infra provider config!",
		},
		{
			name: "CAPI_manifest_lost",
			input: map[string][]byte{
				"ep-params": []byte(`{"ekconfig": {"Cluster": {"provider": "clusterapi"}, "Parameters": {"Extensions": ["capi-metal3"]}}}`),
			},
			funcBeforeTest: func_patch_create_folder,
			expectError:    true,
			expectErrorMsg: "CAPI manifest Lost!",
		},
		{
			name: "CertManage_config_lost",
			input: map[string][]byte{
				"ep-params": []byte(`{"ekconfig": {"Cluster": {"provider": "clusterapi"}, "Parameters": {"Extensions": ["capi-metal3"]}}}`),
				"cluster-manifest": []byte(`{
					"clusterapi":{"configs":[{"name": "metal3", "bin_url": "aaa", "images": ["test:test"]}]}}`),
			},
			funcBeforeTest: func_patch_create_folder,
			expectError:    true,
			expectErrorMsg: "CAPI provider config invalidate",
		},
		{
			name: "Provider lost",
			input: map[string][]byte{
				"ep-params": []byte(`{"ekconfig": {"Cluster": {"provider": "clusterapi"}, "Parameters": {"Extensions": ["capi-metal3"]}}}`),
				"cluster-manifest": []byte(`{
					"clusterapi":{"configs":[{"name": "metal3", "bin_url": "aaa", "images": ["test:test"]}], "cert-manager": {"url": "bbb/cert-manager.yaml"}}}`),
			},
			funcBeforeTest: func_patch_create_folder,
			expectError:    true,
			expectErrorMsg: "CAPI provider config invalidate",
		},
		{
			name: "config in manefest lost",
			input: map[string][]byte{
				"ep-params": []byte(`{"ekconfig": {"Cluster": {"provider": "clusterapi"}, "Parameters": {"Extensions": ["capi-metal3"]}}}`),
				"cluster-manifest": []byte(`{
					"clusterapi":{ "cert-manager": {"url": "bbb/cert-manager.yaml"}}}`),
			},
			funcBeforeTest: func_patch_create_folder,
			expectError:    true,
			expectErrorMsg: "CAPI provider config invalidate",
		},
		{
			name: "Provider parameter lost.",
			input: map[string][]byte{
				"ep-params": []byte(`{"ekconfig": {"Cluster": {"provider": "clusterapi"}, "Parameters": {"Extensions": ["capi-metal3"]}}}`),
				"cluster-manifest": []byte(`{
					"clusterapi":{"configs":[{"name": "metal3", "bin_url": "aaa", "images": ["test:test"], "providers": [{"provider_type": "CoreProvider", "name": "" }]}], "cert-manager": {"url": "bbb/cert-manager.yaml"}}}`),
			},
			funcBeforeTest: func_patch_create_folder,
			expectError:    true,
			expectErrorMsg: "CAPI provider config invalidate",
		},
		{
			name: "Provider config lost.",
			input: map[string][]byte{
				"ep-params": []byte(`{"ekconfig": {"Cluster": {"provider": "clusterapi"}, "Parameters": {"Extensions": ["capi-metal3"]}}}`),
				"cluster-manifest": []byte(`{
					"clusterapi":{"configs":[{"name": "metal3", "bin_url": "aaa", "images": ["test:test"], "providers": [{"provider_type": "CoreProvider", "name": "", "parameters" : {"version": "", "provider_label": ""}}]}], "cert-manager": {"url": "bbb/cert-manager.yaml"}}}`),
			},
			funcBeforeTest: func_patch_create_folder,
			expectError:    true,
			expectErrorMsg: "CAPI provider config invalidate",
		},
		{
			name: "Failed to pull files from registry",
			input: map[string][]byte{
				"ep-params": []byte(`{"ekconfig": {"Cluster": {"provider": "clusterapi"}, "Parameters": {"Extensions": ["capi-metal3"]}}}`),
				"cluster-manifest": []byte(`{
					"clusterapi":{"configs":[{"name": "metal3", "bin_url": "aaa", "images": ["test:test"], "providers": [{"provider_type": "CoreProvider", "name": "", "parameters" : {"version": "", "provider_label": ""}}, {"provider_type": "BootstrapProvider", "name": "", "parameters" : {"version": "", "provider_label": ""}},{"provider_type": "ControlPlaneProvider", "name": "", "parameters" : {"version": "", "provider_label": ""}},{"provider_type": "InfrastructureProvider", "name": "", "parameters" : {"version": "", "provider_label": ""}}]}], "cert-manager": {"url": "bbb/cert-manager.yaml"}}}`),
				"files": []byte(`{"files":[{"url":"core-provider","mirrorurl": "oci:/cluster-api/core-provider"}]}`),
			},
			funcBeforeTest: func_patch_create_folder,
			expectError:    true,
			expectErrorMsg: "Failed to generate local provider repo for clusterctl!",
		},
		{
			name: "Failed to create management cluster",
			input: map[string][]byte{
				"ep-params": []byte(`{"ekconfig": {"Cluster": {"provider": "clusterapi"}, "Parameters": {"Extensions": ["capi-metal3"]}}}`),
				"cluster-manifest": []byte(`{
					"clusterapi":{"configs":[{"name": "metal3", "bin_url": "aaa", "images": ["test:test"], "providers": [{"provider_type": "CoreProvider", "name": "", "parameters" : {"version": "", "provider_label": ""}}, {"provider_type": "BootstrapProvider", "name": "", "parameters" : {"version": "", "provider_label": ""}},{"provider_type": "ControlPlaneProvider", "name": "", "parameters" : {"version": "", "provider_label": ""}},{"provider_type": "InfrastructureProvider", "name": "", "parameters" : {"version": "", "provider_label": ""}}]}], "cert-manager": {"url": "bbb/cert-manager.yaml"}}}`),
				"files": []byte(`{"files":[{"url":"core-provider","mirrorurl": "/cluster-api/core-provider"}]}`),
			},
			funcBeforeTest: func_patch_os_file,
			expectError:    true,
			expectErrorMsg: "Failed to launch management cluster",
		},
		{
			name: "config of ekconfig missing",
			input: map[string][]byte{
				"ep-params": []byte(`{"ekconfig": {"Cluster": {"provider": "clusterapi"}, "Parameters": {"extensions": ["capi-metal3"]}}}`),
				"cluster-manifest": []byte(`{
					"clusterapi":{"configs":[{"name": "metal3", "bin_url": "aaa", "images": ["test:test"], "providers": [{"provider_type": "CoreProvider", "name": "", "parameters" : {"version": "", "provider_label": ""}}, {"provider_type": "BootstrapProvider", "name": "", "parameters" : {"version": "", "provider_label": ""}},{"provider_type": "ControlPlaneProvider", "name": "", "parameters" : {"version": "", "provider_label": ""}},{"provider_type": "InfrastructureProvider", "name": "", "parameters" : {"version": "", "provider_label": ""}}]}], "cert-manager": {"url": "bbb/cert-manager.yaml"}},
					"cluster_providers":[{"name":"kind","images":[{"name":"img_node","repo_tag":""},{"name":"img_haproxy","repo_tag":""}],"binaries":[{"name":"kindtool","url":"","sha256":""}]}]}`),
				"files": []byte(`{"files":[{"url":"core-provider","mirrorurl": "/cluster-api/core-provider"}, {"mirrorurl": "capi/kind"}]}`),
			},
			funcBeforeTest: func_patch_os_file,
			expectError:    true,
			expectErrorMsg: "Failed to launch management cluster",
		},
		{
			name: "Failed to init clusterctl",
			input: map[string][]byte{
				"ep-params": []byte(`{"ekconfig": {"Cluster": {"provider": "clusterapi"}, "Parameters": {"extensions": ["capi-metal3"], "global_settings": {"provider_ip": "", "registry_port": ""}}}}`),
				"cluster-manifest": []byte(`{
					"clusterapi":{"configs":[{"name": "metal3", "bin_url": "aaa", "images": ["test:test"], "providers": [{"provider_type": "CoreProvider", "name": "", "parameters" : {"version": "", "provider_label": ""}}, {"provider_type": "BootstrapProvider", "name": "", "parameters" : {"version": "", "provider_label": ""}},{"provider_type": "ControlPlaneProvider", "name": "", "parameters" : {"version": "", "provider_label": ""}},{"provider_type": "InfrastructureProvider", "name": "", "parameters" : {"version": "", "provider_label": ""}}]}], "cert-manager": {"url": "bbb/cert-manager.yaml"}},
					"cluster_providers":[{"name":"kind","images":[{"name":"img_node","repo_tag":""},{"name":"img_haproxy","repo_tag":""}],"binaries":[{"name":"kindtool","url":"","sha256":""}]}]}`),
				"files": []byte(`{"files":[{"url":"core-provider","mirrorurl": "/cluster-api/core-provider"}, {"mirrorurl": "capi/kind"}]}`),
			},
			funcBeforeTest: func_failed_pull_clusterctl,
			expectError:    true,
			expectErrorMsg: "Failed to init clusterctl!",
		},
		{
			name: "provider launch success",
			input: map[string][]byte{
				"ep-params": []byte(`{"ekconfig": {"Cluster": {"provider": "clusterapi"}, "Parameters": {"extensions": ["capi-metal3"], "global_settings": {"provider_ip": "", "registry_port": ""}}}, "extensions": [{"name": "capi-metal3", "extension": {"extension": [{"name": "Infra-provider", "config": [{"name": "Management-cluster-kubeconfig", "value": ""}]}]}}]}`),
				"cluster-manifest": []byte(`{
					"clusterapi":{"configs":[{"name": "metal3", "bin_url": "aaa", "images": ["test:test"], "providers": [{"provider_type": "CoreProvider", "name": "", "parameters" : {"version": "", "provider_label": "core"}}, {"provider_type": "BootstrapProvider", "name": "", "parameters" : {"version": "", "provider_label": "boot"}},{"provider_type": "ControlPlaneProvider", "name": "", "parameters" : {"version": "", "provider_label": "control"}},{"provider_type": "InfrastructureProvider", "name": "", "parameters" : {"version": "", "provider_label": "infra"}}]}], "cert-manager": {"url": "bbb/cert-manager.yaml"}},
					"cluster_providers":[{"name":"kind","images":[{"name":"img_node","repo_tag":""},{"name":"img_haproxy","repo_tag":""}],"binaries":[{"name":"kindtool","url":"","sha256":""}]}]}`),
				"files": []byte(`{"files":[{"url":"core-provider","mirrorurl": "/cluster-api/core-provider"}, {"mirrorurl": "capi/kind"}, {"url":"kubeadm","mirrorurl": "/bootstrap-kubeadm/kubeadm"}, {"url":"kubeadm","mirrorurl": "/control-plane-kubeadm/kubeadm"}, {"url":"metal3","mirrorurl": "/infrastructure-metal3/metal3"}]}`),
			},
			funcBeforeTest: func_patch_os_file,
			expectError:    false,
			expectErrorMsg: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			input := generateInput(tc.input)
			if input == nil {
				t.Fatalf("Failed to generateInput %s", tc.input)
			}

			if tc.funcBeforeTest != nil {
				plist := tc.funcBeforeTest()
				for _, p := range plist {
					defer func(t *testing.T, m *mpatch.Patch) {
						err := m.Unpatch()
						if err != nil {
							t.Fatal(err)
						}
					}(t, p)
				}
			}

			testOutput := generateOutput(nil)

			if result := PluginMain(input, &testOutput); result != nil {
				if tc.expectError {
					if fmt.Sprint(result) == tc.expectErrorMsg {
						t.Logf("Expected error: {%s} catched, done.", tc.expectErrorMsg)
						return
					} else if tc.expectErrorMsg == "" {
						return
					} else {
						t.Fatal("Unexpected error occurred.")
					}

				}
				t.Logf("Failed to run PluginMain when input is %s.", tc.input)
				t.Error(result)
			}

			_ = testOutput
		})
	}
}
