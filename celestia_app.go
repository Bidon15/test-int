package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
)

func main() {
	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	net, err := pool.CreateNetwork("localnet", func(config *dc.CreateNetworkOptions) {
		config.Driver = "bridge"
		config.IPAM = &dc.IPAMOptions{
			Driver: "default",
			Config: []dc.IPAMConfig{
				{
					Subnet: "192.168.10.0/24",
				},
			},
		}
	})

	if err != nil {
		log.Fatalf("Could not create network: %s", err)
	}
	defer net.Close()

	node0 := &dockertest.RunOptions{
		Name:      "node0",
		NetworkID: net.Network.ID,
		PortBindings: map[dc.Port][]dc.PortBinding{
			"26656/tcp": {{HostIP: "", HostPort: "26656"}},
			"26657/tcp": {{HostIP: "", HostPort: "26667"}},
		},
	}

	node1 := &dockertest.RunOptions{
		Name:      "node1",
		NetworkID: net.Network.ID,
		PortBindings: map[dc.Port][]dc.PortBinding{
			"26656/tcp": {{HostIP: "", HostPort: "26659"}},
			"26657/tcp": {{HostIP: "", HostPort: "26660"}},
		},
	}

	res, err := pool.BuildAndRunWithOptions("/home/nnv/go/src/github.com/celestiaorg/test-int/celestia-app/Dockerfile", node0)

	if err != nil {
		log.Fatalf("Could not start resource %s", err)
	}

	fmt.Println("First Resource IP", res.GetIPInNetwork(net))

	busybox, err := pool.BuildAndRunWithOptions("/home/nnv/go/src/github.com/celestiaorg/test-int/alpine/Dockerfile", &dockertest.RunOptions{
		Name:      "cli",
		NetworkID: net.Network.ID,
		Tty:       true,
	})

	if err != nil {
		log.Fatalf("Could not start 2nd resource %s", err)
	}
	fmt.Println("2nd Resource IP ", busybox.GetIPInNetwork(net))

	res3, err := pool.BuildAndRunWithOptions("/home/nnv/go/src/github.com/celestiaorg/test-int/celestia-app/Dockerfile", node1)

	if err != nil {
		log.Fatalf("Could not start resource %s", err)
	}

	fmt.Println("First Resource IP", res3.GetIPInNetwork(net))

	// curling node0 using busybox
	if err = pool.Retry(func() error {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		// notice that inside the same network, containers need to reach to the exported ports defined internally(like 26657), 
		// rather then defined for host's (e.g. res.GetPort("26657/tcp") == 26660). Still a question mark but works, so ok with that
		url0 := res.GetIPInNetwork(net) + ":" + "26657"
		fmt.Println(url0)
		exitCode, err := busybox.Exec(
			[]string{
				"curl", url0,
				// Uncomment to play around with ping, wget, etc.
				// "time","ping","-w2", url,
				// "wget", "-O", "-", res.GetIPInNetwork(net) + ":26657",
			},
			dockertest.ExecOptions{
				StdOut: &stdout,
				StdErr: &stderr,
			},
		)

		fmt.Println("Exit code ", exitCode)
		fmt.Println("Stdout ", strings.TrimRight(stdout.String(), "\n"))
		fmt.Println("Stderr ", strings.TrimRight(stderr.String(), "\n"))
		// retry method stops re-executing if one of 2 points are met:
		// 1. err is equal to nil
		// 2. timeout is reached (defined internally in dockertest libs)
		if exitCode != 0 {
			err = errors.New("Command failed. Retrying until deadline time")
		}

		return err

	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// curling node1 using busybox
	if err = pool.Retry(func() error {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		url1 := res3.GetIPInNetwork(net) + ":" + "26657"
		exitCode, err := busybox.Exec(
			[]string{
				"curl", url1,
			},
			dockertest.ExecOptions{
				StdOut: &stdout,
				StdErr: &stderr,
			},
		)

		fmt.Println("Exit code ", exitCode)
		fmt.Println("Stdout ", strings.TrimRight(stdout.String(), "\n"))
		fmt.Println("Stderr ", strings.TrimRight(stderr.String(), "\n"))
		if exitCode != 0 {
			err = errors.New("Command failed. Retrying until deadline time")
		}

		return err

	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// testing how http.get request should work (no busybox here)
	if err = pool.Retry(func() error {
		url1 := res3.GetIPInNetwork(net) + ":" + "26657"
		resp, err := http.Get(fmt.Sprintf("http://%s/health", url1))
		if err != nil {
			return err
		}
		fmt.Println(resp)

		return err

	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(res); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
