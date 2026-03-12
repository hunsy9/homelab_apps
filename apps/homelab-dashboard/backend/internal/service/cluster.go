package service

import (
	"context"
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"

	"homelab-dashboard/internal/model"
)

type nodeStatsSummary struct {
	Node struct {
		Fs struct {
			CapacityBytes  int64 `json:"capacityBytes"`
			UsedBytes      int64 `json:"usedBytes"`
			AvailableBytes int64 `json:"availableBytes"`
		} `json:"fs"`
	} `json:"node"`
}

type ClusterService struct {
	client        *kubernetes.Clientset
	metricsClient *metricsv.Clientset
}

func NewClusterService(client *kubernetes.Clientset, metricsClient *metricsv.Clientset) *ClusterService {
	return &ClusterService{client: client, metricsClient: metricsClient}
}

func (s *ClusterService) GetNodes(ctx context.Context) ([]model.Node, error) {
	nodes, err := s.client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Fetch CPU/Memory metrics
	metricsMap := make(map[string]struct {
		cpuUsage    int64
		memoryUsage int64
	})

	if s.metricsClient != nil {
		nodeMetrics, err := s.metricsClient.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
		if err == nil {
			for _, m := range nodeMetrics.Items {
				metricsMap[m.Name] = struct {
					cpuUsage    int64
					memoryUsage int64
				}{
					cpuUsage:    m.Usage.Cpu().MilliValue(),
					memoryUsage: m.Usage.Memory().Value(),
				}
			}
		}
	}

	// Fetch storage stats from kubelet
	storageMap := make(map[string]model.StorageInfo)
	for _, n := range nodes.Items {
		stats := s.getNodeStats(ctx, n.Name)
		if stats != nil {
			storageMap[n.Name] = model.StorageInfo{
				Capacity: stats.Node.Fs.CapacityBytes,
				Usage:    stats.Node.Fs.UsedBytes,
			}
		}
	}

	result := make([]model.Node, 0, len(nodes.Items))
	for _, n := range nodes.Items {
		metrics := metricsMap[n.Name]
		storage := storageMap[n.Name]

		node := model.Node{
			Name:      n.Name,
			Status:    getNodeStatus(n),
			Roles:     getNodeRoles(n),
			IP:        getNodeIP(n),
			Labels:    n.Labels,
			CreatedAt: n.CreationTimestamp.Format("2006-01-02T15:04:05Z"),
			CPU: model.ResourceInfo{
				Allocatable: n.Status.Allocatable.Cpu().MilliValue(),
				Usage:       metrics.cpuUsage,
			},
			Memory: model.ResourceInfo{
				Allocatable: n.Status.Allocatable.Memory().Value(),
				Usage:       metrics.memoryUsage,
			},
			Storage: storage,
		}
		result = append(result, node)
	}
	return result, nil
}

func (s *ClusterService) getNodeStats(ctx context.Context, nodeName string) *nodeStatsSummary {
	data, err := s.client.CoreV1().RESTClient().Get().
		Resource("nodes").
		Name(nodeName).
		SubResource("proxy").
		Suffix("stats/summary").
		Do(ctx).Raw()
	if err != nil {
		return nil
	}

	var stats nodeStatsSummary
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil
	}
	return &stats
}

func (s *ClusterService) GetNamespaces(ctx context.Context) ([]model.Namespace, error) {
	namespaces, err := s.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make([]model.Namespace, 0, len(namespaces.Items))
	for _, ns := range namespaces.Items {
		result = append(result, model.Namespace{
			Name:      ns.Name,
			Status:    string(ns.Status.Phase),
			CreatedAt: ns.CreationTimestamp.Format("2006-01-02T15:04:05Z"),
		})
	}
	return result, nil
}

func (s *ClusterService) GetPods(ctx context.Context, namespace string) ([]model.Pod, error) {
	pods, err := s.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make([]model.Pod, 0, len(pods.Items))
	for _, p := range pods.Items {
		containers := make([]model.Container, 0, len(p.Spec.Containers))
		for i, c := range p.Spec.Containers {
			ready := false
			if i < len(p.Status.ContainerStatuses) {
				ready = p.Status.ContainerStatuses[i].Ready
			}
			containers = append(containers, model.Container{
				Name:  c.Name,
				Image: c.Image,
				Ready: ready,
			})
		}

		result = append(result, model.Pod{
			Name:       p.Name,
			Namespace:  p.Namespace,
			Status:     string(p.Status.Phase),
			Node:       p.Spec.NodeName,
			IP:         p.Status.PodIP,
			Containers: containers,
			Labels:     p.Labels,
			CreatedAt:  p.CreationTimestamp.Format("2006-01-02T15:04:05Z"),
		})
	}
	return result, nil
}

func (s *ClusterService) GetServices(ctx context.Context, namespace string) ([]model.Service, error) {
	services, err := s.client.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make([]model.Service, 0, len(services.Items))
	for _, svc := range services.Items {
		ports := make([]model.ServicePort, 0, len(svc.Spec.Ports))
		for _, p := range svc.Spec.Ports {
			ports = append(ports, model.ServicePort{
				Name:       p.Name,
				Port:       p.Port,
				TargetPort: p.TargetPort.String(),
				Protocol:   string(p.Protocol),
			})
		}

		result = append(result, model.Service{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			Type:      string(svc.Spec.Type),
			ClusterIP: svc.Spec.ClusterIP,
			Ports:     ports,
			Labels:    svc.Labels,
			CreatedAt: svc.CreationTimestamp.Format("2006-01-02T15:04:05Z"),
		})
	}
	return result, nil
}

func getNodeStatus(node corev1.Node) string {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			if condition.Status == corev1.ConditionTrue {
				return "Ready"
			}
			return "NotReady"
		}
	}
	return "Unknown"
}

func getNodeRoles(node corev1.Node) []string {
	roles := []string{}
	for label := range node.Labels {
		if label == "node-role.kubernetes.io/control-plane" {
			roles = append(roles, "control-plane")
		}
		if label == "node-role.kubernetes.io/worker" {
			roles = append(roles, "worker")
		}
	}
	if len(roles) == 0 {
		roles = append(roles, "worker")
	}
	return roles
}

func getNodeIP(node corev1.Node) string {
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeInternalIP {
			return addr.Address
		}
	}
	return ""
}
