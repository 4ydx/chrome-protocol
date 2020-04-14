package cdp

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"syscall"
	"time"
)

// Browser contains information required to stop an exec'd browser at a later point in time.
type Browser struct {
	Port        int
	PID         int
	TempDir     string
	LogFile     *os.File
	Log         *log.Logger
	ConsoleFile *os.File
	Console     *log.Logger
}

// NewBrowser accepts the path to the browser's binary, the port, and any arguments that need to be passed to the binary.
func NewBrowser(path string, port int, logfile string, args ...string) *Browser {
	b := &Browser{}

	var err error
	b.LogFile, err = os.Create(logfile)
	if err != nil {
		panic(err)
	}
	b.Log = log.New(b.LogFile, "", log.Llongfile|log.LstdFlags|log.Lmicroseconds)

	b.ConsoleFile, err = os.Create(logfile + ".console")
	if err != nil {
		panic(err)
	}
	b.Console = log.New(b.ConsoleFile, "", log.Lshortfile|log.LstdFlags)

	b.Port = port

	// Add required values
	required := []string{"--no-first-run", "--no-default-browser-check"}
	for _, req := range required {
		sort.Strings(args)
		at := sort.SearchStrings(args, req)
		if at == len(args) || args[at] != req {
			args = append(args, req)
		}
	}

	// Temp directory
	dir, err := ioutil.TempDir("", "cdp-")
	if err != nil {
		panic(err)
	}
	b.TempDir = dir
	args = append(args, fmt.Sprintf("--user-data-dir=%s", b.TempDir))

	// Debugging port of the browser
	debuggingPort := fmt.Sprintf("--remote-debugging-port=%d", b.Port)
	args = append(args, debuggingPort)
	log.Printf("Args %+v", args)

	cmd := exec.Command(path, args...)

	// Prepare to send output to the log file
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	// Start the browser
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	b.PID = cmd.Process.Pid

	// The browser does need a moment to start up.  This will probably differ on different systems.
	time.Sleep(time.Second * 1)

	// Feed output to the log
	go func() {
		sout := bufio.NewScanner(stdout)
		for sout.Scan() {
			log.Print(sout.Text())
		}
		if err := sout.Err(); err != nil {
			log.Printf("error: %s", err)
		}
	}()
	go func() {
		serr := bufio.NewScanner(stderr)
		for serr.Scan() {
			log.Print(serr.Text())
		}
		if err := serr.Err(); err != nil {
			log.Printf("error: %s", err)
		}
	}()
	return b
}

// Stop kills the running browser process.
func (b *Browser) Stop() {
	log.Print("stopping the browser")
	if b.PID == 0 {
		log.Print("no process id for the browser")
		return
	}
	err := syscall.Kill(b.PID, syscall.SIGKILL)
	if err != nil {
		panic(err)
	}
	if b.TempDir != "" {
		err := os.RemoveAll(b.TempDir)
		if err != nil {
			log.Printf("error: %s", err)
		}
	}
	err = b.LogFile.Close()
	if err != nil {
		log.Printf("error: %s", err)
	}
	err = b.ConsoleFile.Close()
	if err != nil {
		log.Printf("error: %s", err)
	}
}
