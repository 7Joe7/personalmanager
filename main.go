package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"encoding/json"
	"net"

	"github.com/7joe7/personalmanager/resources"
	rutils "github.com/7joe7/personalmanager/resources/utils"
	"io"
	"strings"
)

var (
	// parameters
	action, id, name, projectId, goalId, taskId, habitId, repetition, deadline, estimate, scheduled, taskType, note *string
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

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panicked. %v %s\n", r, string(debug.Stack()))
			log.Fatalf("Panicked. %v %s", r, string(debug.Stack()))
			os.Exit(3)
		}
	}()

	flag.Parse()

	f, err := os.OpenFile(fmt.Sprintf("%s/%s", rutils.GetAppSupportFolderPath(), resources.LOG_FILE_NAME), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)

	logBinaryCall()

	cmd := resources.NewCommand(*action, *id, *name, *projectId, *goalId, *taskId,
		*repetition, *deadline, *estimate, *scheduled,
		*taskType, *note, *noneAllowed, *activeFlag,
		*doneFlag, *donePrevious, *undonePrevious, *negativeFlag,
		*learnedFlag, *basePoints, *habitRepetitionGoal)

	cmdBytes, err := json.Marshal(cmd)
	if err != nil {
		panic(err)
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
