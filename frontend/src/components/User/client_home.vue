<script>
export default {
  name: 'User',
  data() {
    return {
      userName: 'Имя Пользователя', // Можно получить из бэкенда или сессии
      orders: [], // Текущие заказы пользователя
      loading: true, // Флаг загрузки заказов
      orderLimit: 3 // Ограничение на количество отображаемых заказов
    };
  },
  mounted() {
    this.fetchOrders(); // Загружаем заказы при загрузке страницы
  },
  computed: {
    limitedOrders() {
      return this.orders.slice(0, this.orderLimit); // Ограничиваем количество заказов
    }
  },
  methods: {
    async fetchOrders() {
      try {
        const response = await fetch('http://localhost:8080/api/client/orders', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer your-auth-token' // Токен авторизации, если требуется
          }
        });

        if (!response.ok) {
          throw new Error('Ошибка при получении заказов');
        }

        const data = await response.json();
        this.orders = data; // Записываем полученные заказы
      } catch (error) {
        console.error('Ошибка:', error.message);
        alert('Не удалось загрузить заказы.');
      } finally {
        this.loading = false; // Отключаем флаг загрузки
      }
    },
    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleDateString('ru-RU', { year: 'numeric', month: 'long', day: 'numeric' });
    },
    goToCreateOrder() {
      this.$router.push({ name: 'ClientOrders' }); // Переход на страницу создания заказа
    },
    goToPayInvoices() {
      this.$router.push({ name: 'Payment' }); // Переход на страницу оплаты счетов
    }
  }
};
</script>

<template>
  <div class="user-dashboard">
    <h2>Панель пользователя</h2>

    <h3>Добро пожаловать, {{ userName }}!</h3>

    <!-- Список текущих заказов -->
    <div class="orders-section section">
      <h3>Ваши заказы</h3>
      <div v-if="loading" class="loading">Загрузка ваших заказов...</div>

      <!-- Таблица заказов всегда отображается -->
      <table class="orders-table">
        <thead>
        <tr>
          <th>ID Заказа</th>
          <th>Услуга</th>
          <th>Дата заказа</th>
          <th>Дата получения</th>
          <th>Текущий статус</th>
        </tr>
        </thead>
        <tbody>
        <!-- Показать данные, если они загружены -->
        <tr v-if="!loading && limitedOrders.length > 0" v-for="order in limitedOrders" :key="order.OrderID">
          <td>{{ order.OrderID }}</td>
          <td>{{ order.ServiceName }}</td>
          <td>{{ formatDate(order.OrderDate) }}</td>
          <td>{{ formatDate(order.ReceiptDate) }}</td>
          <td>{{ order.Status }}</td>
        </tr>
        <!-- Показать пустую строку, если данных нет -->
        <tr v-if="!loading && limitedOrders.length === 0">
          <td colspan="5" class="no-orders">У вас нет текущих заказов.</td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Кнопки для навигации -->
    <div class="navigation-buttons">
      <button @click="goToCreateOrder" class="btn">Сделать заказ</button>
      <button @click="goToPayInvoices" class="btn">Оплатить счета</button>
    </div>
  </div>
</template>


<style scoped>
.user-dashboard {
  padding: 20px;
  margin: 0 auto;
  max-width: 900px;
}

.section {
  background-color: #f9f9f9;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  padding: 20px;
  margin-bottom: 20px;
}

.orders-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
}

.orders-table th,
.orders-table td {
  padding: 12px;
  border: 1px solid #ddd;
  text-align: left;
}

.orders-table th {
  background-color: #f2f2f2;
}

.loading {
  text-align: center;
  color: #555;
  font-size: 1.2em;
}

.no-orders {
  text-align: center;
  color: #999;
  font-size: 1.2em;
}

.btn {
  background-color: #4CAF50; /* Зеленый цвет фона */
  color: white; /* Белый текст */
  padding: 10px 20px; /* Внутренние отступы */
  border: none; /* Убираем рамку */
  border-radius: 5px; /* Скругляем углы */
  cursor: pointer; /* Указатель мыши при наведении */
  font-size: 1em;
  transition: background-color 0.3s ease; /* Плавная смена цвета при наведении */
}

.btn:hover {
  background-color: #4CAF50;
  color: white;
}

.navigation-buttons {
  display: flex;
  justify-content: space-between;
}

</style>
