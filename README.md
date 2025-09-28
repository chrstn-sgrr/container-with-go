# Container with Go

A self-learning experiment that recreates the core ideas behind containers (in Linux) in under 100 lines of Go.

## Inspiration
- [build-your-own-x](https://github.com/codecrafters-io/build-your-own-x)
- [InfoQ: Building a Container from Scratch in Go](https://www.infoq.com/articles/build-a-container-golang/)

## What This Does
- Forks the running binary (`/proc/self/exe`) so the child can re-enter Go land
- Creates new UTS, PID, and mount namespaces for process isolation
- Rebinds and pivots the root filesystem into a lightweight `rootfs`
- Executes the user-provided command inside this isolated environment

## What This Doesn't Have (Limitations)
- No cgroups or resource quotas
- Minimal error handling (intentionally kept simple)
- Assumes a pre-populated root filesystem

## Requirements
- Linux host (or WSL2/VM). The code relies on Linux-specific namespace syscalls.
- Go 1.21+

## Quick Start (Linux)
```bash
# Build the binary
go build -o gontainer

# Prepare a minimal root filesystem
mkdir -p rootfs/{bin,lib,lib64}
# Copy any binaries you plan to run into rootfs (e.g., busybox or sh)

# Run a command inside the container
sudo ./gontainer run /bin/sh
```

> The namespaces won't work on native Windows, use a Linux VM.

## How to Read the Code
- `main.go` orchestrates the parent (`run`) and child (`child`) roles.
- `parent()` configures namespaces via `SysProcAttr` and launches the child process.
- `child()` mounts the new root, performs `pivot_root`, and finally execs the requested command.
- `must()` is a small helper that panics on unexpected errors—useful for educational prototypes.

