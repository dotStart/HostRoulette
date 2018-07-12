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

export default {
  props: ['value'],
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
    });
  },
  data: function () {
    return {
      selected: [],
    }
  },
  watch: {
    value: function (newValue, oldValue) {
      if (!newValue) {
        newValue = [];
      } else if (newValue === this.selected) {
        return
      }

      this.selected = newValue;
      $(this.$el).dropdown('set selected', this.selected);
      this.initialized = true;
    }
  },
  template: `
<div class="ui fluid search multiple selection dropdown">
  <input type="hidden" />
  <i class="dropdown icon"></i>
  <div class="default text">Search for Languages ...</div>
  <div class="menu">
    <div class="item" data-value="en"><i class="gb flag"></i> English</div>
    <div class="item" data-value="da"><i class="dk flag"></i> Dansk</div>
    <div class="item" data-value="de"><i class="de flag"></i> Deutsch</div>
    <div class="item" data-value="es"><i class="es flag"></i> Español</div>
    <div class="item" data-value="fr"><i class="fr flag"></i> Français</div>
    <div class="item" data-value="it"><i class="it flag"></i> Italiano</div>
    <div class="item" data-value="hu"><i class="hu flag"></i> Magyar</div>
    <div class="item" data-value="nl"><i class="nl flag"></i> Nederlands</div>
    <div class="item" data-value="no"><i class="no flag"></i> Norsk</div>
    <div class="item" data-value="pl"><i class="pl flag"></i> Polski</div>
    <div class="item" data-value="pt"><i class="pt flag"></i> Português</div>
    <div class="item" data-value="ro"><i class="ro flag"></i> Română</div>
    <div class="item" data-value="sk"><i class="sk flag"></i> Slovenčina</div>
    <div class="item" data-value="fi"><i class="fi flag"></i> Suomi</div>
    <div class="item" data-value="sv"><i class="sv flag"></i> Svenska</div>
    <div class="item" data-value="vi"><i class="vi flag"></i> Tiếng Việt</div>
    <div class="item" data-value="tr"><i class="tr flag"></i> Türkçe</div>
    <div class="item" data-value="cs"><i class="cs flag"></i> Čeština</div>
    <div class="item" data-value="el"><i class="gr flag"></i> Ελληνικά</div>
    <div class="item" data-value="bg"><i class="bg flag"></i> Български</div>
    <div class="item" data-value="ru"><i class="ru flag"></i> Русский</div>
    <div class="item" data-value="ar"><i class="ar flag"></i> العربية</div>
    <div class="item" data-value="th"><i class="th flag"></i> ภาษาไทย</div>
    <div class="item" data-value="zh"><i class="cn flag"></i> 中文</div>
    <div class="item" data-value="zh-hk"><i class="cn flag"></i> 中文(粵語)</div>
    <div class="item" data-value="ja"><i class="jp flag"></i> 日本語</div>
    <div class="item" data-value="ko"><i class="kr flag"></i> 한국어</div>
    <div class="item" data-value="asl"><i class="sign language icon"></i> American Sign Language</div>
  </div>
</div>`
};
