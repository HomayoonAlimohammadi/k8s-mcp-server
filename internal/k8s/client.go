package k8s

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client wraps the Kubernetes client with additional functionality
type Client struct {
	clientset kubernetes.Interface
}

// NewClient creates a new Kubernetes client
func NewClient(kubeconfig string, inCluster bool) (*Client, error) {
	config, err := buildConfig(kubeconfig, inCluster)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubernetes config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return &Client{
		clientset: clientset,
	}, nil
}

// buildConfig builds the Kubernetes configuration
func buildConfig(kubeconfig string, inCluster bool) (*rest.Config, error) {
	if inCluster {
		return rest.InClusterConfig()
	}

	if kubeconfig == "" {
		return nil, fmt.Errorf("kubeconfig path is required when not running in cluster")
	}

	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

// ListPods lists pods in the specified namespace
func (c *Client) ListPods(ctx context.Context, namespace string) ([]PodInfo, error) {
	pods, err := c.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods in namespace %s: %w", namespace, err)
	}

	var podInfos []PodInfo
	for _, pod := range pods.Items {
		podInfo := buildPodInfo(&pod)
		podInfos = append(podInfos, podInfo)
	}

	return podInfos, nil
}

// GetPod gets a specific pod by name and namespace
func (c *Client) GetPod(ctx context.Context, namespace, name string) (*PodInfo, error) {
	pod, err := c.clientset.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod %s in namespace %s: %w", name, namespace, err)
	}

	podInfo := buildPodInfo(pod)
	return &podInfo, nil
}

// ListServices lists services in the specified namespace
func (c *Client) ListServices(ctx context.Context, namespace string) ([]ServiceInfo, error) {
	services, err := c.clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services in namespace %s: %w", namespace, err)
	}

	var serviceInfos []ServiceInfo
	for _, service := range services.Items {
		serviceInfo := buildServiceInfo(&service)
		serviceInfos = append(serviceInfos, serviceInfo)
	}

	return serviceInfos, nil
}

// GetService gets a specific service by name and namespace
func (c *Client) GetService(ctx context.Context, namespace, name string) (*ServiceInfo, error) {
	service, err := c.clientset.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get service %s in namespace %s: %w", name, namespace, err)
	}

	serviceInfo := buildServiceInfo(service)
	return &serviceInfo, nil
}

// ListDeployments lists deployments in the specified namespace
func (c *Client) ListDeployments(ctx context.Context, namespace string) ([]DeploymentInfo, error) {
	deployments, err := c.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments in namespace %s: %w", namespace, err)
	}

	var deploymentInfos []DeploymentInfo
	for _, deployment := range deployments.Items {
		deploymentInfo := buildDeploymentInfo(&deployment)
		deploymentInfos = append(deploymentInfos, deploymentInfo)
	}

	return deploymentInfos, nil
}

// GetDeployment gets a specific deployment by name and namespace
func (c *Client) GetDeployment(ctx context.Context, namespace, name string) (*DeploymentInfo, error) {
	deployment, err := c.clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment %s in namespace %s: %w", name, namespace, err)
	}

	deploymentInfo := buildDeploymentInfo(deployment)
	return &deploymentInfo, nil
}

// ListNamespaces lists all namespaces
func (c *Client) ListNamespaces(ctx context.Context) ([]NamespaceInfo, error) {
	namespaces, err := c.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	var namespaceInfos []NamespaceInfo
	for _, namespace := range namespaces.Items {
		namespaceInfo := buildNamespaceInfo(&namespace)
		namespaceInfos = append(namespaceInfos, namespaceInfo)
	}

	return namespaceInfos, nil
}

// GetPodLogs retrieves logs for a specific pod
func (c *Client) GetPodLogs(ctx context.Context, namespace, name string, tailLines *int64) (string, error) {
	opts := &corev1.PodLogOptions{}
	if tailLines != nil {
		opts.TailLines = tailLines
	}

	req := c.clientset.CoreV1().Pods(namespace).GetLogs(name, opts)
	logs, err := req.Stream(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get logs for pod %s in namespace %s: %w", name, namespace, err)
	}
	defer logs.Close()

	buf := make([]byte, 1024*16) // 16KB buffer
	var result []byte
	for {
		n, err := logs.Read(buf)
		if n > 0 {
			result = append(result, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	return string(result), nil
}

// Helper functions to build info structs

func buildPodInfo(pod *corev1.Pod) PodInfo {
	var readyContainers, totalContainers int32
	var restarts int32

	for _, containerStatus := range pod.Status.ContainerStatuses {
		totalContainers++
		if containerStatus.Ready {
			readyContainers++
		}
		restarts += containerStatus.RestartCount
	}

	return PodInfo{
		Name:      pod.Name,
		Namespace: pod.Namespace,
		Status:    string(pod.Status.Phase),
		Ready:     fmt.Sprintf("%d/%d", readyContainers, totalContainers),
		Restarts:  restarts,
		Age:       formatAge(pod.CreationTimestamp.Time),
		Labels:    pod.Labels,
		NodeName:  pod.Spec.NodeName,
		PodIP:     pod.Status.PodIP,
	}
}

func buildServiceInfo(service *corev1.Service) ServiceInfo {
	var ports []string
	for _, port := range service.Spec.Ports {
		portStr := fmt.Sprintf("%d/%s", port.Port, port.Protocol)
		if port.NodePort != 0 {
			portStr += fmt.Sprintf(":%d", port.NodePort)
		}
		ports = append(ports, portStr)
	}

	var externalIPs []string
	for _, ip := range service.Status.LoadBalancer.Ingress {
		if ip.IP != "" {
			externalIPs = append(externalIPs, ip.IP)
		}
		if ip.Hostname != "" {
			externalIPs = append(externalIPs, ip.Hostname)
		}
	}

	return ServiceInfo{
		Name:        service.Name,
		Namespace:   service.Namespace,
		Type:        string(service.Spec.Type),
		ClusterIP:   service.Spec.ClusterIP,
		ExternalIPs: externalIPs,
		Ports:       ports,
		Age:         formatAge(service.CreationTimestamp.Time),
		Labels:      service.Labels,
		Selector:    service.Spec.Selector,
	}
}

func buildDeploymentInfo(deployment *appsv1.Deployment) DeploymentInfo {
	return DeploymentInfo{
		Name:      deployment.Name,
		Namespace: deployment.Namespace,
		Ready:     fmt.Sprintf("%d/%d", deployment.Status.ReadyReplicas, *deployment.Spec.Replicas),
		UpToDate:  deployment.Status.UpdatedReplicas,
		Available: deployment.Status.AvailableReplicas,
		Age:       formatAge(deployment.CreationTimestamp.Time),
		Labels:    deployment.Labels,
		Replicas:  *deployment.Spec.Replicas,
	}
}

func buildNamespaceInfo(namespace *corev1.Namespace) NamespaceInfo {
	return NamespaceInfo{
		Name:   namespace.Name,
		Status: string(namespace.Status.Phase),
		Age:    formatAge(namespace.CreationTimestamp.Time),
		Labels: namespace.Labels,
	}
}

// formatAge formats the age of a Kubernetes resource
func formatAge(creationTime time.Time) string {
	age := time.Since(creationTime)
	
	if age < time.Minute {
		return fmt.Sprintf("%ds", int(age.Seconds()))
	} else if age < time.Hour {
		return fmt.Sprintf("%dm", int(age.Minutes()))
	} else if age < 24*time.Hour {
		return fmt.Sprintf("%dh", int(age.Hours()))
	} else {
		return fmt.Sprintf("%dd", int(age.Hours()/24))
	}
}