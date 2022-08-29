package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"tunl-cli/cmd/tui"
)

type Env = []string

func runPrint(cmd string, env Env, args ...string) {
	log.Println(cmd, strings.Join(args, " "))
	ecmd := exec.Command(cmd, args...)
	ecmd.Env = os.Environ()
	ecmd.Env = append(ecmd.Env, env...)
	ecmd.Stdout = os.Stdout
	ecmd.Stderr = os.Stderr
	err := ecmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func buildCli() {
	os.MkdirAll("./build", 0750)

	runPrint("npm", nil, "-v")
	runPrint("npm", nil, "i", "--prefix", "./ui/monitor")
	runPrint("npm", nil, "run", "build", "--prefix", "./ui/monitor")
	runPrint("go", Env{"GOOS=windows", "GOARCH=amd64"}, "build", "-o", "./build/tunl-amd64.exe", "./cmd")
	runPrint("go", Env{"GOOS=windows", "GOARCH=386"}, "build", "-o", "./build/tunl-386.exe", "./cmd")
	runPrint("go", Env{"GOOS=darwin", "GOARCH=amd64"}, "build", "-o", "./build/tunl-amd64-darwin", "./cmd")
	runPrint("go", Env{"GOOS=linux", "GOARCH=amd64"}, "build", "-o", "./build/tunl-amd64-linux", "./cmd")
	runPrint("go", Env{"GOOS=linux", "GOARCH=386"}, "build", "-o", "./build/tunl-386-linux", "./cmd")
}

func main() {
	tui.PrintInfo("Build tunl.online cli...")
	buildCli()
	tui.PrintInfo("Done!")
}
