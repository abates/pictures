// Copyright 2018 Andrew Bates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import { PolymerElement, html } from '@polymer/polymer/polymer-element.js';
import '@polymer/iron-ajax/iron-ajax.js';
import './shared-styles.js';
import './thumb-card.js';

class HomeView extends PolymerElement {
  static get template() {
    return html`
      <style include="shared-styles">
        :host {
          display: block;
          padding: 10px;
        }

        .grid {
          display: flex;
          flex-wrap: wrap;
          justify-content: center;
        }
      </style>
      <iron-ajax auto url="/api/photos/index.json" handle-as="json" last-response="{{photos}}" debounce-duration="300"></iron-ajax>
      <div class="grid">
        <template is="dom-repeat" items="[[photos]]" as="photo">
          <thumb-card photo="[[photo]]"></thumb-card>
        </template>
      </div>
    `;
  }


}

window.customElements.define('home-view', HomeView);
