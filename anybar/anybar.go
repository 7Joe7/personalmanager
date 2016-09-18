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

func RemoveAndQuit(id string, t resources.Transaction) {
	var port int
	activePorts := GetActivePorts(t)
	for i := 0; i < len(activePorts); i++ {
		if string(activePorts[i].Id) == id {
			port = activePorts[i].Port
			activePorts = append(activePorts[:i], activePorts[i+1:]...)
			break
		}
	}
	saveActivePorts(t, activePorts)
	resources.WaitGroup.Add(1)
	go Quit(port)
}

func AddToActivePorts(title, icon string, id string, t resources.Transaction) {
	activePorts := GetActivePorts(t)
	port := GetNewPort(activePorts)
	activePorts = append(activePorts, &resources.ActivePort{
		Port:       port,
		Name:       title,
		Colour:     icon,
		BucketName: resources.DB_DEFAULT_HABITS_BUCKET_NAME,
		Id:         id})
	sort.Sort(activePorts)
	saveActivePorts(t, activePorts)
	resources.WaitGroup.Add(1)
	go StartWithIcon(port, title, icon)
}

func EnsureActivePorts(activePorts resources.ActivePorts) {
	defer resources.WaitGroup.Done()
	for i := 0; i < len(activePorts); i++ {
		resources.WaitGroup.Add(1)
		StartWithIcon(activePorts[i].Port, activePorts[i].Name, activePorts[i].Colour)
	}
}

func StartWithIcon(port int, title, icon string) {
	defer resources.WaitGroup.Done()
	if err := startWithIcon(port, title, icon); err != nil {
		panic(err)
	}
}

func startWithIcon(port int, title, icon string) error {
	_, err := start(port, title)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Millisecond * 100)
	if err = sendCommand(port, icon); err != nil {
		panic(err)
	}
	return nil
}

func StartNew(port int, title string) {
	defer resources.WaitGroup.Done()
	if _, err := start(port, title); err != nil {
		panic(err)
	}
}

func ChangeIcon(port int, colour string) {
	defer resources.WaitGroup.Done()
	if err := sendCommand(port, colour); err != nil {
		panic(err)
	}
}

func GetNewPort(activePorts []*resources.ActivePort) int {
	for i := 0; i < len(activePorts); i++ {
		if activePorts[i].Port != resources.ANY_PORTS_RANGE_BASE+i {
			return resources.ANY_PORTS_RANGE_BASE + i
		}
	}
	return len(activePorts) + resources.ANY_PORTS_RANGE_BASE
}

func Quit(port int) {
	defer resources.WaitGroup.Done()
	if port != 0 {
		if err := sendCommand(port, resources.ANY_CMD_QUIT); err != nil {
			panic(err)
		}
	}
}

func GetActivePorts(t resources.Transaction) resources.ActivePorts {
	var activePorts resources.ActivePorts
	activePortsB := t.GetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ANYBAR_ACTIVE_PORTS)
	err := json.Unmarshal(activePortsB, &activePorts)
	if err != nil {
		panic(err)
	}
	return activePorts
}

func Ping(port int) bool {
	return ping(port)
}

func saveActivePorts(t resources.Transaction, activePorts resources.ActivePorts) {
	activePortsB, err := json.Marshal(activePorts)
	if err != nil {
		panic(err)
	}
	err = t.SetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ANYBAR_ACTIVE_PORTS, activePortsB)
	if err != nil {
		panic(err)
	}
}

func start(port int, title string) (string, error) {
	cmd := exec.Command("open", "-n", "./AnyBar.app")
	cmd.Env = []string{fmt.Sprintf("ANYBAR_PORT=%d", port), fmt.Sprintf("ANYBAR_TITLE=%s", title)}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func sendCommand(port int, command string) error {
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

func ping(port int) bool {
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
	go func () {
		toReceive := make([]byte, 5)
		_, _, err := conn.ReadFromUDP(toReceive)
		if err != nil {
			return
		}
		received = true
	}()
	err = sendCommand(port, "ping")
	if err != nil {
		return false
	}
	time.Sleep(100 * time.Millisecond)
	return received
}
