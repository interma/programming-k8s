/*
Copyright 2018 interma.

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
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/interma/programming-k8s/pkg/apis/stats/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// PodStatsLister helps list PodStatses.
type PodStatsLister interface {
	// List lists all PodStatses in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.PodStats, err error)
	// PodStatses returns an object that can list and get PodStatses.
	PodStatses(namespace string) PodStatsNamespaceLister
	PodStatsListerExpansion
}

// podStatsLister implements the PodStatsLister interface.
type podStatsLister struct {
	indexer cache.Indexer
}

// NewPodStatsLister returns a new PodStatsLister.
func NewPodStatsLister(indexer cache.Indexer) PodStatsLister {
	return &podStatsLister{indexer: indexer}
}

// List lists all PodStatses in the indexer.
func (s *podStatsLister) List(selector labels.Selector) (ret []*v1alpha1.PodStats, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.PodStats))
	})
	return ret, err
}

// PodStatses returns an object that can list and get PodStatses.
func (s *podStatsLister) PodStatses(namespace string) PodStatsNamespaceLister {
	return podStatsNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// PodStatsNamespaceLister helps list and get PodStatses.
type PodStatsNamespaceLister interface {
	// List lists all PodStatses in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.PodStats, err error)
	// Get retrieves the PodStats from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.PodStats, error)
	PodStatsNamespaceListerExpansion
}

// podStatsNamespaceLister implements the PodStatsNamespaceLister
// interface.
type podStatsNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all PodStatses in the indexer for a given namespace.
func (s podStatsNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.PodStats, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.PodStats))
	})
	return ret, err
}

// Get retrieves the PodStats from the indexer for a given namespace and name.
func (s podStatsNamespaceLister) Get(name string) (*v1alpha1.PodStats, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("podstats"), name)
	}
	return obj.(*v1alpha1.PodStats), nil
}
