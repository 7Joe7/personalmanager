package test

import (
	"github.com/7joe7/personalmanager/resources"
	"github.com/stretchr/testify/mock"
)

func NewAnybarManagerMock() *anybarManagerMock {
	return &anybarManagerMock{}
}

type anybarManagerMock struct {
	mock.Mock
}

func (m *anybarManagerMock) RemoveAndQuit(bucketName []byte, id string, t resources.Transaction) {
	m.Called(bucketName, id, t)
	return
}

func (m *anybarManagerMock) AddToActivePorts(title, icon string, bucketName []byte, id string, t resources.Transaction) {
	m.Called(title, icon, bucketName, id, t)
}

func (m *anybarManagerMock) EnsureActivePorts(activePorts resources.ActivePorts) {
	m.Called(activePorts)
}

func (m *anybarManagerMock) StartWithIcon(port int, title, icon string) {
	m.Called(port, title, icon)
}

func (m *anybarManagerMock) StartNew(port int, title string) {
	m.Called(port, title)
}

func (m *anybarManagerMock) ChangeIcon(port int, colour string) {
	m.Called(port, colour)
}

func (m *anybarManagerMock) GetNewPort(activePorts []*resources.ActivePort) int {
	args := m.Called(activePorts)
	return args.Get(0).(int)
}

func (m *anybarManagerMock) Quit(port int) {
	m.Called(port)
}

func (m *anybarManagerMock) GetActivePorts(t resources.Transaction) resources.ActivePorts {
	args := m.Called(t)
	return args.Get(0).(resources.ActivePorts)
}

func (m *anybarManagerMock) Ping(port int) bool {
	args := m.Called(port)
	return args.Get(0).(bool)
}
