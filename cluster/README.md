# K3d Cluster Setup for Mobile Device Testing
Quickly Set Up a Local K3s Cluster for USB Mobile Device Projects

## Overview
This guide provides instructions to set up a local K3s cluster using K3d, tailored for testing the **Mobile Device Controller** and **Smart Device Manager** projects. The cluster, named `k3d-deviceorbit`, includes a local registry (`deviceorbit:5050`), USB device access, and port mappings for HTTP/gRPC APIs. It’s ideal for developing and testing workflows involving Android and iOS devices in Kubernetes.

### Purpose
- Create a lightweight K3s cluster for local development
- Enable USB device passthrough for mobile device testing
- Integrate with Skaffold and Helm workflows from related projects

## Prerequisites
To set up and use this K3d cluster, ensure the following requirements are met:

- **Operating System**: A Linux distribution (e.g., Ubuntu) with `sudo` privileges
- **Required Tools**:
  - [Docker](https://www.docker.com/) (`docker.io`) for container runtime
  - [K3d](https://k3d.io) (`k3d`) for managing the K3s cluster
  - `kubectx` and `kubectl` for Kubernetes context and cluster management
  - `make` for executing Makefile targets
- **Hardware**: A machine equipped with USB ports and connected Android/iOS devices
- **Dependencies**: `libusb` installed on the host for USB device support
- **Docker Registry Configuration**: 
  The local registry operates on `deviceorbit:5050` over HTTP. Configure Docker to recognize it as an insecure registry:
  - Edit (or create) `/etc/docker/daemon.json`:
    ```json
    {
      "insecure-registries": ["deviceorbit:5050"]
    }
    ```

## Cluster Configuration
The cluster is defined in `k3d.yaml`:
- **Name**: `k3d-deviceorbit`
- **API**: Exposed at `deviceorbit:6445` (mapped to `0.0.0.0:6445`)
- **Registry**: 
  - Name: `deviceorbit`
  - Host: `0.0.0.0:5050`
  - Mirror: `deviceorbit.local:5050` (endpoint: `http://deviceorbit:5000`)
- **Volumes**: 
  - `/dev:/dev` mapped to the server node for USB access
- **Ports**:
  - `8080-8090:8080-8090` (HTTP range)
  - `30050-30100:30050-30100` (NodePort range)
- **K3s Options**: 
  - Allows unsafe sysctls: `net.ipv6.conf.all.disable_ipv6`, `net.ipv4.ip_forward`

## Setup Instructions
### Makefile Targets
The repository includes a Makefile with a single target for initialization:
- **`init`**: Runs `./setup.sh` with `IP_ADDRESS` and `HOSTNAME` arguments
  - Default: `IP_ADDRESS=127.0.0.2`, `HOSTNAME=deviceorbit`
