<!-- src/App.vue -->
<template>
  <div id="app">
    <!-- Основной Контент -->
    <main class="main-content container">
      <router-view></router-view>
    </main>

  </div>
</template>

<script>
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
