package image

import (
	"fmt"
	"github.com/Arnobkumarsaha/kubectl-utils/client"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/restmapper"
)

func getGVRFromOwnerRef(own metav1.OwnerReference) (schema.GroupVersionResource, meta.RESTScope, error) {
	return getGVR(own.APIVersion, own.Kind)
}

func getGVR(apiVersion, kind string) (schema.GroupVersionResource, meta.RESTScope, error) {
	gv, err := schema.ParseGroupVersion(apiVersion)
	if err != nil {
		return schema.GroupVersionResource{}, nil, fmt.Errorf("failed to parse GroupVersion: %v", err)
	}
	gvk := gv.WithKind(kind)

	// Retrieve the RESTMapper from the discovery client
	groupResources, err := restmapper.GetAPIGroupResources(client.DiscoveryClient)
	if err != nil {
		return schema.GroupVersionResource{}, nil, fmt.Errorf("failed to get API group resources: %v", err)
	}

	// Create a RESTMapper from the discovered group resources
	restMapper := restmapper.NewDiscoveryRESTMapper(groupResources)
	mapping, err := restMapper.RESTMapping(gvk.GroupKind(), gv.Version)
	if err != nil {
		return schema.GroupVersionResource{}, nil, fmt.Errorf("failed to map GVK to GVR: %v", err)
	}
	return mapping.Resource, mapping.Scope, nil
}
