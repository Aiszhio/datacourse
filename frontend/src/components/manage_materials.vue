<template>
  <div class="materials-overview">
    <h2>Управление материалами</h2>

    <!-- Таблица закупок материалов с кнопками управления -->
    <div class="purchases-section">
      <h3>История закупок</h3>
      <button @click="openAddPurchaseModal" class="btn">Добавить закупку</button>
      <button @click="sortPurchases" class="btn">Сортировать</button>

      <table>
        <thead>
        <tr>
          <th>ID Закупки</th>
          <th>Материал</th>
          <th>Поставщик</th>
          <th>Количество</th>
          <th>Стоимость</th>
          <th>Дата закупки</th>
          <th>Действия</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="purchase in purchases" :key="purchase.id">
          <td>{{ purchase.id }}</td>
          <td>{{ purchase.material }}</td>
          <td>{{ purchase.supplier }}</td>
          <td>{{ purchase.quantity }}</td>
          <td>{{ purchase.cost }}</td>
          <td>{{ formatDate(purchase.date) }}</td>
          <td>
            <button @click="openEditPurchaseModal(purchase)" class="btn">Редактировать</button>
            <button @click="deletePurchase(purchase.id)" class="btn danger">Удалить</button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Модальные окна для добавления и редактирования закупок -->
    <div v-if="showPurchaseModal" class="modal-overlay">
      <div class="modal">
        <h3>{{ editingPurchase ? "Редактировать закупку" : "Добавить закупку" }}</h3>
        <input v-model="currentPurchase.material" class="input" placeholder="Материал" required />
        <input v-model="currentPurchase.supplier" class="input" placeholder="Поставщик" required />
        <input type="number" v-model="currentPurchase.quantity" class="input" placeholder="Количество" required />
        <input type="number" v-model="currentPurchase.cost" class="input" placeholder="Стоимость" required />
        <input type="date" v-model="currentPurchase.date" class="input" required />
        <button @click="savePurchase" class="btn">Сохранить</button>
        <button @click="closePurchaseModal" class="btn danger">Отмена</button>
      </div>
    </div>

    <!-- Кнопки навигации по другим страницам -->
    <div class="navigation-buttons">
      <button @click="goToAdminHome" class="btn">На главную администратора</button>
      <button @click="goToEmployeesPage" class="btn">Управление сотрудниками</button>
      <button @click="goToOrdersPage" class="btn">История заказов</button>
    </div>

  </div>
</template>

<script>
export default {
  name: 'MaterialsOverview',
  data() {
    return {
      purchases: [],
      expenditures: [],
      showPurchaseModal: false,
      showExpenditureModal: false,
      editingPurchase: false,
      editingExpenditure: false,
      currentPurchase: {},
      currentExpenditure: {}
    };
  },
  mounted() {
    this.fetchPurchases();
    this.fetchExpenditures();
  },
  methods: {
    async fetchPurchases() {
      try {
        const response = await fetch('http://localhost:8080/api/materials/purchases');
        this.purchases = await response.json();
      } catch (error) {
        console.error('Ошибка при загрузке закупок:', error);
      }
    },
    async fetchExpenditures() {
      try {
        const response = await fetch('http://localhost:8080/api/materials/expenditures');
        this.expenditures = await response.json();
      } catch (error) {
        console.error('Ошибка при загрузке расхода материалов:', error);
      }
    },
    openAddPurchaseModal() {
      this.showPurchaseModal = true;
      this.editingPurchase = false;
      this.currentPurchase = { material: '', supplier: '', quantity: 0, cost: 0, date: '' };
    },
    openEditPurchaseModal(purchase) {
      this.showPurchaseModal = true;
      this.editingPurchase = true;
      this.currentPurchase = { ...purchase };
    },
    deletePurchase(id) {
      this.purchases = this.purchases.filter(p => p.id !== id);
    },
    savePurchase() {
      if (this.editingPurchase) {
        const index = this.purchases.findIndex(p => p.id === this.currentPurchase.id);
        this.purchases[index] = { ...this.currentPurchase };
      } else {
        this.purchases.push({ ...this.currentPurchase, id: Date.now() });
      }
      this.closePurchaseModal();
    },
    sortPurchases() {
      this.purchases.sort((a, b) => a.material.localeCompare(b.material));
    },
    closePurchaseModal() {
      this.showPurchaseModal = false;
      this.currentPurchase = {};
    },
    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleDateString('ru-RU', { year: 'numeric', month: 'long', day: 'numeric' });
    },
    goToAdminHome(){
      this.$router.push({ name: 'AdminHome' });
    },
    goToEmployeesPage(){
      this.$router.push({ name: 'ManageEmp' });
    },
    goToOrdersPage(){
      this.$router.push({ name: 'OrderHistory' });
    },
  }
};
</script>

<style scoped>
.purchases-section{
  margin-bottom: 5vh;
}

.navigation-buttons{
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}

.materials-overview {
  padding: 20px;
  max-width: 900px;
  margin: 0 auto;
}

h2 {
  text-align: center;
  margin-bottom: 20px;
}

table {
  width: 100%;
  margin-top: 10px;
  border-collapse: collapse;
}

table th,
table td {
  padding: 12px;
  border: 1px solid #ddd;
  text-align: left;
}

.btn {
  background-color: #4CAF50;
  color: white;
  padding: 10px 15px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 1em;
  transition: background-color 0.3s ease;
  margin-right: 5px;
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
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.input {
  width: 100%;
  padding: 12px;
  margin: 8px 0;
  border: 1px solid #ddd;
  border-radius: 5px;
  box-sizing: border-box;
  font-size: 1em;
  transition: border-color 0.3s ease;
}

.input:focus {
  border-color: #4CAF50;
  outline: none;
}

.input[type="date"] {
  font-family: inherit;
  color: #555;
}
</style>