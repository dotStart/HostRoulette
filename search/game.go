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
package search

import (
  "context"
  "encoding/json"
  "fmt"
  "github.com/dotStart/HostRoulette/twitch"
  "github.com/olivere/elastic"
  "strconv"
  "strings"
)

func (c *Client) UpdateGames(games []*twitch.Game) error {
  c.logger.Debugf("indexing %d games", len(games))

  request := c.search.Bulk().Index("game").Type("doc")
  for _, game := range games {
    request.Add(elastic.NewBulkIndexRequest().Id(fmt.Sprintf("%d", game.Id)).Doc(game))
  }
  result, err := request.Do(context.Background())
  if err != nil {
    return err
  }

  for _, failed := range result.Failed() {
    c.logger.Errorf("index query for game %s has failed: %s", failed.Id, failed.Error.Reason)
  }
  return nil
}

func (c *Client) GetGames(ids []string) ([]*twitch.Game, error) {
  c.logger.Debugf("retrieving games from database: %s", strings.Join(ids, ", "))

  req := c.search.MultiGet()
  for _, id := range ids {
    req.Add(elastic.NewMultiGetItem().Index("game").Id(id))
  }
  res, err := req.Do(context.Background())
  if err != nil {
    return nil, err
  }

  games := make([]*twitch.Game, 0)
  for _, doc := range res.Docs {
    if doc.Source == nil {
      c.logger.Debugf("could not find game \"%s\"", doc.Id)
      continue
    }

    game := &twitch.Game{}
    err := json.Unmarshal(*doc.Source, &game)
    if err != nil {
      return nil, fmt.Errorf("failed to decode game %s: %s", doc.Id, err)
    }

    id, err := strconv.ParseUint(doc.Id, 10, 64)
    if err != nil {
      return nil, err
    }

    game.Id = id
    games = append(games, game)
  }
  return games, nil
}

func (c *Client) SearchGame(query string) (*Result, error) {
  c.logger.Debugf("searching database for game matching query \"%s\"", query)
  q := elastic.NewCompletionSuggester("name").Field("name").Prefix(query)

  result, err := c.search.Search("game").Suggester(q).Do(context.Background())
  if err != nil {
    return nil, err
  }

  return NewSearchResult(result, "name"), nil
}
