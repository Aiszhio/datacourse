<template>
  <div class="worker-dashboard">
    <h2>Панель фотографа</h2>
    <h3>Добро пожаловать, {{ workerName }}!</h3>

    <!-- Раздел "Предстоящие заказы" -->
    <div class="orders-section">
      <h3>Предстоящие заказы</h3>
      <div v-if="loading" class="loading">Загрузка заказов...</div>

      <table v-else-if="upcomingOrders.length > 0" class="orders-table">
        <thead>
        <tr>
          <th>№</th>
          <th>Имя клиента</th>
          <th>Название услуги</th>
          <th>Дата оформления</th>
          <th>Дата получения</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(order, index) in displayedUpcomingOrders" :key="order.OrderID">
          <td>{{ index + 1 }}</td>
          <td>{{ order.ClientName }}</td>
          <td>{{ order.ServiceName }}</td>
          <td>{{ formatDate(order.OrderDate) }}</td>
          <td>{{ formatDate(order.ReceiptDate) }}</td>
        </tr>
        </tbody>
      </table>
      <div v-else class="no-orders">Нет предстоящих заказов</div>

      <!-- Кнопка "Показать ещё" для предстоящих заказов -->
      <button
          v-if="canLoadMoreUpcoming"
          @click="loadMoreUpcoming"
          class="btn show-more-btn">
        Показать ещё
      </button>
    </div>

    <!-- Раздел "История выполненных заказов" -->
    <div class="orders-section">
      <h3>История выполненных заказов</h3>
      <table v-if="completedOrders.length > 0" class="orders-table">
        <thead>
        <tr>
          <th>№</th>
          <th>Имя клиента</th>
          <th>Название услуги</th>
          <th>Дата оформления</th>
          <th>Дата получения</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(order, index) in displayedCompletedOrders" :key="order.OrderID">
          <td>{{ index + 1 }}</td>
          <td>{{ order.ClientName }}</td>
          <td>{{ order.ServiceName }}</td>
          <td>{{ formatDate(order.OrderDate) }}</td>
          <td>{{ formatDate(order.ReceiptDate) }}</td>
        </tr>
        </tbody>
      </table>
      <div v-else class="no-orders">Нет выполненных заказов</div>

      <!-- Кнопка "Показать ещё" для истории заказов -->
      <button
          v-if="canLoadMoreCompleted"
          @click="loadMoreCompleted"
          class="btn show-more-btn">
        Показать ещё
      </button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'WorkerDashboard',
  data() {
    return {
      workerName: '',      // Имя сотрудника
      orders: [],          // Список заказов
      loading: true,       // Индикатор загрузки заказов
      displayCountUpcoming: 10,  // Количество отображаемых предстоящих заказов
      displayCountCompleted: 10, // Количество отображаемых выполненных заказов
    };
  },
  computed: {
    // Фильтрация предстоящих заказов
    upcomingOrders() {
      const today = new Date();
      today.setHours(0, 0, 0, 0); // Нормализация времени до начала дня
      return this.orders.filter(order => new Date(order.ReceiptDate) >= today);
    },
    // Фильтрация выполненных заказов
    completedOrders() {
      const today = new Date();
      today.setHours(0, 0, 0, 0);
      return this.orders.filter(order => new Date(order.ReceiptDate) < today);
    },
    // Отображаемые предстоящие заказы с учётом пагинации
    displayedUpcomingOrders() {
      return this.upcomingOrders.slice(0, this.displayCountUpcoming);
    },
    // Отображаемые выполненные заказы с учётом пагинации
    displayedCompletedOrders() {
      return this.completedOrders.slice(0, this.displayCountCompleted);
    },
    // Проверка, можно ли загрузить ещё предстоящие заказы
    canLoadMoreUpcoming() {
      return this.displayCountUpcoming < this.upcomingOrders.length;
    },
    // Проверка, можно ли загрузить ещё выполненные заказы
    canLoadMoreCompleted() {
      return this.displayCountCompleted < this.completedOrders.length;
    }
  },
  async mounted() {
    await this.fetchWorkerData();
    await this.fetchOrders();
  },
  methods: {
    // Получение данных сотрудника
    async fetchWorkerData() {
      try {
        const response = await fetch('http://localhost:8080/api/user', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json'
          },
          credentials: 'include'
        });

        if (!response.ok) {
          throw new Error('Ошибка при получении имени сотрудника');
        }

        const data = await response.json();
        this.workerName = data.name;
      } catch (error) {
        console.error('Ошибка при получении имени сотрудника:', error.message);
        alert('Не удалось загрузить имя сотрудника.');
      }
    },
    // Получение заказов
    // Получение заказов
    async fetchOrders() {
      try {
        const response = await fetch('http://localhost:8080/api/orders/worker', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json'
          },
          credentials: 'include'
        });

        if (!response.ok) {
          throw new Error('Ошибка при получении заказов');
        }

        const data = await response.json();

        // Убедитесь, что data.data существует и является массивом
        if (!data.data || !Array.isArray(data.data)) {
          throw new Error('Некорректный формат данных');
        }

        // Преобразование полей из snake_case в PascalCase
        this.orders = data.data.map(order => ({
          OrderID: order.order_id,
          ClientName: order.client_name,
          ServiceName: order.service_name,
          OrderDate: order.order_date,
          ReceiptDate: order.receipt_date
        }));

      } catch (error) {
        console.error('Ошибка при загрузке заказов:', error.message);
        alert('Не удалось загрузить заказы.');
      } finally {
        this.loading = false;
      }
    },
    // Форматирование даты с использованием встроенных средств JavaScript
    formatDate(dateString) {
      const options = {year: 'numeric', month: 'long', day: 'numeric'};
      const date = new Date(dateString);
      return date.toLocaleDateString('ru-RU', options);
    },
    // Загрузка дополнительных предстоящих заказов
    loadMoreUpcoming() {
      this.displayCountUpcoming += 10;
    },
    // Загрузка дополнительных выполненных заказов
    loadMoreCompleted() {
      this.displayCountCompleted += 10;
    }
  }
};
</script>

<style scoped>
.worker-dashboard {
  padding: 20px;
  margin: 0 auto;
  max-width: 900px;
  max-height: 100vh;
  background-color: #ffffff;
  border-radius: 10px;
}

.orders-section {
  background-color: #f9f9f9;
  padding: 20px;
  margin: 20px 0;
  border-radius: 10px;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
}

.orders-section h3 {
  margin-bottom: 15px;
  color: #333333;
}

.orders-table {
  width: 100%;
  border-collapse: collapse;
}

.orders-table th,
.orders-table td {
  padding: 12px;
  border: 1px solid #dddddd;
  text-align: left;
}

.orders-table th {
  background-color: #eaeaea;
  font-weight: bold;
  color: #555555;
}

.orders-table tbody tr:nth-child(even) {
  background-color: #f2f2f2;
}

.no-orders {
  text-align: center;
  color: #888888;
  font-style: italic;
  margin-top: 10px;
}

.loading {
  text-align: center;
  color: #555555;
  font-size: 1.2em;
  margin-top: 10px;
}

.btn {
  background-color: #4CAF50;
  color: white;
  padding: 10px 16px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 1em;
  transition: background-color 0.3s ease;
  margin: 15px auto 0;
  display: block;
}

.btn:hover {
  background-color: #45a049;
}

.btn:disabled {
  background-color: #a5d6a7;
  cursor: not-allowed;
}
</style>
