<template>
  <div class="order-history">
    <h2>Заказы</h2>

    <!-- Кнопка для добавления нового заказа -->
    <div class="add-order-button">
      <button @click="openAddModal" class="btn add">Добавить заказ</button>
    </div>

    <!-- Раздел фильтрации заказов -->
    <div class="filter-section">
      <h3>Фильтровать заказы</h3>
      <div class="filters">
        <label>
          Имя клиента:
          <input
              v-model="searchClientName"
              placeholder="Поиск по имени клиента"
              class = "clientNameInput"
          />
        </label>
        <label>
          Название услуги:
          <select v-model="searchService">
            <option value="">Все услуги</option>
            <option
                v-for="service in uniqueServices"
                :key="service.name"
                :value="service.name"
            >
              {{ service.name }}
            </option>
          </select>
        </label>
      </div>
    </div>

    <!-- Раздел истории заказов -->
    <div v-if="!isLoading">
      <h3>Список заказов</h3>
      <table class = "data-table">
        <thead>
        <tr>
          <th @click="sortBy('clientName')">
            Имя клиента
            <span v-if="currentSortKey === 'clientName'">
                {{ sortOrders.clientName === 'asc' ? '▲' : '▼' }}
              </span>
          </th>
          <th @click="sortBy('employeeName')">
            Сотрудник
            <span v-if="currentSortKey === 'employeeName'">
                {{ sortOrders.employeeName === 'asc' ? '▲' : '▼' }}
              </span>
          </th>
          <th @click="sortBy('service')">
            Название услуги
            <span v-if="currentSortKey === 'service'">
                {{ sortOrders.service === 'asc' ? '▲' : '▼' }}
              </span>
          </th>
          <th @click="sortBy('orderDate')">
            Дата оформления
            <span v-if="currentSortKey === 'orderDate'">
                {{ sortOrders.orderDate === 'asc' ? '▲' : '▼' }}
              </span>
          </th>
          <th @click="sortBy('receiveDate')">
            Дата завершения
            <span v-if="currentSortKey === 'receiveDate'">
                {{ sortOrders.receiveDate === 'asc' ? '▲' : '▼' }}
              </span>
          </th>
          <th @click="sortBy('issueDate')">
            Дата выдачи
            <span v-if="currentSortKey === 'issueDate'">
                {{ sortOrders.issueDate === 'asc' ? '▲' : '▼' }}
              </span>
          </th>
          <th>Действие</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="order in paginatedOrders" :key="order.id">
          <td>{{ order.clientName }}</td>
          <td>{{ order.employeeName }}</td>
          <td>{{ order.service }}</td>
          <td>{{ formatDate(order.orderDate) }}</td>
          <td>{{ formatDate(order.receiveDate) }}</td>
          <td>
            <!-- Если issueDate есть — форматируем, иначе показываем прочерк -->
            {{ order.issueDate ? formatDate(order.issueDate) : '—' }}
          </td>
          <td>
            <!-- Кнопка "Выдать" -->
            <button
                @click="openIssueModal(order)"
                :disabled="!isUpcoming(order)"
                class = "btn"
            >
              Выдать
            </button>
            <!-- Кнопка "Удалить" -->
            <button
                @click="deleteOrder(order.id)"
                :disabled="!isUpcoming(order)"
                class = "btn"
            >
              Удалить
            </button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>
    <div v-else>Загрузка заказов...</div>

    <!-- Пагинация -->
    <div v-if="totalPages > 1 && !isLoading">
      <button @click="prevPage" :disabled="currentPage === 1" class = "btn">Назад</button>
      <span>Страница {{ currentPage }} из {{ totalPages }}</span>
      <button @click="nextPage" :disabled="currentPage === totalPages" class = "btn">
        Вперед
      </button>
    </div>

    <!-- Карточки навигации -->
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
        <p>Услуги и оборудование</p>
      </div>
      <div class="card" @click="goToClients">
        <h3>Клиенты</h3>
        <p>Учёт и управление</p>
      </div>
    </div>

    <!-- Модальное окно (Добавление/Редактирование) -->
    <b-modal
        id="order-modal"
        v-model="showModal"
        :title="isEditing ? 'Редактировать заказ' : 'Добавить заказ'"
        hide-footer
        @hide="resetCurrentItem"
    >
      <form @submit.prevent="isEditing ? submitEdit() : submitAdd()">
        <div>
          <label for="phoneNumber">Номер телефона:</label>
          <input
              id="phoneNumber"
              type="tel"
              v-model="currentItem.phoneNumber"
              required
              @focus="showPhoneSuggestions = !!phoneSuggestions.length"
              @blur="hideSuggestions"
          />
          <ul v-if="showPhoneSuggestions && phoneSuggestions.length">
            <li
                v-for="phone in phoneSuggestions"
                :key="phone"
                @mousedown.prevent="selectPhone(phone)"
            >
              {{ phone }}
            </li>
          </ul>
        </div>

        <div>
          <label for="service">Название услуги:</label>
          <select
              id="service"
              v-model="currentItem.service"
              required
          >
            <option value="">Выберите услугу</option>
            <option
                v-for="service in uniqueServices"
                :key="service.name"
                :value="service.name"
            >
              {{ service.name }}
            </option>
          </select>
        </div>

        <div>
          <label for="orderDate">Дата оформления:</label>
          <input
              id="orderDate"
              type="datetime-local"
              v-model="currentItem.orderDate"
              :min="getCurrentDateTime()"
              required
          />
        </div>
        <div>
          <button type="submit" class = "btn">Сохранить</button>
          <button type="button" @click="closeModal" class = "btn">Отмена</button>
        </div>
      </form>
    </b-modal>

    <!-- Модальное окно для выдачи заказа -->
    <b-modal
        id="issue-modal"
        v-model="showIssueModal"
        title="Выдать заказ"
        hide-footer
        @hide="resetIssueModal"
    >
      <form @submit.prevent="submitIssue">
        <div>
          <label for="issueTime">Дата и время выдачи:</label>
          <input
              id="issueTime"
              type="datetime-local"
              v-model="issueTime"
              :min="getIssueMinTime()"
              required
          />
        </div>
        <div>
          <button type="submit" class = "btn">Подтвердить</button>
          <button type="button" @click="closeIssueModal" class = "btn">Отмена</button>
        </div>
      </form>
    </b-modal>
  </div>
