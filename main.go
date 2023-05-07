//go:build darwin

package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// setItermBadge sets a badge in iTerm2: https://iterm2.com/documentation-badges.html
func setItermBadge(msg string) {
	fmt.Printf("\033]1337;SetBadgeFormat=%s\a", base64.StdEncoding.EncodeToString([]byte(msg)))
}

// clearItermBadge clear the current iTerm2 badge.
func clearItermBadge() {
	fmt.Printf("\033]1337;SetBadgeFormat=\a")
}

func run(text string, ppid uint64) error {
	setItermBadge(text)

	// Set up a kqueue file descriptor.
	kq, err := syscall.Kqueue()
	if err != nil {
		return fmt.Errorf("failed kqueue(): %w", err)
	}

	// Define a kevent: wait for a pid to terminate and send an event.
	event := syscall.Kevent_t{
		Ident:  uint64(ppid),
		Filter: syscall.EVFILT_PROC,
		Flags:  syscall.EV_ADD | syscall.EV_ONESHOT,
		Fflags: syscall.NOTE_EXIT,
		Data:   0,
		Udata:  nil,
	}

	// wait for the kevent to trigger.
	events := make([]syscall.Kevent_t, 1)
	nev, err := syscall.Kevent(kq, []syscall.Kevent_t{event}, events, nil)
	if err != nil {
		return fmt.Errorf("failed kevent(): %w", err)
	}

	if nev < 1 {
		return errors.New("no events returned")
	}

	clearItermBadge()
	return nil
}

func main() {
	// exit early if Stdout is not a terminal; if we don't do this, we risk of injecting
	// the iTerm2 escape sequence in other programs, like for example git called by magit.
	if !term.IsTerminal(syscall.Stdout) {
		return
	}

	var pidFlag int
	flag.IntVar(&pidFlag, "pid", 0, "pid of the parent OpenSSH process")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("error: you must specify a message to print in iTerm2")
		os.Exit(1)
	}

	badgeText := strings.Join(args, " ")

	// If pidFlag is 0 it means we are in the primary process; the primary process gets the
	// pid of its parent process (ssh) and call itself passing the "-pid" flag.
	// We're doing this because Go can't just call fork().
	if pidFlag == 0 {
		// the parent process is ssh.
		ppid := os.Getppid()

		cmd := exec.Command(os.Args[0], "-pid", strconv.Itoa(ppid), badgeText)
		// it's VERY important to link stdout, or it will be linked to /dev/null
		cmd.Stdout = os.Stdout
		// detach the child process from its parent
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		if err := cmd.Start(); err != nil {
			log.Fatalf("cannot spawn child: %s", err)
		}
		return
	}

	// this is the secondary process
	if err := run(badgeText, uint64(pidFlag)); err != nil {
		// nobody will ever see these log messages...
		log.Fatalf("fatal error: %s", err)
	}
}
