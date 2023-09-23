# Setup

## Golang API

- `go mod init github.com/LeonLow97`

## VueJS

- `npm install vue@3`
- `cd` into project folder
- `npm install`
- `npm run dev`

---

#### Vitest

- Add to `package.json`: `"test:unit": "vitest --environment jsdom",`
- Test file must end with `.test.js` or `.spec.js`
- To run the test, `npm run test:unit`

---

#### Vitest Coverage

- Add to `package.json`: `"test:unit": "vitest --environment jsdom --coverage"`
- Run `npm run test:unit` and terminal will prompt to install `@vitest/coverage-v8`. Hit 'Y'.
- Run `npm install --save-dev eslint-plugin-vitest-globals`
- In `.eslintrc.cjs`, add:

```js
require('@rushstack/eslint-patch/modern-module-resolution');

module.exports = {
  root: true,
  extends: [
    'plugin:vue/vue3-recommended',
    'eslint:recommended',
    '@vue/eslint-config-prettier/skip-formatting',
    'plugin:vitest-globals/recommended', // add this
  ],
  parserOptions: {
    ecmaVersion: 'latest',
  },
  env: {
    'vitest-globals/env': true, // add this
  },
};
```

---

#### Vue Testing Library

- `npm install --save-dev @testing-library/vue @testing-library/jest-dom @testing-library/user-event`
- Create a `setup.js` in the root directory of your project to import some things from Vue Testing Library.
  - Add this `setup.js` file into a `tests` folder.
- Add the `setup.js` file into `vite.config.js` to make it global.

---

#### Vite Plugin Vuetify, Vue Router, and others

- `npm install vite-plugin-vuetify --save-dev`
- `npm install vue-router`
- `npm install @mdi/font axios core-js roboto-fontface vuex webfontloader`
- `npm install --save-dev sass`

---
