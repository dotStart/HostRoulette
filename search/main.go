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
  "fmt"
  "github.com/dotStart/HostRoulette/config"
  "github.com/olivere/elastic"
  "github.com/op/go-logging"
)

type Client struct {
  logger *logging.Logger
  search *elastic.Client
}

func New(cfg *config.SearchConfig) (*Client, error) {
  search, err := elastic.NewClient(elastic.SetURL(cfg.Url))
  if err != nil {
    return nil, fmt.Errorf("cannot establish connection to elasticsearch: %s", err)
  }

  cl := &Client{
    logger: logging.MustGetLogger("search"),
    search: search,
  }

  pong, code, err := search.Ping(cfg.Url).Do(context.Background())
  if err != nil {
    return nil, fmt.Errorf("cannot ping document store: %s", err)
  }
  cl.logger.Infof("Established connection to document store (code: %d)", code)
  cl.logger.Infof("Connected to ElasticSearch v%s of cluster %s (Commit Hash: %s; Timestamp: %s)", pong.Version.Number, pong.ClusterName, pong.Version.BuildHash, pong.Version.BuildTimestamp)

  err = cl.init()
  if err != nil {
    return nil, err
  }

  return cl, nil
}

func (c *Client) init() error {
  err := c.createIndex("community", communityMapping)
  if err != nil {
    return err
  }

  return c.createIndex("game", gameMapping)
}

func (c *Client) createIndex(index string, mapping string) error {
  exists, err := c.search.IndexExists(index).Do(context.Background())
  if err != nil {
    return fmt.Errorf("failed to query for index \"%s\": %s", index, err)
  }

  if !exists {
    c.logger.Infof("Creating index: %s", index)

    request := c.search.CreateIndex(index)
    if mapping != "" {
      request = request.BodyJson(mapping)
    }
    _, err := request.Do(context.Background())
    if err != nil {
      return fmt.Errorf("cannot create index \"%s\": %s", index, err)
    }
  }

  return nil
}
