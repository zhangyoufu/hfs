// +build plan9 || windows

package main

const umaskAvailable = false

func umask(mask int) (oldmask int) {
	return 0
}
