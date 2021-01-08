package main

import (
	"os"
)

func IsOwnerReadable(info os.FileInfo) bool {
	m := info.Mode()
	return m&(1<<8) != 0
}

func IsOwnerWritable(info os.FileInfo) bool {
	m := info.Mode()
	return m&(1<<7) != 0
}

func IsOwnerExecutable(info os.FileInfo) bool {
	m := info.Mode()
	return m&(1<<6) != 0
}

func IsGroupReadable(info os.FileInfo) bool {
	m := info.Mode()
	return m&(1<<5) != 0
}

func IsGroupWritable(info os.FileInfo) bool {
	m := info.Mode()
	return m&(1<<4) != 0
}

func IsGroupExecutable(info os.FileInfo) bool {
	m := info.Mode()
	return m&(1<<3) != 0
}

func IsOtherReadable(info os.FileInfo) bool {
	m := info.Mode()
	return m&(1<<2) != 0
}

func IsOtherWritable(info os.FileInfo) bool {
	m := info.Mode()
	return m&(1<<1) != 0
}

func IsOtherExecutable(info os.FileInfo) bool {
	m := info.Mode()
	return m&(1) != 0
}
