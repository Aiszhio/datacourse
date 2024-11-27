<template>
  <div class="order-history">
    <h2>Заказы</h2>

    <!-- Кнопка для добавления нового заказа -->
    <div class="add-order-button">
      <button @click="openAddModal" class="btn add">Добавить заказ</button>
    </div>

    <!-- Раздел истории заказов -->
    <div class="order-list-section" v-if="!isLoading">
      <h3>Список заказов</h3>
      <table class="data-table">
        <thead>
        <tr>
          <th @click="sortBy('id')">
            Номер заказа
            <span v-if="currentSortKey === 'id'">
                {{ sortOrders.id === 'asc' ? '▲' : '▼' }}
              </span>
          </th>
          <th @click="sortBy('clientName')">
            Имя клиента
            <span v-if="currentSortKey === 'clientName'">
                {{ sortOrders.clientName === 'asc' ? '▲' : '▼' }}
              </span>
          </th>
          <th @click="sortBy('employeeName')">
            Имя сотрудника
            <span v-if="currentSortKey === 'employeeName'">
                {{ sortOrders.employeeName === 'asc' ? '▲' : '▼' }}
              </span>
          </th>
          <th @click="sortBy('service')">
            Название услуги
            <span v-if="currentSortKey === 'service'">
                {{ sortOrders.service === 'asc' ? '▲' : '▼' }}
              </span>
          </th>
          <th @click="sortBy('orderDate')">
            Дата оформления
            <span v-if="currentSortKey === 'orderDate'">
                {{ sortOrders.orderDate === 'asc' ? '▲' : '▼' }}
              </span>
          </th>
          <th @click="sortBy('receiveDate')">
            Дата получения
            <span v-if="currentSortKey === 'receiveDate'">
                {{ sortOrders.receiveDate === 'asc' ? '▲' : '▼' }}
              </span>
          </th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="order in paginatedOrders" :key="order.id">
          <td>{{ order.id }}</td>
          <td>{{ order.clientName }}</td>
          <td>{{ order.employeeName }}</td>
          <td>{{ order.service }}</td>
          <td>{{ formatDate(order.orderDate) }}</td>
          <td>{{ formatDate(order.receiveDate) }}</td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Индикатор загрузки -->
    <div v-else class="loading">
      Загрузка заказов...
    </div>

    <!-- Пагинация -->
    <div class="pagination" v-if="totalPages > 1 && !isLoading">
      <button
          @click="prevPage"
          :disabled="currentPage === 1"
          class="btn"
      >
        Назад
      </button>
      <span>Страница {{ currentPage }} из {{ totalPages }}</span>
      <button
          @click="nextPage"
          :disabled="currentPage === totalPages"
          class="btn"
      >
        Вперед
      </button>
    </div>

    <!-- Панель навигации с карточками -->
    <div class="card-panel">
      <div class="card" @click="goToAdminHome">
        <h4>На главную</h4>
        <p>Панель администратора</p>
      </div>
      <div class="card" @click="goToEmployeesPage">
        <h4>Сотрудники</h4>
        <p>Управление персоналом</p>
      </div>
      <div class="card" @click="goToMaterialsPage">
        <h4>Материалы</h4>
        <p>Управление материалами</p>
      </div>
      <div class="card" @click="goToBookingsPage">
        <h4>Бронирование</h4>
        <p>Управление бронированиями</p>
      </div>
      <div class="card" @click="goToServicesPage">
        <h4>Услуги</h4>
        <p>Управление услугами</p>
      </div>
    </div>

    <!-- Модальное окно для добавления заказа -->
    <div v-if="showModal" class="modal-overlay">
      <div class="modal">
        <h3>{{ isEditing ? modalTitle : 'Добавить заказ' }}</h3>
        <form @submit.prevent="isEditing ? submitEdit() : submitAdd()">
          <label>
            Имя клиента:
            <input v-model="currentItem.clientName" required />
          </label>
          <label>
            Имя сотрудника:
            <input v-model="currentItem.employeeName" required />
          </label>
          <label>
            Название услуги:
            <input v-model="currentItem.service" required />
          </label>
          <label>
            Дата оформления:
            <input type="date" v-model="currentItem.orderDate" required />
          </label>
          <label>
            Дата получения:
            <input type="date" v-model="currentItem.receiveDate" required />
          </label>
          <div class="modal-actions">
            <button type="submit" class="btn save">Сохранить</button>
            <button type="button" @click="closeModal" class="btn cancel">Отмена</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'; // Убедитесь, что axios настроен правильно

