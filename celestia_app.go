package main

import (
	"fmt"
	"log"
	"net/http"

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
					Subnet: "192.168.10.0/20",
				},
			},
		}
	})

	if err != nil {
		log.Fatalf("Could not create network: %s", err)
	}
	defer net.Close()

	ds := &dockertest.RunOptions{
		Name: "celestia-app",
		Networks: []*dockertest.Network{
			net,
		},
		Cmd:          []string{"--port", "26657"},
		ExposedPorts: []string{"26657"},
	}

	res, err := pool.BuildAndRunWithOptions("/Users/bidon4/go/src/github.com/celestiaorg/test-int/celestia-app/Dockerfile", ds)

	// res, err := pool.BuildAndRun("celestia-app0", "/Users/bidon4/go/src/github.com/celestiaorg/test-int/celestia-app/Dockerfile", []string{
	// "--port", "1317:1317", "--port", "26656:26656", "--port", "26657:26657", "--port", "9090:9090"})
	if err != nil {
		log.Fatalf("Could not start resource %s", err)
	}

	fmt.Println(res.GetIPInNetwork(net))

	if err = pool.Retry(func() error {

		// if we have res.GetIPinNetwork -> it's not connecting...
		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s/health", res.GetPort("26657/tcp")))
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
