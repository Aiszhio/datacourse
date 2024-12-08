<template>
  <div class="materials-overview">
    <h2>Управление материалами</h2>

    <!-- Таблица списка материалов -->
    <div class="materials-section">
      <h3>Список материалов</h3>
      <p>Для чернил указаны флаконы по 30мл штука</p>
      <table>
        <thead>
        <tr>
          <th>Название материала</th>
          <th>Количество</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(material) in sortedMaterials" :key="material.material_id">
          <td>{{ material.material_name }}</td>
          <td>{{ material.quantity }} <span>{{ getUnit(material.material_name) }}</span></td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Таблица расхода материалов -->
    <div class="expenditures-section">
      <h3>Расход материалов</h3>
      <button @click="openAddExpenditureModal" class="btn">Добавить расход</button>
      <table>
        <thead>
        <tr>
          <th>Название материала</th>
          <th @click="sortExpenditures">Дата расхода
            <span v-if="sortExpenditureDirection === 'asc'">▲</span>
            <span v-else>▼</span>
          </th>
          <th>Количество расхода</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(expenditure) in paginatedExpenditures" :key="expenditure.expenditure_id">
          <td>{{ expenditure.material_name }}</td>
          <td>{{ formatDate(expenditure.expenditure_date) }}</td>
          <td>{{ expenditure.quantity }} <span>{{ getUnit(expenditure.material_name) }}</span></td>
        </tr>
        </tbody>
      </table>
      <div class="pagination" v-if="maxPages(expenditures) > 1">
        <button @click="changeExpendituresPage(expendituresPage - 1)" :disabled="expendituresPage <= 1">Назад</button>
        <span>Страница {{ expendituresPage }} из {{ maxPages(expenditures) }}</span>
        <button @click="changeExpendituresPage(expendituresPage + 1)" :disabled="expendituresPage >= maxPages(expenditures)">Вперёд</button>
      </div>
    </div>

    <!-- Модальное окно для добавления расхода -->
    <div v-if="showExpenditureModal" class="modal-overlay">
      <div class="modal">
        <h3>Добавить расход материала</h3>

        <!-- Поле для выбора материала -->
        <label for="expenditureMaterialName">Название материала</label>
        <select id="expenditureMaterialName" v-model="currentExpenditure.material_id" class="input" required>
          <option value="" disabled>Выберите материал</option>
          <option v-for="material in materials" :key="material.material_id" :value="material.material_id">
            {{ material.material_name }}
          </option>
        </select>

        <!-- Поле для ввода количества -->
        <label for="expenditureQuantity">Количество</label>
        <div class="quantity-container">
          <input
              id="expenditureQuantity"
              type="number"
              v-model="currentExpenditure.quantity"
              class="input"
              placeholder="Количество"
              min="1"
              max="999"
              required
          />
          <span class="unit">шт</span>
        </div>

        <!-- Поле для ввода даты расхода -->
        <label for="expenditureDate">Дата расхода</label>
        <input
            id="expenditureDate"
            type="datetime-local"
            v-model="currentExpenditure.expenditure_date"
            :max="todayDate"
            class="input"
            required
        />

        <button @click="saveExpenditure" class="btn">Сохранить</button>
        <button @click="closeExpenditureModal" class="btn danger">Отмена</button>
      </div>
    </div>

    <!-- Таблица закупок материалов -->
    <div class="purchases-section">
      <h3>Закупка материалов</h3>
      <button @click="openAddPurchaseModal" class="btn">Добавить закупку</button>
      <table>
        <thead>
        <tr>
          <th>Название материала</th>
          <th>Стоимость</th>
          <th>Поставщик</th>
          <th>Количество</th>
          <th @click="sortPurchases">Дата поставки
            <span v-if="sortPurchaseDirection === 'asc'">▲</span>
            <span v-else>▼</span>
          </th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(purchase) in paginatedPurchases" :key="purchase.purchase_id">
          <td>{{ purchase.material_name }}</td>
          <td>{{ purchase.cost }} <span class="unit">руб</span></td>
          <td>{{ purchase.supplier }}</td>
          <td>{{ purchase.quantity }} <span>{{ getUnit(purchase.material_name) }}</span></td>
          <td>{{ formatDate(purchase.supply_date) }}</td>
        </tr>
        </tbody>
      </table>
      <div class="pagination" v-if="maxPages(purchases) > 1">
        <button @click="changePurchasesPage(purchasesPage - 1)" :disabled="purchasesPage <= 1">Назад</button>
        <span>Страница {{ purchasesPage }} из {{ maxPages(purchases) }}</span>
        <button @click="changePurchasesPage(purchasesPage + 1)" :disabled="purchasesPage >= maxPages(purchases)">Вперёд</button>
      </div>
    </div>

    <!-- Модальное окно для добавления закупки -->
    <div v-if="showPurchaseModal" class="modal-overlay">
      <div class="modal">
        <h3>Добавить закупку материала</h3>

        <!-- Поле для выбора материала -->
        <label for="purchaseMaterialName">Название материала</label>
        <select id="purchaseMaterialName" v-model="currentPurchase.material_id" class="input" required>
          <option value="" disabled>Выберите материал</option>
          <option v-for="material in materials" :key="material.material_id" :value="material.material_id">
            {{ material.material_name }}
          </option>
        </select>

        <!-- Поле для ввода поставщика -->
        <label for="purchaseSupplier">Поставщик</label>
        <input
            id="purchaseSupplier"
            type="text"
            v-model="currentPurchase.supplier"
            class="input"
            placeholder="Введите поставщика"
            required
            maxlength="40"
        />

        <!-- Поле для ввода количества -->
        <label for="purchaseQuantity">Количество</label>
        <div class="quantity-container">
          <input
              id="purchaseQuantity"
              type="number"
              v-model="currentPurchase.quantity"
              class="input"
              placeholder="Количество"
              min="1"
              max="999"
              required
          />
          <span class="unit">шт</span>
        </div>

        <!-- Поле для ввода цены -->
        <label for="purchaseCost">Цена</label>
        <div class="quantity-container">
          <input
              id="purchaseCost"
              type="number"
              v-model="currentPurchase.cost"
              class="input"
              placeholder="Цена"
              min="0"
              required
          />
          <span class="unit">руб</span>
        </div>

        <!-- Поле для ввода даты поставки -->
        <label for="supplyDate">Дата поставки</label>
        <input
            id="supplyDate"
            type="datetime-local"
            v-model="currentPurchase.supply_date"
            :min="threeDaysAgo"
            :max="todayDate"
            class="input"
            required
        />

        <button @click="savePurchase" class="btn">Сохранить</button>
        <button @click="closePurchaseModal" class="btn danger">Отмена</button>
      </div>
    </div>

    <!-- Панель навигации с карточками -->
    <div class="card-panel">
      <div class="card" @click="goToAdminHome">
        <h4>Главная</h4>
        <p>Панель администратора</p>
      </div>
      <div class="card" @click="goToBookingsPage">
        <h4>Бронирование</h4>
        <p>Управление бронированием</p>
      </div>
      <div class="card" @click="goToOrderHistory">
        <h4>Заказы</h4>
        <p>История заказов</p>
      </div>
      <div class="card" @click="goToServicesPage">
        <h4>Услуги</h4>
        <p>Услуги и оборудование</p>
      </div>
      <div class="card" @click="goToEmployeesPage">
        <h4>Сотрудники</h4>
        <p>Управление персоналом</p>
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
import axios from 'axios';

