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
  "strings"
)

func (c *Client) UpdateCommunities(communities []*twitch.Community) error {
  c.logger.Debugf("indexing %d communities", len(communities))

  request := c.search.Bulk().Index("community").Type("doc")
  for _, community := range communities {
    request.Add(elastic.NewBulkIndexRequest().Id(community.Id).Doc(community))
  }
  result, err := request.Do(context.Background())
  if err != nil {
    return err
  }

  for _, failed := range result.Failed() {
    c.logger.Errorf("index query for community %s has failed: %s", failed.Id, failed.Error.Reason)
  }
  return nil
}

func (c *Client) GetCommunities(ids []string) ([]*twitch.Community, error) {
  c.logger.Debugf("retrieving communities from database: %s", strings.Join(ids, ", "))

  req := c.search.MultiGet()
  for _, id := range ids {
    req.Add(elastic.NewMultiGetItem().Index("community").Id(id))
  }
  res, err := req.Do(context.Background())
  if err != nil {
    return nil, err
  }

  communities := make([]*twitch.Community, len(res.Docs))
  for i, doc := range res.Docs {
    community := &twitch.Community{}
    err := json.Unmarshal(*doc.Source, &community)
    if err != nil {
      return nil, fmt.Errorf("failed to decode community %s: %s", doc.Id, err)
    }
    community.Id = doc.Id
    communities[i] = community
  }
  return communities, nil
}

func (c *Client) SearchCommunity(query string) (*Result, error) {
  c.logger.Debugf("searching database for community matching query \"%s\"", query)
  q := elastic.NewCompletionSuggester("name").Field("display_name").Prefix(query)

  result, err := c.search.Search("community").Suggester(q).Do(context.Background())
  if err != nil {
    return nil, err
  }

  return NewSearchResult(result, "name"), nil
}
