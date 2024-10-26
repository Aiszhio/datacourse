<template>
  <div class="manage-employees">
    <h2>Управление сотрудниками</h2>

    <!-- Раздел добавления сотрудника -->
    <div class="add-employee-section">
      <h3>Добавить сотрудника</h3>
      <form @submit.prevent="addEmployee">
        <input type="text" v-model="newEmployee.name" placeholder="ФИО сотрудника" required />
        <input type="text" v-model="newEmployee.position" placeholder="Должность" required />
        <input type="date" v-model="newEmployee.hireDate" required />
        <button type="submit" class="btn">Добавить сотрудника</button>
      </form>
    </div>

    <!-- Раздел фильтра сотрудников -->
    <div class="filter-section">
      <h3>Фильтр по должности</h3>
      <input type="text" v-model="positionFilter" placeholder="Введите должность" />
    </div>

    <!-- Раздел списка сотрудников -->
    <div class="employee-list-section">
      <h3>Список сотрудников</h3>
      <table>
        <thead>
        <tr>
          <th>ID</th>
          <th>ФИО</th>
          <th>Должность</th>
          <th>Дата найма</th>
          <th>Действия</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="employee in filteredEmployees" :key="employee.id">
          <td>{{ employee.id }}</td>
          <td>{{ employee.name }}</td>
          <td>{{ employee.position }}</td>
          <td>{{ formatDate(employee.hireDate) }}</td>
          <td>
            <button @click="openEditEmployeeModal(employee)" class="btn">Редактировать</button>
            <button @click="deleteEmployee(employee.id)" class="btn danger">Удалить</button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Модальное окно для редактирования сотрудника -->
    <div v-if="showEditEmployeeModal" class="modal-overlay">
      <div class="modal">
        <h3>Редактировать сотрудника</h3>
        <input v-model="currentEmployee.name" class="input" placeholder="ФИО" required />
        <input v-model="currentEmployee.position" class="input" placeholder="Должность" required />
        <input type="date" v-model="currentEmployee.hireDate" class="input" required />
        <button @click="saveEmployee" class="btn">Сохранить</button>
        <button @click="closeEditEmployeeModal" class="btn danger">Отмена</button>
      </div>
    </div>

    <!-- Кнопки для навигации по другим страницам -->
    <div class="navigation-buttons">
      <button @click="goToAdminHome" class="btn">На главную администратора</button>
      <button @click="goToOrderHistory" class="btn">История заказов</button>
      <button @click="goToMaterialsManagement" class="btn">Управление материалами</button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ManageEmployees',
  data() {
    return {
      newEmployee: { name: '', position: '', hireDate: '' },
      employees: [
        { id: 1, name: 'Иван Иванов', position: 'Фотограф', hireDate: '2023-01-15' },
        { id: 2, name: 'Анна Петрова', position: 'Редактор', hireDate: '2022-05-20' }
      ],
      showEditEmployeeModal: false,
      currentEmployee: {},
      positionFilter: ''
    };
  },
  computed: {
    filteredEmployees() {
      if (this.positionFilter) {
        return this.employees.filter(employee =>
            employee.position.toLowerCase().includes(this.positionFilter.toLowerCase())
        );
      }
      return this.employees;
    }
  },
  methods: {
    addEmployee() {
      const newId = this.employees.length + 1;
      const employee = {
        id: newId,
        name: this.newEmployee.name,
        position: this.newEmployee.position,
        hireDate: this.newEmployee.hireDate
      };

      this.employees.push(employee);
      this.newEmployee = { name: '', position: '', hireDate: '' };
    },
    deleteEmployee(id) {
      this.employees = this.employees.filter(employee => employee.id !== id);
      alert(`Сотрудник с ID #${id} удален.`);
    },
    openEditEmployeeModal(employee) {
      this.currentEmployee = { ...employee }; // Копируем данные сотрудника в текущий объект
      this.showEditEmployeeModal = true;
    },
    closeEditEmployeeModal() {
      this.showEditEmployeeModal = false;
      this.currentEmployee = {};
    },
    saveEmployee() {
      const index = this.employees.findIndex(employee => employee.id === this.currentEmployee.id);
      if (index !== -1) {
        this.employees.splice(index, 1, { ...this.currentEmployee });
        alert(`Сотрудник с ID #${this.currentEmployee.id} обновлен.`);
      }
      this.closeEditEmployeeModal();
    },
    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleDateString('ru-RU', { year: 'numeric', month: 'long', day: 'numeric' });
    },
    goToAdminHome() {
      this.$router.push({ name: 'AdminHome' });
    },
    goToOrderHistory() {
      this.$router.push({ name: 'OrderHistory' });
    },
    goToMaterialsManagement() {
      this.$router.push({ name: 'MaterialsOverview' });
    }
  }
};
</script>

<style scoped>
.manage-employees {
  padding: 20px;
  margin: 0 auto;
  max-width: 900px;
  max-height: 100vh;
}

h2 {
  text-align: center;
  margin-bottom: 20px;
}

.add-employee-section,
.filter-section,
.employee-list-section {
  background-color: #f9f9f9;
  padding: 20px;
  margin-bottom: 20px;
  border-radius: 10px;
  box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.1);
}

input[type="text"],
input[type="date"],
.input {
  padding: 10px;
  margin: 10px 0;
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

.navigation-buttons {
  display: flex;
  justify-content: space-around;
  margin-top: 20px;
}

.btn {
  background-color: #4CAF50;
  color: white;
  padding: 10px 15px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 1em;
  margin: 10px;
  transition: background-color 0.3s ease;
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
