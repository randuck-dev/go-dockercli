package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
)

type DockerClient struct {
	socket  string
	request string
}

func NewDockerClient(socket string) DockerClient {
	return DockerClient{
		socket: socket,
	}
}

func (dc *DockerClient) getHttpClient() http.Client {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", dc.socket)
			},
		},
	}
	return httpc
}

func (dc *DockerClient) Get(path string) ([]byte, error) {
	client := dc.getHttpClient()

	resp, err := client.Get("http://localhost" + path)

	if err != nil {
		return make([]byte, 0), err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	return body, nil
}

func (dc *DockerClient) GetContainers() ([]Container, error) {

	resp, err := dc.Get("/containers/json")

	var containers []Container
	err = json.Unmarshal(resp, &containers)

	if err != nil {
		fmt.Println("Error: ", err)
		return []Container{}, err
	}

	return containers, nil
}

func (dc *DockerClient) GetRunningProcesses(id string) ([]Proccess, error) {

	resp, err := dc.Get("/containers/" + id + "/top")

	if err != nil {
		return make([]Proccess, 0), err
	}

	fmt.Println(string(resp))

	var processesResponse ProcessesResponse

	err = json.Unmarshal(resp, &processesResponse)

	if err != nil {
		return make([]Proccess, 0), err
	}

	processes := make([]Proccess, len(processesResponse.Processes))

	for _, v := range processesResponse.Processes {
		proces := Proccess{
			UniqueId:  v[0],
			ProcessId: v[1],
			PPID:      v[2],
			C:         v[3],
			TTY:       v[4],
			TIME:      v[5],
			CMD:       v[6],
		}

		processes = append(processes, proces)
	}

	return processes, nil
}
