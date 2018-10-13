// Code generated by go-bindata.
// sources:
// migrations/1_initial.down.sql
// migrations/1_initial.up.sql
// DO NOT EDIT!

package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var __1_initialDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\xf0\xf4\x73\x71\x8d\x50\xc8\x4c\x49\xcd\x2d\xc8\x2f\x49\xcd\x4b\xae\x8c\xcf\x4e\xad\x2c\x8e\xcf\x4c\xa9\xb0\xe6\x02\x2b\x08\x71\x74\xf2\x71\x55\x28\x48\xac\xcc\x4d\xcd\x2b\x29\xb6\xe6\x02\x04\x00\x00\xff\xff\xd9\x58\xf7\x5c\x36\x00\x00\x00")

func _1_initialDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__1_initialDownSql,
		"1_initial.down.sql",
	)
}

func _1_initialDownSql() (*asset, error) {
	bytes, err := _1_initialDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1_initial.down.sql", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var __1_initialUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x93\x4f\x73\xda\x3c\x10\xc6\xef\xfe\x14\xcf\x81\x77\x5e\xc8\x38\x9d\xf4\x4c\x73\x30\xce\x62\xd4\x38\x12\x95\xe5\x02\x27\xc6\x0d\x1a\xf0\x24\x91\x53\x5b\x29\xe5\xdb\x77\x6c\xe3\x3f\x98\x9e\xea\xe3\xee\xb3\xcf\x6a\x7f\xbb\xf6\x25\x79\x8a\xa0\xbc\x59\x48\x60\x73\x70\xa1\x40\x6b\x16\xa9\x08\xef\xc9\xe9\x4d\x1b\x5b\x60\xec\x00\x40\xba\x83\xa2\xb5\xc2\x52\xb2\x27\x4f\x6e\xf0\x48\x1b\xf7\x9c\xd0\x6f\xef\x99\xd5\xe6\xf9\xb4\x7d\xd1\xa7\x5a\x55\xfa\xf0\x38\x0c\x2b\x45\xf9\xdd\xde\x62\xa5\xb1\xcb\xcc\xff\x16\xc7\xc4\x58\xd8\x83\xc6\x47\xa1\x73\x68\x63\x75\x9e\x9a\x3d\x12\x2c\xe2\x80\xf0\xa2\x4f\x6d\x95\xbf\x20\xff\x11\xe3\xe7\x43\x92\x6f\x5f\xb5\xd9\xdb\xc3\x78\xd0\x6e\x82\x2f\xf7\xf8\x7c\x77\x37\xa9\x1f\xf3\x4b\xe7\x45\x9a\x19\xcc\x58\xc0\x78\xf7\x8c\x3a\x99\xe5\xfb\xc4\xa4\x45\x62\xd3\xcc\x6c\x9b\x79\x2e\x25\xff\xd4\x34\xb1\x36\x4f\x7f\x7c\x58\x5d\xe0\x6b\x24\xf8\xac\x9b\x7e\x32\x75\x9c\x33\xe2\x98\xb3\x6f\x31\x81\xf1\x07\x5a\x0f\x48\x0f\xdc\x8b\x6d\xba\xfb\x0d\xc1\x7b\x1b\x18\xf6\xef\x6c\x85\x84\xa4\x65\xe8\xf9\x84\x79\xcc\x7d\xc5\x04\x07\x33\x85\xce\xed\xb2\xae\x66\x4d\xa9\xad\xf7\x58\x0e\xde\x18\x7f\x4a\x77\xff\xa9\xcd\x92\xdc\x26\x73\xb9\xc9\x9e\xec\x22\xd1\xaf\x69\x80\xb7\xda\x73\xa0\xaf\x19\x72\x6f\xb5\x83\x44\xbf\xa6\xc7\xb4\x95\x77\xb1\x4a\xe9\x4c\x1c\x49\x2a\x96\xbc\x77\xab\x5e\x84\xd1\xa8\x72\x58\x31\xb5\x80\xd1\xc7\x6d\x9e\x1d\xcb\xf0\xb8\xdd\x2f\xe3\x11\x49\x05\xc6\x95\xb8\x40\xec\x0e\x17\xe1\x36\xe7\xe4\x0e\x4f\xc7\xed\xed\x7c\xd2\x1a\x7f\xf7\xc2\x98\x22\x8c\x2b\xc1\x90\xa6\xdb\xb2\x72\xaf\x88\xb8\xfd\x79\x3b\x3f\xc1\xe1\x0b\x3e\x0f\x99\xaf\xae\x4f\x00\x0f\xa2\xbc\xa1\x05\xe3\x41\x5b\x50\xe3\x60\x3c\xc0\x4d\x15\xab\xad\x22\x0a\xc9\x57\xb8\xc1\x5c\x8a\xa7\x06\x49\x95\x89\x39\x13\xfc\x2f\x9a\x96\xca\x6a\x41\x92\xae\xfe\xf0\xfb\xab\xe1\xa6\xce\x68\x84\xd0\xe3\x41\xec\x05\x84\xe2\xe7\xeb\xd4\xf9\x13\x00\x00\xff\xff\x53\xbc\xcd\xd2\x5a\x04\x00\x00")

func _1_initialUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__1_initialUpSql,
		"1_initial.up.sql",
	)
}

func _1_initialUpSql() (*asset, error) {
	bytes, err := _1_initialUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "1_initial.up.sql", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"1_initial.down.sql": _1_initialDownSql,
	"1_initial.up.sql":   _1_initialUpSql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"1_initial.down.sql": &bintree{_1_initialDownSql, map[string]*bintree{}},
	"1_initial.up.sql":   &bintree{_1_initialUpSql, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
