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
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/99nil/arceus/pkg/types"
	"github.com/pkgms/go/ctr"
	apiextensionsV1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/yaml"

	"github.com/99nil/arceus/global"
	"github.com/99nil/arceus/pkg/util"
	"github.com/99nil/arceus/static"
)

func upload() http.HandlerFunc {
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
		if err := UploadResource(fileData); err != nil {
			ctr.BadRequest(w, err)
			return
		}
		ctr.Success(w)
	}
}

func UploadResource(fileData []byte) error {
	var err error
	arr := bytes.Split(fileData, []byte("---\n"))
	for _, v := range arr {
		vb := bytes.TrimSpace(v)
		if len(vb) == 0 {
			continue
		}
		switch checkKind(vb) {
		case global.KindQuickStartRule:
			err = uploadQuickStart(vb)
		case global.KindTemplate:
			err = uploadTemplate(vb)
		case global.KindCustom:
			vb, _ = convertToCustom(vb)
			err = GenerateFile(vb, global.CustomResourcePath)
		default:
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateFile(source []byte, targetDir string) error {
	gvk := apiextensionsV1.SchemeGroupVersion.WithKind("CustomResourceDefinition")
	data := apiextensionsV1.CustomResourceDefinition{
		TypeMeta: v1.TypeMeta{
			APIVersion: gvk.GroupVersion().String(),
			Kind:       gvk.Kind,
		},
	}
	if err := runtime.DecodeInto(scheme.Codecs.UniversalDecoder(), source, &data); err != nil {
		return fmt.Errorf("resource parse failed: %s", err)
	}
	if data.Spec.Group == "" {
		return fmt.Errorf("group not found")
	}
	if data.Spec.Names.Kind == "" {
		return fmt.Errorf("kind not found")
	}
	if len(data.Spec.Versions) == 0 {
		return fmt.Errorf("version not found")
	}
	dir := filepath.Join(targetDir, data.Spec.Group, data.Spec.Names.Kind)
	if err := util.MkdirAll(dir); err != nil {
		return err
	}
	for k, v := range data.Spec.Versions {
		newFile, err := os.Create(filepath.Join(dir, v.Name+".yaml"))
		if err != nil {
			return fmt.Errorf("create file failed: %s", err)
		}
		current := data.DeepCopy()
		current.Spec.Versions = make([]apiextensionsV1.CustomResourceDefinitionVersion, 0, 1)
		current.Spec.Versions = append(current.Spec.Versions, data.Spec.Versions[k])
		// TODO 转为json格式节省空间
		b, err := yaml.Marshal(current)
		if err != nil {
			return fmt.Errorf("generate resource failed: %s", err)
		}
		if _, err := newFile.Write(b); err != nil {
			return fmt.Errorf("save file failed: %s", err)
		}
		newFile.Close()
	}
	return nil
}

func convertToCustom(source []byte) ([]byte, bool) {
	arr := bytes.Split(source, []byte("\n"))
	for k, v := range arr {
		if bytes.HasPrefix(v, []byte("apiVersion: "+types.Group+"/"+types.Version)) {
			arr[k] = []byte("apiVersion: apiextensions.k8s.io/v1")
			return bytes.Join(arr, []byte("\n")), true
		}
	}
	return source, false
}

func checkKind(source []byte) string {
	arr := bytes.Split(source, []byte("\n"))
	kindPrefix := "kind: "
	for _, v := range arr {
		if bytes.HasPrefix(v, []byte(kindPrefix+global.KindQuickStartRule)) {
			return global.KindQuickStartRule
		}
		if bytes.HasPrefix(v, []byte(kindPrefix+global.KindTemplate)) {
			return global.KindTemplate
		}
	}
	return global.KindNull
}

func uploadQuickStart(source []byte) error {
	var data types.QuickStartRule
	if err := yaml.Unmarshal(source, &data); err != nil {
		return err
	}
	if len(data.Spec.Templates) == 0 {
		return errors.New("spec.templates is necessary")
	}
	if len(data.Spec.Settings) == 0 {
		return errors.New("spec.settings is necessary")
	}

	dir := filepath.Join(global.RuleResourcePath, data.Spec.Group, data.Name)
	if err := util.MkdirAll(dir); err != nil {
		return err
	}
	newFile, err := os.Create(filepath.Join(dir, data.Spec.Version+".yaml"))
	if err != nil {
		return err
	}
	defer newFile.Close()

	_, err = newFile.Write(source)
	return nil
}

func uploadTemplate(source []byte) error {
	var template types.Template
	if err := yaml.Unmarshal(source, &template); err != nil {
		return err
	}
	if template.Kind != "Template" {
		return errors.New("kind must be Template")
	}
	if len(template.Spec.Template) == 0 {
		return errors.New("spec.template is necessary")
	}
	for k, v := range template.Spec.Template {
		if v.Name == "" {
			return fmt.Errorf("spec.template.%v.name is necessary", k)
		}
		gvk := schema.FromAPIVersionAndKind(v.APIVersion, v.Kind)
		filePath := filepath.Join(gvk.Group, gvk.Kind, gvk.Version) + ".yaml"
		baseFilePath := filepath.Join(static.KubernetesDir, filePath)
		_, err := fs.Stat(static.Kubernetes, baseFilePath)
		if os.IsNotExist(err) {
			_, err = fs.Stat(os.DirFS(global.CustomResourcePath), filePath)
		}
		if err != nil {
			return fmt.Errorf("resource (%s) not exist", gvk.String())
		}
	}
	dir := filepath.Join(global.TemplateResourcePath, template.Spec.Group, template.Name)
	if err := util.MkdirAll(dir); err != nil {
		return err
	}
	newFile, err := os.Create(filepath.Join(dir, template.Spec.Version+".yaml"))
	if err != nil {
		return err
	}
	defer newFile.Close()

	_, err = newFile.Write(source)
	return err
}
