import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import { BootstrapVue3 } from 'bootstrap-vue-3';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-vue-3/dist/bootstrap-vue-3.css';
import VueTheMask from 'vue-the-mask';
import Swal from 'sweetalert2';
import 'sweetalert2/dist/sweetalert2.min.css';

const app = createApp(App);
app.config.globalProperties.$swal = Swal;

app.use(router);
app.use(VueTheMask);
app.use(BootstrapVue3);

app.mount('#app');
