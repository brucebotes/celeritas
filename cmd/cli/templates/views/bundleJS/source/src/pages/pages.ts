import { html } from 'lit'

export const rootHome = () => html`
  <div>
    <h1>Home page for root</h1>
    <h1 style="color: red;">We should not be here</h1>
  </div>
`
export const subrootHome = () => html`
  <div>
    <h1>Home</h1>
    <simple-greeting></simple-greeting>
  </div>
`
export const projects = () => html`<h1>Projects</h1>`
export const about = () => html`<h1>About</h3>`
