package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"encoding/json"
	"net"
	"io"
	"strings"
	"text/template"

	"github.com/7joe7/personalmanager/resources"
	rutils "github.com/7joe7/personalmanager/resources/utils"
	"github.com/7joe7/personalmanager/utils"
	"os/exec"
)

var (
	// parameters
	action, id, name, projectId, goalId, taskId, habitId, repetition, deadline, alarm, estimate, scheduled, taskType, note *string
	noneAllowed, activeFlag, doneFlag, donePrevious, undonePrevious, negativeFlag, learnedFlag                      *bool
	basePoints, habitRepetitionGoal                                                                                 *int
)

func init() {
	action = flag.String("action", "", fmt.Sprintf("Provide action to be taken from this list: %v.", resources.ACTIONS))
	id = flag.String("id", "", "Provide id of the entity you want to make the action for. Valid for these actions: .")
	name = flag.String("name", "", "Provide name.")
	projectId = flag.String("projectId", "", "Provide project id for project assignment.")
	goalId = flag.String("goalId", "", "Provide goal id for goal assignment.")
	taskId = flag.String("taskId", "", "Provide task id for task assignment.")
	habitId = flag.String("habitId", "", "Provide habit id for habit assignment.")
	repetition = flag.String("repetition", "", "Select repetition period.")
	deadline = flag.String("deadline", "", "Specify deadline in format 'dd.MM.YYYY HH:mm'.")
	alarm = flag.String("alarm", "", "Specify alarm in format 'dd.MM.YYYY HH:mm'.")
	estimate = flag.String("estimate", "", "Specify time estimate in format '2h45m'.")
	scheduled = flag.String("scheduled", "", "Provide schedule period. (NEXT|NONE)")
	taskType = flag.String("taskType", "", "Provide task type. (PERSONAL|WORK)")
	note = flag.String("note", "", "Provide note.")
	noneAllowed = flag.Bool("noneAllowed", false, "Provide information whether list should be retrieved with none value allowed.")
	activeFlag = flag.Bool("active", false, "Toggle active/show active only.")
	doneFlag = flag.Bool("done", false, "Toggle done.")
	donePrevious = flag.Bool("donePrevious", false, "Set done for previous period.")
	undonePrevious = flag.Bool("undonePrevious", false, "Set undone for previous period.")
	negativeFlag = flag.Bool("negative", false, "Set negative flag for habits.")
	learnedFlag = flag.Bool("learned", false, "Set learned flag for habits.")
	basePoints = flag.Int("basePoints", -1, "Set base points for success/failure.")
	habitRepetitionGoal = flag.Int("habitRepetitionGoal", -1, "Set habit goal repetition number.")
}

type plistData struct {
	BinaryAddress string
	SupportFolderAddress string
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panicked. %v %s\n", r, string(debug.Stack()))
			log.Fatalf("Panicked. %v %s", r, string(debug.Stack()))
			os.Exit(3)
		}
	}()

	flag.Parse()

	runningBinary := utils.GetRunningBinaryPath()
	plistAddress := fmt.Sprintf("%s/Library/LaunchAgents/org.erneker.personalmanager.plist", os.Getenv("HOME"))
	if _, err := os.Stat(plistAddress); os.IsNotExist(err) {
		tmpl, err := template.ParseFiles(fmt.Sprintf("%s/org.erneker.personalmanager.plist.tmpl", runningBinary))
		if err != nil {
			log.Fatalln(err)
		}
		plist, err := os.Create(plistAddress)
		if err != nil {
			log.Fatalln(err)
		}
		pd := plistData{
			BinaryAddress: runningBinary,
			SupportFolderAddress: rutils.GetAppSupportFolderPath(),
		}
		err = tmpl.Execute(plist, pd)
		if err != nil {
			log.Fatalln(err)
		}
		out, err := exec.Command("launchctl", "load", plistAddress).CombinedOutput()
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("loaded personal manager daemon", string(out))
	}

	cmd := resources.NewCommand(*action, *id, *name, *projectId, *goalId, *taskId,
		*repetition, *deadline, *alarm, *estimate, *scheduled,
		*taskType, *note, *noneAllowed, *activeFlag,
		*doneFlag, *donePrevious, *undonePrevious, *negativeFlag,
		*learnedFlag, *basePoints, *habitRepetitionGoal)

	cmdBytes, err := json.Marshal(cmd)
	if err != nil {
		log.Fatalln(err)
	}

	addr := fmt.Sprintf("127.0.0.1:%d", resources.PORT)
	conn, err := net.Dial("tcp", addr)

	defer conn.Close()

	if err != nil {
		log.Fatalln(err)
	}

	_, err = conn.Write(cmdBytes)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = conn.Write([]byte(resources.STOP_CHARACTER))
	if err != nil {
		log.Fatalln(err)
	}
	buf := make([]byte, 4096)
	var result string
	for {
		n, err := conn.Read(buf)
		result += string(buf[:n])

		switch err {
		case io.EOF:
			return
		case nil:
			if strings.HasSuffix(result, resources.STOP_CHARACTER) {
				result = strings.TrimSpace(result)
				_, err = fmt.Fprint(os.Stdout, result)
				if err != nil {
					log.Fatalln(err)
				}
			}
		default:
			log.Fatalf("receive of data failed: %v", err)
		}
	}
}

func logBinaryCall() {
	log.Printf(`Called with string parameters:
		action: %s,
		id: %s,
		name: %s,
		projectId: %s,
		goalId: %s,
		taskId: %s,
		habitId: %s,
		repetition: %s,
		deadline: %s,
		estimate: %s,
		scheduled: %s,
		taskType: %s,
		note: %s,
		and with bool parameters:
		noneAllowed: %v,
		activeFlag: %v,
		doneFlag: %v,
		donePrevious: %v,
		undonePrevious: %v,
		and int parameters:
		basePoints: %v,
		repetitionGoal: %v.`, *action, *id, *name, *projectId, *goalId, *taskId,
		*habitId, *repetition, *deadline, *estimate, *scheduled, *taskType, *note, *noneAllowed, *activeFlag,
		*doneFlag, *donePrevious, *undonePrevious, *basePoints, *habitRepetitionGoal)
}
