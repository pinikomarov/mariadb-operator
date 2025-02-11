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

package v1beta1

import (
	condition "github.com/openstack-k8s-operators/lib-common/modules/common/condition"
	"github.com/openstack-k8s-operators/lib-common/modules/common/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Container image fall-back defaults

	// MariaDBContainerImage is the fall-back container image for MariaDB/Galera
	MariaDBContainerImage = "quay.io/podified-antelope-centos9/openstack-mariadb:current-podified"
)

// MariaDBSpec defines the desired state of MariaDB
type MariaDBSpec struct {
	// Secret containing a RootPassword
	// +kubebuilder:validation:Required
	Secret string `json:"secret"`
	// Storage class to host the mariadb databases
	// +kubebuilder:validation:Required
	StorageClass string `json:"storageClass"`
	// Storage size allocated for the mariadb databases
	// +kubebuilder:validation:Required
	StorageRequest string `json:"storageRequest"`
	// ContainerImage - Container Image URL (will be set to environmental default if empty)
	// +kubebuilder:validation:Required
	ContainerImage string `json:"containerImage"`
	// Adoption configuration
	// +kubebuilder:validation:Optional
	AdoptionRedirect AdoptionRedirectSpec `json:"adoptionRedirect,omitempty"`
}

// AdoptionRedirectSpec defines redirection to a different DB instance during Adoption
type AdoptionRedirectSpec struct {
	// MariaDB host to redirect to (IP or name)
	Host string `json:"host,omitempty"`
}

// MariaDBStatus defines the observed state of MariaDB
type MariaDBStatus struct {
	// db init completed
	DbInitHash string `json:"dbInitHash"`

	// Conditions
	Conditions condition.Conditions `json:"conditions,omitempty" optional:"true"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[0].status",description="Status"
//+kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[0].message",description="Message"

// MariaDB is the Schema for the mariadbs API
type MariaDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MariaDBSpec   `json:"spec,omitempty"`
	Status MariaDBStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MariaDBList contains a list of MariaDB
type MariaDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MariaDB `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MariaDB{}, &MariaDBList{})
}

// IsReady - returns true if service is ready to serve requests
func (instance MariaDB) IsReady() bool {
	return instance.Status.Conditions.IsTrue(condition.DeploymentReadyCondition)
}

// RbacConditionsSet - set the conditions for the rbac object
func (instance MariaDB) RbacConditionsSet(c *condition.Condition) {
	instance.Status.Conditions.Set(c)
}

// RbacNamespace - return the namespace
func (instance MariaDB) RbacNamespace() string {
	return instance.Namespace
}

// RbacResourceName - return the name to be used for rbac objects (serviceaccount, role, rolebinding)
func (instance MariaDB) RbacResourceName() string {
	return "mariadb-" + instance.Name
}

// SetupDefaults - initializes any CRD field defaults based on environment variables (the defaulting mechanism itself is implemented via webhooks)
func SetupDefaults() {
	// Acquire environmental defaults and initialize Keystone defaults with them
	mariaDBDefaults := MariaDBDefaults{
		ContainerImageURL: util.GetEnvVar("RELATED_IMAGE_MARIADB_IMAGE_URL_DEFAULT", MariaDBContainerImage),
	}

	SetupMariaDBDefaults(mariaDBDefaults)
}
