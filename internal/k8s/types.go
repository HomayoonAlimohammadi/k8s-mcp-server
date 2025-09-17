package k8s

// PodInfo represents simplified pod information
type PodInfo struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Status    string            `json:"status"`
	Ready     string            `json:"ready"`
	Restarts  int32             `json:"restarts"`
	Age       string            `json:"age"`
	Labels    map[string]string `json:"labels,omitempty"`
	NodeName  string            `json:"node_name,omitempty"`
	PodIP     string            `json:"pod_ip,omitempty"`
}

// ServiceInfo represents simplified service information
type ServiceInfo struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Type        string            `json:"type"`
	ClusterIP   string            `json:"cluster_ip"`
	ExternalIPs []string          `json:"external_ips,omitempty"`
	Ports       []string          `json:"ports"`
	Age         string            `json:"age"`
	Labels      map[string]string `json:"labels,omitempty"`
	Selector    map[string]string `json:"selector,omitempty"`
}

// DeploymentInfo represents simplified deployment information
type DeploymentInfo struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Ready     string            `json:"ready"`
	UpToDate  int32             `json:"up_to_date"`
	Available int32             `json:"available"`
	Age       string            `json:"age"`
	Labels    map[string]string `json:"labels,omitempty"`
	Replicas  int32             `json:"replicas"`
}

// NamespaceInfo represents simplified namespace information
type NamespaceInfo struct {
	Name   string            `json:"name"`
	Status string            `json:"status"`
	Age    string            `json:"age"`
	Labels map[string]string `json:"labels,omitempty"`
}