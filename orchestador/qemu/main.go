package qemu

import (
	"os"
	"os/exec"
)

func Execute_Lab(
	name_iso string,
	path_iso string,
	ram_size_iso string,
	mode string) error {

	var args []string

	if mode == "display" {
		args = []string{"-nographic", "-m", ram_size_iso, "-cdrom", path_iso, "-boot", "d", "-enable-kvm"}
	} else {
		args = []string{"-display", "none", "-m", ram_size_iso, "-cdrom", path_iso, "-boot", "d", "-enable-kvm"}
	}

	cmd := exec.Command("qemu-system-x86_64", args...)

	if mode == "display" {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return cmd.Wait()
}
