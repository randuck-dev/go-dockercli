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

func (dc *DockerClient) GetContainers() ([]Container, error) {

	sd := dc.getHttpClient()
	resp, err := sd.Get("http://localhost/containers/json")
	if err != nil {
		fmt.Println("Failed to parse response")
		return []Container{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	var containers []Container
	err = json.Unmarshal(body, &containers)

	if err != nil {
		fmt.Println("Error: ", err)
		return []Container{}, err
	}

	return containers, nil
}

type Port struct {
	IP          string `json:"IP"`
	PrivatePort int    `json:"PrivatePort"`
	PublicPort  int    `json:"PublicPort"`
	Type        string `json:"Type"`
}

type Container struct {
	ID         string   `json:"Id"`
	Names      []string `json:"Names"`
	Image      string   `json:"Image"`
	ImageID    string   `json:"ImageID"`
	Command    string   `json:"Command"`
	Created    int64    `json:"Created"`
	Ports      []Port   `json:"Ports"`
	State      string   `json:"State"`
	Status     string   `json:"Status"`
	HostConfig struct {
		NetworkMode string `json:"NetworkMode"`
	} `json:"HostConfig"`
}
