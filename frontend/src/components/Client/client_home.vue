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
            <th>Номер клиента</th>
            <th>Номер сотрудника</th>
            <th>Название услуги</th>
            <th>Дата оформления</th>
            <th>Дата получения</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="order in limitedOrders" :key="order.OrderID">
            <td>{{ order.OrderID }}</td>
            <td>{{ order.ClientID }}</td>
            <td>{{ order.EmployeeID || 'Не назначен' }}</td>
            <td>{{ order.ServiceName }}</td>
            <td>{{ formatDate(order.OrderDate) }}</td>
            <td>{{ formatDate(order.ReceiptDate) }}</td>
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
      userName: '',          // Имя пользователя
      userRole: '',          // Роль пользователя
      orders: [],            // Список заказов
      loading: true          // Индикатор загрузки
    };
  },
  computed: {
    limitedOrders() {
      return this.orders.slice(0, 3); // Ограничиваем список заказов до 3, если это нужно
    }
  },
  async mounted() {
    await this.fetchUserData(); // Загружаем данные пользователя
    await this.fetchOrders();   // Загружаем заказы
  },
  methods: {
    async fetchUserData() {
      try {
        const response = await fetch('http://localhost:8080/api/user', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include', // Отправляем куки с запросом
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
          credentials: 'include', // Отправляем куки с запросом
        });

        if (!response.ok) {
          throw new Error('Ошибка при получении заказов');
        }

        const data = await response.json();
        this.orders = data; // Обновляем список заказов
      } catch (error) {
        console.error('Ошибка при загрузке заказов:', error.message);
        alert('Не удалось загрузить заказы.');
      } finally {
        this.loading = false; // Завершаем загрузку
      }
    },
    formatDate(dateString) {
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



