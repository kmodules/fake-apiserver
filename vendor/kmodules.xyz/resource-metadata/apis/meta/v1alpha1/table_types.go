/*
Copyright AppsCode Inc. and Contributors

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

// Table is a tabular representation of a set of API resources. The server transforms the
// object into a set of preferred columns for quickly reviewing the objects.

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Table struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// columns describes each column in the returned items array. The number of cells per row
	// will always match the number of column definitions.
	Columns []ResourceColumn `json:"columns"`
	// rows is the list of items in the table.
	Rows []TableRow `json:"rows"`
}

// TableRow is an individual row in a table.
type TableRow struct {
	// cells will be as wide as the column definitions array and may contain strings, numbers (float64 or
	// int64), booleans, simple maps, lists, or null. See the type field of the column definition for a
	// more detailed description.
	Cells []TableCell `json:"cells"`
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

type TableCell struct {
	// cells will be as wide as the column definitions array and may contain strings, numbers (float64 or
	// int64), booleans, simple maps, lists, or null. See the type field of the column definition for a
	// more detailed description.
	Data interface{} `json:"data"`
	// +optional
	Sort interface{} `json:"sort,omitempty"`
	// +optional
	Link string `json:"link,omitempty"`
	// +optional
	Tooltip string `json:"tooltip,omitempty"`
	// +optional
	Icon string `json:"icon,omitempty"`
	// +optional
	Color string `json:"color,omitempty"`
}

// IncludeObjectPolicy controls which portion of the object is returned with a Table.
type IncludeObjectPolicy string

const (
	// IncludeNone returns no object.
	IncludeNone IncludeObjectPolicy = "None"
	// IncludeMetadata serializes the object containing only its metadata field.
	IncludeMetadata IncludeObjectPolicy = "Metadata"
	// IncludeObject contains the full object.
	IncludeObject IncludeObjectPolicy = "Object"
)

// TableOptions are used when a Table is requested by the caller.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TableOptions struct {
	metav1.TypeMeta `json:",inline"`

	// NoHeaders is only exposed for internal callers. It is not included in our OpenAPI definitions
	// and may be removed as a field in a future release.
	NoHeaders bool `json:"-"`

	// includeObject decides whether to include each object along with its columnar information.
	// Specifying "None" will return no object, specifying "Object" will return the full object contents, and
	// specifying "Metadata" (the default) will return the object's metadata in the PartialObjectMetadata kind
	// in version v1beta1 of the meta.k8s.io API group.
	IncludeObject IncludeObjectPolicy `json:"includeObject,omitempty"`
}
