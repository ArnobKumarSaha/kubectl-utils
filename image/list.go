package image

import (
	"context"
	"fmt"
	"github.com/Arnobkumarsaha/kubectl-utils/client"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"
	"os"
	"strings"
	"text/tabwriter"
)

type info struct {
	image, id, name string
}

var (
	podMap    map[string]map[string][]info
	ownerMap  map[string][]string
	targetGVR schema.GroupVersionResource
)

func listAll() error {
	arr := strings.Split(resource, ",")
	for _, r := range arr {
		err := list(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func list(curRes string) error {
	podGvr := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}
	targetGVR = podGvr
	targetGVR.Group = "apps"
	switch curRes {
	case "deployments", "deployment", "deploy", "dep":
		targetGVR.Resource = "deployments"
	case "replicasets", "replicaset", "rs":
		targetGVR.Resource = "replicasets"
	case "statefulsets", "statefulset", "sts":
		targetGVR.Resource = "statefulsets"
	case "daemonsets", "daemonset", "ds":
		targetGVR.Resource = "daemonsets"
	default:
		klog.Errorf("unknown resource %s", curRes)
	}
	var (
		pods *corev1.PodList
		err  error
	)
	podMap = make(map[string]map[string][]info)
	ownerMap = make(map[string][]string)
	pods, err = client.Client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, pod := range pods.Items {
		podNames, exists := podMap[pod.Namespace]
		if !exists {
			podNames = make(map[string][]info)
		}
		podInfo, exists := podNames[pod.Name]
		if !exists {
			podInfo = make([]info, 0)
		}
		for _, c := range pod.Status.ContainerStatuses {
			inf := info{
				image: c.Image,
				id:    c.ImageID,
				name:  c.Name,
			}
			podInfo = append(podInfo, inf)
		}
		podNames[pod.Name] = podInfo
		podMap[pod.Namespace] = podNames

		res, err := findResource(pod.OwnerReferences, pod.Namespace)
		//klog.Infof("%v/%v %v %v \n", pod.Namespace, pod.Name, res == nil, err)
		if err != nil {
			return err
		}
		if res == nil {
			continue
		}
		s := fmt.Sprintf("%s/%s", res.GetNamespace(), res.GetName())
		if _, exists = ownerMap[s]; !exists {
			ownerMap[s] = make([]string, 0)
		}
		ownerMap[s] = append(ownerMap[s], pod.GetName())
	}
	output()
	return nil
}

func findResource(ownerRefs []metav1.OwnerReference, ns string) (*unstructured.Unstructured, error) {
	for _, ownerRef := range ownerRefs {
		gvr, scope, err := getGVRFromOwnerRef(ownerRef)
		if err != nil {
			return nil, err
		}

		if gvr != targetGVR {
			refs, err := getOwnerRefs(gvr, scope, ns, ownerRef.Name)
			if err != nil {
				klog.Errorf("get owner refs err: %v", err)
				return nil, err
			}
			return findResource(refs, ns)
		}

		target, err := client.DynamicClient.Resource(gvr).Namespace(ns).Get(context.TODO(), ownerRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to get target: %v", err)
		}
		return target, nil
	}

	return nil, nil // No targetRef found in the ownership chain
}

func getOwnerRefs(gvr schema.GroupVersionResource, scope meta.RESTScope, ns, name string) ([]metav1.OwnerReference, error) {
	var (
		rs  *unstructured.Unstructured
		err error
	)
	if scope.Name() == "namespace" {
		rs, err = client.DynamicClient.Resource(gvr).Namespace(ns).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to get %v: %v", gvr.Resource, err)
		}
	} else {
		rs, err = client.DynamicClient.Resource(gvr).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to get in cluster scope %v: %v", gvr.Resource, err)
		}
	}
	slice, found, err := unstructured.NestedSlice(rs.UnstructuredContent(), "metadata", "ownerReferences")
	if err != nil {
		return nil, fmt.Errorf("failed to extract ownerReferences: %v", err)
	}
	if !found {
		return nil, nil
	}

	// Convert the slice to []metav1.OwnerReference
	var ownerReferences []metav1.OwnerReference
	for _, item := range slice {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to cast ownerReference item to map")
		}
		jsonData, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&itemMap)
		if err != nil {
			return nil, fmt.Errorf("failed to convert ownerReference to JSON: %v", err)
		}
		var ownerRef metav1.OwnerReference
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(jsonData, &ownerRef); err != nil {
			return nil, fmt.Errorf("failed to unmarshal ownerReference: %v", err)
		}
		ownerReferences = append(ownerReferences, ownerRef)
	}
	return ownerReferences, nil
}

func output() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "%s\tInfo\t\n", targetGVR.Resource)
	// Iterate over the ownerMap
	for owner, podNames := range ownerMap {
		refs := strings.Split(owner, "/")
		if len(podNames) > 0 {
			infos := podMap[refs[0]][podNames[0]]
			fmt.Fprintf(w, "%v\t%v\t\n", owner, infos[0])
		}
	}
	w.Flush()
}
