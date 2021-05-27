// Package quick

// Created by zc on 2021/5/26.
package quick

import (
	"github.com/zc2638/arceus/pkg/types"
)

var (
	rule = types.QuickStartRule{
		ObjectMeta: types.ObjectMeta{
			Name: "test",
		},
		Spec: types.QuickStartRuleSpec{
			Templates: []types.RuleTemplateDefine{
				{
					Name:     "nginx",
					Group:    "arceus",
					Version:  "v1",
					Template: "nginx",
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
			Rule: []string{
				"nginx-rule",
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
