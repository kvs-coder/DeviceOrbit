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
	$(MAKE) -C src/controller skaffold_run
	$(MAKE) -C src/plugin skaffold_run
	wait

skaffold_clean:
	$(MAKE)	-C src/plugin skaffold_clean
	$(MAKE) -C src/controller skaffold_clean
	wait 

run: check_context skaffold_run

clean: check_context skaffold_clean