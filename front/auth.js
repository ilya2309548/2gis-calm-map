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

  // Простая декодировка payload из JWT (без проверки подписи)
  function decodeJWT(token){
    try {
      const base = token.split('.')[1];
      const json = atob(base.replace(/-/g,'+').replace(/_/g,'/'));
      return JSON.parse(decodeURIComponent(escape(json)));
    } catch { return null; }
  }

  const profilePanel = document.getElementById('user-profile');
  const profName = document.getElementById('prof-name');
  const profEmail = document.getElementById('prof-email');
  const profRole = document.getElementById('prof-role');
  const paramsForm = document.getElementById('user-params-form');
  const saveParamsBtn = document.getElementById('save-user-params');
  const paramsStatus = document.getElementById('user-params-status');

  function clearParamsStatus(){ paramsStatus.className='params-status'; paramsStatus.textContent=''; }
  function setParamsStatus(txt, kind){ paramsStatus.textContent=txt; paramsStatus.className='params-status '+(kind||''); }

  async function fetchUserParams(userId){
    if (!userId) return;
    const token = getToken(); if(!token) return;
    try {
      const resp = await fetch(`${API_BASE}/user-params/${userId}`, { headers:{ Authorization:'Bearer '+token }});
      if (!resp.ok) { return; }
      const data = await resp.json();
      const map = {
        appearance:'appearance', lighting:'lighting', smell:'smell', temperature:'temperature', tactility:'tactility', signage:'signage', intuitiveness:'intuitiveness', staff_attitude:'staff_attitude', people_density:'people_density', self_service:'self_service', calmness:'calmness'
      };
      Object.entries(map).forEach(([formName, apiName]) => {
        const input = paramsForm.querySelector(`input[name="${apiName}"]`);
        if (input && typeof data[apiName] === 'boolean') input.checked = !!data[apiName];
      });
    } catch(e){ /* ignore */ }
  }

  async function saveUserParams(){
    clearParamsStatus();
    const token = getToken(); if(!token){ setParamsStatus('Нет токена','error'); return; }
    const dec = decodeJWT(token); if(!dec || !dec.user_id){ setParamsStatus('Нет user_id','error'); return; }
    const payload = {};
    paramsForm.querySelectorAll('input[type="checkbox"]').forEach(ch => { payload[ch.name] = ch.checked; });
    setParamsStatus('Сохранение...','saving');
    try {
      // Всегда PATCH (запись создаётся автоматически через GET или уже существует); POST только если сервер неожиданно вернул 404
      let resp = await fetch(`${API_BASE}/user-params/${dec.user_id}`, { method:'PATCH', headers:{'Content-Type':'application/json', Authorization:'Bearer '+token}, body: JSON.stringify(payload)});
      if (resp.status === 404) {
        // fallback (теоретически не нужен, но на всякий случай)
        resp = await fetch(`${API_BASE}/user-params`, { method:'POST', headers:{'Content-Type':'application/json', Authorization:'Bearer '+token}, body: JSON.stringify(payload)});
      }
      if(!resp.ok){ const txt = await resp.text(); throw new Error(txt.slice(0,200)||'Ошибка сохранения'); }
      setParamsStatus('Сохранено','ok');
    } catch(err){ setParamsStatus(err.message||'Ошибка','error'); }
  }

  function updateAuthUI() {
    const token = getToken();
    if (token) {
      authLabel.textContent = 'Профиль';
      logoutBtn.style.display = 'inline-flex';
      // Переходим в режим отображения профиля
      profilePanel?.classList.add('visible');
      document.querySelector('.auth-dialog__tabs')?.classList.add('hidden');
      form.style.display = 'none';
      const dec = decodeJWT(token) || {};
      profName.textContent = dec.name || '—';
      profEmail.textContent = dec.email || '—';
      profRole.textContent = dec.role ? ('Роль: '+dec.role) : '';
      fetchUserParams(dec.user_id);
    } else {
      authLabel.textContent = 'Войти';
      logoutBtn.style.display = 'none';
      profilePanel?.classList.remove('visible');
      document.querySelector('.auth-dialog__tabs')?.classList.remove('hidden');
      form.style.display = '';
      form.reset();
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
  saveParamsBtn?.addEventListener('click', saveUserParams);
  document.addEventListener('keydown', (e) => { if (e.key === 'Escape') closeDialog(); });

  // Init
  updateAuthUI();
})();
