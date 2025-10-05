// auth.js — логика авторизации и регистрации
// Бэкенд: внешний сервис
const API_BASE = 'http://81.29.146.35:8080'; // NOTE: при необходимости вынести в конфиг

(function() {
  const authBtn = document.getElementById('auth-btn');
  const authLabel = authBtn?.querySelector('.auth-btn__label');
  const dialog = document.getElementById('auth-dialog');
  const overlay = document.getElementById('auth-overlay');
  const closeBtn = document.getElementById('auth-close');
  const tabs = Array.from(document.querySelectorAll('.auth-tab'));
  const form = document.getElementById('auth-form');
  const errorsBox = document.getElementById('auth-errors');
  const logoutBtn = document.getElementById('auth-logout');
  const nameGroup = form?.querySelector('[data-field="name"]');
  const roleGroup = form?.querySelector('[data-field="role"]');

  let mode = 'login'; // 'login' | 'register'

  function getToken() { return localStorage.getItem('auth_token'); }
  function setToken(t) { localStorage.setItem('auth_token', t); }
  function clearToken() { localStorage.removeItem('auth_token'); }

  function openDialog() {
    dialog?.classList.add('open');
    dialog?.removeAttribute('aria-hidden');
    overlay?.classList.add('visible');
    overlay?.removeAttribute('aria-hidden');
    authBtn?.setAttribute('aria-expanded','true');
  }
  function closeDialog() {
    dialog?.classList.remove('open');
    dialog?.setAttribute('aria-hidden','true');
    overlay?.classList.remove('visible');
    overlay?.setAttribute('aria-hidden','true');
    authBtn?.setAttribute('aria-expanded','false');
  }

  function switchMode(newMode) {
    if (mode === newMode) return;
    mode = newMode;
    tabs.forEach(t => {
      const active = t.dataset.mode === mode;
      t.classList.toggle('active', active);
      t.setAttribute('aria-selected', active ? 'true' : 'false');
    });
    if (mode === 'register') {
      nameGroup.style.display = '';
      roleGroup.style.display = '';
    } else {
      nameGroup.style.display = 'none';
      roleGroup.style.display = 'none';
    }
    errorsBox.textContent = '';
    form.reset();
  }

  function updateAuthUI() {
    const token = getToken();
    if (token) {
      authLabel.textContent = 'Профиль';
      logoutBtn.style.display = 'inline-flex';
    } else {
      authLabel.textContent = 'Войти';
      logoutBtn.style.display = 'none';
    }
  }

  async function request(path, data) {
    const resp = await fetch(API_BASE + path, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    });
    const json = await resp.json().catch(() => ({}));
    if (!resp.ok) {
      throw new Error(json.error || 'Ошибка сервера');
    }
    return json;
  }

  async function handleSubmit(e) {
    e.preventDefault();
    errorsBox.textContent = '';
    const email = form.email.value.trim();
    const password = form.password.value;
    if (!email || !password) {
      errorsBox.textContent = 'Заполните email и пароль';
      return;
    }

    try {
      if (mode === 'login') {
        const res = await request('/login', { email, password });
        if (res.token) setToken(res.token);
      } else {
        const name = form.name.value.trim();
        const role = form.role.value || 'user';
        if (!name) { errorsBox.textContent = 'Введите имя'; return; }
        const res = await request('/register', { name, email, password, role });
        if (res.token) setToken(res.token);
      }
      updateAuthUI();
      closeDialog();
    } catch (err) {
      errorsBox.textContent = err.message;
    }
  }

  function logout() {
    clearToken();
    updateAuthUI();
  }

  // Events
  authBtn?.addEventListener('click', () => {
    if (getToken()) { // если уже авторизован — открываем диалог как профиль
      openDialog();
    } else {
      switchMode('login');
      openDialog();
    }
  });
  closeBtn?.addEventListener('click', closeDialog);
  overlay?.addEventListener('click', closeDialog);
  tabs.forEach(tab => tab.addEventListener('click', () => switchMode(tab.dataset.mode)));
  form?.addEventListener('submit', handleSubmit);
  logoutBtn?.addEventListener('click', () => { logout(); closeDialog(); });
  document.addEventListener('keydown', (e) => { if (e.key === 'Escape') closeDialog(); });

  // Init
  updateAuthUI();
})();
