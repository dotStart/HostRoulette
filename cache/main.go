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
  "crypto/sha1"
  "encoding/hex"
  "fmt"
  "github.com/dotStart/HostRoulette/config"
  "github.com/go-redis/redis"
  "github.com/op/go-logging"
)

type Cache struct {
  logger *logging.Logger
  client *redis.Client
}

func New(cfg *config.CacheConfig) (*Cache, error) {
  var pw string
  if cfg.Password != nil {
    pw = *cfg.Password
  }

  client := redis.NewClient(&redis.Options{
    Addr:     cfg.Address,
    Password: pw,
    DB:       cfg.Database,
  })

  err := client.Ping().Err()
  if err != nil {
    return nil, fmt.Errorf("failed to establish connection to redis server: %s", err)
  }

  return &Cache{
    logger: logging.MustGetLogger("cache"),
    client: client,
  }, nil
}

func calculateHash(key string) string {
  data := []byte(key)
  hash := sha1.Sum(data)
  return hex.EncodeToString(hash[:])
}
