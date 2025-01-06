import { createRouter, createWebHistory } from 'vue-router';
import axios from 'axios';
import home from './components/home.vue';
import User from './components/Client/client_home.vue';
import ClientOrders from "@/components/Client/client_orders.vue";
import Worker from './components/Worker/worker_dashboard.vue';
import AdminDashboard from './components/Admin/admin_home.vue';
import ManageEmployees from './components/Admin/manage_employees.vue';
import OrderHistory from "./components/Admin/order_history.vue";
import MaterialsOverview from './components/Admin/manage_materials.vue';
import Services from "./components/Admin/services.vue";
import Bookings from "./components/Admin/booking.vue";
import Clients from "./components/Admin/manage_clients.vue"

const routes = [
    { path: '/', name: 'home', component: home }, // Главная страница
    { path: '/client', name: 'ClientHome', component: User }, // Страница клиента
    { path: '/client/orders', name: 'ClientOrders', component: ClientOrders }, // Заказы клиента
    { path: '/worker', name: 'WorkerDashboard', component: Worker }, // Панель сотрудника
    { path: '/admin', name: 'AdminHome', component: AdminDashboard },
    { path: '/admin/manage_employee', name: 'ManageEmp', component: ManageEmployees },
    { path: '/admin/orders', name: 'OrderHistory', component: OrderHistory },
    { path: '/admin/MaterialsOverview', name: 'MaterialsOverview', component: MaterialsOverview },
    { path: '/admin/Services', name: 'Services', component: Services },
    { path: '/admin/Bookings', name: 'Bookings', component: Bookings },
    { path: '/admin/Clients', name: 'Clients', component: Clients}
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

router.beforeEach(async (to, from, next) => {
    // Пропускаем проверку роли для главной страницы
    if (to.path === '/') {
        return next();
    }

    try {
        // Запрашиваем роль пользователя из Redis через API
        const response = await axios.get('http://localhost:8080/api/CheckRole', { withCredentials: true });
        const userRole = response.data.role;

        // Логика маршрутов в зависимости от роли
        if (to.path.startsWith('/admin') && userRole !== 'admin') {
            return next({ name: 'home' });
        }

        if (to.path.startsWith('/worker') && userRole !== 'worker') {
            return next({ name: 'home' });
        }

        if (to.path.startsWith('/client') && userRole !== 'client') {
            return next({ name: 'home' });
        }

        // Разрешаем доступ к маршруту, если роль совпадает
        return next();
    } catch (error) {
        // Проверяем, есть ли ответ от сервера
        if (error.response) {
            const { status, data } = error.response;

            if (status === 403 && data.error === 'Доступ запрещен, вы уволены') {
                console.error('Пользователь уволен:', data.error);
                // Перенаправляем на страницу с сообщением о запрете доступа
                return next({ name: 'accessForbidden' });
            }

            if (data && data.error) {
                console.error('Ошибка при проверке роли:', data.error);
                // Можно добавить отображение уведомления пользователю, например:
                // Notify.show(data.error);
            } else {
                console.error('Неизвестная ошибка при проверке роли:', error);
            }
        } else {
            console.error('Ошибка при отправке запроса:', error);
        }

        // Перенаправляем на главную страницу при ошибке
        return next({ name: 'home' });
    }
});


export default router;
