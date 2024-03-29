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
	"io"
	"net/http"

	"github.com/99nil/ditto"
	"github.com/99nil/gopkg/ctr"
)

func convert() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in := r.URL.Query().Get("in")
		out := r.URL.Query().Get("out")
		data, err := io.ReadAll(r.Body)
		if err != nil {
			ctr.InternalError(w, err)
			return
		}
		result, err := ditto.NewTransfer(in, out).Exchange(data)
		if err != nil {
			ctr.InternalError(w, err)
			return
		}
		ctr.Bytes(w, result)
	}
}
