package asset

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

var _config_config_json = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4c\x8d\xc1\x0e\x82\x30\x10\x44\xef\x24\xfc\xc3\x66\xcf\xb6\xec\x52\x0d\xa5\x37\x69\xf9\x0e\x63\x42\x13\x0f\x12\x9a\x6d\xf5\x62\xf8\x77\x03\x1c\xf4\xf8\x32\x6f\x66\x3e\x75\x05\x80\x39\xca\x3b\x0a\x3a\xd8\x11\x00\x1f\xa5\xa4\x5b\x5a\xa4\xa0\x03\x74\x96\xac\xc1\x2d\x59\x4f\x87\x5e\x24\xde\xe7\xfc\xe7\x4f\x71\x5e\xf8\xc7\x00\xf8\x92\xe7\xd6\x95\x92\x93\x6b\x1a\x26\xd6\x2d\x91\xb6\x46\x5f\xb8\xe9\x46\xeb\x4d\x3b\x0e\xca\x87\x9e\xd4\xd9\x70\x50\xd7\x31\x04\x65\x7b\xf2\x9d\xf7\xed\x30\x30\xe9\x3c\x25\x3c\xd6\xd6\xfd\xba\xae\xd6\xba\xfa\x06\x00\x00\xff\xff\x98\x2e\xc7\x74\xaf\x00\x00\x00")

func config_config_json() ([]byte, error) {
	return bindata_read(
		_config_config_json,
		"config/config.json",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
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
var _bindata = map[string]func() ([]byte, error){
	"config/config.json": config_config_json,
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
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func     func() ([]byte, error)
	Children map[string]*_bintree_t
}

var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"config/config.json": &_bintree_t{config_config_json, map[string]*_bintree_t{}},
}}
