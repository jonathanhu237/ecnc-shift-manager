import axios, { AxiosResponse } from "axios";

export type APIResponse<T> = {
  success: boolean;
  message: string;
  data: T;
};

export const api = axios.create({
  baseURL: "/api",
  withCredentials: true,
});

api.interceptors.response.use(
  (response: AxiosResponse<APIResponse<unknown>>) => {
    const { success, message } = response.data;

    if (success === true) {
      return response;
    } else {
      return Promise.reject(new Error(message));
    }
  }
);
