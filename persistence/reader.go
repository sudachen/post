package persistence

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/spacemeshos/go-spacemesh/log"
	"github.com/spacemeshos/merkle-tree/cache"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type (
	LayerReadWriter = cache.LayerReadWriter
)

// NewLabelsReader returns a new labels reader from the initialization files.
// If the initialization was split into multiple files, they will be grouped
// into one unified reader.
func NewLabelsReader(id []byte) (LayerReadWriter, error) {
	if len(id) > 64 {
		err := fmt.Errorf("failed to create labels group reader: id cannot be longer than 64 bytes (got %d bytes)",
			len(id))
		return nil, err
	}

	dir := filepath.Join(GetPostDataPath(), hex.EncodeToString(id))
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("initialization directory not found: %v", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("initialization directory (%v) is empty", dir)
	}
	sort.Sort(numericalSorter(files))

	readers := make([]LayerReadWriter, len(files))
	for i, file := range files {
		reader, err := newReader(filepath.Join(dir, file.Name()), LabelGroupSize)
		if err != nil {
			return nil, err
		}
		readers[i] = reader
	}

	if len(readers) == 1 {
		return readers[0], nil
	} else {
		reader, err := cache.Group(readers)
		if err != nil {
			return nil, err
		}
		return reader, nil
	}
}

type Reader struct {
	f        *os.File
	b        *bufio.Reader
	itemSize uint64
}

func newReader(name string, itemSize uint64) (*Reader, error) {
	f, err := os.OpenFile(name, os.O_RDONLY, OwnerReadWrite)
	if err != nil {
		log.Error("failed to open file for labels reader: %v", err)
		return nil, err
	}
	return &Reader{
		f:        f,
		b:        bufio.NewReader(f),
		itemSize: itemSize,
	}, nil
}

func (r *Reader) Seek(index uint64) error {
	_, err := r.f.Seek(int64(index*r.itemSize), io.SeekStart)
	if err != nil {
		log.Error("failed to seek in labels reader: %v", err)
		return err
	}
	r.b.Reset(r.f)
	return err
}

func (r *Reader) ReadNext() ([]byte, error) {
	ret := make([]byte, r.itemSize)
	_, err := r.b.Read(ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *Reader) Width() (uint64, error) {
	info, err := r.f.Stat()
	if err != nil {
		log.Error("failed to get stats for leaf reader: %v", err)
		return 0, err
	}
	return uint64(info.Size()) / r.itemSize, nil
}

func (r *Reader) Flush() error {
	return nil
}

func (r *Reader) Append(p []byte) (n int, err error) {
	return 0, errors.New("writing not permitted")
}

func (r *Reader) Close() error {
	r.b = nil
	return r.f.Close()
}

type numericalSorter []os.FileInfo

// A compile time check to ensure that widthReader fully implements LayerReadWriter.
var _ sort.Interface = (*numericalSorter)(nil)

func (s numericalSorter) Len() int      { return len(s) }
func (s numericalSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s numericalSorter) Less(i, j int) bool {
	pathA := s[i].Name()
	pathB := s[j].Name()

	// Get the integer values of each filename, placed after the delimiter.
	a, err1 := strconv.ParseInt(pathA[strings.Index(pathA, "-")+1:], 10, 64)
	b, err2 := strconv.ParseInt(pathB[strings.Index(pathB, "-")+1:], 10, 64)

	// If any were not numbers, sort lexicographically.
	if err1 != nil || err2 != nil {
		return pathA < pathB
	}

	return a < b
}