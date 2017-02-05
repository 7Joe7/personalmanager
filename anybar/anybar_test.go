package anybar

import (
	"testing"
	"time"

	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/test"
)

func TestStartNew(t *testing.T) {
	output, err := start(1736, "ahoj")
	test.ExpectSuccess(t, err)
	test.ExpectString("", output, t)
	time.Sleep(100 * time.Millisecond)
	test.ExpectSuccess(t, sendCommand(1736, resources.ANY_CMD_QUIT))
}

func TestChangeIcon(t *testing.T) {
	output, err := start(1737, "ahoj")
	test.ExpectSuccess(t, err)
	test.ExpectString("", output, t)
	time.Sleep(100 * time.Millisecond)
	test.ExpectSuccess(t, sendCommand(1737, resources.ANY_CMD_BLUE))
	test.ExpectSuccess(t, sendCommand(1737, resources.ANY_CMD_QUIT))
}

func TestQuit(t *testing.T) {
	output, err := start(1738, "ahoj")
	test.ExpectSuccess(t, err)
	test.ExpectString("", output, t)
	time.Sleep(100 * time.Millisecond)
	test.ExpectSuccess(t, sendCommand(1738, resources.ANY_CMD_BLUE))
	test.ExpectSuccess(t, sendCommand(1738, resources.ANY_CMD_QUIT))
}
