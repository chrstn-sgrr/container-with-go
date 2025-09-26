package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	// This program works in two phases:
	// 1. "run" - creates a new process with isolated namespaces (parent)
	// 2. "child" - sets up the container filesystem and runs the command

	switch os.Args[1] {
	case "run":
		parent() // Create container with namespaces
	case "child":
		child() // Set up container environment and run command
	default:
		panic("what should I do??") // Invalid usage
	}
}

func parent() {
	// Create a new process that runs this same program with "child" argument
	// /proc/self/exe points to the current executable
	// append([]string{"child"}, os.Args[2:]...) means: ["child", "command", "args..."]
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	// Set up Linux namespaces for isolation (like Docker does)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | // Isolate hostname/domain name
			syscall.CLONE_NEWPID | // Isolate process IDs (PID 1 inside container)
			syscall.CLONE_NEWNS, // Isolate mount points (filesystem view)
	}

	// Connect container's input/output to our terminal
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the child process with namespaces
	if err := cmd.Run(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}

func child() {
	// This function runs INSIDE the container with isolated namespaces

	// Step 1: Set up the container's filesystem root
	// Mount "rootfs" directory as a bind mount (makes it accessible as mount point)
	must(syscall.Mount("rootfs", "rootfs", "", syscall.MS_BIND, ""))

	// Step 2: Create a directory to store the old root filesystem
	must(os.MkdirAll("rootfs/oldrootfs", 0700))

	// Step 3: Switch root filesystems (like chroot but better)
	// PivotRoot swaps the root: "rootfs" becomes "/", old root goes to "rootfs/oldrootfs"
	must(syscall.PivotRoot("rootfs", "rootfs/oldrootfs"))

	// Step 4: Change to the new root directory
	must(os.Chdir("/"))

	// Step 5: Run the actual command inside the container
	// os.Args[2] = command to run, os.Args[3:] = arguments for that command
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command - this is what the user actually wants to run
	if err := cmd.Run(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}

// Helper function: crash the program if there's an error
// Used for operations that should never fail in our container setup
// "must" functions are common in Go, they crash or "panic" if there's an error
func must(err error) {
	if err != nil {
		panic(err) // Terminate + print error if something goes wrong
	}
}
