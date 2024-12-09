import { createApp } from 'vue';
import App from './App.vue';
import VueTheMask from 'vue-the-mask';
import router from './router';

createApp(App)
    .use(router)
    .use(VueTheMask)
    .mount('#app');
