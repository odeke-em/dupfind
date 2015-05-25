package main

import (
	"fmt"
	"github.com/odeke-em/dupfind"
	"os"
)

func cwd() string {
	cwd_, _ := os.Getwd()
	return cwd_
}

func main() {
	argc := len(os.Args)
	if argc < 2 {
		fmt.Fprintf(os.Stderr, "%s <path>\n", os.Args[0])
		return
	}

	spin := dupfinder.NewPlayable(10)
	done := make(chan bool)

	spin.Play()

	go func() {
		p := os.Args[1]

		if p == "" {
			p = cwd()
		}

		sfIndexPtr, err := dupfinder.GroupFilesByChecksum(p, 1)

		if err != nil {
			fmt.Printf("[err] %v got for %s\n", err, p)
			return
		}

		if sfIndexPtr == nil {
			fmt.Printf("failed to get index", p)
			return
		}

		sfIndex := *sfIndexPtr

		for cksum, buckets := range sfIndex {
			if len(buckets) < 2 {
				continue
			}

			fmt.Println("md5Checksum", cksum)
			for i, cks := range buckets {
				fmt.Println("\t", i, cks.Key)
			}
		}

		done <- true
	}()

	<-done
	spin.Stop()
}
