<!-- src/components/Auth.vue -->
<template>
  <div class="d-flex justify-content-center align-items-center min-vh-100">
    <b-card class="p-4 shadow auth-card">
      <h2 class="text-center mb-4">{{ isRegistering ? 'Регистрация' : 'Вход в аккаунт' }}</h2>

      <b-form @submit.prevent="isRegistering ? register() : login()" novalidate>
        <!-- Форма регистрации -->
        <div v-if="isRegistering">
          <b-form-group label="Полное имя" label-for="fullName">
            <b-form-input
                id="fullName"
                v-model="fullName"
                placeholder="Введите полное имя"
                required
                :state="isValidFullName(fullName)"
                aria-describedby="fullName-feedback"
                class="custom-input"
            ></b-form-input>
            <b-form-invalid-feedback id="fullName-feedback">
              Пожалуйста, введите корректное ФИО (например, Иванов Иван Иванович).
            </b-form-invalid-feedback>
          </b-form-group>

          <b-form-group label="Номер телефона" label-for="phone">
            <b-form-input
                id="phone"
                v-model="phone"
                placeholder="7(999) 999-99-99"
                required
                v-mask="'7(###) ###-####'"
                :state="isValidPhone(phone)"
                aria-describedby="phone-feedback"
                class="custom-input"
            ></b-form-input>
            <b-form-invalid-feedback id="phone-feedback">
              Номер телефона должен начинаться с 7 и содержать 11 цифр.
            </b-form-invalid-feedback>
          </b-form-group>

          <b-form-group label="Электронная почта" label-for="email">
            <b-form-input
                type="email"
                id="email"
                v-model="filters.email"
                placeholder="example@mail.com"
                required
                :state="isValidEmail(filters.email)"
                aria-describedby="email-feedback"
                class="custom-input"
            ></b-form-input>
            <b-form-invalid-feedback id="email-feedback">
              Пожалуйста, введите корректный email.
            </b-form-invalid-feedback>
          </b-form-group>
        </div>

        <!-- Форма входа -->
        <div v-else>
          <b-form-group label="Номер телефона" label-for="phone">
            <b-form-input
                id="phone"
                v-model="phone"
                placeholder="7(999) 999-99-99"
                required
                v-mask="'7(###) ###-####'"
                :state="isValidPhone(phone)"
                aria-describedby="phone-feedback"
                class="custom-input"
            ></b-form-input>
            <b-form-invalid-feedback id="phone-feedback">
              Номер телефона должен начинаться с 7 и содержать 11 цифр.
            </b-form-invalid-feedback>
          </b-form-group>

          <b-form-group label="Пароль" label-for="password">
            <b-form-input
                type="password"
                id="password"
                v-model="password"
                placeholder="Введите пароль"
                required
                class="custom-input"
            ></b-form-input>
            <!-- Без валидации пароля, поэтому нет feedback -->
          </b-form-group>
        </div>

        <b-button type="submit" variant="danger" class="w-100 custom-btn">
          {{ isRegistering ? 'Зарегистрироваться' : 'Войти' }}
        </b-button>

        <p v-if="errorMessage" class="error text-danger mt-2">{{ errorMessage }}</p>
      </b-form>

      <div class="toggle-section mt-3 text-center">
        <span>{{ isRegistering ? 'Уже есть аккаунт?' : 'Нет аккаунта?' }}</span>
        <b-button variant="link" class="btn-toggle p-0" @click="toggleForm">
          {{ isRegistering ? 'Войти' : 'Зарегистрироваться' }}
        </b-button>
      </div>
    </b-card>
  </div>
</template>

<script>
export default {
  name: 'Auth',
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

      // Поскольку валидация пароля отключена, можно добавить базовую проверку на наличие пароля
      if (!this.password) {
        this.errorMessage = 'Пожалуйста, введите пароль.';
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

        this.$swal.fire({
          title: 'Успех!',
          text: 'Авторизация прошла успешно!',
          icon: 'success',
          timer: 3000,
          showConfirmButton: false,
          toast: true,
          position: 'top-end'
        });

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
        this.$swal.fire({
          title: 'Ошибка!',
          text: this.errorMessage,
          icon: 'error',
          timer: 3000,
          showConfirmButton: false,
          toast: true,
          position: 'top-end'
        });
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
        this.$swal.fire({
          title: 'Успех!',
          text: 'Регистрация прошла успешно! Теперь вы можете войти в систему.',
          icon: 'success',
          timer: 3000,
          showConfirmButton: false,
          toast: true,
          position: 'top-end'
        });
        this.resetForm();
        this.isRegistering = false;  // Переключаемся на форму входа
      } catch (error) {
        this.errorMessage = error.message || 'Ошибка регистрации. Попробуйте снова.';
        this.$swal.fire({
          title: 'Ошибка!',
          text: this.errorMessage,
          icon: 'error',
          timer: 3000,
          showConfirmButton: false,
          toast: true,
          position: 'top-end'
        });
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
@import url('https://fonts.googleapis.com/css2?family=Roboto:wght@400;500;700&display=swap');

.d-flex {
  font-family: 'Roboto', sans-serif;
}

.auth-card {
  border: none;
  background: #ffffff;
  border-radius: 15px;
}

h2 {
  color: #343a40;
  font-weight: 500;
}

.custom-btn {
  background-color: #ff6b6b;
  color: white;
  border: none;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 1.1em;
  transition: background-color 0.3s;
}

.custom-btn:hover {
  background-color: #ff4d4d;
}

.custom-toggle-btn {
  color: #ff6b6b;
  text-decoration: underline;
}

.custom-toggle-btn:hover {
  color: #ff4d4d;
}

.custom-input {
  font-family: 'Roboto', sans-serif;
}

.invalid {
  border-color: #dc3545 !important;
}

.error {
  font-size: 0.9em;
}

.toggle-section span {
  color: #555;
}
</style>
