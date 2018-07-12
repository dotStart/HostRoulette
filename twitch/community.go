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
  "net/url"
)

type Community struct {
  Id          string `json:"_id"`
  Name        string `json:"name"`
  DisplayName string `json:"display_name"`
}

type CommunityPage struct {
  OldPage
  Content []*Community `json:"communities"`
}

// queries the servers for a list of top communities
func (c *Client) GetTopCommunities(cursor string) (*CommunityPage, error) {
  c.logger.Debugf("requesting top communities from offset \"%s\"", cursor)

  page := &CommunityPage{}
  err := c.executeOld("GET", fmt.Sprintf("/communities/top?limit=100&cursor=%s", cursor), nil, &page)
  if err != nil {
    return nil, err
  }
  return page, nil
}

// queries the servers for a specific community based on its name
func (c *Client) GetCommunity(name string) (*Community, error) {
  community := &Community{}
  err := c.executeOld("GET", fmt.Sprintf("/communities?name=%s", url.QueryEscape(name)), nil, community)
  if err != nil {
    return nil, err
  }
  return community, nil
}

func (c *Community) MarshalJSON() ([]byte, error) {
  return json.Marshal(&struct {
    Id          string `json:"id"`
    Name        string `json:"name"`
    DisplayName string `json:"display_name"`
  }{
    Id:          c.Id,
    Name:        c.Name,
    DisplayName: c.DisplayName,
  })
}
