<!-- src/App.vue -->
<template>
  <div id="app">
    <!-- Панель с кнопкой "Вход" -->
    <header class="header">
      <div class="container d-flex justify-content-end p-3">
        <button class="btn btn-primary" @click="handleLogin">Вход</button>
      </div>
    </header>

    <!-- Основной Контент -->
    <main class="main-content container">
      <router-view></router-view>
    </main>

  </div>
</template>

<script>
import { showSuccess, showError } from '@/utils/alerts';
// Если вы используете модальное окно для входа, раскомментируйте импорт
// import LoginModal from './components/LoginModal.vue';

export default {
  name: 'App',
  // components: {
  //   LoginModal,
  // },
  methods: {
    handleLogin() {
      // Выполняем редирект на маршрут 'home'
      this.$router.push({ name: 'home' });
    },
    // Если используете модальное окно для входа, добавьте метод для обработки успешного входа
    // handleLoginSuccess() {
    //   showSuccess('Вы вошли в систему.');
    //   this.$router.push({ name: 'home' });
    // },
    testAlert() {
      showSuccess('Это тестовое уведомление успеха!');
    },
  },
};
</script>

<style>
#app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

/* Стили для заголовка */
.header {
  background-color: #f8f9fa;
  border-bottom: 1px solid #dee2e6;
}

/* Стили для основного контента */
.main-content {
  flex: 1;
  padding: 20px 0;
}
</style>
