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

	"github.com/99nil/arceus/global"
	"github.com/99nil/arceus/pkg/types"
	"github.com/99nil/arceus/pkg/util"
	"github.com/99nil/arceus/static"
	"github.com/99nil/gopkg/ctr"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/yaml"
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
		ctr.OK(w, template)
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
		dir := filepath.Join(global.TemplateResourcePath, template.Spec.Group, template.Name)
		if err := util.MkdirAll(dir); err != nil {
			ctr.InternalError(w, err)
			return
		}
		newFile, err := os.Create(filepath.Join(dir, template.Spec.Version+".yaml"))
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
