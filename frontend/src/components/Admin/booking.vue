<template>
  <div class="bookings-dashboard">
    <h2>Управление бронированием</h2>

    <!-- Раздел сортировки бронирований -->
    <div class="sort-buttons">
      <button @click="sortBy('id')" class="btn">
        Сортировать по номеру брони
        <span v-if="currentSortKey === 'id'">
          {{ sortOrders.id === 'asc' ? '▲' : '▼' }}
        </span>
      </button>
      <button @click="sortBy('type')" class="btn">
        Сортировать по типу брони
        <span v-if="currentSortKey === 'type'">
          {{ sortOrders.type === 'asc' ? '▲' : '▼' }}
        </span>
      </button>
      <button @click="sortBy('orderId')" class="btn">
        Сортировать по номеру заказа
        <span v-if="currentSortKey === 'orderId'">
          {{ sortOrders.orderId === 'asc' ? '▲' : '▼' }}
        </span>
      </button>
      <button @click="sortBy('time')" class="btn">
        Сортировать по времени брони
        <span v-if="currentSortKey === 'time'">
          {{ sortOrders.time === 'asc' ? '▲' : '▼' }}
        </span>
      </button>
      <button @click="sortBy('name')" class="btn">
        Сортировать по ФИО бронирующего
        <span v-if="currentSortKey === 'name'">
          {{ sortOrders.name === 'asc' ? '▲' : '▼' }}
        </span>
      </button>
    </div>

    <!-- Таблица с бронированиями -->
    <div class="table-section">
      <h3>Брони</h3>
      <table class="data-table">
        <thead>
        <tr>
          <th>Номер брони</th>
          <th>Тип брони</th>
          <th>Номер заказа</th>
          <th>Время брони</th>
          <th>ФИО бронирующего</th>
          <th>Действия</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="booking in sortedBookings" :key="booking.id">
          <td>{{ booking.id }}</td>
          <td>{{ booking.type }}</td>
          <td>{{ booking.orderId }}</td>
          <td>{{ formatDate(booking.time) }}</td>
          <td>{{ booking.name }}</td>
          <td>
            <button
                v-if="isUpcoming(booking.time)"
                @click="editBooking(booking)"
                class="btn edit"
            >
              Редактировать
            </button>
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
    </div>

    <!-- Модальное окно для редактирования бронирования -->
    <div v-if="showModal" class="modal-overlay">
      <div class="modal">
        <h3>{{ modalTitle }}</h3>
        <form @submit.prevent="submitEdit">
          <label>
            Тип брони:
            <input v-model="currentItem.type" required />
          </label>
          <label>
            Номер заказа:
            <input type="number" v-model="currentItem.orderId" required />
          </label>
          <label>
            Время брони:
            <input type="datetime-local" v-model="formattedBookingTime" required />
          </label>
          <label>
            ФИО бронирующего:
            <input v-model="currentItem.name" required />
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
import axios from 'axios'; // Замените на import из настроенного файла axios.js, если используется

export default {
  name: 'Bookings',
  data() {
    return {
      bookings: [], // Пустой массив для загрузки бронирований
      showModal: false,
      modalTitle: '',
      currentItem: {},
      formattedBookingTime: '',
      sortOrders: {
        id: 'asc',
        type: 'asc',
        orderId: 'asc',
        time: 'asc',
        name: 'asc',
      },
      currentSortKey: '',
    };
  },
  computed: {
    sortedBookings() {
      if (!this.currentSortKey) {
        return this.bookings;
      }
      return [...this.bookings].sort((a, b) => {
        let aVal = a[this.currentSortKey];
        let bVal = b[this.currentSortKey];

        // Для даты преобразуем в объекты Date
        if (this.currentSortKey === 'time') {
          aVal = new Date(aVal);
          bVal = new Date(bVal);
        } else {
          // Приводим к строкам для сравнения
          aVal = aVal.toString().toLowerCase();
          bVal = bVal.toString().toLowerCase();
        }

        if (aVal < bVal) return this.sortOrders[this.currentSortKey] === 'asc' ? -1 : 1;
        if (aVal > bVal) return this.sortOrders[this.currentSortKey] === 'asc' ? 1 : -1;
        return 0;
      });
    },
  },
  methods: {
    async fetchBookings() {
      try {
        const response = await axios.get('http://localhost:8080/api/bookings/admin', { withCredentials: true });
        // Предполагаем, что API возвращает объект с массивом bookings
        this.bookings = response.data.bookings;
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
    editBooking(booking) {
      this.modalTitle = `Редактирование брони #${booking.id}`;
      this.currentItem = { ...booking };
      // Форматируем время для input[type="datetime-local"]
      this.formattedBookingTime = this.formatDateForInput(booking.time);
      this.showModal = true;
    },
    formatDateForInput(dateString) {
      const date = new Date(dateString);
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      return `${year}-${month}-${day}T${hours}:${minutes}`;
    },
    closeModal() {
      this.showModal = false;
      this.currentItem = {};
      this.formattedBookingTime = '';
    },
    async submitEdit() {
      try {
        const updatedBooking = {
          type: this.currentItem.type,
          orderId: this.currentItem.orderId,
          time: new Date(this.formattedBookingTime).toISOString(),
          name: this.currentItem.name
        };
        await axios.put(`/api/bookings/${this.currentItem.id}`, updatedBooking);
        // Обновляем локальные данные
        const index = this.bookings.findIndex(b => b.id === this.currentItem.id);
        if (index !== -1) {
          this.bookings[index] = { ...this.currentItem, time: updatedBooking.time };
        }
        alert('Бронирование успешно обновлено.');
        this.closeModal();
      } catch (error) {
        console.error('Ошибка при обновлении бронирования:', error);
        alert('Не удалось обновить бронирование.');
      }
    },
    async deleteBooking(bookingId) {
      if (!confirm(`Вы уверены, что хотите удалить бронь #${bookingId}?`)) {
        return;
      }
      try {
        await axios.delete(`/api/bookings/${bookingId}`);
        // Удаляем бронирование из локального массива
        this.bookings = this.bookings.filter(b => b.id !== bookingId);
        alert('Бронирование успешно удалено.');
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
    sortBy(key) {
      if (this.currentSortKey === key) {
        // Переключаем направление сортировки
        this.sortOrders[key] = this.sortOrders[key] === 'asc' ? 'desc' : 'asc';
      } else {
        // Устанавливаем новый ключ сортировки и направление по умолчанию
        this.currentSortKey = key;
      }
    },
  },
  mounted() {
    this.fetchBookings();
  }
};
</script>

<style scoped>
.bookings-dashboard {
  padding: 20px;
  max-width: 900px;
  max-height: 100vh;
  margin: 0 auto;
}

.card-panel {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
  justify-content: center;
  margin-bottom: 30px;
}

.card {
  background-color: #4CAF50;
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

.table-section {
  margin-bottom: 20px;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th, .data-table td {
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
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  margin: 5px;
}

.btn:hover {
  background-color: #45a049;
}

.btn.danger {
  background-color: #f44336;
}

.btn.danger:hover {
  background-color: #d32f2f;
}

h2 {
  text-align: center;
  width: 100%;
  margin-bottom: 20px;
  color: #333;
}
</style>
