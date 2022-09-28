// Code generated by "esc -o bindata.go -pkg md5sum -private kernels.hsaco"; DO NOT EDIT.

package md5sum

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
		size:    18424,
		modtime: 1648099862,
		compressed: `
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
///4S955+EcAAA==
`,
	},
}

var _escDirs = map[string][]os.FileInfo{}
