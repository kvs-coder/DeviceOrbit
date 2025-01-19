# Smart Device Manager for Kubernetes
Expose USB Mobile Devices (Android/iOS) to Kubernetes Containers

## Overview
The **Smart Device Manager** is a Kubernetes Device Plugin that enables pods to access USB-connected Android and iOS mobile devices on cluster nodes. It leverages the Kubernetes [Device Plugins](https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/device-plugins/) framework to expose these devices as node resources (e.g., `com.kvs.deviceorbit/<serial>`), ensuring secure isolation. Built with Skaffold and deployed via Helm, it pairs seamlessly with the Mobile Device Controller.

### Motivation
Direct access to USB mobile devices (smartphones, tablets) in Kubernetes typically requires privileged containers, breaking isolation. This plugin:
- Exposes Android and iOS devices as node resources
- Abstracts hardware from pod internals
- Delegates allocation to Kubernetes

## Features
- **USB Mobile Device Detection**: Identifies Android/iOS devices using [GoUSB](https://github.com/google/gousb)
- **Resource Exposure**: Registers devices by serial number (e.g., `com.kvs.deviceorbit/FA82T1A01145`)
- **Skaffold Build**: Custom-tagged images pushed to `deviceorbit:5050`
- **Helm Deployment**: Deploys to `k3d-deviceorbit` cluster

## Architecture
1. **Device Detection**: Scans node USB ports with `libusb` to detect Android/iOS devices
2. **Resource Registration**: Advertises devices to Kubelet as node resources via the Device Plugin API
3. **Device Access**: 
   - Kubelet queries the plugin for device access
   - Mounts devices to pod `/dev` using OCI `--device` and updates cgroups
4. **Build & Deploy**: Skaffold builds the image, Helm deploys it

## Usage
### Prerequisites
- Kubernetes cluster (e.g., `k3d-deviceorbit` via K3s)
- Nodes with USB ports and Android/iOS devices
- `libusb` support on nodes
- Skaffold, Helm, Docker, `kubectx` installed
- Docker context `deviceorbit` pointing to `deviceorbit:5050`

### Makefile Targets
The repository includes a Makefile with the following targets:
- **`build`**: Runs `skaffold build` to create the plugin image
- **`check_context`**: Ensures Docker context is `deviceorbit` and Kubernetes context is `k3d-deviceorbit`
- **`skaffold_run`**: Executes `skaffold run --profile=deviceorbit --default-repo=$SKAFFOLD_DEVICEORBIT_REPO`
- **`skaffold_clean`**: Executes `skaffold delete --profile=deviceorbit --default-repo=$SKAFFOLD_DEVICEORBIT_REPO`
- **`run`**: Combines `check_context` and `skaffold_run`
- **`clean`**: Combines `check_context` and `skaffold_clean`

### Build with Skaffold
Build the plugin image:
```bash
make build