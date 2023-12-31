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

package hub

import (
	"sync"

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"kmodules.xyz/resource-metadata/hub/resourcedescriptors"
)

type KV interface {
	Set(key string, val *v1alpha1.ResourceDescriptor)
	Get(key string) (*v1alpha1.ResourceDescriptor, bool)
	Visit(func(key string, val *v1alpha1.ResourceDescriptor))
}

// KVMap is concurrent safe.
type KVMap struct {
	cache map[string]*v1alpha1.ResourceDescriptor
	m     sync.RWMutex
}

var _ KV = &KVMap{}

func NewKVMap(cache map[string]*v1alpha1.ResourceDescriptor) KV {
	return &KVMap{cache: cache}
}

func (s *KVMap) Set(key string, val *v1alpha1.ResourceDescriptor) {
	s.m.Lock()
	s.cache[key] = val
	s.m.Unlock()
}

func (s *KVMap) Get(key string) (*v1alpha1.ResourceDescriptor, bool) {
	s.m.RLock()
	val, found := s.cache[key]
	s.m.RUnlock()
	return val, found
}

func (s *KVMap) Visit(f func(key string, val *v1alpha1.ResourceDescriptor)) {
	s.m.RLock()
	for k, v := range s.cache {
		f(k, v)
	}
	s.m.RUnlock()
}

// KVLocal is not concurrent safe.
type KVLocal struct {
	known KV
	cache map[string]*v1alpha1.ResourceDescriptor
}

var _ KV = &KVLocal{}

func NewKVLocal() KV {
	return &KVLocal{
		known: &KVMap{
			cache: resourcedescriptors.KnownDescriptors(),
		},
		cache: map[string]*v1alpha1.ResourceDescriptor{},
	}
}

func (s *KVLocal) Set(key string, val *v1alpha1.ResourceDescriptor) {
	if _, found := s.known.Get(key); !found {
		s.cache[key] = val
	}
}

func (s *KVLocal) Get(key string) (*v1alpha1.ResourceDescriptor, bool) {
	val, found := s.known.Get(key)
	if found {
		return val, found
	}
	val, found = s.cache[key]
	return val, found
}

func (s *KVLocal) Visit(f func(key string, val *v1alpha1.ResourceDescriptor)) {
	s.known.Visit(f)
	for k, v := range s.cache {
		f(k, v)
	}
}

const KnownUID = "__known__"
