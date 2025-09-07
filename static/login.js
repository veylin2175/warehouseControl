document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('loginForm');
    const errorMessage = document.getElementById('error-message');

    loginForm.addEventListener('submit', async function(e) {
        e.preventDefault();

        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;

        try {
            const response = await fetch('/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ username, password })
            });

            const data = await response.json();

            if (data.status === 'OK') {
                localStorage.setItem('token', data.data.token);
                localStorage.setItem('role', document.getElementById('role').value);
                window.location.href = '/';
            } else {
                errorMessage.textContent = data.error || 'Ошибка авторизации';
            }
        } catch (error) {
            errorMessage.textContent = 'Ошибка подключения к серверу';
        }
    });
});