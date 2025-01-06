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
        <!-- Использование уникального ключа, так как equipment_id всегда присутствует -->
        <tr v-for="equipment in sortedEquipment" :key="equipment.equipment_id">
          <td>{{ equipment.type }}</td>
          <td>{{ equipment.brand }}</td>
          <td>{{ equipment.model }}</td>
          <td>
            <button @click="openEditEquipmentModal(equipment)" class="btn">Редактировать</button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- Модальное окно для добавления/редактирования услуги -->
    <b-modal
        id="service-modal"
        v-model="showServiceModal"
        :title="modalTitle"
        hide-footer
        @hide="resetServiceForm"
    >
      <form @submit.prevent="saveService">
        <div class="mb-3">
          <label for="serviceName" class="form-label">Название услуги:</label>
          <input
              id="serviceName"
              type="text"
              v-model="currentService.name"
              class="form-control"
              placeholder="Введите название услуги"
              maxlength="50"
              required
          />
        </div>

        <div class="mb-3">
          <label for="servicePrice" class="form-label">Стоимость:</label>
          <div class="input-group">
            <input
                id="servicePrice"
                type="number"
                v-model="currentService.price"
                class="form-control"
                placeholder="Введите стоимость"
                min="0"
                required
            />
            <span class="input-group-text">₽</span>
          </div>
        </div>

        <div class="mb-3">
          <label for="requiredEquipment" class="form-label">Требуемое оборудование:</label>
          <select
              id="requiredEquipment"
              v-model="currentService.requiredEquipment"
              multiple
              class="form-select"
          >
            <option
                v-for="equipment in equipmentList"
                :key="equipment.equipment_id"
                :value="equipment.equipment_id"
            >
              {{ equipment.type }} - {{ equipment.brand }} {{ equipment.model }}
            </option>
          </select>
        </div>

        <div class="d-flex justify-content-end">
          <button type="submit" class="btn btn-primary me-2">Сохранить</button>
          <button type="button" @click="closeServiceModal" class="btn btn-danger">Отмена</button>
        </div>
      </form>
    </b-modal>

    <!-- Модальное окно для списания оборудования -->
    <b-modal
        id="subtract-equipment-modal"
        v-model="showSubtractEquipmentModal"
        title="Списать оборудование"
        hide-footer
        @hide="resetSubtractEquipmentForm"
    >
      <div>
        <p>Вы уверены, что хотите списать оборудование:</p>
        <p><strong>{{ currentEquipment.type }} - {{ currentEquipment.brand }} {{ currentEquipment.model }}</strong></p>
        <div class="d-flex justify-content-end mt-3">
          <button @click="subtractEquipment" class="btn btn-danger me-2">Списать</button>
          <button @click="closeSubtractEquipmentModal" class="btn btn-secondary">Отмена</button>
        </div>
      </div>
    </b-modal>

    <!-- Модальное окно для добавления/редактирования оборудования -->
    <b-modal
        id="equipment-modal"
        v-model="showEquipmentModal"
        :title="equipmentModalTitle"
        hide-footer
        @hide="resetEquipmentForm"
    >
      <form @submit.prevent="saveEquipment">
        <div class="mb-3">
          <label for="equipmentType" class="form-label">Тип оборудования:</label>
          <select
              id="equipmentType"
              v-model="currentEquipment.type"
              class="form-control"
              required
          >
            <option disabled value="">Выберите тип</option>
            <option v-for="type in equipmentTypes" :key="type" :value="type">{{ type }}</option>
          </select>
        </div>

        <div class="mb-3">
          <label for="equipmentBrand" class="form-label">Марка:</label>
          <input
              id="equipmentBrand"
              type="text"
              v-model="currentEquipment.brand"
              class="form-control"
              placeholder="Введите марку оборудования"
              required
          />
        </div>

        <div class="mb-3">
          <label for="equipmentModel" class="form-label">Модель:</label>
          <input
              id="equipmentModel"
              type="text"
              v-model="currentEquipment.model"
              class="form-control"
              placeholder="Введите модель оборудования"
              required
          />
        </div>

        <div class="d-flex justify-content-end">
          <button type="submit" class="btn btn-primary me-2">Сохранить</button>
          <button type="button" @click="closeEquipmentModal" class="btn btn-danger">Отмена</button>
        </div>
      </form>
    </b-modal>

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
import dayjs from 'dayjs';
import 'dayjs/locale/ru';
import Swal from 'sweetalert2';

