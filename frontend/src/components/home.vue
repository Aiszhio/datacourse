<script>
export default {
  name: 'home',
};
</script>

<template>
  <div class="login-container">
    <h2>Вход в аккаунт</h2>
    <form @submit.prevent="login">
      <input type="email" v-model="email" placeholder="Электронная почта" required />
      <input type="password" v-model="password" placeholder="Пароль" required />
      <button type="submit" class="btn">Войти</button>
    </form>
  </div>
</template>

<style scoped>
body {
  margin: 0;
  padding: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-size: cover;
  background-position: center;
  font-family: 'Arial', sans-serif;
}

.login-container {
  background: rgba(255, 255, 255, 0.85); /* Чуть более прозрачный белый фон для мягкого эффекта */
  padding: 40px 30px;
  border-radius: 15px;
  box-shadow: 0px 8px 15px rgba(0, 0, 0, 0.3); /* Мягкая тень для улучшенного визуального эффекта */
  width: 360px;
  text-align: center;
}

h2 {
  margin-bottom: 30px;
  color: #333;
  font-size: 1.8em;
  font-weight: bold;
}

input {
  width: 90%;
  padding: 15px;
  margin: 15px 0;
  border: 2px solid #ccc;
  border-radius: 8px;
  font-size: 1em;
  transition: border-color 0.3s;
}

input:focus {
  border-color: #ff6b6b;
}

.btn {
  background-color: #ff6b6b;
  color: white;
  border: none;
  padding: 15px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 1.1em;
  width: 100%;
  transition: background-color 0.3s;
}

.btn:hover {
  background-color: #ff4d4d;
}

.logo img {
  max-width: 150px;
}
</style>
