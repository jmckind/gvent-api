// Copyright [yyyy] [name of copyright owner]
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

package gvent

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// writeJSON will write the given object as JSON.
func writeJSON(w http.ResponseWriter, data interface{}) {
	json, err := json.Marshal(data)
	if err != nil {
		log.Error("Unable to marshal data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(json)
}

// writeList will write the given slice as JSON, using the given label.
func writeList(w http.ResponseWriter, items interface{}, label string) {
	data := make(map[string]interface{})
	data[label] = items

	json, err := json.Marshal(data)
	if err != nil {
		log.Error("Unable to marshal data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(json)
}
