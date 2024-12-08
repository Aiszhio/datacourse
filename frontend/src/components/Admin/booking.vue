<template>
  <div class="bookings-dashboard">
    <h2>Управление бронированием</h2>

    <!-- Раздел фильтрации бронирований -->
    <div class="filter-section">
      <h3>Фильтровать бронирования</h3>
      <div class="filters">
        <label>
          Тип брони:
          <select v-model="filters.type">
            <option value="">Все типы</option>
            <option>Онлайн</option>
            <option>Очный</option>
          </select>
        </label>
        <label>
          ФИО бронирующего:
          <input
              type="text"
              v-model="filters.name"
              placeholder="Поиск по ФИО"
          />
        </label>
        <label>
          Дата брони:
          <input
              type="date"
              v-model="filters.date"
              :max="todayDate"
          />
        </label>
      </div>
    </div>

    <!-- Таблица с бронированиями -->
    <div class="table-section">
      <h3>Брони</h3>
      <button @click="openCreateBookingModal" class="btn">Создать бронь</button>
      <table class="data-table">
        <thead>
        <tr>
          <th @click="sortBy('type')">
            Тип брони
            <span v-if="currentSortKey === 'type'">
                {{ sortOrders.type === 'asc' ? '▲' : '▼' }}
              </span>
          </th>
          <th @click="sortBy('time')">
            Время брони
            <span v-if="currentSortKey === 'time'">
                {{ sortOrders.time === 'asc' ? '▲' : '▼' }}
              </span>
          </th>
          <th @click="sortBy('name')">
            ФИО бронирующего
            <span v-if="currentSortKey === 'name'">
                {{ sortOrders.name === 'asc' ? '▲' : '▼' }}
              </span>
          </th>
          <th>
            Действия
          </th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="booking in paginatedBookings" :key="booking.id">
          <td>{{ booking.type }}</td>
          <td>{{ formatDate(booking.time) }}</td>
          <td>{{ booking.name }}</td>
          <td>
            <button
                v-if="isUpcoming(booking.time)"
                @click="deleteBooking(booking.id)"
                class="btn delete"
            >
              Удалить
            </button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Элементы управления пагинацией -->
    <div class="pagination" v-if="totalPages > 1">
      <button
          @click="prevPage"
          :disabled="currentPage === 1"
          class="btn pagination-btn"
      >
        Назад
      </button>
      <span>Страница {{ currentPage }} из {{ totalPages }}</span>
      <button
          @click="nextPage"
          :disabled="currentPage === totalPages"
          class="btn pagination-btn"
      >
        Вперёд
      </button>
    </div>

    <!-- Панель навигации с карточками -->
    <div class="card-panel">
      <div class="card" @click="goToAdminHome">
        <h4>Главная</h4>
        <p>Панель администратора</p>
      </div>
      <div class="card" @click="goToEmployeesPage">
        <h4>Сотрудники</h4>
        <p>Управление персоналом</p>
      </div>
      <div class="card" @click="goToOrdersPage">
        <h4>Заказы</h4>
        <p>История заказов</p>
      </div>
      <div class="card" @click="goToMaterialsPage">
        <h4>Закупки</h4>
        <p>Материалы и расходы</p>
      </div>
      <div class="card" @click="goToServicesPage">
        <h4>Услуги</h4>
        <p>Управление услугами</p>
      </div>
      <div class="card" @click="goToClients">
        <i class="fas fa-box icon"></i>
        <h3>Клиенты</h3>
        <p>Учёт и управление</p>
      </div>
    </div>

    <!-- Модальное окно для создания бронирования -->
    <div v-if="showCreateModal" class="modal-overlay">
      <div class="modal">
        <h3>Создать бронь</h3>
        <form @submit.prevent="createBooking">
          <label>
            Тип брони:
            <select v-model="newBooking.type" required>
              <option disabled value="">Выберите тип брони</option>
              <option>Онлайн</option>
              <option>Очный</option>
            </select>
          </label>
          <label>
            Время брони:
            <input
                type="datetime-local"
                v-model="newBooking.time"
                :min="minBookingTime"
                required
            />
          </label>
          <label>
            ФИО бронирующего:
            <input
                type="text"
                v-model="newBooking.name"
                required
                maxlength="80"
                placeholder="Введите ФИО"
            />
          </label>
          <div class="modal-actions">
            <button type="submit" class="btn save">Создать</button>
            <button type="button" @click="closeCreateModal" class="btn cancel">Отмена</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'; // Убедитесь, что axios настроен правильно

