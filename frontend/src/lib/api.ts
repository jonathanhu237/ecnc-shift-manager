import axios, { AxiosResponse } from "axios";

export type APIResponse<T> = {
  code: number;
  message: string;
  data: T;
};

export const api = axios.create({
  baseURL: "/api",
  withCredentials: true,
});

api.interceptors.response.use(
  (response: AxiosResponse<APIResponse<unknown>>) => {
    const { code, message } = response.data;

    if (code === 0) {
      return response;
    } else {
      return Promise.reject(new Error(message));
    }
  }
);
