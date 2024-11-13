<template>
  <div class="worker-dashboard">
    <h2>Панель фотографа</h2>
    <h3>Добро пожаловать, {{ workerName }}!</h3>

    <!-- Список текущих заказов -->
    <div class="orders-section">
      <h3>Заказы</h3>
      <div v-if="loading" class="loading">Загрузка заказов...</div>

      <!-- Таблица заказов -->
      <table v-if="!loading && limitedOrders.length > 0" class="orders-table">
        <thead>
        <tr>
          <th>Номер заказа</th>
          <th>Номер клиента</th>
          <th>Номер сотрудника</th>
          <th>Название услуги</th>
          <th>Дата оформления</th>
          <th>Дата получения</th>
          <th>Действия</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="order in limitedOrders" :key="order.OrderID">
          <td>{{ order.OrderID }}</td>
          <td>{{ order.ClientID }}</td>
          <td>{{ order.EmployeeName || "Не назначен" }}</td>
          <td>{{ order.ServiceName }}</td>
          <td>{{ formatDate(order.OrderDate) }}</td>
          <td>{{ formatDate(order.ReceiveDate) }}</td>
          <td>
            <button v-if="!order.isAssigned" @click="openTakeOrderModal(order)" class="btn">Взять</button>
            <button v-else @click="releaseOrder(order)" class="btn danger">Отказаться</button>
          </td>
        </tr>
        </tbody>
      </table>
      <div v-if="!loading && limitedOrders.length === 0" class="no-orders">Нет доступных заказов</div>
    </div>

    <!-- Модальное окно для взятия заказа -->
    <div v-if="showTakeOrderModal" class="modal-overlay">
      <div class="modal">
        <h3>Оформление заказа</h3>
        <label for="workerNameInput">ФИО сотрудника</label>
        <input id="workerNameInput" v-model="workerNameInput" class="input" placeholder="Введите ФИО" required />
        <button @click="confirmTakeOrder" class="btn">Подтвердить</button>
        <button @click="closeTakeOrderModal" class="btn danger">Отмена</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'WorkerDashboard',
  data() {
    return {
      workerName: '',              // Имя сотрудника
      workerNameInput: '',         // Поле для ввода ФИО в модальном окне
      orders: [],                  // Список заказов
      loading: true,               // Индикатор загрузки заказов
      orderLimit: 3,               // Лимит отображаемых заказов
      showTakeOrderModal: false,   // Показ модального окна
      currentOrder: null           // Текущий выбранный заказ
    };
  },
  computed: {
    limitedOrders() {
      return this.orders.slice(0, this.orderLimit);
    }
  },
  async mounted() {
    await this.fetchWorkerData();
    await this.fetchOrders();
  },
  methods: {
    async fetchWorkerData() {
      try {
        const response = await fetch('http://localhost:8080/api/user', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json'
          },
          credentials: 'include'
        });

        if (!response.ok) {
          throw new Error('Ошибка при получении имени сотрудника');
        }

        const data = await response.json();
        this.workerName = data.name;
      } catch (error) {
        console.error('Ошибка при получении имени сотрудника:', error.message);
        alert('Не удалось загрузить имя сотрудника.');
      }
    },
    async fetchOrders() {
      try {
        const response = await fetch('http://localhost:8080/api/orders', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json'
          },
          credentials: 'include'
        });

        if (!response.ok) {
          throw new Error('Ошибка при получении заказов');
        }

        const data = await response.json();
        this.orders = data;
      } catch (error) {
        console.error('Ошибка при загрузке заказов:', error.message);
        alert('Не удалось загрузить заказы.');
      } finally {
        this.loading = false;
      }
    },
    openTakeOrderModal(order) {
      this.currentOrder = order;
      this.showTakeOrderModal = true;
    },
    closeTakeOrderModal() {
      this.workerNameInput = '';
      this.showTakeOrderModal = false;
    },
    async confirmTakeOrder() {
      if (!this.workerNameInput) {
        alert('Пожалуйста, введите ФИО');
        return;
      }

      try {
        const response = await fetch(`http://localhost:8080/api/orders/${this.currentOrder.OrderID}/assign`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ employeeName: this.workerNameInput })
        });

        if (!response.ok) {
          throw new Error('Ошибка при взятии заказа');
        }

        alert(`Заказ #${this.currentOrder.OrderID} взят на исполнение`);
        this.currentOrder.isAssigned = true;
        this.currentOrder.EmployeeName = this.workerNameInput;
        this.closeTakeOrderModal();

      } catch (error) {
        console.error('Ошибка при взятии заказа:', error.message);
        alert('Не удалось взять заказ.');
      }
    },
    async releaseOrder(order) {
      try {
        const response = await fetch(`http://localhost:8080/api/orders/${order.OrderID}/unassign`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json'
          }
        });

        if (!response.ok) {
          throw new Error('Ошибка при отказе от заказа');
        }

        alert(`Вы отказались от заказа #${order.OrderID}`);
        order.isAssigned = false;
        order.EmployeeName = null;

      } catch (error) {
        console.error('Ошибка при отказе от заказа:', error.message);
        alert('Не удалось отказаться от заказа.');
      }
    },
    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleDateString('ru-RU', { year: 'numeric', month: 'long', day: 'numeric' });
    }
  }
};
</script>

<style scoped>
.worker-dashboard {
  padding: 20px;
  margin: 0 auto;
  max-width: 900px;
}

.orders-section {
  background-color: #f9f9f9;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.orders-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
}

.orders-table th,
.orders-table td {
  padding: 12px;
  border: 1px solid #ddd;
  text-align: left;
}

.orders-table th {
  background-color: #f2f2f2;
}

.no-orders {
  text-align: center;
  color: #999;
}

.btn {
  background-color: #4CAF50;
  color: white;
  padding: 10px;
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

.btn.danger {
  background-color: #f44336;
}

.btn.danger:hover {
  background-color: #d32f2f;
}

.loading {
  text-align: center;
  color: #555;
  font-size: 1.2em;
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

.input {
  width: 100%;
  padding: 10px;
  margin: 10px 0;
  border: 1px solid #ccc;
  border-radius: 5px;
  box-sizing: border-box;
}
</style>
