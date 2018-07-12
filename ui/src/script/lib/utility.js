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

/**
 * Computes the actual location of an arbitrary endpoint.
 *
 * Specifically, this method respects the location of a test server in
 * development mode.
 *
 * @param path an arbitrary endpoint path (beginning with a slash character).
 * @returns string an actual endpoint url.
 */
function endpoint(path) {
  if (typeof SERVER_ADDR !== "undefined") {
    return SERVER_ADDR + path;
  }

  return path;
}

export {
  endpoint
}
