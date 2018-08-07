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
import '@polymer/paper-icon-button/paper-icon-button.js';
import './photo-backdrop.js';

class ThumbCard extends PolymerElement {
  static get template() {
    return html`
      <style>
        :host {
          display: block;
          margin: 4px;
          overflow: hidden;
          border: 2px solid rgba(0,0,0,0);
        }

        :host(:hover) {
          border: 2px solid var(--app-highlight-color);
        }

        :host(:hover) img {
          transform: scale(1.1);
        }

        img {
          transition: all 0.2s linear;
          margin: 0;
          border: 0;
          padding: 0;
        }

      </style>

      <img id="thumbnail" src="[[thumbSrc]]" on-load="thumbnailLoaded" on-click="showOverlay"></img>
      <photo-backdrop id="overlay" src="[[src]]" on-click="hideOverlay"></photo-backdrop>
    `;
  }

  _thumbUpdated(newValue, oldValue) {
    if (this.thumbSrc == undefined) {
      this.thumbSrc = newValue;
    }
  }

  static get properties() {
    return {
      src: String,
      thumb: {
        type: String,
        observer: "_thumbUpdated",
      }
    }
  }

  constructor() {
    super();
  }

  thumbnailLoaded() {
    this.$.thumbnail.setAttribute("height", 200);
    var width = this.$.thumbnail.naturalWidth * 200 / this.$.thumbnail.naturalHeight;
    this.setAttribute("style", "width: " + width + "px; height: 200px;");
  }

  showOverlay() {
    this.$.overlay.open();
  }

  hideOverlay() {
    this.$.overlay.close();
  }
}

window.customElements.define('thumb-card', ThumbCard);
