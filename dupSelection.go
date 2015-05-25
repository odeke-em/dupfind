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
type checksumFileInfoIndex map[string][]*KeyFileInfo

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

var GroupFilesBySize = groupFilesBySize

func groupFilesByChecksum(p string, depth int) (*checksumFileInfoIndex, error) {
	sfIndexPtr, err := groupFilesBySize(p, depth)
	if err != nil {
		return nil, err
	}

	sfIndex := *sfIndexPtr
	mapping := make(checksumFileInfoIndex)

	for _, bucket := range sfIndex {
		clashes := checksumBucket(bucket)
		for cksum, ckmap := range clashes {
			retr, rOk := mapping[cksum]
			if !rOk {
				retr = []*KeyFileInfo{}
			}

			for _, kf := range ckmap {
				retr = append(retr, kf)
			}
			mapping[cksum] = retr
		}
	}

	return &mapping, nil
}

var GroupFilesByChecksum = groupFilesByChecksum

func checksumBucket(bucket map[string]*KeyFileInfo) map[string]map[string]*KeyFileInfo {
	waits := uint64(0)
	waitChan := make(chan *pair)

	for p, kf := range bucket {
		if kf == nil || kf.Value == nil {
			continue
		}

		waits += 1
		go func(pp string, kff *KeyFileInfo) {
			cksum, ckErr := md5Checksum(pp)
			var pr *pair
			if ckErr == nil {
				pr = &pair{
					left: cksum,
					right: pair{
						left:  pp,
						right: kff,
					},
				}
			}

			waitChan <- pr
		}(p, kf)
	}

	kfmap := make(map[string]map[string]*KeyFileInfo)

	for i := uint64(0); i < waits; i += 1 {
		pr := <-waitChan
		if pr == nil {
			continue
		}

		cksum, cOk := pr.left.(string)
		if !cOk {
			continue
		}

		ppVal, pOk := pr.right.(pair)
		if !pOk {
			continue
		}

		pp, ppOk := ppVal.left.(string)
		if !ppOk {
			continue
		}

		kff, kOk := ppVal.right.(*KeyFileInfo)
		if !kOk {
			continue
		}

		retr, rOk := kfmap[cksum]
		if !rOk {
			retr = make(map[string]*KeyFileInfo)
		}

		retr[pp] = kff
		kfmap[cksum] = retr
	}

	return kfmap
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
	// TODO: skip if dir

	defer handle.Close()
	hash := md5.New()
	_, err = io.Copy(hash, handle)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