</template>

<script>
import axios from 'axios';
import dayjs from 'dayjs';
import 'dayjs/locale/ru';
dayjs.locale('ru'); // русская локаль

export default {
  name: 'OrderHistory',
  data() {
    return {
      orders: [],
      uniqueServices: [],
      sortOrders: {
        clientName: 'asc',
        employeeName: 'asc',
        service: 'asc',
        orderDate: 'asc',
        receiveDate: 'asc',
        issueDate: 'asc',
      },
      currentSortKey: 'clientName',
      isLoading: true,
      currentPage: 1,
      pageSize: 20,

      // Модалка добавления/редактирования
      showModal: false,
      currentItem: {},
      isEditing: false,

      // Модалка выдачи заказа
      showIssueModal: false,
      issueTime: '',
      issueItem: {},

      // Фильтры
      searchClientName: '',
      searchService: '',

      // Автокомплит телефона
      phoneSuggestions: [],
      showPhoneSuggestions: false,
    };
  },
  computed: {
    filteredOrders() {
      return this.orders.filter((order) => {
        if (!order.clientName || !order.service) {
          return false;
        }
        const nameMatch = order.clientName
            .toLowerCase()
            .includes(this.searchClientName.toLowerCase());
        const serviceMatch = this.searchService
            ? order.service === this.searchService
            : true;
        return nameMatch && serviceMatch;
      });
    },
    sortedOrders() {
      if (!this.currentSortKey) {
        return this.filteredOrders;
      }
      return [...this.filteredOrders].sort((a, b) => {
        let aVal = a[this.currentSortKey] ?? '';
        let bVal = b[this.currentSortKey] ?? '';

        if (
            this.currentSortKey === 'orderDate' ||
            this.currentSortKey === 'receiveDate' ||
            this.currentSortKey === 'issueDate'
        ) {
          aVal = new Date(aVal);
          bVal = new Date(bVal);
        } else {
          aVal = aVal.toString().toLowerCase();
          bVal = bVal.toString().toLowerCase();
        }

        if (aVal < bVal) return this.sortOrders[this.currentSortKey] === 'asc' ? -1 : 1;
        if (aVal > bVal) return this.sortOrders[this.currentSortKey] === 'asc' ? 1 : -1;
        return 0;
      });
    },
    paginatedOrders() {
      const start = (this.currentPage - 1) * this.pageSize;
      const end = start + this.pageSize;
      return this.sortedOrders.slice(start, end);
    },
  },
  methods: {
    // Предстоящий заказ: если issueDate > сейчас ИЛИ orderDate/receiveDate > сейчас
    isUpcoming(order) {
      const now = dayjs();
      const futureIssue = order.issueDate && dayjs(order.issueDate).isAfter(now);
      const futureOrder = order.orderDate && dayjs(order.orderDate).isAfter(now);
      const futureReceive = order.receiveDate && dayjs(order.receiveDate).isAfter(now);
      return futureIssue || futureOrder || futureReceive;
    },

    // Открыть модалку выдачи
    openIssueModal(order) {
      if (!this.isUpcoming(order)) {
        this.showAlert('Выдача недоступна для этого заказа.', 'warning');
        return;
      }
      // Сохраняем заказ, для которого делаем "выдачу"
      this.issueItem = order;

      // Устанавливаем issueTime как максимум из текущего времени и receiveDate
      const receiveDateLocal = dayjs(order.receiveDate).format('YYYY-MM-DDTHH:mm');
      const nowLocal = dayjs().format('YYYY-MM-DDTHH:mm');
      this.issueTime = dayjs(receiveDateLocal).isAfter(nowLocal) ? receiveDateLocal : nowLocal;

      this.showIssueModal = true;
    },

    // Подтверждение выдачи — POST на сервер
    async submitIssue() {
      try {
        if (!this.issueTime) {
          this.showAlert('Не выбрано время выдачи!', 'error');
          return;
        }
        // Преобразуем локальное время в ISO с часовым поясом Москвы
        const issueDate = dayjs(this.issueTime).format();

        // POST запрос на сервер
        const response = await axios.post(
            `http://localhost:8080/api/orders/issue/${this.issueItem.id}`,
            { issueDate },
            { withCredentials: true }
        );

        if (response.data.message === 'Заказ успешно выдан') {
          // Обновляем локально поле issueDate, чтобы сразу отобразилось
          const idx = this.orders.findIndex((o) => o.id === this.issueItem.id);
          if (idx !== -1) {
            this.orders[idx].issueDate = issueDate;
          }
          this.showAlert('Заказ выдан.', 'success');
        } else {
          this.showAlert('Неизвестная ошибка при выдаче.', 'error');
        }
      } catch (error) {
        console.error('Ошибка при выдаче заказа:', error);
        if (error.response && error.response.data && error.response.data.error) {
          this.showAlert(`Ошибка: ${error.response.data.error}`, 'error');
        } else {
          this.showAlert('Произошла ошибка при выдаче заказа.', 'error');
        }
      } finally {
        this.closeIssueModal();
      }
    },

    // Закрыть модалку выдачи
    closeIssueModal() {
      this.showIssueModal = false;
      this.resetIssueModal();
    },
    resetIssueModal() {
      this.issueTime = '';
      this.issueItem = {};
    },

    // Удаление заказа
    async deleteOrder(orderId) {
      const order = this.orders.find((o) => o.id === orderId);
      if (!order) {
        this.showAlert('Заказ не найден.', 'error');
        return;
      }
      if (!this.isUpcoming(order)) {
        this.showAlert('Удаление недоступно для этого заказа.', 'warning');
        return;
      }

      try {
        const result = await this.$swal.fire({
          title: 'Вы уверены?',
          text: `Вы хотите удалить заказ клиента ${order.clientName}?`,
          icon: 'warning',
          showCancelButton: true,
          confirmButtonText: 'Да, удалить',
          cancelButtonText: 'Отмена',
          reverseButtons: true,
        });

        if (result.isConfirmed) {
          const response = await axios.delete(
              `http://localhost:8080/api/orders/${orderId}`,
              { withCredentials: true }
          );
          if (response.data.message === 'Заказ успешно удалён') {
            this.orders = this.orders.filter((o) => o.id !== orderId);
            this.showAlert('Заказ успешно удалён.', 'success');
            this.calculateTotalPages();
          } else {
            this.showAlert('Неизвестная ошибка при удалении заказа.', 'error');
          }
        } else if (result.dismiss === this.$swal.DismissReason.cancel) {
          this.showAlert('Удаление заказа отменено.', 'info');
        }
      } catch (error) {
        console.error('Ошибка при удалении заказа:', error);
        this.showAlert('Не удалось удалить заказ.', 'error');
      }
    },

    // Логика добавления/редактирования (осталась без изменений)
    openAddModal() {
      this.isEditing = false;
      this.currentItem = {
        phoneNumber: '',
        service: '',
        orderDate: this.getCurrentDateTime(),
        receiveDate: this.getCurrentDateTime(),
        issueDate: '',
      };
      this.showModal = true;
    },
    openEditModal(order) {
      this.isEditing = true;
      this.currentItem = { ...order };
      if (!this.currentItem.issueDate) {
        this.currentItem.issueDate = '';
      }
      this.showModal = true;
    },
    closeModal() {
      this.showModal = false;
      this.resetCurrentItem();
    },
    resetCurrentItem() {
      this.currentItem = {};
      this.isEditing = false;
    },

    async submitAdd() {
      try {
        const moscowOrderDate = this.toMoscowTimeISOString(
            new Date(this.currentItem.orderDate)
        );
        const moscowReceiveDate = this.toMoscowTimeISOString(
            new Date(this.currentItem.receiveDate)
        );
        const moscowIssueDate = this.currentItem.issueDate
            ? this.toMoscowTimeISOString(new Date(this.currentItem.issueDate))
            : null;

        const newOrder = {
          phoneNumber: this.currentItem.phoneNumber.replace(/\D/g, ''),
          service: this.currentItem.service,
          orderDate: moscowOrderDate,
          receiveDate: moscowReceiveDate,
          issueDate: moscowIssueDate,
        };

        const response = await axios.post(
            'http://localhost:8080/api/orders/admin',
            newOrder,
            { withCredentials: true }
        );

        if (response.data.message === 'Заказ успешно сохранен') {
          this.orders.push(response.data.order);
          this.showAlert('Заказ успешно добавлен.', 'success');
          this.closeModal();
          this.calculateTotalPages();
        } else {
          this.showAlert(
              'Неизвестная ошибка: ' + (response.data.message || 'Попробуйте позже.'),
              'error'
          );
        }
      } catch (error) {
        console.error('Ошибка при добавлении заказа:', error);
        if (error.response) {
          this.showAlert(
              'Ошибка сервера: ' + (error.response.data.error || 'Неизвестная ошибка'),
              'error'
          );
        } else if (error.request) {
          this.showAlert(
              'Ошибка соединения с сервером. Проверьте интернет-соединение.',
              'error'
          );
        } else {
          this.showAlert('Неизвестная ошибка: ' + error.message, 'error');
        }
      }
    },
    async submitEdit() {
      try {
        const moscowOrderDate = this.toMoscowTimeISOString(
            new Date(this.currentItem.orderDate)
        );
        const moscowReceiveDate = this.toMoscowTimeISOString(
            new Date(this.currentItem.receiveDate)
        );
        const moscowIssueDate = this.currentItem.issueDate
            ? this.toMoscowTimeISOString(new Date(this.currentItem.issueDate))
            : null;

        const updatedOrder = {
          phoneNumber: this.currentItem.phoneNumber.replace(/\D/g, ''),
          service: this.currentItem.service,
          orderDate: moscowOrderDate,
          receiveDate: moscowReceiveDate,
          issueDate: moscowIssueDate,
        };

        const response = await axios.put(
            `http://localhost:8080/api/orders/${this.currentItem.id}`,
            updatedOrder,
            { withCredentials: true }
        );

        if (response.data.message === 'Заказ успешно обновлён') {
          const index = this.orders.findIndex(
              (ord) => ord.id === this.currentItem.id
          );
          if (index !== -1) {
            this.orders.splice(index, 1, response.data.order);
          }
          this.showAlert('Заказ успешно обновлён.', 'success');
          this.closeModal();
          this.calculateTotalPages();
        } else {
          this.showAlert(
              'Неизвестная ошибка: ' + (response.data.message || 'Попробуйте позже.'),
              'error'
          );
        }
      } catch (error) {
        console.error('Ошибка при редактировании заказа:', error);
        if (error.response) {
          this.showAlert(
              'Ошибка сервера: ' + (error.response.data.error || 'Неизвестная ошибка'),
              'error'
          );
        } else if (error.request) {
          this.showAlert(
              'Ошибка соединения с сервером. Проверьте интернет-соединение.',
              'error'
          );
        } else {
          this.showAlert('Неизвестная ошибка: ' + error.message, 'error');
        }
      }
    },

    // Форматирование/преобразование дат
    toMoscowTimeISOString(date) {
      const options = {
        timeZone: 'Europe/Moscow',
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false,
      };
      const formatter = new Intl.DateTimeFormat('sv-SE', options);
      const parts = formatter.formatToParts(date);
      const dateParts = {};
      parts.forEach(({ type, value }) => {
        dateParts[type] = value;
      });
      return `${dateParts.year}-${dateParts.month}-${dateParts.day}T${dateParts.hour}:${dateParts.minute}:${dateParts.second}+03:00`;
    },
    formatDate(dateString) {
      if (!dateString) return '';
      return dayjs(dateString).format('D MMMM YYYY HH:mm');
    },
    getCurrentDateTime() {
      return dayjs().format('YYYY-MM-DDTHH:mm');
    },
    getIssueMinTime() {
      if (!this.issueItem || !this.issueItem.receiveDate) {
        return this.getCurrentDateTime();
      }
      const receiveDate = dayjs(this.issueItem.receiveDate).format('YYYY-MM-DDTHH:mm');
      const now = dayjs().format('YYYY-MM-DDTHH:mm');
      return dayjs(receiveDate).isAfter(now) ? receiveDate : now;
    },

    // Сортировка
    sortBy(key) {
      if (this.currentSortKey === key) {
        this.sortOrders[key] =
            this.sortOrders[key] === 'asc' ? 'desc' : 'asc';
      } else {
        this.currentSortKey = key;
        this.sortOrders[key] = 'asc';
      }
      this.calculateTotalPages();
    },
    // Пагинация
    nextPage() {
      if (this.currentPage < this.totalPages) {
        this.currentPage += 1;
      }
    },
    prevPage() {
      if (this.currentPage > 1) {
        this.currentPage -= 1;
      }
    },
    calculateTotalPages() {
      this.totalPages = Math.ceil(this.sortedOrders.length / this.pageSize) || 1;
      if (this.currentPage > this.totalPages) {
        this.currentPage = this.totalPages;
      }
    },

    // Уникальные услуги
    extractUniqueServices() {
      const servicesSet = new Set();
      this.orders.forEach((o) => {
        servicesSet.add(o.service);
      });
      this.uniqueServices = Array.from(servicesSet).map((name) => ({ name }));
    },

    // Загрузка
    async fetchOrders() {
      this.isLoading = true;
      try {
        const response = await axios.get('http://localhost:8080/api/orders/admin', {
          withCredentials: true,
        });
        this.orders = response.data.orders || [];
        console.log('Список заказов с сервера:', this.orders);
        this.extractUniqueServices();
        this.calculateTotalPages();
      } catch (error) {
        console.error('Ошибка при загрузке заказов:', error);
        this.showAlert('Не удалось загрузить заказы.', 'error');
      } finally {
        this.isLoading = false;
      }
    },

    // Уведомления
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

    // Навигация
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
    goToClients() {
      this.$router.push({ name: 'Clients' });
    },

    // Автокомплит
    hideSuggestions() {
      setTimeout(() => {
        this.showPhoneSuggestions = false;
      }, 200);
    },
    selectPhone(phone) {
      this.currentItem.phoneNumber = phone;
      this.showPhoneSuggestions = false;
    },
  },
  watch: {
    // Форматирование номера телефона + подсказки
    'currentItem.phoneNumber'(val) {
      if (!val) {
        this.phoneSuggestions = [];
        return;
      }
      let digits = val.replace(/\D/g, '');
      if (!digits.startsWith('7')) {
        digits = '7' + digits.replace(/^7+/, '');
      }
      if (digits.length > 1) {
        digits = digits[0] + '(' + digits.slice(1);
      }
      if (digits.length > 4) {
        digits = digits.slice(0, 5) + ')-' + digits.slice(5);
      }
      if (digits.length > 9) {
        digits = digits.slice(0, 10) + '-' + digits.slice(10);
      }
      if (digits.length > 12) {
        digits = digits.slice(0, 13) + '-' + digits.slice(13);
      }
      if (digits.length > 15) {
        digits = digits.slice(0, 16);
      }
      this.currentItem.phoneNumber = digits;

      const plainDigits = digits.replace(/\D/g, '');
      if (plainDigits.length >= 3) {
        let foundPhones = this.orders
            .map((o) => (o.phoneNumber ? o.phoneNumber.toString() : ''))
            .filter((p) => p.startsWith(plainDigits));
        foundPhones = [...new Set(foundPhones)];
        this.phoneSuggestions = foundPhones.slice(0, 5);
      } else {
        this.phoneSuggestions = [];
      }
      this.showPhoneSuggestions = this.phoneSuggestions.length > 0;
    },

    searchClientName() {
      this.currentPage = 1;
      this.calculateTotalPages();
    },
    searchService() {
      this.currentPage = 1;
      this.calculateTotalPages();
    },
    orders() {
      this.calculateTotalPages();
    },
  },
  mounted() {
    this.fetchOrders();
  },
};
</script>

