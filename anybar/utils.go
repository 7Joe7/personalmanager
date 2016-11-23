package anybar

import (
	"time"
	"net"
	"os/exec"
	"fmt"
	"encoding/json"

	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/utils"
)

func startWithIcon(port int, title, icon string) error {
	_, err := start(port, title)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Millisecond * resources.ANY_SLEEP_TIME)
	if err = sendCommand(port, icon); err != nil {
		panic(err)
	}
	return nil
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
	cmd := exec.Command("open", "-n", fmt.Sprintf("%s/AnyBar.app", utils.GetRunningBinaryPath()))
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
	time.Sleep(resources.ANY_SLEEP_TIME * time.Millisecond)
	return received
}
