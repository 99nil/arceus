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

	apiextensionsV1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

const (
	NodeRoot       = "root"
	NodeAPIVersion = "apiVersion"
	NodeKind       = "kind"
)

const (
	TypeObject = "object"
	TypeArray  = "array"
	TypeString = "string"
)

type TNodeDesc struct {
	Locale string `json:"locale"` // 区域
	Desc   string `json:"desc"`   // 描述
}

type TNode struct {
	Key      string      `json:"key"`      // 全局唯一，格式为 index.path 如 pod.spec.containers.0.name
	Name     string      `json:"name"`     // 名称
	Title    string      `json:"title"`    // 节点标题
	Type     string      `json:"type"`     // 值类型
	Value    string      `json:"value"`    // 默认值
	Descs    []TNodeDesc `json:"descs"`    // 描述（默认描述、中文描述）
	Required []string    `json:"required"` // 关联，仅在object节点定义
	Children []TNode     `json:"children"` // 子节点
}

func BuildNode(prop *apiextensionsV1.JSONSchemaProps, node *TNode, extras ...string) *TNode {
	if prop.Type != TypeObject {
		return nil
	}
	if node == nil {
		node = &TNode{
			Key:      NodeRoot,
			Name:     NodeRoot,
			Title:    NodeRoot,
			Type:     TypeObject,
			Required: prop.Required,
		}
	}
	for k, v := range prop.Properties {
		cNode := &TNode{}
		cNode.Key = node.Key + "." + k
		cNode.Name = k
		cNode.Title = cNode.Name
		cNode.Type = v.Type
		cNode.Required = v.Required
		if v.Description != "" {
			cNode.Descs = append(cNode.Descs, TNodeDesc{
				Desc: v.Description,
			})
		}
		if v.Default != nil && v.Default.Raw != nil {
			cNode.Value = string(bytes.Trim(v.Default.Raw, "\""))
		}
		switch v.Type {
		// 对象的时候，需要向下解析properties
		// 类型为object时，并且无children，则默认为string/string
		case TypeObject:
			cNode = BuildNode(&v, cNode)

		// 数组的时候，需要向下解析items，需要加入一个空的数组节点
		case TypeArray:
			array := buildArrayNode(&v, cNode)
			cNode.Children = append(cNode.Children, *array)
		}
		node.Children = append(node.Children, *cNode)
	}
	completeAPIVersion(node, extras...)
	return node
}

func buildArrayNode(prop *apiextensionsV1.JSONSchemaProps, pNode *TNode) *TNode {
	v := prop.Items.Schema
	name := "-"
	node := &TNode{
		Key:      pNode.Key + "." + name,
		Name:     name,
		Title:    name,
		Required: v.Required,
		Type:     v.Type,
	}
	if v.Type == TypeObject {
		node = BuildNode(v, node)
	}
	return node
}

func completeAPIVersion(node *TNode, extras ...string) {
	if node.Name != NodeRoot {
		return
	}
	if len(extras) < 2 {
		return
	}
	apiVersion := extras[0]
	kind := extras[1]
	var versionNode, kindNode *TNode
	for k, v := range node.Children {
		if v.Name == NodeAPIVersion {
			versionNode = &node.Children[k]
		}
		if v.Name == NodeKind {
			kindNode = &node.Children[k]
		}
	}
	if versionNode == nil {
		node.Children = append(node.Children, TNode{
			Key:   node.Key + "." + NodeAPIVersion,
			Name:  NodeAPIVersion,
			Title: NodeAPIVersion,
			Type:  TypeString,
			Value: apiVersion,
		})
	} else {
		versionNode.Value = apiVersion
	}
	if kindNode == nil {
		node.Children = append(node.Children, TNode{
			Key:   node.Key + "." + NodeKind,
			Name:  NodeKind,
			Title: NodeKind,
			Type:  TypeString,
			Value: kind,
		})
	} else {
		kindNode.Value = kind
	}
	node.Required = append(node.Required, NodeAPIVersion, NodeKind)
}
