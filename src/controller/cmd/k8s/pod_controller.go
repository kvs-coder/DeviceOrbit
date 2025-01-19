package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodController struct {
	ctx         context.Context
	kubeClient  *kubernetes.Clientset
	namespace   string
	podProvider PodProvider
}

func NewPodController(ctx context.Context, kubeClient *kubernetes.Clientset, namespace string, podProvider PodProvider) *PodController {
	return &PodController{ctx, kubeClient, namespace, podProvider}
}

func (pc *PodController) ValidateRole() error {
	_, err := pc.kubeClient.CoreV1().Pods(pc.namespace).List(pc.ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to validate k8s pod: %s", err)
	}

	return nil
}

func (pc *PodController) CreatePod(deviceSerial string, platform int) error {
	pod, err := pc.podProvider.DevicePod(deviceSerial, platform)
	if err != nil {
		return fmt.Errorf("failed to create pod config for device serial: %s, cause: %s", deviceSerial, err)
	}

	_, err = pc.kubeClient.CoreV1().Pods(pc.namespace).Create(pc.ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create k8s pod for device serial: %s, cause: %s", deviceSerial, err)
	}

	return nil
}

func (pc *PodController) DeletePod(deviceSerial string) error {
	podName := podName(deviceSerial)
	var gracePeriod int64 = 5
	err := pc.kubeClient.CoreV1().Pods(pc.namespace).Delete(pc.ctx, podName, metav1.DeleteOptions{
		GracePeriodSeconds: &gracePeriod,
	})
	if err != nil {
		return fmt.Errorf("failed to delete k8s pod for device serial: %s, cause: %s", deviceSerial, err)
	}
	return nil
}
