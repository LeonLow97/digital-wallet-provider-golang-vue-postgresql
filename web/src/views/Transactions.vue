<template>
  <div>
    <v-card class="text-center">
      <v-card-actions>
        <v-row align="center">
          <v-col cols="2">
            <router-link class="router-link" to="/">
              <h2>Home</h2>
            </router-link>
          </v-col>
          <v-col cols="8">
            <h1 class="display-1">Transaction History</h1>
          </v-col>
        </v-row>
      </v-card-actions>
    </v-card>

    <v-table class="table-container" fixed-header height="300px" width="80%">
      <thead>
        <tr>
          <th v-for="header in headers" :key="header">
            {{ header }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="transaction in computedTransactions" :key="transaction">
          <td>{{ formatDateTime(transaction.date_transferred) }}</td>
          <td>{{ transaction.sender_name }}</td>
          <td>{{ transaction.amount_transferred }}</td>
          <td>{{ transaction.beneficiary_name }}</td>
          <td>{{ transaction.amount_received }}</td>
          <td>{{ transaction.status }}</td>
        </tr>
      </tbody>
    </v-table>
    <h2 v-if="transactions.length === 0" class="text-center">Transaction History is Empty!</h2>
  </div>
</template>

<script>
import getTransactions from '@/api/transactions/getTransactions'
import formatDateTime from '@/utils/formatDateTime'

export default {
  name: 'TransactionList',
  data() {
    return {
      transactions: [],
      headers: [
        'Date Transferred',
        'Sender',
        'Amount Transferred',
        'Recipient',
        'Amount Received',
        'Transaction Status'
      ]
    }
  },
  computed: {
    computedTransactions() {
      return this.transactions.map((item) => {
        return {
          ...item,
          amount_transferred: `${item.amount_transferred} ${item.amount_transferred_currency}`,
          amount_received: `${item.amount_received} ${item.amount_received_currency}`
        }
      })
    }
  },
  // Pull transactions when the component is created
  mounted() {
    this.fetchTransactions()
  },
  methods: {
    formatDateTime,
    async fetchTransactions() {
      try {
        const data = await getTransactions()
        if (data.transactions != null) {
          this.transactions = data.transactions
        }
      } catch (error) {
        if (error.response.status === 401) {
          this.$router.push('/login')
          localStorage.removeItem('leon_access_token')
        }
      }
    }
  }
}
</script>

<style>
.table-container {
  margin: 50px auto;
  max-width: 80%;
  overflow-x: auto;
}

.v-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

.table-header {
  background-color: #f5f5f5;
  color: #333;
  font-weight: bold;
  padding: 12px 16px;
}

.table-row:nth-child(even) {
  background-color: #f9f9f9;
}

.table-row:hover {
  background-color: #f0f0f0;
}

.v-table td,
.v-table th {
  text-align: left;
  padding: 12px 16px;
}

.router-link {
  display: inline-block;
  padding: 10px 20px;
  background-color: mediumseagreen;
  color: #ffffff; /* White text color */
  border-radius: 5px;
  text-decoration: none;
  transition: background-color 0.3s;
}

.router-link:hover {
  color: beige;
}
</style>