export default {
  name: 'Bookings',
  data() {
    return {
      bookings: [], // Массив для загрузки бронирований
      showCreateModal: false,
      newBooking: {
        type: '',
        time: '',
        name: ''
      },
      sortOrders: {
        type: 'asc',
        time: 'asc',
        name: 'asc'
      },
      currentSortKey: '',
      // Пагинация
      currentPage: 1,
      pageSize: 10,
      totalPages: 1,
      // Фильтры
      filters: {
        type: '',
        name: '',
        date: ''
      },
      todayDate: this.getTodayDate()
    };
  },
  computed: {
    filteredBookings() {
      return this.bookings.filter(booking => {
        const matchesType = this.filters.type
            ? booking.type === this.filters.type
            : true;
        const matchesName = this.filters.name
            ? booking.name.toLowerCase().includes(this.filters.name.toLowerCase())
            : true;
        const matchesDate = this.filters.date
            ? new Date(booking.time).toISOString().split('T')[0] === this.filters.date
            : true;
        return matchesType && matchesName && matchesDate;
      });
    },
    sortedBookings() {
      if (!this.currentSortKey) {
        return this.filteredBookings;
      }
      return [...this.filteredBookings].sort((a, b) => {
        let aVal = a[this.currentSortKey];
        let bVal = b[this.currentSortKey];

        // Для даты преобразуем в объекты Date
        if (this.currentSortKey === 'time') {
          aVal = new Date(aVal);
          bVal = new Date(bVal);
        } else if (typeof aVal === 'string') {
          aVal = aVal.toLowerCase();
          bVal = bVal.toLowerCase();
        }

        if (aVal < bVal) return this.sortOrders[this.currentSortKey] === 'asc' ? -1 : 1;
        if (aVal > bVal) return this.sortOrders[this.currentSortKey] === 'asc' ? 1 : -1;
        return 0;
      });
    },
    paginatedBookings() {
      const start = (this.currentPage - 1) * this.pageSize;
      const end = start + this.pageSize;
      return this.sortedBookings.slice(start, end);
    }
  },
  methods: {
    async fetchBookings() {
      try {
        const response = await axios.get('http://localhost:8080/api/bookings/admin', { withCredentials: true });
        // Предполагаем, что API возвращает объект с массивом bookings
        this.bookings = response.data.bookings;
        this.calculateTotalPages();
      } catch (error) {
        console.error('Ошибка при загрузке бронирований:', error);
        alert('Не удалось загрузить бронирования.');
      }
    },
    formatDate(dateString) {
      const options = {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      };
      return new Date(dateString).toLocaleDateString('ru-RU', options);
    },
    isUpcoming(bookingTime) {
      return new Date(bookingTime) > new Date();
    },
    async deleteBooking(bookingId) {
      if (!confirm(`Вы уверены, что хотите удалить бронь #${bookingId}?`)) {
        return;
      }
      try {
        await axios.delete(`http://localhost:8080/api/bookings/admin/${bookingId}`, { withCredentials: true });
        // Удаляем бронирование из локального массива
        this.bookings = this.bookings.filter(b => b.id !== bookingId);
        alert('Бронирование успешно удалено.');
        this.calculateTotalPages();
        // Если текущая страница пуста после удаления, перейдём на предыдущую
        if (this.paginatedBookings.length === 0 && this.currentPage > 1) {
          this.currentPage--;
        }
      } catch (error) {
        console.error('Ошибка при удалении бронирования:', error);
        alert('Не удалось удалить бронирование.');
      }
    },
    goToAdminHome() {
      this.$router.push({ name: 'AdminHome' });
    },
    goToEmployeesPage() {
      this.$router.push({ name: 'ManageEmp' });
    },
    goToOrdersPage() {
      this.$router.push({ name: 'OrderHistory' });
    },
    goToMaterialsPage() {
      this.$router.push({ name: 'MaterialsOverview' });
    },
    goToServicesPage() {
      this.$router.push({ name: 'Services' });
    },
    goToClients() {
      this.$router.push({ name: 'Clients' });
    },
    sortBy(key) {
      if (this.currentSortKey === key) {
        // Переключаем направление сортировки
        this.sortOrders[key] = this.sortOrders[key] === 'asc' ? 'desc' : 'asc';
      } else {
        // Устанавливаем новый ключ сортировки и направление по возрастанию
        this.currentSortKey = key;
        this.sortOrders[key] = 'asc';
      }
      this.currentPage = 1; // Сбросить на первую страницу при изменении сортировки
    },
    // Пагинация
    calculateTotalPages() {
      this.totalPages = Math.ceil(this.sortedBookings.length / this.pageSize) || 1;
      if (this.currentPage > this.totalPages) {
        this.currentPage = this.totalPages;
      }
    },
    nextPage() {
      if (this.currentPage < this.totalPages) {
        this.currentPage++;
      }
    },
    prevPage() {
      if (this.currentPage > 1) {
        this.currentPage--;
      }
    },
    // Создание бронирования
    openCreateBookingModal() {
      this.showCreateModal = true;
      this.newBooking = {
        type: '',
        time: '',
        name: ''
      };
    },
    closeCreateModal() {
      this.showCreateModal = false;
      this.newBooking = {
        type: '',
        time: '',
        name: ''
      };
    },
    async createBooking() {
      // Проверка, что время брони не в прошлом
      const selectedTime = new Date(this.newBooking.time);
      const now = new Date();
      if (selectedTime < now) {
        alert('Время брони не может быть в прошлом.');
        return;
      }

      // Проверка длины ФИО
      if (this.newBooking.name.length > 80) {
        alert('ФИО бронирующего не может превышать 80 символов.');
        return;
      }

      try {
        const newBookingData = {
          booking_type: this.newBooking.type, // Изменено
          booking_time: new Date(this.newBooking.time).toISOString(), // Изменено
          booker_full_name: this.newBooking.name // Изменено
        };
        const response = await axios.post('http://localhost:8080/api/createOrder/admin', newBookingData, { withCredentials: true });

        // Проверяем, есть ли данные бронирования в ответе
        if (response.data.booking) {
          this.bookings.push(response.data.booking);
          this.bookings = this.sortedBookings; // Обновляем сортировку после добавления
        }

        alert('Бронирование успешно создано.');
        this.closeCreateModal();
        this.calculateTotalPages();
      } catch (error) {
        console.error('Ошибка при создании бронирования:', error);
        if (error.response && error.response.data && error.response.data.error) {
          alert(`Не удалось создать бронирование: ${error.response.data.error}`);
        } else {
          alert('Не удалось создать бронирование.');
        }
      }
    },
    getTodayDate() {
      const date = new Date();
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      return `${year}-${month}-${day}`;
    }
  },
  watch: {
    // Обновляем пагинацию при изменении фильтров или бронирований
    filters: {
      handler() {
        this.currentPage = 1;
        this.calculateTotalPages();
      },
      deep: true
    },
    bookings() {
      this.calculateTotalPages();
    }
  },
  mounted() {
    this.fetchBookings();
  }
};
</script>

