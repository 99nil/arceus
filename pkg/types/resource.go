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

package types

import (
	"io/fs"
	"path/filepath"
	"strings"
)

type Resource struct {
	Value    string     `json:"value"`
	Label    string     `json:"label"`
	Children []Resource `json:"children,omitempty"`
}

func BuildResourceList(base fs.FS, basePath string) ([]Resource, error) {
	dirs, err := fs.ReadDir(base, basePath)
	if err != nil {
		return nil, err
	}
	var list []Resource
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		kindPath := filepath.Join(basePath, dir.Name())
		kindDirs, err := fs.ReadDir(base, kindPath)
		if err != nil {
			continue
		}
		group := Resource{
			Value: dir.Name(),
			Label: dir.Name(),
		}
		for _, kindDir := range kindDirs {
			if !kindDir.IsDir() {
				continue
			}
			versionPath := filepath.Join(kindPath, kindDir.Name())
			versionFiles, err := fs.ReadDir(base, versionPath)
			if err != nil {
				continue
			}
			kind := Resource{
				Value: kindDir.Name(),
				Label: kindDir.Name(),
			}
			for _, versionFile := range versionFiles {
				if versionFile.IsDir() {
					continue
				}
				version := strings.TrimSuffix(versionFile.Name(), ".yaml")
				kind.Children = append(kind.Children, Resource{
					Value: version,
					Label: version,
				})
			}
			group.Children = append(group.Children, kind)
		}
		list = append(list, group)
	}
	return list, nil
}
