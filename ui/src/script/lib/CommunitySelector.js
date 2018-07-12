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

export default {
  props: ['value'],
  data: function () {
    return {
      initialized: false,
      loading: true,
      knownValues: [],
      selected: [],
    }
  },
  mounted: function () {
    $(this.$el).dropdown({
      maxSelections: 3,
      onAdd: (addedValue, addedText, $addedChoice) => {
        if (this.selected.indexOf(addedValue) === -1) {
          this.selected.push(addedValue);
          this.$emit('input', this.selected);
          this.$emit('change');
        }
      },
      onRemove: (removedValue, removedText, $removedChoice) => {
        let index = this.selected.indexOf(removedValue);
        if (index !== -1) {
          this.selected.splice(index, 1);
          this.$emit('input', this.selected);
          this.$emit('change');
        }
      },
      apiSettings: {
        url: endpoint('/api/search/community/{query}'),
        saveRemoteData: false, // GDPR
      },
    });
  },
  watch: {
    value: function (newValue, oldValue) {
      if (!newValue) {
        newValue = [];
      } else if (newValue === this.selected) {
        return
      }

      this.selected = newValue;
      if (!this.initialized) {
        if (newValue.length !== 0) {
          this.loading = true;
          axios.post(endpoint('/api/community'), this.value).then(
              (response) => {
                this.loading = false;

                let values = [];
                for (let community of response.data) {
                  values.push({
                    name: community.display_name,
                    value: community.id,
                  })
                }

                $(this.$el).dropdown('change values', values);
                $(this.$el).dropdown('set selected', this.selected);
              });
        } else {
          this.loading = false;
          $(this.$el).dropdown('set selected', this.selected);
        }
      }
      this.initialized = true;
    }
  },
  destroy: function () {
    $(this.$el).off();
  },
  template: `
<div :class="'ui fluid search multiple selection dropdown' + (loading ? ' loading' : '')">
  <input type="hidden" />
  <i class="dropdown icon"></i>
  <div class="default text">Search for Communities ...</div>
  <div class="menu"></div>
</div>`,
}
