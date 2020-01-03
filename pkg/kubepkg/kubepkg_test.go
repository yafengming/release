/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kubepkg

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPackageVersionSuccess(t *testing.T) {
	testcases := []struct {
		name        string
		packageName string
		version     string
		kubeVersion string
		expected    string
	}{
		{
			name:        "Kubernetes version supplied",
			kubeVersion: "1.17.0",
			expected:    "1.17.0",
		},
		{
			name:     "Kubernetes version not supplied",
			expected: "",
		},
		{
			name:        "CNI version",
			packageName: "kubernetes-cni",
			version:     "0.8.3",
			kubeVersion: "1.17.0",
			expected:    "0.8.3",
		},
		{
			name:        "CRI tools version",
			packageName: "cri-tools",
			kubeVersion: "1.17.0",
			expected:    "1.17.0",
		},
	}

	for _, tc := range testcases {
		actual, err := getPackageVersion(
			&PackageDefinition{
				Name:              tc.packageName,
				Version:           tc.version,
				KubernetesVersion: tc.kubeVersion,
			},
		)

		if err != nil {
			t.Fatalf("did not expect an error: %v", err)
		}

		assert.Equal(t, tc.expected, actual)
	}
}

func TestGetPackageVersionFailure(t *testing.T) {
	a := assert.New(t)

	_, err := getPackageVersion(nil)
	a.Error(err)
}

//nolint:godox
// TODO: Figure out how we want to test success of this function.
//       When channel type is provided, we return a func() (string, error), instead of (string, error).
//       Additionally, those functions have variable output depending on when we run the test cases.
func TestGetKubernetesVersionSuccess(t *testing.T) {
	testcases := []struct {
		name        string
		version     string
		kubeVersion string
		channel     ChannelType
		expected    string
	}{
		{
			name:        "Kubernetes version supplied",
			kubeVersion: "1.17.0",
			expected:    "1.17.0",
		},
	}

	for _, tc := range testcases {
		actual, err := getKubernetesVersion(
			&PackageDefinition{
				Version:           tc.version,
				KubernetesVersion: tc.kubeVersion,
			},
		)

		if err != nil {
			t.Fatalf("did not expect an error: %v", err)
		}

		assert.Equal(t, tc.expected, actual)
	}
}

func TestGetKubernetesVersionFailure(t *testing.T) {
	a := assert.New(t)

	_, err := getKubernetesVersion(nil)
	a.Error(err)
}

func TestFetchVersionSuccess(t *testing.T) {
	testcases := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "Release URL",
			url:      "https://dl.k8s.io/release/stable-1.14.txt",
			expected: "1.14.10",
		},
		{
			name:     "CI URL",
			url:      "https://dl.k8s.io/ci/latest-1.14.txt",
			expected: "1.14.11-beta.1",
		},
	}

	for _, tc := range testcases {
		actual, err := fetchVersion(tc.url)

		if err != nil {
			t.Fatalf("did not expect an error: %v", err)
		}

		assert.Equal(t, tc.expected, actual)
	}
}

func TestFetchVersionFailure(t *testing.T) {
	testcases := []struct {
		name string
		url  string
	}{
		{
			name: "Empty URL string",
			url:  "",
		},
		{
			name: "Bad URL",
			url:  "https://fake.url",
		},
	}

	for _, tc := range testcases {
		_, err := fetchVersion(tc.url)

		require.Error(t, err)
	}
}

func TestGetCNIVersionSuccess(t *testing.T) {
	testcases := []struct {
		name        string
		version     string
		kubeVersion string
		expected    string
	}{
		{
			name:        "CNI version supplied, Kubernetes version < 1.17",
			version:     "0.8.3",
			kubeVersion: "1.16.0",
			expected:    pre117CNIVersion,
		},
		{
			name:        "CNI version supplied, Kubernetes version >= 1.17",
			version:     "0.8.3",
			kubeVersion: "1.17.0",
			expected:    "0.8.3",
		},
		{
			name:        "CNI version not supplied",
			kubeVersion: "1.17.0",
			expected:    minimumCNIVersion,
		},
	}

	for _, tc := range testcases {
		actual, err := getCNIVersion(
			&PackageDefinition{
				Version:           tc.version,
				KubernetesVersion: tc.kubeVersion,
			},
		)

		if err != nil {
			t.Fatalf("did not expect an error: %v", err)
		}

		assert.Equal(t, tc.expected, actual)
	}
}

