import axios from 'axios'

const getBeneficiaries = async () => {
  const jwt_token = localStorage.getItem('leon_access_token')
  const config = {
    headers: { Authorization: `Bearer ${jwt_token}` }
  }
  const url = 'http://localhost:4000/beneficiaries'
  const response = await axios.get(url, config)
  return response.data
}

export default getBeneficiaries