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

    <div class="table-container">
      <table class="v-table">
        <thead>
          <tr>
            <th v-for="header in headers" :key="header.value" class="table-header">
              {{ header.text }}
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in computedTransactionData" :key="item.id" class="table-row">
            <td>{{ formatDateTime(item.date_transferred) }}</td>
            <td>{{ item.sender_name }}</td>
            <td>
              {{ item.sent_amount }}
            </td>
            <td>{{ item.beneficiary_name }}</td>
            <td>
              {{ item.received_amount }}
            </td>
            <td>{{ item.status }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <h2 v-if="transactionData.length === 0" class="text-center">Transaction History is Empty!</h2>
  </div>
</template>

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

<script>
import axios from 'axios'

export default {
  name: 'TransactionList',
  data() {
    return {
      transactionData: [],
      headers: [
        { text: 'Date Transferred', value: 'date_transferred' },
        { text: 'Sender', value: 'sender_name' },
        { text: 'Amount Transferred', value: 'sent_amount' },
        { text: 'Beneficiary', value: 'beneficiary_name' },
        { text: 'Amount Received', value: 'received_amount' },
        { text: 'Transaction Status', value: 'status' }
      ]
    }
  },
  computed: {
    computedTransactionData() {
      return this.transactionData.map((item) => {
        return {
          ...item,
          sent_amount: `${item.amount_transferred} ${item.amount_transferred_currency}`,
          received_amount: `${item.amount_received} ${item.amount_received_currency}`
        }
      })
    }
  },
  // Pull transactions when the component is created
  created() {
    this.fetchData()
  },
  methods: {
    formatDateTime(dateTime) {
      const options = {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      }
      return new Date(dateTime).toLocaleString(undefined, options)
    },
    fetchData() {
      const jwt_token = localStorage.getItem('leon_token')
      const config = {
        headers: { Authorization: `Bearer ${jwt_token}` }
      }
      3

      axios
        .get('http://localhost:4000/transactions', config)
        .then((response) => {
          if (response.data.transactions) {
            this.transactionData = response.data.transactions
          }
        })
        .catch((error) => {
          if (error.response.status == 401) {
            this.$router.push('/login')
            localStorage.removeItem('leon_token')
          }
        })
    }
  }
}
</script>
