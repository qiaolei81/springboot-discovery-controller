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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SpringBootServerSpec defines the desired state of SpringBootServer
type SpringBootServerSpec struct {
	// Server is the target server name or ip address to discover of SpringBootServer.
	Server string `json:"server,omitempty"`
}

// SpringBootServerStatus defines the observed state of SpringBootServer
type SpringBootServerStatus struct {
	// Status is the discovery state which could be Succeeded or Failed
	Status string `json:"status,omitempty"`
	// Message is the useful message to describe discovery state, especially it's in failed state
	Message string `json:"message,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Created",type="date",JSONPath=".metadata.creationTimestamp"
//+kubebuilder:printcolumn:name="Server",type="string",JSONPath=".spec.server"
//+kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status"
//+kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.message"

// SpringBootServer is the Schema for the springbootservers API
type SpringBootServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SpringBootServerSpec   `json:"spec,omitempty"`
	Status SpringBootServerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SpringBootServerList contains a list of SpringBootServer
type SpringBootServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SpringBootServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SpringBootServer{}, &SpringBootServerList{})
}
