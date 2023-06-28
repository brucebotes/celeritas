import { LitElement, html } from 'lit'
import { customElement } from 'lit/decorators.js'
import './timer/my-timer'

@customElement('buncha-timers')
export class BunchaTimers extends LitElement {
  render() {
    return html`
      <div>
        <my-timer duration="7"></my-timer>
        <my-timer duration="60"></my-timer>
        <my-timer duration="300"></my-timer>
      </div>
    `
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'buncha-timers': BunchaTimers
  }
}
