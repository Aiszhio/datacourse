<template>
  <div class="manage-employees">
    <h2>Управление сотрудниками</h2>

    <!-- Секция добавления нового сотрудника -->
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
            :max="maxBirthDate"
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

    <!-- Секция списка сотрудников -->
    <div class="employee-list-section">
      <h3>Список сотрудников</h3>

      <div v-if="loadingEmployees" class="loading">Загрузка сотрудников...</div>
      <div v-else-if="employees.length === 0" class="no-employees">Нет сотрудников</div>

      <table v-else class="employees-table">
        <thead>
        <tr>
          <th @click="sortBy('name')" class="sortable">
            ФИО
            <span v-if="sortKey === 'name'">
                <span v-if="sortOrder === 'asc'">&#9650;</span>
                <span v-else>&#9660;</span>
              </span>
          </th>
          <th>Должность</th>
          <th @click="sortBy('hireDate')" class="sortable">
            Дата оформления
            <span v-if="sortKey === 'hireDate'">
                <span v-if="sortOrder === 'asc'">&#9650;</span>
                <span v-else>&#9660;</span>
              </span>
          </th>
          <th @click="sortBy('birthDate')" class="sortable">
            День рождения
            <span v-if="sortKey === 'birthDate'">
                <span v-if="sortOrder === 'asc'">&#9650;</span>
                <span v-else>&#9660;</span>
              </span>
          </th>
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
        <tr v-for="(employee) in displayedEmployees" :key="employee.id">
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

    <!-- Модальное окно редактирования сотрудника с использованием BootstrapVue3 -->
    <b-modal
        id="edit-employee-modal"
        v-model="showEditEmployeeModal"
        title="Редактировать сотрудника"
        hide-footer
        @hide="resetCurrentEmployee"
    >
      <form @submit.prevent="saveEmployee">
        <label for="edit-name">ФИО</label>
        <input
            id="edit-name"
            type="text"
            v-model="currentEmployee.name"
            required
            maxlength="80"
            class="form-control"
        />

        <label for="edit-position">Должность</label>
        <select
            class="position form-control"
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
            class="form-control"
        />

        <label for="edit-birthDate">День рождения</label>
        <input
            id="edit-birthDate"
            type="date"
            v-model="currentEmployee.birthDate"
            :max="maxBirthDate"
            required
            class="form-control"
        />

        <label for="edit-passport">Паспортные данные</label>
        <input
            id="edit-passport"
            type="text"
            v-model="currentEmployee.passport_data"
            required
            maxlength="10"
            class="form-control"
        />

        <label for="edit-phone">Номер телефона</label>
        <input
            id="edit-phone"
            type="tel"
            v-model="currentEmployee.phone_number"
            required
            maxlength="12"
            class="form-control"
        />

        <label for="edit-status">Статус</label>
        <select
            id="edit-status"
            v-model="currentEmployee.status"
            required
            class="form-control"
        >
          <option value="Работает">Работает</option>
          <option value="Уволен">Уволен</option>
        </select>

        <div class="mt-3 d-flex justify-content-end">
          <button type="submit" class="btn btn-primary mr-2">Сохранить</button>
          <button type="button" @click="closeEditEmployeeModal" class="btn btn-danger">Отмена</button>
        </div>
      </form>
    </b-modal>

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
      <div class="card" @click="goToClients">
        <i class="fas fa-box icon"></i>
        <h3>Клиенты</h3>
        <p>Учёт и управление</p>
      </div>
    </div>
  </div>
</template>

<script>
import dayjs from 'dayjs';
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';
import localizedFormat from 'dayjs/plugin/localizedFormat';
import 'dayjs/locale/ru';

