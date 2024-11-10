<script>
export default {
  name: 'home',
  data() {
    return {
      phone: '', // Используем телефон вместо email
      password: '',
      errorMessage: ''
    };
  },
  methods: {
    async login() {
      try {
        const response = await fetch('http://localhost:8080/api/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            phone: this.phone, // Используем телефон
            password: this.password
          }),
          credentials: 'include' // Включаем отправку куки
        });

        if (!response.ok) {
          throw new Error('Ошибка аутентификации');
        }

        const data = await response.json();

        // Получаем роль из данных ответа
        const role = data.role;

        // Показываем сообщение об успешном входе
        alert('Авторизация прошла успешно!');

        // Выполняем маршрутизацию на основе роли
        if (role === 'admin') {
          this.$router.push({ name: 'AdminHome' });
        } else if (role === 'worker') {
          this.$router.push({ name: 'WorkerDashboard' });
        } else if (role === 'client') {
          this.$router.push({ name: 'ClientHome' });
        } else {
          throw new Error('Неизвестная роль');
        }

      } catch (error) {
        this.errorMessage = 'Ошибка входа. Проверьте данные и попробуйте снова.';
        console.error(error);
      }
    }
  }
};
</script>

<template>
  <div class="login-container">
    <h2>Вход в аккаунт</h2>
    <form @submit.prevent="login">
      <input type="tel" v-model="phone" placeholder="Номер телефона" required />
      <input type="password" v-model="password" placeholder="Пароль" required />
      <button type="submit" class="btn">Войти</button>
      <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
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
  background: rgba(255, 255, 255, 0.85);
  padding: 40px 30px;
  border-radius: 15px;
  box-shadow: 0px 8px 15px rgba(0, 0, 0, 0.3);
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
