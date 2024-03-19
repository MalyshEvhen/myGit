package object

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func StoreFromFile(name, typ string, dryRun bool) (Hash, error) {
	if !IsSupportedType(typ) {
		return Hash{}, fmt.Errorf("invalid object type: %s", typ)
	}

	content, err := os.ReadFile(name)
	if err != nil {
		return Hash{}, err
	}

	size := int64(len(content))
	reader := bytes.NewReader(content)

	return Store(reader, typ, size, dryRun)
}

func Store(r io.Reader, typ string, size int64, dryRun bool) (Hash, error) {
	if !IsSupportedType(typ) {
		return Hash{}, fmt.Errorf("invalid object type: %s", typ)
	}

	var buf bytes.Buffer

	if _, err := EncodeObject(&buf, r, typ, size); err != nil {
		return Hash{}, err
	}

	var fileContent bytes.Buffer
	err := Compress(&fileContent, buf.Bytes())
	if err != nil {
		return Hash{}, err
	}

	sum := sha1.Sum(buf.Bytes())
	name := hex.EncodeToString(sum[:])

	if !dryRun {
		err = writeFile(name, fileContent)
		if err != nil {
			return Hash{}, err
		}
	}
	return Hash(sum), nil
}

func writeFile(name string, fileContent bytes.Buffer) error {
	objPath := filepath.Join(".git", "objects", name[:2], name[2:])
	dirPath := filepath.Dir(objPath)

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	if err := os.WriteFile(objPath, fileContent.Bytes(), 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

func EncodeObject(dst io.Writer, src io.Reader, typ string, size int64) (int64, error) {
	if !IsSupportedType(typ) {
		return 0, fmt.Errorf("invalid object type: %s", typ)
	}

	if size < 0 {
		return 0, fmt.Errorf("invalid size: %d", size)
	}

	_, encodeErr := fmt.Fprintf(dst, "%v %d\000", typ, size)
	if encodeErr != nil {
		return 0, encodeErr
	}

	var bytesWritten int64
	var err error
	if s, ok := src.(io.WriterTo); ok {
		bytesWritten, err = s.WriteTo(dst)
	} else {
		bytesWritten, err = io.Copy(dst, src)
	}
	if err != nil {
		return bytesWritten, err
	}

	if bytesWritten != size {
		return bytesWritten, fmt.Errorf("write size mismatch: %d != %d", bytesWritten, size)
	}

	return bytesWritten, nil
}

func Compress(w io.Writer, data []byte) error {
	if len(data) == 0 {
		w.Write([]byte{})
		return nil
	}

	zw := zlib.NewWriter(w)
	defer zw.Close()

	n, err := zw.Write(data)
	if err != nil {
		if err == io.EOF {
			log.Printf("Unexpected EOF from zlib writer")
		}
		return err
	}

	if n < len(data) {
		log.Printf("Did not write all data to zlib writer. Wrote %d bytes out of %d", n, len(data))
	}

	if err := zw.Flush(); err != nil {
		log.Printf("Error flushing zlib writer: %v", err)
		return err
	}

	return zw.Close()
}

func IsSupportedType(typ string) bool {
	return typ == "blob" || typ == "tree"
}
