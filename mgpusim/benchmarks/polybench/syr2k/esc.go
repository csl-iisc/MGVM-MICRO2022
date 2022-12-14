// Code generated by "esc -o esc.go -pkg syr2k -private kernels.hsaco"; DO NOT EDIT.

package syr2k

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

	"/kernels.hsaco": {
		name:    "kernels.hsaco",
		local:   "kernels.hsaco",
		size:    9552,
		modtime: 1601880096,
		compressed: `
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
Heb/awAAAP//IUoQr1AlAAA=
`,
	},
}

var _escDirs = map[string][]os.FileInfo{}
