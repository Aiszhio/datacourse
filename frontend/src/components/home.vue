<template>
  <div class="auth-container">
    <h2>{{ isRegistering ? 'Регистрация' : 'Вход в аккаунт' }}</h2>

    <form @submit.prevent="isRegistering ? register() : login()">
      <!-- Форма регистрации -->
      <div v-if="isRegistering">
        <input
            type="text"
            v-model="fullName"
            placeholder="Полное имя"
            required
            :class="{'invalid': !isValidFullName(fullName)}"
        />
        <input
            type="tel"
            v-model="phone"
            placeholder="Номер телефона"
            required
            v-mask="'7(###) ###-####'"
            :class="{'invalid': !isValidPhone(phone)}"
        />
        <input
            type="email"
            v-model="filters.email"
            placeholder="Электронная почта"
            required
            :class="{'invalid': !isValidEmail(filters.email)}"
        />
      </div>

      <!-- Форма входа -->
      <div v-else>
        <input
            type="tel"
            v-model="phone"
            placeholder="Номер телефона"
            required
            v-mask="'7(###) ###-####'"
            :class="{'invalid': !isValidPhone(phone)}"
        />
        <input
            type="password"
            v-model="password"
            placeholder="Пароль"
            required
        />
      </div>

      <button type="submit" class="btn">
        {{ isRegistering ? 'Зарегистрироваться' : 'Войти' }}
      </button>

      <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
    </form>

    <div class="toggle-section">
      <span>{{ isRegistering ? 'Уже есть аккаунт?' : 'Нет аккаунта?' }}</span>
      <button @click="toggleForm" class="btn-toggle">
        {{ isRegistering ? 'Войти' : 'Зарегистрироваться' }}
      </button>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';

export default {
  name: 'Home',
  data() {
    return {
      isRegistering: false,  // Флаг для переключения между входом и регистрацией
      phone: '',
      password: '',
      fullName: '',
      errorMessage: '',
      filters: {
        email: ''
      },
    };
  },
  methods: {
    // Валидация электронной почты
    isValidEmail(email) {
      const regex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$/;
      return regex.test(email);
    },

    // Валидация номера телефона
    isValidPhone(phone) {
      const regex = /^7\d{10}$/;
      return regex.test(this.cleanPhone(phone));
    },

    // Валидация ФИО
    isValidFullName(name) {
      const regex = /^[A-Za-zА-Яа-яЁё]+ [A-Za-zА-Яа-яЁё]+ [A-Za-zА-Яа-яЁё]+$/;
      return regex.test(name);
    },

    // Метод для очистки маски телефона
    cleanPhone(phone) {
      return phone.replace(/\D/g, '');
    },

    async login() {
      if (!this.isValidPhone(this.phone)) {
        this.errorMessage = 'Пожалуйста, введите корректный номер телефона.';
        return;
      }

      try {
        const response = await fetch('http://localhost:8080/api/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            phone: this.cleanPhone(this.phone),
            password: this.password
          }),
          credentials: 'include'
        });

        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.error || 'Ошибка аутентификации');
        }

        const data = await response.json();
        const role = data.role;

        alert('Авторизация прошла успешно!');

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
        this.errorMessage = error.message || 'Ошибка входа. Проверьте данные и попробуйте снова.';
        console.error(error);
      }
    },

    async register() {
      if (!this.fullName || !this.isValidEmail(this.filters.email) || !this.isValidPhone(this.phone)) {
        this.errorMessage = 'Пожалуйста, проверьте данные.';
        return;
      }

      try {
        const response = await fetch('http://localhost:8080/api/register', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            FullName: this.fullName,
            PhoneNumber: this.cleanPhone(this.phone),
            Email: this.filters.email,
          }),
          credentials: 'include'
        });

        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.error || 'Ошибка регистрации');
        }

        const data = await response.json();
        alert('Регистрация прошла успешно! Теперь вы можете войти в систему.');
        this.resetForm();
        this.isRegistering = false;  // Переключаемся на форму входа
      } catch (error) {
        this.errorMessage = error.message || 'Ошибка регистрации. Попробуйте снова.';
        console.error(error);
      }
    },

    toggleForm() {
      this.isRegistering = !this.isRegistering;
      this.resetForm();
    },

    resetForm() {
      this.phone = '';
      this.password = '';
      this.fullName = '';
      this.filters.email = '';
      this.errorMessage = '';
    }
  }
};
</script>

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

.auth-container {
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
  margin: 10px 0;
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
  margin-top: 15px;
}

.btn:hover {
  background-color: #ff4d4d;
}

.btn-toggle {
  background-color: transparent;
  color: #ff6b6b;
  border: none;
  cursor: pointer;
  font-size: 1em;

  text-decoration: underline;
}

.btn-toggle:hover {
  color: #ff4d4d;
}

.error {
  color: red;
  font-size: 0.9em;
  margin-top: 10px;
}

.toggle-section {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.toggle-section span {
  margin-right: 10px;
  color: #555;
}
</style>
