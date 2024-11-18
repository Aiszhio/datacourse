<template>
  <div class="user-dashboard">
    <h2>Панель пользователя</h2>
    <h3>Добро пожаловать, {{ userName }}!</h3>

    <!-- Раздел "Ваши заказы" -->
    <div class="orders-section">
      <h3>Ваши заказы</h3>

      <div v-if="loading" class="loading">Загрузка ваших заказов...</div>

      <template v-else>
        <table v-if="filteredOrders.length" class="order-table">
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
          <tr v-for="(order, index) in filteredOrders" :key="order.order_id">
            <td>{{ index + 1 }}</td>
            <td>{{ order.employee_name || 'Не назначен' }}</td>
            <td>{{ order.service_name }}</td>
            <td>{{ formatDate(order.order_date) }}</td>
            <td>{{ formatDate(order.receipt_date) }}</td>
          </tr>
          </tbody>
        </table>

        <div v-else class="no-orders">
          У вас нет текущих заказов.
        </div>
      </template>
    </div>

    <!-- Навигационные кнопки -->
    <div class="card-panel">
      <div @click="goToCreateOrder" class="card">
        <h4>Сделать заказ</h4>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'User',
  data() {
    return {
      userName: '',     // Имя пользователя
      userRole: '',     // Роль пользователя
      orders: [],       // Список всех заказов
      loading: true     // Индикатор загрузки
    };
  },
  computed: {
    filteredOrders() {
      const today = new Date();
      today.setHours(0, 0, 0, 0); // Нормализация времени до начала дня
      return this.orders.filter(order => new Date(order.receipt_date) >= today);
    }
  },
  mounted() {
    this.initializeData();
  },
  methods: {
    async initializeData() {
      try {
        await Promise.all([this.fetchUserData(), this.fetchOrders()]);
      } catch (error) {
        console.error('Ошибка инициализации данных:', error.message);
        this.$router.push({ name: 'home' });
      }
    },
    async fetchUserData() {
      try {
        const response = await fetch('http://localhost:8080/api/user', {
          method: 'GET',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
        });

        if (response.status === 401) {
          this.$router.push({ name: 'home' });
          return;
        }

        if (!response.ok) {
          throw new Error('Ошибка при получении данных пользователя');
        }

        const data = await response.json();
        this.userName = data.name;
        this.userRole = data.role;
      } catch (error) {
        console.error('Ошибка при получении данных пользователя:', error.message);
        this.$router.push({ name: 'home' });
      }
    },
    async fetchOrders() {
      try {
        const response = await fetch('http://localhost:8080/api/orders', {
          method: 'GET',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
        });

        if (!response.ok) {
          throw new Error('Ошибка при получении заказов');
        }

        const data = await response.json();
        this.orders = data;
      } catch (error) {
        console.error('Ошибка при загрузке заказов:', error.message);
        alert('Не удалось загрузить историю заказов.');
      } finally {
        this.loading = false;
      }
    },
    formatDate(dateString) {
      if (!dateString) return "Дата не указана";
      return new Date(dateString).toLocaleDateString("ru-RU", {
        year: "numeric",
        month: "long",
        day: "numeric",
      });
    },
    goToCreateOrder() {
      this.$router.push({ name: "ClientOrders" });
    }
  }
};
</script>

<style scoped>
.user-dashboard {
  padding: 20px;
  margin: 0 auto;
  max-width: 900px;
}

.orders-section {
  margin-bottom: 20px;
}

.order-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
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

.card-panel {
  display: flex;
  gap: 20px;
  justify-content: center;
  margin-top: 20px;
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
</style>