dayjs.extend(utc);
dayjs.extend(timezone);
dayjs.extend(localizedFormat);
dayjs.locale('ru'); // Устанавливаем русскую локаль

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
      sortKey: '', // Ключ для сортировки ('name', 'hireDate', 'birthDate', 'status')
      sortOrder: 'asc', // Порядок сортировки ('asc' или 'desc')
    };
  },
  computed: {
    // Отсортированные сотрудники
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
    maxBirthDate() {
      return dayjs().subtract(18, 'year').format('YYYY-MM-DD'); // Максимальная дата рождения (старше 18 лет)
    },
    todayDate() {
      return dayjs().format('YYYY-MM-DD');
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
        this.newEmployee.passport_data = value.replace(/\D/g, '').slice(0, 11);
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
    // Вспомогательная функция для форматирования даты в ISO8601
    formatDateForInput(dateString) {
      if (!dateString) return '';
      const date = dayjs(dateString);
      if (!date.isValid()) return ''; // Проверка на валидность даты
      return date.toISOString(); // Возвращает "YYYY-MM-DDTHH:mm:ssZ"
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

        // Изменённый маппинг: используем employee_id вместо id и сохраняем даты в формат для input[type="date"]
        this.employees = data.employees.map(employee => ({
          id: employee.employee_id,
          name: employee.full_name,
          position: employee.position,
          hireDate: dayjs(employee.hire_date).format('YYYY-MM-DD'), // Формат для input[type="date"]
          birthDate: dayjs(employee.birth_date).format('YYYY-MM-DD'), // Формат для input[type="date"]
          passport_data: employee.passport_data,
          phone_number: employee.phone_number,
          status: employee.status
        }));

      } catch (error) {
        console.error('Ошибка при загрузке сотрудников:', error.message);
        this.showAlert('Не удалось загрузить сотрудников.', 'error');
      } finally {
        this.loadingEmployees = false;
      }
    },

    // Метод для добавления нового сотрудника
    async addEmployee() {
      const { full_name, position, hire_date, birth_date, passport_data, phone_number } = this.newEmployee;

      // Валидация возраста
      const birthDateObj = dayjs(birth_date);
      const today = dayjs().subtract(18, 'year');
      if (birthDateObj.isAfter(today)) {
        this.showAlert('Сотрудник должен быть старше 18 лет.', 'warning');
        return;
      }

      // Форматируем данные для отправки на сервер в ISO8601
      const formattedEmployee = {
        full_name: full_name,
        position: position,
        hire_date: this.formatDateForInput(hire_date), // ISO8601
        birth_date: this.formatDateForInput(birth_date), // ISO8601
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
          id: addedEmployee.worker.employee_id,
          name: addedEmployee.worker.full_name,
          position: addedEmployee.worker.position,
          hireDate: dayjs(addedEmployee.worker.hire_date).format('YYYY-MM-DD'), // Формат для input[type="date"]
          birthDate: dayjs(addedEmployee.worker.birth_date).format('YYYY-MM-DD'), // Формат для input[type="date"]
          passport_data: addedEmployee.worker.passport_data,
          phone_number: addedEmployee.worker.phone_number,
          status: addedEmployee.worker.status
        });

        // Уведомление об успешном добавлении
        this.showAlert('Сотрудник успешно добавлен!', 'success');

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
        this.showAlert(`Не удалось добавить сотрудника: ${error.message}`, 'error');
      }
    },

    // Метод для увольнения сотрудника
    async fireEmployee(id) {
      if (!id) {
        this.showAlert('Невозможно уволить сотрудника: ID отсутствует.', 'error');
        return;
      }

      const employee = this.employees.find(emp => emp.id === id);
      if (!employee) {
        this.showAlert('Сотрудник не найден.', 'error');
        return;
      }

      if (employee.status === 'Уволен') {
        this.showAlert('Сотрудник уже уволен.', 'info');
        return;
      }

      // Используем SweetAlert2 для подтверждения действия
      const result = await this.$swal.fire({
        title: 'Вы уверены?',
        text: `Вы хотите уволить сотрудника ${employee.name}?`,
        icon: 'warning',
        showCancelButton: true,
        confirmButtonText: 'Да, уволить',
        cancelButtonText: 'Отмена',
        reverseButtons: true
      });

      if (result.isConfirmed) {
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

          this.showAlert(`Сотрудник ${employee.name} успешно уволен.`, 'success');
        } catch (error) {
          console.error('Ошибка при увольнении сотрудника:', error.message);
          this.showAlert(`Не удалось уволить сотрудника: ${error.message}`, 'error');
        }
      } else if (result.dismiss === this.$swal.DismissReason.cancel) {
        this.showAlert('Увольнение отменено.', 'info');
      }
    },

    // Метод для открытия модального окна редактирования сотрудника
    openEditEmployeeModal(employee) {
      console.log('openEditEmployeeModal вызван с сотрудником:', employee);
      if (!employee) {
        this.showAlert('Невозможно открыть модальное окно: данные сотрудника отсутствуют.', 'error');
        return;
      }

      // Копируем данные сотрудника с преобразованием формата даты
      this.currentEmployee = {
        ...employee,
        hireDate: this.formatDateForDateInput(employee.hireDate),
        birthDate: this.formatDateForDateInput(employee.birthDate)
      };

      console.log('currentEmployee установлен:', this.currentEmployee);

      this.showEditEmployeeModal = true;
    },

    // Метод для закрытия модального окна редактирования сотрудника
    closeEditEmployeeModal() {
      this.showEditEmployeeModal = false;
      this.resetCurrentEmployee();
    },

    // Метод для сброса данных текущего сотрудника
    resetCurrentEmployee() {
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

      const birthDateObj = dayjs(birthDate);
      const today = dayjs().subtract(18, 'year');
      if (birthDateObj.isAfter(today)) {
        this.showAlert('Сотрудник должен быть старше 18 лет.', 'warning');
        return;
      }

      // Форматируем данные для отправки на сервер в ISO8601
      const formattedEmployee = {
        full_name: name,
        position: position,
        hire_date: this.formatDateForInput(hireDate), // ISO8601
        birth_date: this.formatDateForInput(birthDate), // ISO8601
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
        const index = this.employees.findIndex(emp => emp.id === updatedEmployee.worker.employee_id);
        if (index !== -1) {
          this.employees[index] = {
            id: updatedEmployee.worker.employee_id,
            name: updatedEmployee.worker.full_name,
            position: updatedEmployee.worker.position,
            hireDate: dayjs(updatedEmployee.worker.hire_date).format('YYYY-MM-DD'), // Формат для input[type="date"]
            birthDate: dayjs(updatedEmployee.worker.birth_date).format('YYYY-MM-DD'), // Формат для input[type="date"]
            passport_data: updatedEmployee.worker.passport_data,
            phone_number: updatedEmployee.worker.phone_number,
            status: updatedEmployee.worker.status
          };
        }

        this.showAlert(`Данные сотрудника ${updatedEmployee.worker.full_name} успешно обновлены.`, 'success');
        this.closeEditEmployeeModal();
      } catch (error) {
        console.error('Ошибка при обновлении сотрудника:', error.message);
        this.showAlert(`Не удалось обновить сотрудника: ${error.message}`, 'error');
      }
    },

    // Метод для форматирования даты для отображения в input[type="date"]
    formatDateForDateInput(dateString) {
      if (!dateString) return '';
      const date = dayjs(dateString);
      if (!date.isValid()) return ''; // Проверка на валидность даты
      return date.format('YYYY-MM-DD'); // Возвращает "YYYY-MM-DD"
    },

    // Метод для форматирования даты для отображения пользователю
    formatDate(dateString) {
      if (!dateString) return '';
      return dayjs(dateString).format('D MMMM YYYY'); // Формат: "8 января 2024"
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

    // Вспомогательный метод для отображения уведомлений
    showAlert(message, type) {
      this.$swal.fire({
        title:
            type === 'error'
                ? 'Ошибка!'
                : type === 'success'
                    ? 'Успех!'
                    : type === 'warning'
                        ? 'Предупреждение!'
                        : 'Внимание!',
        text: message,
        icon: type,
        timer: 3000,
        showConfirmButton: false,
        toast: true,
        position: 'top-end',
      });
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
    },
    goToClients(){
      this.$router.push({ name: 'Clients' })
    },
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
  /* Убрали max-height: 100vh, чтобы модальное окно не было ограничено высотой контейнера */
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

.sortable {
  cursor: pointer;
  user-select: none;
}

.status-working {
  color: green;
  font-weight: bold;
}

.status-fired {
  color: red;
  font-weight: bold;
}

.btn {
  background-color: #4CAF50;
  color: white;
  padding: 10px 15px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 1em;
  margin: 10px 5px 0 0;
  transition: background-color 0.3s ease;
}

.btn.danger {
  background-color: #f44336;
}

.btn.show-more-btn {
  margin-top: 15px;
}

.btn:hover {
  background-color: #45a049;
}

.btn:disabled {
  background-color: #a5d6a7;
  cursor: not-allowed;
}

/* Удаляем кастомные стили модального окна, т.к. используем BootstrapVue3 */
.modal-overlay {
  /* Удалены или закомментированы стили */
}

.modal {
  /* Удалены или закомментированы стили */
}

.modal-buttons {
  display: flex;
  justify-content: space-between;
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
  padding: 20px;
  margin-bottom: 10px;
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

.loading {
  text-align: center;
  color: #555555;
  font-size: 1.2em;
  margin-top: 10px;
}

.no-employees {
  text-align: center;
  color: #888888;
  font-style: italic;
  margin-top: 10px;
}

@media (max-width: 768px) {
  .filters-container {
    flex-direction: column;
  }
}
</style>
