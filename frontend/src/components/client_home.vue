<script>
export default {
  name: 'User',
  data() {
    return {
      userName: 'Имя Пользователя', // Можно получить из бэкенда или сессии
      orders: [], // Текущие заказы пользователя
      loading: true // Флаг загрузки заказов
    };
  },
  mounted() {
    this.fetchOrders(); // Загружаем заказы при загрузке страницы
  },
  methods: {
    async fetchOrders() {
      try {
        // Пример запроса для получения заказов пользователя
        const response = await fetch('http://localhost:8080/api/client/orders', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer your-auth-token' // Токен авторизации, если требуется
          }
        });

        if (!response.ok) {
          throw new Error('Ошибка при получении заказов');
        }

        const data = await response.json();
        this.orders = data; // Записываем полученные заказы
      } catch (error) {
        console.error('Ошибка:', error.message);
        alert('Не удалось загрузить заказы.');
      } finally {
        this.loading = false; // Отключаем флаг загрузки
      }
    },
    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleDateString('ru-RU', { year: 'numeric', month: 'long', day: 'numeric' });
    },
    goToCreateOrder() {
      this.$router.push({ name: 'ClientOrders' }); // Переход на страницу создания заказа
    },
    goToPayInvoices() {
      this.$router.push({ name: 'Payment' }); // Переход на страницу оплаты счетов
    }
  }
};
</script>

<template>
  <div class="user-dashboard">
    <h2>Панель пользователя</h2>

    <!-- Приветственное сообщение -->
    <div class="user-info">
      <h3>Добро пожаловать, {{ userName }}!</h3>
    </div>

    <!-- Список текущих заказов -->
    <div class="orders-section">
      <h3>Ваши заказы</h3>
      <div v-if="loading" class="loading">Загрузка ваших заказов...</div>
      <div v-else>
        <div v-if="orders.length > 0" class="orders-list">
          <div class="order" v-for="order in orders" :key="order.OrderID">
            <h4>Заказ #{{ order.OrderID }}</h4>
            <p><strong>Услуга:</strong> {{ order.ServiceName }}</p>
            <p><strong>Дата заказа:</strong> {{ formatDate(order.OrderDate) }}</p>
            <p><strong>Дата получения:</strong> {{ formatDate(order.ReceiptDate) }}</p>
            <p><strong>Текущий статус:</strong> {{ order.Status }}</p>
          </div>
        </div>
        <div v-else class="no-orders">У вас нет текущих заказов.</div>
      </div>
    </div>

    <!-- Кнопки для навигации -->
    <div class="navigation-buttons">
      <button @click="goToCreateOrder" class="btn">Создать заказ</button>
      <button @click="goToPayInvoices" class="btn">Оплатить счета</button>
    </div>
  </div>
</template>

<style scoped>
.user-dashboard {
  padding: 20px;
  margin: 0 auto;
  max-width: 900px;
}

.user-info {
  background-color: #f9f9f9;
  padding: 20px;
  margin-bottom: 20px;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.orders-section {
  background-color: #f9f9f9;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
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

.order h4 {
  margin-bottom: 10px;
  color: #4CAF50;
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
  margin: 95px;
  width: 30vh;
}

.btn:hover {
  background-color: #45a049;
}

.navigation-buttons {
  display: flex;
  justify-content: space-around;
  margin-top: 20px;
}
</style>
