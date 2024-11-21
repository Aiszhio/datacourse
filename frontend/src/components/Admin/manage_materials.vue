<template>
  <div class="materials-overview">
    <h2>Управление материалами</h2>

    <!-- Таблица списка материалов -->
    <div class="materials-section">
      <h3>Список материалов</h3>
      <table>
        <thead>
        <tr>
          <th>Номер материала</th>
          <th>Название материала</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="material in materials" :key="material.material_id">
          <td>{{ material.material_id }}</td>
          <td>{{ material.material_name }}</td>
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
          <th>Дата расхода</th>
          <th>Количество расхода</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="expenditure in expenditures" :key="expenditure.expenditure_id">
          <td>{{ expenditure.expenditure_id }}</td>
          <td>{{ expenditure.material_name }}</td>
          <td>{{ formatDate(expenditure.expenditure_date) }}</td>
          <td>{{ expenditure.quantity }}</td>
        </tr>
        </tbody>
      </table>
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
          <th>Дата поставки</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="purchase in purchases" :key="purchase.purchase_id">
          <td>{{ purchase.purchase_id }}</td>
          <td>{{ purchase.material_name }}</td>
          <td>{{ purchase.cost }}</td>
          <td>{{ purchase.supplier }}</td>
          <td>{{ purchase.quantity }}</td>
          <td>{{ formatDate(purchase.supply_date) }}</td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Модальное окно для добавления расхода материалов -->
    <div v-if="showExpenditureModal" class="modal-overlay">
      <div class="modal">
        <h3>Добавить расход материала</h3>

        <label for="materialId">Номер материала</label>
        <input
            id="materialId"
            v-model="currentExpenditure.material_id"
            class="input"
            placeholder="Номер материала"
            required
        />

        <label for="expenditureDate">Дата расхода</label>
        <input
            id="expenditureDate"
            type="date"
            v-model="currentExpenditure.expenditure_date"
            class="input"
            required
        />

        <label for="quantity">Количество расхода</label>
        <input
            id="quantity"
            type="number"
            v-model="currentExpenditure.quantity"
            class="input"
            placeholder="Количество расхода"
            min="1"
            required
        />

        <button @click="saveExpenditure" class="btn">Сохранить</button>
        <button @click="closeExpenditureModal" class="btn danger">Отмена</button>
      </div>
    </div>

    <!-- Модальное окно для добавления закупки материалов -->
    <div v-if="showPurchaseModal" class="modal-overlay">
      <div class="modal">
        <h3>Добавить закупку материала</h3>

        <label for="purchaseMaterialId">Номер материала</label>
        <input
            id="purchaseMaterialId"
            v-model="currentPurchase.material_id"
            class="input"
            placeholder="Номер материала"
            required
        />

        <label for="supplier">Поставщик</label>
        <input
            id="supplier"
            v-model="currentPurchase.supplier"
            class="input"
            placeholder="Поставщик"
            required
        />

        <label for="purchaseQuantity">Количество</label>
        <input
            id="purchaseQuantity"
            type="number"
            v-model="currentPurchase.quantity"
            class="input"
            placeholder="Количество"
            min="1"
            required
        />

        <label for="cost">Стоимость</label>
        <input
            id="cost"
            type="number"
            v-model="currentPurchase.cost"
            class="input"
            placeholder="Стоимость"
            min="0"
            required
        />

        <label for="supplyDate">Дата поставки</label>
        <input
            id="supplyDate"
            type="date"
            v-model="currentPurchase.supply_date"
            class="input"
            required
        />

        <button @click="savePurchase" class="btn">Сохранить</button>
        <button @click="closePurchaseModal" class="btn danger">Отмена</button>
      </div>
    </div>

    <!-- Панель навигации в стиле карточек -->
    <div class="card-panel">
      <div class="card" @click="goToAdminHome">
        <h4>На главную</h4>
        <p>Панель администратора</p>
      </div>
      <div class="card" @click="goToEmployeesPage">
        <h4>Сотрудники</h4>
        <p>Управление сотрудниками</p>
      </div>
      <div class="card" @click="goToOrdersPage">
        <h4>Заказы</h4>
        <p>История заказов</p>
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
      currentPurchase: {}
    };
  },
  created() {
    this.fetchMaterials();
    this.fetchExpenditures();
    this.fetchPurchases();
  },
  methods: {
    // Получение списка материалов
    async fetchMaterials() {
      try {
        const response = await axios.get('http://localhost:8080/api/materials', {
          withCredentials: true
        });
        this.materials = response.data.materials;
      } catch (error) {
        console.error('Ошибка при загрузке материалов:', error);
        alert('Не удалось загрузить материалы.');
      }
    },

    // Получение списка расходов
    async fetchExpenditures() {
      try {
        const response = await axios.get('http://localhost:8080/api/expenditures', {
          withCredentials: true
        });
        this.expenditures = response.data.expenditures;
      } catch (error) {
        console.error('Ошибка при загрузке расходов:', error);
        alert('Не удалось загрузить расходы.');
      }
    },

    // Получение списка закупок
    async fetchPurchases() {
      try {
        const response = await axios.get('http://localhost:8080/api/purchases', {
          withCredentials: true
        });
        this.purchases = response.data.purchases;
      } catch (error) {
        console.error('Ошибка при загрузке закупок:', error);
        alert('Не удалось загрузить закупки.');
      }
    },

    // Открытие модального окна для добавления расхода
    openAddExpenditureModal() {
      this.showExpenditureModal = true;
      this.currentExpenditure = { material_id: '', expenditure_date: '', quantity: 0 };
    },

    // Открытие модального окна для добавления закупки
    openAddPurchaseModal() {
      this.showPurchaseModal = true;
      this.currentPurchase = { material_id: '', supplier: '', quantity: 0, cost: 0, supply_date: '' };
    },

    // Сохранение расхода
    async saveExpenditure() {
      if (this.currentExpenditure.quantity > 0) {
        try {
          const response = await axios.post('http://localhost:8080/api/expenditures', this.currentExpenditure, {
            withCredentials: true
          });
          this.expenditures.push(response.data.expenditure);
          this.closeExpenditureModal();
          alert('Расход успешно добавлен.');
        } catch (error) {
          console.error('Ошибка при добавлении расхода:', error);
          alert('Не удалось добавить расход.');
        }
      } else {
        alert('Количество расхода должно быть больше 0');
      }
    },

    // Сохранение закупки
    async savePurchase() {
      if (this.currentPurchase.quantity > 0 && this.currentPurchase.cost >= 0) {
        try {
          const response = await axios.post('http://localhost:8080/api/purchases', this.currentPurchase, {
            withCredentials: true
          });
          this.purchases.push(response.data.purchase);
          this.closePurchaseModal();
          alert('Закупка успешно добавлена.');
        } catch (error) {
          console.error('Ошибка при добавлении закупки:', error);
          alert('Не удалось добавить закупку.');
        }
      } else {
        alert('Количество должно быть больше 0, а стоимость не может быть отрицательной');
      }
    },

    // Закрытие модального окна расхода
    closeExpenditureModal() {
      this.showExpenditureModal = false;
      this.currentExpenditure = {};
    },

    // Закрытие модального окна закупки
    closePurchaseModal() {
      this.showPurchaseModal = false;
      this.currentPurchase = {};
    },

    // Форматирование даты
    formatDate(dateString) {
      const date = new Date(dateString);
      if (isNaN(date)) return 'Invalid Date';
      return date.toLocaleDateString('ru-RU', { year: 'numeric', month: 'long', day: 'numeric' });
    },

    // Навигация
    goToAdminHome() {
      this.$router.push({ name: 'AdminHome' });
    },
    goToEmployeesPage() {
      this.$router.push({ name: 'ManageEmp' });
    },
    goToOrdersPage() {
      this.$router.push({ name: 'OrderHistory' });
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

.navigation-buttons {
  display: flex;
  justify-content: space-between;
  margin-top: 20px;
}
</style>
