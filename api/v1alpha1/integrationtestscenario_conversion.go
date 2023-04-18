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
// ConvertTo converts this ITS to the Hub version (v1beta1).
func (src *IntegrationTestScenario) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.IntegrationTestScenario)
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec.Application = src.Spec.Application
	if &src.Spec.Environment != nil {
		dst.Spec.Environment = v1beta1.TestEnvironment(src.Spec.Environment)
	}
	if src.Spec.Params != nil {
		for _, par := range src.Spec.Params {
			dst.Spec.Params = append(dst.Spec.Params, v1beta1.PipelineParameter(par))
		}
	}
	if src.Spec.Contexts != nil {
		for _, par := range src.Spec.Contexts {
			dst.Spec.Contexts = append(dst.Spec.Contexts, v1beta1.TestContext(par))
		}
	}
	dst.Spec.ResolverRef = v1beta1.ResolverRef{
		Resolver: "bundle",
		Params: []v1beta1.ResolverParameter{
			{
				Name:  "bundle",
				Value: src.Spec.Bundle,
			},
			{
				Name:  "name",
				Value: src.Spec.Pipeline,
			},
			{
				Name:  "kind",
				Value: "pipeline",
			},
		},
	}
	return nil
}

func (dst *IntegrationTestScenario) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.IntegrationTestScenario)
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec.Application = src.Spec.Application
	if &dst.Spec.Environment != nil {
		src.Spec.Environment = v1beta1.TestEnvironment(dst.Spec.Environment)
	}
	if src.Spec.Params != nil {
		for _, par := range src.Spec.Params {
			src.Spec.Params = append(src.Spec.Params, par)
		}
	}
	if src.Spec.Contexts != nil {
		for _, par := range src.Spec.Contexts {
			src.Spec.Contexts = append(src.Spec.Contexts, par)
		}
	}
	src.Spec.ResolverRef = v1beta1.ResolverRef{
		Resolver: "bundle",
		Params: []v1beta1.ResolverParameter{
			{
				Name:  "bundle",
				Value: dst.Spec.Bundle,
			},
			{
				Name:  "name",
				Value: dst.Spec.Pipeline,
			},
			{
				Name:  "kind",
				Value: "pipeline",
			},
		},
	}
	return nil
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
