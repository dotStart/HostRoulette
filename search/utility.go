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

import "github.com/olivere/elastic"

type Result struct {
  Success bool     `json:"success"`
  Matches []*Match `json:"results"`
}

type Match struct {
  Name  string `json:"name"`
  Value string `json:"value"`
}

func NewSearchResult(result *elastic.SearchResult, key string) *Result {
  suggestions := result.Suggest[key]
  if len(suggestions) == 0 {
    return &Result{
      Success: true,
      Matches: []*Match{},
    }
  }

  suggestion := suggestions[0]
  matches := make([]*Match, len(suggestion.Options))
  for i, option := range suggestion.Options {
    matches[i] = &Match{
      Name:  option.Text,
      Value: option.Id,
    }
  }

  return &Result{
    Success: true,
    Matches: matches,
  }
}
