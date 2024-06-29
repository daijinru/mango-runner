package runner

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	"github.com/daijinru/mango-runner/utils"
	"github.com/ttacon/chalk"
)

type Execution struct {
	// text will be printed line by line
	CommandStr string
	PrintLine  bool
	Pipeline   *Pipeline
}

func (ex *Execution) RunCommand(command string, args ...string) (string, error) {
	// fmt.Println(command, args[0])
	cmd := exec.Command(command, args...)
	combine := command + " " + utils.ConvertArrayToStr(args)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		message := fmt.Sprintf("unable to get stdout pipe: %s\n", err)
		ex.Pipeline.WriteError(err, combine)
		return "", fmt.Errorf(message)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		message := fmt.Sprintf("unable to get stderr pipe: %s\n", err)
		ex.Pipeline.WriteError(err, combine)
		return "", fmt.Errorf(message)
	}

	if err := cmd.Start(); err != nil {
		message := fmt.Sprintf("unable to start execution: %s\n", err)
		ex.Pipeline.WriteError(err, combine)
		return "", fmt.Errorf(message)
	}

	scanner := bufio.NewScanner(stdoutPipe)
	var output string
	for scanner.Scan() {
		text := scanner.Text()
		output = text + "\n"
		if ex.Pipeline != nil {
			ex.Pipeline.WriteInfo(output)
		}
		if ex.PrintLine {
			msg := fmt.Sprintf("[%s] %s", utils.TimeNow(), text)
			fmt.Println(chalk.White, msg)
		}
	}
	if err := scanner.Err(); err != nil {
		message := fmt.Sprintf("error output while scanning stdout %s\n", err)
		ex.Pipeline.WriteError(err, combine)
		return "", fmt.Errorf(message)
	}

	scannerErr := bufio.NewScanner(stderr)
	for scannerErr.Scan() {
		text := scannerErr.Text()
		output = text + "\n"
		if ex.Pipeline != nil {
			ex.Pipeline.WriteInfo(output)
		}
		if ex.PrintLine {
			msg := fmt.Sprintf("[%s] %s", utils.TimeNow(), text)
			fmt.Println(chalk.Red, msg)
		}
	}
	if err := scannerErr.Err(); err != nil {
		message := fmt.Sprintf("error output while scanning stderr %s\n", err)
		ex.Pipeline.WriteError(err, combine)
		return "", fmt.Errorf(message)
	}

	if err := cmd.Wait(); err != nil {
		message := fmt.Sprintf("error occured: %s", combine)
		ex.Pipeline.WriteInfo(message + "\n")
		ex.Pipeline.WriteError(err, combine)
		return "", fmt.Errorf(message)
	}

	return output, err
}

func (ex *Execution) RunCommandSplit() (string, error) {
	arr := strings.Split(ex.CommandStr, " ")
	output, err := ex.RunCommand(arr[0], arr[1:]...)
	if err != nil {
		return "", err
	}
	return output, err
}
