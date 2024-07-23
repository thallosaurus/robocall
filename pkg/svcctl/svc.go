package svcctl

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

type ServiceControl = string

const (
	Restart     ServiceControl = "core restart now"
	Stop        ServiceControl = "core stop now"
	ReloadPJSIP ServiceControl = "module reload res_pjsip.so"
)

func notRunning() bool {
	_, err := os.Stat("/var/run/asterisk.pid")

	return errors.Is(err, os.ErrNotExist)
}

type SvcSignals = int

const (
	StartSvc SvcSignals = 0
	StopSvc  SvcSignals = 1
)

func start(wg *sync.WaitGroup, stdout *bytes.Buffer, ret *chan int) {
	defer wg.Done()
	cmd := exec.Command("asterisk", "-f")

	//cmd.Cancel()

	stdoutpipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stdoutpipe.Close()

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	wg.Add(1)

	for {

		// Read Stdin
		b := make([]byte, 1024)

		_, err := stdoutpipe.Read(b)
		if err != nil {

			// Program Closed an Pipe got drained
			if err.Error() == "EOF" {
				log.Println("Reached EOF")
				cmd.Wait()
				//fmt.Println("child process exited")
				*ret <- 1
				break
			} else {
				log.Fatal("stdout err ", err)
			}
		}

		_, err = stdout.Write(b)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(string(b))
	}
}

func RunService(wg *sync.WaitGroup) (chan int, error) {
	if notRunning() {
		c := make(chan int)
		var stdout bytes.Buffer
		go start(wg, &stdout, &c)

		return c, nil
	} else {
		pid, err := os.ReadFile("/var/run/asterisk.pid")
		if err != nil {
			log.Fatal(err)
		}
		return nil, fmt.Errorf("asterisk is already running (%s)", string(pid))
	}
}

func sendToSocket(ctl ServiceControl) error {
	//log.Println("Reloading asterisk")
	cmd := exec.Command("asterisk", "-rx", ctl)
	//err := cmd.Run()
	b, err := cmd.Output()

	fmt.Print(string(b))

	return err
}

func StopService() {
	sendToSocket(Stop)
}

func ReloadSIPModule() {
	sendToSocket(ReloadPJSIP)
}