export default {
  name: 'OrderHistory',
  data() {
    return {
      orders: [],
      sortOrders: {
        id: 'asc',
        clientName: 'asc',
        employeeName: 'asc',
        service: 'asc',
        orderDate: 'asc',
        receiveDate: 'asc',
      },
      currentSortKey: 'id', // Сортировка по ID по умолчанию
      isLoading: true, // Состояние загрузки
      // Пагинация
      currentPage: 1,
      pageSize: 20, // Количество записей на странице
      totalPages: 1,
      // Модальное окно
      showModal: false,
      modalTitle: '',
      currentItem: {},
      isEditing: false, // Флаг для режима редактирования
    };
  },
  computed: {
    sortedOrders() {
      if (!this.currentSortKey) {
        return this.orders;
      }
      return [...this.orders].sort((a, b) => {
        let aVal = a[this.currentSortKey];
        let bVal = b[this.currentSortKey];

        // Для сортировки по ID используем числовое сравнение
        if (this.currentSortKey === 'id') {
          aVal = Number(aVal);
          bVal = Number(bVal);
        }
        // Для дат преобразуем в объекты Date
        else if (this.currentSortKey === 'orderDate' || this.currentSortKey === 'receiveDate') {
          aVal = new Date(aVal);
          bVal = new Date(bVal);
        }
        else {
          // Приводим к строкам для сравнения
          aVal = aVal.toString().toLowerCase();
          bVal = bVal.toString().toLowerCase();
        }

        if (aVal < bVal) return this.sortOrders[this.currentSortKey] === 'asc' ? -1 : 1;
        if (aVal > bVal) return this.sortOrders[this.currentSortKey] === 'asc' ? 1 : -1;
        return 0;
      });
    },
    paginatedOrders() {
      const start = (this.currentPage - 1) * this.pageSize;
      const end = start + this.pageSize;
      return this.sortedOrders.slice(start, end);
    },
  },
  methods: {
    async fetchOrders() {
      this.isLoading = true; // Начинаем загрузку
      try {
        const response = await axios.get('http://localhost:8080/api/orders/admin', {
          withCredentials: true, // Отправка куки с запросом
        });
        this.orders = response.data.orders || [];
        this.calculateTotalPages();
      } catch (error) {
        console.error('Ошибка при загрузке заказов:', error);
        alert('Не удалось загрузить заказы.');
      } finally {
        this.isLoading = false; // Завершаем загрузку
      }
    },
    calculateTotalPages() {
      this.totalPages = Math.ceil(this.sortedOrders.length / this.pageSize) || 1;
      // Если текущая страница превышает общее количество, устанавливаем на последнюю страницу
      if (this.currentPage > this.totalPages) {
        this.currentPage = this.totalPages;
      }
    },
    formatDate(dateString) {
      const options = {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
      };
      return new Date(dateString).toLocaleDateString('ru-RU', options);
    },
    sortBy(key) {
      if (this.currentSortKey === key) {
        // Переключаем направление сортировки
        this.sortOrders[key] = this.sortOrders[key] === 'asc' ? 'desc' : 'asc';
      } else {
        // Устанавливаем новое поле сортировки
        this.currentSortKey = key;
        // Устанавливаем направление сортировки по возрастанию
        this.sortOrders[key] = 'asc';
      }
      this.calculateTotalPages();
    },
    // Пагинация
    nextPage() {
      if (this.currentPage < this.totalPages) {
        this.currentPage += 1;
      }
    },
    prevPage() {
      if (this.currentPage > 1) {
        this.currentPage -= 1;
      }
    },
    goToPage(page) {
      if (page >= 1 && page <= this.totalPages) {
        this.currentPage = page;
      }
    },
    // Модальное окно
    openAddModal() {
      this.modalTitle = 'Добавить заказ';
      this.currentItem = {
        clientName: '',
        employeeName: '',
        service: '',
        orderDate: '',
        receiveDate: '',
      };
      this.isEditing = false;
      this.showModal = true;
    },
    closeModal() {
      this.showModal = false;
      this.currentItem = {};
      this.isEditing = false;
    },
    async submitAdd() {
      try {
        const newOrder = {
          clientName: this.currentItem.clientName,
          employeeName: this.currentItem.employeeName,
          service: this.currentItem.service,
          orderDate: this.currentItem.orderDate, // Формат YYYY-MM-DD
          receiveDate: this.currentItem.receiveDate, // Формат YYYY-MM-DD
        };
        const response = await axios.post('http://localhost:8080/api/orders/admin', newOrder, {
          withCredentials: true,
        });
        // Предполагаем, что сервер возвращает добавленный заказ
        this.orders.push(response.data.order);
        alert('Заказ успешно добавлен.');
        this.closeModal();
        this.calculateTotalPages();
      } catch (error) {
        console.error('Ошибка при добавлении заказа:', error);
        alert('Не удалось добавить заказ.');
      }
    },
    // Навигация
    goToAdminHome() {
      this.$router.push({ name: 'AdminHome' });
    },
    goToEmployeesPage() {
      this.$router.push({ name: 'ManageEmp' });
    },
    goToMaterialsPage() {
      this.$router.push({ name: 'MaterialsOverview' });
    },
    goToBookingsPage() {
      this.$router.push({ name: 'Bookings' });
    },
    goToServicesPage() {
      this.$router.push({ name: 'Services' });
    },
  },
  watch: {
    // Следим за изменениями в sortedOrders для пересчета общего количества страниц
    sortedOrders() {
      this.calculateTotalPages();
    },
  },
  mounted() {
    this.fetchOrders();
  },
};
</script>

