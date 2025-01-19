# Mobile Device Controller for Kubernetes
Manage USB Mobile Devices (Android/iOS) in Kubernetes with gRPC and HTTP APIs

## Overview
The **Mobile Device Controller** is a Kubernetes controller designed to manage USB-connected Android and iOS devices within a cluster. It creates and deletes pods linked to these devices, leveraging Kubernetes Device Plugins to expose them as resources (e.g., `com.kvs.deviceorbit/<serial>`). With gRPC and HTTP APIs, Skaffold-driven builds, and Helm-based deployment, it’s tailored for automated mobile device workflows.

### Motivation
Accessing physical Android and iOS devices in Kubernetes is critical for mobile app development and testing. This controller:
- Provides secure USB device access without privileged containers
- Ties pods to device serial numbers
- Offers flexible APIs for programmatic control

## Features
- **USB Mobile Device Support**: Targets Android (platform 2) and iOS (platform 1)
- **Dual APIs**:
  - gRPC: Type-safe pod management
  - HTTP: REST-like pod control
- **Pod Lifecycle**: Creates/deletes pods dynamically
- **Skaffold Build**: Custom-tagged images pushed to `deviceorbit:5050`
- **Helm Deployment**: Deploys to `k3d-deviceorbit` cluster

## Architecture
1. **Kubernetes Integration**:
   - Uses `k8s.io/client-go` for API access
   - Depends on a Device Plugin for `com.kvs.deviceorbit/<serial>` resources
2. **Pod Management**:
   - Configures pods with device-specific mounts (`/dev` for Android, `/var/run` for iOS)
   - Deletes pods with a 5-second grace period
3. **Servers**:
   - gRPC: `CreateDevice`/`DeleteDevice`
   - HTTP: `/pods` (POST/DELETE), `/health`
4. **Build & Deploy**:
   - Skaffold builds and tags images
   - Helm deploys to `k3d-deviceorbit`

## Usage
### Prerequisites
- Kubernetes cluster (e.g., k3d with `k3d-deviceorbit` context)
- USB-capable nodes with Android/iOS devices
- Device Plugin exposing `com.kvs.deviceorbit/<serial>` resources
- `libusb` on nodes
- Tools: Skaffold, Helm, Docker, `kubectx`
- Docker context `deviceorbit` pointing to `deviceorbit:5050`

### Makefile Targets
The repository includes a Makefile with the following targets:
- **`build`**: Runs `skaffold build` to create the controller image
- **`check_context`**: Ensures Docker context is `deviceorbit` and Kubernetes context is `k3d-deviceorbit`
- **`skaffold_run`**: Executes `skaffold run --profile=deviceorbit --default-repo=$SKAFFOLD_DEVICEORBIT_REPO`
- **`skaffold_clean`**: Executes `skaffold delete --profile=deviceorbit --default-repo=$SKAFFOLD_DEVICEORBIT_REPO`
- **`run`**: Combines `check_context` and `skaffold_run`
- **`clean`**: Combines `check_context` and `skaffold_clean`

### Build with Skaffold
Build the controller image:
```bash
make build