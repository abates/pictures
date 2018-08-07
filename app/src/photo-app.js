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
import { setPassiveTouchGestures, setRootPath } from '@polymer/polymer/lib/utils/settings.js';
import '@polymer/app-layout/app-drawer/app-drawer.js';
import '@polymer/app-layout/app-drawer-layout/app-drawer-layout.js';
import '@polymer/app-layout/app-header/app-header.js';
import '@polymer/app-layout/app-header-layout/app-header-layout.js';
import '@polymer/app-layout/app-scroll-effects/app-scroll-effects.js';
import '@polymer/app-layout/app-toolbar/app-toolbar.js';
import '@polymer/app-route/app-location.js';
import '@polymer/app-route/app-route.js';
import '@polymer/iron-pages/iron-pages.js';
import '@polymer/iron-selector/iron-selector.js';
import '@polymer/paper-dialog/paper-dialog.js';
import '@polymer/paper-dialog-scrollable/paper-dialog-scrollable.js';
import '@polymer/paper-icon-button/paper-icon-button.js';
import '@polymer/iron-icon/iron-icon.js';
import '@vaadin/vaadin-upload/vaadin-upload.js'
import '@vaadin/vaadin-lumo-styles/color.js';
import './photo-icons.js';

// Gesture events like tap and track generated from touch will not be
// preventable, allowing for better scrolling performance.
setPassiveTouchGestures(true);

// Set Polymer's root path to the same value we passed to our service worker
// in `index.html`.
setRootPath(PhotoAppGlobals.rootPath);

class PhotoApp extends PolymerElement {
  static get template() {
    return html`
      <style include="lumo-color">
        :host {
          --app-primary-color: #2d3035;
          --app-primary-text-color: #db6574;
          --app-secondary-color: #8a8d93;
          --app-highlight-color: #e76073;
          --app-card-background-color: #2d3035;  
          display: block;
        }

        app-header {
          color: #fff;
          background-color: var(--app-primary-color);
        }

        app-header paper-icon-button {
          --paper-icon-button-ink-color: white;
        }

        app-drawer {
          --app-drawer-content-container: {
            background-color: #2d3035;
          }
        }

        .drawer-list {
          margin: 0 20px;
        }

        .drawer-list a {
          display: block;
          padding: 0 16px;
          text-decoration: none;
          color: var(--app-secondary-color);
          line-height: 40px;
        }

        .drawer-list a:focus {
          outline: none;
        }

        .drawer-list a.iron-selected {
          color: #bfc1c4;
          font-weight: bold;
          background-color: #34373d
        }

        .text-primary {
          color: var(--app-primary-text-color);
        }

        paper-dialog {
          background-color: var(--app-primary-color);
        }
      </style>

      <app-location route="{{route}}" url-space-regex="^[[rootPath]]">
      </app-location>

      <app-route route="{{route}}" pattern="[[rootPath]]:page" data="{{routeData}}" tail="{{subroute}}">
      </app-route>

      <app-drawer-layout fullbleed="" force-narrow="">
        <!-- Drawer content -->
        <app-drawer id="drawer" slot="drawer">
          <iron-selector selected="[[page]]" attr-for-selected="name" class="drawer-list" role="navigation">
            <a name="home" href="[[rootPath]]">Home</a>
            <a name="upload" href="[[rootPath]]upload">Upload</a>
            <a name="view1" href="[[rootPath]]view1">View One</a>
            <a name="view2" href="[[rootPath]]view2">View Two</a>
            <a name="view3" href="[[rootPath]]view3">View Three</a>
          </iron-selector>
        </app-drawer>

        <!-- Main content -->
        <app-header-layout>
          <app-header slot="header">
            <app-toolbar>
              <paper-icon-button icon="photo-icons:menu" drawer-toggle=""></paper-icon-button>
              <div main-title=""><strong class="text-primary">Photo</strong><strong>Manager</strong></div>
              <paper-icon-button icon="photo-icons:file-upload" on-click="_upload"></paper-icon-button>
            </app-toolbar>
          </app-header>

          <iron-pages selected="[[page]]" attr-for-selected="name" role="main">
            <home-view name="home" path="[[path]]"></home-view>
            <upload-view name="upload"></upload-view>
            <my-view1 name="view1"></my-view1>
            <my-view2 name="view2"></my-view2>
            <my-view3 name="view3"></my-view3>
          </iron-pages>
        </app-header-layout>
      </app-drawer-layout>
      <paper-dialog id="upload" modal="">
        <paper-dialog-scrollable>
          <div theme="dark">
            <vaadin-upload accept="image/*" target="/api/pictures" method="POST"></vaading-upload>
          </div>
        </paper-dialog-scrollable>
        <div class="buttons">
          <paper-button dialog-confirm autofocus>Tap me to close</paper-button>
        </div>
      </paper-dialog>
    `;
  }

  static get properties() {
    return {
      path: String,
      page: {
        type: String,
        reflectToAttribute: true,
        observer: '_pageChanged'
      },
      routeData: Object,
      subroute: Object
    };
  }

  static get observers() {
    return [
      '_routePageChanged(routeData.page)'
    ];
  }

  _routePageChanged(page) {
    if (['upload', 'view1', 'view2', 'view3'].indexOf(page) !== -1) {
      this.page = page;
    } else {
      this.path = window.location.pathname.replace(/\/$/, "");
      this.page = 'home';
    }

    // Close a non-persistent drawer when the page & route are changed.
    if (!this.$.drawer.persistent) {
      this.$.drawer.close();
    }
  }

  _pageChanged(page) {
    // Import the page component on demand.
    //
    // Note: `polymer build` doesn't like string concatenation in the import
    // statement, so break it up.
    switch (page) {
      case 'home':
        import('./home-view.js');
        break;
      case 'upload':
        import('./upload-view.js');
        break;
      case 'view1':
        import('./my-view1.js');
        break;
      case 'view2':
        import('./my-view2.js');
        break;
      case 'view3':
        import('./my-view3.js');
        break;
    }
  }

  _upload() {
    this.$.
    this.$.upload.open();
  }
}

window.customElements.define('photo-app', PhotoApp);
