import { shallowMount } from '@vue/test-utils'
import Transactions from '@/views/Transactions.vue'
import { describe, it, expect, test, vi } from 'vitest'
import axios from 'axios'
import { createVuetify } from 'vuetify'
import { RouterLinkStub } from '@vue/test-utils'

vi.mock('axios')

describe('Transactions.vue', () => {
  const vuetify = createVuetify()

  describe('fetchUsers', () => {
    it('GET Request, fetches the transaction history', async () => {
      const transactionsMock = [
        {
          date_transferred: '2023-07-02T12:00:00Z',
          sender_name: 'Alice',
          amount_transferred: 100,
          amount_transferred_currency: 'SGD',
          beneficiary_name: 'Bob',
          amount_received: 100,
          amount_received_currency: 'SGD',
          status: 'CONFIRMED'
        },
        {
          date_transferred: '2023-07-03T12:00:00Z',
          sender_name: 'Alice',
          amount_transferred: 100,
          amount_transferred_currency: 'SGD',
          beneficiary_name: 'Charlie',
          amount_received: 74,
          amount_received_currency: 'USD',
          status: 'CONFIRMED'
        }
      ]

      localStorage.setItem('leon_access_token', 'JWT_Token')

      axios.get.mockImplementation((url, config) => {
        expect(url).toBe('http://localhost:4000/transactions')
        expect(config.headers.Authorization).toBe('Bearer JWT_Token')

        return Promise.resolve({
          data: {
            transactions: transactionsMock
          }
        })
      })

      const wrapper = shallowMount(Transactions, {
        global: {
          plugins: [vuetify],
          stubs: {
            'router-link': RouterLinkStub
          }
        },
        mocks: {
          localStorage: {
            getItem: () => 'leon_access_token'
          }
        }
      })

      await new Promise((resolve) => setTimeout(resolve, 100))

      expect(axios.get).toHaveBeenCalled()
      expect(wrapper.vm.transactionData).toEqual(transactionsMock)
    })
  })

  it('formats date and time correctly', () => {
    const wrapper = shallowMount(Transactions, {
      global: {
        plugins: [vuetify],
        stubs: {
          'router-link': RouterLinkStub
        }
      }
    })

    const formattedDateTime = wrapper.vm.formatDateTime('2023-07-02T12:34:56')
    expect(formattedDateTime).toEqual('July 2, 2023 at 12:34:56 PM')
  })

  it.only('compute transaction data correctly', () => {
    const transactionData = [
      {
        date_transferred: '2023-07-02T12:00:00Z',
        sender_name: 'Alice',
        amount_transferred: 100,
        amount_transferred_currency: 'SGD',
        beneficiary_name: 'Bob',
        amount_received: 100,
        amount_received_currency: 'SGD',
        status: 'CONFIRMED'
      },
      {
        date_transferred: '2023-07-03T12:00:00Z',
        sender_name: 'Alice',
        amount_transferred: 100,
        amount_transferred_currency: 'SGD',
        beneficiary_name: 'Charlie',
        amount_received: 74,
        amount_received_currency: 'USD',
        status: 'CONFIRMED'
      }
    ]

    const wrapper = shallowMount(Transactions, {
      global: {
        plugins: [vuetify],
        stubs: {
          'router-link': RouterLinkStub
        }
      },
      data() {
        return {
          transactionData
        }
      }
    })

    const computedTransactionData = wrapper.vm.computedTransactionData
    expect(computedTransactionData).toHaveLength(transactionData.length)

    // Add more expectations for the computed data as needed
    expect(computedTransactionData[0].sent_amount).toBe('100 SGD')
    expect(computedTransactionData[0].received_amount).toBe('100 SGD')

    expect(computedTransactionData[1].sent_amount).toBe('100 SGD')
    expect(computedTransactionData[1].received_amount).toBe('74 USD')
  })
})
