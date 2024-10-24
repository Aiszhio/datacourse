<script>
export default {
  name: 'Payment',
  data() {
    return {
      invoices: [], // Все счета (неоплаченные и оплаченные)
      loading: true // Флаг загрузки счетов
    };
  },
  mounted() {
    this.fetchInvoices(); // Загружаем счета при загрузке страницы
  },
  methods: {
    async fetchInvoices() {
      try {
        // Пример запроса для получения счетов пользователя
        const response = await fetch('http://localhost:8080/api/client/invoices', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer your-auth-token' // Токен авторизации, если требуется
          }
        });

        if (!response.ok) {
          throw new Error('Ошибка при получении счетов');
        }

        const data = await response.json();
        this.invoices = data; // Записываем полученные счета
      } catch (error) {
        console.error('Ошибка:', error.message);
        alert('Не удалось загрузить счета.');
      } finally {
        this.loading = false; // Отключаем флаг загрузки
      }
    },
    async payInvoice(invoiceID) {
      try {
        const response = await fetch(`http://localhost:8080/api/client/invoices/${invoiceID}/pay`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer your-auth-token'
          }
        });

        if (!response.ok) {
          throw new Error('Ошибка при оплате счета');
        }

        alert(`Счет #${invoiceID} успешно оплачен.`);
        this.fetchInvoices(); // Обновляем список счетов после успешной оплаты
      } catch (error) {
        console.error('Ошибка:', error.message);
        alert('Не удалось оплатить счет.');
      }
    },
    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleDateString('ru-RU', { year: 'numeric', month: 'long', day: 'numeric' });
    },
    goToOrders() {
      this.$router.push({ name: 'ClientOrders' }); // Переход на страницу заказов
    },
    goToHome() {
      this.$router.push({ name: 'ClientHome' }); // Переход на главную страницу
    }
  }
};
</script>

<template>
  <div class="payment-dashboard">
    <h2>Оплата счетов</h2>

    <!-- Список неоплаченных счетов -->
    <div class="invoices-section">
      <h3>Ваши неоплаченные счета</h3>
      <div v-if="loading" class="loading">Загрузка ваших счетов...</div>
      <div v-else>
        <div v-if="invoices.some(invoice => invoice.Status === 'Неоплачен')" class="invoices-list">
          <div class="invoice" v-for="invoice in invoices.filter(invoice => invoice.Status === 'Неоплачен')" :key="invoice.InvoiceID">
            <h4>Счет #{{ invoice.InvoiceID }}</h4>
            <p><strong>Услуга:</strong> {{ invoice.ServiceName }}</p>
            <p><strong>Сумма:</strong> {{ invoice.Amount }} руб.</p>
            <p><strong>Дата создания:</strong> {{ formatDate(invoice.InvoiceDate) }}</p>
            <button @click="payInvoice(invoice.InvoiceID)" class="btn">Оплатить</button>
          </div>
        </div>
        <div v-else class="no-invoices">У вас нет неоплаченных счетов.</div>
      </div>
    </div>

    <!-- Список оплаченных счетов (история) -->
    <div class="invoices-section">
      <h3>История оплаченных счетов</h3>
      <div v-if="loading" class="loading">Загрузка ваших счетов...</div>
      <div v-else>
        <div v-if="invoices.some(invoice => invoice.Status === 'Оплачен')" class="invoices-list">
          <div class="invoice" v-for="invoice in invoices.filter(invoice => invoice.Status === 'Оплачен')" :key="invoice.InvoiceID">
            <h4>Счет #{{ invoice.InvoiceID }}</h4>
            <p><strong>Услуга:</strong> {{ invoice.ServiceName }}</p>
            <p><strong>Сумма:</strong> {{ invoice.Amount }} руб.</p>
            <p><strong>Дата создания:</strong> {{ formatDate(invoice.InvoiceDate) }}</p>
            <p><strong>Дата оплаты:</strong> {{ formatDate(invoice.PaymentDate) }}</p> <!-- Дата оплаты -->
          </div>
        </div>
        <div v-else class="no-invoices">У вас нет оплаченных счетов.</div>
      </div>
    </div>

    <!-- Кнопки для навигации -->
    <div class="navigation-buttons">
      <button @click="goToOrders" class="btn">Посмотреть заказы</button>
      <button @click="goToHome" class="btn">На главную</button>
    </div>
  </div>
</template>

<style scoped>
.payment-dashboard {
  padding: 20px;
  margin: 0 auto;
  max-width: 900px;
}

.invoices-section {
  background-color: #f9f9f9;
  padding: 20px;
  margin-bottom: 20px;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.invoices-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.invoice {
  background-color: #ffffff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0px 2px 4px rgba(0, 0, 0, 0.1);
}

.invoice h4 {
  margin-bottom: 10px;
  color: #4CAF50;
}

.loading {
  text-align: center;
  color: #555;
  font-size: 1.2em;
}

.no-invoices {
  text-align: center;
  color: #999;
  font-size: 1.2em;
}

.btn {
  background-color: #4CAF50;
  color: white;
  padding: 12px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 1em;
  transition: background-color 0.3s ease;
  margin: 95px;
  width: 30vh;
}

.btn:hover {
  background-color: #45a049;
}

.navigation-buttons {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
