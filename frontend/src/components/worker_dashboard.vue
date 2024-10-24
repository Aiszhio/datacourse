<script>
export default {
  name: 'worker',
  data() {
    return {
      workerName: 'Имя Сотрудника', // Имя сотрудника, можно получить из бэкенда
      workerPosition: 'Фотограф', // Позиция сотрудника
      orders: [], // Заказы, назначенные сотруднику
      loading: true // Флаг загрузки заказов
    };
  },
  mounted() {
    this.fetchOrders(); // Загружаем заказы при загрузке страницы
  },
  methods: {
    async fetchOrders() {
      try {
        // Пример запроса для получения заказов сотрудника
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
      } catch (error) {
        console.error('Ошибка:', error.message);
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

    <!-- Информация о сотруднике -->
    <div class="worker-info">
      <h3>Добро пожаловать, {{ workerName }}!</h3>
      <p><strong>Ваша позиция:</strong> {{ workerPosition }}</p>
    </div>

    <!-- Список текущих заказов -->
    <div class="orders-section">
      <h3>Предстоящие заказы</h3>
      <div v-if="loading" class="loading">Загрузка заказов...</div>
      <div v-else>
        <div v-if="orders.length > 0" class="orders-list">
          <div class="order" v-for="order in orders" :key="order.OrderID">
            <h4>Заказ #{{ order.OrderID }}</h4>
            <p><strong>Услуга:</strong> {{ order.ServiceName }}</p>
            <p><strong>Дата заказа:</strong> {{ formatDate(order.OrderDate) }}</p>
            <p><strong>Статус:</strong> {{ order.Status }}</p>

            <!-- Обновление статуса заказа -->
            <div class="status-update">
              <label for="status">Обновить статус:</label>
              <select v-model="order.updatedStatus" @change="updateOrderStatus(order)">
                <option value="В процессе">В процессе</option>
                <option value="Завершён">Завершён</option>
              </select>
            </div>
          </div>
        </div>
        <div v-else class="no-orders">У вас нет назначенных заказов.</div>
      </div>
    </div>

    <!-- Кнопки для навигации -->
    <div class="navigation-buttons">
      <button @click="goToWorkerOrders" class="btn">Просмотреть заказы</button> <!-- Добавленная кнопка -->
    </div>
  </div>
</template>

<style scoped>
.worker-dashboard {
  padding: 20px;
  margin: 0 auto;
  max-width: 900px;
}

.worker-info {
  background-color: #f9f9f9;
  padding: 20px;
  margin-bottom: 20px;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.orders-section {
  background-color: #f9f9f9;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.orders-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.order {
  background-color: #ffffff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0px 2px 4px rgba(0, 0, 0, 0.1);
}

.order h4 {
  margin-bottom: 10px;
  color: #4CAF50;
}

.status-update {
  margin-top: 15px;
}

select {
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
  font-size: 1em;
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
  padding: 12px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 1em;
  transition: background-color 0.3s ease;
  margin: 50px;
  width: 50vh;
}

.btn:hover {
  background-color: #45a049;
}

.navigation-buttons {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
