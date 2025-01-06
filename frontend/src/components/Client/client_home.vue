<template>
  <div class="user-dashboard">
    <h2>Панель пользователя</h2>
    <h3>Добро пожаловать, {{ userName }}!</h3>

    <!-- Раздел "Ваши бронирования" -->
    <div class="bookings-section">
      <h3>Ваши бронирования</h3>

      <div v-if="loadingBookings" class="loading">Загрузка ваших бронирований...</div>

      <template v-else>
        <table v-if="bookings.length" class="booking-table">
          <thead>
          <tr>
            <th>№</th>
            <th>Тип бронирования</th>
            <th>Время бронирования</th>
            <th>Действия</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="(booking, index) in bookings" :key="booking.id">
            <td>{{ index + 1 }}</td>
            <!-- Исправлено: используем правильное имя поля для типа бронирования -->
            <td>{{ booking.booking_type || 'Не указано' }}</td>
            <!-- Исправлено: используем правильное имя поля для времени бронирования -->
            <td>{{ formatDateTime(booking.booking_time) }}</td>
            <td>
              <button @click="deleteBooking(booking.id)" class="delete-button">
                Удалить
              </button>
            </td>
          </tr>
          </tbody>
        </table>

        <div v-else class="no-bookings">
          У вас нет бронирований.
        </div>
      </template>
    </div>

    <!-- Навигационные кнопки -->
    <div class="card-panel">
      <div @click="goToCreateBooking" class="card">
        <h4>Сделать бронирование</h4>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'User',
  data() {
    return {
      userName: '',
      bookings: [],
      loadingBookings: true,
    };
  },
  mounted() {
    this.initializeData();
  },
  methods: {
    async initializeData() {
      try {
        await Promise.all([
          this.fetchUserData(),
          this.fetchBookings()
        ]);
      } catch (error) {
        console.error('Ошибка инициализации данных:', error.message);
        this.$router.push({ name: 'home' });
      }
    },
    async fetchUserData() {
      try {
        const response = await fetch('http://localhost:8080/api/user', {
          method: 'GET',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
        });

        if (response.status === 401) {
          this.$router.push({ name: 'home' });
          return;
        }

        if (!response.ok) {
          throw new Error('Ошибка при получении данных пользователя');
        }

        const data = await response.json();
        this.userName = data.name;
      } catch (error) {
        console.error('Ошибка при получении данных пользователя:', error.message);
        this.$router.push({ name: 'home' });
      }
    },
    async fetchBookings() {
      try {
        const response = await fetch('http://localhost:8080/api/bookings', {
          method: 'GET',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
        });

        if (!response.ok) {
          throw new Error('Ошибка при получении бронирований');
        }

        const data = await response.json();
        // Предполагаем, что сервер возвращает объект с полем 'data', которое содержит массив бронирований
        this.bookings = data.data || [];
      } catch (error) {
        console.error('Ошибка при загрузке бронирований:', error.message);
        alert('Не удалось загрузить ваши бронирования.');
      } finally {
        this.loadingBookings = false;
      }
    },
    async deleteBooking(bookingId) {
      try {
        const response = await fetch(`http://localhost:8080/api/bookings/${bookingId}`, {
          method: 'DELETE',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
        });

        if (!response.ok) {
          throw new Error('Не удалось удалить бронирование');
        }

        // Удаление из локального массива
        this.bookings = this.bookings.filter(booking => booking.id !== bookingId);
        alert('Бронирование успешно удалено');
      } catch (error) {
        console.error('Ошибка при удалении бронирования:', error.message);
        alert('Ошибка при удалении бронирования. Попробуйте снова.');
      }
    },
    goToCreateBooking() {
      this.$router.push({ name: "ClientOrders" });
    },
    formatDateTime(dateTime) {
      const options = { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' };
      return new Date(dateTime).toLocaleString('ru-RU', options);
    }
  }
};
</script>

<style scoped>
.user-dashboard {
  padding: 20px;
  margin: 0 auto;
  max-width: 900px;
  max-height: 100vh;
}

.bookings-section {
  margin-bottom: 20px;
}

.booking-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.booking-table th,
.booking-table td {
  padding: 12px;
  border: 1px solid #ddd;
  text-align: left;
}

.booking-table th {
  background-color: #f2f2f2;
  font-weight: bold;
}

.loading,
.no-bookings {
  text-align: center;
  font-size: 1.2em;
  color: #555;
  margin-top: 20px;
}

.delete-button {
  background-color: #f44336; /* Красный цвет */
  color: white;
  border: none;
  padding: 8px 12px;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.delete-button:hover {
  background-color: #d32f2f; /* Темно-красный при наведении */
}

.card-panel {
  display: flex;
  gap: 20px;
  justify-content: center;
  margin-top: 20px;
}

.card {
  background-color: #4CAF50;
  color: white;
  border-radius: 10px;
  padding: 15px;
  width: 380px;
  text-align: center;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  cursor: pointer;
  transition: transform 0.3s, box-shadow 0.3s;
}

.card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
}
</style>
