const API_BASE = 'http://localhost:8080/api';

const todoForm = document.getElementById('todo-form');
const todoInput = document.getElementById('todo-input');
const todoList = document.getElementById('todo-list');
const statusDiv = document.getElementById('status');
const initBtn = document.getElementById('init-btn');

// Show status message
function showStatus(message, isError = false) {
  statusDiv.textContent = message;
  statusDiv.className = 'status ' + (isError ? 'error' : 'success');
  setTimeout(() => {
    statusDiv.textContent = '';
    statusDiv.className = 'status';
  }, 3000);
}

// API helpers
async function api(endpoint, options = {}) {
  const res = await fetch(`${API_BASE}${endpoint}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  });
  if (!res.ok && res.status !== 204) {
    const error = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(error.error || error.message || 'Request failed');
  }
  if (res.status === 204) return null;
  return res.json();
}

// Initialize database
async function initDatabase() {
  try {
    initBtn.disabled = true;
    await api('/init', { method: 'POST' });
    showStatus('Database initialized!');
    loadTodos();
  } catch (err) {
    showStatus('Init failed: ' + err.message, true);
  } finally {
    initBtn.disabled = false;
  }
}

// Load all todos
async function loadTodos() {
  try {
    todoList.classList.add('loading');
    const todos = await api('/todos');
    renderTodos(todos);
  } catch (err) {
    showStatus('Failed to load: ' + err.message, true);
  } finally {
    todoList.classList.remove('loading');
  }
}

// Render todos
function renderTodos(todos) {
  todoList.innerHTML = todos.map(todo => `
    <li class="todo-item ${todo.completed ? 'completed' : ''}" data-id="${todo.id}">
      <input type="checkbox" ${todo.completed ? 'checked' : ''} onchange="toggleTodo(${todo.id}, this.checked)">
      <span class="title">${escapeHtml(todo.title)}</span>
      <button class="delete-btn" onclick="deleteTodo(${todo.id})">Delete</button>
    </li>
  `).join('');
}

// Escape HTML
function escapeHtml(text) {
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
}

// Add todo
async function addTodo(title) {
  try {
    todoInput.disabled = true;
    await api('/todos', {
      method: 'POST',
      body: JSON.stringify({ title }),
    });
    todoInput.value = '';
    loadTodos();
  } catch (err) {
    showStatus('Failed to add: ' + err.message, true);
  } finally {
    todoInput.disabled = false;
    todoInput.focus();
  }
}

// Toggle todo
async function toggleTodo(id, completed) {
  try {
    const todo = await api(`/todos/${id}`);
    await api(`/todos/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ ...todo, completed }),
    });
    loadTodos();
  } catch (err) {
    showStatus('Failed to update: ' + err.message, true);
    loadTodos();
  }
}

// Delete todo
async function deleteTodo(id) {
  try {
    await api(`/todos/${id}`, { method: 'DELETE' });
    loadTodos();
  } catch (err) {
    showStatus('Failed to delete: ' + err.message, true);
  }
}

// Event listeners
todoForm.addEventListener('submit', (e) => {
  e.preventDefault();
  const title = todoInput.value.trim();
  if (title) addTodo(title);
});

initBtn.addEventListener('click', initDatabase);

// Initial load
loadTodos();
