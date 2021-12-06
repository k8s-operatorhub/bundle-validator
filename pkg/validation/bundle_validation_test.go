// Copyright 2021 The K8s Community Validator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validation

import (
	"testing"

	"github.com/operator-framework/api/pkg/manifests"
	"github.com/stretchr/testify/require"

)

func Test_Test_checkMaxKubeVersionAnnotation(t *testing.T) {
	type args struct {
		annotations   map[string]string
		bundleDir     string
	}
	tests := []struct {
		name        string
		args        args
		wantError   bool
		wantWarning bool
		errStrings  []string
		warnStrings []string
	}{
		{
			name:      "should work successfully when has no deprecated apis",
			wantError: false,
			args: args{
				bundleDir: "./testdata/bundle_v1",
			},
		},
		{
			name: "should pass when has deprecated apis and kubeMaxVersion annotation is set",
			wantError:   false,
			args: args{
				bundleDir: "./testdata/bundle_v1beta1",
				annotations: map[string]string{
					KubeMaxAnnotation: K8sVerV1betavSupported,
				},
			},
		},
		{
			name: "should fail when has deprecated apis and kubeMaxVersion annotation is not set",
			wantError:   true,
			errStrings: []string{
				"Error: Value : (memcached-operator.v0.0.1) operators.operatorframework.io/maxKubeVersion metadata.annotation is not infomed. This distributions still using the removed APIs then, you **MUST** ensure that its CSV has the informative metadata annotation `operators.operatorframework.io/maxKubeVersion`. More info: this bundle is using APIs which were deprecated and removed in v1.22. More info: https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-22. Migrate the API(s) for CRD: ([\"memcacheds.cache.example.com\"])",
			},
			args: args{
				bundleDir: "./testdata/bundle_v1beta1",
			},
		},
		{
			name: "should fail when has deprecated apis and kubeMaxVersion set is > 1.21",
			wantError:   true,
			errStrings: []string{
				"Error: Value : (memcached-operator.v0.0.1) invalid value for operators.operatorframework.io/maxKubeVersion. The K8s version value 1.22.0 is >= of 1.22.0. Note that this bundle is using APIs which were deprecated and removed in v1.22. More info: https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-22. Migrate the API(s) for CRD: ([\"memcacheds.cache.example.com\"])",
			},
			args: args{
				bundleDir: "./testdata/bundle_v1beta1",
				annotations: map[string]string{
					KubeMaxAnnotation: K8sVerV1betav1Unsupported,
				},
			},
		},
		{
			name: "should fail when kubeMaxVersion is invalid",
			wantError:   true,
			errStrings: []string{
				"Error: Value : (memcached-operator.v0.0.1) operators.operatorframework.io/maxKubeVersion metadata.annotation value (invalid) is invalid. Error: Invalid character(s) found in major number \"invalid\" ",
			},
			args: args{
				bundleDir: "./testdata/bundle_v1",
				annotations: map[string]string{
					KubeMaxAnnotation: "invalid",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Validate the bundle object
			bundle, err := manifests.GetBundleFromDir(tt.args.bundleDir)
			require.NoError(t, err)

			if len(tt.args.annotations) > 0 {
				bundle.CSV.Annotations = tt.args.annotations
			}

			results := validateK8sCommunityBundle(bundle)
			require.Equal(t, tt.wantWarning, len(results.Warnings) > 0)
			if tt.wantWarning {
				require.Equal(t, len(tt.warnStrings), len(results.Warnings))
				for _, w := range results.Warnings {
					wString := w.Error()
					require.Contains(t, tt.warnStrings, wString)
				}
			}

			require.Equal(t, tt.wantError, len(results.Errors) > 0)
			if tt.wantError {
				require.Equal(t, len(tt.errStrings), len(results.Errors))
				for _, err := range results.Errors {
					errString := err.Error()
					require.Contains(t, tt.errStrings, errString)
				}
			}
		})
	}
}
