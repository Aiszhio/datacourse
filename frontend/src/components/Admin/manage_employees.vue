<template>
  <div class="manage-employees">
    <h2>Управление сотрудниками</h2>

    <!-- Раздел добавления сотрудника -->
    <div class="add-employee-section">
      <h3>Добавить сотрудника</h3>
      <form @submit.prevent="addEmployee">
        <label for="name">ФИО</label>
        <input type="text" v-model="newEmployee.name" placeholder="ФИО сотрудника" required />
        <label for="name">Должность</label>
        <input type="text" v-model="newEmployee.position" placeholder="Должность" required />
        <label for="name">Дата оформления</label>
        <input type="date" v-model="newEmployee.hireDate" placeholder="Дата оформления" required />
        <label for="name">День рождения</label>
        <input type="date" v-model="newEmployee.birthDate" placeholder="День рождения" required />
        <label for="name">Паспортные данные</label>
        <input type="text" v-model="newEmployee.passport" placeholder="Паспортные данные" required />
        <label for="name">Номер телефона</label>
        <input type="tel" v-model="newEmployee.phone" placeholder="Номер телефона" required />
        <button type="submit" class="btn">Добавить сотрудника</button>
      </form>
    </div>

    <!-- Раздел списка сотрудников -->
    <div class="employee-list-section">
      <h3>Список сотрудников</h3>
      <table>
        <thead>
        <tr>
          <th>Номер</th>
          <th>ФИО</th>
          <th>Должность</th>
          <th>Дата оформления</th>
          <th>День рождения</th>
          <th>Паспортные данные</th>
          <th>Номер телефона</th>
          <th>Действия</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="employee in filteredEmployees" :key="employee.id">
          <td>{{ employee.id }}</td>
          <td>{{ employee.name }}</td>
          <td>{{ employee.position }}</td>
          <td>{{ formatDate(employee.hireDate) }}</td>
          <td>{{ formatDate(employee.birthDate) }}</td>
          <td>{{ employee.passport }}</td>
          <td>{{ employee.phone }}</td>
          <td>
            <button @click="openEditEmployeeModal(employee)" class="btn">Редактировать</button>
            <button @click="deleteEmployee(employee.id)" class="btn danger">Удалить</button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Модальное окно для редактирования сотрудника с подсказками -->
    <div v-if="showEditEmployeeModal" class="modal-overlay">
      <div class="modal">
        <h3>Редактировать сотрудника</h3>

        <label for="name">ФИО</label>
        <input id="name" v-model="currentEmployee.name" class="input" placeholder="ФИО" required />

        <label for="position">Должность</label>
        <input id="position" v-model="currentEmployee.position" class="input" placeholder="Должность" required />

        <label for="hireDate">Дата оформления</label>
        <input id="hireDate" type="date" v-model="currentEmployee.hireDate" class="input" required />

        <label for="birthDate">День рождения</label>
        <input id="birthDate" type="date" v-model="currentEmployee.birthDate" class="input" required />

        <label for="passport">Паспортные данные</label>
        <input id="passport" v-model="currentEmployee.passport" class="input" placeholder="Паспортные данные" required />

        <label for="phone">Номер телефона</label>
        <input id="phone" type="tel" v-model="currentEmployee.phone" class="input" placeholder="Номер телефона" required />

        <button @click="saveEmployee" class="btn">Сохранить</button>
        <button @click="closeEditEmployeeModal" class="btn danger">Отмена</button>
      </div>
    </div>

    <!-- Панель навигации -->
    <div class="card-panel">
      <div class="card" @click="goToAdminHome">
        <h4>На главную</h4>
        <p>Панель администратора</p>
      </div>
      <div class="card" @click="goToOrderHistory">
        <h4>Заказы</h4>
        <p>История заказов</p>
      </div>
      <div class="card" @click="goToMaterialsManagement">
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
export default {
  name: 'ManageEmployees',
  data() {
    return {
      newEmployee: { name: '', position: '', hireDate: '', birthDate: '', passport: '', phone: '' },
      employees: [
        {
          id: 1,
          name: 'Иван Иванов',
          position: 'Фотограф',
          hireDate: '2023-01-15',
          birthDate: '1990-06-01',
          passport: '1234 567890',
          phone: '+7 912 345 6789'
        },
        {
          id: 2,
          name: 'Анна Петрова',
          position: 'Редактор',
          hireDate: '2022-05-20',
          birthDate: '1988-11-15',
          passport: '9876 543210',
          phone: '+7 903 876 5432'
        }
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
        hireDate: this.newEmployee.hireDate,
        birthDate: this.newEmployee.birthDate,
        passport: this.newEmployee.passport,
        phone: this.newEmployee.phone
      };
      this.employees.push(employee);
      this.newEmployee = { name: '', position: '', hireDate: '', birthDate: '', passport: '', phone: '' };
    },
    deleteEmployee(id) {
      this.employees = this.employees.filter(employee => employee.id !== id);
      alert(`Сотрудник с ID #${id} удален.`);
    },
    openEditEmployeeModal(employee) {
      this.currentEmployee = { ...employee };
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
    },
    goToBookingsPage() {
      this.$router.push({ name: 'Bookings' });
    },
    goToServicesPage() {
      this.$router.push({ name: 'Services' });
    }
  }
};
</script>

<style scoped>
.manage-employees {
  padding: 20px;
  margin: 0 auto;
  max-width: 1500px;
  max-height: 100vh;
}

h2 {
  text-align: center;
  margin-bottom: 20px;
}

.add-employee-section,
.employee-list-section {
  background-color: #f9f9f9;
  padding: 20px;
  margin-bottom: 20px;
  border-radius: 10px;
  box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
}

label {
  display: block;
  font-weight: bold;
  margin-top: 10px;
}

input[type="text"],
input[type="date"],
input[type="tel"],
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
  margin-bottom: 5px;
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
