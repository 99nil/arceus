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

package quick

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/99nil/arceus/global"
	"github.com/99nil/arceus/pkg/types"
	"github.com/99nil/gopkg/ctr"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/tidwall/gjson"
	"sigs.k8s.io/yaml"
)

func list() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := types.BuildResourceList(os.DirFS(global.ResourcePath), "rule")
		if err != nil {
			ctr.InternalError(w, err)
			return
		}
		ctr.OK(w, result)
	}
}

func quickstart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data types.QuickStart
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			ctr.BadRequest(w, err)
			return
		}
		result, err := Parse(&data)
		if err != nil {
			ctr.InternalError(w, err)
			return
		}
		ctr.OK(w, result)
	}
}

// 根据规则名称，获取所有规则资源
func getRules(data *types.QuickStart) ([]types.QuickStartRule, error) {
	rules := make([]types.QuickStartRule, 0, len(data.Spec.Rule))
	for _, r := range data.Spec.Rule {
		filePath := filepath.Join(r.Group, r.Name, r.Version) + ".yaml"
		fileData, err := fs.ReadFile(os.DirFS(global.RuleResourcePath), filePath)
		if err != nil {
			return nil, err
		}
		var rule types.QuickStartRule
		if err := yaml.Unmarshal(fileData, &rule); err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func dataPatch(data []byte, pairs []types.KValuePair) ([]byte, error) {
	for _, pair := range pairs {
		p := strings.ReplaceAll(pair.Key, ".", "/")
		ops := []types.JSONOperation{
			{
				Op:    "replace",
				Path:  "/" + p,
				Value: pair.Value,
			},
		}
		marshal, err := json.Marshal(&ops)
		if err != nil {
			return nil, err
		}
		patch, err := jsonpatch.DecodePatch(marshal)
		if err != nil {
			return nil, err
		}
		data, err = patch.Apply(data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func Parse(data *types.QuickStart, pairs ...types.KValuePair) ([]interface{}, error) {
	rules, err := getRules(data)
	if err != nil {
		return nil, err
	}
	// 解析数据
	jsonData, err := yaml.YAMLToJSON([]byte(data.Spec.Data))
	if err != nil {
		return nil, err
	}

	// 过滤数据
	jsonData, err = dataPatch(jsonData, pairs)
	if err != nil {
		return nil, err
	}

	var result []interface{}
	for _, rule := range rules {
		single, err := ParseSingle(jsonData, &rule)
		if err != nil {
			return nil, err
		}
		result = append(result, single...)
	}
	return result, nil
}

func ParseSingle(data []byte, rule *types.QuickStartRule) ([]interface{}, error) {
	jsonResult := gjson.ParseBytes(data)

	// 处理模板
	var count int
	templates := make(map[string]map[string][]byte)
	for _, v := range rule.Spec.Templates {
		filePath := filepath.Join(v.Template.Group, v.Template.Name, v.Template.Version) + ".yaml"
		fileData, err := fs.ReadFile(os.DirFS(global.TemplateResourcePath), filePath)
		if err != nil {
			return nil, err
		}
		var template types.Template
		if err := yaml.Unmarshal(fileData, &template); err != nil {
			return nil, err
		}
		for _, temp := range template.Spec.Template {
			jsonData, err := yaml.YAMLToJSON([]byte(temp.Data))
			if err != nil {
				return nil, err
			}
			if templates[v.Name] == nil {
				templates[v.Name] = make(map[string][]byte)
			}
			templates[v.Name][temp.Name] = jsonData
			count++
		}
	}

	// 合并path定义
	defines := make(map[string]interface{})
	for _, v := range rule.Spec.Defines {
		src := make([]interface{}, 0, len(v.Src))
		for _, s := range v.Src {
			p := strings.TrimPrefix(s, "/")
			p = strings.ReplaceAll(p, "/", ".")
			val := jsonResult.Get(p).Value()
			src = append(src, val)
		}
		val := fmt.Sprintf(v.Value, src...)
		defines[v.Path] = val
	}

	// 处理规则
	for _, v := range rule.Spec.Settings {
		// 根据path获取值
		path := strings.TrimPrefix(v.Path, "/")
		path = strings.ReplaceAll(path, "/", ".")
		patchResult := jsonResult.Get(path)
		var patchValue interface{}
		if !patchResult.Exists() {
			var ok bool
			patchValue, ok = defines[v.Path]
			if !ok {
				continue
			}
		} else {
			patchValue = patchResult.Value()
		}

		// 根据target填充值
		for _, target := range v.Targets {
			templateData := templates[target.Name]
			patchs := make([]jsonpatch.Patch, 0, len(target.Fields))
			for _, field := range target.Fields {
				ops := []types.JSONOperation{
					{
						Op:    "replace",
						Path:  field.Path,
						Value: patchValue,
					},
				}
				if field.Op != "" {
					ops[0].Op = field.Op
				}
				switch pv := patchValue.(type) {
				case string:
					var opVal interface{}
					if err := json.Unmarshal([]byte(pv), &opVal); err == nil {
						ops[0].Value = opVal
					}
				}
				marshal, err := json.Marshal(&ops)
				if err != nil {
					return nil, err
				}
				patch, err := jsonpatch.DecodePatch(marshal)
				if err != nil {
					return nil, err
				}
				patchs = append(patchs, patch)
			}

			// 判断sub，如果不存在处理全部
			initialData := templateData
			if target.Sub != "" && templateData[target.Sub] != nil {
				initialData = map[string][]byte{
					target.Sub: templateData[target.Sub],
				}
			}
			for _, p := range patchs {
				for k, v := range initialData {
					patchData, err := p.Apply(v)
					if err != nil {
						continue
					}
					initialData[k] = patchData
					templates[target.Name][k] = patchData
				}
			}
		}
	}

	result := make([]interface{}, 0, count)
	for _, temp := range templates {
		for _, v := range temp {
			var item interface{}
			err := json.Unmarshal(v, &item)
			if err != nil {
				return nil, err
			}
			result = append(result, item)
		}
	}
	return result, nil
}
