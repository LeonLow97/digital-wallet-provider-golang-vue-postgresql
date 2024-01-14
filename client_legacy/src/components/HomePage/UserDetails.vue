<template>
  <v-row justify="center" class="mt-5">
    <v-col cols="12" sm="8" md="6">
      <v-card class="text-center" elevation="3" outlined>
        <v-card-text>
          <h2 v-if="user.username">Welcome back {{ user.username }}!</h2>
          <br />
          <h3 v-if="user.balance">My Balance: {{ user.balance }} {{ user.currency }}</h3>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<script>
import getUser from '@/api/users/getUser'

export default {
  name: 'UserDetails',
  data() {
    return {
      user: {},
      amountTransferredCurrency: null
    }
  },
  mounted() {
    this.fetchUser()
  },
  methods: {
    async fetchUser() {
      try {
        const data = await getUser()
        this.user = data
        this.amountTransferredCurrency = data.currency
      } catch (error) {
        if (error.response.status == 401) {
          this.$router.push('/login')
          localStorage.removeItem('leon_access_token')
        }
      }
    }
  }
}
</script>
