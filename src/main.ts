import { mount } from 'svelte'
import './app.css'
import App from './App.svelte'
import { initializePlugins } from './lib/plugins/init.js'

// Initialize plugins before mounting the app
await initializePlugins();

const app = mount(App, {
  target: document.getElementById('app')!,
})

export default app
