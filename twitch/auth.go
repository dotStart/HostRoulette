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
)

type AppToken struct {
  AccessToken string `json:"access_token"`
}

func (c *Client) tickRequestAppToken() {
  for range c.appTokenTicker.C {
    err := c.requestAppToken()
    if err != nil {
      c.logger.Errorf("failed to refresh app token: %s", err)
    }
  }
}

func (c *Client) getAppToken() string {
  c.appTokenLock.RLock()
  defer c.appTokenLock.RUnlock()

  return c.appToken
}

func (c *Client) requestAppToken() error {
  c.logger.Debugf("requesting a new app token")

  params := params()
  params.add("client_id", c.cfg.ClientId)
  params.add("client_secret", c.cfg.ClientSecret)
  params.add("grant_type", "client_credentials")

  token := &AppToken{}
  req, err := c.execute("POST", "https://id.twitch.tv/oauth2/token"+params.String(), nil)
  if err != nil {
    return err
  }
  res, err := c.httpClient.Do(req)
  if err != nil {
    return err
  }
  err = checkResponse(res)
  if err != nil {
    return err
  }

  defer res.Body.Close()
  err = json.NewDecoder(res.Body).Decode(token)
  if err != nil {
    return err
  }

  c.appTokenLock.Lock()
  defer c.appTokenLock.Unlock()

  c.appToken = token.AccessToken
  c.logger.Debugf("acquired app token: %s", c.appToken)
  return nil
}
