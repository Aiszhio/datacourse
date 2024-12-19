<template>
  <div class="order-history">
    <h2>Заказы</h2>

    <!-- Кнопка для добавления нового заказа -->
    <div class="add-order-button">
      <button @click="openAddModal" class="btn add">Добавить заказ</button>
    </div>

    <!-- Раздел фильтрации заказов -->
    <div class="filter-section">
      <h3>Фильтровать заказы</h3>
      <div class="filters">
        <label>
          Имя клиента:
          <input v-model="searchClientName" class="clientNameInput" placeholder="Поиск по имени клиента" />
        </label>
        <label>
          Название услуги:
          <select v-model="searchService">
            <option value="">Все услуги</option>
            <option v-for="service in uniqueServices" :key="service.name" :value="service.name">
              {{ service.name }}
            </option>
          </select>
        </label>
      </div>
    </div>

    <!-- Раздел истории заказов -->
    <div class="order-list-section" v-if="!isLoading">
      <h3>Список заказов</h3>
      <table class="data-table">
        <thead>
        <tr>
          <th @click="sortBy('clientName')">
            Имя клиента
            <span v-if="currentSortKey === 'clientName'">
              {{ sortOrders.clientName === 'asc' ? '▲' : '▼' }}
            </span>
          </th>
          <th @click="sortBy('employeeName')">
            Сотрудник
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
            Дата завершения
            <span v-if="currentSortKey === 'receiveDate'">
              {{ sortOrders.receiveDate === 'asc' ? '▲' : '▼' }}
            </span>
          </th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="order in paginatedOrders" :key="order.id">
          <td>{{ order.clientName }}</td>
          <td>{{ order.employeeName }}</td>
          <td>{{ order.service }}</td>
          <td>{{ formatDateTime(order.orderDate) }}</td>
          <td>{{ formatDateTime(order.receiveDate) }}</td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Индикатор загрузки -->
    <div v-else>
      Загрузка заказов...
    </div>

    <!-- Пагинация -->
    <div v-if="totalPages > 1 && !isLoading">
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
        <p>Услуги и оборудование</p>
      </div>
      <div class="card" @click="goToClients">
        <i class="fas fa-box icon"></i>
        <h3>Клиенты</h3>
        <p>Учёт и управление</p>
      </div>
    </div>

    <!-- Модальное окно для добавления заказа -->
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <h3>{{ isEditing ? modalTitle : 'Добавить заказ' }}</h3>
        <form @submit.prevent="isEditing ? submitEdit() : submitAdd()">
          <label>
            Номер телефона:
            <input
                v-model="currentItem.phoneNumber"
                required
            />
          </label>

          <label>
            Название услуги:
            <input v-model="currentItem.service" placeholder="Введите название услуги" maxlength="40" required />
          </label>

          <label>
            Дата оформления:
            <input
                type="datetime-local"
                v-model="currentItem.orderDate"
                :min="getCurrentDateTime()"
                required
            />
          </label>

          <div>
            <button type="submit" class="btn save">Сохранить</button>
            <button type="button" @click="closeModal" class="btn cancel">Отмена</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'OrderHistory',
  data() {
    return {
      orders: [],
      uniqueServices: [],
      sortOrders: {
        clientName: 'asc',
        employeeName: 'asc',
        service: 'asc',
        orderDate: 'asc',
        receiveDate: 'asc',
      },
      currentSortKey: 'clientName',
      isLoading: true,
      currentPage: 1,
      pageSize: 20,
      totalPages: 1,
      showModal: false,
      modalTitle: '',
      currentItem: {},
      isEditing: false,
      searchClientName: '',
      searchService: '',
    };
  },
  computed: {
    filteredOrders() {
      return this.orders.filter(order => {
        if (!order.clientName || !order.service) {
          return false;
        }

        const matchesClientName = order.clientName
            .toLowerCase()
            .includes(this.searchClientName.toLowerCase());
        const matchesService = this.searchService
            ? order.service === this.searchService
            : true;
        return matchesClientName && matchesService;
      });
    },
    sortedOrders() {
      if (!this.currentSortKey) {
        return this.filteredOrders;
      }
      return [...this.filteredOrders].sort((a, b) => {
        let aVal = a[this.currentSortKey] ?? '';
        let bVal = b[this.currentSortKey] ?? '';

        if (this.currentSortKey === 'orderDate' || this.currentSortKey === 'receiveDate') {
          aVal = new Date(aVal);
          bVal = new Date(bVal);
        } else {
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
    toMoscowTimeISOString(date) {
      const options = {
        timeZone: 'Europe/Moscow',
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false,
      };
      const formatter = new Intl.DateTimeFormat('sv-SE', options);
      const parts = formatter.formatToParts(date);
      const dateParts = {};

      parts.forEach(({ type, value }) => {
        dateParts[type] = value;
      });

      return `${dateParts.year}-${dateParts.month}-${dateParts.day}T${dateParts.hour}:${dateParts.minute}:${dateParts.second}+03:00`;
    },

    formatDateTime(dateString) {
      const date = new Date(dateString);
      const moscowTime = date.toLocaleString('ru-RU', {
        timeZone: 'Europe/Moscow',
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false
      });

      return moscowTime;
    },

    getCurrentDateTime() {
      const now = new Date();
      const year = now.getFullYear();
      const month = String(now.getMonth() + 1).padStart(2, '0');
      const day = String(now.getDate()).padStart(2, '0');
      const hours = String(now.getHours()).padStart(2, '0');
      const minutes = String(now.getMinutes()).padStart(2, '0');
      return `${year}-${month}-${day}T${hours}:${minutes}:00`;
    },

    openAddModal() {
      this.modalTitle = 'Добавить заказ';
      this.currentItem = {
        phoneNumber: '',
        service: '',
        orderDate: this.getCurrentDateTime(),
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
        const moscowOrderDate = this.toMoscowTimeISOString(new Date(this.currentItem.orderDate));
        const newOrder = {
          phoneNumber: this.currentItem.phoneNumber,
          service: this.currentItem.service,
          orderDate: moscowOrderDate,
        };

        const response = await axios.post('http://localhost:8080/api/orders/admin', newOrder, {
          withCredentials: true,
        });

        if (response.data.message === 'Заказ успешно сохранен') {
          this.orders.push(response.data.order);
          alert('Заказ успешно добавлен.');
          this.closeModal();
          this.calculateTotalPages();
        } else {
          alert('Неизвестная ошибка: ' + (response.data.message || 'Попробуйте позже.'));
        }

      } catch (error) {
        console.error('Ошибка при добавлении заказа:', error);

        if (error.response) {
          alert('Ошибка сервера: ' + (error.response.data.error || 'Неизвестная ошибка'));
        } else if (error.request) {
          alert('Ошибка соединения с сервером. Пожалуйста, проверьте ваше интернет-соединение.');
        } else {
          alert('Неизвестная ошибка: ' + error.message);
        }
      }
    },

    sortBy(key) {
      if (this.currentSortKey === key) {
        this.sortOrders[key] = this.sortOrders[key] === 'asc' ? 'desc' : 'asc';
      } else {
        this.currentSortKey = key;
        this.sortOrders[key] = 'asc';
      }
      this.calculateTotalPages();
    },

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

    openEditModal(order) {
      this.modalTitle = 'Редактировать заказ';
      this.currentItem = { ...order };
      this.isEditing = true;
      this.showModal = true;
    },

    async submitEdit() {
      try {
        const moscowOrderDate = this.toMoscowTimeISOString(new Date(this.currentItem.orderDate));
        const updatedOrder = {
          clientName: this.currentItem.clientName,
          service: this.currentItem.service,
          orderDate: moscowOrderDate,
        };

        const response = await axios.put(`http://localhost:8080/api/orders/${this.currentItem.id}`, updatedOrder, {
          withCredentials: true,
        });

        if (response.data.message === 'Заказ успешно обновлён') {
          const index = this.orders.findIndex(order => order.id === this.currentItem.id);
          if (index !== -1) {
            this.$set(this.orders, index, response.data.order);
          }
          alert('Заказ успешно обновлён.');
          this.closeModal();
          this.calculateTotalPages();
        } else {
          alert('Неизвестная ошибка: ' + (response.data.message || 'Попробуйте позже.'));
        }

      } catch (error) {
        console.error('Ошибка при редактировании заказа:', error);

        if (error.response) {
          alert('Ошибка сервера: ' + (error.response.data.error || 'Неизвестная ошибка'));
        } else if (error.request) {
          alert('Ошибка соединения с сервером. Пожалуйста, проверьте ваше интернет-соединение.');
        } else {
          alert('Неизвестная ошибка: ' + error.message);
        }
      }
    },

    async fetchOrders() {
      this.isLoading = true;
      try {
        const response = await axios.get('http://localhost:8080/api/orders/admin', {
          withCredentials: true,
        });
        this.orders = response.data.orders || [];
        this.extractUniqueServices();
        this.calculateTotalPages();
      } catch (error) {
        console.error('Ошибка при загрузке заказов:', error);
        alert('Не удалось загрузить заказы.');
      } finally {
        this.isLoading = false;
      }
    },

    extractUniqueServices() {
      const servicesSet = new Set();
      this.orders.forEach(order => {
        servicesSet.add(order.service);
      });
      this.uniqueServices = Array.from(servicesSet).map(serviceName => ({
        name: serviceName
      }));
    },

    calculateTotalPages() {
      this.totalPages = Math.ceil(this.sortedOrders.length / this.pageSize) || 1;
      if (this.currentPage > this.totalPages) {
        this.currentPage = this.totalPages;
      }
    },

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
    goToClients() {
      this.$router.push({ name: 'Clients' });
    },
  },
  watch: {
    'currentItem.phoneNumber'(val) {
      // Если значение пустое или undefined, просто выходим
      if (!val) return;

      let digits = val.replace(/\D/g, '');
      if (!digits.startsWith('7')) {
        digits = '7' + digits.replace(/^7+/, '');
      }

      if (digits.length > 1) {
        digits = digits[0] + '(' + digits.slice(1);
      }
      if (digits.length > 4) {
        digits = digits.slice(0,5) + ')-' + digits.slice(5);
      }
      if (digits.length > 9) {
        digits = digits.slice(0,10) + '-' + digits.slice(10);
      }
      if (digits.length > 12) {
        digits = digits.slice(0,13) + '-' + digits.slice(13);
      }
      if (digits.length > 15) {
        digits = digits.slice(0,16);
      }

      this.currentItem.phoneNumber = digits;
    },

    searchClientName() {
      this.currentPage = 1;
      this.calculateTotalPages();
    },
    searchService() {
      this.currentPage = 1;
      this.calculateTotalPages();
    },
    orders() {
      this.calculateTotalPages();
    }
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
  background-color: red; /* Зеленый цвет */
}

.btn.cancel:hover {
  background-color: darkred;
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

.clientNameInput{
  width: 98%;
  padding: 12px; /* Увеличим padding для большего комфорта */
  margin: 10px 0;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 16px; /* Увеличим размер шрифта */
  transition: border-color 0.3s ease-in-out;
}

input[type="text"],
input[type="datetime-local"],
select {
  width: 100%;
  padding: 12px; /* Увеличим padding для большего комфорта */
  margin: 10px 0;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 16px; /* Увеличим размер шрифта */
  transition: border-color 0.3s ease-in-out;
}

input[type="text"]:focus,
input[type="datetime-local"]:focus,
select:focus {
  border-color: #007bff;
  outline: none;
}

input[type="text"]:hover,
input[type="datetime-local"]:hover,
select:hover {
  border-color: #aaa;
}

label {
  font-size: 16px; /* Увеличим размер шрифта для label */
  margin-bottom: 5px;
  display: block;
  font-weight: bold;
  color: #333;
}

/* Стили для модальных окон */
.modal input[type="text"],
.modal input[type="datetime-local"],
.modal select {
  margin: 10px 0;
  padding: 12px;
  font-size: 16px; /* Увеличиваем размер шрифта */
  border-radius: 4px;
  width: 100%;
  box-sizing: border-box;
  border: 1px solid #ddd;
}

select[name="service"] {
  max-width: 200px; /* Ограничим максимальную ширину поля */
  margin: 10px 0;
}

.modal input[type="text"]:focus,
.modal input[type="datetime-local"]:focus,
.modal select:focus {
  border-color: #0056b3;
  outline: none;
}
</style>