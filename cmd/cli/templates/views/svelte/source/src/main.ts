import App from './App.svelte'

const app = new App({
  target: document.body,
  props: {
    basePath: '/svh/${MOD_NAME}'
  }
})

export default app
