<script>
export default {
  name: 'Bookings',
  data() {
    return {
      bookings: [
        { id: 1, type: 'Персональная', orderId: 101, time: '2024-11-15 14:00', name: 'Иван Иванов', status: 'Ожидает подтверждения' },
        { id: 2, type: 'Групповая', orderId: 102, time: '2024-11-16 10:00', name: 'Анна Петрова', status: 'Ожидает подтверждения' }
      ], // Пример списка бронирований
      showModal: false,
      modalType: '',
      modalTitle: '',
      currentItem: {}
    };
  },
  methods: {
    confirmBooking(bookingId) {
      const booking = this.bookings.find(b => b.id === bookingId);
      if (booking) {
        booking.status = 'Подтверждена';
        alert(`Бронь #${bookingId} подтверждена`);
      }
    },
    cancelBooking(bookingId) {
      const booking = this.bookings.find(b => b.id === bookingId);
      if (booking) {
        booking.status = 'Отменена';
        alert(`Бронь #${bookingId} отменена`);
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
  }
};
</script>

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
          <th>Статус</th>
          <th>Действия</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="booking in bookings" :key="booking.id">
          <td>{{ booking.id }}</td>
          <td>{{ booking.type }}</td>
          <td>{{ booking.orderId }}</td>
          <td>{{ booking.time }}</td>
          <td>{{ booking.name }}</td>
          <td>{{ booking.status }}</td>
          <td>
            <button v-if="booking.status === 'Ожидает подтверждения'" @click="confirmBooking(booking.id)" class="btn">Подтвердить</button>
            <button v-if="booking.status === 'Ожидает подтверждения'" @click="cancelBooking(booking.id)" class="btn danger">Отменить</button>
            <span v-else>{{ booking.status }}</span>
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

  </div>
</template>

<style scoped>
.bookings-dashboard {
  padding: 20px;
  max-width: 900px;
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
