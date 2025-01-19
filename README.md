# DeviceOrbit Monorepo
Streamlined USB Mobile Device Management in Kubernetes

## Overview
**DeviceOrbit** is a monorepo that integrates USB-connected Android and iOS devices into Kubernetes environments. It combines a Device Plugin, a pod management controller, and a local K3s cluster setup to enable secure, containerized access to mobile devices for development and testing. Built with Skaffold, Helm, and K3s, it’s designed for simplicity and flexibility.

### Purpose
DeviceOrbit simplifies working with physical mobile devices in Kubernetes by:
- Exposing USB devices as node resources
- Managing device-specific pods programmatically
- Providing a lightweight local testing environment

## Components
The monorepo includes three key components, each with its own subdirectory:

1. **`src/plugin`**: A Kubernetes Device Plugin that detects USB mobile devices (Android/iOS) and exposes them as resources (e.g., `com.kvs.deviceorbit/<serial>`).
2. **`src/controller`**: A controller that creates and deletes pods for these devices via gRPC and HTTP APIs.
3. **`cluster/`**: A K3d configuration for a local K3s cluster with USB support and a registry (`deviceorbit:5050`).

For detailed functionality, setup, and usage, refer to the `README.md` in each subdirectory.