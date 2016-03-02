package bindings

type PingResponse struct {
	Pong string `json:"pong"`
}

type StatusResponse struct {
	DockerHostStatus string `json:"docker_host_status"`
	Status           string `json:"status"`
}

type VersionResponse struct {
	DockerHostVersion string `json:"docker_host_version"`
	Version           string `json:"version"`
}

type PushEventResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
