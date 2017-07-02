package anybar

import (
	"testing"
	"time"

	"github.com/7joe7/personalmanager/resources"
	"github.com/stretchr/testify/assert"
)

var (
	m anybarManager = anybarManager{"./"}
)

func TestStartNew(t *testing.T) {
	resources.WaitGroup.Add(1)
	nm := NewAnybarManager("./")
	nm.StartNew(1736, "ahoj")
	time.Sleep(100 * time.Millisecond)
	assert.Nil(t, m.sendCommand(1736, resources.ANY_CMD_QUIT))
}

func TestChangeIcon(t *testing.T) {
	nm := NewAnybarManager("./")
	m.start(1737, "ahoj")
	time.Sleep(100 * time.Millisecond)
	resources.WaitGroup.Add(1)
	nm.ChangeIcon(1737, resources.ANY_CMD_BLUE)
	assert.Nil(t, m.sendCommand(1737, resources.ANY_CMD_QUIT))
}

func TestQuit(t *testing.T) {
	nm := NewAnybarManager("./")
	m.start(1738, "ahoj")
	time.Sleep(100 * time.Millisecond)
	assert.Nil(t, m.sendCommand(1738, resources.ANY_CMD_BLUE))
	resources.WaitGroup.Add(1)
	nm.Quit(1738)
}
