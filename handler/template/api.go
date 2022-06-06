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

	pkgtypes "github.com/99nil/arceus/pkg/types"

	"github.com/zc2638/swag"
	"github.com/zc2638/swag/endpoint"
	"github.com/zc2638/swag/types"
)

// Route handle template routing related
func Route(doc *swag.API) {
	doc.WithTags(swag.Tag{
		Name:        "template",
		Description: "模板",
	}).AddEndpoint(
		endpoint.New(
			http.MethodGet, "/template",
			endpoint.Handler(list()),
			endpoint.Summary("资源列表"),
			endpoint.ResponseSuccess(endpoint.Schema([]pkgtypes.Resource{})),
		),
		endpoint.New(
			http.MethodGet, "/template/info",
			endpoint.Handler(info()),
			endpoint.Summary("模板详情"),
			endpoint.Query("group", types.String, "资源分组", true),
			endpoint.Query("kind", types.String, "资源名称", true),
			endpoint.Query("version", types.String, "资源版本", true),
			endpoint.ResponseSuccess(),
		),
		endpoint.New(
			http.MethodPost, "/template",
			endpoint.Handler(create()),
			endpoint.Summary("模板创建"),
			endpoint.FormData("file", types.File, "file to upload", true),
			endpoint.ResponseSuccess(),
		),
	)
}
