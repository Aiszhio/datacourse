<template>
  <div class="services-dashboard">
    <h2>Управление услугами и оборудованием</h2>

    <!-- Таблица с услугами -->
    <div class="table-section">
      <h3>Услуги</h3>
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
        <!-- Используем sortedServices для сортировки и индексы для нумерации -->
        <tr v-for="(service, index) in sortedServices" :key="service.service_id">
          <td>{{ index + 1 }}</td> <!-- Индекс + 1 для начала нумерации с 1 -->
          <td>{{ service.name }}</td>
          <td>{{ service.price }} ₽</td>
          <td>
            <button @click="openEditServiceModal(service)" class="btn">Редактировать</button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Таблица с оборудованием -->
    <div class="table-section">
      <h3>Оборудование</h3>
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
        <!-- Используем sortedEquipment для сортировки -->
        <tr v-for="(equipment, index) in sortedEquipment" :key="equipment?.equipment_id">
          <td>{{ index + 1 }}</td>
          <td>{{ equipment?.brand }}</td>
          <td>{{ equipment?.model }}</td>
          <td>
            <button @click="openEditEquipmentModal(equipment)" class="btn">Редактировать</button>
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
          <!-- Модальное окно для услуг -->
          <div v-if="modalType === 'service'">
            <label class="form-label">Название услуги:</label>
            <input v-model="currentItem.name" required class="form-input" maxlength="50" />

            <label class="form-label">Стоимость:</label>
            <div class="price-input">
              <input type="number" v-model="currentItem.price" required class="form-input" min="0" />
              <span>₽</span>
            </div>
          </div>
          <!-- Модальное окно для оборудования -->
          <div v-if="modalType === 'equipment'">
            <label class="form-label">Марка:</label>
            <input v-model="currentItem.brand" required class="form-input" />

            <label class="form-label">Модель:</label>
            <input v-model="currentItem.model" required class="form-input" />
          </div>

          <button type="submit" class="btn">Сохранить</button>
          <button @click.prevent="closeModal" type="button" class="btn danger">Отмена</button>
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
        <h4>Материалы</h4>
        <p>Материалы и расходы</p>
      </div>
      <div class="card" @click="goToBookingsPage">
        <h4>Бронирования</h4>
        <p>Управление бронированиями</p>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'Services',
  data() {
    return {
      services: [],
      equipmentList: [],
      showModal: false,
      modalType: '',
      modalTitle: '',
      currentItem: {}
    };
  },
  created() {
    this.fetchServicesAndEquipment();
  },
  computed: {
    sortedServices() {
      return this.services ? [...this.services].sort((a, b) => a.service_id - b.service_id) : [];
    },
    sortedEquipment() {
      return this.equipmentList ? [...this.equipmentList].sort((a, b) => a.equipment_id - b.equipment_id) : [];
    }
  },
  methods: {
    // Получение списка услуг и оборудования
    async fetchServicesAndEquipment() {
      try {
        const response = await axios.get('http://localhost:8080/api/services/admin', {
          withCredentials: true
        });
        this.services = response.data.services;
        this.equipmentList = response.data.equipment;
      } catch (error) {
        console.error('Ошибка при загрузке данных:', error);
        alert('Не удалось загрузить данные об услугах и оборудовании.');
      }
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
    goToMaterialsPage() {
      this.$router.push({ name: 'MaterialsOverview' });
    },
    goToBookingsPage() {
      this.$router.push({ name: 'Bookings' });
    },

    // Открытие модального окна для редактирования услуги
    openEditServiceModal(service) {
      this.modalType = 'service';
      this.modalTitle = 'Редактировать услугу';
      // Проверка service_id
      console.log("Service ID:", service.service_id);  // Убедитесь, что это значение корректно
      this.currentItem = { ...service };
      this.showModal = true;
    },
    // Открытие модального окна для редактирования оборудования
    openEditEquipmentModal(equipment) {
      this.modalType = 'equipment';
      this.modalTitle = 'Редактировать оборудование';
      this.currentItem = { ...equipment };
      this.showModal = true;
    },

    // Сохранение услуги или оборудования
    async saveItem() {
      if (this.modalType === 'service') {
        if (this.currentItem.service_id) {
          // Редактирование услуги
          try {
            const response = await axios.put(
                `http://localhost:8080/api/services/${this.currentItem.service_id}`,
                {
                  name: this.currentItem.name,
                  price: this.currentItem.price
                },
                { withCredentials: true }
            );
            // Обновляем услугу в списке
            const index = this.services.findIndex(s => s.service_id === this.currentItem.service_id);
            if (index !== -1) {
              this.services[index] = { ...response.data.service };
            }
            alert('Услуга успешно обновлена.');

            // Закрытие модального окна
            this.closeModal();
          } catch (error) {
            console.error('Ошибка при редактировании услуги:', error);
            if (error.response && error.response.data && error.response.data.error) {
              alert(`Ошибка: ${error.response.data.error}`);
            } else {
              alert('Не удалось обновить услугу.');
            }
          }
        }
      } else if (this.modalType === 'equipment') {
        if (this.currentItem.equipment_id) {
          // Редактирование оборудования
          try {
            const response = await axios.put(
                `http://localhost:8080/api/equipment/${this.currentItem.equipment_id}`,
                {
                  brand: this.currentItem.brand,
                  model: this.currentItem.model
                },
                { withCredentials: true }
            );
            const index = this.equipmentList.findIndex(e => e.equipment_id === this.currentItem.equipment_id);
            if (index !== -1) {
              this.equipmentList.splice(index, 1, response.data.equipment);
            }

            // Перезагружаем данные с сервера, если нужно
            await this.fetchServicesAndEquipment();  // Перезапросить все данные

            alert('Оборудование успешно обновлено.');
            this.closeModal(); // Закрытие модального окна

          } catch (error) {
            console.error('Ошибка при редактировании оборудования:', error);
            alert('Не удалось обновить оборудование.');
          }
        }
      }
    },
    // Закрытие модального окна
    closeModal() {
      this.showModal = false;
      this.modalType = '';
      this.modalTitle = '';
      this.currentItem = {};
    },
  }
};
</script>

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

