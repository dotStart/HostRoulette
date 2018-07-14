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
  "net/http"
)

func (s *Server) tickUpdateTwitchStatistics() {
  for range s.twitchStatisticsTicker.C {
    s.updateTwitchStatistics()
  }
}

func (s *Server) updateTwitchStatistics() {
  s.statisticsLock.Lock()
  defer s.statisticsLock.Unlock()

  stats, err := s.twitchClient.GetStatistics()
  if err != nil {
    s.logger.Errorf("failed to retrieve twitch statistics: %s", err)
    return
  }

  s.twitchStatistics = stats
}

func (s *Server) HandleStatistics(w http.ResponseWriter, req *http.Request) {
  if req.Method != "GET" || req.URL.Path != "/api/statistics" {
    http.NotFound(w, req)
    return
  }

  s.statisticsLock.RLock()
  twitchStats := s.twitchStatistics
  s.statisticsLock.RUnlock()

  spins, err := s.cacheClient.GetSpinStatistic()
  if err != nil {
    s.logger.Errorf("failed to retrieve spin statistic: %s", err)
  }

  stats := &struct {
    Spins    uint64 `json:"spins"`
    Channels uint64 `json:"channels"`
    Viewers  uint64 `json:"viewers"`
  }{
    Spins:    spins,
    Channels: twitchStats.Channels,
    Viewers:  twitchStats.Viewers,
  }

  s.writeHeaders(w)
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(stats)
}
