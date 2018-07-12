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
package updater

import (
  "github.com/dotStart/HostRoulette/search"
  "github.com/dotStart/HostRoulette/twitch"
  "github.com/op/go-logging"
  "time"
)

type Agent struct {
  logger       *logging.Logger
  twitchClient *twitch.Client
  search       *search.Client

  communityUpdateTicker *time.Ticker
  gameUpdateTicker *time.Ticker
}

func New(search *search.Client, twitchClient *twitch.Client) *Agent {
  a := &Agent{
    logger:       logging.MustGetLogger("updater"),
    twitchClient: twitchClient,
    search:       search,

    communityUpdateTicker: time.NewTicker(24 * time.Hour),
    gameUpdateTicker: time.NewTicker(3 * 24 * time.Hour),
  }

  a.logger.Infof("Performing initial update of all cached records (this may take a while)")
  a.updateCommunities()
  a.updateGames()

  go a.tickUpdateCommunities()
  go a.tickUpdateGames()
  return a
}

func (a *Agent) Close() {
  a.communityUpdateTicker.Stop()
  a.gameUpdateTicker.Stop()
}
