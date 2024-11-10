import { createRouter, createWebHistory } from 'vue-router';
import home from './components/home.vue';
import User from './components/Client/client_home.vue';
import ClientOrders from "@/components/Client/client_orders.vue";
import Worker from './components/Worker/worker_dashboard.vue';
import AdminDashboard from './components/Admin/admin_home.vue';
import ManageEmployees from './components/Admin/manage_employees.vue';
import Order_history from "./components/Admin/order_history.vue";
import MaterialsOverview from './components/Admin/manage_materials.vue';
import Services from "./components/Admin/services.vue";
import Bookings from "./components/Admin/booking.vue";

const routes = [
    { path: '/', name: 'home', component: home }, // Главная страница
    { path: '/client', name: 'ClientHome', component: User }, // Страница клиента
    { path: '/client/orders', name: 'ClientOrders', component: ClientOrders }, // Заказы клиента
    { path: '/worker', name: 'WorkerDashboard', component: Worker }, // Панель сотрудника
    { path: '/admin', name: 'AdminHome', component: AdminDashboard },
    { path: '/admin/manage_employee', name: 'ManageEmp', component: ManageEmployees },
    { path: '/admin/orders', name: 'OrderHistory', component: Order_history },
    { path: '/admin/MaterialsOverview', name: 'MaterialsOverview', component: MaterialsOverview },
    { path: '/admin/Services', name: 'Services', component: Services },
    { path: '/admin/Bookings', name: 'Bookings', component: Bookings }
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;
