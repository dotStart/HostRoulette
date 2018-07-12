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

import "github.com/dotStart/HostRoulette/twitch"

func (a *Agent) tickUpdateGames() {
  for range a.gameUpdateTicker.C {
    a.updateGames()
  }
}

// re-indexes the list of games
func (a *Agent) updateGames() {
  var offset uint64 = 0
  for true {
    page, err := a.twitchClient.GetTopGames(offset) // we're only polling for the top 100 atm
    if err != nil {
      a.logger.Errorf("failed to update top communities: %s", err)
      return
    }

    offset += 100
    if offset > page.TotalElements {
      break
    }

    games := make([]*twitch.Game, len(page.Content))
    for i, topGame := range page.Content {
      games[i] = topGame.Game
    }

    err = a.search.UpdateGames(games)
    if err != nil {
      a.logger.Errorf("failed to index top games: %s", err)
    }
  }
}
