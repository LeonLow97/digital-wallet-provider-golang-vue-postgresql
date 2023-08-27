<template>
  <div>
    <v-card class="text-center">
      <v-card-text>
        <v-row align="center">
          <v-col cols="2">
            <v-btn color="primary" @click="logout">Logout</v-btn>
          </v-col>
          <v-col cols="8">
            <h1 class="display-1">Home</h1>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <user-details />
    <transfer-funds />
  </div>
</template>

<script>
import UserDetails from '@/components/HomePage/UserDetails.vue'
import TransferFunds from '@/components/HomePage/TransferFunds.vue'

export default {
  name: 'HomeView',
  components: {
    UserDetails,
    TransferFunds
  },
  methods: {
    logout() {
      this.$router.push('/login')
      localStorage.removeItem('leon_access_token')
    }
  }
}
</script>
