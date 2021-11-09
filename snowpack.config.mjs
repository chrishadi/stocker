/** @type {import("snowpack").SnowpackUserConfig } */
export default {
  root: 'assets',
  mount: {
    /* ... */
  },
  plugins: [
    /* ... */
  ],
  routes: [
    /* Enable an SPA Fallback in development: */
    // {"match": "routes", "src": ".*", "dest": "/index.html"},
  ],
  optimize: {
    /* Example: Bundle your final build: */
    'bundle': true,
    'entrypoints': ['assets/chart.jsx'],
    'minify': true,
    'target': 'es2018'
  },
  packageOptions: {
    /* ... */
  },
  devOptions: {
    /* ... */
  },
  buildOptions: {
    out: 'static'
  },
};
