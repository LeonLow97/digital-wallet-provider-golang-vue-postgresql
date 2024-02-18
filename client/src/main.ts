import { createApp } from 'vue';
import { createPinia } from 'pinia';

import App from './App.vue';
import './style.css';
import router from '@/router';
import axios from '@/plugins/axios';

const pinia = createPinia();

const app = createApp(App).use(router).use(axios).use(pinia);

app.mount('#app');
