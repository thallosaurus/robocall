package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

func ConvertToGSM(filename string, d []byte) (*os.File, error) {
	fname := fmt.Sprintf("/tmp/%s-%s", string(time.Now().Unix()), filename)
	input_file, err := os.Create(fname)
	/*input_file, err := os.CreateTemp("", "in_sox")
	 */
	if err != nil {
		return nil, err
	}
	defer input_file.Close()

	io.Copy(input_file, bytes.NewBuffer(d))

	fname_out := fmt.Sprintf("/tmp/%s-%s.gsm", string(time.Now().Unix()))
	output_file, err := os.Create(fname_out)

	cmd := exec.Command("sox", "-r", "8000", input_file.Name(), output_file.Name())
	stdoutpipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderrpipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	cmd.Start()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	io.Copy(&stdout, stdoutpipe)
	io.Copy(&stderr, stderrpipe)

	fmt.Println(stdout.String(), stderr.String())

	//log.Print(data.String())
	return output_file, nil
}
