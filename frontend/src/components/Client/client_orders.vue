<template>
  <div class="client-orders-container">
    <h2>История заказов</h2>

    <!-- Секция с заказами -->
    <div class="orders-section">
      <div v-if="loading" class="loading">Загрузка истории заказов...</div>

      <table v-else-if="pastOrders.length" class="order-table">
        <thead>
        <tr>
          <th>№</th>
          <th>Имя фотографа</th>
          <th>Название услуг</th>
          <th>Дата оформления</th>
          <th>Дата завершения</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(order, index) in pastOrders" :key="order.order_id">
          <td>{{ index + 1 }}</td>
          <td>{{ order.employee_name || 'Не назначен' }}</td>
          <td>{{ order.service_name || 'Нет услуг' }}</td>
          <td>{{ formatDate(order.order_date) }}</td>
          <td>{{ formatDate(order.receipt_date) }}</td>
        </tr>
        </tbody>
      </table>

      <div v-else class="no-orders">У вас нет завершенных заказов.</div>
    </div>

    <!-- Секция бронирования -->
    <div class="booking-section">
      <h3>Забронировать время</h3>
      <form @submit.prevent="submitBooking">
        <div class="form-group">
          <label for="bookingDate">Дата и время бронирования:</label>
          <input
              type="datetime-local"
              v-model="bookingDate"
              id="bookingDate"
              :min="currentDateTime"
              required
          />
        </div>
        <button type="submit" class="btn" :disabled="!bookingDate">Забронировать</button>
      </form>
    </div>

    <!-- Навигация -->
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
      bookingDate: '',
      currentDateTime: '',
      bookerFullName: '',
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
    this.fetchUserData();
    this.fetchOrders();
    this.interval = setInterval(this.updateCurrentDateTime, 60000);
  },
  beforeUnmount() {
    clearInterval(this.interval);
  },
  methods: {
    updateCurrentDateTime() {
      this.currentDateTime = new Date().toISOString().slice(0, 16);
    },
    async fetchOrders() {
      this.loading = true;
      try {
        const response = await fetch("http://localhost:8080/api/orders", {
          method: 'GET',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
        });

        if (!response.ok) throw new Error('Ошибка при получении истории заказов');

        const data = await response.json();
        this.orders = data;
      } catch (error) {
        console.error("Ошибка:", error.message);
        this.$swal.fire({
          title: 'Ошибка!',
          text: "Не удалось загрузить историю заказов.",
          icon: 'error',
          timer: 3000,
          showConfirmButton: false,
          toast: true,
          position: 'top-end',
        });
      } finally {
        this.loading = false;
      }
    },
    async fetchUserData() {
      try {
        const response = await fetch("http://localhost:8080/api/user", {
          method: "GET",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
        });

        if (!response.ok) throw new Error("Не удалось получить данные пользователя.");

        const userData = await response.json();
        this.bookerFullName = userData.name || '';
      } catch (error) {
        console.error("Ошибка при получении данных пользователя:", error.message);
        this.$swal.fire({
          title: 'Ошибка!',
          text: "Не удалось получить информацию о пользователе. Пожалуйста, попробуйте снова позже.",
          icon: 'error',
          timer: 3000,
          showConfirmButton: false,
          toast: true,
          position: 'top-end',
        });
      }
    },
    async submitBooking() {
      if (!this.bookingDate) {
        this.$swal.fire({
          title: 'Предупреждение',
          text: "Пожалуйста, выберите дату бронирования.",
          icon: 'warning',
          timer: 3000,
          showConfirmButton: false,
          toast: true,
          position: 'top-end',
        });
        return;
      }

      try {
        // Добавляем секунды и смещение часового пояса для Москвы
        const bookingTime = `${this.bookingDate}:00+03:00`;

        const payload = {
          booking_type: "Онлайн",
          booking_time: bookingTime,
          booker_full_name: this.bookerFullName,
        };

        const response = await fetch("http://localhost:8080/api/createOrder", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify(payload),
        });

        const result = await response.json();

        if (!response.ok) {
          throw new Error(result.error || "Не удалось забронировать время.");
        }

        // Отображение успешного бронирования
        this.$swal.fire({
          title: 'Успех!',
          text: "Вы успешно сделали бронь!",
          icon: 'success',
          timer: 3000,
          showConfirmButton: false,
          toast: true,
          position: 'top-end',
        });

        this.resetBookingForm();
      } catch (error) {
        console.error("Ошибка при бронировании:", error.message);
        this.$swal.fire({
          title: 'Ошибка!',
          text: error.message || "Не удалось забронировать время.",
          icon: 'error',
          timer: 3000,
          showConfirmButton: false,
          toast: true,
          position: 'top-end',
        });
      }
    },
    resetBookingForm() {
      this.bookingDate = '';
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
  margin: 20px;
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
