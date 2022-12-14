// Code generated by "esc -o bindata.go -pkg blackscholes -private kernels.hsaco"; DO NOT EDIT.

package blackscholes

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
		size:    17784,
		modtime: 1653557453,
		compressed: `
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
C9TVeEUAAA==
`,
	},
}

var _escDirs = map[string][]os.FileInfo{}
