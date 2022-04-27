import type { Channel } from 'pusher-js'

export type InitialData = {
  id: string | undefined
  token: string
  publicChannel: Channel
  privateChannel: Channel
  attention: {
    confirm: ({}: any) => {}
    alert: ({}: any) => {}
    promptConfirm: ({}: any) => {}
    prompt: ({}: any) => {}
    toast: ({}: any) => {}
    loading: () => {}
  }
}

// declare variable within the global scope
declare global {
  var __INITIAL_DATA__: InitialData
}
export const initialData = __INITIAL_DATA__
