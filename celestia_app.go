package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	// "net/http"

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

	ds := &dockertest.RunOptions{
		Name:         "celestia-app",
		NetworkID:    "localnet",
		Cmd:          []string{"--port", "26657"},
		ExposedPorts: []string{"26657"},
	}

	res, err := pool.BuildAndRunWithOptions("/tmp/test-int/celestia-app/Dockerfile", ds)

	// res, err := pool.BuildAndRun("celestia-app0", "/Users/bidon4/go/src/github.com/celestiaorg/test-int/celestia-app/Dockerfile", []string{
	// "--port", "1317:1317", "--port", "26656:26656", "--port", "26657:26657", "--port", "9090:9090"})
	if err != nil {
		log.Fatalf("Could not start resource %s", err)
	}

	fmt.Println(res.GetIPInNetwork(net))

	res2, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       "cli",
		Repository: "busybox",
		NetworkID:  "localnet",
		Tty:        true,
	})

	if err != nil {
		log.Fatalf("Could not start 2nd resource %s", err)
	}
	fmt.Println(res2.GetIPInNetwork(net))

	if err = pool.Retry(func() error {
		var stdout bytes.Buffer
		exitCode, err := res2.Exec(
			[]string{
				//"time", "ping", "-w2", res.GetIPInNetwork(net),
				"wget", "-O", "-", res.GetIPInNetwork(net) + ":26657",
			},
			dockertest.ExecOptions{
				//		TTY:    true,
				StdOut: &stdout,
			},
		)
		fmt.Println("Exit code ", exitCode)
		fmt.Println("Stdout ", stdout.String())
		// if we have res.GetIPinNetwork -> it's not connecting...
		// resp, err := http.Get(fmt.Sprintf("http://%s:%s/health", res.GetIPInNetwork(net), res.GetPort("26657/tcp")))
		// if err != nil {
		// 	return err
		// }
		// fmt.Println(resp)
		if exitCode != 0 && err == nil {
			fmt.Println("lol what?!")
		}
		if exitCode != 0 {
			err = errors.New("command failed")
		}
		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(res); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
