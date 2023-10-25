/*
Copyright AppsCode Inc. and Contributors.

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

package resources

import (
	"fmt"

	"kmodules.xyz/client-go/apiextensions"

	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/rest"
	driversapi "x-helm.dev/apimachinery/apis/drivers/v1alpha1"
)

func RegisterCRDs(restcfg *rest.Config) error {
	// register AppRelease CRD
	crds := []*apiextensions.CustomResourceDefinition{
		driversapi.AppRelease{}.CustomResourceDefinition(),
	}
	crdClient, err := crd_cs.NewForConfig(restcfg)
	if err != nil {
		return fmt.Errorf("failed to create crd client, reason %v", err)
	}
	err = apiextensions.RegisterCRDs(crdClient, crds)
	if err != nil {
		return fmt.Errorf("failed to register appRelease crd, reason %v", err)
	}
	return nil
}
