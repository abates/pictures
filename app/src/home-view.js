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
import '@polymer/paper-button/paper-button.js';
import '@polymer/iron-icon/iron-icon.js';
import '@polymer/iron-ajax/iron-ajax.js';
import './shared-styles.js';
import './thumb-card.js';
import './bread-crumbs.js';

class HomeView extends PolymerElement {
  static get template() {
    return html`
      <style include="shared-styles">
        :host {
          display: block;
          padding: 10px;
        }

        .breadcrumbs {
          font-size: 18px;
          font-weight: bold;
        }

        .grid {
          display: flex;
          flex-wrap: wrap;
          justify-content: center;
        }

        .folder {
          font-size: 24px;
          text-align: center;
          color: var(--app-secondary-color);
          text-decoration: none;
        }

        .folder:hover {
          cursor:pointer;
          color: var(--app-highlight-color);
        }

        .folder-icon {
          width: 200px;
          height: 200px;
        }

        .folder-label {
          position: relative;
          top: -24px;
        }
      </style>
      <iron-ajax auto 
        url="/api/pictures{{path}}/index.json" 
        handle-as="json" 
        last-response="{{pictures}}" 
        debounce-duration="300" 
        on-response="_handleResponse"
        on-error="_handleErrorResponse"></iron-ajax>
      <template is="dom-if" if="[[found]]">
        <bread-crumbs class="breadcrumbs" path="[[path]]"></bread-crumbs>
        <div class="grid">
          <template is="dom-repeat" items="[[pictures]]" as="picture">
            <template is="dom-if" if="[[picture.dir]]">
              <a class="folder" href="[[path]]/[[picture.name]]">
                  <iron-icon class="folder-icon" icon="photo-icons:folder-open"></iron-icon>
                  <div class="folder-label">[[picture.name]]</div>
              </a>
            </template>
            <template is="dom-if" if="[[!picture.dir]]">
              <thumb-card thumb="/api/thumbs[[path]]/[[picture.name]]" src="/api/pictures[[path]]/[[picture.name]]"></thumb-card>
            </template>
          </template>
        </div>
      </template>
      <template is="dom-if" if="[[!found]]">
        <h1>Oops you hit a 404. <a href="[[rootPath]]">Head back to home.</a></h1>
      </template>
    `;
  }

  _handleResponse(event) {
      this.found = true;
  }

  _handleErrorResponse(event) {
    if (event.detail.request.xhr.status == 404) {
      this.found = false;
    }
  }

  _redirect(path) {
    this.dispatchEvent(new CustomEvent('redirect-route', {bubbles: true, composed: true, detail: path }));
  }

  _pathChanged() {
    this.pictures = [];
  }

  static get properties() {
    return {
      pictures: {
        type: Array,
      },

      path: {
        type: String,
        observer: "_pathChanged",
      },

      found: {
        type: Boolean,
        value: true
      }
    }
  }
}

window.customElements.define('home-view', HomeView);
