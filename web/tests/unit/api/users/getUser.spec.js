import axios from 'axios'

import getUser from '@/api/users/getUser'

vi.mock('axios')

describe('getUser', () => {
  const config = {
    headers: { Authorization: `Bearer JWT_Token_Test_Value` }
  }

  beforeEach(() => {
    axios.get.mockResolvedValue({
      data: [
        {
          username: 'Alice',
          balance: 68677,
          currency: 'SGD'
        }
      ]
    })

    localStorage.setItem('leon_access_token', 'JWT_Token_Test_Value')
  })

  afterEach(() => {
    localStorage.removeItem('leon_access_token')
  })

  it('fetches user details for the user', async () => {
    await getUser()
    expect(axios.get).toHaveBeenCalledWith('http://localhost:4000/user', config)
  })

  it('fetches user details to display in the home page his balance', async () => {
    const user = await getUser()

    expect(user).toEqual([
      {
        username: 'Alice',
        balance: 68677,
        currency: 'SGD'
      }
    ])
  })
})
