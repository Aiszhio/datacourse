<template>
  <div class="manage-employees">
    <h2>Управление сотрудниками</h2>

    <div class="add-employee-section">
      <h3>Добавить сотрудника</h3>
      <form @submit.prevent="addEmployee">
        <label for="name">ФИО</label>
        <input
            type="text"
            v-model="newEmployee.full_name"
            placeholder="ФИО сотрудника"
            required
            maxlength="50"
        />

        <label for="position">Должность</label>
        <select
            class="position"
            v-model="newEmployee.position"
            required
        >
          <option disabled value="">Выберите должность</option>
          <option>Фотограф</option>
          <option>Администратор</option>
        </select>

        <label for="hireDate">Дата оформления</label>
        <input
            type="date"
            id="hireDate"
            v-model="newEmployee.hire_date"
            :max="todayDate"
            required
        />

        <label for="birthDate">День рождения</label>
        <input
            type="date"
            v-model="newEmployee.birth_date"
            :max="minBirthDate"
            required
        />

        <label for="passport">Паспортные данные</label>
        <input
            type="text"
            v-model="formattedPassportData"
            placeholder="Паспортные данные"
            required
            maxlength="10"
        />

        <label for="phone">Номер телефона</label>
        <input
            type="tel"
            v-model="formattedPhoneNumber"
            placeholder="Номер телефона"
            required
            maxlength="12"
        />

        <button type="submit" class="btn">Добавить сотрудника</button>
      </form>
    </div>

    <div class="employee-list-section">
      <h3>Список сотрудников</h3>

      <div v-if="loadingEmployees" class="loading">Загрузка сотрудников...</div>
      <div v-else-if="employees.length === 0" class="no-employees">Нет сотрудников</div>

      <table v-else class="employees-table">
        <thead>
        <tr>
          <th>№</th>
          <th @click="sortBy('name')" class="sortable">
            ФИО
            <span v-if="sortKey === 'name'">
                <span v-if="sortOrder === 'asc'">&#9650;</span>
                <span v-else>&#9660;</span>
              </span>
          </th>
          <th>Должность</th>
          <th>Дата оформления</th>
          <th>День рождения</th>
          <th>Паспортные данные</th>
          <th>Номер телефона</th>
          <th @click="sortBy('status')" class="sortable">
            Статус
            <span v-if="sortKey === 'status'">
                <span v-if="sortOrder === 'asc'">&#9650;</span>
                <span v-else>&#9660;</span>
              </span>
          </th>
          <th>Действия</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(employee, index) in displayedEmployees" :key="employee.id">
          <td>{{ index + 1 }}</td>
          <td>{{ employee.name }}</td>
          <td>{{ employee.position }}</td>
          <td>{{ formatDate(employee.hireDate) }}</td>
          <td>{{ formatDate(employee.birthDate) }}</td>
          <td>{{ employee.passport_data }}</td>
          <td>{{ formatPhoneNumber(employee.phone_number) }}</td>
          <td>
              <span :class="{
                'status-working': employee.status === 'Работает',
                'status-fired': employee.status === 'Уволен'
              }">
                {{ employee.status }}
              </span>
          </td>
          <td>
            <!-- Кнопка редактирования только для сотрудников со статусом "Работает" -->
            <button
                v-if="employee.status === 'Работает'"
                @click="openEditEmployeeModal(employee)"
                class="btn"
            >
              Редактировать
            </button>
            <!-- Кнопка "Уволить" только для сотрудников со статусом "Работает" -->
            <button
                v-if="employee.status === 'Работает'"
                @click="fireEmployee(employee.id)"
                class="btn danger"
            >
              Уволить
            </button>
          </td>
        </tr>
        </tbody>
      </table>

      <button
          v-if="canLoadMoreEmployees"
          @click="loadMoreEmployees"
          class="btn show-more-btn"
      >
        Показать ещё
      </button>
    </div>

    <div v-if="showEditEmployeeModal" class="modal-overlay">
      <div class="modal">
        <h3>Редактировать сотрудника</h3>

        <form @submit.prevent="saveEmployee">
          <label for="edit-name">ФИО</label>
          <input
              id="edit-name"
              type="text"
              v-model="currentEmployee.name"
              required
              maxlength="80"
          />

          <label for="edit-position">Должность</label>
          <select
              class="position"
              id="edit-position"
              v-model="currentEmployee.position"
              required
          >
            <option disabled value="">Выберите должность</option>
            <option>Фотограф</option>
            <option>Администратор</option>
          </select>

          <label for="edit-hireDate">Дата оформления</label>
          <input
              id="edit-hireDate"
              type="date"
              v-model="currentEmployee.hireDate"
              :max="todayDate"
              required
          />

          <label for="edit-birthDate">День рождения</label>
          <input
              id="edit-birthDate"
              type="date"
              v-model="currentEmployee.birthDate"
              :max="minBirthDate"
              required
          />

          <label for="edit-passport">Паспортные данные</label>
          <input
              id="edit-passport"
              type="text"
              v-model="formattedEditPassportData"
              required
              maxlength="10"
          />

          <label for="edit-phone">Номер телефона</label>
          <input
              id="edit-phone"
              type="tel"
              v-model="formattedEditPhoneNumber"
              required
              maxlength="12"
          />

          <button type="submit" class="btn">Сохранить</button>
          <button type="button" @click="closeEditEmployeeModal" class="btn danger">Отмена</button>
        </form>
      </div>
    </div>

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
      newEmployee: {
        full_name: '',
        position: '',
        hire_date: '',
        birth_date: '',
        passport_data: '',
        phone_number: '',
        status: 'Работает' // Добавлено поле 'status' с дефолтным значением
      },
      employees: [],
      loadingEmployees: true,
      showEditEmployeeModal: false,
      currentEmployee: {
        id: '',
        name: '',
        position: '',
        hireDate: '',
        birthDate: '',
        passport_data: '',
        phone_number: '',
        status: ''
      },
      displayCountEmployees: 10,
      sortKey: '', // Ключ для сортировки ('name' или 'status')
      sortOrder: 'asc', // Порядок сортировки ('asc' или 'desc')
    };
  },
  computed: {
    sortedEmployees() {
      if (!this.sortKey) {
        return this.employees;
      }

      return [...this.employees].sort((a, b) => {
        let aKey = a[this.sortKey];
        let bKey = b[this.sortKey];

        // Приведение к одному регистру для строк
        if (typeof aKey === 'string') aKey = aKey.toLowerCase();
        if (typeof bKey === 'string') bKey = bKey.toLowerCase();

        if (aKey < bKey) return this.sortOrder === 'asc' ? -1 : 1;
        if (aKey > bKey) return this.sortOrder === 'asc' ? 1 : -1;
        return 0;
      });
    },
    displayedEmployees() {
      return this.sortedEmployees.slice(0, this.displayCountEmployees);
    },
    canLoadMoreEmployees() {
      return this.displayCountEmployees < this.employees.length;
    },
    minBirthDate() {
      const today = new Date();
      today.setFullYear(today.getFullYear() - 18);
      return today.toISOString().split('T')[0];
    },
    todayDate() {
      const today = new Date();
      const year = today.getFullYear();
      const month = String(today.getMonth() + 1).padStart(2, '0');
      const day = String(today.getDate()).padStart(2, '0');
      return `${year}-${month}-${day}`;
    },
    // Вычисляемое свойство для форматирования номера телефона при добавлении
    formattedPhoneNumber: {
      get() {
        const raw = this.newEmployee.phone_number;
        if (!raw) return '';
        // Формат: "7-910-157-88-91"
        return raw.replace(/(\d)(\d{3})(\d{3})(\d{2})(\d{2})/, '$1-$2-$3-$4-$5');
      },
      set(value) {
        // Удаляем все нецифровые символы
        this.newEmployee.phone_number = value.replace(/\D/g, '').slice(0, 11);
      }
    },
    // Вычисляемое свойство для форматирования паспортных данных при добавлении
    formattedPassportData: {
      get() {
        const raw = this.newEmployee.passport_data;
        if (!raw) return '';
        // Формат: "1234 567890"
        return raw.replace(/(\d{4})(\d{6})/, '$1 $2');
      },
      set(value) {
        // Удаляем все нецифровые символы
        this.newEmployee.passport_data = value.replace(/\D/g, '').slice(0, 10);
      }
    },
    // Вычисляемое свойство для форматирования номера телефона при редактировании
    formattedEditPhoneNumber: {
      get() {
        const raw = this.currentEmployee.phone_number;
        if (!raw) return '';
        return raw.replace(/(\d)(\d{3})(\d{3})(\d{2})(\d{2})/, '$1-$2-$3-$4-$5');
      },
      set(value) {
        this.currentEmployee.phone_number = value.replace(/\D/g, '').slice(0, 11);
      }
    },
    // Вычисляемое свойство для форматирования паспортных данных при редактировании
    formattedEditPassportData: {
      get() {
        const raw = this.currentEmployee.passport_data;
        if (!raw) return '';
        return raw.replace(/(\d{4})(\d{6})/, '$1 $2');
      },
      set(value) {
        this.currentEmployee.passport_data = value.replace(/\D/g, '').slice(0, 10);
      }
    }
  },
  methods: {
    // Метод для сортировки по выбранному ключу
    sortBy(key) {
      if (this.sortKey === key) {
        // Если уже сортируем по этому ключу, переключаем порядок
        this.sortOrder = this.sortOrder === 'asc' ? 'desc' : 'asc';
      } else {
        // Иначе устанавливаем новый ключ сортировки и порядок по возрастанию
        this.sortKey = key;
        this.sortOrder = 'asc';
      }
    },
    // Вспомогательная функция для форматирования даты в RFC3339
    formatDateForInput(dateString) {
      if (!dateString) return '';
      const date = new Date(dateString);
      if (isNaN(date.getTime())) return ''; // Проверка на валидность даты
      return date.toISOString(); // Возвращает дату в формате RFC 3339
    },

    // Метод для получения списка сотрудников
    async fetchEmployees() {
      try {
        const response = await fetch('http://localhost:8080/api/employees', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json'
          },
          credentials: 'include'
        });

        if (!response.ok) {
          throw new Error('Ошибка при получении списка сотрудников');
        }

        const data = await response.json();

        // Измененный маппинг: используем employee_id вместо id и сохраняем даты в "YYYY-MM-DD"
        this.employees = data.employees.map(employee => ({
          id: employee.employee_id,
          name: employee.full_name,
          position: employee.position,
          hireDate: employee.hire_date, // "YYYY-MM-DD"
          birthDate: employee.birth_date, // "YYYY-MM-DD"
          passport_data: employee.passport_data,
          phone_number: employee.phone_number,
          status: employee.status
        }));

      } catch (error) {
        console.error('Ошибка при загрузке сотрудников:', error.message);
        alert('Не удалось загрузить сотрудников.');
      } finally {
        this.loadingEmployees = false;
      }
    },

    // Метод для добавления нового сотрудника
    async addEmployee() {
      const { full_name, position, hire_date, birth_date, passport_data, phone_number } = this.newEmployee;

      // Валидация возраста
      const birthDateObj = new Date(birth_date);
      const today = new Date();
      today.setFullYear(today.getFullYear() - 18);
      if (birthDateObj > today) {
        alert('Сотрудник должен быть старше 18 лет.');
        return;
      }

      // Форматируем данные для отправки на сервер в RFC3339
      const formattedEmployee = {
        full_name: full_name,
        position: position,
        hire_date: this.formatDateForInput(hire_date), // Формат RFC3339
        birth_date: this.formatDateForInput(birth_date), // Формат RFC3339
        passport_data: passport_data.replace(/\s/g, ''),
        phone_number: phone_number.replace(/-/g, ''),
        status: 'Работает' // Добавляем статус по умолчанию
      };

      console.log('Отправляемые данные:', formattedEmployee);

      try {
        const response = await fetch('http://localhost:8080/api/employees', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include',
          body: JSON.stringify(formattedEmployee),
        });

        console.log('Ответ сервера (добавление):', response);

        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.message || 'Ошибка при добавлении сотрудника');
        }

        const addedEmployee = await response.json();

        // Проверяем, что addedEmployee содержит данные о сотруднике
        if (!addedEmployee || !addedEmployee.worker) {
          throw new Error('Некорректный ответ от сервера при добавлении сотрудника.');
        }

        console.log('Добавленный сотрудник:', addedEmployee.worker);

        // Добавляем нового сотрудника в список
        this.employees.push({
          id: addedEmployee.worker.id,
          name: addedEmployee.worker.full_name,
          position: addedEmployee.worker.position,
          hireDate: addedEmployee.worker.hire_date, // "YYYY-MM-DD"
          birthDate: addedEmployee.worker.birth_date, // "YYYY-MM-DD"
          passport_data: addedEmployee.worker.passport_data,
          phone_number: addedEmployee.worker.phone_number,
          status: addedEmployee.worker.status
        });

        // Уведомление об успешном добавлении
        alert('Сотрудник успешно добавлен!');

        // Сброс формы
        this.newEmployee = {
          full_name: '',
          position: '',
          hire_date: '',
          birth_date: '',
          passport_data: '',
          phone_number: '',
          status: 'Работает' // Сброс с дефолтным значением
        };
      } catch (error) {
        console.error('Ошибка при добавлении сотрудника:', error.message);
        alert(`Не удалось добавить сотрудника: ${error.message}`);
      }
    },

    // Метод для увольнения сотрудника
    async fireEmployee(id) {
      if (!id) {
        alert('Невозможно уволить сотрудника: ID отсутствует.');
        return;
      }

      const employee = this.employees.find(emp => emp.id === id);
      if (!employee) {
        alert('Сотрудник не найден.');
        return;
      }

      if (!confirm(`Вы уверены, что хотите уволить сотрудника ${employee.name}?`)) {
        return;
      }

      try {
        const response = await fetch(`http://localhost:8080/api/employees/${id}/fire`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json'
          },
          credentials: 'include'
        });

        console.log('Ответ сервера (увольнение):', response);

        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.message || 'Ошибка при увольнении сотрудника');
        }

        // Обновление статуса сотрудника в локальном списке
        this.employees = this.employees.map(emp => {
          if (emp.id === id) {
            return { ...emp, status: 'Уволен' };
          }
          return emp;
        });

        alert(`Сотрудник ${employee.name} успешно уволен.`);
      } catch (error) {
        console.error('Ошибка при увольнении сотрудника:', error.message);
        alert(`Не удалось уволить сотрудника: ${error.message}`);
      }
    },

    // Метод для открытия модального окна редактирования сотрудника
    openEditEmployeeModal(employee) {
      // Копируем данные сотрудника с преобразованием формата даты
      this.currentEmployee = {
        ...employee,
        hireDate: this.formatDateForDateInput(employee.hireDate),
        birthDate: this.formatDateForDateInput(employee.birthDate)
      };

      this.showEditEmployeeModal = true;
    },
    // В методах или computed свойствах
    formatDateForDateInput(dateString) {
      if (!dateString) return '';
      const date = new Date(dateString);
      if (isNaN(date.getTime())) return ''; // Проверка на валидность даты
      return date.toISOString().split('T')[0]; // Возвращает "YYYY-MM-DD"
    },



    // Метод для закрытия модального окна редактирования сотрудника
    closeEditEmployeeModal() {
      this.showEditEmployeeModal = false;
      this.currentEmployee = {
        id: '',
        name: '',
        position: '',
        hireDate: '',
        birthDate: '',
        passport_data: '',
        phone_number: '',
        status: ''
      };
    },

    // Метод для сохранения изменений сотрудника
    async saveEmployee() {
      // Валидация: убедимся, что сотрудник старше 18 лет
      const { name, position, hireDate, birthDate, passport_data, phone_number, status } = this.currentEmployee;

      const birthDateObj = new Date(birthDate);
      const today = new Date();
      today.setFullYear(today.getFullYear() - 18);
      if (birthDateObj > today) {
        alert('Сотрудник должен быть старше 18 лет.');
        return;
      }

      // Форматируем данные для отправки на сервер в RFC3339
      const formattedEmployee = {
        full_name: name,
        position: position,
        hire_date: this.formatDateForInput(hireDate), // Формат RFC3339
        birth_date: this.formatDateForInput(birthDate), // Формат RFC3339
        passport_data: passport_data.replace(/\s/g, ''),
        phone_number: phone_number.replace(/-/g, ''),
        status: status
      };

      console.log('Отправляемые данные (редактирование):', formattedEmployee);

      try {
        const response = await fetch(`http://localhost:8080/api/employees/${this.currentEmployee.id}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json'
          },
          credentials: 'include',
          body: JSON.stringify(formattedEmployee)
        });

        console.log('Ответ сервера (редактирование):', response);

        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.message || 'Ошибка при обновлении сотрудника');
        }

        const updatedEmployee = await response.json();

        console.log('Обновленный сотрудник:', updatedEmployee.worker);

        // Обновляем данные сотрудника в локальном списке
        const index = this.employees.findIndex(emp => emp.id === updatedEmployee.worker.id);
        if (index !== -1) {
          this.employees[index] = {
            id: updatedEmployee.worker.id,
            name: updatedEmployee.worker.full_name,
            position: updatedEmployee.worker.position,
            hireDate: updatedEmployee.worker.hire_date, // "YYYY-MM-DD"
            birthDate: updatedEmployee.worker.birth_date, // "YYYY-MM-DD"
            passport_data: updatedEmployee.worker.passport_data,
            phone_number: updatedEmployee.worker.phone_number,
            status: updatedEmployee.worker.status
          };
        }

        alert(`Данные сотрудника ${updatedEmployee.worker.full_name} успешно обновлены.`);
        this.closeEditEmployeeModal();
      } catch (error) {
        console.error('Ошибка при обновлении сотрудника:', error.message);
        alert(`Не удалось обновить сотрудника: ${error.message}`);
      }
    },

    // Метод для форматирования даты для отображения пользователю
    formatDate(dateString) {
      const options = { year: 'numeric', month: 'long', day: 'numeric' };
      const date = new Date(dateString);
      return date.toLocaleDateString('ru-RU', options);
    },

    // Метод для форматирования номера телефона для отображения в таблице
    formatPhoneNumber(rawPhone) {
      if (!rawPhone) return '';
      return rawPhone.replace(/(\d)(\d{3})(\d{3})(\d{2})(\d{2})/, '$1-$2-$3-$4-$5');
    },

    // Метод для загрузки дополнительных сотрудников
    loadMoreEmployees() {
      this.displayCountEmployees += 10;
    },

    // Методы навигации по маршрутам
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
  },
  mounted() {
    this.fetchEmployees();
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
.position {
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

.btn:hover {
  background-color: #45a049;
}

.status-working {
  color: green;
  font-weight: bold;
}

.status-fired {
  color: red;
  font-weight: bold;
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
</style>