<style scoped>
.order-history {
  padding: 20px;
  margin: 0 auto;
  max-width: 1200px;
}

h2 {
  text-align: center;
  margin-bottom: 20px;
  color: #333;
}

/* Кнопка "Добавить заказ" */
.add-order-button {
  text-align: right;
  margin-bottom: 15px;
}

.add-order-button .btn.add {
  background-color: #4caf50; /* Зеленый цвет */
  color: white;
  padding: 10px 20px;
  font-size: 1em;
  border-radius: 5px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  transition: background-color 0.3s ease, transform 0.2s ease;
}

.add-order-button .btn.add:hover {
  background-color: #45a049;
  transform: translateY(-2px);
}

/* Стили для таблицы заказов */
.order-list-section {
  background-color: #f9f9f9;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.1);
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 10px;
}

.data-table th,
.data-table td {
  padding: 12px;
  border: 1px solid #ddd;
  text-align: left;
}

.data-table th {
  background-color: #f2f2f2;
  position: relative;
  user-select: none;
}

.data-table th:hover {
  background-color: #e0e0e0;
}

.data-table th span {
  margin-left: 5px;
  font-size: 0.8em;
}

.status-working {
  color: green;
  font-weight: bold;
}

.status-fired {
  color: red;
  font-weight: bold;
}

