package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func findContainer(c *client.Client, nameToSearch string) (*types.Container, error) {
	containers, err := c.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Could not list containers: %w", err)
	}

	for _, c := range containers {
		container := c
		for _, n := range c.Names {
			if strings.HasSuffix(n, nameToSearch) {
				return &container, nil
			}
		}
	}
	return nil, fmt.Errorf("Container '%s' does not exist.", nameToSearch)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s CONTAINER COMMAND\n", os.Args[0])
		return
	}

	name := os.Args[1]
	command := strings.Join(os.Args[2:], " ")

	dockerC, err := client.NewClientWithOpts()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	container, err := findContainer(dockerC, name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = SendInput(dockerC, container, command+"\n")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
