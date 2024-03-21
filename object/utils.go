package object

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func ReadFromFile(name string, typ Kind) (io.Reader, int64, error) {
	if !IsSupportedType(typ) {
		return nil, 0, fmt.Errorf("invalid object type: %s", typ)
	}

	content, err := os.ReadFile(name)
	if err != nil {
		return nil, 0, err
	}

	size := int64(len(content))
	reader := bytes.NewReader(content)

	return reader, size, nil
}

func Store(r io.Reader, kind Kind, size int64, dryRun bool) (Hash, error) {
	if !IsSupportedType(kind) {
		return Hash{}, fmt.Errorf("invalid object type: %s", kind)
	}

	var buf bytes.Buffer

	if _, err := EncodeObject(&buf, r, kind, size); err != nil {
		return Hash{}, err
	}

	var fileContent bytes.Buffer
	err := Compress(&fileContent, buf.Bytes())
	if err != nil {
		return Hash{}, err
	}

	sum := sha1.Sum(buf.Bytes())
	hash := Hash(sum)

	if !dryRun {
		err = WriteFile(hash, fileContent)
		if err != nil {
			return Hash{}, err
		}
	}
	return hash, nil
}

func LoadFileByHash(hash Hash) (io.ReadCloser, error) {
	path := MakeObjectPath(hash)

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	defer file.Close()

	return Decompress(file)
}

func WriteFile(hash Hash, content bytes.Buffer) error {
	objPath := MakeObjectPath(hash)
	dirPath := filepath.Dir(objPath)

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	if err := os.WriteFile(objPath, content.Bytes(), 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

func MakeObjectPath(hash Hash) string {
	name := hash.String()
	return filepath.Join(".git", "objects", name[:2], name[2:])
}

func EncodeObject(dst io.Writer, src io.Reader, kind Kind, size int64) (int64, error) {
	if !IsSupportedType(kind) {
		return 0, fmt.Errorf("invalid object type: %s", kind)
	}

	if size < 0 {
		return 0, fmt.Errorf("invalid size: %d", size)
	}

	_, err := fmt.Fprintf(dst, "%v %d\000", kind, size)
	if err != nil {
		return 0, err
	}

	var bytesWritten int64
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

func Decompress(r io.Reader) (io.ReadCloser, error) {
	zr, err := zlib.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("new zlib reader %w", err)
	}
	defer zr.Close()

	return zr, nil
}

func IsSupportedType(kind Kind) bool {
	return kind == "blob" || kind == "tree"
}
