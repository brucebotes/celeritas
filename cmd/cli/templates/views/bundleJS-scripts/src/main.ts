import { LitElement, css, html } from 'lit'
import { customElement, property, state, query } from 'lit/decorators.js'
import 'urlpattern-polyfill'
import { Router, Routes } from '@lit-labs/router'
import './style.css'

import { Drawer } from '@material/mwc-drawer'
import './applications.ts'
import '@material/mwc-top-app-bar'
import '@material/mwc-icon-button'
import '@material/mwc-drawer'
import '@material/mwc-list'

import { rootHome, subrootHome, projects, about } from './pages/pages'

const modPath = '/jsb/${MOD_NAME}'

@customElement('module-router')
class ModuleRouter extends LitElement {
  static styles = css`
    p,
    .time {
      margin: 5px;
    }

    footer {
      text-align: center;
    }

    mwc-drawer[open] mwc-top-app-bar {
      /* Default width of drawer is 256px. See CSS Custom Properties below */
      --mdc-top-app-bar-width: calc(100% - var(--mdc-drawer-width, 256px));
    }
  `
  private _router = new Router(this, [
    {
      path: modPath,
      render: rootHome,
      enter: async () => {
        // TODO: this ensures that we have "/" added to the base url (ie modPath)
        // But it does destroy the histrory. I need to fix this design or find
        // an alternative approach
        //history.pushState({}, '', modPath + '/')
        await this._router.goto(modPath + '/')
        return false
      }
    },
    {
      path: modPath + '/*',
      render: () => html`<module-routes></module-routes>`
    }
  ])
  @state()
  private _isOpen: boolean = false

  @query('mwc-drawer')
  drawer!: Drawer

  _toggleDrawer() {
    this._isOpen = !this._isOpen
    console.log('isOpen \u2192', this._isOpen)
    if (this.drawer) {
      console.log('Setting drawer from ', this.drawer.open)
      this.drawer.open = !this.drawer.open
    }
  }
  _unSelectMWC_list_item() {
    const aList = this.renderRoot.querySelectorAll('mwc-list-item')
    aList.forEach((i) => {
      i.removeAttribute('selected')
      i.removeAttribute('activated')
    })
  }

  async _toggleActive(e: Event, href?: string) {
    this._unSelectMWC_list_item()
    const mwcListElem = e.currentTarget! as HTMLElement
    mwcListElem.setAttribute('selected', '')
    mwcListElem.setAttribute('activated', '')
    if (href !== undefined) {
      const route = this._router.link() + href
      history.pushState({}, '', route)
      await this._router.goto(route)
    }
  }

  render() {
    return html`
      <mwc-drawer type="dismissible" ?open=${this._isOpen}>
        <div>
          <mwc-list>
            <mwc-list-item
              selected
              activated
              @click=${(e: Event) => this._toggleActive(e, '')}
            >
              <p>Home</p>
            </mwc-list-item>
            <mwc-list-item
              @click=${(e: Event) => this._toggleActive(e, 'projects')}
            >
              <p>Projects</p>
            </mwc-list-item>
            <mwc-list-item
              @click=${(e: Event) => this._toggleActive(e, 'todos')}
            >
              <p>Todos</p>
            </mwc-list-item>
            <mwc-list-item
              @click=${(e: Event) => this._toggleActive(e, 'timers')}
            >
              <p>Timers</p>
            </mwc-list-item>
            <mwc-list-item
              @click=${(e: Event) => this._toggleActive(e, 'about')}
            >
              <p>About</p>
            </mwc-list-item>
          </mwc-list>
        </div>
        <div slot="appContent">
          <mwc-top-app-bar>
            <mwc-icon-button
              icon="menu"
              slot="navigationIcon"
              @click=${this._toggleDrawer}
            ></mwc-icon-button>
            <div slot="title">Title</div>
            <mwc-icon-button
              icon="file_download"
              slot="actionItems"
            ></mwc-icon-button>
            <mwc-icon-button icon="print" slot="actionItems"></mwc-icon-button>
            <mwc-icon-button
              icon="favorite"
              slot="actionItems"
            ></mwc-icon-button>
            <div>
              <main>${this._router.outlet()}</main>
              <hr />
              <footer>
                <my-clock-element class="time"></my-clock-element>
                <p>Routing is contolled by @lit-labs/router</p>
              </footer>
            </div>
          </mwc-top-app-bar>
        </div>
      </mwc-drawer>
    `
  }
}

@customElement('module-routes')
class ModuleRoutes extends LitElement {
  static styles = css`
    .body {
      font-family: sans-serif;
      display: flex;
      justify-content: center;
    }

    motion-todos {
      width: 600px;
    }

    @media screen and (max-width: 600px) {
      motion-todos {
        width: 100%;
      }
    }
  `

  private _routes = new Routes(this, [
    {
      path: '',
      render: subrootHome
    },
    {
      pattern: new URLPattern({
        pathname: 'projects'
      }),
      render: projects
    },
    {
      pattern: new URLPattern({
        pathname: 'about'
      }),
      render: about
    },
    {
      pattern: new URLPattern({
        pathname: 'timers'
      }),
      render: () => html`<buncha-timers></buncha-timers>`
    },
    {
      pattern: new URLPattern({
        pathname: 'todos'
      }),
      render: () => html`<motion-todos></motion-todos>`
    }
  ])
  render() {
    return html` <div class="body">${this._routes.outlet()}</div> `
  }
}
