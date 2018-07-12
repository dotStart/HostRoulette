/*
 * Copyright 2018 Johannes Donath <johannesd@torchmind.com>
 * and other copyright owners as documented in the project's IP log.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package server

import (
  "encoding/json"
  "github.com/dotStart/HostRoulette/search"
  "net/http"
  "strings"
)

func (s *Server) HandleGetCommunity(w http.ResponseWriter, req *http.Request) {
  if req.Method == "OPTIONS" {
    s.writeHeaders(w)
    return
  }
  if req.Method != "POST" || req.URL.Path != "/api/community" {
    http.NotFound(w, req)
    return
  }

  ids := make([]string, 0)
  defer req.Body.Close()
  err := json.NewDecoder(req.Body).Decode(&ids)
  if err != nil {
    s.logger.Errorf("failed to decode request body: %s", err)
    http.Error(w, "cannot read request body", http.StatusBadRequest)
    return
  }

  result, err := s.search.GetCommunities(ids)
  if err != nil {
    s.logger.Errorf("failed to retrieve communities: %s", err)
    http.Error(w, "failed to execute query", http.StatusServiceUnavailable)
  }

  s.writeHeaders(w)
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(result)
}

func (s *Server) HandleCommunitySearch(w http.ResponseWriter, req *http.Request) {
  if req.Method == "OPTIONS" {
    s.writeHeaders(w)
    return
  }

  query := req.URL.Path[22:]
  if strings.Contains(query, "/") {
    http.NotFound(w, req)
    return
  }

  var err error
  var result *search.Result
  if query == "" {
    result = &search.Result{
      Success: true,
      Matches: make([]*search.Match, 0),
    }
  } else {
    result, err = s.search.SearchCommunity(query)
    if err != nil {
      s.logger.Errorf("failed to execute search query for community using query \"%s\": %s", query, err)
      http.Error(w, "failed to execute query", http.StatusServiceUnavailable)
      return
    }
  }

  s.writeHeaders(w)
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(result)
}
