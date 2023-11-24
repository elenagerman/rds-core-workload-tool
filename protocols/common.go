package protocols

import (
	"fmt"
	"os/exec"
)

// CommonTest keeps common vars from connectivity tests
type CommonTest struct {
	MTU             int
	ServerIP        string
	ProtocolVersion int
	Negative        bool
}

//RunCommand runs command and return output
func (ct *CommonTest) RunCommand(cmd string) (string, error) {
	commandOutput, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("command execution failed - %s due to the error - %s", cmd, err)
	}
	fmt.Printf("Executed command: %s", cmd)
	return string(commandOutput), nil
}

func totalPackageLoss(total int, loss int) int {
	if loss != 0 {
		return int(float64(loss) / float64(total) * 100)
	}
	return 0
}
