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

package static

import "embed"

const (
	PatchDir      = "patch"
	KubernetesDir = "kubernetes"
	UIDir         = "ui"
)

// Patch is a dir for patch resources
//
//go:embed patch
var Patch embed.FS

// Kubernetes is a dir for static content
//
//go:embed kubernetes
var Kubernetes embed.FS

// UI is a dir for static content
//
//go:embed ui
var UI embed.FS
