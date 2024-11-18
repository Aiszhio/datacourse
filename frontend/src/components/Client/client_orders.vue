<template>
  <div class="client-orders-container">
    <h2>История заказов</h2>

    <!-- Таблица с историей заказов -->
    <div class="orders-section">
      <div v-if="loading" class="loading">Загрузка истории заказов...</div>

      <table v-else-if="pastOrders.length" class="order-table">
        <thead>
        <tr>
          <th>№</th>
          <th>Имя фотографа</th>
          <th>Название услуги</th>
          <th>Дата оформления</th>
          <th>Дата получения</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(order, index) in pastOrders" :key="order.order_id">
          <td>{{ index + 1 }}</td>
          <td>{{ order.employee_name || 'Не назначен' }}</td>
          <td>{{ order.service_name }}</td>
          <td>{{ formatDate(order.order_date) }}</td>
          <td>{{ formatDate(order.receipt_date) }}</td>
        </tr>
        </tbody>
      </table>

      <div v-else class="no-orders">
        У вас нет завершенных заказов.
      </div>
    </div>

    <!-- Форма для создания нового заказа -->
    <div class="new-order-section">
      <h3>Сделать новый заказ</h3>
      <form @submit.prevent="createOrder">
        <label for="service">Выберите услугу:</label>
        <select v-model="newOrder.service" id="service" required>
          <option v-for="service in services" :key="service.service_id" :value="service.service_id">
            {{ service.service_name }} — {{ service.price }}₽
          </option>
        </select>

        <label for="orderDate">Дата и время оформления:</label>
        <input
            type="datetime-local"
            v-model="newOrder.orderDate"
            id="orderDate"
            :min="currentDateTime"
            required
        />

        <div v-if="busyDates.length" class="busy-dates-warning">
          <p>Выбранная дата занята! Попробуйте другую:</p>
          <ul>
            <li v-for="date in busyDates" :key="date">{{ formatDate(date) }}</li>
          </ul>
        </div>

        <button type="submit" class="btn" :disabled="busyDates.length">
          Оформить заказ
        </button>
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
      orders: [],
      loading: false,
      newOrder: {
        service: '',
        orderDate: '',
      },
      services: [],
      busyDates: [],
      currentDateTime: '',
    };
  },
  computed: {
    pastOrders() {
      const now = new Date();
      return this.orders.filter(order => new Date(order.receipt_date) < now);
    },
  },
  mounted() {
    this.updateCurrentDateTime();
    this.fetchData('orders');
    this.fetchData('services');
    this.interval = setInterval(this.updateCurrentDateTime, 60000);
  },
  beforeUnmount() {
    clearInterval(this.interval);
  },
  methods: {
    updateCurrentDateTime() {
      const now = new Date();
      this.currentDateTime = now.toISOString().slice(0, 16);
    },
    async fetchData(type) {
      this.loading = type === 'orders';
      try {
        const endpoint = type === 'orders' ? 'orders' : 'services';
        const response = await fetch(`http://localhost:8080/api/${endpoint}`, {
          method: 'GET',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
        });

        if (!response.ok) throw new Error(`Ошибка при получении ${type === 'orders' ? 'заказов' : 'услуг'}`);

        const data = await response.json();
        this[type] = data;
      } catch (error) {
        console.error("Ошибка:", error.message);
        alert(`Не удалось загрузить ${type === 'orders' ? 'историю заказов' : 'список услуг'}.`);
      } finally {
        if (type === 'orders') this.loading = false;
      }
    },
    async createOrder() {
      const { service, orderDate } = this.newOrder;
      if (!service || !orderDate) {
        alert("Пожалуйста, заполните все поля.");
        return;
      }

      try {
        const payload = {
          service,
          orderDate: new Date(orderDate).toISOString(),
        };

        const response = await fetch("http://localhost:8080/api/createOrder", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify(payload),
        });

        const result = await response.json();

        if (!response.ok) {
          throw new Error(result.error || "Не удалось создать заказ.");
        }

        alert("Заказ успешно создан!");
        this.resetForm();
        this.fetchData('orders');
      } catch (error) {
        console.error("Ошибка при создании заказа:", error.message);
        alert(error.message || "Не удалось создать заказ.");
      }
    },
    resetForm() {
      this.newOrder = { service: '', orderDate: '' };
      this.busyDates = [];
    },
    formatDate(dateString) {
      if (!dateString) return "Дата не указана";
      return new Date(dateString).toLocaleString("ru-RU", {
        year: "numeric",
        month: "long",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      });
    },
    goToHome() {
      this.$router.push({ name: "ClientHome" });
    },
  },
};
</script>

<style scoped>
.client-orders-container {
  padding: 20px;
  margin: 0 auto;
  max-width: 900px;
}

.orders-section,
.new-order-section {
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

.order-table th,
.order-table td {
  padding: 12px;
  border: 1px solid #ddd;
  text-align: left;
}

.order-table th {
  background-color: #f2f2f2;
  font-weight: bold;
}

.loading,
.no-orders {
  text-align: center;
  font-size: 1.2em;
  color: #555;
  margin-top: 20px;
}

form label {
  font-weight: bold;
  margin-top: 10px;
}

form input,
form select {
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
  padding: 15px;
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

.busy-dates-warning {
  color: red;
  margin-top: 10px;
}

.busy-dates-warning ul {
  list-style-type: disc;
  margin-left: 20px;
}
</style>
