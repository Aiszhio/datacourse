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
        </tr>
        </thead>
        <tbody>
        <tr v-for="client in paginatedClients" :key="client.id">
          <td>{{ client.FullName }}</td>
          <td>{{ formatPhone(client.PhoneNumber) }}</td>
          <td>{{ client.Email }}</td>
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
    <div v-if="showCreateClientModal" class="modal-overlay">
      <div class="modal">
        <h3>Создать клиента</h3>
        <form @submit.prevent="createClient">
          <label>
            ФИО:
            <input
                type="text"
                v-model="newClient.FullName"
                placeholder="Введите ФИО"
                required
                maxlength="80"
            />
          </label>
          <label>
            Номер телефона:
            <input
                type="text"
                v-model="newClient.PhoneNumber"
                placeholder="Введите номер телефона"
                required
                maxlength="11"
                @input="formatPhoneInput"
            />
          </label>
          <label>
            Почта:
            <input
                type="email"
                v-model="newClient.Email"
                placeholder="Введите почту"
                required
                maxlength="100"
            />
          </label>
          <div class="modal-actions">
            <button type="submit" class="btn save">Создать</button>
            <button type="button" @click="closeCreateClientModal" class="btn cancel">Отмена</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'AdminClientsDashboard',
  data() {
    return {
      clients: [],
      showCreateClientModal: false,
      newClient: {
        FullName: '',
        PhoneNumber: '',
        Email: ''
      },
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
          aVal = Number(aVal);
          bVal = Number(bVal);
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
          alert('Ошибка загрузки клиентов: Некорректный формат данных.');
        }
      } catch (error) {
        console.error('Ошибка при загрузке клиентов:', error);
        alert('Не удалось загрузить клиентов.');
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
    formatPhone(phone) {
      if (!phone) return phone;
      if (phone.length === 11) {
        return `+7 (${phone.substring(1, 4)}) ${phone.substring(4, 7)}-${phone.substring(7, 9)}-${phone.substring(9, 11)}`;
      }
      return phone;
    },
    async createClient() {
      if (!/^7\d{10}$/.test(this.newClient.PhoneNumber)) {
        alert('Номер телефона должен начинаться с 7 и содержать 11 цифр.');
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
          alert(response.data.message);
          this.closeCreateClientModal();
          this.fetchClients();
        } else if (response.data.error) {
          alert(response.data.error);
        }
      } catch (error) {
        console.error('Ошибка при создании клиента:', error);
        if (error.response && error.response.data) {
          alert(error.response.data.error || error.response.data.message || 'Неизвестная ошибка.');
        } else {
          alert('Не удалось создать клиента.');
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
    },
    resetFilters() {
      this.filters = {
        name: '',
        phone: '',
        email: ''
      };
      this.currentPage = 1;
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

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal {
  background-color: white;
  padding: 25px;
  border-radius: 8px;
  width: 400px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.modal h3 {
  margin-bottom: 15px;
  color: #333;
}

.modal label {
  display: block;
  margin-bottom: 10px;
  color: #555;
}

.modal input[type="text"],
.modal input[type="email"],
.modal input[type="date"],
.modal select {
  width: 100%;
  padding: 12px;
  margin-top: 5px;
  box-sizing: border-box;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 16px;
  transition: border-color 0.3s ease-in-out;
}

.modal input[type="text"]:focus,
.modal input[type="email"]:focus,
.modal input[type="date"]:focus,
.modal select:focus {
  border-color: #0056b3;
  outline: none;
}

.modal .modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 15px;
}

.modal .btn.save {
  background-color: #4CAF50;
  color: white;
}

.modal .btn.save:hover {
  background-color: #45a049;
}

.modal .btn.cancel {
  background-color: #F44336;
  color: white;
}

.modal .btn.cancel:hover {
  background-color: #D32F2F;
}
</style>
