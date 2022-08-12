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
	"github.com/99nil/arceus/pkg/types"
)

var (
	rule = types.QuickStartRule{
		ObjectMeta: types.ObjectMeta{
			Name: "test",
		},
		Spec: types.QuickStartRuleSpec{
			Templates: []types.RuleTemplateDefine{
				{
					Name: "nginx",
					Template: types.RuleTemplateResourceDefine{
						Name:    "nginx",
						Group:   "arceus",
						Version: "v1",
					},
				},
			},
			Settings: []types.RuleSetting{
				{
					Path: "/type",
					Targets: []types.SettingTarget{
						{
							Name: "nginx",
							Sub:  "service",
							Fields: []types.SettingTargetField{
								{
									Path: "/spec/type",
								},
							},
						},
					},
				},
				{
					Path: "/match",
					Targets: []types.SettingTarget{
						{
							Name: "nginx",
							Sub:  "deploy",
							Fields: []types.SettingTargetField{
								{
									Path: "/spec/template/spec/affinity/nodeAffinity/requiredDuringSchedulingIgnoredDuringExecution/nodeSelectorTerms/0/matchExpressions",
								},
							},
						},
					},
				},
			},
		},
	}

	data = types.QuickStart{
		ObjectMeta: types.ObjectMeta{
			Name: "test",
		},
		Spec: types.QuickStartSpec{
			Rule: []types.QuickStartSpecRule{
				{
					Group:   "arceus",
					Version: "v1",
					Name:    "nginx-rule",
				},
			},
			Data: `{"type": "ClusterIP", "match": {"key": "node-role.kubernetes.io/edge", "operator": "NotExists"}}`,
		},
	}
)

//func TestParseSingle(t *testing.T) {
//	jsonData, err := yaml.YAMLToJSON([]byte(data.Spec.Data))
//	if err != nil {
//		t.Fatal(err)
//	}
//	result, err := ParseSingle(jsonData, &rule)
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log(result)
//}
//
//func TestParse(t *testing.T) {
//	result, err := Parse(&data)
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log(result)
//}
