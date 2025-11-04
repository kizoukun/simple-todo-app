const BASE_API_URL = "http://localhost:8080/api/v1";


document.addEventListener("DOMContentLoaded", onLoad)

let todos = [];

async function onLoad() {
    await refetchTodoItems();
}

async function refetchTodoItems() {
    const response = await fetchTodoItems();
    if (response) {
        todos = response.data;
        renderTodos();
        updateStats();
    }
}

async function fetchTodoItems() {
    try {
        const response = await fetch(`${BASE_API_URL}/todos`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
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
        const response = await fetch(`${BASE_API_URL}/todos`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
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
        const response = await fetch(`${BASE_API_URL}/todos/${id}/complete`, {
            method: 'POST',
            body: JSON.stringify({ completed: !todo.completed }),
            headers: {
                'Content-Type': 'application/json'
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
        const response = await fetch(`${BASE_API_URL}/todos/${id}`, {
            method: 'DELETE'
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