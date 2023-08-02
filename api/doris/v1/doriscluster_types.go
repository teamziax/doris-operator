/*
Copyright 2023.

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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DorisClusterSpec defines the desired state of DorisCluster
type DorisClusterSpec struct {
	//defines the pod that will created from feSpec template.
	FeSpec *FeSpec `json:"feSpec,omitempty"`

	//defines the pod that will created from beSpec template.
	BeSpec *BeSpec `json:"beSpec,omitempty"`

	//defines the pod that will created from cnSpec template.
	CnSpec *CnSpec `json:"cnSpec,omitempty"`

	BrokerSpec *BrokerSpec `json:"brokerSpec,omitempty"`

	//components register or drop self in doris cluster.
	AdminUser *AdminUser `json:"adminUser,omitempty"`
}

// AdminUser manage the software service nodes in doris cluster.
type AdminUser struct {
	//the user name for admin service's node.
	Name string `json:"name,omitempty"`

	//login to doris db.
	Password string `json:"password,omitempty"`
}

// describes a template for creating copies of a fe software service.
type FeSpec struct {
	//the number of fe in election. electionNumber <= replicas, left as observers. default value=3
	ElectionNumber *int32 `json:"electionNumber,omitempty"`

	//the foundation spec for creating be software services.
	BaseSpec `json:",inline"`
}

// describes a template for creating copies of a be software service.
type BeSpec struct {

	//the foundation spec for creating be software services.
	BaseSpec `json:",inline"`
}

// Fe address for other components access, if not config generate default.
type FeAddress struct {
	//the service name that proxy fe on k8s. the service must in same namespace with fe.
	ServiceName string `json:"ServiceName,omitempty"`

	//the fe addresses if not deploy by crd, user can use k8s deploy fe observer.
	Endpoints Endpoints `json:"endpoints,omitempty"`
}

type Endpoints struct {
	Address []string `json:":address,omitempty"`
	Port    int      `json:"port,omitempty"`
}

// describes a template for creating copies of a cn software service. cn, the service for external table.
type CnSpec struct {
	//the foundation spec for creating cn software services.
	BaseSpec `json:",inline"`

	//AutoScalingPolicy auto scaling strategy
	AutoScalingPolicy *AutoScalingPolicy `json:"autoScalingPolicy,omitempty"`
}

type BrokerSpec struct {
	//expose the cn listen ports
	Service ExportService `json:"service,omitempty"`

	//the foundation spec for creating cn software services.
	BaseSpec `json:"baseSpec,omitempty"`
}

// describe the foundation spec of component about doris.
type BaseSpec struct {
	//annotation for fe pods. user can config monitor annotation for collect to monitor system.
	Annotations map[string]string `json:"annotations,omitempty"`

	//serviceAccount for cn access cloud service.
	ServiceAccount string `json:"serviceAccount,omitempty"`

	//expose the be listen ports
	Service *ExportService `json:"service,omitempty"`

	//A special supplemental group that applies to all containers in a pod.
	// Some volume types allow the Kubelet to change the ownership of that volume
	// to be owned by the pod:
	FsGroup *int64 `json:"fsGroup,omitempty"`
	// specify register fe addresses
	FeAddress *FeAddress `json:"feAddress,omitempty"`

	//Replicas is the number of desired cn Pod.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:feSpecMinimum=3
	//+optional
	Replicas *int32 `json:"replicas,omitempty"`

	//Image for a starrocks cn deployment.
	Image string `json:"image"`

	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.
	// If specified, these secrets will be passed to individual puller implementations for them to use.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,15,rep,name=imagePullSecrets"`

	//+optional
	//set the fe service for register cn, when not set, will use the fe config to find.
	//Deprecated,
	//FeServiceName string `json:"feServiceName,omitempty"`

	//the reference for cn configMap.
	//+optional
	ConfigMapInfo ConfigMapInfo `json:"configMapInfo,omitempty"`

	//defines the specification of resource cpu and mem.
	corev1.ResourceRequirements `json:",inline"`
	// (Optional) If specified, the pod's nodeSelector，displayName="Map of nodeSelectors to match when scheduling pods on nodes"
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	//+optional
	//cnEnvVars is a slice of environment variables that are added to the pods, the default is empty.
	EnvVars []corev1.EnvVar `json:"envVars,omitempty"`

	//+optional
	//If specified, the pod's scheduling constraints.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// (Optional) Tolerations for scheduling pods onto some dedicated nodes
	//+optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	//+optional
	// podLabels for user selector or classify pods
	PodLabels map[string]string `json:"podLabels,omitempty"`

	// HostAliases is an optional list of hosts and IPs that will be injected into the pod's hosts
	// file if specified. This is only valid for non-hostNetwork pods.
	// +optional
	HostAliases []corev1.HostAlias `json:"hostAliases,omitempty"`

	PersistentVolumes []PersistentVolume `json:"persistentVolumes,omitempty"`
}

type PersistentVolume struct {
	// volumeClaimTemplates is a list of claims that pods are allowed to reference.
	corev1.PersistentVolumeClaim `json:"persistentVolumeClaim"`
	MountPath                    string `json:"mountPath"`
}

type ConfigMapInfo struct {
	//the config info for start progress.
	ConfigMapName string `json:"configMapName,omitempty"`

	//the config response key in configmap.
	ResolveKey string `json:"resolveKey,omitempty"`
}

// ExportService consisting of expose ports for user access to software service.
type ExportService struct {
	//type of service,the possible value for the service type are : ClusterIP, NodePort, LoadBalancer,ExternalName.
	//More info: https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types
	// +optional
	Type corev1.ServiceType `json:"type,omitempty"`

	// Only applies to Service Type: LoadBalancer.
	// This feature depends on whether the underlying cloud-provider supports specifying
	// the loadBalancerIP when a load balancer is created.
	// This field will be ignored if the cloud-provider does not support the feature.
	// This field was under-specified and its meaning varies across implementations,
	// and it cannot support dual-stack.
	// As of Kubernetes v1.24, users are encouraged to use implementation-specific annotations when available.
	// This field may be removed in a future API version.
	// +optional
	LoadBalancerIP string `json:"loadBalancerIP,omitempty"`
}

// DorisClusterStatus defines the observed state of DorisCluster
type DorisClusterStatus struct {
	FEStatus *ComponentStatus `json:"feStatus,omitempty"`

	BEStatus *ComponentStatus `json:"beStatus,omitempty"`

	CnStatus *ComponentStatus `json:"cnStatus,omitempty"`

	BrokerStatus *ComponentStatus `json:"brokerStatus,omitempty"`
}

type ComponentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// StarRocksComponentStatus represents the status of a starrocks component.
	//the name of fe service exposed for user.
	AccessService string `json:"accessService,omitempty"`

	//FailedInstances failed pod names.
	FailedMembers []string `json:"failedInstances,omitempty"`

	//CreatingInstances in creating pod names.
	CreatingMembers []string `json:"creatingInstances,omitempty"`

	//RunningInstances in running status pod names.
	RunningMembers []string `json:"runningInstances,omitempty"`

	ComponentCondition ComponentCondition `json:"componentCondition"`
}

type ComponentCondition struct {
	SubResourceName string `json:"subResourceName,omitempty"`
	// Phase of statefulset condition.
	Phase ComponentPhase `json:"phase"`
	// The last time this condition was updated.
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
	// The reason for the condition's last transition.
	Reason string `json:"reason"`
	// A human readable message indicating details about the transition.
	Message string `json:"message"`
}

type ComponentPhase string

const (
	Reconciling      ComponentPhase = "reconciling"
	WaitScheduling   ComponentPhase = "waitScheduling"
	HaveMemberFailed ComponentPhase = "haveMemberFailed"
	Available        ComponentPhase = "available"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=dcr
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="FeStatus",type=string,JSONPath=`.status.feStatus.componentCondition.phase`
// +kubebuilder:printcolumn:name="BeStatus",type=string,JSONPath=`.status.beStatus.componentCondition.phase`
// +kubebuilder:storageversion
// +genclient
// DorisCluster is the Schema for the dorisclusters API
type DorisCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DorisClusterSpec   `json:"spec,omitempty"`
	Status DorisClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DorisClusterList contains a list of DorisCluster
type DorisClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DorisCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DorisCluster{}, &DorisClusterList{})
}