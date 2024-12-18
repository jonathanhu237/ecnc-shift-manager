import axios, { AxiosResponse } from "axios";

export type APIResponse = {
    code: number;
    message: string;
    data: unknown;
};

export const api = axios.create({
    baseURL: process.env.API_SERVER_URL,
    withCredentials: true,
});

api.interceptors.response.use((response: AxiosResponse<APIResponse>) => {
    const { code, message } = response.data;

    if (code === 0) {
        return response;
    } else {
        return Promise.reject(new Error(message));
    }
});