func TestGetCNIVersionFailure(t *testing.T) {
	a := assert.New(t)

	_, err := getCNIVersion(nil)
	a.Error(err)
}

func TestGetCRIToolsVersionSuccess(t *testing.T) {
	testcases := []struct {
		name        string
		version     string
		kubeVersion string
		expected    string
	}{
		{
			name:     "User-supplied CRI tools version",
			version:  "1.17.0",
			expected: "1.17.0",
		},
		{
			name:        "Pre-release or CI Kubernetes version",
			kubeVersion: "1.18.0-alpha.1",
			expected:    "1.17.0",
		},
	}

	for _, tc := range testcases {
		actual, err := getCRIToolsVersion(
			&PackageDefinition{
				Version:           tc.version,
				KubernetesVersion: tc.kubeVersion,
			},
		)

		if err != nil {
			t.Fatalf("did not expect an error: %v", err)
		}

		assert.Equal(t, tc.expected, actual)
	}
}

func TestGetCRIToolsVersionFailure(t *testing.T) {
	a := assert.New(t)

	_, err := getCRIToolsVersion(nil)
	a.Error(err)
}

func TestGetDownloadLinkBaseSuccess(t *testing.T) {
	testcases := []struct {
		name        string
		kubeVersion string
		channel     ChannelType
		expected    string
	}{
		{
			name:        "CI version",
			kubeVersion: "1.18.0-alpha.1.277+2099c00290d262",
			channel:     ChannelNightly,
			expected:    "https://dl.k8s.io/ci/v1.18.0-alpha.1.277+2099c00290d262",
		},
		{
			name:        "non-CI version",
			kubeVersion: "1.18.0-alpha.1",
			expected:    "https://dl.k8s.io/v1.18.0-alpha.1",
		},
	}

	for _, tc := range testcases {
		actual, err := getDownloadLinkBase(
			&PackageDefinition{
				KubernetesVersion: tc.kubeVersion,
				Channel:           tc.channel,
			},
		)

		if err != nil {
			t.Fatalf("did not expect an error: %v", err)
		}

		assert.Equal(t, tc.expected, actual)
	}
}

func TestGetDownloadLinkBaseFailure(t *testing.T) {
	a := assert.New(t)

	_, err := getDownloadLinkBase(nil)
	a.Error(err)
}

func TestGetCIBuildsDownloadLinkBaseSuccess(t *testing.T) {
	testcases := []struct {
		name        string
		kubeVersion string
		expected    string
	}{
		{
			name:        "CI version",
			kubeVersion: "1.18.0-alpha.1.277+2099c00290d262",
			expected:    "https://dl.k8s.io/ci/v1.18.0-alpha.1.277+2099c00290d262",
		},
	}

	for _, tc := range testcases {
		actual, err := getCIBuildsDownloadLinkBase(
			&PackageDefinition{
				KubernetesVersion: tc.kubeVersion,
			},
		)

		if err != nil {
			t.Fatalf("did not expect an error: %v", err)
		}

		assert.Equal(t, tc.expected, actual)
	}
}

func TestGetCIBuildsDownloadLinkBaseFailure(t *testing.T) {
	a := assert.New(t)

	_, err := getCIBuildsDownloadLinkBase(nil)
	a.Error(err)
}

