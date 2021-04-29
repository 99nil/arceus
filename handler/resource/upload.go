/*
Copyright © 2021 zc2638 <zc2638@qq.com>.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package resource

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/zc2638/arceus/pkg/util"

	"github.com/pkgms/go/ctr"
	"github.com/tidwall/gjson"
	apiextensionsV1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"

	"github.com/zc2638/arceus/global"
	"github.com/zc2638/arceus/pkg/types"
)

func upload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			ctr.BadRequest(w, err)
			return
		}
		defer file.Close()

		fileData, err := io.ReadAll(file)
		if err != nil {
			ctr.InternalError(w, err)
			return
		}
		arr := bytes.Split(fileData, []byte("---\n"))
		for _, v := range arr {
			vb := bytes.TrimSpace(v)
			if len(vb) == 0 {
				continue
			}
			if err := GenerateFile(vb, global.CustomResourcePath); err != nil {
				ctr.BadRequest(w, err)
				return
			}
		}
		ctr.Success(w)
	}
}

func GenerateFile(source []byte, targetDir string) error {
	var data apiextensionsV1.CustomResourceDefinition
	if err := runtime.DecodeInto(scheme.Codecs.UniversalDecoder(), source, &data); err != nil {
		return fmt.Errorf("resource parse failed: %s", err)
	}
	if data.Spec.Group == "" {
		return fmt.Errorf("group not found")
	}
	if data.Spec.Names.Kind == "" {
		return fmt.Errorf("kind not found")
	}
	if len(data.Spec.Versions) == 0 {
		return fmt.Errorf("version not found")
	}
	dir := filepath.Join(targetDir, data.Spec.Group, data.Spec.Names.Kind)
	if _, err := os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}
	for k, v := range data.Spec.Versions {
		newFile, err := os.Create(filepath.Join(dir, v.Name+".yaml"))
		if err != nil {
			return fmt.Errorf("create file failed: %s", err)
		}
		current := data.DeepCopy()
		current.Spec.Versions = make([]apiextensionsV1.CustomResourceDefinitionVersion, 0, 1)
		current.Spec.Versions = append(current.Spec.Versions, data.Spec.Versions[k])
		// TODO 转为json格式节省空间
		b, err := yaml.Marshal(current)
		if err != nil {
			return fmt.Errorf("generate resource failed: %s", err)
		}
		if _, err := newFile.Write(b); err != nil {
			return fmt.Errorf("save file failed: %s", err)
		}
		newFile.Close()
	}
	return nil
}

func generate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			ctr.BadRequest(w, err)
			return
		}
		defer file.Close()

		fileData, err := io.ReadAll(file)
		if err != nil {
			ctr.InternalError(w, err)
			return
		}
		jsonData, err := yaml.YAMLToJSON(fileData)
		if err != nil {
			ctr.BadRequest(w, err)
			return
		}

		data := types.ArceusResourceDefinition{
			TypeMeta:   types.TypeMeta{},
			ObjectMeta: types.ObjectMeta{},
			Spec: types.ArceusResourceDefinitionSpec{
				Group: "custom.arceus",
				Names: types.ArceusResourceDefinitionNames{
					Kind: "kind-" + util.RandomStr(6),
				},
			},
		}
		result := gjson.ParseBytes(jsonData)
		jsonSchema := dealSchema(result)
		version := types.ArceusResourceDefinitionVersion{}
		version.Name = "v1"
		version.Schema = &types.ArceusResourceValidation{
			OpenAPIV3Schema: jsonSchema,
		}
		data.Spec.Versions = []types.ArceusResourceDefinitionVersion{version}
		// 转yaml
		b, err := yaml.Marshal(&data)
		if err != nil {
			ctr.InternalError(w, err)
			return
		}
		ctr.Bytes(w, b)
	}
}

func dealSchema(data gjson.Result) *types.JSONSchemaProps {
	props := &types.JSONSchemaProps{}
	if data.IsArray() {
		// array handle
		props.Type = TypeArray
		arr := data.Array()
		if len(arr) == 0 {
			return props
		}
		set := make(map[string]types.JSONSchemaProps)
		var itemProps types.JSONSchemaProps
		for _, v := range arr {
			current := dealSchema(v)
			if itemProps.Type == "" {
				itemProps.Type = current.Type
			}
			if itemProps.Type != current.Type {
				continue
			}
			if current.Type != TypeObject {
				itemProps = *dealSchema(v)
				break
			}
			for ik, iv := range current.Properties {
				set[ik] = iv
			}
		}
		if itemProps.Type == TypeObject {
			itemProps.Properties = set
		}
		props.Items = &itemProps
	} else if data.IsObject() {
		// object handle
		props.Type = TypeObject
		obj := data.Map()
		props.Properties = make(map[string]types.JSONSchemaProps)
		props.Required = make([]string, 0, len(obj))
		for k, v := range obj {
			props.Properties[k] = *dealSchema(v)
			props.Required = append(props.Required, k)
		}
	} else {
		switch data.Type {
		case gjson.String:
			props.Type = TypeString
		case gjson.Number:
			props.Type = TypeNumber
		case gjson.True, gjson.False:
			props.Type = TypeBoolean
		default:
			props.Type = TypeString
		}
		val := data.String()
		props.Default = &val
	}
	return props
}
