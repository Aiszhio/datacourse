<template>
  <div class="materials-overview">
    <h2>Управление материалами</h2>

    <!-- Таблица списка материалов -->
    <div class="materials-section">
      <h3>Список материалов</h3>
      <p>Для чернил указаны флаконы по 50мл штука</p>
      <table>
        <thead>
        <tr>
          <th>Номер материала</th>
          <th>Название материала</th>
          <th>Количество</th> <!-- Новый столбец для количества -->
        </tr>
        </thead>
        <tbody>
        <tr v-for="(material, index) in materials" :key="material.material_id">
          <td>{{ index + 1 }}</td>  <!-- Индекс будет начинаться с 1 -->
          <td>{{ material.material_name }}</td>
          <td>{{ material.quantity }} <span>{{ getUnit(material.material_name) }}</span></td> <!-- Добавляем единицу измерения -->
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
          <th>Номер расхода</th>
          <th>Название материала</th>
          <th @click="sortExpenditures('expenditure_date')">Дата расхода</th>
          <th>Количество расхода</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(expenditure, index) in paginatedExpenditures" :key="expenditure.expenditure_id">
          <td>{{ index + 1 + (expendituresPage - 1) * pageSize }}</td> <!-- Индекс будет начинаться с 1, для текущей страницы -->
          <td>{{ expenditure.material_name }}</td>
          <td>{{ formatDate(expenditure.expenditure_date) }}</td>
          <td>{{ expenditure.quantity }} <span>{{ getUnit(expenditure.material_name) }}</span></td> <!-- Единица измерения для расхода -->
        </tr>
        </tbody>
      </table>
      <div class="pagination">
        <button @click="changeExpendituresPage(expendituresPage - 1)" :disabled="expendituresPage <= 1">Назад</button>
        <span>Страница {{ expendituresPage }} из {{ maxPages(expenditures) }}</span>
        <button @click="changeExpendituresPage(expendituresPage + 1)" :disabled="expendituresPage >= maxPages(expenditures)">Вперёд</button>
      </div>
    </div>

    <!-- Таблица закупок материалов -->
    <div class="purchases-section">
      <h3>Закупка материалов</h3>
      <button @click="openAddPurchaseModal" class="btn">Добавить закупку</button>
      <table>
        <thead>
        <tr>
          <th>Номер закупки</th>
          <th>Название материала</th>
          <th>Стоимость</th>
          <th>Поставщик</th>
          <th>Количество</th>
          <th @click="sortPurchases('supply_date')">Дата поставки</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(purchase, index) in paginatedPurchases" :key="purchase.purchase_id">
          <td>{{ index + 1 + (purchasesPage - 1) * pageSize }}</td>  <!-- Индекс будет начинаться с 1, для текущей страницы -->
          <td>{{ purchase.material_name }}</td>
          <td>{{ purchase.cost }}</td>
          <td>{{ purchase.supplier }}</td>
          <td>{{ purchase.quantity }} <span>{{ getUnit(purchase.material_name) }}</span></td> <!-- Единица измерения для закупки -->
          <td>{{ formatDate(purchase.supply_date) }}</td>
        </tr>
        </tbody>
      </table>
      <div class="pagination">
        <button @click="changePurchasesPage(purchasesPage - 1)" :disabled="purchasesPage <= 1">Назад</button>
        <span>Страница {{ purchasesPage }} из {{ maxPages(purchases) }}</span>
        <button @click="changePurchasesPage(purchasesPage + 1)" :disabled="purchasesPage >= maxPages(purchases)">Вперёд</button>
      </div>
    </div>

    <!-- Модальные окна для добавления расхода и закупки -->
    <div v-if="showExpenditureModal" class="modal-overlay">
      <div class="modal">
        <h3>Добавить расход материала</h3>
        <label for="materialName">Название материала</label>
        <select id="materialName" v-model="currentExpenditure.material_id" class="input" required>
          <option value="" disabled>Выберите материал</option>
          <option v-for="material in materials" :key="material.material_id" :value="material.material_id">
            {{ material.material_name }}
          </option>
        </select>

        <label for="expenditureDate">Дата расхода</label>
        <input id="expenditureDate" type="date" v-model="currentExpenditure.expenditure_date" class="input" :max="maxDate" required />

        <label for="quantity">Количество расхода</label>
        <div class="quantity-container">
          <input
              id="quantity"
              type="number"
              v-model="currentExpenditure.quantity"
              class="input"
              placeholder="Количество расхода"
              min="1"
              max="500"
              required
              maxlength="4"
          />
          <span class="unit">шт</span> <!-- Здесь добавляем единицу измерения -->
        </div>

        <button @click="saveExpenditure" class="btn">Сохранить</button>
        <button @click="closeExpenditureModal" class="btn danger">Отмена</button>
      </div>
    </div>

    <div v-if="showPurchaseModal" class="modal-overlay">
      <div class="modal">
        <h3>Добавить закупку материала</h3>
        <label for="purchaseMaterialName">Название материала</label>
        <select id="purchaseMaterialName" v-model="currentPurchase.material_id" class="input" required>
          <option value="" disabled>Выберите материал</option>
          <option v-for="material in materials" :key="material.material_id" :value="material.material_id">
            {{ material.material_name }}
          </option>
        </select>

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
              maxlength="4"
          />
          <span class="unit">шт</span> <!-- Добавляем единицу измерения -->
        </div>


        <label for="supplyDate">Дата поставки</label>
        <input id="supplyDate" type="date" v-model="currentPurchase.supply_date" class="input" :min="threeDaysAgo" required />

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
      <div class="card" @click="goToEmployeesPage">
        <h4>Сотрудники</h4>
        <p>Управление персоналом</p>
      </div>
      <div class="card" @click="goToOrdersPage">
        <h4>Заказы</h4>
        <p>История заказов</p>
      </div>
      <div class="card" @click="goToServicesPage">
        <h4>Услуги</h4>
        <p>Услуги и оборудование</p>
      </div>
      <div class="card" @click="goToBookingPage">
        <h4>Бронирования</h4>
        <p>Управление бронированиями</p>
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
      expendituresPage: 1, // Пагинация для расхода
      purchasesPage: 1, // Пагинация для закупки
      pageSize: 10,
      sortExpenditureDirection: 'asc',
      sortPurchaseDirection: 'asc',
      maxDate: this.getMaxDate(), // Максимальная дата для расхода (не в будущем)
      threeDaysAgo: this.getThreeDaysAgo() // Дата для закупки (не раньше 3 дней назад)
    };
  },
  created() {
    this.fetchMaterials();
    this.fetchExpenditures();
    this.fetchPurchases();
  },
  computed: {
    paginatedExpenditures() {
      const startIndex = (this.expendituresPage - 1) * this.pageSize;
      return this.sortedExpenditures.slice(startIndex, startIndex + this.pageSize);
    },
    paginatedPurchases() {
      const startIndex = (this.purchasesPage - 1) * this.pageSize;
      return this.sortedPurchases.slice(startIndex, startIndex + this.pageSize);
    },
    sortedExpenditures() {
      return this.expenditures.sort((a, b) => {
        const dateA = new Date(a.expenditure_date);
        const dateB = new Date(b.expenditure_date);
        return this.sortExpenditureDirection === 'asc' ? dateA - dateB : dateB - dateA;
      });
    },
    sortedPurchases() {
      return this.purchases.sort((a, b) => {
        const dateA = new Date(a.supply_date);
        const dateB = new Date(b.supply_date);
        return this.sortPurchaseDirection === 'asc' ? dateA - dateB : dateB - dateA;
      });
    }
  },
  methods: {
    // Получение списка материалов
    async fetchMaterials() {
      try {
        const response = await axios.get('http://localhost:8080/api/materials', { withCredentials: true });
        this.materials = response.data.materials;
      } catch (error) {
        console.error('Ошибка при загрузке материалов:', error);
        alert('Не удалось загрузить материалы.');
      }
    },

    // Получение списка расходов
    async fetchExpenditures() {
      try {
        const response = await axios.get('http://localhost:8080/api/expenditures', { withCredentials: true });
        this.expenditures = response.data.expenditures;
      } catch (error) {
        console.error('Ошибка при загрузке расходов:', error);
        alert('Не удалось загрузить расходы.');
      }
    },

    // Получение списка закупок
    async fetchPurchases() {
      try {
        const response = await axios.get('http://localhost:8080/api/purchases', { withCredentials: true });
        this.purchases = response.data.purchases;
      } catch (error) {
        console.error('Ошибка при загрузке закупок:', error);
        alert('Не удалось загрузить закупки.');
      }
    },

    // Открытие модального окна для добавления расхода
    openAddExpenditureModal() {
      this.showExpenditureModal = true;
      this.currentExpenditure = {
        material_id: '',
        expenditure_date: '',
        quantity: ''
      };
    },

    // Закрытие модального окна расхода
    closeExpenditureModal() {
      this.showExpenditureModal = false;
    },

    // Сохранение расхода
    async saveExpenditure() {
      try {
        await axios.post('http://localhost:8080/api/expenditures', this.currentExpenditure, { withCredentials: true });
        this.closeExpenditureModal();
        this.fetchExpenditures();
        this.fetchMaterials(); // Обновить количество материалов
      } catch (error) {
        console.error('Ошибка при сохранении расхода:', error);
        alert('Не удалось сохранить расход.');
      }
    },

    // Открытие модального окна для добавления закупки
    openAddPurchaseModal() {
      this.showPurchaseModal = true;
      this.currentPurchase = {
        material_id: '',
        supplier: '',
        quantity: '',
        cost: '',
        supply_date: ''
      };
    },

    // Закрытие модального окна закупки
    closePurchaseModal() {
      this.showPurchaseModal = false;
    },

    // Сохранение закупки
    async savePurchase() {
      try {
        await axios.post('http://localhost:8080/api/purchases', this.currentPurchase, { withCredentials: true });
        this.closePurchaseModal();
        this.fetchPurchases();
        this.fetchMaterials(); // Обновить количество материалов
      } catch (error) {
        console.error('Ошибка при сохранении закупки:', error);
        alert('Не удалось сохранить закупку.');
      }
    },

    // Форматирование даты
    formatDate(dateString) {
      const date = new Date(dateString);
      return `${date.getDate().toString().padStart(2, '0')}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getFullYear()}`;
    },

    // Пагинация для расходов
    changeExpendituresPage(newPage) {
      this.expendituresPage = newPage;
    },

    // Пагинация для закупок
    changePurchasesPage(newPage) {
      this.purchasesPage = newPage;
    },

    // Максимальное количество страниц
    maxPages(data) {
      return Math.ceil(data.length / this.pageSize);
    },

    // Сортировка расходов
    sortExpenditures(column) {
      this.sortExpenditureDirection = this.sortExpenditureDirection === 'asc' ? 'desc' : 'asc';
    },

    // Сортировка закупок
    sortPurchases(column) {
      this.sortPurchaseDirection = this.sortPurchaseDirection === 'asc' ? 'desc' : 'asc';
    },

    // Получение максимальной даты
    getMaxDate() {
      const date = new Date();
      return date.toISOString().split('T')[0]; // Возвращает строку в формате YYYY-MM-DD
    },

    // Получение даты, не ранее 3 дней назад
    getThreeDaysAgo() {
      const date = new Date();
      date.setDate(date.getDate() - 3);
      return date.toISOString().split('T')[0]; // Возвращает строку в формате YYYY-MM-DD
    },

    // Получение единицы измерения в зависимости от названия материала
    getUnit(materialName) {
      switch (materialName) {
        case 'Фото бумага премиум':
        case 'Холст для печати':
        case 'Картон для фотокниг':
        case 'Фото рамки':
        case 'Обложки для фотокниг':
        case 'Клеевые полоски для альбомов':
        case 'Защитные пленки для фотографий':
        case 'Чернила для принтера':
        case 'Пленка для фотоальбомов':
        case 'Фотографическая пленка':
          return 'шт';
        default:
          return '';
      }
    }
  }
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
  margin-right: 5px; /* Отступ между полем и единицей измерения */
  width: 80px; /* Устанавливаем ширину поля */
}

.quantity-container .unit {
  font-size: 14px; /* Размер шрифта для единицы измерения */
}
</style>
