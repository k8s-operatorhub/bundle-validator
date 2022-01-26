// Copyright 2021 K8S Community Bundle Validator Authors.
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
	"fmt"

	"github.com/blang/semver"
	"github.com/operator-framework/api/pkg/manifests"
	"github.com/operator-framework/api/pkg/validation"
	"github.com/operator-framework/api/pkg/validation/errors"
	interfaces "github.com/operator-framework/api/pkg/validation/interfaces"
)

// K8S Community Bundle Validator validate the bundle manifests against the required and specific criteria to publish
// the projects on the K8S community catalog
//
// Be aware that this validator is in alpha stage and can be changed. Also, the intention here is to decouple
// this validator and move it out of this project. Following its current checks:
//
// - Ensure that when found the usage of the removed APIs on 1.22/OCP 4.9 the CSV has the annotation
// operatorhub.io/ui-metadata-max-k8s-version with a value < 1.22.
var K8sCommunityBundleValidator interfaces.Validator = interfaces.ValidatorFunc(k8sCommunityBundleValidator)

// KubeMaxAnnotation define the annotation that will be checked
const KubeMaxAnnotation = "operatorhub.io/ui-metadata-max-k8s-version"

// K8sVerV1betav1Unsupported version where the apis v1betav1 is no longer supported
const K8sVerV1betav1Unsupported = "1.22.0"

// K8sVerV1betavSupported version where the apis v1betav1 is supported
const K8sVerV1betavSupported = "1.21.0"

func k8sCommunityBundleValidator(objs ...interface{}) (results []errors.ManifestResult) {
	for _, obj := range objs {
		switch v := obj.(type) {
		case *manifests.Bundle:
			results = append(results, validateK8sCommunityBundle(v))
		}
	}

	return results
}

// K8sCommunityBundleValidatorChecks defines the attributes used to perform the checks
type K8sCommunityBundleValidatorChecks struct {
	bundle           manifests.Bundle
	deprecateAPIsMsg string
	errs             []error
	warns            []error
}

// validateK8sCommunityBundle will check the bundle against the criteria to publish into K8s Community Catalog
func validateK8sCommunityBundle(bundle *manifests.Bundle) errors.ManifestResult {
	result := errors.ManifestResult{}
	if bundle == nil {
		result.Add(errors.ErrInvalidBundle("Bundle is nil", nil))
		return result
	}
	result.Name = bundle.Name

	if bundle.CSV == nil {
		result.Add(errors.ErrInvalidBundle("Bundle csv is nil", bundle.Name))
		return result
	}

	checks := K8sCommunityBundleValidatorChecks{bundle: *bundle, errs: []error{}, warns: []error{}}

	objs := bundle.ObjectsToValidate()
	for _, obj := range bundle.Objects {
		objs = append(objs, obj)
	}

	// pass the objects to the validator
	resultDeprecation := validation.AlphaDeprecatedAPIsValidator.Validate(objs...)

	for _, res := range resultDeprecation {
		for _, res := range res.Warnings {
			checks.deprecateAPIsMsg = res.Detail
		}
	}

	checks = checkMaxKubeVersionAnnotation(checks)

	for _, err := range checks.errs {
		result.Add(errors.ErrInvalidCSV(err.Error(), bundle.CSV.GetName()))
	}
	for _, warn := range checks.warns {
		result.Add(errors.WarnInvalidCSV(warn.Error(), bundle.CSV.GetName()))
	}

	return result
}

// checkMaxVersionAnnotation will verify if the annotation was informed and has a valid value according
// the findings found using the AlphaDeprecatedAPI validator implemented in the operator-framework/api
func checkMaxKubeVersionAnnotation(checks K8sCommunityBundleValidatorChecks) K8sCommunityBundleValidatorChecks {
	kubeMax := checks.bundle.CSV.Annotations[KubeMaxAnnotation]
	semVerKube1_22, _ := semver.ParseTolerant(K8sVerV1betav1Unsupported)

	semVerKubeMax, err := semver.ParseTolerant(kubeMax)
	if len(kubeMax) > 0 {
		if err != nil {
			checks.errs = append(checks.errs, fmt.Errorf("%s metadata.annotation value (%s) is invalid. Error: %s ",
				KubeMaxAnnotation,
				kubeMax,
				err))
			return checks
		}
	}

	if len(checks.deprecateAPIsMsg) > 0 {
		if len(kubeMax) == 0 {
			if err != nil {
				checks.errs = append(checks.errs, fmt.Errorf("%s metadata.annotation is not infomed. "+
					"This distributions still using the removed APIs then, you **MUST** ensure that its " +
					"CSV has the informative metadata annotation `%s`. More info: %s",
					KubeMaxAnnotation,
					KubeMaxAnnotation,
					checks.deprecateAPIsMsg))
				return checks
			}
		}

		if semVerKubeMax.GE(semVerKube1_22) {
			checks.errs = append(checks.errs, fmt.Errorf("invalid value for %s. "+
				"The K8s version value %s is >= of %s. Note that %s",
				KubeMaxAnnotation,
				kubeMax,
				K8sVerV1betav1Unsupported,
				checks.deprecateAPIsMsg))
			return checks
		}

		return checks
	}

	return checks
}
