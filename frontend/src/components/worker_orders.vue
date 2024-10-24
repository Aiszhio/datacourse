<script>
export default {
  name: 'worker_orders',
  data() {
    return {
      orders: [],  // Список доступных заказов для фотографа
      loading: false  // Флаг для индикации загрузки
    };
  },
  mounted() {
    this.fetchOrders();  // Загружаем доступные заказы при загрузке страницы
  },
  methods: {
    async fetchOrders() {
      this.loading = true;  // Устанавливаем флаг загрузки
      try {
        const response = await fetch('http://localhost:8080/api/orders/available', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer your-auth-token'  // Если требуется токен авторизации
          }
        });

        if (!response.ok) {
          throw new Error('Ошибка при получении списка заказов');
        }

        const data = await response.json();
        this.orders = data;  // Записываем полученные данные
      } catch (error) {
        console.error('Ошибка:', error.message);
        alert('Не удалось загрузить список заказов. Попробуйте позже.');
      } finally {
        this.loading = false;  // Отключаем флаг загрузки после завершения запроса
      }
    },
    selectOrder(order) {
      alert(`Заказ #${order.OrderID} выбран для выполнения.`);
      // Здесь можно добавить логику для закрепления заказа за фотографом
    },
    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleDateString('ru-RU', {year: 'numeric', month: 'long', day: 'numeric'});
    },
    goToHome() {
      this.$router.push({name: 'WorkerDashboard'});
    },
  }
};
</script>

<template>
  <div class="photographer-orders-container">
    <h2> История заказов </h2>
    <div v-if="loading" class="loading">Загрузка доступных заказов...</div>
    <div v-else>
      <div v-if="orders.length > 0" class="orders-list">
        <div class="order" v-for="order in orders" :key="order.OrderID">
          <h3>Заказ #{{ order.OrderID }}</h3>
          <p><strong>Клиент:</strong> {{ order.ClientName }}</p>
          <p><strong>Услуга:</strong> {{ order.ServiceName }}</p>
          <p><strong>Дата заказа:</strong> {{ formatDate(order.OrderDate) }}</p>
          <p><strong>Дата выдачи:</strong> {{ formatDate(order.ReceiptDate) }}</p>
          <button @click="selectOrder(order)" class="btn">Выбрать заказ</button>
        </div>
      </div>
      <div v-else class="no-orders">Нет доступных заказов.</div>
    </div>

    <!-- Кнопки для перемещения по сайту -->
    <div class="navigation-buttons">
      <button @click="goToHome" class="btn">На главную</button>
    </div>
  </div>
</template>

<style scoped>
.photographer-orders-container {
  background-color: #f9f9f9;
  padding: 20px;
  margin: 20px auto;
  width: 100%;
  max-width: 700px;
  box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.1);
  border-radius: 10px;
  text-align: left;
}

h2 {
  text-align: center;
  color: #333;
  margin-bottom: 20px;
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

.order h3 {
  margin-bottom: 10px;
  color: #4CAF50;
}

.order p {
  margin: 5px 0;
  color: #555;
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
  text-align: center;
}

.btn:hover {
  background-color: #45a049;
}

.navigation-buttons {
  display: flex;
  justify-content: space-around;
  margin-top: 20px;
}

.navigation-buttons .btn {
  width: 150px;
}
</style>