dayjs.locale('ru'); // Устанавливаем русскую локаль

export default {
  name: 'Services',
  data() {
    return {
      services: [], // Список услуг
      equipmentList: [], // Список оборудования
      showServiceModal: false, // Модальное окно для услуг
      showSubtractEquipmentModal: false, // Модальное окно для списания оборудования
      showEquipmentModal: false, // Модальное окно для оборудования
      modalTitle: '',
      equipmentModalTitle: '',
      currentService: {},
      currentEquipment: {},
      equipmentTypes: ['Принтер', 'Фотокамера', 'Видеокамера'], // Типы оборудования
      expandedServiceIds: [], // Отслеживание развёрнутых услуг
    };
  },
  computed: {
    sortedServices() {
      return this.services
          ? [...this.services].sort((a, b) => a.service_id - b.service_id)
          : [];
    },
    sortedEquipment() {
      return this.equipmentList
          ? [...this.equipmentList].sort((a, b) => a.equipment_id - b.equipment_id)
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
      } catch (error) {
        console.error('Ошибка при загрузке данных:', error);
        this.showAlert('Не удалось загрузить данные об услугах и оборудовании.', 'error');
      }
    },

    // Открытие модального окна для добавления услуги
    openAddServiceModal() {
      console.log('Открытие модального окна для добавления услуги');
      this.modalTitle = 'Добавить услугу';
      this.currentService = {
        name: '',
        price: '',
        requiredEquipment: [],
      };
      this.showServiceModal = true;
    },

    // Открытие модального окна для редактирования услуги
    openEditServiceModal(service) {
      console.log('Открытие модального окна для редактирования услуги:', service);
      this.modalTitle = 'Редактировать услугу';
      // Копируем объект, чтобы избежать прямого изменения оригинала
      this.currentService = {
        ...service,
        requiredEquipment: service.RequiredEquipment.map(eq => eq.equipment_id),
      };
      this.showServiceModal = true;
    },

    // Открытие модального окна для списания оборудования
    openSubtractEquipmentModal(equipment) {
      console.log('Открытие модального окна для списания оборудования:', equipment);
      this.currentEquipment = { ...equipment };
      this.showSubtractEquipmentModal = true;
    },

    // Открытие модального окна для добавления оборудования
    openAddEquipmentModal() {
      console.log('Открытие модального окна для добавления оборудования');
      this.equipmentModalTitle = 'Добавить оборудование';
      this.currentEquipment = {
        type: '',
        brand: '',
        model: '',
      };
      this.showEquipmentModal = true;
    },

    // Открытие модального окна для редактирования оборудования
    openEditEquipmentModal(equipment) {
      console.log('Открытие модального окна для редактирования оборудования:', equipment);
      this.equipmentModalTitle = 'Редактировать оборудование';
      this.currentEquipment = { ...equipment };
      this.showEquipmentModal = true;
    },

    // Сохранение услуги (добавление или редактирование)
    async saveService() {
      if (this.currentService.name.trim() === '' || this.currentService.price === '') {
        this.showAlert('Пожалуйста, заполните все обязательные поля.', 'warning');
        return;
      }

      try {
        if (this.currentService.service_id) {
          // Редактирование услуги
          const response = await axios.put(
              `http://localhost:8080/api/services/${this.currentService.service_id}`,
              {
                name: this.currentService.name,
                price: this.currentService.price,
                requiredEquipment: this.currentService.requiredEquipment, // Массив equipment_id
              },
              { withCredentials: true }
          );

          if (response.data.message) {
            this.showAlert(response.data.message, 'success');
          } else {
            this.showAlert('Услуга успешно обновлена.', 'success');
          }

          // Обновляем услугу в списке
          this.services = this.services.map(service =>
              service.service_id === this.currentService.service_id ? response.data : service
          );
        } else {
          // Добавление новой услуги
          const response = await axios.post(
              'http://localhost:8080/api/services/admin',
              {
                name: this.currentService.name,
                price: this.currentService.price,
                requiredEquipment: this.currentService.requiredEquipment, // Массив equipment_id
              },
              { withCredentials: true }
          );

          if (response.data.message) {
            this.showAlert(response.data.message, 'success');
          } else {
            this.showAlert('Услуга успешно добавлена.', 'success');
          }

          // Добавляем новую услугу в список
          this.services.push(response.data);
        }

        // Закрытие модального окна
        this.closeServiceModal();
      } catch (error) {
        console.error('Ошибка при сохранении услуги:', error);
        if (error.response && error.response.data && error.response.data.error) {
          this.showAlert(`Ошибка: ${error.response.data.error}`, 'error');
        } else {
          this.showAlert('Не удалось сохранить услугу.', 'error');
        }
      }
    },

    // Сохранение оборудования (добавление или редактирование)
    async saveEquipment() {
      if (this.currentEquipment.type.trim() === '' ||
          this.currentEquipment.brand.trim() === '' ||
          this.currentEquipment.model.trim() === '') {
        this.showAlert('Пожалуйста, заполните все обязательные поля.', 'warning');
        return;
      }

      try {
        if (this.currentEquipment.equipment_id) {
          // Редактирование оборудования
          const response = await axios.put(
              `http://localhost:8080/api/equipment/${this.currentEquipment.equipment_id}`,
              {
                type: this.currentEquipment.type,
                brand: this.currentEquipment.brand,
                model: this.currentEquipment.model,
              },
              { withCredentials: true }
          );

          if (response.data.message) {
            this.showAlert(response.data.message, 'success');
          } else {
            this.showAlert('Оборудование успешно обновлено.', 'success');
          }

          // Обновляем оборудование в списке
          this.equipmentList = this.equipmentList.map(equipment =>
              equipment.equipment_id === this.currentEquipment.equipment_id ? response.data : equipment
          );
        } else {
          // Добавление нового оборудования
          const response = await axios.post(
              'http://localhost:8080/api/equipment/admin',
              {
                type: this.currentEquipment.type,
                brand: this.currentEquipment.brand,
                model: this.currentEquipment.model,
              },
              { withCredentials: true }
          );

          if (response.data.message) {
            this.showAlert(response.data.message, 'success');
          } else {
            this.showAlert('Оборудование успешно добавлено.', 'success');
          }

          // Добавляем новое оборудование в список
          this.equipmentList.push(response.data);
        }

        // Закрытие модального окна
        this.closeEquipmentModal();
      } catch (error) {
        console.error('Ошибка при сохранении оборудования:', error);
        if (error.response && error.response.data && error.response.data.error) {
          this.showAlert(`Ошибка: ${error.response.data.error}`, 'error');
        } else {
          this.showAlert('Не удалось сохранить оборудование.', 'error');
        }
      }
    },

    // Списание оборудования
    async subtractEquipment() {
      try {
        console.log('Списание оборудования:', this.currentEquipment);
        const equipmentID = this.currentEquipment.equipment_id;

        // Отправка DELETE запроса на сервер
        const response = await axios.delete(
            `http://localhost:8080/api/equipment/${equipmentID}`,
            { withCredentials: true }
        );

        if (response.data.message) {
          this.showAlert(response.data.message, 'success');
        } else {
          this.showAlert('Оборудование успешно списано.', 'success');
        }

        // Обновление списка оборудования в UI
        this.equipmentList = this.equipmentList.filter(e => e.equipment_id !== equipmentID);
        this.closeSubtractEquipmentModal();
      } catch (error) {
        console.error('Ошибка при списании оборудования:', error);
        if (error.response && error.response.data && error.response.data.error) {
          this.showAlert(`Ошибка: ${error.response.data.error}`, 'error');
        } else {
          this.showAlert('Не удалось списать оборудование.', 'error');
        }
      }
    },

    // Открытие модального окна для удаления услуги
    openDeleteServiceModal(service) {
      console.log('Открытие модального окна для удаления услуги:', service);
      Swal.fire({
        title: 'Вы уверены?',
        text: `Вы хотите удалить услугу "${service.name}"?`,
        icon: 'warning',
        showCancelButton: true,
        confirmButtonText: 'Да, удалить',
        cancelButtonText: 'Отмена',
        reverseButtons: true,
      }).then(async (result) => {
        if (result.isConfirmed) {
          await this.deleteService(service.service_id);
        } else if (result.dismiss === Swal.DismissReason.cancel) {
          this.showAlert('Удаление услуги отменено.', 'info');
        }
      });
    },

    // Удаление услуги
    async deleteService(service_id) {
      try {
        const response = await axios.delete(
            `http://localhost:8080/api/services/${service_id}`,
            { withCredentials: true }
        );

        if (response.data.message) {
          this.showAlert(response.data.message, 'success');
        } else {
          this.showAlert('Услуга успешно удалена.', 'success');
        }

        // Обновление списка услуг
        this.services = this.services.filter(service => service.service_id !== service_id);
      } catch (error) {
        console.error('Ошибка при удалении услуги:', error);
        if (error.response && error.response.data && error.response.data.error) {
          this.showAlert(`Ошибка: ${error.response.data.error}`, 'error');
        } else {
          this.showAlert('Не удалось удалить услугу.', 'error');
        }
      }
    },

    // Закрытие модального окна для услуг
    closeServiceModal() {
      this.showServiceModal = false;
    },

    // Закрытие модального окна для оборудования
    closeEquipmentModal() {
      this.showEquipmentModal = false;
    },

    // Закрытие модального окна для списания оборудования
    closeSubtractEquipmentModal() {
      this.showSubtractEquipmentModal = false;
    },

    // Сброс формы для услуг
    resetServiceForm() {
      this.currentService = {};
      this.modalTitle = '';
      this.showServiceModal = false;
    },

    // Сброс формы для оборудования
    resetEquipmentForm() {
      this.currentEquipment = {};
      this.equipmentModalTitle = '';
      this.showEquipmentModal = false;
    },

    // Сброс данных списания оборудования
    resetSubtractEquipmentForm() {
      this.currentEquipment = {};
      this.showSubtractEquipmentModal = false;
    },

    // Разворачивание/Скрытие оборудования для услуги
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
      return requiredEquipment.map(eq => {
        if (eq.type && eq.brand && eq.model) {
          return `${eq.type} - ${eq.brand} ${eq.model}`;
        } else {
          return 'Некорректные данные оборудования';
        }
      }).join(', ');
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

    // Метод для отображения уведомлений с помощью SweetAlert2
    showAlert(message, type) {
      Swal.fire({
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
  },
  created() {
    this.fetchServicesAndEquipment();
  },
};
</script>

<style scoped>
.services-dashboard {
  padding: 20px;
  max-width: 1000px;
  margin: 0 auto;
}

h2 {
  text-align: center;
  margin-bottom: 20px;
}

.table-section {
  margin-bottom: 30px;
}

.add-service-button,
.add-equipment-button {
  text-align: right;
  margin-bottom: 10px;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 12px;
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
  margin: 2px;
  transition: background-color 0.3s;
}

.btn:hover {
  background-color: #45a049;
}

.btn.danger {
  background-color: #f44336;
}

.btn.danger:hover {
  background-color: #e53935;
}

.btn.add {
  background-color: #2196F3;
}

.btn.add:hover {
  background-color: #1976D2;
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
  color: white;
  border-radius: 10px;
  padding: 15px;
  width: 180px;
  text-align: center;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  cursor: pointer;
  transition: transform 0.3s, box-shadow 0.3s;
}

.card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
}

.input-group-text {
  width: 50px;
  text-align: center;
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

.form-select {
  width: 100%;
  padding: 8px 12px;
  margin-bottom: 15px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.input-group {
  display: flex;
  align-items: center;
}

.input-group .form-control {
  flex: 1;
}

.d-flex {
  display: flex;
}

.justify-content-end {
  justify-content: flex-end;
}
</style>
