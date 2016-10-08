package anybar

import (
	"github.com/7joe7/personalmanager/resources"
)

var (
	am resources.AnybarManager
)

type anybarManager struct{}

func NewAnybarManager() resources.AnybarManager {
	return &anybarManager{}
}

func Start(anybarManager resources.AnybarManager) {
	am = anybarManager
}

func RemoveAndQuit(bucketName []byte, id string, t resources.Transaction) {
	am.RemoveAndQuit(bucketName, id, t)
}

func AddToActivePorts(title, icon string, bucketName []byte, id string, t resources.Transaction) {
	am.AddToActivePorts(title, icon, bucketName, id, t)
}

func EnsureActivePorts(activePorts resources.ActivePorts) {
	am.EnsureActivePorts(activePorts)
}

func StartWithIcon(port int, title, icon string) {
	am.StartWithIcon(port, title, icon)
}

func StartNew(port int, title string) {
	am.StartNew(port, title)
}

func ChangeIcon(port int, colour string) {
	am.ChangeIcon(port, colour)
}

func GetNewPort(activePorts []*resources.ActivePort) int {
	return am.GetNewPort(activePorts)
}

func Quit(port int) {
	am.Quit(port)
}

func GetActivePorts(t resources.Transaction) resources.ActivePorts {
	return am.GetActivePorts(t)
}

func Ping(port int) bool {
	return am.Ping(port)
}
