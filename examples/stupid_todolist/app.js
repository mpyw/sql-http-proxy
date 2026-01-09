const API_BASE = 'http://localhost:8090/api';

const todoForm = document.getElementById('todo-form');
const todoInput = document.getElementById('todo-input');
const searchInput = document.getElementById('search-input');
const filterButtons = document.querySelectorAll('.filter-btn');
const todoList = document.getElementById('todo-list');
const statusDiv = document.getElementById('status');

let currentFilter = 'all';

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

// Load todos with optional search and filter
async function loadTodos(search = '', filter = 'all') {
  try {
    todoList.classList.add('loading');
    const params = new URLSearchParams();
    if (search) params.set('q', search);
    if (filter && filter !== 'all') params.set('filter', filter);
    const query = params.toString() ? `?${params}` : '';
    const todos = await api(`/todos${query}`);
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
    loadTodos(searchInput.value, currentFilter);
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
    loadTodos(searchInput.value, currentFilter);
  } catch (err) {
    showStatus('Failed to update: ' + err.message, true);
    loadTodos(searchInput.value, currentFilter);
  }
}

// Delete todo
async function deleteTodo(id) {
  try {
    await api(`/todos/${id}`, { method: 'DELETE' });
    loadTodos(searchInput.value, currentFilter);
  } catch (err) {
    showStatus('Failed to delete: ' + err.message, true);
  }
}

// Debounce helper
function debounce(fn, delay) {
  let timeout;
  return (...args) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => fn(...args), delay);
  };
}

// Event listeners
todoForm.addEventListener('submit', (e) => {
  e.preventDefault();
  const title = todoInput.value.trim();
  if (title) addTodo(title);
});

searchInput.addEventListener('input', debounce((e) => {
  loadTodos(e.target.value, currentFilter);
}, 300));

filterButtons.forEach(btn => {
  btn.addEventListener('click', () => {
    filterButtons.forEach(b => b.classList.remove('active'));
    btn.classList.add('active');
    currentFilter = btn.dataset.filter;
    loadTodos(searchInput.value, currentFilter);
  });
});

// Initial load
loadTodos();