<style scoped>
.order-history {
  padding: 20px;
  margin: 0 auto;
  max-width: 1200px;
  max-height: 100vh;
}

h2 {
  text-align: center;
  margin-bottom: 20px;
  color: #333;
}

/* Кнопка "Добавить заказ" */
.add-order-button {
  text-align: right;
  margin-bottom: 15px;
}

.add-order-button .btn.add {
  background-color: #4CAF50; /* Зеленый цвет */
  color: white;
  padding: 10px 20px;
  font-size: 1em;
  border-radius: 5px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  transition: background-color 0.3s ease, transform 0.2s ease;
}

.add-order-button .btn.add:hover {
  background-color: #45a049;
  transform: translateY(-2px);
}

/* Стили для таблицы заказов */
.order-list-section {
  background-color: #f9f9f9;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.1);
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 10px;
}

.data-table th,
.data-table td {
  padding: 12px;
  border: 1px solid #ddd;
  text-align: left;
  cursor: pointer; /* Указатель при наведении */
}

.data-table th {
  background-color: #f2f2f2;
  position: relative;
  user-select: none;
}

.data-table th:hover {
  background-color: #e0e0e0;
}

.data-table th span {
  margin-left: 5px;
  font-size: 0.8em;
}

/* Стили для кнопок */
.btn {
  background-color: #4CAF50; /* Зеленый цвет */
  color: white;
  padding: 8px 12px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 1em;
  transition: background-color 0.3s ease;
  margin: 5px;
}

.btn:hover {
  background-color: #45a049;
}

.btn.save {
  background-color: #4CAF50; /* Зеленый цвет */
}

.btn.save:hover {
  background-color: #45a049;
}

.btn.cancel {
  background-color: #4CAF50; /* Зеленый цвет */
}

.btn.cancel:hover {
  background-color: #45a049;
}

/* Стили для пагинации */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 15px;
  margin-top: 15px;
}

.pagination .btn {
  background-color: #4CAF50; /* Зеленый цвет */
}

.pagination .btn:hover:not(:disabled) {
  background-color: #45a049;
}

.pagination .btn:disabled {
  background-color: #a5d6a7;
  cursor: not-allowed;
}

.pagination span {
  font-size: 1em;
  color: #333;
}

/* Стили для модального окна */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal {
  background-color: white;
  padding: 25px;
  border-radius: 8px;
  width: 400px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.modal h3 {
  margin-bottom: 15px;
  color: #333;
}

.modal label {
  display: block;
  margin-bottom: 10px;
  color: #555;
}

.modal input,
.modal select {
  width: 100%;
  padding: 8px;
  margin-top: 5px;
  box-sizing: border-box;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 15px;
}

/* Стили для панели навигации с карточками */
.card-panel {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
  justify-content: center;
  margin-top: 30px;
}

.card {
  background-color: #4CAF50; /* Зеленый цвет */
  border-radius: 10px;
  padding: 15px;
  width: 180px;
  text-align: center;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  cursor: pointer;
  color: white;
  transition: transform 0.3s, box-shadow 0.3s;
}

.card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
}
</style>
