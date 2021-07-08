package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

// go run main.go run <container> <cmd> <args>
func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("help")
	}
}

func run() {
	fmt.Printf("Running parent as %d\n", os.Getpid())

	// https://github.com/opencontainers/runc/blob/master/libcontainer/factory_linux.go#L199
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// https://man7.org/linux/man-pages/man2/unshare.2.html
	// https://github.com/opencontainers/runc/blob/master/libcontainer/configs/namespaces_syscall.go
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	must(cmd.Run())
}

func child() {
	fmt.Printf("Running %v as PID %d\n", os.Args[3:], os.Getpid())

	cg()

	cmd := exec.Command(os.Args[3], os.Args[4:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//https://github.com/opencontainers/runc/blob/master/libcontainer/standard_init_linux.go#L114
	must(syscall.Sethostname([]byte("HoraDeK8S")))

	// https://github.com/opencontainers/runc/blob/master/libcontainer/rootfs_linux.go#L942
	must(syscall.Chroot("/home/vagrant/containers/" + os.Args[2]))
	must(os.Chdir("/"))

	// https://github.com/opencontainers/runc/blob/master/libcontainer/rootfs_linux.go#L389
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	must(cmd.Run())

	must(syscall.Unmount("proc", 0))
}

// https://github.com/opencontainers/runc/blob/master/libcontainer/cgroups/utils.go
func cg() {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	os.Mkdir(filepath.Join(pids, "horadek8s"), 0755)
	must(ioutil.WriteFile(filepath.Join(pids, "horadek8s/pids.max"), []byte("20"), 0700))

	// Removes the new cgroup in place after the container exits
	must(ioutil.WriteFile(filepath.Join(pids, "horadek8s/notify_on_release"), []byte("1"), 0700))
	must(ioutil.WriteFile(filepath.Join(pids, "horadek8s/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
