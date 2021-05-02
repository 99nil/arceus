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
package template

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkgms/go/ctr"
	apiextensionsV1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/yaml"

	"github.com/zc2638/arceus/global"
	"github.com/zc2638/arceus/pkg/data/resp"
	"github.com/zc2638/arceus/pkg/types"
	"github.com/zc2638/arceus/static"
)

func list() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := types.BuildResourceList(os.DirFS(global.ResourcePath), "template")
		if err != nil {
			ctr.InternalError(w, err)
			return
		}
		ctr.OK(w, result)
	}
}

func info() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		group := r.URL.Query().Get("group")
		kind := r.URL.Query().Get("kind")
		version := r.URL.Query().Get("version")

		filePath := filepath.Join(group, kind, version) + ".yaml"
		fileData, err := fs.ReadFile(os.DirFS(global.TemplateResourcePath), filePath)
		if err != nil {
			ctr.BadRequest(w, errors.New("resource not exist"))
			return
		}
		var template types.Template
		if err := yaml.Unmarshal(fileData, &template); err != nil {
			ctr.InternalError(w, err)
			return
		}

		var result = make([]resp.TemplateDataItem, 0, len(template.Spec.Template))
		for _, v := range template.Spec.Template {
			gvk := schema.FromAPIVersionAndKind(v.APIVersion, v.Kind)
			filePath := filepath.Join(gvk.Group, gvk.Kind, gvk.Version) + ".yaml"
			baseFilePath := filepath.Join(static.KubernetesDir, filePath)
			fileData, err := fs.ReadFile(static.Kubernetes, baseFilePath)
			if os.IsNotExist(err) {
				fileData, err = fs.ReadFile(os.DirFS(global.CustomResourcePath), filePath)
			}
			if err != nil {
				continue
			}
			var data apiextensionsV1.CustomResourceDefinition
			if err := runtime.DecodeInto(scheme.Codecs.UniversalDecoder(), fileData, &data); err != nil {
				ctr.InternalError(w, fmt.Errorf("resource parse failed: %s", err))
				return
			}
			apiVersion := data.Spec.Group + "/" + data.Spec.Versions[0].Name
			node := types.BuildNode(data.Spec.Versions[0].Schema.OpenAPIV3Schema, nil, apiVersion, data.Spec.Names.Kind)
			result = append(result, resp.TemplateDataItem{
				Template: v.Data,
				Data:     node,
			})
		}
		ctr.OK(w, result)
	}
}

func create() http.HandlerFunc {
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
		var template types.Template
		if err := yaml.Unmarshal(fileData, &template); err != nil {
			ctr.BadRequest(w, err)
			return
		}
		if template.Kind != "Template" {
			ctr.BadRequest(w, errors.New("kind must be Template"))
			return
		}
		if len(template.Spec.Template) == 0 {
			ctr.BadRequest(w, errors.New("spec.template is necessary"))
			return
		}
		for k, v := range template.Spec.Template {
			if v.Name == "" {
				ctr.BadRequest(w, fmt.Errorf("spec.template.%v.name is necessary", k))
				return
			}
			gvk := schema.FromAPIVersionAndKind(v.APIVersion, v.Kind)
			filePath := filepath.Join(gvk.Group, gvk.Kind, gvk.Version) + ".yaml"
			baseFilePath := filepath.Join(static.KubernetesDir, filePath)
			_, err := fs.Stat(static.Kubernetes, baseFilePath)
			if os.IsNotExist(err) {
				_, err = fs.Stat(os.DirFS(global.CustomResourcePath), filePath)
			}
			if err != nil {
				ctr.BadRequest(w, fmt.Errorf("resource (%s) not exist", gvk.String()))
				return
			}
		}
		newFilePath := filepath.Join(global.TemplateResourcePath,
			template.Spec.Group, template.Name, template.Spec.Version+".yaml")
		newFile, err := os.Create(newFilePath)
		if err != nil {
			ctr.InternalError(w, err)
			return
		}
		defer newFile.Close()
		if _, err := newFile.Write(fileData); err != nil {
			ctr.InternalError(w, err)
			return
		}
		ctr.Success(w)
	}
}