func TestGetDefaultReleaseDownloadLinkBaseSuccess(t *testing.T) {
	testcases := []struct {
		name        string
		kubeVersion string
		expected    string
	}{
		{
			name:        "Release version",
			kubeVersion: "1.17.0",
			expected:    "https://dl.k8s.io/v1.17.0",
		},
		{
			name:        "Pre-release version",
			kubeVersion: "1.18.0-alpha.1",
			expected:    "https://dl.k8s.io/v1.18.0-alpha.1",
		},
	}

	for _, tc := range testcases {
		actual, err := getDefaultReleaseDownloadLinkBase(
			&PackageDefinition{
				KubernetesVersion: tc.kubeVersion,
			},
		)

		if err != nil {
			t.Fatalf("did not expect an error: %v", err)
		}

		assert.Equal(t, tc.expected, actual)
	}
}

func TestGetDefaultReleaseDownloadLinkBaseFailure(t *testing.T) {
	a := assert.New(t)

	_, err := getDefaultReleaseDownloadLinkBase(nil)
	a.Error(err)
}

func TestGetKubeadmDependenciesSuccess(t *testing.T) {
	testcases := []struct {
		name    string
		version string
		deps    []string
	}{
		{
			name:    "minimum supported kubernetes",
			version: "1.15.0",
			deps: []string{
				"kubelet (>= 1.13.0)",
				"kubectl (>= 1.13.0)",
				"kubernetes-cni (>= 0.7.5)",
				"cri-tools (>= 1.13.0)",
				"${misc:Depends}",
			},
		},
		{
			name:    "latest stable minor kubernetes",
			version: "1.17.0",
			deps: []string{
				"kubelet (>= 1.13.0)",
				"kubectl (>= 1.13.0)",
				"kubernetes-cni (>= 0.7.5)",
				"cri-tools (>= 1.13.0)",
				"${misc:Depends}",
			},
		},
		{
			name:    "latest alpha kubernetes",
			version: "1.18.0-alpha.1",
			deps: []string{
				"kubelet (>= 1.13.0)",
				"kubectl (>= 1.13.0)",
				"kubernetes-cni (>= 0.7.5)",
				"cri-tools (>= 1.13.0)",
				"${misc:Depends}",
			},
		},
		{
			name:    "next stable minor kubernetes",
			version: "1.18.0",
			deps: []string{
				"kubelet (>= 1.13.0)",
				"kubectl (>= 1.13.0)",
				"kubernetes-cni (>= 0.7.5)",
				"cri-tools (>= 1.13.0)",
				"${misc:Depends}",
			},
		},
	}

	for _, tc := range testcases {
		expected := strings.Join(tc.deps, ", ")

		actual, err := GetKubeadmDependencies(
			&PackageDefinition{
				Version: tc.version,
			},
		)

		if err != nil {
			t.Fatalf("did not expect an error: %v", err)
		}

		assert.Equal(t, expected, actual)
	}
}

func TestGetKubeadmDependenciesFailure(t *testing.T) {
	a := assert.New(t)

	_, err := GetKubeadmDependencies(nil)
	a.Error(err)
}

func TestGetCNIDownloadLinkSuccess(t *testing.T) {
	testcases := []struct {
		name     string
		version  string
		arch     string
		expected string
	}{
		{
			name:     "CNI <= 0.7.5",
			version:  "0.7.5",
			arch:     "amd64",
			expected: "https://github.com/containernetworking/plugins/releases/download/v0.7.5/cni-plugins-amd64-v0.7.5.tgz",
		},
		{
			name:     "CNI > 0.8.3",
			version:  "0.8.3",
			arch:     "amd64",
			expected: "https://github.com/containernetworking/plugins/releases/download/v0.8.3/cni-plugins-linux-amd64-v0.8.3.tgz",
		},
	}

	for _, tc := range testcases {
		actual, err := getCNIDownloadLink(
			&PackageDefinition{
				Version: tc.version,
			},
			tc.arch,
		)

		if err != nil {
			t.Fatalf("did not expect an error: %v", err)
		}

		assert.Equal(t, tc.expected, actual)
	}
}

func TestGetCNIDownloadLinkFailure(t *testing.T) {
	a := assert.New(t)

	_, err := getCNIDownloadLink(nil, "amd64")
	a.Error(err)
}
