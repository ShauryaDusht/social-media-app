const CONFIG = {
    API_URL: '/api',
    getToken() {
        return localStorage.getItem('token');
    },
    getUser() {
        return JSON.parse(localStorage.getItem('user'));
    },
    checkAuth() {
        const token = this.getToken();
        if (!token) {
            window.location.href = '/login';
            return false;
        }
        return true;
    },
    logout() {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = '/login';
    }
};