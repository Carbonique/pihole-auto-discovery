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
	IP   string
	Name string
}

type Domains []Domain

var outputFile = "./test"

func main() {

	log.Println("Run started")

	containers, err := getContainers()
	if err != nil {
		log.Fatal(err.Error())
	}

	domains, err := getPiholeLabels(containers)
	if err != nil {
		log.Fatal(err.Error())
	}

	domains.sort()
	err = writeDNSList(outputFile, domains)
	if err != nil {
		log.Fatal(err.Error())
	}
	//	readFile(outputFile, domains)
	log.Println("Run finished")
}

func (d Domains) sort() Domains {
	sort.SliceStable(d, func(i, j int) bool {
		return d[i].Name < d[j].Name
	})
	return Domains{}
}

//getPiholeLabels retrieves Docker containers with pihole labels
func getPiholeLabels(containers []types.Container) (Domains, error) {

	domains := Domains{}
	for _, container := range containers {

		domain := newDomain(container)

		// If app is empty, we do not append
		if domain == (Domain{}) {
			log.Println("Container has no sui labels")
			continue
		}
		domains = append(domains, domain)
	}
	return domains, nil
}

func newDomain(c types.Container) Domain {
	log.Printf("Parsing labels from container: %+q", c.Names)

	domain := Domain{}

	domain.Name = parseName(c.Labels)
	domain.IP = parseIP(c.Labels)

	return domain

}

func writeDNSList(outputFile string, domains []Domain) error {

	f, err := createFileIfNotExists(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, domain := range domains {
		f.WriteString(domain.IP + " " + domain.Name + "\n")
	}

	return nil
}

func getContainers() ([]types.Container, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return []types.Container{}, err
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return []types.Container{}, err
	}

	return containers, nil

}

func parseName(m map[string]string) string {
	if val, ok := m["pihole.domain.name"]; ok {
		log.Printf("Container label pihole.domain.name: %s\n", val)
		return val
	}
	return ""
}

func parseIP(m map[string]string) string {
	if val, ok := m["pihole.domain.ip"]; ok {
		log.Printf("Container label pihole.domain.ip: %s\n", val)
		return val
	}
	return ""
}

func createFileIfNotExists(file string) (*os.File, error) {

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

	if err != nil {
		return nil, err
	}

	return f, nil
}

func readFile(fileName string, domains Domains) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		os.Exit(1)
	}

	s := string(b)
	sp := strings.Split(s, "\n")

	for _, v := range sp {
		fmt.Println(v)
		for _, v2 := range domains {
			fmt.Println(v2)
		}
	}
}
