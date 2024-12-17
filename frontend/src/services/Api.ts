import Cookies from 'js-cookie';
import axios, {AxiosInstance, AxiosResponse} from 'axios';
import {ShortProject} from "@/types/ShortProject.ts";
import {Project} from "@/types/Project.ts";


class Api {
    private axiosInstance: AxiosInstance;
    private COOKIE_NAME = 'access';

    constructor() {
        this.axiosInstance = axios.create({
            baseURL: import.meta.env.VITE_RECEIVER_HOST,
            timeout: 50000,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + Cookies.get(this.COOKIE_NAME)
            },
        });

        this.axiosInstance.interceptors.request.use((config) => {
                const token = Cookies.get(this.COOKIE_NAME);
                if (token) {
                    config.headers.Authorization = 'Bearer ' + token;
                }
                return config;
            },
            (error) => {
                return Promise.reject(error);
            }
        )

        this.axiosInstance.interceptors.response.use(
            (response) => {
                return response;
            },
            (error) => {
                if (error.response.status === 401 && window.location.pathname !== "/") {
                    window.location.href = "/";
                }
                return Promise.reject(error);
            }
        )
    }

    public createProject(name: string, startUrl: string, numberOfLinks: number, depth: number) {
        return this.axiosInstance.post('/project/create', {
            "name": name,
            "start_url": startUrl,
            "number_of_links": numberOfLinks,
            "depth": depth
        });
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

    public getUser(): Promise<AxiosResponse<any>> {
        return this.axiosInstance.get('/user/get');
    }

    public registerUser(username: string, email: string, password: string) {
        return this.axiosInstance.post('/user/register', {"username": username, "email": email, "password": password});
    }

    public loginUser(login: string, password: string) {
        return this.axiosInstance.post('/user/login', {"login": login, "password": password});
    }

    public updateUser(username: string, email: string, newPassword: string, currentPassword: string) {
        return this.axiosInstance.put('/user/update', {
            "username": username,
            "email": email,
            "new_password": newPassword,
            "current_password": currentPassword
        });
    }
}

export default new Api();