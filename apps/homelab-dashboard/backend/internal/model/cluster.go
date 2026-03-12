package model

type Node struct {
	Name      string            `json:"name"`
	Status    string            `json:"status"`
	Roles     []string          `json:"roles"`
	IP        string            `json:"ip"`
	CPU       ResourceUsage     `json:"cpu"`
	Memory    ResourceUsage     `json:"memory"`
	Labels    map[string]string `json:"labels"`
	CreatedAt string            `json:"createdAt"`
}

type ResourceUsage struct {
	Capacity    string `json:"capacity"`
	Allocatable string `json:"allocatable"`
}

type Namespace struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

type Pod struct {
	Name       string            `json:"name"`
	Namespace  string            `json:"namespace"`
	Status     string            `json:"status"`
	Node       string            `json:"node"`
	IP         string            `json:"ip"`
	Containers []Container       `json:"containers"`
	Labels     map[string]string `json:"labels"`
	CreatedAt  string            `json:"createdAt"`
}

type Container struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	Ready bool   `json:"ready"`
}

type Service struct {
	Name       string            `json:"name"`
	Namespace  string            `json:"namespace"`
	Type       string            `json:"type"`
	ClusterIP  string            `json:"clusterIP"`
	Ports      []ServicePort     `json:"ports"`
	Labels     map[string]string `json:"labels"`
	CreatedAt  string            `json:"createdAt"`
}

type ServicePort struct {
	Name       string `json:"name"`
	Port       int32  `json:"port"`
	TargetPort string `json:"targetPort"`
	Protocol   string `json:"protocol"`
}
