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

package web

import (
	"net/http"

	"github.com/99nil/arceus/static"

	"github.com/go-chi/chi"
	"github.com/zc2638/swag"
)

func New() http.Handler {
	mux := chi.NewMux()
	mux.Mount("/", http.StripPrefix("/web", fileSystem()))
	return mux
}

func fileSystem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" {
			localRedirect(w, r, "/web/")
			return
		}
		http.FileServer(swag.DirFS(static.UIDir, static.UI)).ServeHTTP(w, r)
	}
}

func localRedirect(w http.ResponseWriter, r *http.Request, newPath string) {
	if q := r.URL.RawQuery; q != "" {
		newPath += "?" + q
	}
	w.Header().Set("Location", newPath)
	w.WriteHeader(http.StatusMovedPermanently)
}
