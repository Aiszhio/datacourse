<template>
  <div class="client-orders-container">
    <h2>История заказов</h2>

    <!-- Таблица с историей заказов -->
    <div class="orders-section">
      <div v-if="loading" class="loading">Загрузка истории заказов...</div>

      <table v-if="!loading && orders.length > 0" class="order-table">
        <thead>
        <tr>
          <th>Название услуги</th>
          <th>Дата оформления</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="order in orders" :key="order.OrderID">
          <td>{{ order.ServiceName }}</td>
          <td>{{ formatDate(order.OrderDate) }}</td>
        </tr>
        </tbody>
      </table>

      <div v-if="!loading && orders.length === 0" class="no-orders">
        У вас нет завершенных заказов.
      </div>
    </div>

    <!-- Форма для создания нового заказа -->
    <div class="new-order-section">
      <h3>Сделать новый заказ</h3>
      <form @submit.prevent="createOrder">
        <label for="service">Выберите услугу:</label>
        <select v-model="newOrder.service" id="service" required>
          <option v-for="service in services" :key="service.id" :value="service.name">
            {{ service.name }}
          </option>
        </select>

        <label for="orderDate">Дата и время оформления:</label>
        <input type="datetime-local" v-model="newOrder.orderDate" id="orderDate" required />

        <button type="submit" class="btn">Оформить заказ</button>
      </form>
    </div>

    <!-- Навигационные кнопки -->
    <div class="card-panel">
      <div @click="goToHome" class="card">
        <h4>На главную</h4>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ClientOrders',
  data() {
    return {
      orders: [], // Список завершенных заказов клиента
      loading: false, // Флаг для индикации загрузки
      newOrder: { // Данные для нового заказа
        service: '',
        orderDate: '',
      },
      services: [] // Доступные услуги
    };
  },
  mounted() {
    this.fetchOrders(); // Загружаем историю заказов при загрузке страницы
    this.fetchServices(); // Загружаем доступные услуги для формы
  },
  methods: {
    async fetchOrders() {
      this.loading = true;
      try {
        const response = await fetch('http://localhost:8080/api/client/orders/history', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer your-auth-token'
          }
        });

        if (!response.ok) throw new Error('Ошибка при получении истории заказов');
        this.orders = await response.json();
      } catch (error) {
        console.error('Ошибка:', error.message);
        alert('Не удалось загрузить историю заказов.');
      } finally {
        this.loading = false;
      }
    },
    async fetchServices() {
      try {
        const response = await fetch('http://localhost:8080/api/services', {
          method: 'GET',
          headers: { 'Content-Type': 'application/json' }
        });

        if (!response.ok) throw new Error('Ошибка при получении списка услуг');
        this.services = await response.json();
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
            'Authorization': 'Bearer your-auth-token'
          },
          body: JSON.stringify(this.newOrder)
        });

        if (!response.ok) throw new Error('Ошибка при создании заказа');
        alert('Заказ успешно создан!');
        this.newOrder = { service: '', orderDate: '' };
        this.fetchOrders();
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
      this.$router.push({ name: 'ClientHome' });
    },
  }
};
</script>

<style scoped>
.client-orders-container {
  padding: 20px;
  margin: 0 auto;
  max-width: 900px;
}

.orders-section, .new-order-section {
  background-color: #f9f9f9;
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

.order-table {
  width: 100%;
  border-collapse: collapse;
}

.order-table th, .order-table td {
  padding: 12px;
  border: 1px solid #ddd;
  text-align: left;
}

.order-table th {
  background-color: #f2f2f2;
  font-weight: bold;
}

.loading, .no-orders {
  text-align: center;
  font-size: 1.2em;
  color: #555;
  margin-top: 20px;
}

form label {
  font-weight: bold;
  margin-top: 10px;
}

form input, form select {
  padding: 10px;
  margin-top: 5px;
  margin-bottom: 15px;
  border: 1px solid #ccc;
  border-radius: 5px;
  width: 100%;
  box-sizing: border-box;
}

.btn {
  background-color: #4CAF50;
  color: white;
  padding: 10px 20px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.btn:hover {
  background-color: #45a049;
}

.card-panel {
  display: flex;
  gap: 20px;
  justify-content: center;
}

.card {
  background-color: #4CAF50;
  color: white;
  border-radius: 10px;
  padding: 5px;
  width: 180px;
  text-align: center;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  cursor: pointer;
  transition: transform 0.3s, box-shadow 0.3s;
}

.card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
}
</style>
