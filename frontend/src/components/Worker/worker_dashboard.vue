<template>
  <div class="worker-dashboard">
    <h2>Панель фотографа</h2>
    <h3>Добро пожаловать, {{ workerName }}!</h3>

    <!-- Раздел "Будущие заказы" -->
    <div class="orders-section">
      <h3>Будущие заказы</h3>

      <!-- Фильтры для будущих заказов -->
      <div class="filters-container">
        <div class="filter-item">
          <label for="filterFutureName">Название заказа:</label>
          <input
              id="filterFutureName"
              v-model="filterFutureName"
              type="text"
              placeholder="По названию заказа"
              class="filter-input"
          />
        </div>
        <div class="filter-item">
          <label for="filterFutureDate">Дата оформления:</label>
          <input
              id="filterFutureDate"
              v-model="filterFutureDate"
              type="date"
              class="filter-input"
          />
        </div>
        <div class="filter-item">
          <label for="filterFutureEndDate">Дата окончания:</label>
          <input
              id="filterFutureEndDate"
              v-model="filterFutureEndDate"
              type="date"
              class="filter-input"
          />
        </div>
      </div>

      <div v-if="loading" class="loading">Загрузка заказов...</div>

      <table v-else-if="filteredFutureOrders.length > 0" class="orders-table">
        <thead>
        <tr>
          <th>
            Название заказа
            <button @click="sortFuture('serviceName')">
              Сортировать {{ getSortIcon('future', 'serviceName') }}
            </button>
          </th>
          <th>
            Дата оформления
            <button @click="sortFuture('orderDate')">
              Сортировать {{ getSortIcon('future', 'orderDate') }}
            </button>
          </th>
          <th>
            Дата окончания
            <button @click="sortFuture('receiptDate')">
              Сортировать {{ getSortIcon('future', 'receiptDate') }}
            </button>
          </th>
          <th>
            Дата завершения
          </th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="order in displayedFutureOrders" :key="order.id">
          <td>{{ order.serviceName }}</td>
          <td>{{ formatDate(order.orderDate) }}</td>
          <td>{{ formatDate(order.receiptDate) }}</td>
          <td>{{ formatDate(order.remainingTime) }}</td>
        </tr>
        </tbody>
      </table>
      <div v-else class="no-orders">Нет будущих заказов</div>

      <!-- Кнопка "Показать ещё" для будущих заказов -->
      <button
          v-if="canLoadMoreFuture"
          @click="loadMoreFuture"
          class="btn show-more-btn"
      >
        Показать ещё
      </button>
    </div>

    <!-- Раздел "История выполненных заказов" -->
    <div class="orders-section">
      <h3>История выполненных заказов</h3>

      <!-- Фильтры для выполненных заказов -->
      <div class="filters-container">
        <div class="filter-item">
          <label for="filterExpiredClient">Имя клиента:</label>
          <input
              id="filterExpiredClient"
              v-model="filterExpiredClient"
              type="text"
              placeholder="По имени клиента"
              class="filter-input"
          />
        </div>
        <div class="filter-item">
          <label for="filterExpiredService">Название услуги:</label>
          <input
              id="filterExpiredService"
              v-model="filterExpiredService"
              type="text"
              placeholder="По названию услуги"
              class="filter-input"
          />
        </div>
        <div class="filter-item">
          <label for="filterExpiredDate">Дата оформления:</label>
          <input
              id="filterExpiredDate"
              v-model="filterExpiredDate"
              type="date"
              class="filter-input"
          />
        </div>
        <div class="filter-item">
          <label for="filterExpiredReceive">Дата окончания:</label>
          <input
              id="filterExpiredReceive"
              v-model="filterExpiredReceiveDate"
              type="date"
              class="filter-input"
          />
        </div>
      </div>

      <table v-if="filteredPastOrders.length > 0" class="orders-table">
        <thead>
        <tr>
          <th>
            Имя клиента
            <button @click="sortPast('clientName')">
              Сортировать {{ getSortIcon('past', 'clientName') }}
            </button>
          </th>
          <th>
            Название услуги
            <button @click="sortPast('serviceName')">
              Сортировать {{ getSortIcon('past', 'serviceName') }}
            </button>
          </th>
          <th>
            Дата оформления
            <button @click="sortPast('orderDate')">
              Сортировать {{ getSortIcon('past', 'orderDate') }}
            </button>
          </th>
          <th>
            Дата окончания
            <button @click="sortPast('receiptDate')">
              Сортировать {{ getSortIcon('past', 'receiptDate') }}
            </button>
          </th>
          <th>
            Дата завершения
          </th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="order in displayedPastOrders" :key="order.id">
          <td>{{ order.clientName }}</td>
          <td>{{ order.serviceName }}</td>
          <td>{{ formatDate(order.orderDate) }}</td>
          <td>{{ formatDate(order.receiptDate) }}</td>
          <td>{{ formatDate(order.remainingTime) }}</td>
        </tr>
        </tbody>
      </table>
      <div v-else class="no-orders">Нет выполненных заказов</div>

      <!-- Кнопка "Показать ещё" для истории заказов -->
      <button
          v-if="canLoadMorePast"
          @click="loadMorePast"
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
      workerName: '',               // Имя сотрудника
      pastOrders: [],               // Список прошедших заказов
      futureOrders: [],             // Список будущих заказов
      loading: true,                // Индикатор загрузки заказов
      displayCountPast: 10,         // Количество отображаемых прошедших заказов
      displayCountFuture: 10,       // Количество отображаемых будущих заказов
      sortConfigPast: {             // Конфигурация сортировки для прошедших заказов
        key: 'orderDate',
        order: 'asc', // 'asc' или 'desc'
      },
      sortConfigFuture: {           // Конфигурация сортировки для будущих заказов
        key: 'orderDate',
        order: 'asc',
      },
      // Фильтры для прошедших заказов
      filterExpiredClient: '',
      filterExpiredService: '',
      filterExpiredDate: '',
      filterExpiredReceiveDate: '',
      // Фильтры для будущих заказов
      filterFutureName: '',
      filterFutureDate: '',
      filterFutureEndDate: '',
    };
  },
  computed: {
    // Фильтрованные прошедшие заказы
    filteredPastOrders() {
      return this.pastOrders.filter(order => {
        const matchesClient = this.filterExpiredClient
            ? order.clientName.toLowerCase().includes(this.filterExpiredClient.toLowerCase())
            : true;
        const matchesService = this.filterExpiredService
            ? order.serviceName.toLowerCase().includes(this.filterExpiredService.toLowerCase())
            : true;
        const matchesDate = this.filterExpiredDate
            ? this.formatDate(order.orderDate, 'yyyy-MM-dd') === this.filterExpiredDate
            : true;
        const matchesReceiveDate = this.filterExpiredReceiveDate
            ? this.formatDate(order.receiptDate, 'yyyy-MM-dd') === this.filterExpiredReceiveDate
            : true;
        return matchesClient && matchesService && matchesDate && matchesReceiveDate;
      });
    },
    // Отображаемые прошедшие заказы с учётом пагинации и сортировки
    displayedPastOrders() {
      return this.sortedPastOrders.slice(0, this.displayCountPast);
    },
    // Фильтрованные будущие заказы
    filteredFutureOrders() {
      return this.futureOrders.filter(order => {
        const matchesName = this.filterFutureName
            ? order.serviceName.toLowerCase().includes(this.filterFutureName.toLowerCase())
            : true;
        const matchesDate = this.filterFutureDate
            ? this.formatDate(order.orderDate, 'yyyy-MM-dd') === this.filterFutureDate
            : true;
        const matchesEndDate = this.filterFutureEndDate
            ? this.formatDate(order.receiptDate, 'yyyy-MM-dd') === this.filterFutureEndDate
            : true;
        return matchesName && matchesDate && matchesEndDate;
      });
    },
    // Отображаемые будущие заказы с учётом пагинации и сортировки
    displayedFutureOrders() {
      return this.sortedFutureOrders.slice(0, this.displayCountFuture);
    },
    // Отсортированные прошедшие заказы
    sortedPastOrders() {
      return this.filteredPastOrders.slice().sort((a, b) => {
        let modifier = this.sortConfigPast.order === 'asc' ? 1 : -1;
        if (a[this.sortConfigPast.key] < b[this.sortConfigPast.key]) return -1 * modifier;
        if (a[this.sortConfigPast.key] > b[this.sortConfigPast.key]) return 1 * modifier;
        return 0;
      });
    },
    // Отсортированные будущие заказы
    sortedFutureOrders() {
      return this.filteredFutureOrders.slice().sort((a, b) => {
        let modifier = this.sortConfigFuture.order === 'asc' ? 1 : -1;
        if (a[this.sortConfigFuture.key] < b[this.sortConfigFuture.key]) return -1 * modifier;
        if (a[this.sortConfigFuture.key] > b[this.sortConfigFuture.key]) return 1 * modifier;
        return 0;
      });
    },
    // Проверка, можно ли загрузить ещё прошедшие заказы
    canLoadMorePast() {
      return this.displayCountPast < this.filteredPastOrders.length;
    },
    // Проверка, можно ли загрузить ещё будущие заказы
    canLoadMoreFuture() {
      return this.displayCountFuture < this.filteredFutureOrders.length;
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
          // Попробуем получить сообщение об ошибке из ответа
          const errorData = await response.json();
          throw new Error(errorData.error || 'Ошибка при получении заказов');
        }

        const data = await response.json();

        console.log('Полученные данные:', JSON.stringify(data, null, 2)); // Для отладки

        // Проверка формата данных
        if (!data.orders || typeof data.orders !== 'object') {
          throw new Error('Некорректный формат данных для заказов');
        }

        // Проверка pastOrders
        if (!Array.isArray(data.orders.pastOrders)) {
          throw new Error('Некорректный формат данных для заказов: pastOrders не массив');
        }

        // Маппинг прошедших заказов
        this.pastOrders = data.orders.pastOrders.map(order => ({
          id: order.id,
          clientName: order.clientName,
          serviceName: order.serviceName,
          orderDate: order.orderDate,
          receiptDate: order.adjustedReceiptDate,
          remainingTime: order.remainingTime || '' // Добавлено для pastOrders
        }));

        // Маппинг будущих заказов
        if (data.orders.futureOrders === null) {
          this.futureOrders = [];
        } else if (Array.isArray(data.orders.futureOrders)) {
          this.futureOrders = data.orders.futureOrders.map(order => ({
            id: order.id,
            clientName: order.clientName,
            serviceName: order.serviceName,
            orderDate: order.orderDate,
            receiptDate: order.adjustedReceiptDate,
            remainingTime: order.remainingTime || '' // Теперь это строка типа "2024-12-17T10:00:00Z"
          }));
        } else {
          throw new Error('Некорректный формат данных для заказов: futureOrders не массив и не null');
        }

      } catch (error) {
        console.error('Ошибка при загрузке заказов:', error.message);
        alert('Не удалось загрузить заказы: ' + error.message);
      } finally {
        this.loading = false;
      }
    },
    // Форматирование даты с возможностью выбора формата
    formatDate(dateString, format = 'full') {
      if (!dateString) return 'Неизвестно';
      const date = new Date(dateString);
      if (isNaN(date)) return 'Неизвестно';

      if (format === 'yyyy-MM-dd') {
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');
        return `${year}-${month}-${day}`;
      }
      const options = { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' };
      return date.toLocaleString('ru-RU', options);
    },
    // Загрузка дополнительных прошедших заказов
    loadMorePast() {
      this.displayCountPast += 10;
    },
    // Загрузка дополнительных будущих заказов
    loadMoreFuture() {
      this.displayCountFuture += 10;
    },
    // Метод для сортировки прошедших заказов
    sortPast(key) {
      if (this.sortConfigPast.key === key) {
        // Если уже сортируем по этому ключу, переключаем порядок
        this.sortConfigPast.order = this.sortConfigPast.order === 'asc' ? 'desc' : 'asc';
      } else {
        // Иначе устанавливаем новый ключ и порядок по возрастанию
        this.sortConfigPast.key = key;
        this.sortConfigPast.order = 'asc';
      }
    },
    // Метод для сортировки будущих заказов
    sortFuture(key) {
      if (this.sortConfigFuture.key === key) {
        // Если уже сортируем по этому ключу, переключаем порядок
        this.sortConfigFuture.order = this.sortConfigFuture.order === 'asc' ? 'desc' : 'asc';
      } else {
        // Иначе устанавливаем новый ключ и порядок по возрастанию
        this.sortConfigFuture.key = key;
        this.sortConfigFuture.order = 'asc';
      }
    },
    // Метод для отображения иконки сортировки
    getSortIcon(table, key) {
      let config;
      if (table === 'past') {
        config = this.sortConfigPast;
      } else if (table === 'future') {
        config = this.sortConfigFuture;
      }

      if (config.key !== key) {
        return '';
      }

      return config.order === 'asc' ? '↑' : '↓';
    },
    // Функция для форматирования оставшегося времени как даты
    formatRemainingTime(dateString, isPast = false) {
      return this.formatDate(dateString);
    }
  }
};
</script>

