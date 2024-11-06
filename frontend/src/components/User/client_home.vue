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
            <th>Статус</th>
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
            <td>{{ order.Status }}</td>
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
      userName: 'Имя Пользователя', // Имя пользователя
      orders: [
        {
          OrderID: 123,
          ClientID: 456,
          EmployeeID: 789, // Номер сотрудника, назначенного на заказ
          ServiceName: 'Фотосессия',
          OrderDate: '2024-10-01',
          ReceiptDate: '2024-10-10',
          Status: 'В процессе'
        }
      ], // Пример текущего заказа пользователя
      loading: false, // Устанавливаем на false, чтобы отключить индикатор загрузки
      orderLimit: 3 // Ограничение на количество отображаемых заказов
    };
  },
  computed: {
    limitedOrders() {
      return this.orders.slice(0, this.orderLimit);
    }
  },
  methods: {
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
