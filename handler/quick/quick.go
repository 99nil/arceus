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
package quick

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/pkgms/go/ctr"
	"github.com/tidwall/gjson"
	"sigs.k8s.io/yaml"

	"github.com/zc2638/arceus/global"
	"github.com/zc2638/arceus/pkg/types"
)

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
	for _, name := range data.Spec.Rule {
		filePath := name + ".yaml"
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

func Parse(data *types.QuickStart) ([]interface{}, error) {
	rules, err := getRules(data)
	if err != nil {
		return nil, err
	}
	// 解析数据
	jsonData, err := yaml.YAMLToJSON([]byte(data.Spec.Data))
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
		filePath := filepath.Join(v.Group, v.Template, v.Version) + ".yaml"
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

	// 处理规则
	for _, v := range rule.Spec.Settings {
		// 根据path获取值
		path := strings.TrimPrefix(v.Path, "/")
		path = strings.ReplaceAll(path, "/", ".")
		patchValue := jsonResult.Get(path).Value()
		// 根据target填充值
		for _, target := range v.Targets {
			templateData := templates[target.Name]
			var ops []types.JSONOperation
			for _, field := range target.Fields {
				op := "replace"
				if field.Op != "" {
					op = field.Op
				}
				ops = append(ops, types.JSONOperation{
					Op:    op,
					Path:  field.Path,
					Value: patchValue,
				})
			}
			marshal, err := json.Marshal(&ops)
			if err != nil {
				return nil, err
			}
			patch, err := jsonpatch.DecodePatch(marshal)
			if err != nil {
				return nil, err
			}

			// 判断sub，如果不存在处理全部
			initialData := templateData
			if target.Sub != "" && templateData[target.Sub] != nil {
				initialData = map[string][]byte{
					target.Sub: templateData[target.Sub],
				}
			}
			for k, v := range initialData {
				patchData, err := patch.Apply(v)
				if err != nil {
					continue
				}
				templates[target.Name][k] = patchData
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
