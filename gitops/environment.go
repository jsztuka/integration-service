package gitops

import (
	applicationapiv1alpha1 "github.com/redhat-appstudio/application-api/api/v1alpha1"
	"github.com/redhat-appstudio/integration-service/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterCredentials struct {
	applicationapiv1alpha1.KubernetesClusterCredentials
}

//newly created environment needs to point to Snapshot nad IntegrationTestScenario, hence it needs to be set in the same environment
func (t *ClusterCredentials) SetTargetNamespace(namespace string) *ClusterCredentials {
	t.TargetNamespace = namespace
	return t
}

type CopiedEnvironment struct {
	applicationapiv1alpha1.Environment
}

func NewCopyOfExistingEnvironment(existingEnvironment applicationapiv1alpha1.Environment, namespace string, integrationTestScenario *v1alpha1.IntegrationTestScenario) *CopiedEnvironment {

	existingTargetNamespace := existingEnvironment.Spec.UnstableConfigurationFields.TargetNamespace
	existingApiURL := existingEnvironment.Spec.UnstableConfigurationFields.APIURL
	existingClusterCreds := existingEnvironment.Spec.UnstableConfigurationFields.ClusterCredentialsSecret

	copyOfEnvironment := applicationapiv1alpha1.Environment{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: existingEnvironment.Name + "-",
			Namespace:    namespace,
		},
		Spec: applicationapiv1alpha1.EnvironmentSpec{
			Type:               applicationapiv1alpha1.EnvironmentType_POC,
			DisplayName:        existingEnvironment.Name,
			Tags:               []string{"ephemeral"},
			DeploymentStrategy: applicationapiv1alpha1.DeploymentStrategy_Manual,
			UnstableConfigurationFields: &applicationapiv1alpha1.UnstableEnvironmentConfiguration{ //TO-DO change naming
				KubernetesClusterCredentials: applicationapiv1alpha1.KubernetesClusterCredentials{
					TargetNamespace:          existingTargetNamespace,
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

func (e *CopiedEnvironment) WithApplicationSnapshot(snapshot *applicationapiv1alpha1.Snapshot) *CopiedEnvironment {

	if e.ObjectMeta.Labels == nil {
		e.ObjectMeta.Labels = map[string]string{}
	}
	e.ObjectMeta.Labels["test.appstudio.openshift.io/snapshot"] = snapshot.Name

	return e
}
