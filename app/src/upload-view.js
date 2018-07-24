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
import '@vaadin/vaadin-upload/vaadin-upload.js'
import '@vaadin/vaadin-lumo-styles/color.js';
import './shared-styles.js';

class UploadView extends PolymerElement {
  static get template() {
    return html`
      <style include="shared-styles lumo-color">
        :host {
          display: block;

          padding: 10px;
        }
      </style>
      <div theme="dark">
        <vaadin-upload accept="image/*" target="/api/pictures" method="POST"></vaading-upload>
      </div>
    `;
  }
}

window.customElements.define('upload-view', UploadView);
