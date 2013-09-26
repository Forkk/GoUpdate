// Copyright 2013 MultiMC Contributors
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

package repo

// Structures that represent a GoUpdate index file.

const IndexFileName = "index.json"

type VersionSummary struct {
	Id   int
	Name string
}

type Channel struct {
	Id   string
	Name string
	CurrentVersion int
}

type Index struct {
	ApiVersion int
	Versions   []VersionSummary
	Channels   []Channel
}

// Returns a new, blank index struct.
func NewBlankIndex() Index {
	return Index{ApiVersion: 0, Versions: []VersionSummary{}, Channels: []Channel{}}
}
