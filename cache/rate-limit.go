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

import (
  "fmt"
  "github.com/go-redis/redis"
  "time"
)

// retrieves the rate limit usage for the given remote address
func (c *Cache) GetRateLimitUsage(addr string) (uint64, error) {
  c.logger.Debugf("retrieving rate limit usage for address \"%s\"", addr)

  value, err := c.client.Get(fmt.Sprintf("rate_limit_%s", calculateHash(addr))).Uint64()
  if err == redis.Nil {
    return 0, nil
  }
  return value, err
}

// increments teh rate limit usage for the given address
func (c *Cache) IncrementRateLimitUsage(addr string, ttl time.Duration) error {
  c.logger.Debugf("incrementing rate limit usage for address \"%s\"", addr)

  key := fmt.Sprintf("rate_limit_%s", calculateHash(addr))
  err := c.client.Incr(key).Err()
  if err != nil {
    return err
  }

  return c.client.Expire(key, ttl).Err()
}
