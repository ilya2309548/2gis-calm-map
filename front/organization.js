// organization.js — создание организации и загрузка файлов
// Требует наличия auth.js (для токена) и заполненного address (клик по карте)

const ORG_API_BASE = 'http://81.29.146.35:8080';

(function(){
  const openFeedbackBtn = document.getElementById('open-feedback'); // можем использовать для запуска после выбора адреса
  const orgDialog = document.getElementById('org-dialog');
  const orgOverlay = document.getElementById('org-overlay');
  const closeBtn = document.getElementById('org-close');
  const cancelBtn = document.getElementById('org-cancel');
  const form = document.getElementById('org-form');
  const errorsBox = document.getElementById('org-errors');
  const addrInput = document.getElementById('org-address');
  const latInput = document.getElementById('org-lat');
  const lonInput = document.getElementById('org-lon');
  const mapFileInput = document.getElementById('org-map-file');
  const picFileInput = document.getElementById('org-picture-file');
  const submitBtn = document.getElementById('org-submit');
  const addressPill = document.getElementById('selected-address');
  const openOrgBtn = document.getElementById('open-org-dialog');

  function getToken(){ return localStorage.getItem('auth_token'); }

  function openDialog(){
    orgDialog.classList.add('open');
    orgDialog.removeAttribute('aria-hidden');
    orgOverlay.classList.add('visible');
    orgOverlay.removeAttribute('aria-hidden');
  }
  function closeDialog(){
    orgDialog.classList.remove('open');
    orgDialog.setAttribute('aria-hidden','true');
    orgOverlay.classList.remove('visible');
    orgOverlay.setAttribute('aria-hidden','true');
  }

  // Публичная функция для открытия с предзаполнением (может вызываться из других скриптов при необходимости)
  window.openOrganizationDialog = function(){
    const text = addressPill?.textContent || '';
    if (!text || text.startsWith('Выберите') || addressPill.dataset.error === '1') {
      errorsBox.textContent = 'Сначала выберите корректный адрес на карте.';
      return;
    }
    // В координатах лежит последний маркер
    if (window.APP?.selectedMarker) {
      const coords = window.APP.selectedMarker.getCoordinates(); // [lon, lat]
      lonInput.value = coords[0].toFixed(6);
      latInput.value = coords[1].toFixed(6);
    }
    addrInput.value = text;
    errorsBox.textContent = '';
    openDialog();
  }

  // Можно привязать, например, к клику на адрес
  // Кнопка явного открытия
  openOrgBtn?.addEventListener('click', () => {
    window.openOrganizationDialog();
  });

  orgOverlay.addEventListener('click', closeDialog);
  closeBtn.addEventListener('click', closeDialog);
  cancelBtn.addEventListener('click', closeDialog);
  document.addEventListener('keydown', (e)=>{ if(e.key==='Escape') closeDialog(); });

  async function apiJSON(path, method, body){
    const token = getToken();
    const res = await fetch(ORG_API_BASE + path, {
      method,
      headers: {
        'Content-Type': 'application/json',
        ...(token ? { 'Authorization': 'Bearer ' + token } : {})
      },
      body: body ? JSON.stringify(body) : undefined
    });
    let json = null;
    try { json = await res.json(); } catch { json = {}; }
    if (!res.ok) throw new Error(json.error || res.statusText || 'Ошибка запроса');
    return json;
  }

  async function apiUpload(path, file, fieldName){
    if (!file) return;
    const token = getToken();
    const fd = new FormData();
    fd.append(fieldName, file);
    const res = await fetch(ORG_API_BASE + path, {
      method: 'POST',
      headers: {
        ...(token ? { 'Authorization': 'Bearer ' + token } : {})
      },
      body: fd
    });
    if (!res.ok) {
      let txt = await res.text();
      throw new Error('Upload failed: ' + txt.slice(0,200));
    }
  }

  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    errorsBox.textContent='';
    submitBtn.disabled = true;
    submitBtn.textContent = 'Создание...';

    const address = addrInput.value.trim();
    const organization_type = form.organization_type.value;
    const lat = parseFloat(latInput.value);
    const lon = parseFloat(lonInput.value);

    if (!getToken()) {
      errorsBox.textContent = 'Нужна авторизация.';
      submitBtn.disabled = false; submitBtn.textContent='Создать';
      return;
    }
    if (!address) {
      errorsBox.textContent = 'Адрес пуст.';
      submitBtn.disabled = false; submitBtn.textContent='Создать';
      return;
    }

    try {
      const org = await apiJSON('/organization', 'POST', {
        address,
        organization_type,
        latitude: isFinite(lat)? lat : undefined,
        longitude: isFinite(lon)? lon : undefined
      });
      const orgId = org.id;
      // Загрузка файлов (не обязательны)
      const mapFile = mapFileInput.files[0];
      const picFile = picFileInput.files[0];
      if (mapFile) {
        await apiUpload(`/organization/${orgId}/map/upload`, mapFile, 'file');
      }
      if (picFile) {
        await apiUpload(`/organization/${orgId}/picture/upload`, picFile, 'file');
      }
      submitBtn.textContent = 'Готово';
      setTimeout(closeDialog, 600);
    } catch (err) {
      errorsBox.textContent = err.message;
    } finally {
      setTimeout(()=>{ submitBtn.disabled = false; submitBtn.textContent='Создать'; }, 800);
    }
  });
})();
