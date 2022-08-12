// Copyright © 2021 zc2638 <zc2638@qq.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resource

import (
	"io"
	"net/http"

	"github.com/99nil/arceus/pkg/types"
	"github.com/99nil/arceus/pkg/util"
	"github.com/99nil/gopkg/ctr"

	"github.com/tidwall/gjson"
	"sigs.k8s.io/yaml"
)

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

		kind := "kind-" + util.RandomStr(6)
		data := types.ArceusResourceDefinition{
			TypeMeta: types.TypeMeta{
				APIVersion: types.Group + "/" + types.Version,
				Kind:       "CustomResourceDefinition",
			},
			ObjectMeta: types.ObjectMeta{
				Name: kind + "." + types.CustomGroup,
			},
			Spec: types.ArceusResourceDefinitionSpec{
				Group: types.CustomGroup,
				Names: types.ArceusResourceDefinitionNames{
					Kind: kind,
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
		props.Type = types.TypeArray
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
			if current.Type != types.TypeObject {
				itemProps = *dealSchema(v)
				break
			}
			for ik, iv := range current.Properties {
				set[ik] = iv
			}
		}
		if itemProps.Type == types.TypeObject {
			itemProps.Properties = set
		}
		props.Items = &itemProps
	} else if data.IsObject() {
		// object handle
		props.Type = types.TypeObject
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
			props.Type = types.TypeString
		case gjson.Number:
			props.Type = types.TypeNumber
		case gjson.True, gjson.False:
			props.Type = types.TypeBoolean
		default:
			props.Type = types.TypeString
		}
		val := data.String()
		props.Default = &val
	}
	return props
}
