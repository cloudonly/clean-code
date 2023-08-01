package main

import (
	"github.com/luomu/clean-code/pkg/k8s/client"
	"github.com/luomu/clean-code/pkg/k8s/clustercache"
	log "github.com/luomu/clean-code/pkg/logging/zerolog"
)

func getAllResource() {
	cache := getClusterCache()
	nodes := cache.GetAllNodes()
	for _, node := range nodes {
		log.Infof("[Node] - %s", node.GetName())
	}

	pods := cache.GetAllPods()
	for _, pod := range pods {
		log.Infof("[Pod] - %s/%s", pod.GetNamespace(), pod.GetName())
	}
}

func getClusterCache() clustercache.ClusterCache {
	kubeClient, err := client.LoadKubeClient("")
	if err != nil {
		log.Errorf("failed to load kubernetes client.")
		return nil
	}
	cache := clustercache.NewKubernetesClusterCache(kubeClient)
	cache.Run()
	return cache
}
