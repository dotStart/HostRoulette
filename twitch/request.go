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

import "net/url"

type requestParameters struct {
  params map[string][]string
}

func params() *requestParameters {
  return &requestParameters{
    params: make(map[string][]string),
  }
}

func (p *requestParameters) add(name string, value string) *requestParameters {
  valueMap := p.params[name]
  if valueMap == nil {
    valueMap = make([]string, 0)
  }

  valueMap = append(valueMap, value)
  p.params[name] = valueMap

  return p
}

func (p *requestParameters) remove(name string) *requestParameters {
  delete(p.params, name)
  return p
}

func (p *requestParameters) String() string {
  enc := ""

  for key, values := range p.params {
    for _, val := range values {
      if len(enc) == 0 {
        enc += "?"
      } else {
        enc += "&"
      }

      enc += url.QueryEscape(key)
      enc += "="
      enc += url.QueryEscape(val)
    }
  }

  return enc
}
