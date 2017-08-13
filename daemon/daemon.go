package main

import (
    "os"
    "os/signal"
    "syscall"
    "fmt"
    "time"
    "strings"
    "strconv"
    "bufio"
    "io"
    "net"
    "log"
    "runtime/debug"
    "encoding/json"

    rutils "github.com/7joe7/personalmanager/resources/utils"
    "github.com/7joe7/personalmanager/resources"
    "github.com/7joe7/personalmanager/operations"
    "github.com/7joe7/personalmanager/operations/alfred"
    "github.com/7joe7/personalmanager/operations/anybar"
    "github.com/7joe7/personalmanager/db"
    "github.com/7joe7/personalmanager/utils"
    "github.com/7joe7/personalmanager/operations/goals"
    "github.com/7joe7/personalmanager/operations/configuration"

    "github.com/everdev/mack"
    "github.com/pkg/errors"
    "github.com/7joe7/personalmanager/operations/exporter"
    "github.com/7joe7/personalmanager/operations/alarm"
)

func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Panicked. %v %s\n", r, string(debug.Stack()))
            log.Fatalf("Panicked. %v %s", r, string(debug.Stack()))
            os.Exit(3)
        }
    }()

    // initialize catching the interrupt and terminate signal
    c := make(chan os.Signal)
    signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

    logFilePath := fmt.Sprintf("%s/%s", rutils.GetAppSupportFolderPath(), resources.LOG_DAEMON_FILE_NAME)
    log.Printf("logging into: %v", logFilePath)
    f, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    log.SetOutput(f)

    err = ensureAppSupportFolder(rutils.GetAppSupportFolderPath())
    if err != nil {
        panic(err)
    }

    runningBinary := utils.GetRunningBinaryPath()
    log.Println("running binary:", runningBinary)

    resources.Alf = alfred.NewAlfred()
    resources.Abr = anybar.NewAnybarManager(runningBinary)
    resources.Alr = alarm.NewAlarmManager()

    go resources.Alr.Run()

    db.Open()
    t := db.NewTransaction()
    operations.InitializeBuckets(t)
    operations.EnsureValues(t)
    operations.Synchronize(t)
    t.Execute()

    resources.Alr.Sync()

    listen, err := net.Listen("tcp4", ":" + strconv.Itoa(resources.PORT))
    if err != nil {
        log.Fatalf("socket listen port %d failed, %v", resources.PORT, err)
        os.Exit(1)
    }
    defer listen.Close()
    log.Printf("beginning to listen on port: %d", resources.PORT)

    // the daemon routine
    go func () {
        defer func() {
            if r := recover(); r != nil {
                fmt.Printf("Panicked. %v %s\n", r, string(debug.Stack()))
                log.Fatalf("Panicked. %v %s", r, string(debug.Stack()))
                os.Exit(3)
            }
        }()
        log.Println("starting to accept messages")
        for {
            // accepting any messages coming through
            conn, err := listen.Accept()
            if err != nil {
                log.Fatalln(err)
                continue
            }
            log.Println("acknowledged a connection")
            err = handleMessage(conn)
            if err != nil {
                log.Printf("unable to handle command. %v\n", err)
            }
        }
    }()

    mack.Tell("Alfred 3", "run trigger \"SyncAnybarPorts\" in workflow \"org.erneker.personalmanager\"")

    // waiting for the signals
    <-c
    signal.Stop(c)
    close(c)
    os.Exit(0)
}

func handleMessage(conn net.Conn) error {
    defer conn.Close()
    buf := make([]byte, 4096)
    r   := bufio.NewReader(conn)
    var data string

    for {
        n, err := r.Read(buf)
        data += string(buf[:n])

        switch err {
        case io.EOF:
            return nil
        case nil:
            log.Println("received:", data)
            if strings.HasSuffix(data, resources.STOP_CHARACTER) {
                data = strings.TrimSpace(data)
                cmd := &resources.Command{}
                err = json.Unmarshal([]byte(data), cmd)
                if err != nil {
                    return err
                }
                log.Printf("received command: %v\n", cmd)
                err = handleCommand(cmd, conn)
                if err != nil {
                    return err
                }
                if cmd.Alarm != "" {
                    resources.Alr.Sync()
                }
                return nil
            }
        default:
            return err
        }
    }
}

