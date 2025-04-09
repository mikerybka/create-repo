package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mikerybka/util"
	"gopkg.in/yaml.v2"
)

func main() {
	err := createRepo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Hosts map[string]Host

type Host struct {
	User string `yaml:"user"`
}

func readGithubHostsConfig() (map[string]Host, error) {
	b, err := os.ReadFile(filepath.Join(util.HomeDir(), ".config/gh/hosts.yml"))
	if err != nil {
		return nil, err
	}
	hosts := map[string]Host{}
	err = yaml.Unmarshal(b, &hosts)
	if err != nil {
		return nil, err
	}
	return hosts, nil
}

func getGithubUser() (string, error) {
	hosts, err := readGithubHostsConfig()
	if err != nil {
		return "", err
	}
	if len(hosts) != 1 {
		panic("houston we have a problem")
	}
	host, ok := hosts["github.com"]
	if !ok {
		return "", fmt.Errorf("non-github host in ~/.config/gh/hosts.yml")
	}
	return host.User, nil
}

func createRepo(id string) error {
	// Get user id from gh config file
	ghUser, err := getGithubUser()
	if err != nil {
		return err
	}

	// gh repo create
	cmd := exec.Command("gh", "repo", "create", id, "--public", "--license", "gpl-3.0")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, out)
	}

	// gh repo clone
	cmd = exec.Command("gh", "repo", "clone", id)
	cmd.Dir = filepath.Join(util.HomeDir(), "src/github.com", ghUser)
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, out)
	}

	return nil
}
