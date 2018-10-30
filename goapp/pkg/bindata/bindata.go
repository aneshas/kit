// Code generated by go-bindata.
// sources:
// data/template/app/default.json
// data/template/service/service.json
// data/template/service/service.proto
// DO NOT EDIT!

package bindata

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

var _dataTemplateAppDefaultJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xaa\xe6\x52\x50\x50\x50\x50\x4a\x2c\x28\x50\xb2\x52\x80\x70\xc0\x02\x79\x89\xb9\xa9\x4a\x56\x0a\x4a\xd5\xd5\x7a\x7e\x89\xb9\xa9\xb5\xb5\x4a\x3a\x08\xc9\x12\x25\x2b\x05\x03\x24\x7e\x5e\x7e\x4a\x6a\xb1\x92\x95\x42\x34\x5c\x08\x04\xaa\x51\x78\x28\xa6\x26\xe7\xa6\x20\x99\x87\x6a\x2e\x8a\x70\xad\x0e\xb1\x46\x16\x14\xe5\x97\xe4\x53\xdb\xd0\xe2\xd4\xa2\xb2\xcc\xe4\x54\x6a\x1b\x9b\x5f\x90\x9a\x97\x58\x90\x49\x6d\x63\x83\x5c\x1d\x5d\x7c\x5d\xf5\x70\x87\xad\x21\xaa\xc1\x70\x5e\x2c\x17\x84\x5f\xcb\x05\x08\x00\x00\xff\xff\x8a\x34\x05\xe9\x0f\x02\x00\x00")

func dataTemplateAppDefaultJsonBytes() ([]byte, error) {
	return bindataRead(
		_dataTemplateAppDefaultJson,
		"data/template/app/default.json",
	)
}

func dataTemplateAppDefaultJson() (*asset, error) {
	bytes, err := dataTemplateAppDefaultJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/template/app/default.json", size: 527, mode: os.FileMode(420), modTime: time.Unix(1534254300, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dataTemplateServiceServiceJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x96\xdf\x6e\x83\x20\x14\xc6\xef\x7d\x0a\xc2\x75\x63\xb6\x5b\x1f\x62\x2f\xb0\x2c\x0d\x69\xcf\x1c\x49\xf9\x13\xa4\xdb\x85\xe1\xdd\x97\x3a\x54\x98\x88\x58\x17\xdb\x6d\xf6\x8e\xe3\xf1\xe3\xeb\xf9\x7e\x12\xea\x0c\x21\x84\x30\x91\x12\x17\xe8\x6b\xd1\x14\x38\x61\x80\x0b\x84\x73\xbc\xeb\x8b\x1a\x17\xe8\xc1\x59\x73\x71\x84\x0a\x17\xe8\xb9\x2b\x5d\x7e\xb5\xb7\xf2\xd4\x0e\xec\xe8\xe8\x8d\xe9\x4e\xe8\x8f\xef\x33\xd8\xaf\xae\xf3\x27\xc2\xc0\x98\xc0\xae\x53\xbb\x27\xba\x98\x76\x33\x70\xc5\x08\xe5\x79\x29\x22\x9e\x3c\x6f\x8f\x09\x7d\x0c\x34\xf1\x02\x8c\xab\x02\x93\x27\xa2\x1b\x33\x15\xa8\x77\x7a\x80\xbd\x35\x95\x6b\x79\xc2\x93\x2a\x26\xda\x31\xfe\xf4\x25\xf8\x64\xd8\xef\xf7\x99\x5d\x2a\x5f\x52\x09\x1d\x9a\xeb\x7f\x23\xac\x73\x95\x8f\x4d\x24\xe8\x71\x15\xd2\xac\xa5\x5f\xcb\x98\xfd\x1b\x1b\x65\x0e\x65\x77\x76\x98\x5d\xfc\x2c\xe5\x2b\xee\x73\xc6\x94\xb4\x22\xbc\x92\x42\xe9\xd4\x11\x45\x32\xec\xc5\x93\xb2\x4c\x77\x3b\x70\x3d\x33\xdb\xee\xfd\xc4\x8c\xbb\xfe\x59\x59\xf7\xbb\x7c\x50\x25\xf7\x6e\xf2\xb6\xd2\xce\x7a\x9e\xeb\x46\xf3\x15\x88\x3e\xab\x6f\x91\x25\x4b\xc4\x69\x4a\xef\x0a\x9f\x5e\xdd\xfb\x3f\x46\x25\x03\xad\xe8\xa1\xda\x98\x0c\xf4\x5f\xc7\xa4\x54\x82\x81\x7e\x83\x73\xe5\x81\xe9\x94\xed\xcc\x97\xb1\xd9\x06\xf7\x67\xc9\x3c\x89\xb2\xa4\xbc\xdc\xc8\x0c\xf4\x5f\x79\x5a\x3a\x38\xda\xe9\x2e\x63\xb0\x8d\xe8\xde\x18\xbc\xd9\x9d\x50\x48\xe0\x44\xd2\x55\xef\x84\x53\x37\xfb\x1b\xde\x07\x93\xbf\xdd\x75\xe3\xcc\xfc\xba\xc9\x4c\xf6\x19\x00\x00\xff\xff\xc8\xee\x58\x34\xff\x10\x00\x00")

func dataTemplateServiceServiceJsonBytes() ([]byte, error) {
	return bindataRead(
		_dataTemplateServiceServiceJson,
		"data/template/service/service.json",
	)
}

