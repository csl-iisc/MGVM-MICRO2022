// Code generated by "esc -private -o esc.go -pkg mergesort ./kernel.hsaco"; DO NOT EDIT.

package mergesort

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	if !f.isDir {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is not directory", f.name)
	}

	fis, ok := _escDirs[f.local]
	if !ok {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is directory, but we have no info about content of this dir, local=%s", f.name, f.local)
	}
	limit := count
	if count <= 0 || limit > len(fis) {
		limit = len(fis)
	}

	if len(fis) == 0 && count > 0 {
		return nil, io.EOF
	}

	return fis[0:limit], nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// _escFS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func _escFS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// _escDir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func _escDir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// _escFSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func _escFSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// _escFSMustByte is the same as _escFSByte, but panics if name is not present.
func _escFSMustByte(useLocal bool, name string) []byte {
	b, err := _escFSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// _escFSString is the string version of _escFSByte.
func _escFSString(useLocal bool, name string) (string, error) {
	b, err := _escFSByte(useLocal, name)
	return string(b), err
}

// _escFSMustString is the string version of _escFSMustByte.
func _escFSMustString(useLocal bool, name string) string {
	return string(_escFSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/kernel.hsaco": {
		name:    "kernel.hsaco",
		local:   "./kernel.hsaco",
		size:    18216,
		modtime: 1648625292,
		compressed: `
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
`,
	},
}

var _escDirs = map[string][]os.FileInfo{}