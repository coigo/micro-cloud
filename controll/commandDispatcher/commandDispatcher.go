package commandDispatcher

import (
	"fmt"
	"os/exec"
)

func UpCommand () string {
	cmd := exec.Command("docker", "exec", "-i", "a195e30364d5", "docker", "run", "-d", "alpine", "sleep", "infinity")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	return string(out)
	
	
}

func DownCommand (id string) {
		cmd := exec.Command("docker", "exec", "-i", "a195e30364d5", "docker", "stop", id)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}