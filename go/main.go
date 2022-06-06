package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"os"
)

type Domain struct {
		IP  	string
		Name   string
	}

func main(){
	var domains []Domain

	containers := getContainers()

	for _, container := range containers {
		hasPiholeLabels, containerLabels := getContainerLabels(container)
		if (hasPiholeLabels) {
			domain := containerLabels
			domains = append(domains, domain)
		}
	}

	writeFile("./dit nog aanpassen.txt", domains)

}

func writeFile(filename string, domains []Domain){
// this does not write to the file, but only to the temp file currently
	f, err := os.CreateTemp("./", "example")
		if err != nil {
			panic(err)
		}
		//defer os.Remove(f.Name()) // defer clean up

		for _, domain := range domains {

			f.WriteString(domain.IP + " " + domain.Name + "\n")

			if err != nil {
		    panic(err)
		  }

		}



}

func getContainers() []types.Container {
  ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

  containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
  if err != nil {
    panic(err)
  }
  return containers
}

func getContainerLabels(container types.Container) (bool, Domain) {
	domain := Domain{}

		for key, value := range container.Labels{

			if key == "pihole.domain.ip" {
				fmt.Printf("Container label pihole.domain.ip: %s\n", value)
				domain.IP = value
			}

			if key == "pihole.domain.name" {
				domain.Name = value
				fmt.Printf("Container label pihole.domain.name: %s\n", value)
			}
		}

		if (Domain{} != domain) {
			return true, domain
		}

		fmt.Println("Container has no pihole labels")
		return false, domain
}
