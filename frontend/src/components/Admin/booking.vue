<template>
  <div class="bookings-dashboard">
    <h2>Управление бронированием</h2>

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
          <th>Действия</th> <!-- Добавляем столбец Действия -->
        </tr>
        </thead>
        <tbody>
        <tr v-for="booking in bookings" :key="booking.id">
          <td>{{ booking.id }}</td>
          <td>{{ booking.type }}</td>
          <td>{{ booking.orderId }}</td>
          <td>{{ formatDate(booking.time) }}</td>
          <td>{{ booking.name }}</td>
          <td>
            <!-- Проверяем, является ли бронирование предстоящим -->
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
      formattedBookingTime: ''
    };
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
