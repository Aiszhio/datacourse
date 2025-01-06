<template>
  <div class="admin-clients-dashboard">
    <h2>Управление клиентами</h2>

    <!-- Раздел фильтрации клиентов -->
    <div class="filter-wrapper">
      <label>
        ФИО:
        <input
            type="text"
            v-model="filters.name"
            placeholder="Поиск по ФИО"
            class="filter-input"
        />
      </label>
      <label>
        Номер телефона:
        <input
            type="text"
            v-model="filters.phone"
            placeholder="Поиск по телефону"
            class="filter-input"
        />
      </label>
      <label>
        Почта:
        <input
            type="text"
            v-model="filters.email"
            placeholder="Поиск по почте"
            class="filter-input"
        />
      </label>
      <button @click="resetFilters" class="btn filter">Сбросить</button>
    </div>

    <!-- Кнопка для добавления нового клиента -->
    <div class="create-client-button">
      <button @click="openCreateClientModal" class="btn create">Добавить клиента</button>
    </div>

    <!-- Таблица с клиентами -->
    <div class="table-section">
      <h3>Список клиентов</h3>
      <table class="data-table">
        <thead>
        <tr>
          <th @click="sortBy('FullName')">
            ФИО
            <span v-if="currentSortKey === 'FullName'">
                {{ sortOrders.FullName === 'asc' ? '▲' : '▼' }}
              </span>
          </th>
          <th @click="sortBy('phone')">
            Номер телефона
            <span v-if="currentSortKey === 'phone'">
                {{ sortOrders.phone === 'asc' ? '▲' : '▼' }}
              </span>
          </th>
          <th @click="sortBy('email')">
            Почта
            <span v-if="currentSortKey === 'email'">
                {{ sortOrders.email === 'asc' ? '▲' : '▼' }}
              </span>
          </th>
          <th>Действия</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="client in paginatedClients" :key="client.id">
          <td>{{ client.FullName }}</td>
          <td>{{ formatPhone(client.PhoneNumber) }}</td>
          <td>{{ client.Email }}</td>
          <td>
            <button @click="openEditClientModal(client)" class="btn edit">Редактировать</button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Пагинация -->
    <div class="pagination" v-if="totalPages > 1">
      <button
          @click="prevPage"
          :disabled="currentPage === 1"
          class="btn pagination-btn"
      >
        Назад
      </button>
      <span>Страница {{ currentPage }} из {{ totalPages }}</span>
      <button
          @click="nextPage"
          :disabled="currentPage === totalPages"
          class="btn pagination-btn"
      >
        Вперед
      </button>
    </div>

    <!-- Панель навигации с карточками -->
    <div class="card-panel">
      <div class="card" @click="goToAdminHome">
        <i class="fas fa-home icon"></i>
        <h3>Главная</h3>
        <p>Панель администратора</p>
      </div>
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
      <div class="card" @click="goToServicesPage">
        <i class="fas fa-concierge-bell icon"></i>
        <h3>Услуги</h3>
        <p>Управление услугами</p>
      </div>
      <div class="card" @click="goToBookingsPage">
        <i class="fas fa-calendar-check icon"></i>
        <h3>Бронирование</h3>
        <p>Управление бронированием</p>
      </div>
    </div>

    <!-- Модальное окно для создания клиента -->
    <b-modal
        id="create-client-modal"
        v-model="showCreateClientModal"
        title="Создать клиента"
        hide-footer
        @hide="resetCreateClientForm"
    >
      <form @submit.prevent="createClient">
        <div class="mb-3">
          <label for="clientFullName" class="form-label">ФИО:</label>
          <input
              id="clientFullName"
              type="text"
              v-model="newClient.FullName"
              placeholder="Введите ФИО"
              required
              maxlength="80"
              class="form-control"
          />
        </div>
        <div class="mb-3">
          <label for="clientPhone" class="form-label">Номер телефона:</label>
          <input
              id="clientPhone"
              type="text"
              v-model="newClient.PhoneNumber"
              placeholder="Введите номер телефона"
              required
              maxlength="11"
              @input="formatPhoneInput"
              class="form-control"
          />
        </div>
        <div class="mb-3">
          <label for="clientEmail" class="form-label">Почта:</label>
          <input
              id="clientEmail"
              type="email"
              v-model="newClient.Email"
              placeholder="Введите почту"
              required
              maxlength="100"
              class="form-control"
          />
        </div>
        <div class="d-flex justify-content-end">
          <button type="submit" class="btn btn-primary me-2">Создать</button>
          <button type="button" @click="closeCreateClientModal" class="btn btn-secondary">Отмена</button>
        </div>
      </form>
    </b-modal>

    <!-- Модальное окно для редактирования клиента -->
    <b-modal
        id="edit-client-modal"
        v-model="showEditClientModal"
        title="Редактировать клиента"
        hide-footer
        @hide="resetEditClientForm"
    >
      <form @submit.prevent="updateClient">
        <div class="mb-3">
          <label for="editClientFullName" class="form-label">ФИО:</label>
          <input
              id="editClientFullName"
              type="text"
              v-model="editClient.FullName"
              placeholder="Введите ФИО"
              required
              maxlength="80"
              class="form-control"
          />
        </div>
        <div class="mb-3">
          <label for="editClientPhone" class="form-label">Номер телефона:</label>
          <input
              id="editClientPhone"
              type="text"
              v-model="editClient.PhoneNumber"
              placeholder="Введите номер телефона"
              required
              maxlength="11"
              @input="formatPhoneInputEdit"
              class="form-control"
          />
        </div>
        <div class="mb-3">
          <label for="editClientEmail" class="form-label">Почта:</label>
          <input
              id="editClientEmail"
              type="email"
              v-model="editClient.Email"
              placeholder="Введите почту"
              required
              maxlength="100"
              class="form-control"
          />
        </div>
        <div class="d-flex justify-content-end">
          <button type="submit" class="btn btn-primary me-2">Сохранить</button>
          <button type="button" @click="closeEditClientModal" class="btn btn-secondary">Отмена</button>
        </div>
      </form>
    </b-modal>

    <!-- Модальное окно подтверждения удаления клиента -->
    <b-modal
        id="delete-confirm-modal"
        v-model="showDeleteConfirmModal"
        title="Подтверждение удаления"
        hide-footer
        @hide="resetDeleteConfirmModal"
    >
      <p>Вы уверены, что хотите удалить клиента <strong>{{ deleteClient.FullName }}</strong>?</p>
      <div class="d-flex justify-content-end">
        <button @click="deleteClientConfirmed" class="btn btn-danger me-2">Удалить</button>
        <button @click="closeDeleteConfirmModal" class="btn btn-secondary">Отмена</button>
      </div>
    </b-modal>
  </div>
