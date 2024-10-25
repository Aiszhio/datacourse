<template>
  <div class="admin-dashboard">
    <h2>Панель администратора</h2>

    <!-- Раздел со списком сотрудников -->
    <div class="section">
      <h3>Список сотрудников</h3>
      <table>
        <thead>
        <tr>
          <th>ID</th>
          <th>Имя</th>
          <th>Должность</th>
          <th>Дата найма</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="employee in displayedEmployees" :key="employee.id">
          <td>{{ employee.id }}</td>
          <td>{{ employee.name }}</td>
          <td>{{ employee.position }}</td>
          <td>{{ formatDate(employee.hireDate) }}</td>
        </tr>
        </tbody>
      </table>
      <button @click="goToEmployeesPage" class="btn">Управление сотрудниками</button>
    </div>

    <!-- Раздел с историей заказов -->
    <div class="section">
      <h3>История заказов</h3>
      <table>
        <thead>
        <tr>
          <th>ID Заказа</th>
          <th>Клиент</th>
          <th>Услуга</th>
          <th>Дата заказа</th>
          <th>Статус</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="order in displayedOrders" :key="order.id">
          <td>{{ order.id }}</td>
          <td>{{ order.client }}</td>
          <td>{{ order.service }}</td>
          <td>{{ formatDate(order.orderDate) }}</td>
          <td>{{ order.status }}</td>
        </tr>
        </tbody>
      </table>
      <button @click="goToOrdersPage" class="btn">История заказов</button>
    </div>

    <!-- Раздел с историей закупок -->
    <div class="section">
      <h3>История закупок</h3>
      <table>
        <thead>
        <tr>
          <th>ID Закупки</th>
          <th>Материал</th>
          <th>Поставщик</th>
          <th>Количество</th>
          <th>Стоимость</th>
          <th>Дата закупки</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="purchase in displayedPurchases" :key="purchase.id">
          <td>{{ purchase.id }}</td>
          <td>{{ purchase.material }}</td>
          <td>{{ purchase.supplier }}</td>
          <td>{{ purchase.quantity }}</td>
          <td>{{ purchase.cost }}</td>
          <td>{{ formatDate(purchase.date) }}</td>
        </tr>
        </tbody>
      </table>
      <button @click="goToMaterialsPage" class="btn">Просмотр закупок и расходов материалов</button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AdminDashboard',
  data() {
    return {
      employees: [
        { id: 1, name: 'Иван Иванов', position: 'Фотограф', hireDate: '2023-01-15' },
        { id: 2, name: 'Анна Петрова', position: 'Редактор', hireDate: '2022-05-20' },
        // Дополнительные записи
      ],
      orders: [
        { id: 101, client: 'Мария Сидорова', service: 'Фотосессия', orderDate: '2024-08-15', status: 'Выполнен' },
        { id: 102, client: 'Петр Иванов', service: 'Фотоальбом', orderDate: '2024-09-01', status: 'В процессе' },
        // Дополнительные записи
      ],
      purchases: [
        { id: 301, material: 'Фотобумага', supplier: 'Art Supplies Inc.', quantity: 100, cost: 1500, date: '2024-08-01' },
        { id: 302, material: 'Картридж для принтера', supplier: 'Print Solutions Ltd.', quantity: 50, cost: 5000, date: '2024-08-05' },
        // Дополнительные записи
      ]
    };
  },
  computed: {
    displayedEmployees() {
      return this.employees.slice(0, 5); // Показываем только первых 5 сотрудников
    },
    displayedOrders() {
      return this.orders.slice(0, 5); // Показываем только первые 5 заказов
    },
    displayedPurchases() {
      return this.purchases.slice(0, 5); // Показываем только первые 5 закупок
    }
  },
  methods: {
    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleDateString('ru-RU', { year: 'numeric', month: 'long', day: 'numeric' });
    },
    goToEmployeesPage() {
      this.$router.push({ name: 'ManageEmp' });
    },
    goToOrdersPage() {
      this.$router.push({ name: 'OrderHistory' });
    },
    goToMaterialsPage() {
      this.$router.push({ name: 'MaterialsOverview' });
    }
  }
};
</script>

<style scoped>
.section {
  margin-bottom: 20px;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 10px;
}

table th,
table td {
  padding: 8px;
  border: 1px solid #ddd;
  text-align: left;
}

.btn {
  background-color: #4CAF50;
  color: white;
  padding: 8px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  margin: 10px;
}

.btn:hover {
  background-color: #45a049;
}
</style>
