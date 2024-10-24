import { createRouter, createWebHistory } from 'vue-router';
import home from './components/home.vue';
import User from './components/client_home.vue';
import Payment from './components/payment.vue';
import ClientOrders from "./components/client_orders.vue";
import WorkerOrders from "./components/worker_orders.vue";
import Worker from './components/worker_dashboard.vue';
import Admin from './components/admin_home.vue'

const routes = [
    { path: '/', name: 'home', component: home }, // Главная страница
    { path: '/client', name: 'ClientHome', component: User }, // Страница клиента
    { path: '/client/orders', name: 'ClientOrders', component: ClientOrders }, // Заказы клиента
    { path: '/client/payment', name: 'Payment', component: Payment }, // Оплата счетов клиента
    { path: '/worker', name: 'WorkerDashboard', component: Worker }, // Панель сотрудника
    { path: '/worker/orders', name: 'WorkerOrders', component: WorkerOrders }, // Заказы сотрудника
    { path: '/admin', name: 'AdminHome', component: Admin }
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;
