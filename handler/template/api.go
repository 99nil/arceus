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
	"net/http"

	"github.com/zc2638/arceus/pkg/data/resp"

	"github.com/zc2638/arceus/pkg/types"

	"github.com/zc2638/swag/endpoint"
	"github.com/zc2638/swag/swagger"
)

const tag = "template"

// Route handle template routing related
func Route(doc *swagger.API) {
	doc.Tags = append(doc.Tags, swagger.Tag{
		Name:        tag,
		Description: "模板",
	})
	doc.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/template",
			endpoint.Handler(list()),
			endpoint.Summary("资源列表"),
			endpoint.ResponseSuccess(endpoint.Schema([]types.Resource{})),
			endpoint.Tags(tag),
		),
		endpoint.New(
			http.MethodGet, "/template/info",
			endpoint.Handler(info()),
			endpoint.Summary("模板详情"),
			endpoint.Query("group", swagger.TypeString, "资源分组", true),
			endpoint.Query("kind", swagger.TypeString, "资源名称", true),
			endpoint.Query("version", swagger.TypeString, "资源版本", true),
			endpoint.ResponseSuccess(endpoint.Schema([]resp.TemplateDataItem{})),
			endpoint.Tags(tag),
		),
		endpoint.New(
			http.MethodPost, "/template",
			endpoint.Handler(create()),
			endpoint.Summary("模板创建"),
			func(b *endpoint.Builder) {
				if b.Endpoint.Parameters == nil {
					b.Endpoint.Parameters = []swagger.Parameter{}
				}
				b.Endpoint.Parameters = append(b.Endpoint.Parameters, swagger.Parameter{
					In:          "formData",
					Type:        "file",
					Name:        "file",
					Description: "file to upload",
					Required:    true,
				})
			},
			endpoint.ResponseSuccess(),
			endpoint.Tags(tag),
		),
	)
}
