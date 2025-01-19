package docker

type ContainerConfig struct {
	Image   string
	Tag     string
	Command []string
	Args    []string
}

func NewContainerConfig(image, tag string, command, args []string) ContainerConfig {
	return ContainerConfig{Image: image, Tag: tag, Command: command, Args: args}
}
