<script>
export default {
  name: 'worker',
  data() {
    return {
      workerName: 'Имя Сотрудника', // Имя сотрудника, можно получить из бэкенда
      orders: [], // Заказы, назначенные сотруднику
      loading: true, // Флаг загрузки заказов
      orderLimit: 3 // Количество отображаемых заказов
    };
  },
  mounted() {
    this.fetchOrders(); // Загружаем заказы при загрузке страницы
  },
  computed: {
    limitedOrders() {
      // Ограничиваем количество заказов до трех
      return this.orders.slice(0, this.orderLimit);
    }
  },
  methods: {
    async fetchOrders() {
      try {
        console.log("Начинается загрузка заказов...");
        const response = await fetch('http://localhost:8080/api/worker/orders', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer your-auth-token' // Если требуется токен авторизации
          }
        });

        if (!response.ok) {
          throw new Error('Ошибка при получении заказов');
        }

        const data = await response.json();
        this.orders = data.map(order => ({ ...order, updatedStatus: order.Status }));
        console.log("Загруженные заказы:", this.orders); // Вывод загруженных данных

      } catch (error) {
        console.error('Ошибка при загрузке заказов:', error.message);
        alert('Не удалось загрузить заказы.');
      } finally {
        this.loading = false;
      }
    },
    async updateOrderStatus(order) {
      try {
        const response = await fetch(`http://localhost:8080/api/orders/${order.OrderID}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer your-auth-token'
          },
          body: JSON.stringify({ status: order.updatedStatus })
        });

        if (!response.ok) {
          throw new Error('Ошибка при обновлении статуса заказа');
        }

        alert(`Статус заказа #${order.OrderID} обновлён на "${order.updatedStatus}"`);
        order.Status = order.updatedStatus; // Обновляем локальный статус
      } catch (error) {
        console.error('Ошибка:', error.message);
        alert('Не удалось обновить статус заказа.');
      }
    },
    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleDateString('ru-RU', { year: 'numeric', month: 'long', day: 'numeric' });
    },
    goToWorkerOrders() {
      this.$router.push({ name: 'WorkerOrders' }); // Переход на страницу с заказами фотографа
    }
  }
};
</script>

<template>
  <div class="worker-dashboard">
    <h2>Панель сотрудника</h2>
    <h3>Добро пожаловать, {{ workerName }}!</h3>

    <!-- Список текущих заказов -->
    <div class="orders-section">
      <h3>Предстоящие заказы</h3>
      <div v-if="loading" class="loading">Загрузка заказов...</div>

      <!-- Таблица заказов (выводится всегда) -->
      <table class="orders-table">
        <thead>
        <tr>
          <th>ID Заказа</th>
          <th>Услуга</th>
          <th>Дата заказа</th>
          <th>Статус</th>
          <th>Обновить статус</th>
        </tr>
        </thead>
        <tbody>
        <tr v-if="!loading && limitedOrders.length > 0" v-for="order in limitedOrders" :key="order.OrderID">
          <td>{{ order.OrderID }}</td>
          <td>{{ order.ServiceName }}</td>
          <td>{{ formatDate(order.OrderDate) }}</td>
          <td>{{ order.Status }}</td>
          <td>
            <select v-model="order.updatedStatus" @change="updateOrderStatus(order)">
              <option value="В процессе">В процессе</option>
              <option value="Завершён">Завершён</option>
            </select>
          </td>
        </tr>

        <!-- Пустые строки для отображения таблицы даже без данных -->
        <tr v-else>
          <td colspan="5" class="no-orders">Нет доступных заказов</td>
        </tr>
        </tbody>
      </table>
      <!-- Кнопки для навигации -->
      <div class="navigation-buttons">
        <button @click="goToWorkerOrders" class="btn">Просмотреть все заказы</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.worker-dashboard {
  padding: 20px;
  margin: 0 auto;
  max-width: 900px;
}

.orders-section {
  background-color: #f9f9f9;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
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

.no-orders {
  text-align: center;
  color: #999;
}

.status-update {
  margin-top: 15px;
}

.btn {
  background-color: #4CAF50;
  color: white;
  padding: 10px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 1em;
  transition: background-color 0.3s ease;
}

.btn:hover {
  background-color: #45a049;
}

.loading {
  text-align: center;
  color: #555;
  font-size: 1.2em;
}

.navigation-buttons {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
