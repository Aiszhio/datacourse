<template>
  <div class="order-history">
    <h2>История заказов</h2>

    <!-- Раздел сортировки заказов -->
    <div class="sort-buttons">
      <button @click="sortBy('clientName')" class="btn">
        Сортировать по имени клиента
        <span v-if="currentSortKey === 'clientName'">
          {{ sortOrders.clientName === 'asc' ? '▲' : '▼' }}
        </span>
      </button>
      <button @click="sortBy('employeeName')" class="btn">
        Сортировать по имени сотрудника
        <span v-if="currentSortKey === 'employeeName'">
          {{ sortOrders.employeeName === 'asc' ? '▲' : '▼' }}
        </span>
      </button>
      <button @click="sortBy('service')" class="btn">
        Сортировать по услуге
        <span v-if="currentSortKey === 'service'">
          {{ sortOrders.service === 'asc' ? '▲' : '▼' }}
        </span>
      </button>
      <button @click="sortBy('orderDate')" class="btn">
        Сортировать по дате оформления
        <span v-if="currentSortKey === 'orderDate'">
          {{ sortOrders.orderDate === 'asc' ? '▲' : '▼' }}
        </span>
      </button>
      <button @click="sortBy('receiveDate')" class="btn">
        Сортировать по дате получения
        <span v-if="currentSortKey === 'receiveDate'">
          {{ sortOrders.receiveDate === 'asc' ? '▲' : '▼' }}
        </span>
      </button>
    </div>

    <!-- Индикатор загрузки -->
    <div v-if="isLoading" class="loading">
      Загрузка заказов...
    </div>

    <!-- Раздел истории заказов -->
    <div class="order-list-section" v-else>
      <h3>Список заказов</h3>
      <table class="data-table">
        <thead>
        <tr>
          <th>Номер заказа</th>
          <th>Имя клиента</th>
          <th>Имя сотрудника</th>
          <th>Название услуги</th>
          <th>Дата оформления</th>
          <th>Дата получения</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="order in sortedOrders" :key="order.id">
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

    <!-- Панель навигации -->
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
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'OrderHistory',
  data() {
    return {
      orders: [],
      sortOrders: {
        clientName: 'asc',
        employeeName: 'asc',
        service: 'asc',
        orderDate: 'asc',
        receiveDate: 'asc',
      },
      currentSortKey: '',
      isLoading: true, // Состояние загрузки
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

        // Для дат преобразуем в объекты Date
        if (this.currentSortKey === 'orderDate' || this.currentSortKey === 'receiveDate') {
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
    async fetchOrders() {
      this.isLoading = true; // Начинаем загрузку
      try {
        const response = await axios.get('http://localhost:8080/api/orders/admin', {
          withCredentials: true, // Отправка куки с запросом
        });
        this.orders = response.data.orders || [];
        // Устанавливаем начальную сортировку после загрузки данных
        if (this.currentSortKey) {
          this.sortBy(this.currentSortKey);
        }
      } catch (error) {
        console.error('Ошибка при загрузке заказов:', error);
        alert('Не удалось загрузить заказы.');
      } finally {
        this.isLoading = false; // Завершаем загрузку
      }
    },
    sortBy(key) {
      if (this.currentSortKey === key) {
        // Переключаем направление сортировки
        this.sortOrders[key] = this.sortOrders[key] === 'asc' ? 'desc' : 'asc';
      } else {
        // Устанавливаем новое поле сортировки
        this.currentSortKey = key;
      }
    },
    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleDateString('ru-RU', { year: 'numeric', month: 'long', day: 'numeric' });
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
  max-width: 900px;
  max-height: 100vh;
}

h2 {
  text-align: center;
  margin-bottom: 20px;
}

.sort-buttons {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.order-list-section {
  background-color: #f9f9f9;
  padding: 20px;
  margin-bottom: 20px;
  border-radius: 10px;
  box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.1);
}

label {
  display: block;
  font-weight: bold;
  margin-top: 10px;
}

.input {
  padding: 10px;
  margin: 5px 0;
  border: 1px solid #ccc;
  border-radius: 5px;
  width: 90%;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 10px;
}

table th,
table td {
  padding: 12px;
  border: 1px solid #ddd;
  text-align: left;
}

.card-panel {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
  justify-content: center;
  margin-top: 20px;
}

.card {
  background-color: #4CAF50;
  border-radius: 10px;
  padding: 10px;
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

.btn {
  background-color: #4CAF50;
  color: white;
  padding: 10px 15px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 1em;
  transition: background-color 0.3s ease;
  margin: 5px;
}

.btn.danger {
  background-color: #f44336;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal {
  background: white;
  padding: 20px;
  border-radius: 10px;
  max-width: 400px;
  width: 100%;
}

.btn:hover {
  background-color: #45a049;
}
</style>
