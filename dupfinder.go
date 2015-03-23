package dupfinder

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func spliceChannels(from, to chan *KeyFileInfo) {
	for f := range from {
		to <- f
	}
}

func symlinkResolver(symPath string) (os.FileInfo, error) {
	resolvedPath, err := filepath.EvalSymlinks(symPath)
	if err != nil {
		return nil, err
	}
	return os.Stat(resolvedPath)
}

func list(p string, depth int, pass checker) (chan *KeyFileInfo, error) {
	localInfo, err := os.Stat(p)
	if err != nil || localInfo == nil {
		return nil, err
	}

	kvChan := make(chan *KeyFileInfo)
	go func(fCh *chan *KeyFileInfo) {
		fChan := *fCh
		defer func() {
			close(fChan)
		}()

		if !localInfo.IsDir() {
			fChan <- &KeyFileInfo{Key: p, Value: &localInfo}
			return
		}

		if depth == 0 {
			return
		}

		if depth >= 1 {
			depth -= 1
		}

		listing, err := ioutil.ReadDir(p)
		if err != nil {
			return
		}

		for _, file := range listing {
			fullPath := filepath.Join(p, file.Name())
			if symlink(file) {
				resolvedFile, fErr := symlinkResolver(fullPath)
				if fErr != nil {
					continue
				}
				file = resolvedFile
			}

			isDir := file.IsDir()

			if pass(isDir) {
				fChan <- &KeyFileInfo{Key: fullPath, Value: &file}
			}

			if isDir {
				children, childErr := list(fullPath, depth, pass)
				if childErr != nil {
					continue
				}
				spliceChannels(children, fChan)
			}
		}
	}(&kvChan)

	return kvChan, nil
}

func listFiles(p string, depth int) (chan *KeyFileInfo, error) {
	return list(p, depth, file)
}

func listDirs(p string, depth int) (chan *KeyFileInfo, error) {
	return list(p, depth, dir)
}

func listAll(p string, depth int) (chan *KeyFileInfo, error) {
	return list(p, depth, all)
}
