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
          <th>Название услуг</th>
          <th>Дата оформления</th>
          <th>Дата получения</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(order, index) in pastOrders" :key="order.order_id">
          <td>{{ index + 1 }}</td>
          <td>{{ order.employee_name || 'Не назначен' }}</td>
          <td>
            <ul>
              <li v-for="content in order.order_contents" :key="content.id">
                {{ content.service.name }} — {{ content.service.price }}₽
              </li>
            </ul>
          </td>
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
      <form @submit.prevent="addToCart">
        <label for="service">Выберите услугу:</label>
        <select v-model="selectedService" id="service" required>
          <option disabled value="">-- Выберите услугу --</option>
          <option v-for="service in services" :key="service.service_id" :value="service.service_id">
            {{ service.service_name }} — {{ service.price }}₽
          </option>
        </select>

        <button type="button" class="btn" @click="addToCart" :disabled="!selectedService">
          Добавить в корзину
        </button>
      </form>

      <!-- Выбор даты бронирования -->
      <div class="date-selection">
        <label for="bookingDate">Дата и время бронирования:</label>
        <input
            type="datetime-local"
            v-model="bookingDate"
            id="bookingDate"
            :min="currentDateTime"
            required
        />
      </div>

      <!-- Корзина заказов -->
      <div v-if="cart.length" class="cart-section">
        <h4>Корзина</h4>
        <table class="cart-table">
          <thead>
          <tr>
            <th>№</th>
            <th>Название услуги</th>
            <th>Цена</th>
            <th>Действия</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="(item, index) in cart" :key="index">
            <td>{{ index + 1 }}</td>
            <td>{{ getServiceName(item.service_id) }}</td>
            <td>{{ getServicePrice(item.service_id) }}₽</td>
            <td>
              <button type="button" class="remove-btn" @click="removeFromCart(index)">Удалить</button>
            </td>
          </tr>
          </tbody>
        </table>
        <button type="button" class="btn submit-btn" @click="submitOrder" :disabled="!bookingDate">
          Оформить заказ
        </button>
      </div>
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
      selectedService: '',
      bookingDate: '',
      services: [],
      cart: [],
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
        if (type === 'orders') {
          // Предполагается, что сервер возвращает Order вместе с OrderContent
          this.orders = data.map(order => ({
            ...order,
            order_contents: order.order_contents || [],
          }));
        } else {
          this.services = data;
        }
      } catch (error) {
        console.error("Ошибка:", error.message);
        alert(`Не удалось загрузить ${type === 'orders' ? 'историю заказов' : 'список услуг'}.`);
      } finally {
        if (type === 'orders') this.loading = false;
      }
    },
    addToCart() {
      if (!this.selectedService) {
        alert("Пожалуйста, выберите услугу.");
        return;
      }

      // Проверка, добавлена ли уже эта услуга в корзину
      const exists = this.cart.find(item => item.service_id === this.selectedService);
      if (exists) {
        alert("Эта услуга уже добавлена в корзину.");
        return;
      }

      this.cart.push({
        service_id: this.selectedService,
      });

      // Очистка поля выбора
      this.selectedService = '';
    },
    removeFromCart(index) {
      this.cart.splice(index, 1);
    },
    async submitOrder() {
      if (!this.cart.length) {
        alert("Корзина пуста.");
        return;
      }

      if (!this.bookingDate) {
        alert("Пожалуйста, выберите дату бронирования.");
        return;
      }

      try {
        const payload = {
          booking_date: new Date(this.bookingDate).toISOString(),
          services: this.cart.map(item => ({
            service_id: item.service_id,
          })),
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
        this.resetCart();
        this.fetchData('orders');
      } catch (error) {
        console.error("Ошибка при создании заказа:", error.message);
        alert(error.message || "Не удалось создать заказ.");
      }
    },
    resetCart() {
      this.cart = [];
      this.bookingDate = '';
    },
    getServiceName(service_id) {
      const service = this.services.find(s => s.service_id === service_id);
      return service ? service.service_name : 'Неизвестная услуга';
    },
    getServicePrice(service_id) {
      const service = this.services.find(s => s.service_id === service_id);
      return service ? service.price : '0';
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
  max-height: 100vh;
}

.orders-section,
.new-order-section,
.cart-section {
  background-color: #f9f9f9;
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

.order-table,
.cart-table {
  width: 100%;
  border-collapse: collapse;
}

.order-table th,
.order-table td,
.cart-table th,
.cart-table td {
  padding: 12px;
  border: 1px solid #ddd;
  text-align: left;
}

.order-table th,
.cart-table th {
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

form label,
.date-selection label {
  font-weight: bold;
  margin-top: 10px;
  display: block;
}

form input,
form select,
.date-selection input {
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
  margin-top: 10px;
}

.btn:hover {
  background-color: #45a049;
}

.remove-btn {
  background-color: #f44336;
  color: white;
  padding: 5px 10px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.remove-btn:hover {
  background-color: #d32f2f;
}

.submit-btn {
  background-color: #2196F3;
}

.submit-btn:hover {
  background-color: #1976D2;
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
