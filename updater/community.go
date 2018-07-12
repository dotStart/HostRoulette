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

func (a *Agent) tickUpdateCommunities() {
  for range a.communityUpdateTicker.C {
    a.updateCommunities()
  }
}

// re-indexes the list of top communities
func (a *Agent) updateCommunities() {
  var cursor string
  for i := 0; i < 5; i++ { // we're only indexing the top 500 communities atm
    page, err := a.twitchClient.GetTopCommunities(cursor)
    if err != nil {
      a.logger.Errorf("failed to update top communities: %s", err)
      return
    }

    err = a.search.UpdateCommunities(page.Content)
    if err != nil {
      a.logger.Errorf("failed to index top communities: %s", err)
    }

    cursor = page.Cursor
  }
}
