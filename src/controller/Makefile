SKAFFOLD_DEVICEORBIT_REPO=deviceorbit:5050

build:
	skaffold build

check_context:
	@if [ "$$(docker context show)" != "deviceorbit" ]; then \
		echo "Switching Docker context to 'deviceorbit'"; \
		docker context use deviceorbit; \
	fi
	@if [ "$$(kubectx --current)" != "k3d-deviceorbit" ]; then \
		echo "Switching Kubernetes context to 'k3d-deviceorbit'"; \
		kubectx k3d-deviceorbit; \
	fi

skaffold_run:
	skaffold run --profile=deviceorbit --default-repo=${SKAFFOLD_DEVICEORBIT_REPO}

skaffold_clean:
	skaffold delete --profile=deviceorbit --default-repo=${SKAFFOLD_DEVICEORBIT_REPO}

run: check_context skaffold_run

clean: check_context skaffold_clean