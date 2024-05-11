## Installing Tailwind with Vue 3 (Vite)

- [Install Tailwind CSS with Vue 3 and Vite](https://v2.tailwindcss.com/docs/guides/vue-3-vite)

```
## Install Tailwind and its peer-dependencies using npm:
npm install -D tailwindcss@latest postcss@latest autoprefixer@latest

## Next, generate your tailwind.config.js and postcss.config.js files:
npx tailwindcss init -p

## Add Prettier class sorting in tailwind css
## Link: https://tailwindcss.com/blog/automatic-class-sorting-with-prettier
npm install -D prettier prettier-plugin-tailwindcss

## Create a `.prettierrc` file and add the following
{
    "plugins": [
        "prettier-plugin-tailwindcss"
    ]
}
```

## Installing qrcode.vue for scanning qr code

```
## Package Documentation: https://www.npmjs.com/package/qrcode.vue
npm install --save qrcode.vue # yarn add qrcode.vue
```