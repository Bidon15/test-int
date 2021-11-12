package helper

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/ory/dockertest/v3"
)

func SetupAccount(node *dockertest.Resource, name string) error {
	// Creating a node0
	// create new account
	// gentx ing (register and stake)
	// do steps a and b on node1, manually (like tendermint can't do this for you) send the generated genesis tx to node0
	// install node1's gentx into node0's directory
	// collect-gentxs ing
	// sending the genesis.json to node1
	// starting
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var exitCode int
	var err error
	exitCode, err = node.Exec(
		[]string{
			fmt.Sprintf("NODE_NAME=%s", name), "&&",
			"celestia-appd", "keys", "add", "$NODE_NAME", "--keyring-backend=$KEY_TYPE",
		},
		dockertest.ExecOptions{
			StdOut: &stdout,
			StdErr: &stderr,
		},
	)

	if exitCode != 0 {
		err = errors.New("can't create a new key")
	}

	fmt.Println("Stdout ", strings.TrimRight(stdout.String(), "\n"))
	fmt.Println("Stderr ", strings.TrimRight(stderr.String(), "\n"))

	return err
}

func AddAcount(node *dockertest.Resource, cointype string) error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var exitCode int
	var err error
	exitCode, err = node.Exec(
		[]string{
			"node_addr=$(celestia-appd", "keys", "show", "$NODE_NAME", "-a", "--keyring-backend", "$KEY_TYPE)", "&&",
			"celestia-appd", "add-genesis-account", "$node_addr", cointype, "--keyring-backend", "$KEY_TYPE",
		},
		dockertest.ExecOptions{
			StdOut: &stdout,
			StdErr: &stderr,
		},
	)

	if exitCode != 0 {
		err = errors.New("can't create a new key")
	}

	fmt.Println("Stdout ", strings.TrimRight(stdout.String(), "\n"))
	fmt.Println("Stderr ", strings.TrimRight(stderr.String(), "\n"))

	return err
}

func GenTx(node *dockertest.Resource, amount string) error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var exitCode int
	var err error
	exitCode, err = node.Exec(
		[]string{
			"celestia-appd", "gentx", "$NODE_NAME", amount, "--keyring-backend=$KEY_TYPE", "--chain-id", "$CHAIN_ID",
		},
		dockertest.ExecOptions{
			StdOut: &stdout,
			StdErr: &stderr,
		},
	)

	if exitCode != 0 {
		err = errors.New("can't create a new key")
	}

	fmt.Println("Stdout ", strings.TrimRight(stdout.String(), "\n"))
	fmt.Println("Stderr ", strings.TrimRight(stderr.String(), "\n"))

	return err
}

//send gentx-xxxxxx.json func
//pwd: /root/.celestia-app/config/gentx
// TODO: transfer files between nodes through cat > stdBytes > other resource
func GetFile(node *dockertest.Resource, file_path string) (bytes.Buffer, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var exitCode int
	var err error
	exitCode, err = node.Exec(
		[]string{
			"cat", file_path,
		},
		dockertest.ExecOptions{
			StdOut: &stdout,
			StdErr: &stderr,
		},
	)

	if exitCode != 0 {
		err = errors.New("can't create a new key")
	}

	fmt.Println("Stdout ", strings.TrimRight(stdout.String(), "\n"))
	fmt.Println("Stderr ", strings.TrimRight(stderr.String(), "\n"))

	return stdout, err
}

func SendFile(node *dockertest.Resource, content bytes.Buffer, file_path string) error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var exitCode int
	var err error
	exitCode, err = node.Exec(
		[]string{
			"echo", strings.TrimRight(content.String(), "\n"), ">>", file_path,
		},
		dockertest.ExecOptions{
			StdOut: &stdout,
			StdErr: &stderr,
		},
	)

	if exitCode != 0 {
		err = errors.New("can't create a new key")
	}

	fmt.Println("Stdout ", strings.TrimRight(stdout.String(), "\n"))
	fmt.Println("Stderr ", strings.TrimRight(stderr.String(), "\n"))

	return err
}

func RemoveFile(node *dockertest.Resource, file_path string) error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var exitCode int
	var err error
	exitCode, err = node.Exec(
		[]string{
			"rm", file_path,
		},
		dockertest.ExecOptions{
			StdOut: &stdout,
			StdErr: &stderr,
		},
	)

	if exitCode != 0 {
		err = errors.New("can't create a new key")
	}

	fmt.Println("Stdout ", strings.TrimRight(stdout.String(), "\n"))
	fmt.Println("Stderr ", strings.TrimRight(stderr.String(), "\n"))

	return err
}
