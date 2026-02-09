# Linux Containerization Experiments

Small Go programs to learn Linux container primitives step by step.

## Programs

- `without-constraints`: baseline process re-exec demo (no namespaces/chroot).
- `namespaces`: runs a child process in new UTS/PID/mount namespaces.
- `chroot`: runs a child process in new namespaces and a chroot root filesystem.

## Repository Layout

- `cmd/without-constraints/main.go`
- `cmd/namespaces/main.go`
- `cmd/chroot/main.go`
- `scripts/setup-chroot.sh` (downloads and extracts Alpine rootfs into `chroot/`)

## Setup Chroot Rootfs

From repository root:

```bash
./scripts/setup-chroot.sh
```

This creates:

- `chroot/alpine-minirootfs-3.23.3-x86_64/`

If you want to re-download and replace existing files:

```bash
./scripts/setup-chroot.sh --force
```

## Build

From repository root:

```bash
mkdir -p bin

go build -o bin/without-constraints ./cmd/without-constraints
go build -o bin/namespaces ./cmd/namespaces
go build -o bin/chroot ./cmd/chroot
```

## Run

### 1) Without constraints

```bash
./bin/without-constraints run /bin/sh
```

### 2) Namespaces

```bash
sudo ./bin/namespaces run /bin/sh
```

### 3) Chroot + namespaces

```bash
sudo ./bin/chroot run ./chroot /bin/sh
```

## Notes

- Namespace and chroot operations require Linux and elevated privileges.
- These are learning examples and not production-grade container runtime code.
