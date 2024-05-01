import { createApp } from 'vue';
import { createPinia } from 'pinia';

import App from './App.vue';
import './index.css';
import router from '@/router';
import interceptor from '@/plugins/interceptor';

const pinia = createPinia();

const app = createApp(App).use(router).use(interceptor).use(pinia);

app.mount('#app');
