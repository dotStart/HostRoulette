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
package cache

import "github.com/go-redis/redis"

func (c *Cache) GetSpinStatistic() (uint64, error) {
  c.logger.Debugf("retrieving spin statistic")
  res, err := c.client.Get("statistic_spins").Uint64()
  if err == redis.Nil {
    return 0, nil
  }
  return res, err
}

func (c *Cache) IncrementSpinStatistic() (uint64, error) {
  c.logger.Debugf("incrementing spin statistic")

  script := redis.NewScript(`
    local current
    current = redis.call("incr",KEYS[1])
    if tonumber(current) == 1 then
      redis.call("expire",KEYS[1],86400)
    end
    return current
  `)

  val, err := script.Run(c.client, []string{"statistic_spins"}).Result()
  if err != nil {
    return 0, err
  }
  return uint64(val.(int64)), nil
}
