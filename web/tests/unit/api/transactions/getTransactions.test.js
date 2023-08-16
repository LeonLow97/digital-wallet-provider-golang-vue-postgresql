import axios from 'axios'

import getTransactions from '@/api/transactions/getTransactions'

vi.mock('axios')

describe('getTransactions', () => {
  const config = {
    headers: { Authorization: `Bearer JWT_Token_Test_Value` }
  }

  beforeEach(() => {
    axios.get.mockResolvedValue({
      data: [
        {
          sender_name: 'Alice',
          amount_transferred: 100,
          amount_transferred_currency: 'SGD',
          beneficiary_name: 'Bob',
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
    })

    localStorage.setItem('leon_access_token', 'JWT_Token_Test_Value')
  })

  afterEach(() => {
    localStorage.removeItem('leon_access_token')
  })

  it('fetches transactions for the user', async () => {
    await getTransactions()
    expect(axios.get).toHaveBeenCalledWith('http://localhost:4000/transactions', config)
  })

  it('extracts transactions from response', async () => {
    const transactions = await getTransactions()

    expect(transactions).toEqual([
      {
        sender_name: 'Alice',
        amount_transferred: 100,
        amount_transferred_currency: 'SGD',
        beneficiary_name: 'Bob',
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
    ])

    expect(transactions).toHaveLength(2)
  })
})
