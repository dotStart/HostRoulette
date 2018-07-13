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
package config

import (
  "fmt"
  "github.com/hashicorp/hcl2/hclparse"
)
import "github.com/hashicorp/hcl2/gohcl"

type ServerConfig struct {
  BindAddress string         `hcl:"bind-address,attr"`
  Auth        Authentication `hcl:"auth,block"`
  Cache       CacheConfig    `hcl:"cache,block"`
  Search      SearchConfig   `hcl:"search,block"`
}

// loads a configuration file
func LoadConfig(path string) (*ServerConfig, error) {
  parser := hclparse.NewParser()
  file, diag := parser.ParseHCLFile(path)
  if diag.HasErrors() {
    return nil, fmt.Errorf("failed to parse configuration file: %s", diag.Error())
  }

  cfg := &ServerConfig{}
  diag = gohcl.DecodeBody(file.Body, nil, cfg)
  if diag.HasErrors() {
    return nil, fmt.Errorf("failed to parse configuration file: %s", diag.Error())
  }
  return cfg, nil
}
