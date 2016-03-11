package gitreceive

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/deis/pkg/log"

	client "k8s.io/kubernetes/pkg/client/unversioned"
)

func readLine(line string) (string, string, string, error) {
	spl := strings.Split(line, " ")
	if len(spl) != 3 {
		return "", "", "", fmt.Errorf("malformed line [%s]", line)
	}
	return spl[0], spl[1], spl[2], nil
}

func Run(conf *Config) error {
	log.Debug("Running git hook")

	kubeClient, err := client.NewInCluster()
	if err != nil {
		return fmt.Errorf("couldn't reach the api server (%s)", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		oldRev, newRev, refName, err := readLine(line)

		if err != nil {
			return fmt.Errorf("reading STDIN (%s)", err)
		}

		log.Debug("read [%s,%s,%s]", oldRev, newRev, refName)

		// if we're processing a receive-pack on an existing repo, run a build
		if strings.HasPrefix(conf.SSHOriginalCommand, "git-receive-pack") {
			if err := build(conf, kubeClient, newRev); err != nil {
				return err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
