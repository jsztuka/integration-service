package gitops

import (
	"github.com/google/uuid"
	applicationapiv1alpha1 "github.com/redhat-appstudio/application-api/api/v1alpha1"
	"github.com/redhat-appstudio/integration-service/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterCredentials struct {
	applicationapiv1alpha1.KubernetesClusterCredentials
}

type CopiedEnvironment struct {
	applicationapiv1alpha1.Environment
}

// AsPipelineRun casts the IntegrationPipelineRun to PipelineRun, so it can be used in the Kubernetes client.
func (r *CopiedEnvironment) AsEnvironment() *applicationapiv1alpha1.Environment {
	return &r.Environment
}

func NewCopyOfExistingEnvironment(existingEnvironment *applicationapiv1alpha1.Environment, namespace string, integrationTestScenario *v1alpha1.IntegrationTestScenario) *CopiedEnvironment {
	id := uuid.New()
	//existingTargetNamespace := existingEnvironment.Spec.UnstableConfigurationFields.TargetNamespace
	existingApiURL := existingEnvironment.Spec.UnstableConfigurationFields.APIURL
	existingClusterCreds := existingEnvironment.Spec.UnstableConfigurationFields.ClusterCredentialsSecret

	//here should be the magic that decides envVars

	copyEnvVar := applicationapiv1alpha1.EnvironmentConfiguration{}
	// if integrationtestScenario has no EnvVars specified
	// use existing EnvVars
	if integrationTestScenario.Spec.Environment.Configuration.Env == nil {
		copyEnvVar.Env = existingEnvironment.Spec.Configuration.Env
	} else {
		for intEnvVars := range integrationTestScenario.Spec.Environment.Configuration.Env {
			match := false
			// if existing environment does not contain EnvVars, copy ones from IntegrationTestScenario
			if existingEnvironment.Spec.Configuration.Env == nil {
				copyEnvVar.Env = integrationTestScenario.Spec.Environment.Configuration.Env
				break
			}
			for existingEnvVar := range existingEnvironment.Spec.Configuration.Env {
				if integrationTestScenario.Spec.Environment.Configuration.Env[intEnvVars].Name == existingEnvironment.Spec.Configuration.Env[existingEnvVar].Name {
					match = true
					copyEnvVar.Env[existingEnvVar].Value = integrationTestScenario.Spec.Environment.Configuration.Env[intEnvVars].Value
				}
				// in case that EnvVar from IntegrationTestScenario is not matching any EnvVar from existingEnv, add this ITS EnvVar to coppied Environment
				if !match && (existingEnvVar == len(existingEnvironment.Spec.Configuration.Env)-1) {
					copyEnvVar.Env = append(copyEnvVar.Env, applicationapiv1alpha1.EnvVarPair{Name: integrationTestScenario.Spec.Environment.Configuration.Env[intEnvVars].Name, Value: integrationTestScenario.Spec.Environment.Configuration.Env[intEnvVars].Value})
				}
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
			UnstableConfigurationFields: &applicationapiv1alpha1.UnstableEnvironmentConfiguration{ //TO-DO change naming
				KubernetesClusterCredentials: applicationapiv1alpha1.KubernetesClusterCredentials{
					TargetNamespace:          integrationTestScenario.Name + "-" + id.String(), //Copied environment needs to hold name of (integrationScenario.Name + snapshot.Name + -UUID)
					APIURL:                   existingApiURL,
					ClusterCredentialsSecret: existingClusterCreds,
				},
			},
		},
	}
	return &CopiedEnvironment{copyOfEnvironment}
}

func (e *CopiedEnvironment) WithIntegrationLabels(integrationTestScenario *v1alpha1.IntegrationTestScenario) *CopiedEnvironment {
	if e.ObjectMeta.Labels == nil {
		e.ObjectMeta.Labels = map[string]string{}
	}
	e.ObjectMeta.Labels["test.appstudio.openshift.io/scenario"] = integrationTestScenario.Name

	return e

}

func (e *CopiedEnvironment) WithApplicationSnapshot(snapshot *applicationapiv1alpha1.ApplicationSnapshot) *CopiedEnvironment {

	if e.ObjectMeta.Labels == nil {
		e.ObjectMeta.Labels = map[string]string{}
	}
	e.ObjectMeta.Labels["test.appstudio.openshift.io/snapshot"] = snapshot.Name

	return e
}

//add lables from environment to binding (environment and binding should have same labels)
