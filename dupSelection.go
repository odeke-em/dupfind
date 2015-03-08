package dupfinder

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

type pair struct {
	left  interface{}
	right interface{}
}

type sizeFileInfoIndex map[int64]map[string]*KeyFileInfo

func groupFilesBySize(p string, depth int) (*sizeFileInfoIndex, error) {
	kvChan, err := listFiles(p, depth)
	if err != nil || kvChan == nil {
		return nil, err
	}

	index := make(sizeFileInfoIndex)

	for kv := range kvChan {
		if kv == nil {
			continue
		}

		v := kv.Value
		vRef := *v
		size := vRef.Size()
		kvMap, ok := index[size]
		if !ok {
			kvMap = make(map[string]*KeyFileInfo)
			index[size] = kvMap
		}
		kvMap[kv.Key] = kv
	}

	return &index, nil
}

func sizeClashes(p1Index, p2Index sizeFileInfoIndex) []int64 {
	clashes := []int64{}
	for size, _ := range p1Index {
		_, ok := p2Index[size]
		if !ok {
			continue
		}
		clashes = append(clashes, size)
	}

	return clashes
}

func match(path1, path2 string) (*map[int64]*pair, error) {
	p1IndexPtr, p1Err := groupFilesBySize(path1, -1)
	if p1Err != nil {
		return nil, p1Err
	}

	p2IndexPtr, p2Err := groupFilesBySize(path2, -1)
	if p2Err != nil {
		return nil, p2Err
	}

	p1Index, p2Index := *p1IndexPtr, *p2IndexPtr
	clashes := sizeClashes(p1Index, p2Index)

	clashMap := map[int64]*pair{}

	for _, size := range clashes {
		p := pair{
			left:  p1Index[size],
			right: p2Index[size],
		}
		left, _ := p1Index[size]
		right, _ := p2Index[size]
		fmt.Printf("\nsize: %d left: %v right: %v\n\n", size, left, right)
		clashMap[size] = &p
	}

	fmt.Println("clashes", clashes) // , clashMap)
	return nil, nil
}

func md5Checksum(p string) (string, error) {
	handle, err := os.Open(p)
	if err != nil || handle == nil {
		return "", err
	}

	defer handle.Close()
	hash := md5.New()
	_, err = io.Copy(hash, handle)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