/* Стили для кнопок */
.btn {
  background-color: #4caf50; /* Зеленый цвет */
  color: white;
  padding: 8px 12px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 1em;
  transition: background-color 0.3s ease;
  margin: 5px;
}

.btn.edit {
  background-color: #007bff; /* Синий цвет для редактирования */
}

.btn.edit:hover {
  background-color: #0069d9;
}

.btn.danger {
  background-color: #f44336; /* Красный цвет для удаления */
}

.btn.danger:hover {
  background-color: #d32f2f;
}

.btn.save {
  background-color: #4caf50; /* Зеленый цвет */
}

.btn.save:hover {
  background-color: #45a049;
}

.btn.cancel {
  background-color: red; /* Красный цвет */
}

.btn.cancel:hover {
  background-color: darkred;
}

/* Стили для пагинации */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 15px;
  margin-top: 15px;
}

.pagination .btn {
  background-color: #4caf50; /* Зеленый цвет */
}

.pagination .btn:hover:not(:disabled) {
  background-color: #45a049;
}

.pagination .btn:disabled {
  background-color: #a5d6a7;
  cursor: not-allowed;
}

.pagination span {
  font-size: 1em;
  color: #333;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 15px;
}

/* Стили для панели навигации с карточками */
.card-panel {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
  justify-content: center;
  margin-top: 30px;
}

