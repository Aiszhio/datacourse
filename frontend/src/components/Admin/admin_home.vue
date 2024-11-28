<template>
  <div class="admin-dashboard">
    <h2>Добро пожаловать, {{ adminName }}!</h2>
    <p>Сегодня: {{ currentDate }}</p>

    <!-- Панель с карточками первого ряда -->
    <div class="card-panel">
      <div class="card" @click="goToEmployeesPage">
        <i class="fas fa-user-tie icon"></i>
        <h3>Сотрудники</h3>
        <p>Управление персоналом</p>
      </div>
      <div class="card" @click="goToOrdersPage">
        <i class="fas fa-clipboard-list icon"></i>
        <h3>Заказы</h3>
        <p>История заказов</p>
      </div>
      <div class="card" @click="goToMaterialsPage">
        <i class="fas fa-box icon"></i>
        <h3>Закупки</h3>
        <p>Материалы и расходы</p>
      </div>
    </div>

    <!-- Панель с карточками второго ряда -->
    <div class="card-panel">
      <div class="card" @click="goToServicesPage">
        <i class="fas fa-concierge-bell icon"></i>
        <h3>Услуги</h3>
        <p>Список доступных услуг</p>
      </div>
      <div class="card" @click="goToBookingPage">
        <i class="fas fa-calendar-check icon"></i>
        <h3>Бронирование</h3>
        <p>Управление бронированием</p>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AdminDashboard',
  data() {
    return {
      adminName: '', // Имя администратора
      currentDate: new Date().toLocaleDateString('ru-RU', {
        weekday: 'long', year: 'numeric', month: 'long', day: 'numeric'
      })
    };
  },
  async mounted() {
    await this.fetchAdminName(); // Загружаем имя администратора при монтировании компонента
  },
  methods: {
    async fetchAdminName() {
      try {
        const response = await fetch('http://localhost:8080/api/user', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json'
          },
          credentials: 'include' // Отправляем куки с запросом
        });

        if (!response.ok) {
          throw new Error('Ошибка при получении имени администратора');
        }

        const data = await response.json();
        this.adminName = data.name; // Обновляем adminName с полученным именем

      } catch (error) {
        console.error('Ошибка при получении имени администратора:', error.message);
        alert('Не удалось загрузить имя администратора.');
      }
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
    goToBookingPage() {
      this.$router.push({ name: 'Bookings' });
    },
  }
};
</script>

<style scoped>
.admin-dashboard {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  text-align: center;
}

h2 {
  color: #333;
  margin-bottom: 10px;
}

p {
  color: #B0C4DE;
  margin-bottom: 20px;
}

.card-panel {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
  justify-content: center;
}

.card {
  background-color: #4CAF50;
  border-radius: 10px;
  padding: 15px;
  margin: 5px;
  width: 180px;
  text-align: center;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  transition: transform 0.3s, box-shadow 0.3s;
  cursor: pointer;
}

.card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
}

.icon {
  font-size: 40px;
  color: #FFFFE0;
  margin-bottom: 10px;
}

h3 {
  color: #FFFFE0;
  margin-bottom: 5px;
  font-size: 18px;
}
</style>
