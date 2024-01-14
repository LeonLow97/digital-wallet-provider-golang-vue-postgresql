import { shallowMount } from '@vue/test-utils'
import Transfer from '@/views/Transfer.vue'
import { describe, expect, test, vi } from 'vitest'
import axios from 'axios'
import { createVuetify } from 'vuetify'

vi.mock('axios')

const fetchBeneficiariesResponse = {
  data: {
    beneficiaries: [
      {
        beneficiary_id: 2,
        beneficiary_name: 'Bob',
        mobile_number: '+65 89230122',
        currency: 'SGD'
      },
      {
        beneficiary_id: 4,
        beneficiary_name: 'David',
        mobile_number: '+49 1234567890',
        currency: 'EUR'
      }
    ]
  }
}

const fetchUserResponse = {
  data: {
    user: {
      username: 'Charlie',
      mobile_number: '+1 555-123-4567',
      currency: 'USD',
      balance: '6091.02'
    }
  }
}

describe('Transfer.vue', () => {
  describe('fetchUser', () => {
    const vuetify = createVuetify()

    it('fetches logged in user details', async () => {
      localStorage.setItem('leon_access_token', 'JWT_Token')

      axios.get
        .mockResolvedValueOnce(fetchUserResponse)
        .mockResolvedValueOnce(fetchBeneficiariesResponse)

      const wrapper = shallowMount(Transfer, {
        global: {
          plugins: [vuetify]
        }
      })

      await new Promise((resolve) => setTimeout(resolve, 100))

      expect(axios.get).toHaveBeenCalledTimes(2)
      expect(axios.get).toHaveBeenCalledWith('http://localhost:4000/user', expect.any(Object))
      expect(axios.get).toHaveBeenCalledWith(
        'http://localhost:4000/beneficiaries',
        expect.any(Object)
      )

      expect(wrapper.vm.beneficiaries).toEqual([
        'Bob,  +65 89230122,  SGD',
        'David,  +49 1234567890,  EUR'
      ])
      expect(wrapper.vm.user).toEqual({
        user: {
          username: 'Charlie',
          mobile_number: '+1 555-123-4567',
          currency: 'USD',
          balance: '6091.02'
        }
      })
    })
  })
})
