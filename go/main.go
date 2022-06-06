package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"os"
	"strings"
	"sort"
)

type Domain struct {
		IP  	string
		Name   string
	}

func main(){

	containers := getContainers()

	domains := getPiholeDomains(containers)

	sort.SliceStable(domains, func(i, j int) bool {
		return domains[i].Name < domains[j].Name
	})

	fmt.Println("after sorting")
	for _, domain := range domains{
		fmt.Println(domain.Name)
	}


	//fmt.Printf("%t", h == h2)
	//writeFile("./dit nog aanpassen.txt", domains)

}

func getPiholeDomains(containers []types.Container) []Domain {
	var domains []Domain

	for _, container := range containers {
		fmt.Println("---")
		fmt.Println("Found container with name:", strings.Trim(fmt.Sprint(container.Names), "/[]"))

		hasPiholeLabels, piholeLabels := getPiholeLabels(container)

		if (hasPiholeLabels) {
			domain := piholeLabels
			domains = append(domains, domain)
		}
	}
	return domains
}

func writeFile(filename string, domains []Domain){
// this does not write to the file, but only to the temp file currently
	f, err := os.CreateTemp("./", "example")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name()) // defer clean up

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

func getPiholeLabels(container types.Container) (bool, Domain) {
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
