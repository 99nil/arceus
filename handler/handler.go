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
package handler

import (
	"net/http"

	"github.com/zc2638/arceus/handler/quick"

	"github.com/zc2638/arceus/handler/template"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/zc2638/arceus/handler/home"
	"github.com/zc2638/arceus/handler/resource"
	"github.com/zc2638/arceus/handler/web"
	"github.com/zc2638/swag"
)

func New() http.Handler {
	mux := chi.NewMux()
	mux.Use(
		middleware.Recoverer,
		middleware.Logger,
		cors.AllowAll().Handler,
	)
	mux.Mount("/web", web.New())
	apiDoc := swag.New(swag.Title("Arceus API Doc"))
	apiDoc.AddEndpointFunc(
		home.Route,
		resource.Route,
		template.Route,
		quick.Route,
	)
	apiDoc.RegisterMuxWithData(mux, false)
	return mux
}
