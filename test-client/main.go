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

package main

import (
	"context"
	"fmt"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2/klogr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

func NewClient() (client.Client, error) {
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)

	ctrl.SetLogger(klogr.New()) // nolint:staticcheck
	cfg := ctrl.GetConfigOrDie()
	cfg.QPS = 100
	cfg.Burst = 100

	// WARNING: must be set to application/json to avoid accidentally using protobuf encoding with fake-apiserver
	cfg.AcceptContentTypes = runtime.ContentTypeJSON

	hc, err := rest.HTTPClientFor(cfg)
	if err != nil {
		return nil, err
	}
	mapper, err := apiutil.NewDynamicRESTMapper(cfg, hc)
	if err != nil {
		return nil, err
	}

	return client.New(cfg, client.Options{
		Scheme: scheme,
		Mapper: mapper,
		//Opts: client.WarningHandlerOptions{
		//	SuppressWarnings:   false,
		//	AllowDuplicateLogs: false,
		//},
	})
}

func main() {
	if err := useGeneratedClient(); err != nil {
		panic(err)
	}
	if err := useKubebuilderClient(); err != nil {
		panic(err)
	}
}

func useGeneratedClient() error {
	fmt.Println("Using Generated client")
	cfg := ctrl.GetConfigOrDie()
	cfg.QPS = 100
	cfg.Burst = 100

	kc, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}

	var list *apps.DeploymentList
	list, err = kc.AppsV1().Deployments(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, db := range list.Items {
		fmt.Println(client.ObjectKeyFromObject(&db))
	}
	return nil
}

func useKubebuilderClient() error {
	fmt.Println("Using kubebuilder client")
	kc, err := NewClient()
	if err != nil {
		return err
	}

	var u unstructured.Unstructured
	raw := `{
  "apiVersion": "v1",
  "data": {
    "password": "UyFCXCpkJHpEc2I9",
    "username": "YWRtaW4="
  },
  "kind": "Secret",
  "metadata": {
    "creationTimestamp": null,
    "name": "db-user-pass2",
    "namespace": "default"
  }
}`
	err = json.Unmarshal([]byte(raw), &u)
	if err != nil {
		return err
	}
	err = kc.Create(context.TODO(), &u)
	if err != nil {
		return err
	}

	var sec core.Secret
	if err = kc.Get(context.TODO(), types.NamespacedName{Name: "db-user-pass2", Namespace: "default"}, &sec); err != nil {
		return err
	}
	fmt.Println(client.ObjectKeyFromObject(&sec))

	var list core.SecretList
	err = kc.List(context.TODO(), &list)
	if err != nil {
		return err
	}
	for _, items := range list.Items {
		fmt.Println(client.ObjectKeyFromObject(&items))
	}
	return nil
}
