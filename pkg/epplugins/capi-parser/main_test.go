/*
* Copyright (c) 2021 Intel Corporation.
*
* SPDX-License-Identifier: Apache-2.0
*
 */

// Template auto-generated once, maintained by plugin owner.

package capiparser

import (
	"fmt"
	"testing"
)

func TestPluginMain(t *testing.T) {
	cases := []struct {
		name                  string
		input, expectedOutput map[string][]byte
		expectError           bool
		expectErrorMsg        string
	}{
		{
			name: "Ekconfigs_lost",
			input: map[string][]byte{
				"ep-params": nil,
			},
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

			expectedOutput: nil,
			expectError:    true,
			expectErrorMsg: "Failed to get CAPI infra provider config!",
		},
		{
			name: "Mutiple_Infra_provider_config",
			input: map[string][]byte{
				"ep-params": []byte(`{"ekconfig": {"Cluster": {"provider": "clusterapi"}, "Parameters": {"Extensions": ["capi-metal3", "capi-byoh"]}}}`),
				"cluster-manifest": []byte(`{
					"clusterapi":{"configs":[{"name": "metal3", "bin_url": "aaa"}]}}`),
			},
			expectedOutput: nil,
			expectError:    true,
			expectErrorMsg: "Failed to get CAPI infra provider config!",
		},
		{
			name: "CAPI_manifest_lost",
			input: map[string][]byte{
				"ep-params": []byte(`{"ekconfig": {"Cluster": {"provider": "clusterapi"}, "Parameters": {"Extensions": ["capi-metal3"]}}}`),
			},

			expectedOutput: nil,
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
			expectedOutput: map[string][]byte{
				"docker-images": []byte(`{"images": []}`),
				"files":         nil,
			},
			expectError:    true,
			expectErrorMsg: "Failed to get Cluster providers yaml list.",
		},
		{
			name: "Kind config lost",
			input: map[string][]byte{
				"ep-params": []byte(`{"ekconfig": {"Cluster": {"provider": "clusterapi"}, "Parameters": {"Extensions": ["capi-metal3"]}}}`),
				"cluster-manifest": []byte(`{
					"clusterapi":{"configs":[{"name": "metal3", "bin_url": "aaa/clusterctl", "images": ["test:test"]}], "cert-manager": {"url": "bbb/cert-manager.yaml"}}}`),
			},
			expectedOutput: map[string][]byte{
				"docker-images": []byte(`{"images": [{"url":"test:test"}]}`),
				"files":         []byte(`{"files":[{"url":"aaa/clusterctl","urlreplacement":{"new":"capi/bin","origin":"://aaa"}},{"url":"bbb/cert-manager.yaml","urlreplacement":{"new":"capi/certManager","origin":"://bbb"}}]}`),
			},
			expectError:    true,
			expectErrorMsg: "CAPI manifest kind info Lost!",
		},
		{
			name: "kind container image  lost",
			input: map[string][]byte{
				"ep-params": []byte(`{"ekconfig": {"Cluster": {"provider": "clusterapi"}, "Parameters": {"Extensions": ["capi-metal3"]}}}`),
				"cluster-manifest": []byte(`{
					"clusterapi":{"configs":[{"name": "metal3", "bin_url": "aaa/clusterctl", "images": ["test:test"]}], "cert-manager": {"url": "bbb/cert-manager.yaml"}},
					"cluster_providers":[{
						"name":"kind",
						"images":[
							{"name":"img_haproxy","repo_tag":""}
							],
						"binaries":[{"name":"kindtool","url":"","sha256":""}]
					}]}`),
			},
			expectedOutput: map[string][]byte{
				"docker-images": []byte(`{"images": [{"url":"test:test"},{"name":"kind"},{"name":"kindhaproxy"}]}`),
				"files":         []byte(`{"files":[{"url":"aaa/clusterctl","urlreplacement":{"new":"capi/bin","origin":"://aaa"}},{"url":"bbb/cert-manager.yaml","urlreplacement":{"new":"capi/cert-manager","origin":"://bbb"}},{"hashtype":"sha256","urlreplacement":{"new":"capi/kind","origin":"://."}}]}`),
			},
			expectError:    true,
			expectErrorMsg: "Failed to get management cluster binary list",
		},
		{
			name: "capi parse success",
			input: map[string][]byte{
				"ep-params": []byte(`{"ekconfig": {"Cluster": {"provider": "clusterapi"}, "Parameters": {"Extensions": ["capi-metal3"]}}}`),
				"cluster-manifest": []byte(`{
					"clusterapi":{"configs":[{"name": "metal3", "bin_url": "aaa/clusterctl", "images": ["test:test"]}], "cert-manager": {"url": "bbb/cert-manager.yaml"}},
					"cluster_providers":[{
						"name":"kind",
						"images":[
							{"name":"img_node","repo_tag":""},
							{"name":"img_haproxy","repo_tag":""}
							],
						"binaries":[{"name":"kindtool","url":"","sha256":""}]
					}]}`),
			},
			expectedOutput: map[string][]byte{
				"docker-images": []byte(`{"images": [{"url":"test:test"},{"name":"kind"},{"name":"kindhaproxy"}]}`),
				"files":         []byte(`{"files":[{"url":"aaa/clusterctl","urlreplacement":{"new":"capi/bin","origin":"://aaa"}},{"url":"bbb/cert-manager.yaml","urlreplacement":{"new":"capi/cert-manager","origin":"://bbb"}},{"hashtype":"sha256","urlreplacement":{"new":"capi/kind","origin":"://."}}]}`),
			},
			expectError:    false,
			expectErrorMsg: "",
		},
	}

	// Optional: add setup for the test series
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			input := generateInput(tc.input)
			if input == nil {
				t.Fatalf("Failed to generateInput %s", tc.input)
			}
			testOutput := generateOutput(nil)

			expectedOutput := generateOutput(tc.expectedOutput)

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

			if testOutput.EqualWith(expectedOutput) {
				t.Log("Output expected.")
			} else {
				testOstr, _ := testOutput.MarshalBinary()
				expectOstr, _ := expectedOutput.MarshalBinary()
				t.Logf("Output is %s, Expectoutput is %s", testOstr, expectOstr)
				t.Errorf("Failed to get expected output when input is %s.", tc.input)
			}

			// Optional: Add additional check conditions here

		})
	}

	// Optional: add teardown for the test series
}
