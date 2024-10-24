<script>
export default {
  name: 'admin',
  data() {
    return {
      newItem: '', // Новая запись
      filterText: '', // Текст для фильтрации данных
      data: [ // Примерные данные
        { id: 1, name: 'Фотосессия' },
        { id: 2, name: 'Фотоальбом' }
      ],
    };
  },
  computed: {
    filteredData() {
      return this.data.filter(item => item.name.toLowerCase().includes(this.filterText.toLowerCase()));
    }
  },
  methods: {
    addData() {
      if (this.newItem) {
        const newItemObj = {
          id: this.data.length + 1, // Генерация нового ID
          name: this.newItem
        };
        this.data.push(newItemObj);
        this.newItem = ''; // Очищаем поле ввода
      }
    },
    deleteItem(id) {
      this.data = this.data.filter(item => item.id !== id);
    },
    editItem(item) {
      // Логика для редактирования элемента
      alert(`Редактирование элемента: ${item.name}`);
    }
  }
};
</script>

<template>
  <div class="admin-dashboard">
    <h2>Панель администратора</h2>

    <!-- Форма для добавления данных -->
    <div class="add-section">
      <h3>Добавить новый элемент</h3>
      <form @submit.prevent="addData">
        <input type="text" v-model="newItem" placeholder="Введите название" required />
        <button type="submit" class="btn">Добавить</button>
      </form>
    </div>

    <!-- Фильтр для поиска -->
    <div class="filter-section">
      <h3>Фильтр данных</h3>
      <input type="text" v-model="filterText" placeholder="Введите текст для поиска" />
    </div>

    <!-- Таблица данных -->
    <div class="table-section">
      <h3>Данные</h3>
      <table>
        <thead>
        <tr>
          <th>ID</th>
          <th>Название</th>
          <th>Действия</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="item in filteredData" :key="item.id">
          <td>{{ item.id }}</td>
          <td>{{ item.name }}</td>
          <td>
            <button @click="editItem(item)" class="btn">Редактировать</button>
            <button @click="deleteItem(item.id)" class="btn delete-btn">Удалить</button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.admin-dashboard {
  padding: 20px;
  margin: 0 auto;
  max-width: 900px;
}

.add-section,
.filter-section,
.table-section {
  margin-bottom: 20px;
  background-color: #f9f9f9;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.1);
}

h2 {
  text-align: center;
  margin-bottom: 20px;
}

input {
  padding: 10px;
  margin: 10px 0;
  border: 1px solid #ccc;
  border-radius: 5px;
  width: 100%;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
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
  padding: 10px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.delete-btn {
  background-color: #e74c3c;
}

.btn:hover {
  background-color: #45a049;
}

.delete-btn:hover {
  background-color: #c0392b;
}
</style>
