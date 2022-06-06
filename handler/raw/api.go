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

	"github.com/zc2638/swag"
	"github.com/zc2638/swag/endpoint"
	"github.com/zc2638/swag/types"
)

// Route handle raw routing related
func Route(doc *swag.API) {
	doc.WithTags(swag.Tag{
		Name:        "raw",
		Description: "数据",
	}).AddEndpoint(
		endpoint.New(
			http.MethodPost, "/raw/convert",
			endpoint.Handler(convert()),
			endpoint.Summary("数据转换"),
			endpoint.Query("in", types.String, "数据格式: json、yaml、xml、toml", true),
			endpoint.Query("out", types.String, "转换格式: json、yaml、xml、toml", true),
			endpoint.Body(struct{}{}, "", true),
			endpoint.ResponseSuccess(),
		),
	)
}
