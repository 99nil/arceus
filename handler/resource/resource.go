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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	apiextensionsV1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/pkgms/go/ctr"
	"github.com/zc2638/arceus/global"
)

type Group struct {
	Value    string `json:"value"`
	Label    string `json:"label"`
	Children []Kind `json:"children"`
}

type Kind struct {
	Value    string    `json:"value"`
	Label    string    `json:"label"`
	Children []Version `json:"children"`
}

type Version struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

func list() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dirs, err := os.ReadDir(global.ResourcePath)
		if err != nil {
			ctr.InternalError(w, err)
			return
		}
		var list []Group
		for _, dir := range dirs {
			if !dir.IsDir() {
				continue
			}
			kindPath := filepath.Join(global.ResourcePath, dir.Name())
			kindDirs, err := os.ReadDir(kindPath)
			if err != nil {
				continue
			}
			group := Group{
				Value: dir.Name(),
				Label: dir.Name(),
			}
			for _, kindDir := range kindDirs {
				if !kindDir.IsDir() {
					continue
				}
				versionPath := filepath.Join(kindPath, kindDir.Name())
				versionFiles, err := os.ReadDir(versionPath)
				if err != nil {
					continue
				}
				kind := Kind{
					Value: kindDir.Name(),
					Label: kindDir.Name(),
				}
				for _, versionFile := range versionFiles {
					if versionFile.IsDir() {
						continue
					}
					version := strings.TrimSuffix(versionFile.Name(), ".yaml")
					kind.Children = append(kind.Children, Version{
						Value: version,
						Label: version,
					})
				}
				group.Children = append(group.Children, kind)
			}
			list = append(list, group)
		}
		ctr.OK(w, list)
	}
}

func info() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		group := r.URL.Query().Get("group")
		kind := r.URL.Query().Get("kind")
		version := r.URL.Query().Get("version")
		filePath := filepath.Join(global.ResourcePath, group, kind, version) + ".yaml"
		data, err := os.ReadFile(filePath)
		if err != nil {
			ctr.BadRequest(w, errors.New("resource not exist"))
			return
		}
		ctr.Bytes(w, data)
	}
}

func tree() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		group := r.URL.Query().Get("group")
		kind := r.URL.Query().Get("kind")
		version := r.URL.Query().Get("version")
		filePath := filepath.Join(global.ResourcePath, group, kind, version) + ".yaml"
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			ctr.BadRequest(w, errors.New("resource not exist"))
			return
		}
		// 解析到结构体
		var data apiextensionsV1.CustomResourceDefinition
		if err := runtime.DecodeInto(scheme.Codecs.UniversalDecoder(), fileData, &data); err != nil {
			ctr.InternalError(w, fmt.Errorf("resource parse failed: %s", err))
			return
		}
		// 转换为自定义的树结构
		if len(data.Spec.Versions) == 0 {
			ctr.OK(w, struct{}{})
			return
		}
		apiVersion := data.Spec.Group + "/" + data.Spec.Versions[0].Name
		node := BuildNode(data.Spec.Versions[0].Schema.OpenAPIV3Schema, nil, apiVersion, data.Spec.Names.Kind)
		b, err := json.Marshal(node)
		if err != nil {
			ctr.InternalError(w, fmt.Errorf("resource parse failed: %s", err))
			return
		}
		ctr.Bytes(w, b)
	}
}
