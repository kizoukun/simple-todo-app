const BASE_API_URL = "http://localhost:8080/api/v1";


document.addEventListener("DOMContentLoaded", onLoad)

let todos = [];
let boardId = null;

let BoardTitle = '';
let BoardDesc = '';

const titleEl = document.getElementById('boardTitle');
const descEl = document.getElementById('boardDescription');

async function onLoad() {
    const urlParams = new URLSearchParams(window.location.search);
    boardId = urlParams.get('id');
    if (!boardId) {
        window.location.href = "./boards/index.html";
        return;
    }
    await refetchTodoItems();
}

async function refetchTodoItems() {
    const response = await fetchTodoItems();
    if (response) {
        const board = response.data.board;
        titleEl.textContent = board.title;
        descEl.textContent = board.description;
        BoardTitle = board.title;
        BoardDesc = board.description;
        todos = response.data.todos;
        renderTodos();
        updateStats();
    }
}

async function fetchTodoItems() {
    try {
        const response = await fetch(`${BASE_API_URL}/todos/${boardId}`, {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`,
            }
        });
        if (response.status === 401) {
            window.location.href = "./auth/auth.html";
            return null;
        }
        if (response.status === 404 || response.status === 400) {
            window.location.href = "./boards/index.html";
            return null;
        }
        const todoItems = await response.json();
        return todoItems;
    } catch (error) {
        console.error("Error fetching todo items:", error);
    }
}

const todoInput = document.getElementById('todoInput');
const todosList = document.getElementById('todosList');
const totalCount = document.getElementById('totalCount');
const completedCount = document.getElementById('completedCount');
const progressFill = document.getElementById('progressFill');

function updateStats() {
    const total = todos.length;
    const completed = todos.filter(t => t.completed).length;
    const progress = total === 0 ? 0 : (completed / total) * 100;

    totalCount.textContent = total;
    completedCount.textContent = completed;
    progressFill.style.width = progress + '%';
}

function renderTodos() {
    if (todos.length === 0) {
        todosList.innerHTML = `
                    <div class="empty-state">
                        <div class="empty-state-icon">✓</div>
                        <p>No tasks yet. Add one to get started!</p>
                    </div>
                `;
        return;
    }

    todosList.innerHTML = todos.map((todo, index) => `
                <div class="todo-item ${todo.completed ? 'completed' : ''}">
                    <input 
                        type="checkbox" 
                        class="checkbox" 
                        ${todo.completed ? 'checked' : ''}
                        onchange="toggleTodo(${todo.id})"
                    >
                    <span class="todo-text">${escapeHtml(todo.task)}</span>
                    <button class="btn-delete" onclick="deleteTodo(${todo.id})">×</button>
                </div>
            `).join('');
}

async function addTodo() {
    const text = todoInput.value.trim();

    if (text === '') return;

    try {
        const response = await fetch(`${BASE_API_URL}/todos/${boardId}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`,
            },
            body: JSON.stringify({ task: text })
        });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
    } catch (error) {
        console.error("Error fetching todo items:", error);
    }
    todoInput.value = '';
    todoInput.focus();
    refetchTodoItems();
}

async function toggleTodo(id) {
    const todo = todos.find(t => t.id === id);
    if (!todo) {
        alert("Todo item not found");
        return;
    }

    try {
        const response = await fetch(`${BASE_API_URL}/todos/${boardId}/${id}/complete`, {
            method: 'POST',
            body: JSON.stringify({ completed: !todo.completed }),
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`,
            }
        });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        refetchTodoItems();
    } catch (error) {
        console.error("Error toggling todo item:", error);
    }
}

async function deleteTodo(id) {
    try {
        const response = await fetch(`${BASE_API_URL}/todos/${boardId}/${id}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`
            }
        });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        refetchTodoItems();
    } catch (error) {
        console.error("Error deleting todo item:", error);
    }
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

todoInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') addTodo();
});

async function addUser() {
    const email = prompt("Enter the email of the user to add:");
    if (!email) return;

    try {
        const response = await fetch(`${BASE_API_URL}/todos/${boardId}/invite`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`,
            },
            body: JSON.stringify({ email: email })
        });
        const data = await response.json();
        if (!data.success) {
            alert("Failed to add user: " + data.message);
            return;
        }
        await new Promise(resolve => setTimeout(resolve, 100));
        await copyclipBoard(data.data.invite_code);
        alert("User invited with this code (dev)!: " + data.data.invite_code);
    } catch (error) {
        console.error("Error adding user:", error);
    }
}


async function copyclipBoard(code) {
    // Try to ensure the document has focus (may help with NotAllowedError)
    try { window.focus(); } catch (e) { /* ignore */ }

    // Primary attempt: Clipboard API
    async function tryClipboardWrite(text) {
        if (navigator.clipboard && window.isSecureContext) {
            return navigator.clipboard.writeText(text);
        }
        // If Clipboard API not available, reject to trigger fallback
        return Promise.reject(new Error('Clipboard API unavailable'));
    }

    // Fallback using a temporary textarea + execCommand('copy')
    function fallbackCopyTextToClipboard(text) {
        const ta = document.createElement('textarea');
        ta.value = text;
        // Prevent scrolling to bottom
        ta.style.position = 'fixed';
        ta.style.left = '-9999px';
        document.body.appendChild(ta);
        ta.focus();
        ta.select();
        let ok = false;
        try {
            ok = document.execCommand('copy');
        } catch (err) {
            ok = false;
        } finally {
            document.body.removeChild(ta);
        }
        return ok;
    }

    // Try clipboard write, then fallback; show UI if all fail
    try {
        await tryClipboardWrite(code);
        alert("User invited — invite code copied to clipboard: " + code);
    } catch (err) {
        // If writeText rejected due to focus or permission, try fallback
        const ok = fallbackCopyTextToClipboard(code);
    }
}

async function deleteBoard() {
    if (!confirm("Are you sure you want to delete this board? This action cannot be undone.")) {
        return;
    }

    try {
        const response = await fetch(`${BASE_API_URL}/boards/${boardId}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`,
            },
        });
        const data = await response.json();
        if (!data.success) {
            alert("Failed to delete board: " + data.message);
            return;
        }
        window.location.href = `./boards/index.html`;
    } catch (error) {
        console.error("Error deleting board:", error);
    }
}


// Title: save on Enter (prevent newline) and blur
titleEl.addEventListener('keydown', (e) => {
    if (e.key === 'Enter') {
        e.preventDefault(); // prevent newline
        titleEl.blur();     // trigger blur -> save
    }
});
titleEl.addEventListener('blur', saveTitle);

// Description: save on blur
descEl.addEventListener('blur', saveDesc);

async function saveTitle() {
    try {
        const title = titleEl.textContent.trim() || '';
        await saveBoard(title, BoardDesc);
        BoardTitle = title;
    } catch (err) {
        console.error("Error saving title:", err);
    }
}
async function saveDesc() {
    try {
        const description = descEl.textContent.trim() || '';
        await saveBoard(BoardTitle, description);
        BoardDesc = description;
    } catch (err) {
        console.error("Error saving description:", err);
    }
}


async function saveBoard(title, description) {
    const response = await fetch(`${BASE_API_URL}/boards/${boardId}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('authToken')}`,
        },
        body: JSON.stringify({ title: title, description: description })
    });
    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }
}