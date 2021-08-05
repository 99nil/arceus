// Package app

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
package app

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/99nil/arceus/handler/resource"
	"github.com/spf13/cobra"
)

func applyCommand() *cobra.Command {
	var filePath string
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply arceus resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			if filePath == "" {
				return errors.New("resource path not found")
			}
			stat, err := os.Stat(filePath)
			if err != nil {
				return err
			}

			if err := initConfig(); err != nil {
				return err
			}

			var fileURLs []string
			if stat.IsDir() {
				// 处理目录下所有文件
				dirs, err := os.ReadDir(filePath)
				if err != nil {
					return err
				}
				for _, dir := range dirs {
					if dir.IsDir() {
						continue
					}
					fileURLs = append(fileURLs, filepath.Join(filePath, dir.Name()))
				}
			} else {
				fileURLs = append(fileURLs, filePath)
			}
			return apply(fileURLs)
		},
	}
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "resource file path")
	return cmd
}

// TODO 支持quickstart资源使用
// TODO 新增create替代，与quickstart命令相对应
// TODO 兼容API和CMD，定义interface
func apply(fileURLs []string) error {
	for _, url := range fileURLs {
		fileData, err := os.ReadFile(url)
		if err != nil {
			return err
		}
		if err := resource.UploadResource(fileData); err != nil {
			return err
		}
	}
	return nil
}
