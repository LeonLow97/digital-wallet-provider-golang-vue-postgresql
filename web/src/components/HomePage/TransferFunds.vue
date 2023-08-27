<template>
  <v-dialog v-model="modalOpen" persistent max-width="800px">
    <v-card>
      <v-card-title>Transfer Money</v-card-title>
      <v-card-text>
        <v-select
          v-model="selectedBeneficiary"
          :items="beneficiaries"
          label="Select Recipient"
          item-title="beneficiary_name"
          item-value="beneficiary_id"
          variant="outlined"
          return-object
        ></v-select>

        <v-row justify="center">
          <v-col cols="12" sm="6">
            <v-form v-if="selectedBeneficiary">
              <v-text-field v-model="amountTransferred" label="Amount Transferred"></v-text-field>
              <v-text-field
                v-model="amountReceived"
                label="Amount Received"
                :readonly="true"
                :class="{ 'readonly-text-field': true }"
              ></v-text-field>
            </v-form>
          </v-col>

          <v-col cols="12" sm="6">
            <v-form v-if="selectedBeneficiary">
              <v-text-field
                v-model="amountTransferredCurrency"
                label="Amount Transferred Currency"
                :readonly="true"
                :class="{ 'readonly-text-field': true }"
              ></v-text-field>
              <v-text-field
                v-model="amountReceivedCurrency"
                label="Amount Received Currency"
                :readonly="true"
                :class="{ 'readonly-text-field': true }"
              ></v-text-field>
            </v-form>
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions>
        <v-btn color="primary" @click="transferFunds" :disabled="!selectedBeneficiary"
          >Transfer</v-btn
        >
        <v-btn @click="closeModal">Cancel</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import getBeneficiaries from '@/api/beneficiaries/getBeneficiaries'

export default {
  name: 'TransferFunds',
  data() {
    return {
      modalOpen: true,
      beneficiaries: [],
      selectedBeneficiary: null,
      amountTransferred: null,
      amountTransferredCurrency: null,
      amountReceived: 0,
      amountReceivedCurrency: null
    }
  },
  watch: {
    selectedBeneficiary(newValue) {
      console.log(newValue)
      // if (newValue) {
      //   this.amountReceivedCurrency = newValue.split(',  ')[2]
      //   this.amountTransferred = 0
      // }
    }
  },
  mounted() {
    this.fetchBeneficiaries()
  },
  methods: {
    async fetchBeneficiaries() {
      try {
        const data = await getBeneficiaries()
        if (data.beneficiaries != null) {
          this.beneficiaries = data.beneficiaries
        }
      } catch (error) {
        if (error.response.status === 401) {
          this.$router.push('/login')
          localStorage.removeItem('leon_access_token')
        }
      }
    },
    transferFunds() {
      console.log('Transferring Funds...')
    },
    closeModal() {
      console.log('Closing Modal...')
    }
  }
}
</script>
