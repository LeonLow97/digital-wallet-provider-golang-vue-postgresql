import axios from 'axios';
import type { User } from '@/types/user';

const LOGIN = async (body: { email: string; password: string }) => {
  try {
    const apiURL = `${import.meta.env.VITE_APP_API_URL}/login`;
    console.log(apiURL);
    const { data, status } = await axios.post<User>(
      apiURL,
      JSON.stringify(body),
      {
        withCredentials: true,
      }
    );

    // Return an object containing both data and status
    return { data, status };
  } catch (error) {
    throw error;
  }
};

export default LOGIN;