export default {
  name: 'MaterialsOverview',
  data() {
    return {
      materials: [],
      expenditures: [],
      purchases: [],
      showExpenditureModal: false,
      showPurchaseModal: false,
      currentExpenditure: {},
      currentPurchase: {},
      expendituresPage: 1,
      purchasesPage: 1,
      pageSize: 10,
      sortExpenditureDirection: 'asc',
      sortPurchaseDirection: 'asc',
      todayDate: this.getTodayDate(),
      threeDaysAgo: this.getThreeDaysAgo(),
    };
  },
  computed: {
    sortedMaterials() {
      return this.materials.slice().sort((a, b) => a.material_name.localeCompare(b.material_name));
    },
    sortedExpenditures() {
      return this.expenditures.slice().sort((a, b) => {
        const dateA = new Date(a.expenditure_date);
        const dateB = new Date(b.expenditure_date);
        return this.sortExpenditureDirection === 'asc' ? dateA - dateB : dateB - dateA;
      });
    },
    sortedPurchases() {
      return this.purchases.slice().sort((a, b) => {
        const dateA = new Date(a.supply_date);
        const dateB = new Date(b.supply_date);
        return this.sortPurchaseDirection === 'asc' ? dateA - dateB : dateB - dateA;
      });
    },
    paginatedExpenditures() {
      const startIndex = (this.expendituresPage - 1) * this.pageSize;
      return this.sortedExpenditures.slice(startIndex, startIndex + this.pageSize);
    },
    paginatedPurchases() {
      const startIndex = (this.purchasesPage - 1) * this.pageSize;
      return this.sortedPurchases.slice(startIndex, startIndex + this.pageSize);
    },
  },
  methods: {
    async fetchMaterials() {
      try {
        const response = await axios.get('http://localhost:8080/api/materials', { withCredentials: true });
        this.materials = response.data.materials;
        console.log('Материалы успешно загружены:', this.materials);
      } catch (error) {
        console.error('Ошибка при загрузке материалов:', error);
        alert('Не удалось загрузить материалы.');
      }
    },

    async fetchExpenditures() {
      try {
        const response = await axios.get('http://localhost:8080/api/expenditures', { withCredentials: true });
        this.expenditures = response.data.expenditures;
        console.log('Расходы успешно загружены:', this.expenditures);
      } catch (error) {
        console.error('Ошибка при загрузке расходов:', error);
        alert('Не удалось загрузить расходы.');
      }
    },

    async fetchPurchases() {
      try {
        const response = await axios.get('http://localhost:8080/api/purchases', { withCredentials: true });
        this.purchases = response.data.purchases;
        console.log('Закупки успешно загружены:', this.purchases);
      } catch (error) {
        console.error('Ошибка при загрузке закупок:', error);
        alert('Не удалось загрузить закупки.');
      }
    },

    openAddExpenditureModal() {
      this.showExpenditureModal = true;
      this.currentExpenditure = {
        material_id: '',
        expenditure_date: '',
        quantity: '',
      };
    },

    closeExpenditureModal() {
      this.showExpenditureModal = false;
    },

    openAddPurchaseModal() {
      this.showPurchaseModal = true;
      this.currentPurchase = {
        material_id: '',
        supplier: '',
        quantity: '',
        cost: '',
        supply_date: '',
      };
    },

    closePurchaseModal() {
      this.showPurchaseModal = false;
    },

    async saveExpenditure() {
      try {
        const formattedExpenditure = {
          ...this.currentExpenditure,
          expenditure_date: new Date(this.currentExpenditure.expenditure_date).toISOString(),
        };

        const response = await axios.post('http://localhost:8080/api/expenditures', formattedExpenditure, { withCredentials: true });

        if (response.data.message) {
          console.log('Сообщение с сервера:', response.data.message);
          alert(`${response.data.message}`);
        } else {
          console.log('Ответ от сервера без сообщения:', response.data);
          alert('Расход успешно добавлен, но нет сообщения от сервера.');
        }

        this.closeExpenditureModal();
        this.fetchExpenditures();
        this.fetchMaterials();
      } catch (error) {
        if (error.response) {
          const errorMessage = error.response.data.error || 'Неизвестная ошибка с сервера';
          console.error('Ошибка при сохранении расхода:', errorMessage);
          alert(`Ошибка: ${errorMessage}`);
        } else if (error.request) {
          console.error('Ошибка при сохранении расхода: нет ответа от сервера');
          alert('Ошибка: нет ответа от сервера');
        } else {
          console.error('Ошибка при настройке запроса:', error.message);
          alert(`Ошибка: ${error.message}`);
        }
      }
    },

    async savePurchase() {
      try {
        if (isNaN(this.currentPurchase.cost) || this.currentPurchase.cost <= 0) {
          alert('Введите правильную цену для закупки!');
          return;
        }

        if (!this.currentPurchase.supplier || this.currentPurchase.supplier.trim() === '') {
          alert('Введите имя поставщика!');
          return;
        }

        const formattedPurchase = {
          ...this.currentPurchase,
          supply_date: new Date(this.currentPurchase.supply_date).toISOString(),
        };

        const response = await axios.post('http://localhost:8080/api/purchases', formattedPurchase, { withCredentials: true });

        if (response.data.message) {
          console.log('Сообщение с сервера:', response.data.message);
          alert(`${response.data.message}`);
        } else {
          console.log('Ответ от сервера без сообщения:', response.data);
          alert('Закупка успешно добавлена, но нет сообщения от сервера.');
        }

        this.closePurchaseModal();
        this.fetchPurchases();
        this.fetchMaterials();
      } catch (error) {
        if (error.response) {
          const errorMessage = error.response.data.error || 'Неизвестная ошибка с сервера';
          console.error('Ошибка при сохранении закупки:', errorMessage);
          alert(`Ошибка: ${errorMessage}`);
        } else if (error.request) {
          console.error('Ошибка при сохранении закупки: нет ответа от сервера');
          alert('Ошибка: нет ответа от сервера');
        } else {
          console.error('Ошибка при настройке запроса:', error.message);
          alert(`Ошибка: ${error.message}`);
        }
      }
    },

    formatDate(dateString) {
      const date = new Date(dateString);
      const day = String(date.getDate()).padStart(2, '0');
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const year = date.getFullYear();
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      return `${day}-${month}-${year} ${hours}:${minutes}`;
    },

    changeExpendituresPage(newPage) {
      if (newPage >= 1 && newPage <= this.maxPages(this.expenditures)) {
        this.expendituresPage = newPage;
      }
    },

    changePurchasesPage(newPage) {
      if (newPage >= 1 && newPage <= this.maxPages(this.purchases)) {
        this.purchasesPage = newPage;
      }
    },

    maxPages(data) {
      return Math.ceil(data.length / this.pageSize);
    },

    sortExpenditures() {
      this.sortExpenditureDirection = this.sortExpenditureDirection === 'asc' ? 'desc' : 'asc';
    },

    sortPurchases() {
      this.sortPurchaseDirection = this.sortPurchaseDirection === 'asc' ? 'desc' : 'asc';
    },

    getTodayDate() {
      const date = new Date();
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      return `${year}-${month}-${day}T${hours}:${minutes}`;
    },

    getThreeDaysAgo() {
      const date = new Date();
      date.setDate(date.getDate() - 3);
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      return `${year}-${month}-${day}T${hours}:${minutes}`;
    },

    getUnit(materialName) {
      const units = {
        'Фото бумага премиум': 'шт',
        'Холст для печати': 'шт',
        'Картон для фотокниг': 'шт',
        'Фото рамки': 'шт',
        'Обложки для фотокниг': 'шт',
        'Клеевые полоски для альбомов': 'шт',
        'Защитные пленки для фотографий': 'шт',
        'Чернила для принтера': 'фл',
        'Пленка для фотоальбомов': 'шт',
        'Фотографическая пленка': 'шт',
      };
      return units[materialName] || '';
    },

    goToAdminHome() {
      this.$router.push({ name: 'AdminHome' });
    },
    goToOrderHistory() {
      this.$router.push({ name: 'OrderHistory' });
    },
    goToEmployeesPage() {
      this.$router.push({ name: 'ManageEmp' });
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
  created() {
    this.fetchMaterials();
    this.fetchExpenditures();
    this.fetchPurchases();
  },
};
</script>

<style scoped>
.materials-overview {
  padding: 20px;
  max-width: 1500px;
  max-height: 100vh;
  margin: 0 auto;
}

h2 {
  text-align: center;
  margin-bottom: 20px;
}

.materials-section, .expenditures-section, .purchases-section {
  background-color: #f9f9f9;
  padding: 20px;
  margin-bottom: 20px;
  border-radius: 10px;
  box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.1);
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

label {
  display: block;
  font-weight: bold;
  margin-top: 10px;
}

.input {
  width: 100%;
  padding: 10px;
  margin: 10px 0;
  border: 1px solid #ccc;
  border-radius: 5px;
  box-sizing: border-box;
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

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.pagination button {
  padding: 10px;
  margin: 0 5px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.pagination button:disabled {
  background-color: #ccc;
}
.quantity-container {
  display: flex;
  align-items: center;
}

.quantity-container input {
  margin-right: 5px;
  width: 80px;
}

.quantity-container .unit {
  font-size: 14px;
}
</style>
