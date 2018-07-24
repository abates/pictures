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
import '@vaadin/vaadin-lumo-styles/color.js';
import '@vaadin/vaadin-date-picker/vaadin-date-picker.js';

class PhotoCard extends PolymerElement {
  static get template() {
    return html`
      <style include="lumo-color">
        :host {
          position: fixed;
          overflow: hidden;
          opacity: 0;
          transition: all 0.5s linear;
        }

        img {
          display: inline;
          transition: all 0.2s linear;
        }

        .close-icon {
          position: absolute;
          color: white;
          top: 0;
          right: 0;
        }

        .date-picker {
          --lumo-base-color: transparent;
          --lumo-contrast-10pct: transparent;
          position: absolute;
          bottom: 10px;
          left: 10px;
        }
      </style>

      <paper-icon-button class="close-icon" icon="photo-icons:close" on-click="close"></paper-icon-button>
      <img id="photo" src$="[[src]]" on-load="loaded"></img>
      <vaadin-date-picker theme="dark" class="date-picker"></vaadin-date-picker>
    `;
  }

  constructor() {
    super();
    window.addEventListener("resize", e => this.resize());
    this.addEventListener("iron-overlay-closed", e => this.onClose());
    this.withBackdrop = true;
  }

  static get properties() {
    return {
      src: String,
    }
  }

  close() {
    this.style.opacity = "0";
    this.dispatchEvent(new CustomEvent('photo-card-closed', {detail: {closed: true}}));
  }

  loaded() {
    this.resize();
    this.style.zIndex = 200;
    this.style.opacity = "1.0";
  }

  resize() {
    var minWidth = window.innerWidth - 20;
    var minHeight = window.innerHeight - 20;
    var width = this.$.photo.naturalWidth;
    var height = this.$.photo.naturalHeight;

    if (minWidth < this.$.photo.naturalWidth) {
      var scale = minWidth / this.$.photo.naturalWidth;
      width = minWidth;
      height = this.$.photo.naturalHeight * scale;
    } else if (minHeight < this.$.photo.naturalHeight) {
      var scale = minHeight / this.$.photo.naturalHeight;
      width = this.$.photo.naturalWidth * scale;
      height = minHeight;
    } 

    this.$.photo.setAttribute("width", width)
    this.$.photo.setAttribute("height", height)

    var top = (window.innerHeight/2) - (height/2);
    var left = (window.innerWidth/2) - (width/2);
    this.style.top = top + "px";
    this.style.left = left + "px";
  }
}

window.customElements.define('photo-card', PhotoCard);
