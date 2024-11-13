<template>
  <div class="user-dashboard">
    <h2>Панель пользователя</h2>
    <h3>Добро пожаловать, {{ userName }}!</h3>

    <!-- Раздел "Ваши заказы" -->
    <div class="orders-section">
      <h3>Ваши заказы</h3>
      <div v-if="loading" class="loading">Загрузка ваших заказов...</div>

      <div v-if="!loading && limitedOrders.length > 0" class="order-table">
        <table>
          <thead>
          <tr>
            <th>Номер заказа</th>
            <th>Имя фотографа</th>
            <th>Название услуги</th>
            <th>Дата оформления</th>
            <th>Дата получения</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="order in limitedOrders" :key="order.order_id">
            <td>{{ order.order_id }}</td>
            <td>{{ order.employee_name || 'Не назначен' }}</td>
            <td>{{ order.service_name }}</td>
            <td>{{ formatDate(order.order_date) }}</td>
            <td>{{ formatDate(order.receipt_date) }}</td>
          </tr>
          </tbody>
        </table>
      </div>

      <div v-if="!loading && limitedOrders.length === 0" class="no-orders">
        У вас нет текущих заказов.
      </div>
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
      userName: '', // User name
      userRole: '', // User role
      orders: [],   // Orders list
      loading: true // Loading indicator
    };
  },
  computed: {
    limitedOrders() {
      return this.orders.slice(0, 3); // Limit orders display to 3
    }
  },
  async mounted() {
    await this.fetchUserData();
    await this.fetchOrders();
  },
  methods: {
    async fetchUserData() {
      try {
        const response = await fetch('http://localhost:8080/api/user', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
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
        console.error('Ошибка:', error.message);
        this.$router.push({ name: 'home' });
      }
    },
    async fetchOrders() {
      try {
        const response = await fetch('http://localhost:8080/api/orders', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
        });

        if (!response.ok) {
          throw new Error('Ошибка при получении заказов');
        }

        const data = await response.json();
        this.orders = data;
      } catch (error) {
        console.error('Ошибка при загрузке заказов:', error.message);
        alert('Не удалось загрузить заказы.');
      } finally {
        this.loading = false;
      }
    },
    formatDate(dateString) {
      if (!dateString) return "Дата не указана";
      const date = new Date(dateString);
      return date.toLocaleDateString('ru-RU', { year: 'numeric', month: 'long', day: 'numeric' });
    },
    goToCreateOrder() {
      this.$router.push({ name: 'ClientOrders' });
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
  margin-top: 20px;
  border-collapse: collapse;
  width: 100%;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.order-table table {
  width: 100%;
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

.loading {
  text-align: center;
  color: #555;
  font-size: 1.2em;
}

.no-orders {
  text-align: center;
  color: #999;
  font-size: 1.2em;
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
  padding: 10px;
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



