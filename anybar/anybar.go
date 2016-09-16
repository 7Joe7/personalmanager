package anybar

import (
	"os/exec"
	"fmt"
	"net"
	"time"
	"sort"
	"encoding/json"

	"github.com/7joe7/personalmanager/resources"
)

func StartNewPort(title, icon string, id []byte, t resources.Transaction) {
	var activePorts resources.ActivePorts
	activePortsB := t.GetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ANYBAR_ACTIVE_HABITS_PORTS)
	err := json.Unmarshal(activePortsB, activePorts)
	if err != nil {
		panic(err)
	}

	port := GetNewPort(activePorts)
	activePorts = append(activePorts, &resources.ActivePort{Port:port,BucketName:resources.DB_DEFAULT_HABITS_BUCKET_NAME,Id:id})
	sort.Sort(activePorts)

	activePortsB, err = json.Marshal(activePorts)
	if err != nil {
		panic(err)
	}
	err = t.SetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ANYBAR_ACTIVE_HABITS_PORTS, activePortsB)
	if err != nil {
		panic(err)
	}
}

func StartWithIcon(port int, title, icon string) {
	_, err := start(port, title)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Millisecond * 100)
	if err = sendCommand(port, icon); err != nil {
		panic(err)
	}
	resources.WaitGroup.Done()
}

func StartNew(port int, title string) {
	if _, err := start(port, title); err != nil {
		panic(err)
	}
	resources.WaitGroup.Done()
}

func ChangeIcon(port int, colour string) {
	if err := sendCommand(port, colour); err != nil {
		panic(err)
	}
	resources.WaitGroup.Done()
}

func GetNewPort(activePorts []*resources.ActivePort) int {
	for i := 0; i < len(activePorts); i++ {
		if activePorts[i].Port != resources.ANY_PORTS_RANGE_BASE + i {
			return resources.ANY_PORTS_RANGE_BASE + i
		}
	}
	return len(activePorts) + resources.ANY_PORTS_RANGE_BASE
}

func Quit(port int) {
	if err := sendCommand(port, resources.ANY_CMD_QUIT); err != nil {
		panic(err)
	}
	resources.WaitGroup.Done()
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
