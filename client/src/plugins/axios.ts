import axios, { AxiosError } from 'axios';
import type { App } from 'vue';

const baseURL = import.meta.env.VITE_APP_BASE_URL;

axios.create({
  baseURL,
  headers: {
    'Content-type': 'application/json',
  },
  validateStatus: () => true,
});

// Axios Interceptor
// https://stackoverflow.com/questions/72370102/axios-interceptors-with-vue-3-and-vite-not-working
export default {
  install: (app: App): void => {
    app.config.globalProperties.$http = axios;
    const $http = app.config.globalProperties.$http;

    const handleResponse = (response: any) => {
      // Here your code
      console.log('response in interceptor', response);
      return response;
    };

    const handleError = (err: AxiosError) => {
      if (err.response?.status === 401) {
        const loginUrl = `${baseURL}/login`;
        window.location.replace(loginUrl);
      }
      return err;
    };

    $http.interceptors.response.use(handleResponse, handleError);
  },
};