.card {
  background-color: #4caf50; /* Зеленый цвет */
  border-radius: 10px;
  padding: 15px;
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

.clientNameInput {
  width: 98%;
  padding: 12px; /* Увеличим padding для большего комфорта */
  margin: 10px 0;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 16px; /* Увеличим размер шрифта */
  transition: border-color 0.3s ease-in-out;
}

input[type="text"],
input[type="datetime-local"],
select {
  width: 100%;
  padding: 12px; /* Увеличим padding для большего комфорта */
  margin: 10px 0;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 16px; /* Увеличим размер шрифта */
  transition: border-color 0.3s ease-in-out;
}

input[type="text"]:focus,
input[type="datetime-local"]:focus,
select:focus {
  border-color: #007bff;
  outline: none;
}

input[type="text"]:hover,
input[type="datetime-local"]:hover,
select:hover {
  border-color: #aaa;
}

label {
  font-size: 16px; /* Увеличим размер шрифта для label */
  margin-bottom: 5px;
  display: block;
  font-weight: bold;
  color: #333;
}

/* Стили для модальных окон */
.modal input[type="tel"],
.modal input[type="text"],
.modal input[type="datetime-local"],
.modal select {
  margin: 10px 0;
  padding: 12px;
  font-size: 16px; /* Увеличиваем размер шрифта */
  border-radius: 4px;
  width: 100%;
  box-sizing: border-box;
  border: 1px solid #ddd;
}

select[name="service"] {
  max-width: 200px; /* Ограничим максимальную ширину поля */
  margin: 10px 0;
}

.modal input[type="tel"]:focus,
.modal input[type="text"]:focus,
.modal input[type="datetime-local"]:focus,
.modal select:focus {
  border-color: #0056b3;
  outline: none;
}
</style>
