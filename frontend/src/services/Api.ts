// import Cookies from 'js-cookie';
import axios from 'axios';

class Api {
    private axiosInstance: axios.AxiosInstance;
    // private COOKIE_NAME = 'token';
    // private LOGIN_PATH = '/login';

    constructor() {
        this.axiosInstance = axios.create({
            baseURL:  process.env.RECEIVER_PORT,
            timeout: 1000,
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
        return this.axiosInstance.post('/project/create', {"name": name, "start_url": startUrl});
    }

    public getProject(id: string) {
        return this.axiosInstance.get('/project/get/' + id);
    }

    public getAllProjectsShort() {
        return this.axiosInstance.get('/project/getAllShort', );
    }

    public deleteProject(id: string) {
        return this.axiosInstance.delete('/project/delete' + id);
    }
}

export default new Api();