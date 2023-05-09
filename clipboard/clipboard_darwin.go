package clipboard

import (
	"os/exec"
	"syscall"
)

func Init() {
	// nothing
}

func Set(data string) {
	cmd := exec.Command("pbcopy")
	if stdin, err := cmd.StdinPipe(); err == nil {
		if cmd.Start() == nil {
			_, _ = stdin.Write([]byte(data))
			_ = stdin.Close()
			_ = cmd.Process.Signal(syscall.SIGPIPE)
		}
	}
}
