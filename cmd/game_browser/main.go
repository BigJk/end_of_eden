package main

import (
	"errors"
	"fmt"
	"github.com/samber/lo"
	"os"
	"os/exec"
	"runtime"
)

const GameExecutable = "end_of_eden"

func locateTtyd() (string, error) {
	locations := []string{
		"ttyd",
		"./ttyd/ttyd.win32.exe",
		"./ttyd/ttyd.aarch64",
		"./ttyd/ttyd.arm",
		"./ttyd/ttyd.armhf",
		"./ttyd/ttyd.i686",
		"./ttyd/ttyd.mips",
		"./ttyd/ttyd.mips64",
		"./ttyd/ttyd.mips64el",
		"./ttyd/ttyd.mipsel",
		"./ttyd/ttyd.s390x",
		"./ttyd/ttyd.x86_64",
	}

	for i := range locations {
		ttydVersion, err := exec.Command(locations[i], "-v").CombinedOutput()
		if err != nil {
			continue
		}
		fmt.Printf("%s version: %s", locations[i], string(ttydVersion))
		return locations[i], nil
	}

	return "", errors.New("ttyd not found")
}

func main() {
	args := os.Args[1:]

	ext := ""
	if runtime.GOOS == "window" {
		ext = ".exe"
	}

	// Check if ttyd exists
	ttyd, err := locateTtyd()
	if err != nil {
		panic(err)
	}

	// Build ttyd command
	game := exec.Command(ttyd, lo.Flatten([][]string{
		{"--check-origin", "--browser", "--once", "./" + GameExecutable + ext},
		args, // pass args to end_of_eden
	})...)
	game.Stdout = os.Stdout
	game.Stderr = os.Stderr

	// Run ttyd
	if err := game.Run(); err != nil {
		panic(err)
	}

	fmt.Println("Shutdown!")
}