<style scoped>
.bookings-dashboard {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
  max-height: 100vh;
}

h2 {
  text-align: center;
  width: 100%;
  margin-bottom: 20px;
  color: #333;
}

.top-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.sort-buttons {
  display: flex;
  gap: 10px;
}

.sort-buttons .btn.sort {
  background-color: #4CAF50; /* Зеленый цвет для кнопок сортировки */
  color: white;
}

.sort-buttons .btn.sort:hover {
  background-color: #45a049;
}

.create-booking-button .btn.create {
  background-color: #4CAF50; /* Зеленый цвет для кнопки создания */
  color: white;
}

.create-booking-button .btn.create:hover {
  background-color: #45a049;
}

.table-section {
  margin-bottom: 20px;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 10px;
  border: 1px solid #ddd;
  text-align: left;
}

.data-table th {
  background-color: #f2f2f2;
}

.btn {
  background-color: #4CAF50;
  color: white;
  padding: 8px 12px;
  margin-bottom: 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn.edit {
  background-color: #4CAF50; /* Зеленый цвет для кнопки редактирования */
  margin-right: 5px;
}

.btn.edit:hover {
  background-color: #45a049;
}

.btn.delete {
  background-color: #F44336; /* Красный цвет для кнопки удаления */
}

.btn.delete:hover {
  background-color: #D32F2F;
}

.btn.pagination-btn {
  background-color: #4CAF50; /* Зеленый цвет для кнопок пагинации */
  color: white;
}

.btn.pagination-btn:hover:not(:disabled) {
  background-color: #45a049;
}

.btn.pagination-btn:disabled {
  background-color: #a5d6a7;
  cursor: not-allowed;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 15px;
  margin-top: 15px;
}

.card-panel {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
  justify-content: center;
  margin-top: 30px;
}

.card {
  background-color: #4CAF50; /* Зеленый цвет */
  color: #fff;
  border-radius: 10px;
  padding: 15px;
  width: 180px;
  text-align: center;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  cursor: pointer;
  transition: transform 0.3s, box-shadow 0.3s;
}

.card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
}

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

