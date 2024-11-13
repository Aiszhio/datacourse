<script>
export default {
  name: 'Bookings',
  data() {
    return {
      services: [
        { id: 1, name: 'Фотосессия', price: 3000 },
        { id: 2, name: 'Видеосъемка', price: 5000 },
        { id: 3, name: 'Обработка фото', price: 1500 }
      ],
      requirements: [
        { id: 1, serviceId: 1, equipmentId: 1 },
        { id: 2, serviceId: 2, equipmentId: 2 },
        { id: 3, serviceId: 3, equipmentId: 1 }
      ],
      equipmentList: [
        { id: 1, brand: 'Canon', model: 'EOS 5D Mark IV' },
        { id: 2, brand: 'Sony', model: 'Alpha a7 III' },
        { id: 3, brand: 'Nikon', model: 'D850' }
      ],
      showModal: false,
      modalType: '',
      modalTitle: '',
      currentItem: {}
    };
  },
  methods: {
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
    goToBookingsPage() {
      this.$router.push({ name: 'Bookings' });
    },
    openAddServiceModal() {
      this.modalType = 'service';
      this.modalTitle = 'Добавить услугу';
      this.currentItem = { name: '', price: '' };
      this.showModal = true;
    },
    openEditServiceModal(service) {
      this.modalType = 'service';
      this.modalTitle = 'Редактировать услугу';
      this.currentItem = { ...service };
      this.showModal = true;
    },
    openAddRequirementModal() {
      this.modalType = 'requirement';
      this.modalTitle = 'Добавить требование к услуге';
      this.currentItem = { serviceId: '', equipmentId: '' };
      this.showModal = true;
    },
    openAddEquipmentModal() {
      this.modalType = 'equipment';
      this.modalTitle = 'Добавить оборудование';
      this.currentItem = { brand: '', model: '' };
      this.showModal = true;
    },
    openEditEquipmentModal(equipment) {
      this.modalType = 'equipment';
      this.modalTitle = 'Редактировать оборудование';
      this.currentItem = { ...equipment };
      this.showModal = true;
    },
    saveItem() {
      if (this.modalType === 'service') {
        // Логика сохранения услуги
      } else if (this.modalType === 'requirement') {
        // Логика сохранения требования
      } else if (this.modalType === 'equipment') {
        // Логика сохранения оборудования
      }
      this.closeModal();
    },
    deleteService(id) {
      this.services = this.services.filter(service => service.id !== id);
    },
    deleteRequirement(id) {
      this.requirements = this.requirements.filter(req => req.id !== id);
    },
    deleteEquipment(id) {
      this.equipmentList = this.equipmentList.filter(eq => eq.id !== id);
    },
    closeModal() {
      this.showModal = false;
      this.currentItem = {};
    }
  }
};
</script>

<template>
  <div class="services-dashboard">
    <h2>Управление услугами</h2>
    <!-- Таблица с услугами -->
    <div class="table-section">
      <h3>Услуги</h3>
      <button @click="openAddServiceModal" class="btn">Добавить услугу</button>
      <table class="data-table">
        <thead>
        <tr>
          <th>Номер услуги</th>
          <th>Название</th>
          <th>Стоимость</th>
          <th>Действия</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="service in services" :key="service.id">
          <td>{{ service.id }}</td>
          <td>{{ service.name }}</td>
          <td>{{ service.price }}</td>
          <td>
            <button @click="openEditServiceModal(service)" class="btn">Редактировать</button>
            <button @click="deleteService(service.id)" class="btn danger">Удалить</button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Таблица с оборудованием -->
    <div class="table-section">
      <h3>Оборудование</h3>
      <button @click="openAddEquipmentModal" class="btn">Добавить оборудование</button>
      <table class="data-table">
        <thead>
        <tr>
          <th>Номер оборудования</th>
          <th>Марка</th>
          <th>Модель</th>
          <th>Действия</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="equipment in equipmentList" :key="equipment.id">
          <td>{{ equipment.id }}</td>
          <td>{{ equipment.brand }}</td>
          <td>{{ equipment.model }}</td>
          <td>
            <button @click="openEditEquipmentModal(equipment)" class="btn">Редактировать</button>
            <button @click="deleteEquipment(equipment.id)" class="btn danger">Удалить</button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Модальные окна для добавления и редактирования -->
    <div v-if="showModal" class="modal-overlay">
      <div class="modal">
        <h3>{{ modalTitle }}</h3>
        <form @submit.prevent="saveItem">
          <div v-if="modalType === 'service'">
            <label class="form-label">Название услуги:</label>
            <input v-model="currentItem.name" required class="form-input" />
            <label class="form-label">Стоимость:</label>
            <input type="number" v-model="currentItem.price" required class="form-input" />
          </div>
          <div v-if="modalType === 'requirement'">
            <label class="form-label">Номер услуги:</label>
            <input type="number" v-model="currentItem.serviceId" required class="form-input" />
            <label class="form-label">Номер оборудования:</label>
            <input type="number" v-model="currentItem.equipmentId" required class="form-input" />
          </div>
          <div v-if="modalType === 'equipment'">
            <label class="form-label">Марка:</label>
            <input v-model="currentItem.brand" required class="form-input" />
            <label class="form-label">Модель:</label>
            <input v-model="currentItem.model" required class="form-input" />
          </div>
          <button type="submit" class="btn">Сохранить</button>
          <button @click="closeModal" class="btn danger">Отмена</button>
        </form>
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
      <div class="card" @click="goToMaterialsPage">
        <h4>Закупки</h4>
        <p>Материалы и расходы</p>
      </div>
      <div class="card" @click="goToBookingsPage">
        <h4>Бронирования</h4>
        <p>Управление бронированиями</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.services-dashboard {
  padding: 20px;
  max-width: 100vh;
  max-height: 100vh;
  margin: 0 auto;
}

.card-panel {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
  justify-content: center;
}

.card {
  background-color: #4CAF50;
  color: #fff;
  border-radius: 10px;
  padding: 10px;
  margin-bottom: 5px;
  width: 180px;
  max-height: 15vh;
  text-align: center;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  cursor: pointer;
  transition: transform 0.3s, box-shadow 0.3s;
}

.card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
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
}

.btn {
  background-color: #4CAF50;
  color: white;
  padding: 8px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin: 5px;
  transition: background-color 0.3s;
}

.btn.danger {
  background-color: #f44336;
}

.btn:hover {
  background-color: #45a049;
}

.btn.danger:hover {
  background-color: #e53935;
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

.form-label {
  display: block;
  font-weight: bold;
  margin-bottom: 5px;
  color: #333;
}

.form-input {
  width: 90%;
  padding: 10px;
  margin-bottom: 15px;
  border: 1px solid #ddd;
  border-radius: 4px;
  transition: border-color 0.3s;
}

.form-input:focus {
  border-color: #4CAF50;
  outline: none;
}
</style>

