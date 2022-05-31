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
package global

import (
	"github.com/99nil/arceus/pkg/util"
	"github.com/99nil/gopkg/ctr"

	"github.com/sirupsen/logrus"
)

func Init(cfg *Config) error {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:            true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
		FullTimestamp:          true,
		TimestampFormat:        "2006/01/02 15:04:05",
	})
	ctr.InitLogger(logrus.StandardLogger())

	buildAllPath(cfg.BasePath)
	if err := util.MkdirAll(CustomResourcePath); err != nil {
		return err
	}
	if err := util.MkdirAll(TemplateResourcePath); err != nil {
		return err
	}
	return util.MkdirAll(RuleResourcePath)
}

type Config struct {
	BasePath string `json:"base_path"`
}

func Environ() *Config {
	cfg := &Config{}
	cfg.BasePath = ResourcePath
	return cfg
}
