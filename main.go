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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kmodules.xyz/client-go/tools/clientcmd"
	"kmodules.xyz/fake-apiserver/pkg"
	"kmodules.xyz/fake-apiserver/pkg/resources"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	s := pkg.NewServer(pkg.NewOptions())
	srv, restcfg, err := s.Run()
	if err != nil {
		klog.Fatalln(err)
	}
	klog.Infoln("Server Started")

	kubecfg, err := clientcmd.BuildKubeConfigBytes(restcfg, metav1.NamespaceDefault)
	if err != nil {
		klog.Fatalln(err)
	}
	err = os.WriteFile("local.kubeconfig", kubecfg, 0o640)
	if err != nil {
		klog.Fatalln(err)
	}

	err = resources.InitCluster(restcfg)
	if err != nil {
		klog.Fatalln(err)
	}
	s.Checkpoint()

	/*
		mapper, err := apiutil.NewDiscoveryRESTMapper(restcfg)
		if err != nil {
			klog.Fatalln(err)
		}
		kinds, err := mapper.RESTMappings(schema.GroupKind{
			Kind: "Namespace",
		})
		if err != nil {
			klog.Fatalln(err)
		}
		fmt.Println(kinds)
	*/

	kc := kubernetes.NewForConfigOrDie(restcfg)
	ns := &core.Namespace{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo",
		},
	}
	ns, err = kc.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	if err != nil {
		panic("failed to create ns" + ns.Name)
	}

	ns.Labels = map[string]string{
		"tes": "abc",
	}
	ns, err = kc.CoreV1().Namespaces().Update(context.TODO(), ns, metav1.UpdateOptions{})
	if err != nil {
		panic("failed to update ns" + ns.Name)
	}

	go func() {
		time.Sleep(10 * time.Second)

		sel, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
			MatchLabels: map[string]string{
				"id": "test",
			},
		})
		if err != nil {
			panic(err)
		}
		err = kc.CoreV1().ConfigMaps("default").DeleteCollection(context.TODO(), metav1.DeleteOptions{}, metav1.ListOptions{
			LabelSelector: sel.String(),
		})
		if err != nil {
			panic(err)
		}
	}()

	<-done
	klog.Infoln("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	klog.Infoln("Server Exited Properly")

	current, deleted := s.Export()
	fmt.Println("CURRENT objects __________________________")
	data, err := json.MarshalIndent(current, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	fmt.Println("DELETED objects __________________________")
	data, err = json.MarshalIndent(deleted, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
