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
  "fmt"
  "github.com/dotStart/HostRoulette/twitch"
  "math/rand"
  "net/http"
)

const WheelRequests = 3
const WheelSelectionPool = 32

type wheelSettings struct {
  Communities []string `json:"communities"`
  Games       []string `json:"games"`
  Languages   []string `json:"languages"`
}

type WheelSelection struct {
  Stream *twitch.Stream `json:"stream"`
  Game   *twitch.Game   `json:"game"`
  User   *twitch.User   `json:"user"`
}

func (s *Server) HandleWheel(w http.ResponseWriter, req *http.Request) {
  if req.Method == "OPTIONS" {
    s.writeHeaders(w)
    return
  }

  if req.Method != "POST" {
    http.NotFound(w, req)
    return
  }

  settings := &wheelSettings{}
  defer req.Body.Close()
  err := json.NewDecoder(req.Body).Decode(settings)
  if err != nil {
    s.logger.Errorf("failed to decode wheel request body: %s", err)
    http.Error(w, "failed to decode request", http.StatusServiceUnavailable)
    return
  }

  streams, err := s.selectFilteredStreams(settings)
  if err != nil {
    s.logger.Errorf("failed to retrieve stream list: %s", err)
    http.Error(w, "failed to retrieve stream list", http.StatusServiceUnavailable)
    return
  }
  s.logger.Debugf("found %d qualifying streams", len(streams))

  if len(streams) == 0 {
    s.writeHeaders(w)
    w.WriteHeader(http.StatusNoContent)
    return
  }

  total := len(streams)
  if total > WheelSelectionPool {
    total = WheelSelectionPool
  }
  streams = streams[len(streams)-total:]

  userIds := make([]string, len(streams))
  gameIds := make([]string, len(streams))
  for i, stream := range streams {
    userIds[i] = stream.UserId

    if stream.GameId != "0" {
      gameIds[i] = stream.GameId
    }
  }

  userPage, err := s.twitchClient.GetUsers(userIds)
  if err != nil {
    s.logger.Errorf("failed to retrieve users for selected streams: %s", err)
    http.Error(w, "failed to retrieve user list", http.StatusServiceUnavailable)
    return
  }
  s.logger.Debugf("selecting from %d streams and %d users", len(streams), len(userPage.Content))

  gameResults, err := s.search.GetGames(gameIds)
  if err != nil {
    s.logger.Errorf("failed to retrieve games for selected streams: %s", err)
    http.Error(w, "failed to retrieve game list", http.StatusServiceUnavailable)
    return
  }

  users := make(map[string]*twitch.User)
  for _, user := range userPage.Content {
    users[user.Id] = user
  }

  games := make(map[string]*twitch.Game)
  for _, game := range gameResults {
    games[fmt.Sprintf("%d", game.Id)] = game
  }

  selection := make([]*WheelSelection, 0)
  for len(selection) < 8 {
    if len(streams) == 0 {
      if len(selection) == 0 {
        s.logger.Debugf("stream list exceeded without alternatives")

        s.writeHeaders(w)
        w.WriteHeader(http.StatusNoContent)
        return
      }

      index := rand.Intn(len(selection))
      selection = append(selection, selection[index])
      continue
    }

    index := rand.Intn(len(streams))
    stream := streams[index]
    user := users[stream.UserId]
    game := games[stream.GameId]
    if user != nil && (stream.GameId != "0" && game != nil) {
      selection = append(selection, &WheelSelection{
        Stream: stream,
        User:   user,
        Game:   game,
      })
    } else {
      s.logger.Debugf("cannot find user \"%s\" or game \"%s\" for stream %s - skipped", stream.UserId, stream.GameId, stream.Id)
    }

    streams = append(streams[:index], streams[index+1:]...)
  }

  s.incrementSpinCount()
  s.writeHeaders(w)
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(&selection)
}

// selects a pre-filtered list of streams using the new API
func (s *Server) selectFilteredStreams(settings *wheelSettings) ([]*twitch.Stream, error) {
  streams := make([]*twitch.Stream, 0)
  var cursor string
  for i := 0; i < WheelRequests; i++ {
    page, err := s.twitchClient.GetStreams(cursor, settings.Communities, settings.Games, settings.Languages)
    if err != nil {
      return nil, err
    }

    cursor = page.Pagination.Cursor

    if len(streams) > 100 {
      streams = append(streams[len(streams)-100:], page.Content...)
    } else {
      streams = append(streams, page.Content...)
    }
    s.logger.Debugf("query yielded %d elements (now using %d in total)", len(page.Content), len(streams))

    if len(page.Content) == 0 || len(page.Content) < 100 {
      break
    }
  }

  return streams, nil
}
