package dupfinder

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

func match(path1, path2 string) error {
	p1IndexPtr, p1Err := groupFilesBySize(path1, -1)
	if p1Err != nil {
		return p1Err
	}

	p2IndexPtr, p2Err := groupFilesBySize(path2, -1)
	if p2Err != nil {
		return p2Err
	}

	p1Index, p2Index := *p1IndexPtr, *p2IndexPtr
	return nil
}
