// Package quick

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
package quick

import (
	"net/http"

	"github.com/99nil/arceus/pkg/types"

	"github.com/zc2638/swag"
	"github.com/zc2638/swag/endpoint"
)

// Route handle template routing related
func Route(doc *swag.API) {
	doc.WithTags(swag.Tag{
		Name:        "quickstart",
		Description: "模板",
	}).AddEndpoint(
		endpoint.New(
			http.MethodGet, "/rule/list",
			endpoint.Handler(list()),
			endpoint.Summary("规则列表"),
			endpoint.ResponseSuccess(endpoint.Schema([]types.Resource{})),
		),
		endpoint.New(
			http.MethodPost, "/quickstart",
			endpoint.Handler(quickstart()),
			endpoint.Summary("快速开始"),
			endpoint.Body(types.QuickStart{}, "快速开始资源定义", true),
			endpoint.ResponseSuccess(),
		),
	)
}
