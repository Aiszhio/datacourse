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
          <th>№</th> <!-- Заголовок для нумерации -->
          <th>
            Дата оформления
            <button @click="sortUpcoming('orderDate')">
              Сортировать {{ getSortIcon('upcoming', 'orderDate') }}
            </button>
          </th>
          <th>
            Дата окончания
            <button @click="sortUpcoming('orderEnd')">
              Сортировать {{ getSortIcon('upcoming', 'orderEnd') }}
            </button>
          </th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(order, index) in displayedUpcomingOrders" :key="order.id">
          <td>{{ index + 1 }}</td> <!-- Нумерация строк -->
          <td>{{ formatDate(order.orderDate) }}</td>
          <td>{{ formatDate(order.orderEnd) }}</td>
        </tr>
        </tbody>
      </table>
      <div v-else class="no-orders">Нет предстоящих заказов</div>

      <!-- Кнопка "Показать ещё" для предстоящих заказов -->
      <button
          v-if="canLoadMoreUpcoming"
          @click="loadMoreUpcoming"
          class="btn show-more-btn"
      >
        Показать ещё
      </button>
    </div>

    <!-- Раздел "История выполненных заказов" -->
    <div class="orders-section">
      <h3>История выполненных заказов</h3>
      <table v-if="expiredOrders.length > 0" class="orders-table">
        <thead>
        <tr>
          <th>№</th> <!-- Заголовок для нумерации -->
          <th>
            Имя клиента
            <button @click="sortExpired('clientName')">
              Сортировать {{ getSortIcon('expired', 'clientName') }}
            </button>
          </th>
          <th>
            Название услуги
            <button @click="sortExpired('service')">
              Сортировать {{ getSortIcon('expired', 'service') }}
            </button>
          </th>
          <th>
            Время оформления
            <button @click="sortExpired('orderDate')">
              Сортировать {{ getSortIcon('expired', 'orderDate') }}
            </button>
          </th>
          <th>
            Время окончания
            <button @click="sortExpired('receiveDate')">
              Сортировать {{ getSortIcon('expired', 'receiveDate') }}
            </button>
          </th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(order, index) in displayedExpiredOrders" :key="order.id">
          <td>{{ index + 1 }}</td> <!-- Нумерация строк -->
          <td>{{ order.clientName }}</td>
          <td>{{ order.service }}</td>
          <td>{{ formatDate(order.orderDate) }}</td>
          <td>{{ formatDate(order.receiveDate) }}</td>
        </tr>
        </tbody>
      </table>
      <div v-else class="no-orders">Нет выполненных заказов</div>

      <!-- Кнопка "Показать ещё" для истории заказов -->
      <button
          v-if="canLoadMoreExpired"
          @click="loadMoreExpired"
          class="btn show-more-btn"
      >
        Показать ещё
      </button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Worker',
  data() {
    return {
      workerName: '',             // Имя сотрудника
      upcomingOrders: [],         // Список предстоящих заказов
      expiredOrders: [],          // Список выполненных заказов
      loading: true,              // Индикатор загрузки заказов
      displayCountUpcoming: 10,   // Количество отображаемых предстоящих заказов
      displayCountExpired: 10,    // Количество отображаемых выполненных заказов
      sortConfigUpcoming: {       // Конфигурация сортировки для предстоящих заказов
        key: 'orderDate',
        order: 'asc', // 'asc' или 'desc'
      },
      sortConfigExpired: {        // Конфигурация сортировки для выполненных заказов
        key: 'orderDate',
        order: 'asc',
      },
    };
  },
  computed: {
    // Отображаемые предстоящие заказы с учётом пагинации и сортировки
    displayedUpcomingOrders() {
      return this.sortedUpcomingOrders.slice(0, this.displayCountUpcoming);
    },
    // Отображаемые выполненные заказы с учётом пагинации и сортировки
    displayedExpiredOrders() {
      return this.sortedExpiredOrders.slice(0, this.displayCountExpired);
    },
    // Отсортированные предстоящие заказы
    sortedUpcomingOrders() {
      return this.upcomingOrders.slice().sort((a, b) => {
        let modifier = this.sortConfigUpcoming.order === 'asc' ? 1 : -1;
        if (a[this.sortConfigUpcoming.key] < b[this.sortConfigUpcoming.key]) return -1 * modifier;
        if (a[this.sortConfigUpcoming.key] > b[this.sortConfigUpcoming.key]) return 1 * modifier;
        return 0;
      });
    },
    // Отсортированные выполненные заказы
    sortedExpiredOrders() {
      return this.expiredOrders.slice().sort((a, b) => {
        let modifier = this.sortConfigExpired.order === 'asc' ? 1 : -1;
        if (a[this.sortConfigExpired.key] < b[this.sortConfigExpired.key]) return -1 * modifier;
        if (a[this.sortConfigExpired.key] > b[this.sortConfigExpired.key]) return 1 * modifier;
        return 0;
      });
    },
    // Проверка, можно ли загрузить ещё предстоящие заказы
    canLoadMoreUpcoming() {
      return this.displayCountUpcoming < this.upcomingOrders.length;
    },
    // Проверка, можно ли загрузить ещё выполненные заказы
    canLoadMoreExpired() {
      return this.displayCountExpired < this.expiredOrders.length;
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

        // Проверка формата данных
        if (!data.upcoming || !Array.isArray(data.upcoming)) {
          throw new Error('Некорректный формат данных для предстоящих заказов');
        }

        if (!data.expired || !Array.isArray(data.expired)) {
          throw new Error('Некорректный формат данных для выполненных заказов');
        }

        // Маппинг предстоящих заказов
        this.upcomingOrders = data.upcoming.map(order => ({
          id: order.id,
          bookingId: order.booking_id,
          employeeId: order.employee_id,
          clientId: order.client_id,
          orderDate: order.orderDate,
          orderEnd: order.orderEnd
        }));

        // Маппинг выполненных заказов
        this.expiredOrders = data.expired.map(order => ({
          id: order.id,
          clientName: order.clientName,
          service: order.service,
          orderDate: order.orderDate,
          receiveDate: order.receiveDate
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
      const options = { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' };
      const date = new Date(dateString);
      return date.toLocaleString('ru-RU', options);
    },
    // Загрузка дополнительных предстоящих заказов
    loadMoreUpcoming() {
      this.displayCountUpcoming += 10;
    },
    // Загрузка дополнительных выполненных заказов
    loadMoreExpired() {
      this.displayCountExpired += 10;
    },
    // Метод для сортировки предстоящих заказов
    sortUpcoming(key) {
      if (this.sortConfigUpcoming.key === key) {
        // Если уже сортируем по этому ключу, переключаем порядок
        this.sortConfigUpcoming.order = this.sortConfigUpcoming.order === 'asc' ? 'desc' : 'asc';
      } else {
        // Иначе устанавливаем новый ключ и порядок по возрастанию
        this.sortConfigUpcoming.key = key;
        this.sortConfigUpcoming.order = 'asc';
      }
    },
    // Метод для сортировки выполненных заказов
    sortExpired(key) {
      if (this.sortConfigExpired.key === key) {
        // Если уже сортируем по этому ключу, переключаем порядок
        this.sortConfigExpired.order = this.sortConfigExpired.order === 'asc' ? 'desc' : 'asc';
      } else {
        // Иначе устанавливаем новый ключ и порядок по возрастанию
        this.sortConfigExpired.key = key;
        this.sortConfigExpired.order = 'asc';
      }
    },
    // Метод для отображения иконки сортировки
    getSortIcon(table, key) {
      let config;
      if (table === 'upcoming') {
        config = this.sortConfigUpcoming;
      } else if (table === 'expired') {
        config = this.sortConfigExpired;
      }

      if (config.key !== key) {
        return '';
      }

      return config.order === 'asc' ? '↑' : '↓';
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

.orders-table th button {
  background-color: #4CAF50;
  color: white;
  padding: 8px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  margin: 1px;
}

.orders-table th button:hover {
  text-decoration: underline;
}
</style>
