<script>
export default {
  name: 'ClientOrders',
  data() {
    return {
      orders: [],  // Список заказов клиента
      loading: false,  // Флаг для индикации загрузки
      newOrder: {  // Данные для нового заказа
        service: '',
        orderDate: '',
      },
      services: [] // Доступные услуги
    };
  },
  mounted() {
    this.fetchOrders();  // Загружаем заказы клиента при загрузке страницы
    this.fetchServices(); // Загружаем доступные услуги для формы
  },
  methods: {
    async fetchOrders() {
      this.loading = true;  // Устанавливаем флаг загрузки
      try {
        const response = await fetch('http://localhost:8080/api/client/orders', {  // Запрос заказов клиента
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer your-auth-token'  // Токен для авторизации
          }
        });

        if (!response.ok) {
          throw new Error('Ошибка при получении списка заказов');
        }

        const data = await response.json();
        this.orders = data;  // Записываем полученные заказы
      } catch (error) {
        console.error('Ошибка:', error.message);
        alert('Не удалось загрузить список заказов. Попробуйте позже.');
      } finally {
        this.loading = false;  // Отключаем флаг загрузки после завершения запроса
      }
    },
    async fetchServices() {
      try {
        const response = await fetch('http://localhost:8080/api/services', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json'
          }
        });

        if (!response.ok) {
          throw new Error('Ошибка при получении списка услуг');
        }

        const data = await response.json();
        this.services = data;  // Записываем полученные услуги
      } catch (error) {
        console.error('Ошибка:', error.message);
        alert('Не удалось загрузить список услуг.');
      }
    },
    async createOrder() {
      if (!this.newOrder.service || !this.newOrder.orderDate) {
        alert('Пожалуйста, заполните все поля.');
        return;
      }

      try {
        const response = await fetch('http://localhost:8080/api/client/orders', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer your-auth-token'  // Токен авторизации, если требуется
          },
          body: JSON.stringify(this.newOrder)
        });

        if (!response.ok) {
          throw new Error('Ошибка при создании заказа');
        }

        alert('Заказ успешно создан!');
        this.newOrder = { service: '', orderDate: ''};  // Сброс формы
        this.fetchOrders();  // Обновляем список заказов
      } catch (error) {
        console.error('Ошибка:', error.message);
        alert('Не удалось создать заказ.');
      }
    },
    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleDateString('ru-RU', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    },
    goToHome() {
      this.$router.push({ name: 'ClientHome' });  // Переход на главную страницу клиента
    },
    goToPayInvoices() {
      this.$router.push({ name: 'Payment' });  // Переход на страницу оплаты счетов
    }
  }
};
</script>

<template>
  <div class="client-orders-container">
    <h2>Ваши заказы</h2>

    <!-- Обертка для таблицы заказов -->
    <div class="orders-section section">
      <!-- Загрузка данных -->
      <div v-if="loading" class="loading">Загрузка ваших заказов...</div>

      <!-- Таблица заказов -->
      <table class="orders-table" v-if="!loading">
        <thead>
        <tr>
          <th>ID Заказа</th>
          <th>Услуга</th>
          <th>Дата заказа</th>
          <th>Текущий статус</th>
        </tr>
        </thead>
        <tbody>
        <!-- Показать заказы -->
        <tr v-if="orders.length > 0" v-for="order in orders" :key="order.OrderID">
          <td>{{ order.OrderID }}</td>
          <td>{{ order.ServiceName }}</td>
          <td>{{ formatDate(order.OrderDate) }}</td>
          <td>{{ order.Status }}</td>
        </tr>
        <!-- Показать пустую строку, если заказов нет -->
        <tr v-if="orders.length === 0">
          <td colspan="5" class="no-orders">У вас нет текущих заказов.</td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Форма для создания нового заказа -->
    <div class="new-order-section section">
      <h3>Сделать новый заказ</h3>
      <form @submit.prevent="createOrder">
        <label for="service">Выберите услугу:</label>
        <select v-model="newOrder.service" id="service" required>
          <option v-for="service in services" :key="service.id" :value="service.name">
            {{ service.name }}
          </option>
        </select>

        <label for="orderDate">Дата и время заказа:</label>
        <input type="datetime-local" v-model="newOrder.orderDate" id="orderDate" required />

        <button type="submit" class="btn">Заказать</button>
      </form>
    </div>

    <!-- Кнопки для навигации -->
    <div class="navigation-buttons">
      <button @click="goToHome" class="btn">На главную</button>
      <button @click="goToPayInvoices" class="btn">Оплатить счета</button>
    </div>
  </div>
</template>

<style scoped>
.client-orders-container {
  padding: 20px;
  margin: 0 auto;
  max-width: 900px;
  max-height: 100vh;
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
  background-color: #4CAF50;
  color: white;
  padding: 10px 20px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 1em;
  transition: background-color 0.3s ease;
}

.btn:hover {
  background-color: #45a049;
}

.navigation-buttons {
  display: flex;
  justify-content: space-between;
  margin-top: 20px;
}

form {
  display: flex;
  flex-direction: column;
}

input, select {
  padding: 10px;
  margin-bottom: 15px;
  border: 1px solid #ccc;
  border-radius: 5px;
}

</style>