<style scoped>
.worker-dashboard {
  padding: 20px;
  margin: 0 auto;
  max-height: 100vh;
  background-color: #ffffff;

  border-left: 1px solid #cccccc;
  border-right: 1px solid #cccccc;

  box-shadow:
      2px 0 8px rgba(0, 0, 0, 0.1),
      -2px 0 8px rgba(0, 0, 0, 0.1);


  border-radius: 8px;

  display: flex;
  flex-direction: column;

  width: 90%;
}

.worker-dashboard h2,
.worker-dashboard h3 {
  text-align: center;
  color: #333333;
}

.orders-section {
  background-color: #f9f9f9;
  padding: 25px 20px;
  margin: 20px 0;
  border-radius: 10px;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);

  /* Позволяем секциям растягиваться */
  flex: 1 1 auto;

  /* Избегаем внутренней прокрутки */
  overflow: visible;
}

.orders-section h3 {
  margin-bottom: 20px;
  color: #333333;
}

.filters-container {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
  margin-bottom: 20px;
}

.filter-item {
  flex: 1;
  min-width: 200px;
  display: flex;
  flex-direction: column;
}

.filter-item label {
  margin-bottom: 5px;
  font-weight: bold;
  color: #555555;
}

.filter-input {
  padding: 8px 12px;
  border: 1px solid #dddddd;
  border-radius: 5px;
  font-size: 1em;
}

.orders-table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 15px;

  /* Делаем таблицу гибкой */
  table-layout: fixed;
}

.orders-table th,
.orders-table td {
  padding: 12px 15px;
  border: 1px solid #dddddd;
  text-align: left;

  /* Позволяем содержимому переноситься */
  word-wrap: break-word;
  word-break: break-all;

  /* Дополнительные стили для предотвращения переполнения */
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: normal;
}

.orders-table th {
  background-color: #eaeaea;
  font-weight: bold;
  color: #555555;
  position: relative;
}

.orders-table th button {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 0.9em;
  margin-left: 8px;
  color: #007BFF;
}

.orders-table th button:hover {
  text-decoration: underline;
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
  padding: 10px 20px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 1em;
  transition: background-color 0.3s ease;
  margin: 0 auto;
  display: block;
}

.btn:hover {
  background-color: #45a049;
}

.btn:disabled {
  background-color: #a5d6a7;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .filters-container {
    flex-direction: column;
  }
}

.orders-section {
  background-color: #f9f9f9;
  padding: 25px 20px;
  margin: 20px 0;
  border-radius: 10px;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
  height: 40vh;
  overflow-x: auto; /* Добавляет горизонтальную прокрутку при необходимости */
}
</style>
