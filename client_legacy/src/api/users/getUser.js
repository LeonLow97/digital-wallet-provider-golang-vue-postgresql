import axios from 'axios'

const getUser = async () => {
  const jwt_token = localStorage.getItem('leon_access_token')
  const config = {
    headers: { Authorization: `Bearer ${jwt_token}` }
  }
  const url = 'http://localhost:4000/user'
  const response = await axios.get(url, config)
  return response.data
}

export default getUser
