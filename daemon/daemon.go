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

    rutils "github.com/7joe7/personalmanager/resources/utils"

    "github.com/everdev/mack"
)

const (
    STOP_CHARACTER = "\r\n\r\n"
    PORT = 7007
    IP = "127.0.0.1"
)

func main() {
    // initialize catching the interrupt and terminate signal
    c := make(chan os.Signal)
    signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

    err := ensureAppSupportFolder(rutils.GetAppSupportFolderPath())
    if err != nil {
        panic(err)
    }

    mack.Tell("Alfred 3", "run trigger \"SyncAnybarPorts\" in workflow \"org.erneker.personalmanager\"")

    listen, err := net.Listen("tcp4", ":" + strconv.Itoa(PORT))
    if err != nil {
        log.Fatalf("Socket listen port %d failed, %v", PORT, err)
        os.Exit(1)
    }
    defer listen.Close()
    log.Printf("Begin listen port: %d", PORT)

    go func () {
        addr := fmt.Sprintf("%s:%d", IP, PORT)
        conn, err := net.Dial("tcp", addr)

        defer conn.Close()

        if err != nil {
            log.Fatalln(err)
        }

        conn.Write([]byte("hi"))
        conn.Write([]byte(STOP_CHARACTER))
        log.Printf("Send: %s", "hi")
    }()

    // the daemon routine
    go func () {
        for {
            fmt.Println("still running")
            conn, err := listen.Accept()
            if err != nil {
                log.Fatalln(err)
                continue
            }
            go handleConnection(conn)
            time.Sleep(time.Second * 60)
        }
    }()

    // waiting for the signals
    <-c
    signal.Stop(c)
    close(c)
    os.Exit(0)
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    buf := make([]byte, 1024)
    r   := bufio.NewReader(conn)

    for {
        n, err := r.Read(buf)
        data := string(buf[:n])

        switch err {
        case io.EOF:
            return
        case nil:
            log.Println("Receive:", data)
            if strings.HasSuffix(data, STOP_CHARACTER) {
                return
            }
        default:
            log.Fatalf("Receive data failed: %v", err)
            return
        }
    }
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