</template>

<script>
import axios from 'axios';
import Swal from 'sweetalert2'; // Импортируем SweetAlert2

export default {
  name: 'AdminClientsDashboard',
  data() {
    return {
      clients: [],
      showCreateClientModal: false,
      showEditClientModal: false,
      showDeleteConfirmModal: false,
      newClient: {
        FullName: '',
        PhoneNumber: '',
        Email: ''
      },
      editClient: {
        id: null,
        FullName: '',
        PhoneNumber: '',
        Email: ''
      },
      deleteClient: {},
      sortOrders: {
        FullName: 'asc',
        phone: 'asc',
        email: 'asc'
      },
      currentSortKey: '',
      filters: {
        name: '',
        phone: '',
        email: ''
      },
      currentPage: 1,
      pageSize: 10,
      totalPages: 1
    };
  },
  computed: {
    filteredClients() {
      return this.clients.filter(client => {
        const matchesName = this.filters.name
            ? client.FullName.toLowerCase().includes(this.filters.name.toLowerCase())
            : true;
        const matchesPhone = this.filters.phone
            ? client.PhoneNumber.includes(this.filters.phone)
            : true;
        const matchesEmail = this.filters.email
            ? client.Email.toLowerCase().includes(this.filters.email.toLowerCase())
            : true;
        return matchesName && matchesPhone && matchesEmail;
      });
    },
    sortedClients() {
      if (!this.currentSortKey) {
        return this.filteredClients;
      }
      return [...this.filteredClients].sort((a, b) => {
        let aVal = a[this.currentSortKey];
        let bVal = b[this.currentSortKey];

        if (this.currentSortKey === 'phone') {
          aVal = aVal.toString();
          bVal = bVal.toString();
        } else {
          aVal = aVal.toString().toLowerCase();
          bVal = bVal.toString().toLowerCase();
        }

        if (aVal < bVal) return this.sortOrders[this.currentSortKey] === 'asc' ? -1 : 1;
        if (aVal > bVal) return this.sortOrders[this.currentSortKey] === 'asc' ? 1 : -1;
        return 0;
      });
    },
    paginatedClients() {
      const start = (this.currentPage - 1) * this.pageSize;
      const end = start + this.pageSize;
      return this.sortedClients.slice(start, end);
    }
  },
  methods: {
    async fetchClients() {
      try {
        const response = await axios.get('http://localhost:8080/api/admin/clients', {
          withCredentials: true
        });
        console.log('Ответ от сервера:', response.data);

        if (response.data && Array.isArray(response.data.clients)) {
          this.clients = response.data.clients;
          this.calculateTotalPages();
        } else {
          console.error('Некорректный формат данных:', response.data);
          Swal.fire({
            title: 'Ошибка!',
            text: 'Ошибка загрузки клиентов: Некорректный формат данных.',
            icon: 'error',
            timer: 3000,
            showConfirmButton: false,
            toast: true,
            position: 'top-end',
          });
        }
      } catch (error) {
        console.error('Ошибка при загрузке клиентов:', error);
        Swal.fire({
          title: 'Ошибка!',
          text: 'Не удалось загрузить клиентов.',
          icon: 'error',
          timer: 3000,
          showConfirmButton: false,
          toast: true,
          position: 'top-end',
        });
      }
    },
    formatPhoneInput() {
      let phone = this.newClient.PhoneNumber.replace(/\D/g, '');
      phone = phone.substring(0, 11);
      if (phone.length > 0 && phone[0] !== '7') {
        phone = '7' + phone.substring(1);
      }
      this.newClient.PhoneNumber = phone;
    },
    formatPhoneInputEdit() {
      let phone = this.editClient.PhoneNumber.replace(/\D/g, '');
      phone = phone.substring(0, 11);
      if (phone.length > 0 && phone[0] !== '7') {
        phone = '7' + phone.substring(1);
      }
      this.editClient.PhoneNumber = phone;
    },
    formatPhone(phone) {
      if (!phone) return phone;
      if (phone.length === 11) {
        return `+7 (${phone.substring(1, 4)}) ${phone.substring(4, 7)}-${phone.substring(7, 9)}-${phone.substring(9, 11)}`;
      }
      return phone;
    },
    async createClient() {
      if (!/^7\d{10}$/.test(this.newClient.PhoneNumber)) {
        Swal.fire({
          title: 'Ошибка!',
          text: 'Номер телефона должен начинаться с 7 и содержать 11 цифр.',
          icon: 'error',
          timer: 3000,
          showConfirmButton: false,
          toast: true,
          position: 'top-end',
        });
        return;
      }

      try {
        const clientData = {
          FullName: this.newClient.FullName,
          PhoneNumber: this.newClient.PhoneNumber,
          Email: this.newClient.Email
        };
        const response = await axios.post('http://localhost:8080/api/admin/clients', clientData, {
          withCredentials: true
        });

        if (response.data.message) {
          Swal.fire({
            title: 'Успех!',
            text: response.data.message,
            icon: 'success',
            timer: 2000,
            showConfirmButton: false,
            toast: true,
            position: 'top-end',
          });
          this.closeCreateClientModal();
          this.fetchClients();
        } else if (response.data.error) {
          Swal.fire({
            title: 'Ошибка!',
            text: response.data.error,
            icon: 'error',
            timer: 3000,
            showConfirmButton: false,
            toast: true,
            position: 'top-end',
          });
        }
      } catch (error) {
        console.error('Ошибка при создании клиента:', error);
        if (error.response && error.response.data) {
          Swal.fire({
            title: 'Ошибка!',
            text: error.response.data.error || error.response.data.message || 'Неизвестная ошибка.',
            icon: 'error',
            timer: 3000,
            showConfirmButton: false,
            toast: true,
            position: 'top-end',
          });
        } else {
          Swal.fire({
            title: 'Ошибка!',
            text: 'Не удалось создать клиента.',
            icon: 'error',
            timer: 3000,
            showConfirmButton: false,
            toast: true,
            position: 'top-end',
          });
        }
      }
    },
    async openEditClientModal(client) {
      this.editClient = { ...client };
      this.showEditClientModal = true;
    },
    async updateClient() {
      if (!/^7\d{10}$/.test(this.editClient.PhoneNumber)) {
        Swal.fire({
          title: 'Ошибка!',
          text: 'Номер телефона должен начинаться с 7 и содержать 11 цифр.',
          icon: 'error',
          timer: 3000,
          showConfirmButton: false,
          toast: true,
          position: 'top-end',
        });
        return;
      }

      try {
        const clientData = {
          FullName: this.editClient.FullName,
          PhoneNumber: this.editClient.PhoneNumber,
          Email: this.editClient.Email
        };
        const response = await axios.put(`http://localhost:8080/api/admin/clients/${this.editClient.id}`, clientData, {
          withCredentials: true
        });

        if (response.data.message) {
          Swal.fire({
            title: 'Успех!',
            text: response.data.message,
            icon: 'success',
            timer: 2000,
            showConfirmButton: false,
            toast: true,
            position: 'top-end',
          });
          this.closeEditClientModal();
          this.fetchClients();
        } else if (response.data.error) {
          Swal.fire({
            title: 'Ошибка!',
            text: response.data.error,
            icon: 'error',
            timer: 3000,
            showConfirmButton: false,
            toast: true,
            position: 'top-end',
          });
        }
      } catch (error) {
        console.error('Ошибка при обновлении клиента:', error);
        if (error.response && error.response.data) {
          Swal.fire({
            title: 'Ошибка!',
            text: error.response.data.error || error.response.data.message || 'Неизвестная ошибка.',
            icon: 'error',
            timer: 3000,
            showConfirmButton: false,
            toast: true,
            position: 'top-end',
          });
        } else {
          Swal.fire({
            title: 'Ошибка!',
            text: 'Не удалось обновить клиента.',
            icon: 'error',
            timer: 3000,
            showConfirmButton: false,
            toast: true,
            position: 'top-end',
          });
        }
      }
    },
    confirmDeleteClient(client) {
      this.deleteClient = client;
      this.showDeleteConfirmModal = true;
    },
    async deleteClientConfirmed() {
      try {
        const response = await axios.delete(`http://localhost:8080/api/admin/clients/${this.deleteClient.id}`, {
          withCredentials: true
        });

        if (response.data.message) {
          Swal.fire({
            title: 'Успех!',
            text: response.data.message,
            icon: 'success',
            timer: 2000,
            showConfirmButton: false,
            toast: true,
            position: 'top-end',
          });
          this.closeDeleteConfirmModal();
          this.fetchClients();
        } else if (response.data.error) {
          Swal.fire({
            title: 'Ошибка!',
            text: response.data.error,
            icon: 'error',
            timer: 3000,
            showConfirmButton: false,
            toast: true,
            position: 'top-end',
          });
        }
      } catch (error) {
        console.error('Ошибка при удалении клиента:', error);
        if (error.response && error.response.data) {
          Swal.fire({
            title: 'Ошибка!',
            text: error.response.data.error || error.response.data.message || 'Неизвестная ошибка.',
            icon: 'error',
            timer: 3000,
            showConfirmButton: false,
            toast: true,
            position: 'top-end',
          });
        } else {
          Swal.fire({
            title: 'Ошибка!',
            text: 'Не удалось удалить клиента.',
            icon: 'error',
            timer: 3000,
            showConfirmButton: false,
            toast: true,
            position: 'top-end',
          });
        }
      }
    },
    sortBy(key) {
      if (this.currentSortKey === key) {
        this.sortOrders[key] = this.sortOrders[key] === 'asc' ? 'desc' : 'asc';
      } else {
        this.currentSortKey = key;
        this.sortOrders[key] = 'asc';
      }
      this.currentPage = 1;
      this.calculateTotalPages();
    },
    resetFilters() {
      this.filters = {
        name: '',
        phone: '',
        email: ''
      };
      this.currentPage = 1;
      this.calculateTotalPages();
    },
    calculateTotalPages() {
      this.totalPages = Math.ceil(this.sortedClients.length / this.pageSize) || 1;
      if (this.currentPage > this.totalPages) {
        this.currentPage = this.totalPages;
      }
    },
    nextPage() {
      if (this.currentPage < this.totalPages) {
        this.currentPage++;
      }
    },
    prevPage() {
      if (this.currentPage > 1) {
        this.currentPage--;
      }
    },
    openCreateClientModal() {
      this.showCreateClientModal = true;
      this.newClient = {
        FullName: '',
        PhoneNumber: '',
        Email: ''
      };
    },
    closeCreateClientModal() {
      this.showCreateClientModal = false;
      this.newClient = {
        FullName: '',
        PhoneNumber: '',
        Email: ''
      };
    },
    resetCreateClientForm() {
      this.newClient = {
        FullName: '',
        PhoneNumber: '',
        Email: ''
      };
    },
    openEditClientModal(client) {
      this.editClient = { ...client };
      this.showEditClientModal = true;
    },
    closeEditClientModal() {
      this.showEditClientModal = false;
      this.editClient = {
        id: null,
        FullName: '',
        PhoneNumber: '',
        Email: ''
      };
    },
    resetEditClientForm() {
      this.editClient = {
        id: null,
        FullName: '',
        PhoneNumber: '',
        Email: ''
      };
    },
    confirmDeleteClient(client) {
      this.deleteClient = client;
      this.showDeleteConfirmModal = true;
    },
    closeDeleteConfirmModal() {
      this.showDeleteConfirmModal = false;
      this.deleteClient = {};
    },
    goToAdminHome() {
      this.$router.push({ name: 'AdminHome' });
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
    goToBookingsPage() {
      this.$router.push({ name: 'Bookings' });
    }
  },
  watch: {
    filters: {
      handler() {
        this.currentPage = 1;
        this.calculateTotalPages();
      },
      deep: true
    },
    clients() {
      this.calculateTotalPages();
    }
  },
  mounted() {
    this.fetchClients();
  }
};
</script>

<style scoped>
.admin-clients-dashboard {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
  max-height: 100vh;
}

h2 {
  text-align: center;
  width: 100%;
  margin-bottom: 20px;
  color: #333;
}

.filter-wrapper {
  display: flex;
  gap: 15px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.filter-wrapper label {
  display: flex;
  flex-direction: column;
  font-weight: bold;
  color: #333;
}

.filter-wrapper .filter-input {
  width: 200px;
  padding: 12px;
  margin-top: 5px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 16px;
  transition: border-color 0.3s ease-in-out;
}

.filter-wrapper .filter-input:focus {
  border-color: #007bff;
  outline: none;
}

.filter-wrapper .filter-input:hover {
  border-color: #aaa;
}

.filter-wrapper .btn.filter {
  align-self: flex-end;
  background-color: #4CAF50;
  color: white;
  padding: 12px 20px;
  border-radius: 4px;
  border: none;
  cursor: pointer;
  transition: background-color 0.3s;
}

.filter-wrapper .btn.filter:hover {
  background-color: #45a049;
}

.create-client-button {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 20px;
}

.create-client-button .btn.create {
  background-color: #4CAF50;
  color: white;
}

.create-client-button .btn.create:hover {
  background-color: #45a049;
}

.table-section {
  margin-bottom: 20px;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 10px;
  border: 1px solid #ddd;
  text-align: left;
}

.data-table th {
  background-color: #f2f2f2;
  cursor: pointer;
  user-select: none;
}

.data-table th:hover {
  background-color: #e0e0e0;
}

.btn {
  padding: 8px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn.edit {
  background-color: #2196F3; /* Синий цвет для кнопки редактирования */
  color: white;
  margin-right: 5px;
}

.btn.edit:hover {
  background-color: #1976D2;
}

.btn.delete {
  background-color: #F44336; /* Красный цвет для кнопки удаления */
  color: white;
}

.btn.delete:hover {
  background-color: #D32F2F;
}

.btn.pagination-btn {
  background-color: #4CAF50; /* Зеленый цвет для кнопок пагинации */
  color: white;
}

.btn.pagination-btn:hover:not(:disabled) {
  background-color: #45a049;
}

.btn.pagination-btn:disabled {
  background-color: #a5d6a7;
  cursor: not-allowed;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 15px;
  margin-top: 15px;
}

.card-panel {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
  justify-content: center;
  margin-top: 30px;
}

.card {
  background-color: #4CAF50;
  color: #fff;
  border-radius: 10px;
  padding: 15px;
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

.card h3 {
  color: #FFFFE0;
  margin-bottom: 5px;
  font-size: 18px;
}

.icon {
  font-size: 40px;
  color: #FFFFE0;
  margin-bottom: 10px;
}

.form-label {
  display: block;
  font-weight: bold;
  margin-bottom: 5px;
  color: #333;
}

.form-control {
  width: 100%;
  padding: 8px 12px;
  margin-bottom: 15px;
  border: 1px solid #ddd;
  border-radius: 4px;
  transition: border-color 0.3s;
}

.form-control:focus {
  border-color: #4CAF50;
  outline: none;
}

.d-flex {
  display: flex;
}

.justify-content-end {
  justify-content: flex-end;
}
</style>
