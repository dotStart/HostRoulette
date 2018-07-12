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

type Stream struct {
  Id          string `json:"id"`
  UserId      string `json:"user_id"`
  GameId      string `json:"game_id"`
  Title       string `json:"title"`
  ViewerCount uint64 `json:"viewer_count"`
  Thumbnail   string `json:"thumbnail_url"`
}

type StreamPage struct {
  NewPage
  Content []*Stream `json:"data"`
}

func (c *Client) GetStreams(cursor string, communityIds []string, gameIds []string, languages []string) (*StreamPage, error) {
  params := params()
  params.add("first", "100")
  if cursor != "" {
    params.add("after", cursor)
  }
  for _, communityId := range communityIds {
    params.add("community_id", communityId)
  }
  for _, gameId := range gameIds {
    params.add("game_id", gameId)
  }
  for _, language := range languages {
    params.add("language", language)
  }

  page := &StreamPage{}
  err := c.executeNew("GET", "/streams"+params.String(), nil, page)
  return page, err
}
