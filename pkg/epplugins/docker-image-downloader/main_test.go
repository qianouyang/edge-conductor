/*
* Copyright (c) 2022 Intel Corporation.
*
* SPDX-License-Identifier: Apache-2.0
*
 */

// Template auto-generated once, maintained by plugin owner.

//nolint: dupl
package dockerimagedownloader

import (
	docker "ep/pkg/eputils/docker"
	dockermock "ep/pkg/eputils/docker/mock"
	restfulcli "ep/pkg/eputils/restfulcli"
	restfulmock "ep/pkg/eputils/restfulcli/mock"
	"errors"
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	mpatch "github.com/undefinedlabs/go-mpatch"
)

var (
	errTest = errors.New("testErr")
)

func unpatch(t *testing.T, m *mpatch.Patch) {
	err := m.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
}

func TestPluginMain(t *testing.T) {
	func_Imagedownloader_successful := func(ctrl *gomock.Controller, ctrl1 *gomock.Controller) []*mpatch.Patch {
		mockDockerCli := dockermock.NewMockDockerClientWrapperImage(ctrl)
		mockRestyCli := restfulmock.NewMockGoharborClientWrapper(ctrl1)
		patchImagePull, err := mpatch.PatchMethod(docker.ImagePull, mockDockerCli.ImagePull)
		if err != nil {
			t.Fatal(err)
		}
		patchMapImageURLCreateHarborProject, err := mpatch.PatchMethod(restfulcli.MapImageURLCreateHarborProject, mockRestyCli.MapImageURLCreateHarborProject)
		if err != nil {
			t.Fatal(err)
		}
		patchTagImageToLocal, err := mpatch.PatchMethod(docker.TagImageToLocal, mockDockerCli.TagImageToLocal)
		if err != nil {
			t.Fatal(err)
		}
		patchImagePush, err := mpatch.PatchMethod(docker.ImagePush, mockDockerCli.ImagePush)
		if err != nil {
			t.Fatal(err)
		}
		mockDockerCli.EXPECT().ImagePull(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
		fakeNewImages := []string{"aaa", "bbb"}
		mockRestyCli.EXPECT().MapImageURLCreateHarborProject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fakeNewImages, nil)
		mockDockerCli.EXPECT().TagImageToLocal(gomock.Any(), gomock.Any()).AnyTimes().Return("", nil)
		mockDockerCli.EXPECT().ImagePush(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
		return []*mpatch.Patch{patchImagePull, patchMapImageURLCreateHarborProject, patchTagImageToLocal, patchImagePush}
	}

	func_GetAuthConf_err := func(ctrl *gomock.Controller, ctrl1 *gomock.Controller) []*mpatch.Patch {
		mockDockerCli := dockermock.NewMockDockerClientWrapperImage(ctrl)
		patchGetAuthConf, err := mpatch.PatchMethod(docker.GetAuthConf, mockDockerCli.GetAuthConf)
		if err != nil {
			t.Fatal(err)
		}
		mockDockerCli.EXPECT().GetAuthConf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, errTest)
		return []*mpatch.Patch{patchGetAuthConf}
	}
	func_ImagePull_err := func(ctrl *gomock.Controller, ctrl1 *gomock.Controller) []*mpatch.Patch {
		mockDockerCli := dockermock.NewMockDockerClientWrapperImage(ctrl)
		patchImagePull, err := mpatch.PatchMethod(docker.ImagePull, mockDockerCli.ImagePull)
		if err != nil {
			t.Fatal(err)
		}
		mockDockerCli.EXPECT().ImagePull(gomock.Any(), gomock.Any()).AnyTimes().Return(errTest)
		return []*mpatch.Patch{patchImagePull}
	}
	func_HarborProject_err := func(ctrl *gomock.Controller, ctrl1 *gomock.Controller) []*mpatch.Patch {
		mockDockerCli := dockermock.NewMockDockerClientWrapperImage(ctrl)
		mockRestyCli := restfulmock.NewMockGoharborClientWrapper(ctrl1)
		patchImagePull, err := mpatch.PatchMethod(docker.ImagePull, mockDockerCli.ImagePull)
		if err != nil {
			t.Fatal(err)
		}
		patchMapImageURLCreateHarborProject, err := mpatch.PatchMethod(restfulcli.MapImageURLCreateHarborProject, mockRestyCli.MapImageURLCreateHarborProject)
		if err != nil {
			t.Fatal(err)
		}
		mockDockerCli.EXPECT().ImagePull(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
		mockRestyCli.EXPECT().MapImageURLCreateHarborProject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, errTest)
		return []*mpatch.Patch{patchImagePull, patchMapImageURLCreateHarborProject}
	}

	func_TagImageToLocal_err := func(ctrl *gomock.Controller, ctrl1 *gomock.Controller) []*mpatch.Patch {
		mockDockerCli := dockermock.NewMockDockerClientWrapperImage(ctrl)
		mockRestyCli := restfulmock.NewMockGoharborClientWrapper(ctrl1)

		patchImagePull, err := mpatch.PatchMethod(docker.ImagePull, mockDockerCli.ImagePull)
		if err != nil {
			t.Fatal(err)
		}
		patchMapImageURLCreateHarborProject, err := mpatch.PatchMethod(restfulcli.MapImageURLCreateHarborProject, mockRestyCli.MapImageURLCreateHarborProject)
		if err != nil {
			t.Fatal(err)
		}
		patchTagImageToLocal, err := mpatch.PatchMethod(docker.TagImageToLocal, mockDockerCli.TagImageToLocal)
		if err != nil {
			t.Fatal(err)
		}
		mockDockerCli.EXPECT().ImagePull(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
		fakeNewImages := []string{"aaa", "bbb"}
		mockRestyCli.EXPECT().MapImageURLCreateHarborProject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fakeNewImages, nil)
		mockDockerCli.EXPECT().TagImageToLocal(gomock.Any(), gomock.Any()).AnyTimes().Return("", errTest)
		return []*mpatch.Patch{patchImagePull, patchMapImageURLCreateHarborProject, patchTagImageToLocal}
	}
	func_ImagePush_err := func(ctrl *gomock.Controller, ctrl1 *gomock.Controller) []*mpatch.Patch {
		mockDockerCli := dockermock.NewMockDockerClientWrapperImage(ctrl)
		mockRestyCli := restfulmock.NewMockGoharborClientWrapper(ctrl1)
		patchImagePull, err := mpatch.PatchMethod(docker.ImagePull, mockDockerCli.ImagePull)
		if err != nil {
			t.Fatal(err)
		}
		patchMapImageURLCreateHarborProject, err := mpatch.PatchMethod(restfulcli.MapImageURLCreateHarborProject, mockRestyCli.MapImageURLCreateHarborProject)
		if err != nil {
			t.Fatal(err)
		}
		patchTagImageToLocal, err := mpatch.PatchMethod(docker.TagImageToLocal, mockDockerCli.TagImageToLocal)
		if err != nil {
			t.Fatal(err)
		}
		patchImagePush, err := mpatch.PatchMethod(docker.ImagePush, mockDockerCli.ImagePush)
		if err != nil {
			t.Fatal(err)
		}
		mockDockerCli.EXPECT().ImagePull(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
		fakeNewImages := []string{"aaa", "bbb"}
		mockRestyCli.EXPECT().MapImageURLCreateHarborProject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fakeNewImages, nil)
		mockDockerCli.EXPECT().TagImageToLocal(gomock.Any(), gomock.Any()).AnyTimes().Return("", nil)
		mockDockerCli.EXPECT().ImagePush(gomock.Any(), gomock.Any()).AnyTimes().Return(errTest)
		return []*mpatch.Patch{patchImagePull, patchMapImageURLCreateHarborProject, patchTagImageToLocal, patchImagePush}
	}
	cases := []struct {
		name                  string
		input, expectedOutput map[string][]byte
		expectError           bool
		wantErr               error
		funcBeforeTest        func(*gomock.Controller, *gomock.Controller) []*mpatch.Patch
	}{
		{
			name: "Epparams_ekconfig_wrong_0",
			input: map[string][]byte{
				"ep-params": []byte(`{
					"ekconfig": {}}`),
				"docker-images": []byte(`{"images": []}`),
			},
			expectError: true,
			wantErr:     nil,
		},
		{
			name: "Epparams_ekconfig_wrong_1",
			input: map[string][]byte{
				"ep-params": []byte(`{
					"ekconfig": {}}`),
				"docker-images": []byte(`{"images": [{"name": "test","url": "hello-world:latest"}]}`),
			},
			expectError: true,
			wantErr:     errEpparamsEkconfig,
		},
		{
			name: "Missing_registry_info",
			input: map[string][]byte{
				"ep-params": []byte(`{
					"ekconfig": {
						"Parameters": {
							"customconfig": {"registry": {"user": "test","password": "test123"}},
							"global_settings": {"provider_ip": "","registry_port": ""}}}}`),
				"docker-images": []byte(`{"images": [{"name": "test","url": "hello-world:latest"}]}`),
			},
			expectError:    true,
			wantErr:        errTest,
			funcBeforeTest: func_GetAuthConf_err,
		},
		{
			name: "Docker_pull_images_err",
			input: map[string][]byte{
				"ep-params": []byte(`{
					"ekconfig": {
						"Parameters": {
							"customconfig": {"registry": {"user": "test","password": "test123"}},
							"global_settings": {"provider_ip": "10.10.10.10","registry_port": "5678"}}}}`),
				"docker-images": []byte(`{"images": [{"name": "test","url": "hello-world:latest"}]}`),
			},
			expectError:    true,
			wantErr:        errTest,
			funcBeforeTest: func_ImagePull_err,
		},
		{
			name: "Create_HarborProject_err",
			input: map[string][]byte{
				"ep-params": []byte(`{
					"ekconfig": {
						"Parameters": {
							"customconfig": {"registry": {"user": "test","password": "test123"}},
							"global_settings": {"provider_ip": "10.10.10.10","registry_port": "5678"}}}}`),
				"docker-images": []byte(`{"images": [{"name": "test","url": "temp/hello-world:latest"}]}`),
			},
			expectError:    true,
			wantErr:        errTest,
			funcBeforeTest: func_HarborProject_err,
		},
		{
			name: "Tag_ImageToLocal_err",
			input: map[string][]byte{
				"ep-params": []byte(`{
					"ekconfig": {
						"Parameters": {
							"customconfig": {"registry": {"user": "test","password": "test123"}},
							"global_settings": {"provider_ip": "10.10.10.10","registry_port": "5678"}}}}`),
				"docker-images": []byte(`{"images": [{"name": "test","url": "temp/hello-world:latest"}]}`),
			},
			expectError:    true,
			wantErr:        errTest,
			funcBeforeTest: func_TagImageToLocal_err,
		},
		{
			name: "Docker_push_images_err",
			input: map[string][]byte{
				"ep-params": []byte(`{
					"ekconfig": {
						"Parameters": {
							"customconfig": {"registry": {"user": "test","password": "test123"}},
							"global_settings": {"provider_ip": "10.10.10.10","registry_port": "5678"}}}}`),
				"docker-images": []byte(`{"images": [{"name": "test","url": "temp/hello-world:latest"}]}`),
			},
			expectError:    true,
			wantErr:        errTest,
			funcBeforeTest: func_ImagePush_err,
		},
		{
			name: "Success_Imagedownloader_ok",
			input: map[string][]byte{
				"ep-params": []byte(`{
					"ekconfig": {
						"Parameters": {
							"customconfig": {"registry": {"user": "test","password": "test123"}},
							"global_settings": {"provider_ip": "10.10.10.10","registry_port": "5678"}}}}`),
				"docker-images": []byte(`{"images": [{"name": "test","url": "temp/hello-world:latest"}]}`),
			},
			expectError:    false,
			wantErr:        errTest,
			funcBeforeTest: func_Imagedownloader_successful,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ctrl1 := gomock.NewController(t)
			defer ctrl1.Finish()

			var plist []*mpatch.Patch
			if tc.funcBeforeTest != nil {
				plist = tc.funcBeforeTest(ctrl, ctrl1)
				for _, p := range plist {
					defer unpatch(t, p)
				}
			}
			input := generateInput(tc.input)
			if input == nil {
				t.Fatalf("Failed to generateInput %s", tc.input)
			}
			testOutput := generateOutput(tc.expectedOutput)
			result := PluginMain(input, &testOutput)

			if tc.expectError {
				if fmt.Sprint(result) == fmt.Sprint(tc.wantErr) {
					t.Log("Error expected.")
				} else {
					t.Error("Expect:", tc.expectError, "; But found:", result)
				}
			} else {
				if result != nil {
					t.Error("Unexpected Error:", result)
				}
			}

			_ = testOutput
		})
	}

}
