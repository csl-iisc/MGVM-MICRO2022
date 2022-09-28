// Code generated by "esc -o esc.go -pkg lu -private kernels.hsaco"; DO NOT EDIT.

package lu

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
		size:    13720,
		modtime: 1601398835,
		compressed: `
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
`,
	},
}

var _escDirs = map[string][]os.FileInfo{}
