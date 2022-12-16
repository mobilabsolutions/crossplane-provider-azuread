/*
Copyright 2022 Upbound Inc.
*/

// Code generated by upjet. DO NOT EDIT.

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

type PermissionGrantObservation struct {
	ID *string `json:"id,omitempty" tf:"id,omitempty"`
}

type PermissionGrantParameters struct {

	// A set of claim values for delegated permission scopes which should be included in access tokens for the resource
	// +kubebuilder:validation:Required
	ClaimValues []*string `json:"claimValues" tf:"claim_values,omitempty"`

	// The object ID of the service principal representing the resource to be accessed
	// +crossplane:generate:reference:type=github.com/upbound/provider-azuread/apis/serviceprincipals/v1beta1.Principal
	// +kubebuilder:validation:Optional
	ResourceServicePrincipalObjectID *string `json:"resourceServicePrincipalObjectId,omitempty" tf:"resource_service_principal_object_id,omitempty"`

	// Reference to a Principal in serviceprincipals to populate resourceServicePrincipalObjectId.
	// +kubebuilder:validation:Optional
	ResourceServicePrincipalObjectIDRef *v1.Reference `json:"resourceServicePrincipalObjectIdRef,omitempty" tf:"-"`

	// Selector for a Principal in serviceprincipals to populate resourceServicePrincipalObjectId.
	// +kubebuilder:validation:Optional
	ResourceServicePrincipalObjectIDSelector *v1.Selector `json:"resourceServicePrincipalObjectIdSelector,omitempty" tf:"-"`

	// The object ID of the service principal for which this delegated permission grant should be created
	// +crossplane:generate:reference:type=github.com/upbound/provider-azuread/apis/serviceprincipals/v1beta1.Principal
	// +kubebuilder:validation:Optional
	ServicePrincipalObjectID *string `json:"servicePrincipalObjectId,omitempty" tf:"service_principal_object_id,omitempty"`

	// Reference to a Principal in serviceprincipals to populate servicePrincipalObjectId.
	// +kubebuilder:validation:Optional
	ServicePrincipalObjectIDRef *v1.Reference `json:"servicePrincipalObjectIdRef,omitempty" tf:"-"`

	// Selector for a Principal in serviceprincipals to populate servicePrincipalObjectId.
	// +kubebuilder:validation:Optional
	ServicePrincipalObjectIDSelector *v1.Selector `json:"servicePrincipalObjectIdSelector,omitempty" tf:"-"`

	// The object ID of the user on behalf of whom the service principal is authorized to access the resource
	// +kubebuilder:validation:Optional
	UserObjectID *string `json:"userObjectId,omitempty" tf:"user_object_id,omitempty"`
}

// PermissionGrantSpec defines the desired state of PermissionGrant
type PermissionGrantSpec struct {
	v1.ResourceSpec `json:",inline"`
	ForProvider     PermissionGrantParameters `json:"forProvider"`
}

// PermissionGrantStatus defines the observed state of PermissionGrant.
type PermissionGrantStatus struct {
	v1.ResourceStatus `json:",inline"`
	AtProvider        PermissionGrantObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// PermissionGrant is the Schema for the PermissionGrants API. <no value>
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,azuread}
type PermissionGrant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PermissionGrantSpec   `json:"spec"`
	Status            PermissionGrantStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PermissionGrantList contains a list of PermissionGrants
type PermissionGrantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PermissionGrant `json:"items"`
}

// Repository type metadata.
var (
	PermissionGrant_Kind             = "PermissionGrant"
	PermissionGrant_GroupKind        = schema.GroupKind{Group: CRDGroup, Kind: PermissionGrant_Kind}.String()
	PermissionGrant_KindAPIVersion   = PermissionGrant_Kind + "." + CRDGroupVersion.String()
	PermissionGrant_GroupVersionKind = CRDGroupVersion.WithKind(PermissionGrant_Kind)
)

func init() {
	SchemeBuilder.Register(&PermissionGrant{}, &PermissionGrantList{})
}
