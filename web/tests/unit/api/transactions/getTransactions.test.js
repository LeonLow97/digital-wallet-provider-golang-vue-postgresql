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
          amount_transferred: 100
        },
        {
          sender_name: 'Alice',
          amount_transferred: 100
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
        amount_transferred: 100
      },
      {
        sender_name: 'Alice',
        amount_transferred: 100
      }
    ])

    expect(transactions).toHaveLength(2)
  })
})
