// Package raw
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
package raw

import (
	"net/http"

	"github.com/zc2638/swag/endpoint"
	"github.com/zc2638/swag/swagger"
)

// Route handle raw routing related
func Route(doc *swagger.API) {
	doc.WithTags(swagger.Tag{
		Name:        "raw",
		Description: "数据",
	}).AddEndpoint(
		endpoint.New(
			http.MethodPost, "/raw/convert",
			endpoint.Handler(convert()),
			endpoint.Summary("数据转换"),
			endpoint.Query("in", swagger.TypeString, "数据格式: json、yaml、xml、toml", true),
			endpoint.Query("out", swagger.TypeString, "转换格式: json、yaml、xml、toml", true),
			endpoint.Body(struct{}{}, "", true),
			endpoint.ResponseSuccess(),
		),
	)
}
