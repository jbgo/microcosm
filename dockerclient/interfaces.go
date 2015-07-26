package dockerclient

type DockerClient interface {
	InspectContainer(id string) (*Container, error)
	AddEventListener(listener chan<- *Event) error
}

type Container struct {
	ID     string
	Labels map[string]string
}

type Event struct {
	ContainerID string
	Status      string
	ImageID     string
	Timestamp   int64
}
