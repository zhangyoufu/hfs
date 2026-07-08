// +build aix || darwin || dragonfly || freebsd || illumos || js || linux || netbsd || openbsd || solaris || wasip1 || zos

package main

import "syscall"

const umaskAvailable = true

func umask(mask int) (oldmask int) {
	return syscall.Umask(mask)
}
