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
import {mixinBehaviors} from '@polymer/polymer/lib/legacy/class.js';
import {IronOverlayBehavior} from '@polymer/iron-overlay-behavior/iron-overlay-behavior.js';
import './photo-card.js';

class PhotoBackdrop extends mixinBehaviors(IronOverlayBehavior, PolymerElement) {
  static get template() {
    return html`
      <style>
        :host {
        }
      </style>
    `;
  }

  constructor() {
    super();
    this.addEventListener("iron-overlay-closed", e => this.onClose());
    this.withBackdrop = true;
    this.noCancelOnOutsideClick = true;
  }

  static get properties() {
    return {
      src: String,
    }
  }

  open() {
    super.open();
    this.photoCard = document.createElement("photo-card");
    this.photoCard.src = this.src;
    this.photoCard.addEventListener("photo-card-closed", e => this.close());
    document.body.append(this.photoCard);
  }

  close() {
    super.close();
  }

  onClose() {
    this.photoCard.remove();
  }
}

window.customElements.define('photo-backdrop', PhotoBackdrop);
