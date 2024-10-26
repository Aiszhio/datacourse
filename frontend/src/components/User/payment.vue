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
      return date.toLocaleDateString('ru-RU', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
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
    <div class="invoices-section section">
      <h3>Ваши неоплаченные счета</h3>
      <div v-if="loading" class="loading">Загрузка ваших счетов...</div>
      <table class="invoices-table" v-if="!loading">
        <thead>
        <tr>
          <th>ID Счета</th>
          <th>Услуга</th>
          <th>Сумма</th>
          <th>Дата создания</th>
          <th>Оплата</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="invoice in invoices.filter(invoice => invoice.Status === 'Неоплачен')" :key="invoice.InvoiceID">
          <td>{{ invoice.InvoiceID }}</td>
          <td>{{ invoice.ServiceName }}</td>
          <td>{{ invoice.Amount }} руб.</td>
          <td>{{ formatDate(invoice.InvoiceDate) }}</td>
          <td><button @click="payInvoice(invoice.InvoiceID)" class="btn">Оплатить</button></td>
        </tr>
        <tr v-if="invoices.filter(invoice => invoice.Status === 'Неоплачен').length === 0">
          <td colspan="5" class="no-invoices">У вас нет неоплаченных счетов.</td>
        </tr>
        </tbody>
      </table>
    </div>

    <!-- История оплаченных счетов -->
    <div class="invoices-section section">
      <h3>История оплаченных счетов</h3>
      <div v-if="loading" class="loading">Загрузка ваших счетов...</div>
      <table class="invoices-table" v-if="!loading">
        <thead>
        <tr>
          <th>ID Счета</th>
          <th>Услуга</th>
          <th>Сумма</th>
          <th>Дата создания</th>
          <th>Дата оплаты</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="invoice in invoices.filter(invoice => invoice.Status === 'Оплачен')" :key="invoice.InvoiceID">
          <td>{{ invoice.InvoiceID }}</td>
          <td>{{ invoice.ServiceName }}</td>
          <td>{{ invoice.Amount }} руб.</td>
          <td>{{ formatDate(invoice.InvoiceDate) }}</td>
          <td>{{ formatDate(invoice.PaymentDate) }}</td>
        </tr>
        <tr v-if="invoices.filter(invoice => invoice.Status === 'Оплачен').length === 0">
          <td colspan="5" class="no-invoices">У вас нет оплаченных счетов.</td>
        </tr>
        </tbody>
      </table>
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

.invoices-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
}

.invoices-table th, .invoices-table td {
  padding: 12px;
  border: 1px solid #ddd;
  text-align: left;
}

.invoices-table th {
  background-color: #f2f2f2;
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
  padding: 10px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.btn:hover {
  background-color: #45a049;
}

.navigation-buttons {
  display: flex;
  justify-content: space-between;
  margin-top: 20px;
}
</style>
