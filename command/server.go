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
package command

import (
  "context"
  "flag"
  "fmt"
  "github.com/dotStart/HostRoulette/config"
  "github.com/dotStart/HostRoulette/search"
  "github.com/dotStart/HostRoulette/server"
  "github.com/dotStart/HostRoulette/twitch"
  "github.com/dotStart/HostRoulette/ui"
  "github.com/dotStart/HostRoulette/updater"
  "github.com/google/subcommands"
  "github.com/op/go-logging"
  "net"
  "net/http"
  "os"
)

type ServerCommand struct {
  logger *logging.Logger

  flagLogLevel   string
  flagDevMode    bool
  flagConfigFile string
}

func (s *ServerCommand) Name() string {
  return "server"
}

func (s *ServerCommand) Synopsis() string {
  return "starts a new server instance"
}

func (s *ServerCommand) Usage() string {
  return `Usage: host-roulette server [options]

This command starts a new Host Roulette server:

  $ host roulette server

Available command specific flags:

`
}

func (s *ServerCommand) SetFlags(f *flag.FlagSet) {
  f.StringVar(&s.flagLogLevel, "log-level", "info", "selects the server log level")
  f.BoolVar(&s.flagDevMode, "dev", false, "en- or disables development mode (e.g. permits API access from localhost)")
  f.StringVar(&s.flagConfigFile, "config-file", "", "selects the location of a server configuration file")
}

func (s *ServerCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
  level, err := logging.LogLevel(s.flagLogLevel)
  if err != nil {
    fmt.Fprintf(os.Stderr, "error: Illegal log level \"%s\": %s", s.flagLogLevel, err)
    return 1
  }

  if s.flagConfigFile == "" {
    fmt.Fprintf(os.Stderr, "error: must specify a configuration file")
    return 1
  }
  cfg, err := config.LoadConfig(s.flagConfigFile)
  if err != nil {
    fmt.Fprintf(os.Stderr, "error: cannot load configuration: %s", err)
    return 1;
  }

  listener, err := net.Listen("tcp", cfg.BindAddress)
  if err != nil {
    fmt.Fprintf(os.Stderr, "error: cannot bind on address %s: %s", cfg.BindAddress, err)
    return 1
  }

  fmt.Printf("==> Host Roulette Configuration\n\n")
  fmt.Printf("      Bind Address: %s\n", cfg.BindAddress)
  fmt.Printf("         Log Level: %s\n", s.flagLogLevel)
  fmt.Printf("  Development Mode: %v\n", s.flagDevMode)
  fmt.Printf("           Version: %s\n", versionFull())
  fmt.Printf("       Commit Hash: %s\n\n", commitHash)

  fmt.Printf("==> Starting Server\n\n")

  format := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} [%{level:.4s}] %{module} : %{color:reset} %{message}`)
  backend := logging.AddModuleLevel(logging.NewBackendFormatter(logging.NewLogBackend(os.Stdout, "", 0), format))
  backend.SetLevel(level, "")
  logging.SetBackend(backend)
  s.logger = logging.MustGetLogger("server")

  srch, err := search.New(&cfg.Search)
  if err != nil {
    fmt.Fprintf(os.Stderr, "error: failed to connect to elasticsearch: %s", err)
    return 1
  }
  twitchClient, err := twitch.NewClient(&cfg.Auth)
  if err != nil {
    fmt.Fprintf(os.Stderr, "error: failed to establish connection with Twitch: %s", err)
    return 1
  }
  updt := updater.New(srch, twitchClient)
  defer updt.Close()

  mux := http.NewServeMux()
  srv := &http.Server{
    Handler: mux,
  }
  defer srv.Close()

  api := server.New(mux, srch, twitchClient)
  if s.flagDevMode {
    api.CorsDisabled = true
    s.logger.Debugf("development mode has been enabled")
  }

  ui.Register(mux)

  s.logger.Infof("Listening for requests on %s", cfg.BindAddress)
  err = srv.Serve(listener)
  if err != nil {
    fmt.Fprintf(os.Stderr, "error: failed to initialize http server: %s", err)
  }
  return 0
}
