/*
Copyright 2025.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type DeploySyncerSpec struct {
    // RawFileUrl is the GitHub repository URL
    RawFileUrl string `json:"RawFileUrl"`
    // IntervalSeconds is the sync interval in seconds
    IntervalSeconds int32 `json:"intervalSeconds"`
}

type DeploySyncerStatus struct {
    // LastStatus shows the last sync status
    LastStatus string `json:"lastStatus"`
    // LastSyncTime shows when the last sync occurred
    LastSyncTime string `json:"lastSyncTime"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// DeploySyncer is the Schema for the deploysyncers API.
type DeploySyncer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeploySyncerSpec   `json:"spec,omitempty"`
	Status DeploySyncerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DeploySyncerList contains a list of DeploySyncer.
type DeploySyncerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DeploySyncer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DeploySyncer{}, &DeploySyncerList{})
}
