/*
Copyright 2022.

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

package v1alpha1

import (
	"github.com/redhat-appstudio/integration-service/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// log is for logging in this package.
var integrationtestscenariolog = logf.Log.WithName("integrationtestscenario-resource")

func (r *IntegrationTestScenario) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// Hub marks this type as a conversion hub.
// ConvertTo converts this Memcached to the Hub version (v2alpha1).
func (src *IntegrationTestScenario) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.IntegrationTestScenario)
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec.Application = src.Spec.Application
	//dst.Spec.Environment = v1beta1.TestEnvironment(src.Spec.Environment)
	dst.Spec.ResolverRef = v1beta1.ResolverRef{
		Resolver: "bundle",
		Params: []v1beta1.ResolverParameter{
			{
				Name:  "bundle",
				Value: src.Spec.Bundle,
				// tektonv1beta1.ParamValue{
				// 	Type:      tektonv1beta1.ParamTypeString,
				// 	StringVal: src.Spec.Bundle,
				// },
			},
			{
				Name:  "name",
				Value: src.Spec.Pipeline,
				// tektonv1beta1.ParamValue{
				// 	Type:      tektonv1beta1.ParamTypeString,
				// 	StringVal: src.Spec.Pipeline,
				// },
			},
			{
				Name:  "kind",
				Value: "pipeline",
				// Value: tektonv1beta1.ParamValue{
				// 	Type:      tektonv1beta1.ParamTypeString,
				// 	StringVal: "pipeline",
				// },
			},
		},
	}
	return nil
}

func (dst *IntegrationTestScenario) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.IntegrationTestScenario)
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec.Application = src.Spec.Application
	//dst.Spec.Environment = TestEnvironment(src.Spec.Environment)
	src.Spec.ResolverRef = v1beta1.ResolverRef{
		Resolver: "bundle",
		Params: []v1beta1.ResolverParameter{
			{
				Name:  "bundle",
				Value: dst.Spec.Bundle,
				// tektonv1beta1.ParamValue{
				// 	Type:      tektonv1beta1.ParamTypeString,
				// 	StringVal: dst.Spec.Bundle,
				// },
			},
			{
				Name:  "name",
				Value: dst.Spec.Pipeline,
				// tektonv1beta1.ParamValue{
				// 	Type:      tektonv1beta1.ParamTypeString,
				// 	StringVal: dst.Spec.Pipeline,
				// },
			},
			{
				Name:  "kind",
				Value: "pipeline",
				// tektonv1beta1.ParamValue{
				// 	Type:      tektonv1beta1.ParamTypeString,
				// 	StringVal: "pipeline",
				// },
			},
		},
	}
	return nil
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
