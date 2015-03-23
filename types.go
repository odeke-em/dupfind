package dupfinder

import "os"

type KeyFileInfo struct {
	Key   string
	Value *os.FileInfo
}

type checker func(bool) bool

func all(isDir bool) bool {
	return true
}

func dir(isDir bool) bool {
	return isDir
}

func file(isDir bool) bool {
	return !isDir
}

func symlink(f os.FileInfo) bool {
	return (f.Mode() & os.ModeSymlink) == 0
}

type lister func(string, int) (chan *KeyFileInfo, error)
