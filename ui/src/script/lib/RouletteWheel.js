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
import axios from 'axios';
import {endpoint} from './utility';

/**
 * Defines the default stream which is substituted before an actual stream is
 * chosen by the wheel.
 */
const DEFAULT_STREAM = {
  stream: {
    viewer_count: '',
  },
  game: {
    name: '',
  },
  user: {
    login: '',
    display_name: '',
  },

  game: '',
  viewers: 0,
  channel: {
    _id: '',
    name: '',
    display_name: '',
    status: '',
    mature: false,
    partner: false,
  },
  preview: {
    small: '',
    medium: '',
    large: '',
  }
};

/**
 * Defines a multiplier which is applied to the wheel's velocity once per frame.
 * @type {number}
 */
const VELOCITY_MULTIPLIER = 0.99;

/**
 * Defines the minimum permitted velocity for the wheel (e.g. the wheel is
 * considered stopped if the velocity drops below this value).
 * @type {number}
 */
const VELOCITY_MINIMUM = 0.1;

/**
 * Defines the minimum amount of spin that is given to the wheel when its spin
 * button is pressed.
 * @type {number}
 */
const VELOCITY_SPIN_MINIMUM = 15;

/**
 * Defines the range from which the target velocity is selected.
 * @type {number}
 */
const VELOCITY_SPIN_RANGE = 10;

export default {
  props: ['settings'],
  data: function () {
    return {
      renderInterval: null,

      demo: true,
      spinning: false,
      loading: false,
      hasResult: false,
      angle: 0,
      velocity: .5,

      result: DEFAULT_STREAM,
      options: []
    }
  },
  mounted: function () {
    this.renderInterval = window.setInterval(() => {
      if (!this.spinning && !this.demo) {
        return;
      }

      this.angle += this.velocity;
      if (this.demo) {
        return;
      }

      this.velocity *= VELOCITY_MULTIPLIER;

      if (this.velocity < VELOCITY_MINIMUM) {
        this.velocity = 0;

        let index = Math.round(-this.angle / 45) % 8;
        if (index < 0) {
          index += 8;
        }

        this.spinning = false;
        this.hasResult = true;
        this.result = this.options[index];

        this.$emit('result', this.result);
      }
    }, 16) // 1000 ms / 60 frames
  },
  destroyed: function () {
    window.clearInterval(this.renderInterval);
  },
  methods: {
    spin: function () {
      this.$emit('spin', () => {
        this.loading = true;

        axios.post(endpoint('/api/wheel'), this.settings).then((response) => {
          this.loading = false;
          if (response.status === 204) {
            this.$emit('no-content');
            return
          }

          this.options = response.data;

          for (let option of this.options) {
            let url = option.stream.thumbnail_url;
            url = url.replace('{width}', 320);
            url = url.replace('{height}', 180);
            option.stream.thumbnail_url = url;
          }

          this.demo = false;
          this.spinning = true;
          this.velocity = Math.random() * VELOCITY_SPIN_RANGE
              + VELOCITY_SPIN_MINIMUM;

          this.$emit('count', parseInt(response.headers['x-spin-count']));
        }).catch((error) => {
          this.loading = false;

          if (!error.response) {
            this.$emit('error', 'unknown');
            return
          }

          let errorType = 'unknown';
          switch (error.response.status) {
            case 429:
              errorType = 'rate-limit';
              break;
          }
          this.$emit('error', errorType);
        });
      });
    }
  },
  template: `
<div id="roulette-wheel" :class="!spinning ? 'finished' : ''">
  <ul :style="{ transform: 'rotate(' + angle + 'deg)' }">
    <li v-if="!demo" v-for="option in options" :class="'option' + (option === result ? ' active' : '')">
      <div class="preview" :style="{ backgroundImage: 'url(' + option.stream.thumbnail_url + ')' }"></div>
    </li>
    <li v-if="demo" v-for="i in 8" class="option">
      <div class="demo preview"></div>
    </li>
  </ul>

  <section :class="'result' + (!spinning && hasResult ? ' active' : '')">
    <h3 class="ui inverted header">
      <a :href="'https://twitch.tv/' + result.user.login" target="_blank">{{ result.user.display_name }}</a>
      <div class="sub header" v-if="result.game != null">Streaming {{ result.game.name }}</div>
    </h3>

    <div class="ui horizontal inverted list">
      <div class="item" :title="result.viewers + ' viewers'">
        <i class="eye icon"></i>
      {{ result.stream.viewer_count }}
      </div>
    </div>
  </section>

  <a class="spin" href="javascript:void(0)" v-on:click="spin">
    <template v-if="!loading">Spin</template>
    <template v-if="loading"><i class="notched circle loading icon"></i></template>
  </a>
</div>`
}
