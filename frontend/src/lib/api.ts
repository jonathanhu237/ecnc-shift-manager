import axios from "axios";

export type APIResponse = {
    code: number;
    message: string;
    data: unknown;
};

export const api = axios.create({
    baseURL: process.env.API_SERVER_URL,
    withCredentials: true,
});
