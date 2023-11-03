package main

import (
	docker "go-dockercli/internal"
	"log/slog"
	"sync"
)

func main() {
	docker_socket := "/var/run/docker.sock"

	wait := sync.WaitGroup{}
	wait.Add(1)

	// go docker.Raw_http_parsing_docker_socket(docker_socket, &wait)

	go docker_http_builtin(docker_socket, &wait)

	wait.Wait()
}

func docker_http_builtin(docker_socket string, wg *sync.WaitGroup) {
	dc := docker.NewDockerClient(docker_socket)

	containers, err := dc.GetContainers()

	if err != nil {
		slog.Info("error while fetching containers", "err", err)
	}

	for _, v := range containers {
		slog.Info("Name", "image", v.Image)
	}

	running_processes, err := dc.GetRunningProcesses(containers[0].ID)

	if err != nil {
		panic(err)
	}
	slog.Info("Found running processes for process", "id", containers[0].ID, "processes", running_processes)

	images, err := dc.ListImages()

	for _, v := range images {
		slog.Info("Image", "image", v)
	}

	wg.Done()
}
