// Copyright 2019 Yunion
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

package webconsole

import (
	"yunion.io/x/onecloud/pkg/mcclient/modulebase"
)

var (
	CommandLog *CommandLogManager
)

func init() {
	CommandLog = NewCommandLogManager()

	modulebase.Register("v1", CommandLog)
}

type CommandLogManager struct {
	modulebase.ResourceManager
}

func NewCommandLogManager() *CommandLogManager {
	return &CommandLogManager{
		modulebase.ResourceManager{
			BaseManager: *modulebase.NewBaseManager("webconsole", "", "webconsole", []string{
				"id", "ops_time", "obj_id", "obj_type", "obj_name", "user", "user_id", "tenant", "tenant_id", "owner_tenant_id", "action", "notes",
				"session_id", "accessed_at", "type", "login_user", "start_time", "ps1", "command",
			}, nil),
			Keyword: "commandlog", KeywordPlural: "commandlogs",
		},
	}
}
