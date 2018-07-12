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
package twitch

import (
  "encoding/json"
  "fmt"
  "github.com/dotStart/HostRoulette/config"
  "github.com/op/go-logging"
  "io"
  "io/ioutil"
  "net/http"
  "sync"
  "time"
)

const oldApiBaseUrl = "https://api.twitch.tv/kraken"
const newApiBaseUrl = "https://api.twitch.tv/helix"

type Client struct {
  logger     *logging.Logger
  httpClient *http.Client
  cfg        *config.Authentication

  appTokenLock   sync.RWMutex
  appToken       string
  appTokenTicker *time.Ticker
}

func NewClient(cfg *config.Authentication) (*Client, error) {
  cl := &Client{
    logger:     logging.MustGetLogger("twitch"),
    httpClient: &http.Client{},
    cfg:        cfg,

    appTokenTicker: time.NewTicker(30 * 24 * time.Hour),
  }
  err := cl.requestAppToken()
  if err != nil {
    return nil, fmt.Errorf("failed to acquire app token: %s", err)
  }

  go cl.tickRequestAppToken()

  return cl, nil
}

func (c *Client) execute(method string, url string, body io.Reader) (*http.Request, error) {
  req, err := http.NewRequest(method, url, body)
  if err != nil {
    return nil, err
  }

  c.logger.Debugf("sending request: %s %s", req.Method, req.URL)
  req.Header.Set("User-Agent", "Host Roulette (+https://host-roulette.dotstart.tv)")

  return req, nil
}

// executes an arbitrary request
func (c *Client) executeOld(method string, endpoint string, body io.Reader, output interface{}) error {
  req, err := c.execute(method, oldApiBaseUrl+endpoint, body)
  if err != nil {
    return err
  }

  req.Header.Set("Client-ID", c.cfg.ClientId)
  req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")

  res, err := c.httpClient.Do(req)
  if err != nil {
    return err
  }

  err = checkResponse(res)
  if err != nil {
    return err
  }

  defer res.Body.Close()
  return json.NewDecoder(res.Body).Decode(output)
}

func (c *Client) executeNew(method string, endpoint string, body io.Reader, output interface{}) error {
  req, err := c.execute(method, newApiBaseUrl+endpoint, body)
  if err != nil {
    return err
  }

  req.Header.Set("Authorization", "Bearer "+c.getAppToken())

  res, err := c.httpClient.Do(req)
  if err != nil {
    return err
  }

  err = checkResponse(res)
  if err != nil {
    return err
  }

  defer res.Body.Close()
  return json.NewDecoder(res.Body).Decode(output)
}

func checkResponse(res *http.Response) error {
  category := res.StatusCode / 100

  if category == 2 {
    return nil
  }

  var message string
  defer res.Body.Close()
  messageBytes, err := ioutil.ReadAll(res.Body)
  if err != nil {
    message = string(messageBytes)
  }

  switch category {
  case 4:
    return fmt.Errorf("illegal response: client error (code %d) - %s", res.StatusCode, message)
  case 5:
    return fmt.Errorf("illegal response: server error (code %d) - %s", res.StatusCode, message)
  }

  return fmt.Errorf("illegal response: unknown error %d of category %d - %s", res.StatusCode, category, message)
}
