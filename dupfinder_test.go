package dupfinder

import (
	"fmt"
	"testing"
)

func TestListing(t *testing.T) {
	maxDepth := 2
	filesChan, err := listFiles("X", maxDepth)
	if err == nil {
		t.Errorf("expected an error trying to list 'X'")
	}
	if filesChan != nil {
		t.Errorf("not expecting a fileChan back")
	}

	filesChan, err = listFiles("/", maxDepth)
	if err != nil {
		t.Errorf("expecting no errors got %v", err)
	}

	if filesChan == nil {
		t.Errorf("expecting a channel of listings")
	}

	for kv := range filesChan {
		if kv.Key == "" {
			t.Errorf("expected a non empty path")
		}
		// fmt.Println(kv.Key)
	}

	dirsChan, dirErr := listDirs("/", maxDepth)
	if dirErr != nil {
		t.Errorf("expecting no errors got %v", dirErr)
	}

	if dirsChan == nil {
		t.Errorf("expecting a channel of listings")
	}

	for kv := range dirsChan {
		if kv.Key == "" {
			t.Errorf("expected a non empty path")
		}
		// fmt.Println(kv.Key)
	}
}

func TestGroupFilesBySize(t *testing.T) {
	var p = "/"
	_, err := groupFilesBySize(p, 2)
	if err != nil {
		t.Errorf("expecting successful discovery instead got %v", err)
	}
}

func TestClashMatch(t *testing.T) {
	_, err := match(".", "..")
	fmt.Println("err", err)
}
