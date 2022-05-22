package main

import (
	"context"
  "os"
	"fmt"
  "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"time"
	"flag"
	"strings"
)

type Domain struct {
  IP  	string
	Name   string
}

type Domains struct {
  Domains []Domain
}

func main(){
	filename := flag.String("local-dns-config", "./test", "Location of local dns config file")
	mode := flag.String("run-mode", "interval", "Run mode (interval vs. once)")
	interval := flag.Int("check-interval", 30, "Interval in seconds for checking container labels")
	flag.Parse()

  checkFileExists(*filename)

	if *mode == "once" {
		fmt.Println()
		fmt.Println("Starting run")
		fmt.Println("---------------------")
		fmt.Println()
		writeLocalDNSFile(*filename)
		fmt.Println()
		fmt.Println("---------------------")
		fmt.Println("Stopping run")
		fmt.Println()
	} else if *mode == "interval" {
		for {
			fmt.Println()
			fmt.Println("Starting run")
			fmt.Println("---------------------")
			fmt.Println()
	    writeLocalDNSFile(*filename)
			fmt.Println()
			fmt.Println("---------------------")
			fmt.Println("Stopping run")
			fmt.Println()
			time.Sleep(time.Duration(*interval) * time.Second)
	  }
	} else {
		fmt.Println("Unknown run mode")
	}
}

func checkFileExists(filename string) error {
    _, err := os.Stat(filename)
        if os.IsNotExist(err) {
            _, err := os.Create(filename)
                if err != nil {
                    return err
                }
        }
        return nil
}


func writeLocalDNSFile(filename string){
	  containers := getContainers()

		domains:= getLabels(containers)

	  writeFile(filename, domains)

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

func getLabels(containers []types.Container) Domains {
	domains_empty := []Domain{}
	domains := Domains{domains_empty}

	for _, container := range containers {
		fmt.Println("---")
		fmt.Println("Found container with name:", strings.Trim(fmt.Sprint(container.Names), "/[]"))
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
		if (Domain{}) != domain  {
				domains.AddItem(domain)
		} else {
			fmt.Println("Container has no pihole labels")
		}
		fmt.Println("---")
	}

	return domains
}

func (domains *Domains) AddItem(domain Domain) []Domain {
	domains.Domains = append(domains.Domains, domain)
	return domains.Domains
}

func writeFile(filename string, domains Domains){
// this does not write to the file, but only to the temp file currently
	f, err := os.CreateTemp("./", "example")
		if err != nil {
			panic(err)
		}
		defer os.Remove(f.Name()) // defer clean up

		for _, domain := range domains.Domains {

			f.WriteString(domain.IP + " " + domain.Name + "\n")

			if err != nil {
		    panic(err)
		  }

		}



}
