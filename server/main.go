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
  "github.com/dotStart/HostRoulette/cache"
  "github.com/dotStart/HostRoulette/search"
  "github.com/dotStart/HostRoulette/twitch"
  "github.com/op/go-logging"
  "math/rand"
  "net/http"
  "strings"
  "sync"
  "time"
)

type Server struct {
  logger       *logging.Logger
  search       *search.Client
  cacheClient  *cache.Cache
  twitchClient *twitch.Client

  CorsDisabled bool

  // statistics
  statisticsLock         sync.RWMutex
  twitchStatisticsTicker *time.Ticker
  twitchStatistics       *twitch.Statistics
}

func New(mux *http.ServeMux, search *search.Client, cacheClient *cache.Cache, twitchClient *twitch.Client) *Server {
  rand.Seed(time.Now().Unix())

  srv := &Server{
    logger:                 logging.MustGetLogger("srv"),
    search:                 search,
    cacheClient:            cacheClient,
    twitchClient:           twitchClient,
    twitchStatisticsTicker: time.NewTicker(time.Minute * 30),
    twitchStatistics:       &twitch.Statistics{},
  }

  srv.updateTwitchStatistics()
  go srv.tickUpdateTwitchStatistics()

  mux.HandleFunc("/api/community", srv.HandleGetCommunity)
  mux.HandleFunc("/api/game", srv.HandleGetGame)
  mux.HandleFunc("/api/wheel", srv.HandleWheel)
  mux.HandleFunc("/api/search/community/", srv.HandleCommunitySearch)
  mux.HandleFunc("/api/statistics", srv.HandleStatistics)
  mux.HandleFunc("/api/search/game/", srv.HandleGameSearch)

  return srv
}

func (s *Server) writeHeaders(w http.ResponseWriter) {
  if s.CorsDisabled {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Access-Control-Expose-Headers", "X-Rate-Limit,X-Rate-Limit-Remaining,X-Spin-Count")
  }
}

func getAddress(req *http.Request) string {
  forward := req.Header.Get("X-Forwarded-For")
  if forward != "" {
    return forward
  }

  index := strings.LastIndex(req.RemoteAddr, ":")
  if index == -1 {
    return req.RemoteAddr
  }

  addr := req.RemoteAddr[:index]
  if !strings.HasPrefix(addr, "[") {
    return addr
  }
  return addr[1 : len(addr)-1]
}
