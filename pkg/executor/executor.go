package executor

import (
	"io"
	"os/exec"
	"strings"
)

func Exec(command string) ([]byte, error) {
	parseCommand := strings.Split(command, " ")
	cmd := exec.Command(parseCommand[0], parseCommand[1:]...)

	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(pipe)
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func Close() error {
	return nil
}