func handleCommand(cmd *resources.Command, conn net.Conn) error {
    t := db.NewTransaction()
    operations.Synchronize(t)
    t.Execute()
    switch cmd.Action {
    case resources.ACT_CREATE_TASK:
        operations.AddTask(cmd)
        resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "task"), conn)
    case resources.ACT_CREATE_PROJECT:
        operations.AddProject(cmd)
        resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "project"), conn)
    case resources.ACT_CREATE_TAG:
        operations.AddTag(cmd)
        resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "tag"), conn)
    case resources.ACT_CREATE_GOAL:
        goals.AddGoal(cmd)
        resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "goal"), conn)
    case resources.ACT_CREATE_HABIT:
        operations.AddHabit(cmd)
        resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_CREATE_SUCCESS, "habit"), conn)
    case resources.ACT_PRINT_TASKS:
        resources.Alf.PrintEntities(resources.Tasks{Tasks: operations.GetTasks(), NoneAllowed: cmd.NoneAllowed, Status: operations.GetStatus(), Sum: true}, conn)
    case resources.ACT_PRINT_PERSONAL_TASKS:
        resources.Alf.PrintEntities(resources.Tasks{Tasks: operations.GetPersonalTasks(), NoneAllowed: cmd.NoneAllowed, Status: operations.GetStatus(), Sum: true}, conn)
    case resources.ACT_PRINT_PERSONAL_NEXT_TASKS:
        resources.Alf.PrintEntities(resources.Tasks{Tasks: operations.GetNextTasks(), NoneAllowed: cmd.NoneAllowed, Status: operations.GetStatus(), Sum: true}, conn)
    case resources.ACT_PRINT_PERSONAL_UNSCHEDULED_TASKS:
        resources.Alf.PrintEntities(resources.Tasks{Tasks: operations.GetUnscheduledTasks(), NoneAllowed: cmd.NoneAllowed, Status: operations.GetStatus(), Sum: true}, conn)
    case resources.ACT_PRINT_SHOPPING_TASKS:
        resources.Alf.PrintEntities(resources.Tasks{Tasks: operations.GetShoppingTasks(), NoneAllowed: cmd.NoneAllowed, Status: operations.GetStatus()}, conn)
    case resources.ACT_PRINT_WORK_NEXT_TASKS:
        resources.Alf.PrintEntities(resources.Tasks{Tasks: operations.GetWorkNextTasks(), NoneAllowed: cmd.NoneAllowed, Status: operations.GetStatus(), Sum: true}, conn)
    case resources.ACT_PRINT_WORK_UNSCHEDULED_TASKS:
        resources.Alf.PrintEntities(resources.Tasks{Tasks: operations.GetWorkUnscheduledTasks(), NoneAllowed: cmd.NoneAllowed, Status: operations.GetStatus(), Sum: true}, conn)
    case resources.ACT_PRINT_TASK_NOTE:
        resources.Alf.PrintResult(operations.GetTask(cmd.ID).Note, conn)
    case resources.ACT_PRINT_PROJECTS:
        resources.Alf.PrintEntities(resources.Projects{operations.GetProjects(), cmd.NoneAllowed, operations.GetStatus()}, conn)
    case resources.ACT_PRINT_ACTIVE_PROJECTS:
        resources.Alf.PrintEntities(resources.Projects{operations.GetActiveProjects(), cmd.NoneAllowed, operations.GetStatus()}, conn)
    case resources.ACT_PRINT_INACTIVE_PROJECTS:
        resources.Alf.PrintEntities(resources.Projects{operations.GetInactiveProjects(), cmd.NoneAllowed, operations.GetStatus()}, conn)
    case resources.ACT_PRINT_TAGS:
        resources.Alf.PrintEntities(resources.Tags{operations.GetTags(), cmd.NoneAllowed, operations.GetStatus()}, conn)
    case resources.ACT_PRINT_GOALS:
        resources.Alf.PrintEntities(resources.Goals{goals.GetGoals(), cmd.NoneAllowed, operations.GetStatus()}, conn)
    case resources.ACT_PRINT_ACTIVE_GOALS:
        resources.Alf.PrintEntities(resources.Goals{goals.GetActiveGoals(), cmd.NoneAllowed, operations.GetStatus()}, conn)
    case resources.ACT_PRINT_NON_ACTIVE_GOALS:
        resources.Alf.PrintEntities(resources.Goals{goals.GetNonActiveGoals(), cmd.NoneAllowed, operations.GetStatus()}, conn)
    case resources.ACT_PRINT_INCOMPLETE_GOALS:
        resources.Alf.PrintEntities(resources.Goals{goals.GetIncompleteGoals(), cmd.NoneAllowed, operations.GetStatus()}, conn)
    case resources.ACT_PRINT_HABITS:
        if cmd.ActiveFlag {
            resources.Alf.PrintEntities(resources.Habits{operations.GetActiveHabits(), cmd.NoneAllowed, operations.GetStatus(), true}, conn)
        } else {
            resources.Alf.PrintEntities(resources.Habits{operations.GetNonActiveHabits(), cmd.NoneAllowed, operations.GetStatus(), false}, conn)
        }
    case resources.ACT_PRINT_HABIT_DESCRIPTION:
        resources.Alf.PrintResult(operations.GetHabit(cmd.ID).Description, conn)
    case resources.ACT_PRINT_REVIEW:
        resources.Alf.PrintEntities(resources.Items{[]*resources.AlfredItem{operations.GetReview().GetItem()}}, conn)
    case resources.ACT_EXPORT_SHOPPING_TASKS:
        exporter.ExportShoppingTasks(resources.CFG_EXPORT_CONFIG_PATH)
        resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_EXPORT_SUCCESS, "shopping tasks"), conn)
    case resources.ACT_DELETE_TASK:
        operations.DeleteTask(cmd.ID)
        resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_DELETE_SUCCESS, "task"), conn)
    case resources.ACT_DELETE_PROJECT:
        operations.DeleteProject(cmd.ID)
        resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_DELETE_SUCCESS, "project"), conn)
    case resources.ACT_DELETE_TAG:
        operations.DeleteTag(cmd.ID)
        resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_DELETE_SUCCESS, "tag"), conn)
    case resources.ACT_DELETE_GOAL:
        goals.DeleteGoal(cmd.ID)
        resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_DELETE_SUCCESS, "goal"), conn)
    case resources.ACT_DELETE_HABIT:
        operations.DeleteHabit(cmd.ID)
        resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_DELETE_SUCCESS, "habit"), conn)
    case resources.ACT_MODIFY_TASK:
        operations.ModifyTask(cmd)
        resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "task"), conn)
    case resources.ACT_MODIFY_PROJECT:
        operations.ModifyProject(cmd)
        resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "project"), conn)
    case resources.ACT_MODIFY_TAG:
        operations.ModifyTag(cmd)
        resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "tag"), conn)
    case resources.ACT_MODIFY_GOAL:
        goals.ModifyGoal(cmd)
        resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "goal"), conn)
    case resources.ACT_MODIFY_HABIT:
        operations.ModifyHabit(cmd)
        resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "habit"), conn)
    case resources.ACT_MODIFY_REVIEW:
        operations.ModifyReview(cmd)
        resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_MODIFY_SUCCESS, "review"), conn)
    case resources.ACT_SYNC_ANYBAR_PORTS:
        t := db.NewTransaction()
        operations.SynchronizeAnybarPorts(t)
        t.Execute()
        resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_SYNC_SUCCESS, "AnyBar ports"), conn)
    case resources.ACT_DEBUG_DATABASE:
        db.PrintoutDbContents(cmd.ID)
    case resources.ACT_SYNC_WITH_JIRA:
        operations.SyncWithJira()
    case resources.ACT_BACKUP_DATABASE:
        db.BackupDatabase()
    case resources.ACT_SET_CONFIG_VALUE:
        switch cmd.ID {
        case string(resources.DB_DEFAULT_EMAIL):
            exporter.SetEmail(cmd.Name)
            resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_SET_SUCCESS, "e-mail", cmd.Name), conn)
        case string(resources.DB_WEEKS_LEFT):
            configuration.SetWeeksLeft(cmd.BasePoints)
            weeksLeft := fmt.Sprint(cmd.BasePoints)
            if !resources.Abr.Ping(resources.ANY_PORT_WEEKS_LEFT) {
                resources.WaitGroup.Add(1)
                go resources.Abr.StartWithIcon(resources.ANY_PORT_WEEKS_LEFT, weeksLeft, resources.ANY_CMD_BROWN)
            }
            resources.Alf.PrintResult(fmt.Sprintf(resources.MSG_SET_SUCCESS, "weeks left", weeksLeft), conn)
        }
    case resources.ACT_CUSTOM:
        db.DeleteEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte("271"))
        //t := db.NewTransaction()
        //t.Add(func() error {
        //    //getNewHabit := func() resources.Entity {
        //    //    return &resources.Habit{}
        //    //}
        //    //err := t.MapEntities(resources.DB_DEFAULT_HABITS_BUCKET_NAME, true, getNewHabit, func(e resources.Entity) func() {
        //    //    return func() {
        //    //        h := e.(*resources.Habit)
        //    //        if !h.Active {
        //    //            return
        //    //        }
        //    //        if h.Repetition == resources.HBT_REPETITION_DAILY {
        //    //            for h.Deadline.Before(time.Now()) {
        //    //                h.Deadline = addPeriod(resources.HBT_REPETITION_DAILY, h.Deadline)
        //    //            }
        //    //        }
        //    //    }
        //    //})
        //    //habit := &resources.Habit{}
        //    //err = t.ModifyEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte("137"), true, habit, func () {
        //    //	habit.Successes = 6
        //    //	habit.LastStreak = 4
        //    //	habit.ActualStreak = 2
        //    //})
        //    if err != nil {
        //        return err
        //    }
        //    return nil
        //})
        //t.Execute()
    default:
        return errors.New("no such action")
    }
    resources.WaitGroup.Wait()
    return nil
}

// Creating application support folder if it doesn't exist
func ensureAppSupportFolder(appSupportFolderPath string) error {
    _, err := os.Stat(appSupportFolderPath)
    if err != nil {
        if os.IsNotExist(err) {
            err = os.Mkdir(appSupportFolderPath, os.FileMode(0744))
            if err != nil {
                return err
            }
        } else {
            return err
        }
    }
    return nil
}

func addPeriod(repetition string, deadline *time.Time) *time.Time {
    if deadline == nil {
        return nil
    }
    switch repetition {
    case resources.HBT_REPETITION_DAILY:
        return utils.GetTimePointer(deadline.Add(24 * time.Hour))
    case resources.HBT_REPETITION_WEEKLY:
        return utils.GetTimePointer(deadline.Add(7 * 24 * time.Hour))
    case resources.HBT_REPETITION_MONTHLY:
        return utils.GetTimePointer(deadline.AddDate(0, 1, 0))
    }
    return nil
}