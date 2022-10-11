// Package jcalendar provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package jcalendar

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xY3W4bNxN9lQW/73JtOW56ozsnsQu1QWz4BwFqCAW9HMlMd0lmOCtVMPTuBcmVd1fi",
	"WlItGTbgm0QWxTPk8JzhGT6wTBdGK1BkWf+B2eweCu4/fkbgBOJ0Aorc3wa1ASQJfnTwxf1LMwOsz6Qi",
	"GAOy+TxdfKXvfkBGbJ4ucG4s4HNgThE1XoI1WllYBboiTqWNgaXsWlIOjSFLKNW4I4zb7iX8LMFGdi04",
	"cff//xFGrM/+16uz16tS1xsoU1LI2nxeI3YtfBPI85IimHY9qCQo7FbojynhiHwWop0hwHephJ6uRjpD",
	"XURS65Ke5xvmvMZ/Xpoa6/TrbpzECp5npcY4/1L2BYjL3EZ31rnlgb1AOeHU5Nqd1jlwFYYvwQCn+OgF",
	"R5KZNLxS4uPRrS6ufUJbpdqnJC7F04LLPL5hiZa+8QKio1/5E4PXsoA/tYLz0cgCbar1gZpIghsjOEGn",
	"Fgf2JMvAEIhYQmOwTaLHCQEnFN3GXtjSAbaGRR2HnbKQrugGYsm41n+DuuCzXHOxmg0/ug3S85TbWozX",
	"ruPo88uwZ3pk0e4rqUbagWRaEc9clHnKBNgMpSGpFeuz3z/zHJTgmJxcDFjKKFwk7GrKx2PAZHl8AmjD",
	"zA+HRy412oDiRrI+++Xw6PADS5nhdO/30IPJ4tI1OmyyHfxECJuEHyVurVhwP+JR0X8eCNZnF9oGTluW",
	"Mgwp+6TFbLG3iu3cmFxmflrvh9Wqvu3XZbJ1Jc597rA6bb/846OjncVqmQ4fqp2U8z9cWj8eH+9udy1f",
	"EQl5owzqDKzldzkkp4ok+ar76w63vXYRA0WAiufJFeAEMPETPLNtWRQcZ6xfmYIlrhAfW9a/ZRXbhm5K",
	"Rb3egxRzt7IxROh3CVSiqhj4JAF/g4p/A+EJjrwAAnRhl0EXdBagSI5crjQyJ0XW98pgKVP+KmFSVGyW",
	"6Ao8YQlpI5nLZWm4R062DVw3KV+SD5+4SB4l6WJ/fLnY3zQlZ7pU4vWqYBsRlBbwr02U0MBM+J0uKQEF",
	"ZJORxsQayByhE4fWrRB3H61XicPYRCPVyrcSSroc7A7GUiV6lDgLghOeL2L9LAFndbCR8zHPigRKbBCH",
	"nL95VdK379p/I9pv3VltwTqRel1NJd0ntXKipUH6/qO+IE0ZKQvBb7so7seJDQ8QK+6spNDNbHA7Vkh7",
	"uh13bwxjbdq7P3y1/tAflFTjiqqhFrsTbMig+iLIINdjqXbToHz1UJvS8J+D6XR64BAPSsxBZVqELr/O",
	"VrshhM63C8OtnWoU8V62ltBthdGYMYy2jW3V7ZPr7b66k+xvuA9xZXhdB7zquqIWy5HsxsPtp9Y13yNe",
	"psZVbxfvJW4L6nnGBPscZZ+nTuDe1D8U2//g+kcIkITp3lYUAK6sxlz/9yrIBoZ/qSu2Hfa4ci72dvik",
	"A1h+Pa4L4tLz/j7dc+Rd/91Bvw0HneXa+t62k+sLSS10NPQRrYcOHC8xZ33W40ay+XD+bwAAAP//Cuii",
	"KWocAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}