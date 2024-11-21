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

    <!-- Раздел истории заказов -->
    <div class="order-list-section">
      <h3>Список заказов</h3>
      <table>
        <thead>
        <tr>
          <th>Номер заказа</th>
          <th>Имя клиента</th>
          <th>Имя сотрудника</th>
          <th>Название услуги</th>
          <th>Дата оформления</th>
          <th>Дата получения</th>
          <th>Действия</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="order in orders" :key="order.id">
          <td>{{ order.id }}</td>
          <td>{{ order.clientName }}</td>
          <td>{{ order.employeeName }}</td>
          <td>{{ order.service }}</td>
          <td>{{ formatDate(order.orderDate) }}</td>
          <td>{{ formatDate(order.receiveDate) }}</td>
          <td>
            <button @click="openEditOrderModal(order)" class="btn">Редактировать</button>
            <button @click="deleteOrder(order.id)" class="btn danger">Удалить</button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Модальное окно для редактирования заказа -->
    <div v-if="showEditOrderModal" class="modal-overlay">
      <div class="modal">
        <h3>Редактировать заказ</h3>

        <!-- Отображение имен клиента и сотрудника как статичных полей -->
        <label>Имя клиента</label>
        <input
            v-model="currentOrder.clientName"
            class="input"
            placeholder="Имя клиента"
            required
            readonly
        />

        <label>Имя сотрудника</label>
        <input
            v-model="currentOrder.employeeName"
            class="input"
            placeholder="Имя сотрудника"
            required
            readonly
        />

        <label for="service">Название услуги</label>
        <input
            id="service"
            v-model="currentOrder.service"
            class="input"
            placeholder="Название услуги"
            required
        />

        <label for="orderDate">Дата оформления</label>
        <input
            id="orderDate"
            type="date"
            v-model="currentOrder.orderDate"
            class="input"
            required
        />

        <label for="receiveDate">Дата получения</label>
        <input
            id="receiveDate"
            type="date"
            v-model="currentOrder.receiveDate"
            class="input"
            required
        />

        <button @click="saveOrder" class="btn">Сохранить</button>
        <button @click="closeEditOrderModal" class="btn danger">Отмена</button>
      </div>
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
      showEditOrderModal: false,
      currentOrder: {},
      sortOrders: {
        clientName: 'asc',
        employeeName: 'asc',
        service: 'asc',
        orderDate: 'asc',
        receiveDate: 'asc',
      },
      currentSortKey: '',
    };
  },
  created() {
    this.fetchOrders();
  },
  methods: {
    async fetchOrders() {
      try {
        const response = await axios.get('http://localhost:8080/api/orders/admin', {
          withCredentials: true, // Отправка куки с запросом
        });
        this.orders = response.data.orders;
      } catch (error) {
        console.error('Ошибка при загрузке заказов:', error);
      }
    },
    async deleteOrder(id) {
      if (!confirm(`Вы уверены, что хотите удалить заказ #${id}?`)) return;
      try {
        await axios.delete(`http://localhost:8080/api/orders/${id}`, {
          withCredentials: true, // Отправка куки с запросом
        });
        this.orders = this.orders.filter((order) => order.id !== id);
        alert(`Заказ #${id} удален.`);
      } catch (error) {
        console.error(`Ошибка при удалении заказа #${id}:`, error);
        alert('Не удалось удалить заказ.');
      }
    },
    openEditOrderModal(order) {
      this.currentOrder = { ...order };
      this.showEditOrderModal = true;
    },
    closeEditOrderModal() {
      this.showEditOrderModal = false;
      this.currentOrder = {};
    },
    async saveOrder() {
      try {
        await axios.put(`http://localhost:8080/api/orders/${this.currentOrder.id}`, this.currentOrder, {
          withCredentials: true, // Отправка куки с запросом
        });
        const index = this.orders.findIndex((order) => order.id === this.currentOrder.id);
        if (index !== -1) {
          this.orders.splice(index, 1, { ...this.currentOrder });
          alert(`Заказ #${this.currentOrder.id} обновлен.`);
        }
        this.closeEditOrderModal();
      } catch (error) {
        console.error(`Ошибка при обновлении заказа #${this.currentOrder.id}:`, error);
        alert('Не удалось обновить заказ.');
      }
    },
    sortBy(key) {
      if (this.currentSortKey === key) {
        // Переключение направления сортировки
        this.sortOrders[key] = this.sortOrders[key] === 'asc' ? 'desc' : 'asc';
      } else {
        // Установка направления сортировки по умолчанию (возрастание) для нового ключа
        this.sortOrders[key] = 'asc';
        this.currentSortKey = key;
      }

      const direction = this.sortOrders[key];
      this.orders.sort((a, b) => {
        let aVal = a[key];
        let bVal = b[key];

        // Для дат преобразуем в объекты Date
        if (key === 'orderDate' || key === 'receiveDate') {
          aVal = new Date(aVal);
          bVal = new Date(bVal);
        } else {
          // Приводим к строкам для сравнения
          aVal = aVal.toString().toLowerCase();
          bVal = bVal.toString().toLowerCase();
        }

        if (aVal < bVal) return direction === 'asc' ? -1 : 1;
        if (aVal > bVal) return direction === 'asc' ? 1 : -1;
        return 0;
      });
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
