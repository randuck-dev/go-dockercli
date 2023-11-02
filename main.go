package main

import (
	docker "go-dockercli/internal"
	"log/slog"
)

func main() {
	slog.Info("Starting to listen to docker socket. Will communicate with the HTTP Client")

	docker_socket := "/var/run/docker.sock"

	dc := docker.NewDockerClient(docker_socket)

	containers, err := dc.GetContainers()

	if err != nil {
		slog.Info("error while getting containers", err)
	}

	for _, v := range containers {
		slog.Info("Name", "image", v.Image)
	}

	running_processes, err := dc.GetRunningProcesses(containers[0].ID)

	if err != nil {
		panic(err)
	}
	slog.Info("Found running processes for process", "id", containers[0].ID, "processes", running_processes)
}
