import { render, screen } from '@testing-library/vue'
import { createVuetify } from 'vuetify'
import { RouterLinkStub } from '@vue/test-utils'
import axios from 'axios'

import Transactions from '@/views/Transactions.vue'

vi.mock('axios')

describe('Transactions', () => {
  beforeEach(() => {
    localStorage.setItem('leon_access_token', 'Bearer JWTToken')
  })

  const renderTransactions = () => {
    const vuetify = createVuetify()

    render(Transactions, {
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
  }

  const mockTransactionsResponse = (transaction = {}) => {
    axios.get.mockResolvedValue({
      data: transaction
    })
  }

  const mockTransactionsErrorResponse = (error = {}) => {
    axios.get.mockRejectedValue(error)
  }

  describe('when the component first loads', () => {
    it('displays no transactions', () => {
      mockTransactionsResponse()
      renderTransactions()

      const heading = screen.getByRole('heading', {
        name: /transaction history is empty!/i
      })

      expect(heading.textContent).toBe('Transaction History is Empty!')
    })

    it('displays 2 transactions', async () => {
      const mockTransactions = [
        {
          sender_name: 'Alice',
          amount_transferred: 100,
          amount_transferred_currency: 'SGD',
          beneficiary_name: 'Daniel Ong',
          amount_received: 100,
          amount_received_currency: 'SGD',
          status: 'CONFIRMED',
          date_transferred: '2023-07-02T12:00:00Z'
        },
        {
          sender_name: 'Alice',
          amount_transferred: 100,
          amount_transferred_currency: 'SGD',
          beneficiary_name: 'Leon Low',
          amount_received: 100,
          amount_received_currency: 'SGD',
          status: 'CONFIRMED',
          date_transferred: '2023-07-02T12:00:00Z'
        }
      ]

      mockTransactionsResponse({
        transactions: mockTransactions
      })

      renderTransactions()

      const tableCellBeneficiary1 = await screen.findByRole('cell', {
        name: /leon low/i
      })
      const tableCellBeneficiary2 = await screen.findByRole('cell', {
        name: /daniel ong/i
      })

      const tableCellAmountTransferred = await screen.findAllByText('100 SGD')

      expect(tableCellBeneficiary1.textContent).toBe('Leon Low')
      expect(tableCellBeneficiary2.textContent).toBe('Daniel Ong')
      expect(tableCellAmountTransferred[0]).toBeInTheDocument()
    })

    it('returns an error when getting transactions', async () => {
      const push = vi.fn()
      const $router = { push }

      mockTransactionsErrorResponse({
        response: {
          status: 401,
          data: 'Unauthorized'
        }
      })

      const vuetify = createVuetify()

      render(Transactions, {
        global: {
          mocks: {
            $router
          },
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

      // expect(push).toHaveBeenCalledWith('/login') // Check that $router.push was called with the correct argument
      // expect(localStorage.removeItem).toHaveBeenCalledWith('leon_access_token')
    })
  })
})
