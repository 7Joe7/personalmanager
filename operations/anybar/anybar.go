package anybar

import (
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"sort"
	"time"

	"github.com/7joe7/personalmanager/resources"
)

type anybarManager struct {
	binPath string
}

func NewAnybarManager(binPath string) resources.Anybar {
	return &anybarManager{binPath}
}

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
	am.saveActivePorts(t, activePorts)
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
	am.saveActivePorts(t, activePorts)
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
	if err := am.startWithIcon(port, title, icon); err != nil {
		panic(err)
	}
}

func (am *anybarManager) StartNew(port int, title string) {
	defer resources.WaitGroup.Done()
	if _, err := am.start(port, title); err != nil {
		panic(err)
	}
}

func (am *anybarManager) ChangeIcon(port int, colour string) {
	defer resources.WaitGroup.Done()
	if err := am.sendCommand(port, colour); err != nil {
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
		if err := am.sendCommand(port, resources.ANY_CMD_QUIT); err != nil {
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
	addr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:3500")
	if err != nil {
		return false
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return false
	}
	defer conn.Close()
	var received bool
	go func() {
		toReceive := make([]byte, 5)
		_, _, err := conn.ReadFromUDP(toReceive)
		if err != nil {
			return
		}
		received = true
	}()
	err = am.sendCommand(port, "ping")
	if err != nil {
		return false
	}
	time.Sleep(resources.ANY_SLEEP_TIME * time.Millisecond)
	return received
}

func (am *anybarManager) startWithIcon(port int, title, icon string) error {
	_, err := am.start(port, title)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Millisecond * resources.ANY_SLEEP_TIME)
	if err = am.sendCommand(port, icon); err != nil {
		panic(err)
	}
	return nil
}

func (am *anybarManager) saveActivePorts(t resources.Transaction, activePorts resources.ActivePorts) {
	activePortsB, err := json.Marshal(activePorts)
	if err != nil {
		panic(err)
	}
	err = t.SetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ANYBAR_ACTIVE_PORTS, activePortsB)
	if err != nil {
		panic(err)
	}
}

func (am *anybarManager) start(port int, title string) (string, error) {
	cmd := exec.Command("open", "-n", fmt.Sprintf("%s/AnyBar.app", am.binPath))
	cmd.Env = []string{fmt.Sprintf("ANYBAR_PORT=%d", port), fmt.Sprintf("ANYBAR_TITLE=%s", title)}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (am *anybarManager) sendCommand(port int, command string) error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write([]byte(command))
	return err
}
