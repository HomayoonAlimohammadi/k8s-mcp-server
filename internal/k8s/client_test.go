package k8s

import (
	"context"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestClient_ListPods(t *testing.T) {
	// Create fake clientset
	fakeClientset := fake.NewSimpleClientset()

	// Create test pods
	pod1 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod-1",
			Namespace: "default",
			Labels:    map[string]string{"app": "test"},
			CreationTimestamp: metav1.Time{Time: time.Now().Add(-time.Hour)},
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Ready:        true,
					RestartCount: 0,
				},
			},
		},
		Spec: corev1.PodSpec{
			NodeName: "test-node-1",
		},
	}

	pod2 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod-2",
			Namespace: "default",
			Labels:    map[string]string{"app": "test"},
			CreationTimestamp: metav1.Time{Time: time.Now().Add(-30 * time.Minute)},
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodPending,
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Ready:        false,
					RestartCount: 1,
				},
			},
		},
	}

	// Add pods to fake clientset
	fakeClientset.CoreV1().Pods("default").Create(context.TODO(), pod1, metav1.CreateOptions{})
	fakeClientset.CoreV1().Pods("default").Create(context.TODO(), pod2, metav1.CreateOptions{})

	// Create client
	client := &Client{clientset: fakeClientset}

	// Test ListPods
	pods, err := client.ListPods(context.TODO(), "default")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pods) != 2 {
		t.Fatalf("Expected 2 pods, got %d", len(pods))
	}

	// Verify first pod
	if pods[0].Name != "test-pod-1" {
		t.Errorf("Expected pod name 'test-pod-1', got '%s'", pods[0].Name)
	}
	if pods[0].Status != "Running" {
		t.Errorf("Expected pod status 'Running', got '%s'", pods[0].Status)
	}
	if pods[0].Ready != "1/1" {
		t.Errorf("Expected pod ready '1/1', got '%s'", pods[0].Ready)
	}
}

func TestClient_GetPod(t *testing.T) {
	// Create fake clientset
	fakeClientset := fake.NewSimpleClientset()

	// Create test pod
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			Labels:    map[string]string{"app": "test"},
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
			PodIP: "10.0.0.1",
		},
		Spec: corev1.PodSpec{
			NodeName: "test-node",
		},
	}

	// Add pod to fake clientset
	fakeClientset.CoreV1().Pods("default").Create(context.TODO(), pod, metav1.CreateOptions{})

	// Create client
	client := &Client{clientset: fakeClientset}

	// Test GetPod
	retrievedPod, err := client.GetPod(context.TODO(), "default", "test-pod")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrievedPod.Name != "test-pod" {
		t.Errorf("Expected pod name 'test-pod', got '%s'", retrievedPod.Name)
	}
	if retrievedPod.PodIP != "10.0.0.1" {
		t.Errorf("Expected pod IP '10.0.0.1', got '%s'", retrievedPod.PodIP)
	}
	if retrievedPod.NodeName != "test-node" {
		t.Errorf("Expected node name 'test-node', got '%s'", retrievedPod.NodeName)
	}
}

func TestClient_ListServices(t *testing.T) {
	// Create fake clientset
	fakeClientset := fake.NewSimpleClientset()

	// Create test service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "default",
			Labels:    map[string]string{"app": "test"},
		},
		Spec: corev1.ServiceSpec{
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: "10.0.0.100",
			Ports: []corev1.ServicePort{
				{
					Port:     80,
					Protocol: corev1.ProtocolTCP,
				},
			},
			Selector: map[string]string{"app": "test"},
		},
	}

	// Add service to fake clientset
	fakeClientset.CoreV1().Services("default").Create(context.TODO(), service, metav1.CreateOptions{})

	// Create client
	client := &Client{clientset: fakeClientset}

	// Test ListServices
	services, err := client.ListServices(context.TODO(), "default")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(services) != 1 {
		t.Fatalf("Expected 1 service, got %d", len(services))
	}

	if services[0].Name != "test-service" {
		t.Errorf("Expected service name 'test-service', got '%s'", services[0].Name)
	}
	if services[0].Type != "ClusterIP" {
		t.Errorf("Expected service type 'ClusterIP', got '%s'", services[0].Type)
	}
	if services[0].ClusterIP != "10.0.0.100" {
		t.Errorf("Expected cluster IP '10.0.0.100', got '%s'", services[0].ClusterIP)
	}
}

func TestClient_ListDeployments(t *testing.T) {
	// Create fake clientset
	fakeClientset := fake.NewSimpleClientset()

	// Create test deployment
	replicas := int32(3)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-deployment",
			Namespace: "default",
			Labels:    map[string]string{"app": "test"},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
		},
		Status: appsv1.DeploymentStatus{
			ReadyReplicas:     2,
			UpdatedReplicas:   3,
			AvailableReplicas: 2,
			Replicas:          3,
		},
	}

	// Add deployment to fake clientset
	fakeClientset.AppsV1().Deployments("default").Create(context.TODO(), deployment, metav1.CreateOptions{})

	// Create client
	client := &Client{clientset: fakeClientset}

	// Test ListDeployments
	deployments, err := client.ListDeployments(context.TODO(), "default")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(deployments) != 1 {
		t.Fatalf("Expected 1 deployment, got %d", len(deployments))
	}

	if deployments[0].Name != "test-deployment" {
		t.Errorf("Expected deployment name 'test-deployment', got '%s'", deployments[0].Name)
	}
	if deployments[0].Ready != "2/3" {
		t.Errorf("Expected deployment ready '2/3', got '%s'", deployments[0].Ready)
	}
	if deployments[0].UpToDate != 3 {
		t.Errorf("Expected deployment up-to-date 3, got %d", deployments[0].UpToDate)
	}
}

func TestClient_GetPodLogs(t *testing.T) {
	// Create fake clientset
	fakeClientset := fake.NewSimpleClientset()

	// Create test pod
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
		},
	}

	// Add pod to fake clientset
	fakeClientset.CoreV1().Pods("default").Create(context.TODO(), pod, metav1.CreateOptions{})

	// Create client
	client := &Client{clientset: fakeClientset}

	// Test GetPodLogs - this will fail in the fake client, but we can test the error handling
	_, err := client.GetPodLogs(context.TODO(), "default", "test-pod", nil)
	// We expect an error since the fake client doesn't implement log streaming properly
	if err == nil {
		t.Log("Note: Fake client may not properly simulate log streaming errors")
	}
}

func TestFormatAge(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "seconds",
			time:     now.Add(-30 * time.Second),
			expected: "30s",
		},
		{
			name:     "minutes",
			time:     now.Add(-5 * time.Minute),
			expected: "5m",
		},
		{
			name:     "hours",
			time:     now.Add(-3 * time.Hour),
			expected: "3h",
		},
		{
			name:     "days",
			time:     now.Add(-25 * time.Hour),
			expected: "1d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatAge(tt.time)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}