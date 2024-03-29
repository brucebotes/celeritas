const esbuild = require('esbuild')
const liveServer = require('live-server') // dev server
const sveltePlugin = require('esbuild-svelte') // esbuild plugin svelte
const sveltePreProcess = require('svelte-preprocess')

function showUsage() {
  console.log('USAGE')
  console.log('node esbuild.js dev')
  console.log('node esbuild.js prod')
  process.exit(0)
}

if (process.argv.length < 3) {
  showUsage()
}

if (!['dev', 'prod'].includes(process.argv[2])) {
  showUsage()
}

// production mode, or not
const production = process.argv[2] === 'prod'

// esbuild watch in dev mode to rebuild out files
let watch = false
if (!production) {
  watch = {
    onRebuild(error) {
      if (error)
        console.error('esbuild: Watch build failed:', error.getMessage())
      else console.log('esbuild: Watch build succeeded')
    }
  }
}

const buildpath = production
  ? '../../public/views/${MOD_NAME}/bundle.js'
  : './public/build/bundle.js'

// esbuild build options
// see: https://esbuild.github.io/api/#build-api
const options = {
  entryPoints: ['./src/main.js'],
  bundle: true,
  watch,
  format: 'iife',
  //minify: production,
  minify: false,
  sourcemap: true,
  outfile: buildpath,
  plugins: [
    sveltePlugin({
      preprocess: sveltePreProcess() // All preprocess processes
    })
  ]
}

// start web dev server
if (!production) {
  const params = {
    port: 8081, // Set the server port. Defaults to 8080.
    root: './public', // Set root directory that's being served. Defaults to cwd.
    open: true, // When false, it won't load your browser by default.
    wait: 1000, // Waits for all changes, before reloading. Defaults to 0 sec.
    logLevel: 2 // 0 = errors only, 1 = some, 2 = lots
  }
  liveServer.start(params)
}

// esbuild dev + prod
esbuild.build(options).catch((err) => {
  console.error(err)
  process.exit(0)
})
