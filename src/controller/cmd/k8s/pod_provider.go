package k8s

import (
	"fmt"
	"strings"

	"com.kvs/deviceorbit/controller/cmd/docker"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type platform = int

const (
	unknown platform = iota
	iOS
	android
)

type PodProvider interface {
	DevicePod(serial string, platform platform) (*v1.Pod, error)
}

type podProvider struct {
	containerConfig docker.ContainerConfig
}

func DefaultPodProvider(containerConfig docker.ContainerConfig) PodProvider {
	return &podProvider{containerConfig}
}

func (p *podProvider) DevicePod(serial string, platform platform) (*v1.Pod, error) {
	name := podName(serial)
	containerImage := containerImage(p.containerConfig.Image, p.containerConfig.Tag)
	resourceName := resourceName(serial)
	resourceQty, err := resource.ParseQuantity("1")
	if err != nil {
		return nil, fmt.Errorf("failed to parse quantity: %s", err)
	}

	switch platform {
	case iOS:
		return &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: name},
			Spec: v1.PodSpec{
				RestartPolicy: v1.RestartPolicyNever,
				Containers: []v1.Container{
					{
						Name:    name,
						Image:   containerImage,
						Command: p.containerConfig.Command,
						Args:    p.containerConfig.Args,
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								v1.ResourceName(resourceName): resourceQty,
							},
							Limits: v1.ResourceList{
								v1.ResourceName(resourceName): resourceQty,
							},
						},
						VolumeMounts: []v1.VolumeMount{
							{Name: "usbmuxd", MountPath: "/var/run"},
							{Name: "dev-net-mount", MountPath: "/dev/net"},
						},
					},
				},
				Volumes: []v1.Volume{
					{
						Name: "dev-net-mount",
						VolumeSource: v1.VolumeSource{
							HostPath: &v1.HostPathVolumeSource{Path: "/dev/net"},
						},
					},
					{
						Name: "usbmuxd",
						VolumeSource: v1.VolumeSource{
							HostPath: &v1.HostPathVolumeSource{Path: "/var/run"},
						},
					},
				},
			},
		}, nil
	case android:
		return &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: name},
			Spec: v1.PodSpec{
				RestartPolicy: v1.RestartPolicyNever,
				Containers: []v1.Container{
					{
						Name:    name,
						Image:   containerImage,
						Command: p.containerConfig.Command,
						Args:    p.containerConfig.Args,
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								v1.ResourceName(resourceName): resourceQty,
							},
							Limits: v1.ResourceList{
								v1.ResourceName(resourceName): resourceQty,
							},
						},
						VolumeMounts: []v1.VolumeMount{
							{Name: "dev-mount", MountPath: "/dev"},
						},
					},
				},
				Volumes: []v1.Volume{
					{
						Name: "dev-mount",
						VolumeSource: v1.VolumeSource{
							HostPath: &v1.HostPathVolumeSource{Path: "/dev"},
						},
					},
				},
			},
		}, nil
	}

	return nil, fmt.Errorf("unknown platform: %d, cannot create pod config for serial: %s", platform, serial)
}

func podName(deviceSerial string) string {
	return fmt.Sprintf("mobile-device-%s", strings.ToLower(deviceSerial))
}

func resourceName(deviceSerial string) string {
	return fmt.Sprintf("com.kvs.deviceorbit/%s", deviceSerial)
}

func containerImage(containerImage, containerTag string) string {
	return fmt.Sprintf("%s:%s", containerImage, containerTag)
}
