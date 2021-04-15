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
package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/zc2638/arceus/global"
	"github.com/zc2638/arceus/handler/resource"

	"github.com/pkgms/go/server"
	"github.com/zc2638/arceus/handler"

	"github.com/spf13/cobra"
)

var cfgFile string

func NewServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "arceus",
		SilenceUsage: true,
		RunE:         run,
	}
	cfgFilePath := os.Getenv("ARCEUS_CONFIG")
	if cfgFilePath == "" {
		cfgFilePath = "config/config.yaml"
	}
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", cfgFilePath, "config file (default is $HOME/config.yaml)")
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	s := server.New(&server.Config{
		Port: 2638,
	})
	s.Handler = handler.New()
	fmt.Println("Listen on", s.Addr)
	return s.Run(ctx)
}

func generate() {
	dirPath := filepath.Join(global.ResourcePath, "crd")
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range dir {
		if v.IsDir() {
			fmt.Printf("generate file skipped: %s\n", v.Name())
			continue
		}
		file, err := os.ReadFile(filepath.Join(dirPath, v.Name()))
		if err != nil {
			fmt.Printf("read file skipped: %s\n", v.Name())
			continue
		}
		if err := resource.GenerateFile(file); err != nil {
			fmt.Printf("generate file failed: %s, %s\n", v.Name(), err)
		}
	}
}