.modal select,
.modal input {
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

.btn.save {
  background-color: #4CAF50; /* Зеленый цвет для кнопки сохранения */
  color: white;
}

.btn.save:hover {
  background-color: #45a049;
}

.btn.cancel {
  background-color: #4CAF50; /* Зеленый цвет для кнопки отмены */
  color: white;
}

.btn.cancel:hover {
  background-color: #45a049;
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

/* Стили для меток (label) */
label {
  font-size: 16px; /* Увеличим размер шрифта для label */
  margin-bottom: 5px;
  display: block;
  font-weight: bold;
  color: #333;
}

/* Стили для блоков фильтрации */
.filter-wrapper {
  display: flex;
  gap: 15px;
  margin-bottom: 20px;
}

.filter-wrapper .filter-input {
  width: 100%;
  padding: 12px;
  margin: 10px 0;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 16px;
  transition: border-color 0.3s ease-in-out;
}

.filter-wrapper .filter-input:focus {
  border-color: #4CAF50;
  outline: none;
}

.filter-wrapper .filter-input:hover {
  border-color: #aaa;
}

.filter-wrapper .filter-select {
  width: 100%;
  padding: 12px;
  margin: 10px 0;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 16px;
  transition: border-color 0.3s ease-in-out;
}

.filter-wrapper .filter-select:focus {
  border-color: #4CAF50;
  outline: none;
}

.filter-wrapper .filter-select:hover {
  border-color: #aaa;
}

/* Стили для кнопок фильтрации */
.filter-wrapper .btn.filter {
  background-color: #4CAF50; /* Зеленый цвет для кнопок фильтрации */
  color: white;
  padding: 8px 12px;
  border-radius: 4px;
  border: none;
  cursor: pointer;
  transition: background-color 0.3s;
}

.filter-wrapper .btn.filter:hover {
  background-color: #45a049;
}

.filter-wrapper .btn.filter:focus {
  outline: none;
  box-shadow: 0 0 5px rgba(76, 175, 80, 0.3);
}

input[type="text"],
input[type="datetime-local"],
select,
input[type="date"] {
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
select:focus,
input[type="date"]:focus {
  border-color: #007bff;
  outline: none;
}

input[type="text"]:hover,
input[type="datetime-local"]:hover,
select:hover,
input[type="date"]:hover {
  border-color: #aaa;
}

/* Стили для меток (label) */
label {
  font-size: 16px; /* Увеличим размер шрифта для label */
  margin-bottom: 5px;
  display: block;
  font-weight: bold;
  color: #333;
}

/* Специальные стили для поля "date" */
input[type="date"] {
  -webkit-appearance: none;
  -moz-appearance: none;
  appearance: none;
  background-color: #fff;
}

input[type="date"]:focus {
  border-color: #4CAF50; /* Зеленый цвет для фокуса */
}

input[type="date"]:hover {
  border-color: #aaa;
}
</style>
