package utils

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
)

func ConvertToGSM(filepath string) {
	f, err := os.CreateTemp("", "sox")
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("sox", "-r", "8000", filepath, f.Name())
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	var data bytes.Buffer
	io.Copy(&data, f)
	log.Print(data.String())
}