func dataTemplateServiceServiceJson() (*asset, error) {
	bytes, err := dataTemplateServiceServiceJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/template/service/service.json", size: 4351, mode: os.FileMode(420), modTime: time.Unix(1534163962, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dataTemplateServiceServiceProto = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x64\x92\xcd\x8e\xda\x30\x10\xc7\xef\x7e\x8a\x51\x4e\xcb\x81\xb8\xb4\xb7\x44\x1c\x2a\xad\xe8\xc7\x61\x8b\x16\x6e\x55\x15\x79\xcd\x60\xac\x4d\x3c\xc6\x9e\xc0\x22\x94\x77\xaf\x12\x43\x00\xad\x2f\x1e\xcf\xd7\xff\x37\x93\xc4\x93\x63\xf5\x01\x73\xc8\x7c\x20\xa6\x6f\x59\x29\x84\x57\xfa\x5d\x19\x84\xf3\x39\x5f\xf6\xce\x65\x7a\xbf\xa8\x06\xbb\x4e\x08\xdb\x78\x0a\x0c\x99\xb1\xbc\x6b\xdf\x72\x4d\x8d\x6c\x8e\x96\xdf\xe9\x28\x0d\x4d\x87\x36\xd3\x83\xaa\xed\x46\x31\x85\x28\x47\x33\x1f\x42\x59\x39\x36\x18\xde\x7a\x6a\xd0\x4d\xe3\x51\x19\x83\x41\x92\x67\x4b\x2e\x4a\xe5\x1c\xb1\x1a\xec\xb1\x4c\xa4\x20\x18\xaa\xae\x84\x73\xc8\xce\xe7\x3c\x81\xdd\x65\x3c\x99\xe0\x75\x6e\x14\xe3\x51\x9d\x52\xbd\xae\x0c\xba\xea\x22\x93\x5f\x64\x72\xf2\xe8\x94\xb7\x87\xaf\xd7\xc8\x04\xe6\x70\x16\x00\x00\xd6\x6d\xa9\xb8\xd8\xfd\x61\xcb\x35\x16\x83\xde\xb3\x8d\xbe\x56\xa7\x24\x0b\x2b\x0c\x07\xab\x11\xd6\x47\x1b\xbc\xfc\xbd\xfa\xf3\x02\xdf\x97\xbf\xe0\x99\x74\xdb\xa0\x4b\x43\x64\xe5\xd8\xe7\x80\x21\x5a\x72\x05\x64\xb3\xfc\xcb\xc5\xdf\xa5\x0b\x3f\x18\x83\x53\x75\xb5\x21\x1d\xef\xb5\xdb\x50\x17\x90\xed\x98\x7d\x2c\xa4\xbc\xdb\x3b\x93\x63\x92\x86\x94\xf7\x77\x12\x1b\x8c\x3a\xd8\x61\xc4\x02\xb2\x1f\x7d\x14\x94\xf7\xb5\xd5\x03\x0c\x5c\xb7\x17\x70\x8b\x01\x9d\xc6\x47\x8c\xa8\x77\xd8\x60\x2c\xe0\xe7\x7a\xbd\x5c\x95\xa2\x2b\x85\x90\x12\x3e\x0d\x1e\xd3\xe0\xe2\x72\x7f\x4e\x48\x03\x48\x09\x0b\x22\x08\x5e\xc3\xf0\xee\x8d\x05\xd1\xd3\x82\xe8\x15\xf7\x13\x08\xc8\x6d\x70\x11\x92\x23\xfa\x49\x29\xba\x41\x70\xa8\xc2\x7d\x8b\x91\x45\x83\x31\xf6\xc8\xa9\xe8\xd6\xf9\x15\xf7\xad\x0d\xb8\x49\xe0\x1c\xac\x33\xf0\xa6\x02\xcc\x61\x06\x7f\xc7\x85\x3c\xdd\xfe\xc1\xad\xc5\x7a\x73\xfb\xca\xd7\x93\x4a\x2b\x47\x5c\x61\xe3\xf9\x04\x05\x70\x68\x71\xcc\xe9\x06\xeb\xdf\x23\x5a\xf4\xe4\x22\x3e\xb2\x45\x7f\x83\x5b\xb1\xe2\x36\xc2\x3d\x5b\x4c\xae\x39\xcc\xfa\x56\xff\x03\x00\x00\xff\xff\xbf\x53\xb6\xf1\x7d\x03\x00\x00")

func dataTemplateServiceServiceProtoBytes() ([]byte, error) {
	return bindataRead(
		_dataTemplateServiceServiceProto,
		"data/template/service/service.proto",
	)
}

func dataTemplateServiceServiceProto() (*asset, error) {
	bytes, err := dataTemplateServiceServiceProtoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/template/service/service.proto", size: 893, mode: os.FileMode(420), modTime: time.Unix(1534163255, 0)}
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
	"data/template/app/default.json": dataTemplateAppDefaultJson,
	"data/template/service/service.json": dataTemplateServiceServiceJson,
	"data/template/service/service.proto": dataTemplateServiceServiceProto,
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
	"data": &bintree{nil, map[string]*bintree{
		"template": &bintree{nil, map[string]*bintree{
			"app": &bintree{nil, map[string]*bintree{
				"default.json": &bintree{dataTemplateAppDefaultJson, map[string]*bintree{}},
			}},
			"service": &bintree{nil, map[string]*bintree{
				"service.json": &bintree{dataTemplateServiceServiceJson, map[string]*bintree{}},
				"service.proto": &bintree{dataTemplateServiceServiceProto, map[string]*bintree{}},
			}},
		}},
	}},
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
