package anybar

import (
	"github.com/7joe7/personalmanager/resources"
)

func NewAnybarManagerMock() resources.AnybarManager {
	return &anybarManagerMock{}
}

type anybarManagerMock struct{
	PingResult bool
	ActivePortsResult resources.ActivePorts
	NewPortResult int
}

func (am *anybarManagerMock) RemoveAndQuit(id string, t resources.Transaction) {

}

func (am *anybarManagerMock) AddToActivePorts(title, icon string, id string, t resources.Transaction) {

}

func (am *anybarManagerMock) EnsureActivePorts(activePorts resources.ActivePorts) {

}

func (am *anybarManagerMock) StartWithIcon(port int, title, icon string) {

}

func (am *anybarManagerMock) StartNew(port int, title string) {

}

func (am *anybarManagerMock) ChangeIcon(port int, colour string) {

}

func (am *anybarManagerMock) GetNewPort(activePorts []*resources.ActivePort) int {
	return am.NewPortResult
}

func (am *anybarManagerMock) Quit(port int) {

}

func (am *anybarManagerMock) GetActivePorts(t resources.Transaction) resources.ActivePorts {
	return am.ActivePortsResult
}

func (am *anybarManagerMock) Ping(port int) bool {
	return am.PingResult
}

