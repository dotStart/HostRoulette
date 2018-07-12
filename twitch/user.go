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

type User struct {
  Id           string `json:"id"`
  Login        string `json:"login"`
  DisplayName  string `json:"display_name"`
  Description  string `json:"description"`
  ProfileImage string `json:"profile_image_url"`
}

type UserPage struct {
  Content []*User `json:"data"`
}

func (c *Client) GetUsers(userIds []string) (*UserPage, error) {
  params := params()
  for _, id := range userIds {
    params.add("id", id)
  }

  page := &UserPage{}
  err := c.executeNew("GET", "/users"+params.String(), nil, page)
  return page, err
}
