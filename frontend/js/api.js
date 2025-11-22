const API_BASE_URL = ''; // Relative path for Nginx proxy

class ApiClient {
    constructor() {
        this.token = localStorage.getItem('token');
    }

    setToken(token) {
        this.token = token;
        localStorage.setItem('token', token);
    }

    clearToken() {
        this.token = null;
        localStorage.removeItem('token');
    }

    isAuthenticated() {
        return !!this.token;
    }

    async request(endpoint, method = 'GET', data = null) {
        const headers = {
            'Content-Type': 'application/json',
        };

        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        const config = {
            method,
            headers,
        };

        if (data) {
            config.body = JSON.stringify(data);
        }

        try {
            const response = await fetch(`${API_BASE_URL}${endpoint}`, config);
            const responseData = await response.json();

            if (!response.ok) {
                throw new Error(responseData.error || responseData.message || 'Something went wrong');
            }

            return responseData;
        } catch (error) {
            throw error;
        }
    }

    async login(username, password) {
        // The backend expects username, email (optional), password
        // Based on models/login.go
        const data = { username, password, email: "" };
        return this.request('/v1/login', 'POST', data);
    }

    async register(userData) {
        return this.request('/v1/register', 'POST', userData);
    }

    async getUser() {
        return this.request('/v1/profile', 'GET');
    }

    async shortenUrl(originalUrl, validForInMonths = 1) {
        return this.request('/shortcode', 'POST', {
            original_url: originalUrl,
            valid_for_in_months: parseInt(validForInMonths)
        });
    }

    async customShortenUrl(originalUrl, shortCode, validForInMonths = 1) {
        return this.request('/custom/shortcode', 'POST', {
            original_url: originalUrl,
            short_code: shortCode,
            valid_for_in_months: parseInt(validForInMonths)
        });
    }
}

const api = new ApiClient();
