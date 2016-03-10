// Package pkg provides common libraries for the Deis builder.
//
// The Deis builder is responsible for packaging Docker images for consumers.
package pkg

import (
	"fmt"

	"github.com/Masterminds/cookoo"
	clog "github.com/Masterminds/cookoo/log"
	"github.com/deis/sa-builder/pkg/sshd"

	"log"
	"os"
)

// Return codes that will be sent to the shell.
const (
	StatusOk = iota
	StatusLocalError
)

// Run starts the Builder service.
//
// The Builder service is responsible for setting up the local container
// environment and then listening for new builds. The main listening service
// is SSH. Builder listens for new Git commands and then sends those on to
// Git.
//
// Run returns on of the Status* status code constants.
func Run(sshHostIP string, sshHostPort int, cmd string) int {
	reg, router, ocxt := cookoo.Cookoo()
	log.SetFlags(0) // Time is captured elsewhere.

	// We layer the context to add better logging and also synchronize
	// access so that goroutines don't get into race conditions.
	cxt := cookoo.SyncContext(ocxt)
	cxt.Put("cookoo.Router", router)
	cxt.AddLogger("stdout", os.Stdout)

	// Build the routes. See routes.go.
	routes(reg)

	// Bootstrap the background services. If this fails, we stop.
	if err := router.HandleRequest("boot", cxt, false); err != nil {
		clog.Errf(cxt, "Fatal errror on boot: %s", err)
		return StatusLocalError
	}

	cxt.Put(sshd.Address, fmt.Sprintf("%s:%d", sshHostIP, sshHostPort))

	// Supply route names for handling various internal routing. While this
	// isn't necessary for Cookoo, it makes it easy for us to mock these
	// routes in tests. c.f. sshd/server.go
	cxt.Put("route.sshd.pubkeyAuth", "pubkeyAuth")
	cxt.Put("route.sshd.sshPing", "sshPing")
	cxt.Put("route.sshd.sshGitReceive", "sshGitReceive")

	// Start the SSH service.
	// TODO: We could refactor Serve to be a command, and then run this as
	// a route.
	if err := sshd.Serve(reg, router, cxt); err != nil {
		clog.Errf(cxt, "SSH server failed: %s", err)
		return StatusLocalError
	}

	return StatusOk
}
