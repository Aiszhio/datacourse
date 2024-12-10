<template>
  <div class="services-dashboard">
    <h2>Управление услугами и оборудованием</h2>

    <!-- Таблица с услугами -->
    <div class="table-section">
      <h3>Услуги</h3>
      <div class="add-service-button">
        <button @click="openAddServiceModal" class="btn add">Добавить услугу</button>
      </div>
      <table class="data-table">
        <thead>
        <tr>
          <th>Название</th>
          <th>Стоимость</th>
          <th>Требуемое оборудование</th>
          <th>Действия</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(service) in sortedServices" :key="service.service_id">
          <td>{{ service.name }}</td>
          <td>{{ service.price }} ₽</td>
          <td>
            <button @click="toggleEquipment(service.service_id)" class="btn">
              {{ isExpanded(service.service_id) ? 'Скрыть' : 'Показать' }}
            </button>
            <div v-if="isExpanded(service.service_id)">
              {{ formatEquipmentList(service.RequiredEquipment) }}
            </div>
          </td>
          <td>
            <button @click="openEditServiceModal(service)" class="btn">Редактировать</button>
            <button @click="openDeleteServiceModal(service)" class="btn danger">Удалить</button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Таблица с оборудованием -->
    <div class="table-section">
      <h3>Оборудование</h3>
      <div class="add-equipment-button">
        <button @click="openAddEquipmentModal" class="btn add">Добавить оборудование</button>
      </div>
      <table class="data-table">
        <thead>
        <tr>
          <th>Тип оборудования</th>
          <th>Марка</th>
          <th>Модель</th>
          <th>Действия</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(equipment) in sortedEquipment" :key="equipment.equipment_id">
          <td>{{ equipment.type }}</td>
          <td>{{ equipment.brand }}</td>
          <td>{{ equipment.model }}</td>
          <td>
            <button @click="openEditEquipmentModal(equipment)" class="btn">Редактировать</button>
            <button @click="openSubtractEquipmentModal(equipment)" class="btn danger">Списать</button>
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

            <label class="form-label">Требуемое оборудование:</label>
            <select v-model="currentItem.requiredEquipment" multiple class="form-input">
              <option v-for="equipment in equipmentList" :key="equipment.equipment_id" :value="equipment.equipment_id">
                {{ equipment.type }} - {{ equipment.brand }} {{ equipment.model }}
              </option>
            </select>
          </div>

          <!-- Модальное окно для оборудования -->
          <div v-if="modalType === 'equipment'">
            <label class="form-label">Тип оборудования:</label>
            <select v-model="currentItem.type" required class="form-input">
              <option disabled value="">Выберите тип</option>
              <option v-for="type in equipmentTypes" :key="type" :value="type">{{ type }}</option>
            </select>

            <label class="form-label">Марка:</label>
            <input v-model="currentItem.brand" required class="form-input" />

            <label class="form-label">Модель:</label>
            <input v-model="currentItem.model" required class="form-input" />
          </div>

          <!-- Модальное окно для списания оборудования -->
          <div v-if="modalType === 'subtractEquipment'">
            <p>Вы уверены, что хотите списать оборудование:</p>
            <p><strong>{{ currentItem.type }} - {{ currentItem.brand }} {{ currentItem.model }}</strong></p>
          </div>

          <!-- Кнопки действия -->
          <div v-if="modalType !== 'subtractEquipment'">
            <button type="submit" class="btn">Сохранить</button>
          </div>
          <div v-else>
            <button @click.prevent="subtractEquipment" type="button" class="btn danger">Списать</button>
          </div>
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
  name: 'Services',
  data() {
    return {
      services: [], // Группированные услуги с оборудованием
      equipmentList: [], // Список всех оборудования
      showModal: false,
      modalType: '',
      modalTitle: '',
      currentItem: {},
      equipmentTypes: ['Принтер', 'Фотокамера', 'Видеокамера'], // Список типов оборудования
      expandedServiceIds: [] // Для отслеживания развёрнутых услуг
    };
  },
  created() {
    this.fetchServicesAndEquipment();
  },
  computed: {
    sortedServices() {
      return this.services
          ? [...this.services].sort((a, b) => a.service_id - b.service_id) // Используйте правильные поля
          : [];
    },
    sortedEquipment() {
      return this.equipmentList
          ? [...this.equipmentList].sort((a, b) => a.equipment_id - b.equipment_id) // Используйте правильные поля
          : [];
    }
  },
  methods: {
    // Получение списка услуг и оборудования
    async fetchServicesAndEquipment() {
      try {
        const response = await axios.get('http://localhost:8080/api/services/admin', {
          withCredentials: true
        });
        this.services = response.data.services || [];
        this.equipmentList = response.data.equipment || [];
        console.log('Полученные услуги:', this.services);
        console.log('Список оборудования:', this.equipmentList);
        // this.calculateTotalPages(); // Удалите или раскомментируйте, если добавите метод
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
    goToClients() {
      this.$router.push({ name: 'Clients' })
    },

    // Открытие модального окна для добавления услуги
    openAddServiceModal() {
      this.modalType = 'service';
      this.modalTitle = 'Добавить услугу';
      this.currentItem = {
        name: '',
        price: '',
        requiredEquipment: [],
      };
      this.showModal = true;
    },

    // Открытие модального окна для редактирования услуги
    openEditServiceModal(service) {
      this.modalType = 'service';
      this.modalTitle = 'Редактировать услугу';
      // Копируем объект, чтобы избежать прямого изменения оригинала
      this.currentItem = {
        ...service,
        requiredEquipment: service.RequiredEquipment.map(eq => eq.equipment_id) // Исправлено на 'RequiredEquipment'
      };
      this.showModal = true;
    },

    // Открытие модального окна для редактирования оборудования
    openEditEquipmentModal(equipment) {
      this.modalType = 'equipment';
      this.modalTitle = 'Редактировать оборудование';
      this.currentItem = { ...equipment };
      this.showModal = true;
    },

    // Открытие модального окна для добавления оборудования
    openAddEquipmentModal() {
      this.modalType = 'equipment';
      this.modalTitle = 'Добавить оборудование';
      this.currentItem = {
        type: '',
        brand: '',
        model: '',
      };
      this.showModal = true;
    },

    // Открытие модального окна для списания оборудования
    openSubtractEquipmentModal(equipment) {
      this.modalType = 'subtractEquipment';
      this.modalTitle = 'Списать оборудование';
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
                  price: this.currentItem.price,
                  requiredEquipment: this.currentItem.requiredEquipment, // Массив equipment_id
                },
                { withCredentials: true }
            );
            // Обновляем услугу в списке
            const index = this.services.findIndex(s => s.service_id === this.currentItem.service_id);
            if (index !== -1) {
              this.services[index] = response.data.service;
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
        } else {
          // Добавление новой услуги
          try {
            const response = await axios.post(
                'http://localhost:8080/api/services/admin',
                {
                  name: this.currentItem.name,
                  price: this.currentItem.price,
                  requiredEquipment: this.currentItem.requiredEquipment, // Массив equipment_id
                },
                { withCredentials: true }
            );
            // Добавляем новую услугу в список
            this.services.push(response.data.service);
            alert('Услуга успешно добавлена.');

            // Закрытие модального окна
            this.closeModal();
          } catch (error) {
            console.error('Ошибка при добавлении услуги:', error);
            if (error.response && error.response.data && error.response.data.error) {
              alert(`Ошибка: ${error.response.data.error}`);
            } else {
              alert('Не удалось добавить услугу.');
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
                  type: this.currentItem.type,
                  brand: this.currentItem.brand,
                  model: this.currentItem.model,
                },
                { withCredentials: true }
            );
            const index = this.equipmentList.findIndex(e => e.equipment_id === this.currentItem.equipment_id);
            if (index !== -1) {
              this.equipmentList.splice(index, 1, response.data.equipment);
            }

            alert('Оборудование успешно обновлено.');
            this.closeModal(); // Закрытие модального окна

          } catch (error) {
            console.error('Ошибка при редактировании оборудования:', error);
            if (error.response && error.response.data && error.response.data.error) {
              alert(`Ошибка: ${error.response.data.error}`);
            } else {
              alert('Не удалось обновить оборудование.');
            }
          }
        } else {
          // Добавление нового оборудования
          try {
            const response = await axios.post(
                'http://localhost:8080/api/equipment/admin',
                {
                  type: this.currentItem.type,
                  brand: this.currentItem.brand,
                  model: this.currentItem.model,
                },
                { withCredentials: true }
            );
            // Добавляем новое оборудование в список
            this.equipmentList.push(response.data.equipment);
            alert('Оборудование успешно добавлено.');

            // Закрытие модального окна
            this.closeModal();
          } catch (error) {
            console.error('Ошибка при добавлении оборудования:', error);
            if (error.response && error.response.data && error.response.data.error) {
              alert(`Ошибка: ${error.response.data.error}`);
            } else {
              alert('Не удалось добавить оборудование.');
            }
          }
        }
      } else if (this.modalType === 'subtractEquipment') {
        // Списание оборудования
        try {
          const equipmentID = this.currentItem.equipment_id;

          // Отправка запроса на списание оборудования
          const response = await axios.post(
              `http://localhost:8080/api/equipment/${equipmentID}/subtract`,
              {},
              { withCredentials: true }
          );

          if (response.data.message === 'Оборудование успешно списано') {
            // Обновление списка оборудования
            this.equipmentList = this.equipmentList.filter(e => e.equipment_id !== equipmentID);
            alert('Оборудование успешно списано.');
            this.closeModal();
          } else {
            alert('Неизвестная ошибка: ' + (response.data.message || 'Попробуйте позже.'));
          }
        } catch (error) {
          console.error('Ошибка при списании оборудования:', error);
          if (error.response && error.response.data && error.response.data.error) {
            alert(`Ошибка: ${error.response.data.error}`);
          } else {
            alert('Не удалось списать оборудование.');
          }
        }
      } else if (this.modalType === 'deleteService') {
        // Удаление услуги
        await this.deleteService();
      }
    },

    // Закрытие модального окна
    closeModal() {
      this.showModal = false;
      this.modalType = '';
      this.modalTitle = '';
      this.currentItem = {};
    },

    // Разворачивание/Скрытие оборудования
    toggleEquipment(service_id) {
      const index = this.expandedServiceIds.indexOf(service_id);
      if (index === -1) {
        this.expandedServiceIds.push(service_id);
      } else {
        this.expandedServiceIds.splice(index, 1);
      }
    },
    isExpanded(service_id) {
      return this.expandedServiceIds.includes(service_id);
    },
    formatEquipmentList(requiredEquipment) {
      console.log('formatEquipmentList called with:', requiredEquipment);
      if (!requiredEquipment || requiredEquipment.length === 0) {
        return 'Нет оборудования';
      }
      return requiredEquipment.map(eq => `${eq.type} - ${eq.brand} ${eq.model}`).join(', ');
    },

    // Удаление услуги
    async openDeleteServiceModal(service) {
      this.modalType = 'deleteService';
      this.modalTitle = 'Удалить услугу';
      this.currentItem = { ...service };
      this.showModal = true;
    },
    async deleteService() {
      try {
        const response = await axios.delete(
            `http://localhost:8080/api/services/${this.currentItem.service_id}`,
            { withCredentials: true }
        );
        if (response.data.message === 'Услуга успешно удалена') {
          this.services = this.services.filter(s => s.service_id !== this.currentItem.service_id);
          alert('Услуга успешно удалена.');
          this.closeModal();
        } else {
          alert('Неизвестная ошибка: ' + (response.data.message || 'Попробуйте позже.'));
        }
      } catch (error) {
        console.error('Ошибка при удалении услуги:', error);
        if (error.response && error.response.data && error.response.data.error) {
          alert(`Ошибка: ${error.response.data.error}`);
        } else {
          alert('Не удалось удалить услугу.');
        }
      }
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

.btn.add {
  background-color: #2196F3;
}

.btn:hover {
  background-color: #45a049;
}

.btn.danger:hover {
  background-color: #e53935;
}

.btn.add:hover {
  background-color: #1976D2;
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
  max-width: 500px;
  width: 100%;
}

.form-label {
  display: block;
  font-weight: bold;
  margin-bottom: 5px;
  color: #333;
}

.form-input {
  width: 100%;
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

.price-input {
  display: flex;
  align-items: center;
}

.price-input input {
  flex: 1;
}

.price-input span {
  margin-left: 5px;
}
</style>
