/*
Copyright 2018 The Kubernetes Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodePortrait is the future resource usage portrait of node
type NodePortrait struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   NodePortraitSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status NodePortraitStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// NodePortraitSpec is the specification of a NodePortrait.
type NodePortraitSpec struct {
	Trait map[string]string `json:"trait,omitempty"`
}

// NodePortraitStatus is the status of a NodePortrait.
type NodePortraitStatus struct {
	// len=48
	CpuDistributions []int64 `json:"cpuDistributions,omitempty"`
	MemDistributions []int64 `json:"memDistributions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodePortraitList is a list of NodePortrait objects.
type NodePortraitList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []NodePortrait `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SduPortrait is an example type with a spec and a status.
type SduPortrait struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   SduPortraitSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status SduPortraitStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SduPortraitList is a list of SduPortrait objects.
type SduPortraitList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []SduPortrait `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// SduPortraitSpec is the specification of a SduPortrait.
type SduPortraitSpec struct {
	ClusterName   string            `json:"clusterName,omitempty"`
	Type          string            `json:"type,omitempty"`
	DataSource    string            `json:"dataSource,omitempty"`
	PromotionName string            `json:"promotionName"`
	Date          string            `json:"date"`
	Peak          string            `json:"peak,omitempty"`
	Trait         map[string]string `json:"trait"`
}

// SduPortraitStatus is the status of a SduPortrait.
type SduPortraitStatus struct {
	ResourceCharacteristics ResourceCharacteristics `json:"resourceCharacteristics,omitempty"`
	Predict                 Predict                 `json:"predict,omitempty"`
}

type ResourceCharacteristics struct {
	// key=containerName
	ContainersRecommend map[string]PercentileContainer `json:"containersRecommend,omitempty"`
}

type PercentileContainer struct {
	// key = percentile value ;  key2 = "cpu/memory"
	Percentile map[string]map[string]int64 `json:"percentile"`
}

type Predict struct {
	// len=48
	CpuDistributions []int64 `json:"cpuDistributions,omitempty"`
	MemDistributions []int64 `json:"memDistributions,omitempty"`
	// len=12
	CpuSumDistributions []int64 `json:"cpuSumDistributions,omitempty"`
	MemSumDistributions []int64 `json:"memSumDistributions,omitempty"`
}
