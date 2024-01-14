<template>
  <v-row justify="center" class="button-container">
    <v-col cols="auto">
      <v-btn @click="openModal" color="success" size="large" variant="elevated"
        >Transfer Money</v-btn
      >

      <v-btn
        color="info"
        size="large"
        type="submit"
        variant="elevated"
        @click="navigateToTransactions"
      >
        Transaction History
      </v-btn>
    </v-col>
  </v-row>
  <v-snackbar v-model="snackbar" :timeout="timeout" :color="color">
    {{ snackbarText }}
  </v-snackbar>
</template>

<script>
import axios from 'axios'
import currencyConversion from '@/utils/currencyConversion'
import getUser from '@/api/users/getUser'

export default {
  name: 'TransferFunds',
  data() {
    return {
      modalOpen: false,
      user: {},
      beneficiaries: [],
      selectedBeneficiary: null,
      amountTransferred: null,
      amountTransferredCurrency: null,
      amountReceived: 0,
      amountReceivedCurrency: null,
      snackbar: false,
      timeout: 3000,
      color: '',
      snackbarText: ''
    }
  },
  watch: {
    selectedBeneficiary(newValue) {
      if (newValue) {
        this.amountReceivedCurrency = newValue.split(',  ')[2]
        this.amountTransferred = 0
      }
    },
    // FOR DEVELOPMENT ONLY!!! IN FUTURE USE AN EXTERNAL API TO RETRIEVE FX DATA
    amountTransferred(newValue) {
      if (!newValue) {
        this.amountReceived = 0
      } else {
        const amountTransferredCurrency = this.amountTransferredCurrency
        const amountReceivedCurrency = this.amountReceivedCurrency

        const convertedAmount = currencyConversion(
          newValue,
          amountTransferredCurrency,
          amountReceivedCurrency
        )
        this.amountReceived = convertedAmount.toFixed(2)
      }
    }
  },
  created() {
    this.fetchUser()
    this.fetchBeneficiaries()
  },
  methods: {
    logout() {
      this.$router.push('/login')
      localStorage.removeItem('leon_access_token')
    },
    navigateToTransactions() {
      this.$router.push('/transactions')
    },
    fetchUser() {
      getUser()
        .then((data) => {
          this.user = data
          this.amountTransferredCurrency = data.currency
        })
        .catch((error) => {
          if (error.response.status == 401) {
            this.$router.push('/login')
            localStorage.removeItem('leon_access_token')
          }
        })
    },
    fetchBeneficiaries() {
      const jwt_token = localStorage.getItem('leon_access_token')
      const config = {
        headers: { Authorization: `Bearer ${jwt_token}` }
      }

      axios
        .get('http://localhost:4000/beneficiaries', config)
        .then((response) => {
          this.beneficiaries = response.data.beneficiaries.map(
            (beneficiary) =>
              beneficiary.beneficiary_name +
              ',  ' +
              beneficiary.mobile_number +
              ',  ' +
              beneficiary.currency
          )
        })
        .catch((error) => {
          console.error(error)
        })
    },
    openModal() {
      this.modalOpen = true
    },
    closeModal() {
      this.selectedBeneficiary = null
      this.modalOpen = false
      this.amountTransferred = null
      this.amountReceived = 0
    },
    transferFunds() {
      const jwt_token = localStorage.getItem('leon_access_token')
      const config = {
        headers: { Authorization: `Bearer ${jwt_token}` }
      }

      const payload = {
        beneficiary_name: this.selectedBeneficiary.split(',  ')[0],
        mobile_number: this.selectedBeneficiary.split(',  ')[1],
        amount_transferred: this.amountTransferred,
        amount_transferred_currency: this.amountTransferredCurrency.toUpperCase()
      }

      axios
        .post('http://localhost:4000/transaction', payload, config)
        .then((response) => {
          this.closeModal()
          if (response.status == 201) {
            this.snackbar = true
            this.color = 'green'
            this.snackbarText = 'Successfully created a transaction!'
            this.fetchUser()
          }
        })
        .catch((error) => {
          this.snackbar = true
          this.color = 'red'
          this.snackbarText = error.response.data
        })
    }
  }
}
</script>

<style>
.button-container {
  margin-top: 30px;
  gap: 10px;
}

.readonly-text-field input {
  background-color: #dad7d7;
  color: #777;
  cursor: not-allowed;
}
</style>
