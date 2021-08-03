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
	"bytes"
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/zc2638/arceus/handler/quick"
	"github.com/zc2638/arceus/pkg/types"
	"github.com/zc2638/arceus/pkg/util"
)

type QuickStartOption struct {
	FilePath   string
	OutputPath string
	Values     []string
}

func quickstartCommand() *cobra.Command {
	opt := &QuickStartOption{}
	cmd := &cobra.Command{
		Use: "quickstart",
		Aliases: []string{
			"qs",
		},
		Short: "Quick start to use",
		RunE: func(cmd *cobra.Command, args []string) error {
			if opt.FilePath == "" {
				return errors.New("QuickStart resource file path not found")
			}
			stat, err := os.Stat(opt.FilePath)
			if err != nil {
				return err
			}

			if err := initConfig(); err != nil {
				return err
			}

			var fileURLs []string
			if stat.IsDir() {
				// 处理目录下所有文件
				dirs, err := os.ReadDir(opt.FilePath)
				if err != nil {
					return err
				}
				for _, dir := range dirs {
					if dir.IsDir() {
						continue
					}
					fileURLs = append(fileURLs, filepath.Join(opt.FilePath, dir.Name()))
				}
			} else {
				fileURLs = append(fileURLs, opt.FilePath)
			}
			return quickstart(opt, fileURLs)
		},
	}
	cmd.Flags().StringVarP(&opt.FilePath, "file", "f", "", "Resource file path")
	cmd.Flags().StringVarP(&opt.OutputPath, "output", "o", "quickstart", "Resource output path")
	cmd.Flags().StringArrayVarP(&opt.Values, "values", "v", nil, "Resource parameters (key-value pair) e.g. app.name=test")
	return cmd
}

func quickstart(opt *QuickStartOption, fileURLs []string) error {
	pairs := types.ParseKValuePairs(opt.Values)
	for _, url := range fileURLs {
		fileData, err := os.ReadFile(url)
		if err != nil {
			return err
		}
		var data types.QuickStart
		if err := yaml.Unmarshal(fileData, &data); err != nil {
			return err
		}
		if data.Kind != types.QuickStartKind {
			continue
		}

		result, err := quick.Parse(&data, pairs...)
		if err != nil {
			return err
		}
		filename := data.Name
		if filename == "" {
			filename = util.RandomStr(6)
		}

		if err := util.MkdirAll(opt.OutputPath); err != nil {
			return err
		}
		newFile, err := os.Create(filepath.Join(opt.OutputPath, filename+".yaml"))
		if err != nil {
			return err
		}

		dataSet := make([][]byte, 0, len(result))
		for _, v := range result {
			b, err := yaml.Marshal(v)
			if err != nil {
				return err
			}
			dataSet = append(dataSet, bytes.TrimRight(b, "\n"))
		}
		yamlData := bytes.Join(dataSet, []byte("\n---\n"))
		_, _ = newFile.Write(yamlData)
		newFile.Close()
	}
	return nil
}
