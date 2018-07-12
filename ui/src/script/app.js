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
import Vue from 'vue';
import CommunitySelector from './lib/CommunitySelector';
import GameSelector from './lib/GameSelector';
import LanguageSelector from './lib/LanguageSelector';
import RouletteWheel from './lib/RouletteWheel';
import {endpoint} from './lib/utility';
import '../style/app.less';

const app = new Vue({
  el: '#roulette',
  components: {
    'roulette-wheel': RouletteWheel,
    'community-selector': CommunitySelector,
    'game-selector': GameSelector,
    'language-selector': LanguageSelector,
  },
  data: {
    // settings
    store: false,
    settings: {
      communities: [],
      games: [],
      languages: [],
    },
    warnings: {
      noFilters: false,
    },

    // statistics
    statistics: {
      spins: 0,
      channels: 0,
      viewers: 0,
    },

    // wheel
    result: null,
  },
  computed: {
    storageAvailable: function () {
      return !!window.localStorage;
    }
  },
  methods: {
    changeStorageSettings: function () {
      if (!this.store) {
        window.localStorage.clear();
        console.log(
            'Storage permission revoked - Cleared all settings from local storage');
        return;
      }

      this.updateSettings();
    },
    updateSettings: function () {
      if (!this.store) {
        return;
      }

      window.localStorage.setItem('settings', JSON.stringify(this.settings));
      console.log('Written settings to local storage');

      this.updateWarningSettings();
    },
    updateWarningSettings: function () {
      if (!this.store) {
        return;
      }

      window.localStorage.setItem('warnings', JSON.stringify(this.warnings));
      console.log('Written warning states to local storage')
    },
    onSpin: function (cb) {
      this.result = null;

      let showWarning = !this.warnings.noFilters;
      if (showWarning) {
        let settingCount = this.settings.communities.length +
            this.settings.games.length +
            this.settings.languages.length;

        if (settingCount === 0) {
          $('#warning-no-filters').modal({
            blurring: true,
            onApprove: () => {
              cb();
            }
          }).modal('show');

          this.warnings.noFilters = true;
          this.updateWarningSettings();
        } else {
          cb();
        }
      } else {
        cb();
      }
    },
    acceptResult: function (result) {
      this.result = result;
    },

    showNoContentWarning: function () {
      $('#warning-no-content').modal({blurring: true}).modal('show');
    }
  }
});

axios.get(endpoint('/api/statistics')).then((response) => {
  app.statistics = response.data;
});

const settings = window.localStorage.getItem('settings');
if (settings != null) {
  console.log('Loading previously stored settings');
  app.store = true;
  app.settings = JSON.parse(settings);
}

const warnings = window.localStorage.getItem('warnings');
if (warnings != null) {
  console.log('Loading previously stored warning states');
  app.warnings = JSON.parse(warnings);
}
