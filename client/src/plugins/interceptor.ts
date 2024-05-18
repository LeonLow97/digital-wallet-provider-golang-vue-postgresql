import axios, { AxiosError } from "axios";
import type { AxiosRequestConfig, AxiosResponse } from "axios";
import type { App } from "vue";
import { useUserStore } from "@/stores/user";

const baseURL = import.meta.env.VITE_APP_BASE_URL;

axios.create({
  baseURL,
  withCredentials: true,
  headers: {
    "Content-type": "application/json",
  },
});

// Axios Interceptor for response
// https://stackoverflow.com/questions/72370102/axios-interceptors-with-vue-3-and-vite-not-working
export default {
  install: (app: App): void => {
    app.config.globalProperties.$http = axios;
    const $http = app.config.globalProperties.$http;

    const handleResponse = (response: AxiosResponse) => {
      // Check if the CSRF Token is in the response headers
      const csrfToken = response.headers["x-csrf-token"];
      if (csrfToken) {
        useUserStore().STORE_CSRF_TOKEN(csrfToken);
      }

      return response;
    };

    const handleResponseError = (err: AxiosError) => {
      if (err.response?.status === 401) {
        useUserStore().LOGOUT_USER();
        const loginUrl = `${baseURL}/login`;
        window.location.replace(loginUrl);
      }
      return Promise.reject(err);
    };

    const handleRequest = (request: AxiosRequestConfig) => {
      const storeCsrfToken = useUserStore().csrfToken;
      if (storeCsrfToken) {
        if (!request.headers) {
          request.headers = {};
        }
        request.headers["X-CSRF-Token"] = storeCsrfToken;
      }

      return request;
    };

    const handleRequestError = (err: AxiosError) => {
      return Promise.reject(err);
    };

    $http.interceptors.response.use(handleResponse, handleResponseError);
    $http.interceptors.request.use(handleRequest, handleRequestError);
  },
};
