<template>
  <div class="order-history">
    <h2>История заказов</h2>

    <!-- Раздел добавления пользователя -->
    <div class="add-user-section">
      <h3>Добавить нового пользователя</h3>
      <form @submit.prevent="addUser">
        <input type="text" v-model="newUser.name" placeholder="Имя пользователя" required />
        <input type="email" v-model="newUser.email" placeholder="Электронная почта" required />
        <input type="tel" v-model="newUser.phone" placeholder="Номер телефона" required />
        <button type="submit" class="btn">Добавить пользователя</button>
      </form>
    </div>

    <!-- Раздел истории заказов -->
    <div class="order-list-section">
      <h3>Список заказов</h3>
      <div class="sort-buttons">
        <button @click="sortBy('client')" class="btn">Сортировать по клиенту</button>
        <button @click="sortBy('service')" class="btn">Сортировать по услуге</button>
        <button @click="sortBy('orderDate')" class="btn">Сортировать по дате</button>
        <button @click="sortBy('status')" class="btn">Сортировать по статусу</button>
      </div>

      <table>
        <thead>
        <tr>
          <th>ID Заказа</th>
          <th>Клиент</th>
          <th>Услуга</th>
          <th>Дата заказа</th>
          <th>Статус</th>
          <th>Действия</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="order in orders" :key="order.id">
          <td>{{ order.id }}</td>
          <td>{{ order.client }}</td>
          <td>{{ order.service }}</td>
          <td>{{ formatDate(order.orderDate) }}</td>
          <td>{{ order.status }}</td>
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
        <input v-model="currentOrder.client" class="input" placeholder="Клиент" required />
        <input v-model="currentOrder.service" class="input" placeholder="Услуга" required />
        <input type="date" v-model="currentOrder.orderDate" class="input" required />
        <select v-model="currentOrder.status" class="input">
          <option value="ожидает подтверждение брони">Ожидает подтверждение брони</option>
          <option value="забронирован">Забронирован</option>
          <option value="отменен">Отменен</option>
          <option value="выполнен">Выполнен</option>
        </select>
        <button @click="saveOrder" class="btn">Сохранить</button>
        <button @click="closeEditOrderModal" class="btn danger">Отмена</button>
      </div>
    </div>

    <!-- Кнопки для навигации по другим страницам -->
    <div class="navigation-buttons">
      <button @click="goToAdminHome" class="btn">На главную администратора</button>
      <button @click="goToEmployeesPage" class="btn">Управление сотрудниками</button>
      <button @click="goToMaterialsPage" class="btn">Управление материалами</button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'OrderHistory',
  data() {
    return {
      newUser: { name: '', email: '', phone: '' },
      orders: [
        { id: 101, client: 'Мария Сидорова', service: 'Фотосессия', orderDate: '2024-08-15', status: 'выполнен' },
        { id: 102, client: 'Петр Иванов', service: 'Фотоальбом', orderDate: '2024-09-01', status: 'ожидает подтверждение брони' }
      ],
      showEditOrderModal: false,
      currentOrder: {}
    };
  },
  methods: {
    addUser() {
      const user = {
        name: this.newUser.name,
        email: this.newUser.email,
        phone: this.newUser.phone
      };
      alert(`Пользователь ${user.name} (${user.email}, тел. ${user.phone}) добавлен.`);
      this.newUser = { name: '', email: '', phone: '' };
    },
    deleteOrder(id) {
      this.orders = this.orders.filter(order => order.id !== id);
      alert(`Заказ #${id} удален.`);
    },
    openEditOrderModal(order) {
      this.currentOrder = { ...order }; // Копируем данные заказа в текущий объект
      this.showEditOrderModal = true;
    },
    closeEditOrderModal() {
      this.showEditOrderModal = false;
      this.currentOrder = {};
    },
    saveOrder() {
      const index = this.orders.findIndex(order => order.id === this.currentOrder.id);
      if (index !== -1) {
        this.orders.splice(index, 1, { ...this.currentOrder });
        alert(`Заказ #${this.currentOrder.id} обновлен.`);
      }
      this.closeEditOrderModal();
    },
    sortBy(key) {
      this.orders.sort((a, b) => {
        if (a[key] < b[key]) return -1;
        if (a[key] > b[key]) return 1;
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
    }
  }
};
</script>

<style scoped>
.order-history {
  display: flex;
  flex-direction: column;
  justify-content: space-between; /* Это раздвигает заголовок, контент и кнопки */
  padding: 20px;
  margin: 0 auto;
  max-width: 900px;
  height: 100vh; /* Используем 100vh для полной высоты */
}

.order-history .content {
  flex-grow: 1; /* Контент будет занимать все доступное пространство */
  overflow-y: auto; /* Если контент превышает высоту, включаем прокрутку */
}

h2 {
  text-align: center;
  margin-bottom: 20px;
}

.add-user-section,
.order-list-section {
  background-color: #f9f9f9;
  padding: 20px;
  margin-bottom: 20px;
  border-radius: 10px;
  box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.1);
}

input[type="text"],
input[type="email"],
input[type="tel"],
select,
input[type="date"] {
  padding: 10px;
  margin: 10px 0;
  border: 1px solid #ccc;
  border-radius: 5px;
  width: 100%;
  box-sizing: border-box;
  font-size: 1em;
  transition: border-color 0.3s ease;
}

/* Стили для input при редактировании */
.input {
  padding: 12px;
  margin-top: 15px;
  border: 1px solid #ddd;
  border-radius: 8px;
  width: 100%;
  box-sizing: border-box;
  font-size: 1em;
  background-color: #f9f9f9;
  transition: border-color 0.3s ease;
}

.input:focus {
  border-color: #4CAF50; /* Зеленый цвет при фокусе */
  outline: none;
  background-color: #ffffff;
}

/* Дополнительные стили для полей ФИО и Услуга */
input[type="text"].input-name,
input[type="text"].input-service {
  font-weight: 500;
  color: #333;
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

select {
  padding: 8px;
  border-radius: 5px;
  border: 1px solid #ccc;
}

.navigation-buttons {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  margin-top: 20px;
  padding-bottom: 20px; /* Добавляем отступ снизу для красоты */
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
  margin-left: 3px;
  margin-top: 3px;
}

.btn:hover {
  background-color: #45a049;
}
</style>
