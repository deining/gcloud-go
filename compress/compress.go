package compress

import (
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
)

type discardCloser struct {
	io.Writer
}

func (discardCloser) Close() error {
	return nil
}

var discardWriter = discardCloser{io.Discard}

func TextSum(h *hash.Hash) string {
	return fmt.Sprintf("%x", (*h).Sum(nil))
}

func hashCompressCopy(target io.WriteCloser, source io.Reader) (*hash.Hash, error) {
	h := sha256.New()
	mWriter := io.MultiWriter(target, h)
	zWriter := gzip.NewWriter(mWriter)
	defer zWriter.Close()
	if _, err := io.Copy(zWriter, source); err != nil {
		return nil, err
	} else {
		return &h, nil
	}
}

// compress file
func HashAndCompressFile(outFile, inFile string) (*hash.Hash, error) {
	if inF, err := os.Open(inFile); err != nil {
		panic(err)
	} else if outF, err := os.Create(outFile); err != nil {
		panic(err)
	} else if h, err := hashCompressCopy(outF, inF); err != nil {
		return nil, err
	} else {
		inF.Close()
		outF.Close()
		return h, nil
	}
}
