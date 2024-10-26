<script>
export default {
  name: 'worker_orders',
  data() {
    return {
      orders: [],  // Список доступных заказов для фотографа
      loading: false,  // Флаг для индикации загрузки
      sortKey: '',  // Ключ для сортировки
      sortOrder: 1  // Порядок сортировки: 1 для возрастания, -1 для убывания
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
      // Логика для закрепления заказа за фотографом
    },
    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleDateString('ru-RU', { year: 'numeric', month: 'long', day: 'numeric' });
    },
    goToHome() {
      this.$router.push({ name: 'WorkerDashboard' });
    },
    sortBy(key) {
      // Если уже сортируем по этому ключу, меняем порядок
      if (this.sortKey === key) {
        this.sortOrder *= -1;
      } else {
        this.sortKey = key;
        this.sortOrder = 1;
      }

      this.orders.sort((a, b) => {
        let valueA = a[key];
        let valueB = b[key];

        if (typeof valueA === 'string') {
          valueA = valueA.toLowerCase();
          valueB = valueB.toLowerCase();
        }

        if (valueA < valueB) return -1 * this.sortOrder;
        if (valueA > valueB) return 1 * this.sortOrder;
        return 0;
      });
    }
  }
};
</script>

<template>
  <div class="photographer-orders-container">
    <h2> История заказов </h2>
    <div v-if="loading" class="loading">Загрузка доступных заказов...</div>

    <div v-else>
      <!-- Кнопки для сортировки -->
      <div class="sort-buttons">
        <button @click="sortBy('OrderID')" class="btn">Сортировать по ID заказа</button>
        <button @click="sortBy('ClientName')" class="btn">Сортировать по клиенту</button>
        <button @click="sortBy('ServiceName')" class="btn">Сортировать по услуге</button>
        <button @click="sortBy('OrderDate')" class="btn">Сортировать по дате заказа</button>
      </div>

      <!-- Таблица заказов -->
      <table class="orders-table">
        <thead>
        <tr>
          <th>ID Заказа</th>
          <th>Клиент</th>
          <th>Услуга</th>
          <th>Дата заказа</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="order in orders" :key="order.OrderID">
          <td>{{ order.OrderID }}</td>
          <td>{{ order.ClientName }}</td>
          <td>{{ order.ServiceName }}</td>
          <td>{{ formatDate(order.OrderDate) }}</td>
          <td>{{ formatDate(order.ReceiptDate) }}</td>
          <td><button @click="selectOrder(order)" class="btn">Выбрать заказ</button></td>
        </tr>
        </tbody>
      </table>
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
  max-width: 900px;
  box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.1);
  border-radius: 10px;
  text-align: left;
}

h2 {
  text-align: center;
  color: #333;
  margin-bottom: 20px;
}

.sort-buttons {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.orders-table {
  width: 100%;
  border-collapse: collapse;
}

.orders-table th,
.orders-table td {
  padding: 12px;
  border: 1px solid #ddd;
  text-align: left;
}

.orders-table th {
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

.navigation-buttons {
  display: flex;
  justify-content: space-around;
  margin-top: 20px;
}
</style>
