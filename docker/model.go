package docker

type Proccess struct {
	UniqueId  string
	ProcessId string
	PPID      string
	C         string
	TTY       string
	TIME      string
	CMD       string
}

type ProcessesResponse struct {
	Titles    []string   `json:"Titles"`
	Processes [][]string `json:"Processes"`
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

type Image struct {
	ID          string      `json:"Id"`
	ParentId    string      `json:"ParentId"`
	RepoTags    []string    `json:"RepoTags"`
	RepoDigests []string    `json:"RepoDigests"`
	Created     int         `json:"Created"`
	Size        int64       `json:"Size"`
	SharedSize  int64       `json:"SharedSize"`
	VirtualSize int64       `json:"VirtualSize"`
	Labels      interface{} `json:"Labels"`
	Containers  int         `json:"Containers"`
}
