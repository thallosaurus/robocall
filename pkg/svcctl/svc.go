package svcctl

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
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

func start(stdout *bytes.Buffer) {
	cmd := exec.Command("asterisk", "-f")

	stdoutpipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stdoutpipe.Close()

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	for {
		b := make([]byte, 1024)

		_, err := stdoutpipe.Read(b)
		if err != nil {
			if err.Error() == "EOF" {
				log.Println("Reached EOF")
				cmd.Wait()
				fmt.Println("child process exited")
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

var s = make(chan int)

func RunService() error {
	if notRunning() {
		var stdout bytes.Buffer
		go start(&stdout)

		return nil
	} else {
		pid, err := os.ReadFile("/var/run/asterisk.pid")
		if err != nil {
			log.Fatal(err)
		}
		return fmt.Errorf("asterisk is already running (%s)", string(pid))
	}
}

func sendToSocket(ctl ServiceControl) error {
	//log.Println("Reloading asterisk")
	cmd := exec.Command("asterisk", "-rx", fmt.Sprintf("%s", ctl))
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
	//if !notRunning() {
	//} else {
	//log.Fatal("service not running")
	//}
}
