package gitops

import (
	"github.com/google/uuid"
	applicationapiv1alpha1 "github.com/redhat-appstudio/application-api/api/v1alpha1"
	"github.com/redhat-appstudio/integration-service/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CopiedEnvironment struct {
	applicationapiv1alpha1.Environment
}

// AsPipelineRun casts the IntegrationPipelineRun to PipelineRun, so it can be used in the Kubernetes client.
func (r *CopiedEnvironment) AsEnvironment() *applicationapiv1alpha1.Environment {
	return &r.Environment
}

// check if Env already contains requested Name
func contains(s applicationapiv1alpha1.EnvironmentConfiguration, e string) bool {
	for a := range s.Env {
		if s.Env[a].Name == e {
			return true
		}
	}
	return false
}

// NewCopyOfExistingEnvironment gets the existing environment from current namespace and makes copy of it
// new name is generated consisting of existing environment name and integrationTestScenario name
// targetNamespace gets name of integrationTestScenario and uuid
func NewCopyOfExistingEnvironment(existingEnvironment *applicationapiv1alpha1.Environment, namespace string, integrationTestScenario *v1alpha1.IntegrationTestScenario, applicationSnapshot *applicationapiv1alpha1.ApplicationSnapshot) *CopiedEnvironment {
	id := uuid.New()
	existingApiURL := existingEnvironment.Spec.UnstableConfigurationFields.KubernetesClusterCredentials.APIURL
	existingClusterCreds := existingEnvironment.Spec.UnstableConfigurationFields.KubernetesClusterCredentials.ClusterCredentialsSecret

	copyEnvVar := applicationapiv1alpha1.EnvironmentConfiguration{}
	copyEnvVar = existingEnvironment.Spec.Configuration

	for intEnvVars := range integrationTestScenario.Spec.Environment.Configuration.Env {
		// if existing environment does not contain EnvVars, copy ones from IntegrationTestScenario
		if existingEnvironment.Spec.Configuration.Env == nil {
			copyEnvVar.Env = integrationTestScenario.Spec.Environment.Configuration.Env
			break
		}
		for existingEnvVar := range existingEnvironment.Spec.Configuration.Env {
			// envVar names are matching? overwrite existing environment with one from ITS
			if integrationTestScenario.Spec.Environment.Configuration.Env[intEnvVars].Name == copyEnvVar.Env[existingEnvVar].Name {
				copyEnvVar.Env[existingEnvVar].Value = integrationTestScenario.Spec.Environment.Configuration.Env[intEnvVars].Value
			} else if (integrationTestScenario.Spec.Environment.Configuration.Env[intEnvVars].Name != copyEnvVar.Env[existingEnvVar].Name) && (!contains(copyEnvVar, integrationTestScenario.Spec.Environment.Configuration.Env[intEnvVars].Name)) {
				// in case that EnvVar from IntegrationTestScenario is not matching any EnvVar from existingEnv, add this ITS EnvVar to coppied Environment
				copyEnvVar.Env = append(copyEnvVar.Env, applicationapiv1alpha1.EnvVarPair{Name: integrationTestScenario.Spec.Environment.Configuration.Env[intEnvVars].Name, Value: integrationTestScenario.Spec.Environment.Configuration.Env[intEnvVars].Value})
			}
		}
	}

	copyOfEnvironment := applicationapiv1alpha1.Environment{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: existingEnvironment.Name + "-" + integrationTestScenario.Name + "-",
			Namespace:    namespace,
		},
		Spec: applicationapiv1alpha1.EnvironmentSpec{
			Type:               applicationapiv1alpha1.EnvironmentType_POC,
			DisplayName:        existingEnvironment.Name + "-" + integrationTestScenario.Name,
			Tags:               []string{"ephemeral"},
			DeploymentStrategy: applicationapiv1alpha1.DeploymentStrategy_Manual,
			Configuration:      copyEnvVar,
			UnstableConfigurationFields: &applicationapiv1alpha1.UnstableEnvironmentConfiguration{
				KubernetesClusterCredentials: applicationapiv1alpha1.KubernetesClusterCredentials{
					TargetNamespace:          integrationTestScenario.Name + "-" + id.String(),
					APIURL:                   existingApiURL,
					ClusterCredentialsSecret: existingClusterCreds,
				},
			},
		},
	}
	return &CopiedEnvironment{copyOfEnvironment}
}

// WithIntegrationLabels adds IntegrationTestScenario name as label to the coppied environment.
func (e *CopiedEnvironment) WithIntegrationLabels(integrationTestScenario *v1alpha1.IntegrationTestScenario) *CopiedEnvironment {
	if e.ObjectMeta.Labels == nil {
		e.ObjectMeta.Labels = map[string]string{}
	}
	e.ObjectMeta.Labels["test.appstudio.openshift.io/scenario"] = integrationTestScenario.Name

	return e

}

//WithApplicationSnapshot adds the name of snapshot as label to the coppied environment.
func (e *CopiedEnvironment) WithApplicationSnapshot(snapshot *applicationapiv1alpha1.ApplicationSnapshot) *CopiedEnvironment {

	if e.ObjectMeta.Labels == nil {
		e.ObjectMeta.Labels = map[string]string{}
	}
	e.ObjectMeta.Labels["test.appstudio.openshift.io/snapshot"] = snapshot.Name

	return e
}
