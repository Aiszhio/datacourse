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
    <div class="links">
      <a href="#">Забыли пароль?</a><br />
      <a href="#">Создать аккаунт</a>
    </div>
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
  background-image: url('@/assets/studio-bg.jpg'); /* Фоновое изображение */
  background-size: cover;
  font-family: 'Arial', sans-serif;
}

.login-container {
  background: rgba(255, 255, 255, 0.9);
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.2);
  width: 320px;
  height: 40%;
  text-align: center;
}

h2 {
  margin-bottom: 20px;
  color: #333;
}

input {
  width: 90%;
  padding: 10px;
  margin: 10px 0;
  border: 1px solid #ccc;
  border-radius: 5px;
}

.btn {
  background-color: #ff6b6b;
  color: white;
  border: none;
  padding: 10px;
  border-radius: 5px;
  cursor: pointer;
  width: 100%;
}

.btn:hover {
  background-color: #ff4d4d;
}

.logo {
  margin-bottom: 20px;
}

.logo img {
  max-width: 150px;
}

.links {
  margin-top: 20px;
}

.links a {
  color: #555;
  text-decoration: none;
}

.links a:hover {
  text-decoration: underline;
}
</style>