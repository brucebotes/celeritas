<script lang="ts">
  import { Router, Route, NotFound, redirect } from './pager'

  export let basePath = '/svh/${MOD_NAME}'

  import Login from './pages/Login.svelte'
  import Home from './pages/Home.svelte'
  import About from './pages/About.svelte'
  import Profile from './pages/Profile.svelte'

  const data = { foo: 'bar', custom: true }

  const guard = (ctx, next) => {
    // check for example if user is authenticated
    if (true) {
      redirect('/login')
    } else {
      // go to the next callback in the chain
      next()
    }
  }
</script>

<main>
  <nav class="navbar is-white" role="navigation" aria-label="Main navigation">
    <div class="container">
      <div class="navbar-brand">
        <a href="./" class="navbar-item">Brand</a>
        <button class="button navbar-burger" data-target="navMenu">
          <span />
          <span />
          <span />
        </button>
      </div>
      <div class="navbar-menu" id="navMenu">
        <div class="navbar-start">
          <a href="./" class="navbar-item">home</a>
          <a href="./about" class="navbar-item">about</a>
          <a href="./profile/joe" class="navbar-item">profile</a>
          <a href="./news" class="navbar-item">news</a>
          <a href="./login" class="navbar-item">login</a>
        </div>
        <div class="navbar-end">
          <div class="navbar-item has-dropdown is-hoverable">
            <a href="#0" class="navbar-link">Afrikaans</a>
            <div class="navbar-dropdown is-boxed">
              <a href="https://companyname.com/language" class="navbar-item"
                >Language</a
              >
              <a href="https://companyname.com/language" class="navbar-item"
                >Language</a
              >
              <a href="https://companyname.com/language" class="navbar-item"
                >Language</a
              >
            </div>
          </div>
        </div>
      </div>
    </div>
  </nav>

  <Router {basePath}>
    <Route path="/" component={Home} {data} />
    <Route path="/about" component={About} />
    <Route path="/login" component={Login} />
    <Route path="/profile/:username" let:params>
      <h2>Hello {params.username}!</h2>
      <p>Here is your profile</p>
    </Route>
    <Route path="/news" middleware={[guard]}>
      <h2>Latest News</h2>
      <p>Finally some good news!</p>
    </Route>
    <NotFound>
      <h2>Sorry. Page not found.</h2>
    </NotFound>
  </Router>
</main>

