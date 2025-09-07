// Утилиты для работы с API
const API_BASE = ''; // Пустой базовый путь

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

    async request(url, options = {}) {
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers
        };

        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        const config = {
            ...options,
            headers
        };

        const response = await fetch(API_BASE + url, config);

        if (response.status === 401) {
            this.clearToken();
            window.location.href = '/login.html';
            throw new Error('Unauthorized');
        }

        return response;
    }

    async get(url) {
        const response = await this.request(url, { method: 'GET' });
        return response.json();
    }

    async post(url, data) {
        const response = await this.request(url, {
            method: 'POST',
            body: JSON.stringify(data)
        });
        return response.json();
    }

    async put(url, data) {
        const response = await this.request(url, {
            method: 'PUT',
            body: JSON.stringify(data)
        });
        return response.json();
    }

    async delete(url) {
        const response = await this.request(url, { method: 'DELETE' });
        return response.json();
    }
}

// Глобальный экземпляр API клиента
const api = new ApiClient();

// Форматирование даты
function formatDate(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleString('ru-RU');
}

// Форматирование действия
function formatAction(action) {
    const actions = {
        'create': 'Создание',
        'update': 'Обновление',
        'delete': 'Удаление'
    };
    return actions[action] || action;
}