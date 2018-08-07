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

class BreadCrumbs extends PolymerElement {
  static get template() {
    return html`
      <style>
        a {
          all: inherit;
          color: var(--app-secondary-color);
        }

        a:hover {
          all: inherit;
          cursor: pointer;
          color: var(--app-highlight-color);
        }
      </style>

      <a href="/">/</a>
      <template is="dom-repeat" items="[[breadcrumbs]]" as="breadcrumb">
        <a href="[[breadcrumb.path]]">[[breadcrumb.name]]&nbsp;/</a>
      </template>
    `;
  }

  _pathChanged(newValue) {
    var path = "/";
    this.set("breadcrumbs", []);
    newValue.split("/").forEach(b => {
      if (b != "") {
        if (path != "/") {
          path += "/";
        }
        path += b;
        this.push("breadcrumbs", {path: path, name: b});
      }
    });
  }

  static get properties() {
    return {
      path: {
        type: String,
        observer: "_pathChanged"
      }
    }
  }
}

window.customElements.define('bread-crumbs', BreadCrumbs);
