/*
Copyright 2024 The Kubernetes Authors.

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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1alpha1 "github.com/openshift/node-feature-discovery/api/nfd/v1alpha1"
)

// NodeFeatureGroupLister helps list NodeFeatureGroups.
// All objects returned here must be treated as read-only.
type NodeFeatureGroupLister interface {
	// List lists all NodeFeatureGroups in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.NodeFeatureGroup, err error)
	// NodeFeatureGroups returns an object that can list and get NodeFeatureGroups.
	NodeFeatureGroups(namespace string) NodeFeatureGroupNamespaceLister
	NodeFeatureGroupListerExpansion
}

// nodeFeatureGroupLister implements the NodeFeatureGroupLister interface.
type nodeFeatureGroupLister struct {
	indexer cache.Indexer
}

// NewNodeFeatureGroupLister returns a new NodeFeatureGroupLister.
func NewNodeFeatureGroupLister(indexer cache.Indexer) NodeFeatureGroupLister {
	return &nodeFeatureGroupLister{indexer: indexer}
}

// List lists all NodeFeatureGroups in the indexer.
func (s *nodeFeatureGroupLister) List(selector labels.Selector) (ret []*v1alpha1.NodeFeatureGroup, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.NodeFeatureGroup))
	})
	return ret, err
}

// NodeFeatureGroups returns an object that can list and get NodeFeatureGroups.
func (s *nodeFeatureGroupLister) NodeFeatureGroups(namespace string) NodeFeatureGroupNamespaceLister {
	return nodeFeatureGroupNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// NodeFeatureGroupNamespaceLister helps list and get NodeFeatureGroups.
// All objects returned here must be treated as read-only.
type NodeFeatureGroupNamespaceLister interface {
	// List lists all NodeFeatureGroups in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.NodeFeatureGroup, err error)
	// Get retrieves the NodeFeatureGroup from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.NodeFeatureGroup, error)
	NodeFeatureGroupNamespaceListerExpansion
}

// nodeFeatureGroupNamespaceLister implements the NodeFeatureGroupNamespaceLister
// interface.
type nodeFeatureGroupNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all NodeFeatureGroups in the indexer for a given namespace.
func (s nodeFeatureGroupNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.NodeFeatureGroup, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.NodeFeatureGroup))
	})
	return ret, err
}

// Get retrieves the NodeFeatureGroup from the indexer for a given namespace and name.
func (s nodeFeatureGroupNamespaceLister) Get(name string) (*v1alpha1.NodeFeatureGroup, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("nodefeaturegroup"), name)
	}
	return obj.(*v1alpha1.NodeFeatureGroup), nil
}
