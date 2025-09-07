document.addEventListener('DOMContentLoaded', function() {
    // Проверка авторизации
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = '/login.html';
        return;
    }

    // Отображение роли пользователя
    const role = localStorage.getItem('role') || 'viewer';
    document.getElementById('user-role').textContent = `Роль: ${getRoleName(role)}`;

    // Обработчик выхода
    document.getElementById('logout-btn').addEventListener('click', function() {
        localStorage.removeItem('token');
        localStorage.removeItem('role');
        window.location.href = '/login.html';
    });

    // Инициализация табов
    initTabs();

    // Загрузка данных
    loadItems();
    loadHistory();

    // Обработчики модальных окон
    initModal();
});

function getRoleName(role) {
    const roles = {
        'admin': 'Администратор',
        'manager': 'Менеджер',
        'viewer': 'Просмотрщик'
    };
    return roles[role] || role;
}

// Табы
function initTabs() {
    const tabBtns = document.querySelectorAll('.tab-btn');
    const tabContents = document.querySelectorAll('.tab-content');

    tabBtns.forEach(btn => {
        btn.addEventListener('click', function() {
            const tabName = this.getAttribute('data-tab');

            // Убираем активный класс у всех кнопок и контентов
            tabBtns.forEach(b => b.classList.remove('active'));
            tabContents.forEach(c => c.classList.remove('active'));

            // Добавляем активный класс текущей кнопке и контенту
            this.classList.add('active');
            document.getElementById(`${tabName}-tab`).classList.add('active');
        });
    });
}

// Загрузка товаров
async function loadItems() {
    try {
        const response = await api.get('/items');
        if (response.status === 'OK') {
            renderItems(response.data);
        } else {
            console.error('Ошибка загрузки товаров:', response.error);
        }
    } catch (error) {
        console.error('Ошибка подключения:', error);
    }
}

function renderItems(items) {
    const tbody = document.querySelector('#items-table tbody');
    tbody.innerHTML = '';

    items.forEach(item => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${item.id}</td>
            <td>${item.name}</td>
            <td>${item.quantity}</td>
            <td>${formatDate(item.created_at)}</td>
            <td>${formatDate(item.updated_at)}</td>
            <td class="action-buttons">
                <button class="btn-warning" onclick="editItem(${item.id})">Редактировать</button>
                <button class="btn-danger" onclick="deleteItem(${item.id})">Удалить</button>
            </td>
        `;
        tbody.appendChild(row);
    });
}

// Загрузка истории
async function loadHistory() {
    try {
        const response = await api.get('/history');
        if (response.status === 'OK') {
            renderHistory(response.data);
        } else {
            console.error('Ошибка загрузки истории:', response.error);
        }
    } catch (error) {
        console.error('Ошибка подключения:', error);
    }
}

function renderHistory(history) {
    const tbody = document.querySelector('#history-table tbody');
    tbody.innerHTML = '';

    history.forEach(record => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${record.id}</td>
            <td>${record.item_id}</td>
            <td>${formatAction(record.action)}</td>
            <td>${record.changed_by}</td>
            <td>${formatDate(record.changed_at)}</td>
            <td>
                <button class="details-btn" onclick="showDetails(${record.id})">Подробности</button>
            </td>
        `;
        tbody.appendChild(row);
    });
}

// Модальное окно для товаров
function initModal() {
    const modal = document.getElementById('item-modal');
    const closeBtn = modal.querySelector('.close');
    const form = document.getElementById('item-form');

    // Закрытие модального окна
    closeBtn.addEventListener('click', function() {
        modal.style.display = 'none';
    });

    // Закрытие при клике вне модального окна
    window.addEventListener('click', function(event) {
        if (event.target === modal) {
            modal.style.display = 'none';
        }
    });

    // Обработчик формы
    form.addEventListener('submit', async function(e) {
        e.preventDefault();

        const id = document.getElementById('item-id').value;
        const name = document.getElementById('item-name').value;
        const quantity = parseInt(document.getElementById('item-quantity').value);

        try {
            let response;
            if (id) {
                // Обновление
                response = await api.put(`/items/${id}`, { name, quantity });
            } else {
                // Создание
                response = await api.post('/items', { name, quantity });
            }

            if (response.status === 'OK') {
                modal.style.display = 'none';
                loadItems();
            } else {
                alert('Ошибка: ' + response.error);
            }
        } catch (error) {
            alert('Ошибка подключения к серверу');
        }
    });

    // Кнопка добавления товара
    document.getElementById('add-item-btn').addEventListener('click', function() {
        openItemModal();
    });
}

function openItemModal(item = null) {
    const modal = document.getElementById('item-modal');
    const title = document.getElementById('modal-title');
    const idInput = document.getElementById('item-id');
    const nameInput = document.getElementById('item-name');
    const quantityInput = document.getElementById('item-quantity');

    if (item) {
        // Редактирование
        title.textContent = 'Редактировать товар';
        idInput.value = item.id;
        nameInput.value = item.name;
        quantityInput.value = item.quantity;
    } else {
        // Создание
        title.textContent = 'Добавить товар';
        idInput.value = '';
        nameInput.value = '';
        quantityInput.value = '';
    }

    modal.style.display = 'block';
}

async function editItem(id) {
    try {
        const response = await api.get(`/items/${id}`);
        if (response.status === 'OK') {
            openItemModal(response.data);
        } else {
            alert('Ошибка загрузки товара: ' + response.error);
        }
    } catch (error) {
        alert('Ошибка подключения к серверу');
    }
}

async function deleteItem(id) {
    if (!confirm('Вы уверены, что хотите удалить этот товар?')) {
        return;
    }

    try {
        const response = await api.delete(`/items/${id}`);
        if (response.status === 'OK') {
            loadItems();
        } else {
            alert('Ошибка удаления: ' + response.error);
        }
    } catch (error) {
        alert('Ошибка подключения к серверу');
    }
}

// Показ деталей истории
function showDetails(id) {
    // Для простоты показываем alert, но можно сделать полноценное модальное окно
    alert('Подробности истории #' + id);
}