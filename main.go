package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Domain struct {
	Name string
	IP   string
}

func main() {

	containers, err := getContainers()
	if err != nil{
		log.Fatal(err.Errors())
	}

	domains := getPiholeDomains(containers)
	domains.sort()

	err := writeFile("./dit nog aanpassen.txt", domains)
	if err != nil{
		log.Fatal(err.Errors())
	}

}

func (d []Domain) sort() []Domain{

	sort.SliceStable(d, func(i, j int) bool {
		return d[i].Name < d[j].Name
	})
}

func getPiholeDomains(containers []types.Container) []Domain {
	var domains []Domain

	for _, container := range containers {
		log.Printf("Found container with name:", strings.Trim(fmt.Sprint(container.Names), "/[]"))

		hasPiholeLabels, piholeLabels := getPiholeLabels(container)

		if hasPiholeLabels {
			domain := piholeLabels
			domains = append(domains, domain)
		}
	}
	return domains
}

func writeFile(filename string, domains []Domain) error {
	// this does not write to the file, but only to the temp file currently
	f, err := os.CreateTemp("./", "example")
	if err != nil {
		return err
	}

	defer os.Remove(f.Name()) // defer clean up

	for _, domain := range domains {

		f.WriteString(domain.IP + " " + domain.Name + "\n")

		if err != nil {
			return err
		}
	}
	return nil
}

func getContainers() ([]types.Container, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return []types.Container, error
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return []types.Container, error
	}
	return containers, nil
}

func getIP(labels container.Labels) string {
	for key, value := range labels {
		if key == "pihole.domain.ip" {
			log.Printf("Container label pihole.domain.ip: %s\n", value)
			return value
		}
	}
	return ""
} 

func getName(labels container.Labels) string {
	for key, value := range labels {
		if key == "pihole.domain.name" {
			log.Printf("Container label pihole.domain.name: %s\n", value)
			return value
		}
	}
	return ""
} 

func getPiholeLabels(container types.Container) (bool, Domain) {
	domain := Domain{}

	domain.Name = getName(container.Labels)
	domain.IP = getIP(container.Labels)

	if (domain == Domain{}) {
		log.Println("Container has no pihole labels")
		return false, domain
	}

	return true, domain
}
