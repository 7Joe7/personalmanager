package anybar

import (
	"encoding/json"
	"sort"

	"github.com/7joe7/personalmanager/resources"
)

func (am *anybarManager) RemoveAndQuit(bucketName []byte, id string, t resources.Transaction) {
	var port int
	activePorts := am.GetActivePorts(t)
	for i := 0; i < len(activePorts); i++ {
		if string(bucketName) == string(activePorts[i].BucketName) && activePorts[i].Id == id {
			port = activePorts[i].Port
			activePorts = append(activePorts[:i], activePorts[i+1:]...)
			break
		}
	}
	saveActivePorts(t, activePorts)
	resources.WaitGroup.Add(1)
	go am.Quit(port)
}

func (am *anybarManager) AddToActivePorts(title, icon string, bucketName []byte, id string, t resources.Transaction) {
	activePorts := am.GetActivePorts(t)
	port := am.GetNewPort(activePorts)
	activePorts = append(activePorts, &resources.ActivePort{
		Port:       port,
		Name:       title,
		Colour:     icon,
		BucketName: bucketName,
		Id:         id})
	sort.Sort(activePorts)
	saveActivePorts(t, activePorts)
	resources.WaitGroup.Add(1)
	go am.StartWithIcon(port, title, icon)
}

func (am *anybarManager) EnsureActivePorts(activePorts resources.ActivePorts) {
	defer resources.WaitGroup.Done()
	for i := 0; i < len(activePorts); i++ {
		if !am.Ping(activePorts[i].Port) {
			resources.WaitGroup.Add(1)
			am.StartWithIcon(activePorts[i].Port, activePorts[i].Name, activePorts[i].Colour)
		}
	}
}

func (am *anybarManager) StartWithIcon(port int, title, icon string) {
	defer resources.WaitGroup.Done()
	if err := startWithIcon(port, title, icon); err != nil {
		panic(err)
	}
}

func (am *anybarManager) StartNew(port int, title string) {
	defer resources.WaitGroup.Done()
	if _, err := start(port, title); err != nil {
		panic(err)
	}
}

func (am *anybarManager) ChangeIcon(port int, colour string) {
	defer resources.WaitGroup.Done()
	if err := sendCommand(port, colour); err != nil {
		panic(err)
	}
}

func (am *anybarManager) GetNewPort(activePorts []*resources.ActivePort) int {
	for i := 0; i < len(activePorts); i++ {
		if activePorts[i].Port != resources.ANY_PORTS_RANGE_BASE+i {
			return resources.ANY_PORTS_RANGE_BASE + i
		}
	}
	return len(activePorts) + resources.ANY_PORTS_RANGE_BASE
}

func (am *anybarManager) Quit(port int) {
	defer resources.WaitGroup.Done()
	if port != 0 {
		if err := sendCommand(port, resources.ANY_CMD_QUIT); err != nil {
			panic(err)
		}
	}
}

func (am *anybarManager) GetActivePorts(t resources.Transaction) resources.ActivePorts {
	var activePorts resources.ActivePorts
	activePortsB := t.GetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ANYBAR_ACTIVE_PORTS)
	err := json.Unmarshal(activePortsB, &activePorts)
	if err != nil {
		panic(err)
	}
	return activePorts
}

func (am *anybarManager) Ping(port int) bool {
	return ping(port)
}
