// import Cookies from 'js-cookie';
import axios, {AxiosInstance, AxiosResponse} from 'axios';
import {ShortProject} from "@/types/ShortProject.ts";
import {Project} from "@/types/Project.ts";


class Api {
    private axiosInstance: AxiosInstance;
    // private COOKIE_NAME = 'token';
    // private LOGIN_PATH = '/login';

    constructor() {
        this.axiosInstance = axios.create({
            baseURL: import.meta.env.VITE_RECEIVER_HOST,
            timeout: 50000,
            headers: {
                'Content-Type': 'application/json',
                // 'Authorization': 'Bearer ' + Cookies.get(this.COOKIE_NAME)
            },
        });

        // this.axiosInstance.interceptors.request.use((config) => {
        //         const token = Cookies.get(this.COOKIE_NAME);
        //         if (token) {
        //             config.headers.Authorization = 'Bearer ' + token;
        //         }
        //         return config;
        //     },
        //     (error) => {
        //         return Promise.reject(error);
        //     }
        // )
        //
        // this.axiosInstance.interceptors.response.use(
        //     (response) => {
        //         return response;
        //     },
        //     (error) => {
        //         if (error.response.status === 401 && window.location.pathname !== this.LOGIN_PATH) {
        //             window.location.href = this.LOGIN_PATH;
        //         }
        //         return Promise.reject(error);
        //     }
        // )
    }

    public createProject(name: string, startUrl: string) {
        console.log(import.meta.env.VITE_RECEIVER_HOST);
        return this.axiosInstance.post('/project/create', {"name": name, "start_url": startUrl});
    }

    public getProject(id: string): Promise<AxiosResponse<Project>> {
        return this.axiosInstance.get('/project/get/' + id);
    }

    public getAllProjectsShort(): Promise<AxiosResponse<ShortProject[]>> {
        return this.axiosInstance.get('/project/getAllShort');
    }

    public deleteProject(id: string): Promise<AxiosResponse<ShortProject>> {
        return this.axiosInstance.delete('/project/delete' + id);
    }
}

export default new Api();