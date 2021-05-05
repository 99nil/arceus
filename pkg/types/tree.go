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
package types

import (
	"bytes"
	"sort"
	"strings"

	apiextensionsV1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

const (
	NodeRoot       = "root"
	NodeAPIVersion = "apiVersion"
	NodeKind       = "kind"
	NodeMetadata   = "metadata"
)

const (
	TypeObject  = "object"
	TypeArray   = "array"
	TypeString  = "string"
	TypeNumber  = "number"
	TypeBoolean = "boolean"
)

// TODO change to #
const TreeNodeArray = "-"

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
	Enums    []string    `json:"enums"`    // 枚举
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
		for _, e := range v.Enum {
			if e.Raw == nil {
				continue
			}
			cNode.Enums = append(cNode.Enums, string(bytes.Trim(e.Raw, "\"")))
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
	sort.SliceStable(node.Children, func(i, j int) bool {
		return strings.Compare(node.Children[i].Name, node.Children[j].Name) < 0
	})
	completeAPIVersion(node, extras...)
	completeMetadata(node)
	return node
}

func buildArrayNode(prop *apiextensionsV1.JSONSchemaProps, pNode *TNode) *TNode {
	v := prop.Items.Schema
	node := &TNode{
		Key:      pNode.Key + "." + TreeNodeArray,
		Name:     TreeNodeArray,
		Title:    TreeNodeArray,
		Required: v.Required,
		Type:     v.Type,
	}
	if v.Type == TypeObject {
		node = BuildNode(v, node)
	}
	return node
}

func completeAPIVersion(node *TNode, extras ...string) {
	if node.Key != NodeRoot {
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

func completeMetadata(node *TNode) {
	if node.Name != NodeMetadata {
		return
	}
	if len(node.Children) == 0 && node.Type == TypeObject {
		node.Required = []string{"name"}
		node.Children = make([]TNode, 0, 5)
		node.Children = append(node.Children, TNode{
			Key:   node.Key + ".name",
			Name:  "name",
			Title: "name",
			Type:  TypeString,
			Descs: []TNodeDesc{
				{
					Locale: "",
					Desc: `Name must be unique within a namespace. 
Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically.
Name is primarily intended for creation idempotence and configuration definition.
Cannot be updated.`,
				},
			},
		}, TNode{
			Key:   node.Key + ".namespace",
			Name:  "namespace",
			Title: "namespace",
			Type:  TypeString,
			Descs: []TNodeDesc{
				{
					Desc: `Namespace defines the space within which each name must be unique. 
An empty namespace is equivalent to the "default" namespace, but "default" is the canonical representation.
Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty.`,
				},
			},
		}, TNode{
			Key:   node.Key + ".clusterName",
			Name:  "clusterName",
			Title: "clusterName",
			Type:  TypeString,
			Descs: []TNodeDesc{
				{
					Desc: `The name of the cluster which the object belongs to.
This is used to distinguish resources with same name and namespace in different clusters.
This field is not set anywhere right now and apiserver is going to ignore it if set in create or update request.`,
				},
			},
		}, TNode{
			Key:   node.Key + ".labels",
			Name:  "labels",
			Title: "labels",
			Type:  TypeObject,
			Descs: []TNodeDesc{
				{
					Desc: `Map of string keys and values that can be used to organize and categorize (scope and select) objects.
May match selectors of replication controllers and services.
More info: http://kubernetes.io/docs/user-guide/labels`,
				},
			},
		}, TNode{
			Key:   node.Key + ".annotations",
			Name:  "annotations",
			Title: "annotations",
			Type:  TypeObject,
			Descs: []TNodeDesc{
				{
					Desc: `Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. 
They are not queryable and should be preserved when modifying objects.
More info: http://kubernetes.io/docs/user-guide/annotations`,
				},
			},
		})
	}
}
