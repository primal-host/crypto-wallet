package server

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Wallet</title>
<style>
  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
    background: #0f1117;
    color: #e4e4e7;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  header {
    width: 100%;
    padding: 1.5rem 2rem;
    background: #16181d;
    border-bottom: 1px solid #27272a;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  header h1 { font-size: 1.25rem; font-weight: 600; }
  .header-right { display: flex; align-items: center; gap: 1rem; }
  .header-right .version { color: #71717a; font-size: 0.875rem; }
  main {
    width: 100%;
    max-width: 72rem;
    padding: 2rem;
    flex: 1;
  }

  /* Wallet identity bar */
  .wallet-bar {
    display: flex;
    align-items: center;
    gap: 1rem;
    margin-bottom: 2rem;
    padding: 1rem 1.25rem;
    background: #16181d;
    border: 1px solid #27272a;
    border-radius: 0.5rem;
  }
  .wallet-bar .bar-left {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    flex: 1;
    min-width: 0;
  }
  .wallet-bar .bar-right {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-shrink: 0;
  }
  .wallet-bar .label { color: #71717a; font-size: 0.875rem; white-space: nowrap; }
  .wallet-bar .address {
    font-family: monospace;
    font-size: 0.875rem;
    color: #a1a1aa;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .wallet-bar .no-wallet { color: #71717a; font-size: 0.875rem; font-style: italic; }
  .wallet-bar .key-badge {
    font-size: 0.6875rem;
    color: #a1a1aa;
    background: #27272a;
    padding: 0.125rem 0.5rem;
    border-radius: 0.75rem;
    white-space: nowrap;
  }
  .wallet-bar .method-badge {
    font-size: 0.6875rem;
    color: #71717a;
    background: #1e1e22;
    padding: 0.125rem 0.5rem;
    border-radius: 0.75rem;
    white-space: nowrap;
  }
  .wallet-bar .lock-icon {
    font-size: 1rem;
    margin-right: 0.25rem;
  }

  /* Key selector */
  .key-selector {
    background: #0f1117;
    border: 1px solid #27272a;
    border-radius: 0.25rem;
    color: #e4e4e7;
    font-size: 0.75rem;
    padding: 0.25rem 0.5rem;
    font-family: monospace;
    max-width: 10rem;
  }

  /* Endpoint cards grid */
  .endpoints {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(20rem, 1fr));
    gap: 1rem;
  }
  .ep-card {
    background: #16181d;
    border: 1px solid #27272a;
    border-radius: 0.5rem;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }
  .ep-card-header {
    padding: 1rem 1.25rem;
    border-bottom: 1px solid #1e1e22;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  .ep-card-header h3 { font-size: 1rem; font-weight: 600; }
  .ep-card-body {
    padding: 1rem 1.25rem;
    display: flex;
    flex-direction: column;
    gap: 0.625rem;
    flex: 1;
  }
  .ep-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 0.8125rem;
  }
  .ep-row .label { color: #71717a; }
  .ep-row .value { font-family: monospace; font-size: 0.8rem; color: #a1a1aa; }
  .ep-row .value.balance { color: #e4e4e7; font-weight: 600; font-size: 0.9rem; }

  /* Status dot */
  .status-dot {
    display: inline-block;
    width: 8px;
    height: 8px;
    border-radius: 50%;
    margin-right: 0.375rem;
  }
  .status-online .status-dot { background: #4ade80; }
  .status-offline .status-dot { background: #f87171; }
  .status-checking .status-dot { background: #facc15; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.4; } }

  .status-text { font-size: 0.75rem; }
  .status-online .status-text { color: #4ade80; }
  .status-offline .status-text { color: #f87171; }
  .status-checking .status-text { color: #facc15; }

  /* URL display */
  .url-display {
    font-family: monospace;
    font-size: 0.75rem;
    color: #52525b;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    max-width: 14rem;
  }

  /* Buttons */
  .btn {
    padding: 0.35rem 0.75rem;
    border: 1px solid #27272a;
    border-radius: 0.25rem;
    background: #27272a;
    color: #e4e4e7;
    font-size: 0.75rem;
    cursor: pointer;
    transition: background 0.15s;
    white-space: nowrap;
  }
  .btn:hover { background: #3f3f46; }
  .btn:disabled { opacity: 0.5; cursor: not-allowed; }
  .btn-primary {
    background: #1d4ed8;
    border-color: #1d4ed8;
  }
  .btn-primary:hover { background: #2563eb; }
  .btn-primary:disabled { background: #1e3a5f; border-color: #1e3a5f; }

  /* Setup choices */
  .setup-choices {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    margin-top: 0.5rem;
  }
  .setup-choice {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 0.75rem 1rem;
    background: #0f1117;
    border: 1px solid #27272a;
    border-radius: 0.375rem;
    cursor: pointer;
    transition: border-color 0.15s;
  }
  .setup-choice:hover { border-color: #3f3f46; }
  .setup-choice .choice-icon { font-size: 1.5rem; flex-shrink: 0; }
  .setup-choice .choice-text h4 { font-size: 0.875rem; margin-bottom: 0.125rem; }
  .setup-choice .choice-text p { font-size: 0.75rem; color: #71717a; }

  /* Empty state */
  .empty-state {
    text-align: center;
    padding: 4rem 2rem;
    color: #71717a;
  }
  .empty-state h2 { font-size: 1.125rem; margin-bottom: 0.5rem; color: #a1a1aa; }
  .empty-state p { font-size: 0.875rem; margin-bottom: 1.5rem; }

  /* Modal */
  .modal-overlay {
    display: none;
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.6);
    z-index: 100;
    justify-content: center;
    align-items: center;
  }
  .modal-overlay.active { display: flex; }
  .modal {
    background: #16181d;
    border: 1px solid #27272a;
    border-radius: 0.5rem;
    padding: 1.5rem;
    width: 26rem;
    max-width: 90vw;
  }
  .modal h3 { margin-bottom: 1rem; font-size: 1rem; }
  .modal p {
    font-size: 0.8125rem;
    color: #a1a1aa;
    margin-bottom: 0.75rem;
    line-height: 1.5;
  }
  .modal label {
    display: block;
    font-size: 0.8125rem;
    color: #71717a;
    margin-bottom: 0.25rem;
    margin-top: 0.75rem;
  }
  .modal label:first-of-type { margin-top: 0; }
  .modal input, .modal select {
    width: 100%;
    padding: 0.5rem 0.75rem;
    background: #0f1117;
    border: 1px solid #27272a;
    border-radius: 0.25rem;
    color: #e4e4e7;
    font-size: 0.875rem;
    font-family: inherit;
  }
  .modal input:focus, .modal select:focus {
    outline: none;
    border-color: #1d4ed8;
  }
  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 0.5rem;
    margin-top: 1.25rem;
  }
  .modal-error {
    color: #f87171;
    font-size: 0.8125rem;
    margin-top: 0.5rem;
    display: none;
  }

  /* Hex block number formatting */
  .mono { font-family: monospace; font-size: 0.8rem; }

  /* Latency */
  .latency { font-size: 0.75rem; color: #52525b; }
  .latency.fast { color: #4ade80; }
  .latency.medium { color: #facc15; }
  .latency.slow { color: #fb923c; }
</style>
</head>
<body>

<header>
  <h1>Wallet</h1>
  <div class="header-right">
    <span class="version">v{{VERSION}}</span>
  </div>
</header>

<main>
  <div class="wallet-bar" id="wallet-bar">
    <div class="bar-left">
      <span id="wallet-status" class="no-wallet">Checking wallet...</span>
    </div>
    <div class="bar-right" id="wallet-actions"></div>
  </div>

  <div id="endpoints-container">
    <div class="empty-state status-checking">
      <span class="status-dot"></span>
      <span class="status-text">Loading endpoints...</span>
    </div>
  </div>
</main>

<!-- Setup Wallet Modal -->
<div class="modal-overlay" id="setup-modal">
  <div class="modal">
    <h3>Setup Wallet</h3>
    <p>Protect your private keys with encryption. Keys are stored locally and never leave this device.</p>
    <div class="setup-choices">
      <div class="setup-choice" onclick="setupBiometric()">
        <span class="choice-icon">&#128275;</span>
        <div class="choice-text">
          <h4>Biometric (Recommended)</h4>
          <p>Face ID / Touch ID encrypts your keys</p>
        </div>
      </div>
      <div class="setup-choice" onclick="showPasswordSetup()">
        <span class="choice-icon">&#128273;</span>
        <div class="choice-text">
          <h4>Password</h4>
          <p>Enter a password to encrypt your keys</p>
        </div>
      </div>
    </div>
    <div class="modal-error" id="setup-error"></div>
    <div class="modal-footer">
      <button class="btn" onclick="hideModal('setup-modal')">Cancel</button>
    </div>
  </div>
</div>

<!-- Password Setup Modal -->
<div class="modal-overlay" id="password-setup-modal">
  <div class="modal">
    <h3>Set Encryption Password</h3>
    <p>This password will be required to unlock your wallet. Choose something strong — there is no recovery if you forget it.</p>
    <label for="setup-password">Password</label>
    <input type="password" id="setup-password" placeholder="Enter password" autocomplete="off">
    <label for="setup-password-confirm">Confirm Password</label>
    <input type="password" id="setup-password-confirm" placeholder="Confirm password" autocomplete="off">
    <div class="modal-error" id="password-setup-error"></div>
    <div class="modal-footer">
      <button class="btn" onclick="hideModal('password-setup-modal'); showModal('setup-modal')">Back</button>
      <button class="btn btn-primary" id="btn-password-setup" onclick="setupWithPassword()">Set Password</button>
    </div>
  </div>
</div>

<!-- Password Unlock Modal -->
<div class="modal-overlay" id="password-unlock-modal">
  <div class="modal">
    <h3>Unlock Wallet</h3>
    <label for="unlock-password">Password</label>
    <input type="password" id="unlock-password" placeholder="Enter password" autocomplete="off">
    <div class="modal-error" id="password-unlock-error"></div>
    <div class="modal-footer">
      <button class="btn" onclick="hideModal('password-unlock-modal')">Cancel</button>
      <button class="btn btn-primary" id="btn-password-unlock" onclick="unlockWithPassword()">Unlock</button>
    </div>
  </div>
</div>

<!-- Add Key Modal (choose generate or import) -->
<div class="modal-overlay" id="addkey-modal">
  <div class="modal">
    <h3>Add Key</h3>
    <p>Generate a new key or import an existing one.</p>
    <div class="setup-choices">
      <div class="setup-choice" onclick="generateKey()">
        <span class="choice-icon">&#9889;</span>
        <div class="choice-text">
          <h4>Generate New Key</h4>
          <p>Create a random private key on this device</p>
        </div>
      </div>
      <div class="setup-choice" onclick="hideModal('addkey-modal'); showImportModal()">
        <span class="choice-icon">&#128229;</span>
        <div class="choice-text">
          <h4>Import Existing Key</h4>
          <p>Paste a private key you already have</p>
        </div>
      </div>
    </div>
    <div class="modal-error" id="addkey-error"></div>
    <div class="modal-footer">
      <button class="btn" onclick="hideModal('addkey-modal')">Cancel</button>
    </div>
  </div>
</div>

<!-- Import Key Modal -->
<div class="modal-overlay" id="import-modal">
  <div class="modal">
    <h3>Import Private Key</h3>
    <label for="import-label">Label</label>
    <input type="text" id="import-label" placeholder="e.g. Main, Test, Hot" autocomplete="off" spellcheck="false">
    <label for="import-key">Private Key (hex)</label>
    <input type="password" id="import-key" placeholder="0x..." autocomplete="off" spellcheck="false">
    <div class="modal-error" id="import-error"></div>
    <div class="modal-footer">
      <button class="btn" onclick="hideModal('import-modal')">Cancel</button>
      <button class="btn btn-primary" id="btn-import-confirm" onclick="doImportKey()">Import</button>
    </div>
  </div>
</div>

<script>
// ── State ──────────────────────────────────────────────
let endpoints = [];
let walletState = 'none';       // 'none' | 'locked' | 'unlocked'
let decryptedKeys = [];          // [{id, label, address, key}] — in-memory only
let activeKeyIndex = 0;
let aesKey = null;               // CryptoKey, held while unlocked
let storedKeyCount = 0;
let credMethod = '';             // 'prf' | 'password'

// ── Constants ──────────────────────────────────────────
const PRF_SALT = new TextEncoder().encode('wallet-encryption-v1');
const HKDF_INFO = new TextEncoder().encode('AES-GCM Wallet Encryption Key V1');
const PBKDF2_ITERATIONS = 600000;
const DB_NAME = 'wallet-vault';
const DB_VERSION = 1;

// ── Init ───────────────────────────────────────────────
(async function init() {
  try {
    const cred = await getCredential();
    if (cred) {
      const keys = await getEncryptedKeys();
      storedKeyCount = keys.length;
      credMethod = cred.method || 'prf';
      walletState = 'locked';
    }
  } catch (e) {
    console.error('init check failed:', e);
  }
  renderWalletBar();
  refresh();
  setInterval(refresh, 10000);
})();

// ── IndexedDB Helpers ──────────────────────────────────
function openVaultDB() {
  return new Promise((resolve, reject) => {
    const req = indexedDB.open(DB_NAME, DB_VERSION);
    req.onupgradeneeded = (e) => {
      const db = e.target.result;
      if (!db.objectStoreNames.contains('credentials')) {
        db.createObjectStore('credentials', { keyPath: 'id' });
      }
      if (!db.objectStoreNames.contains('keys')) {
        db.createObjectStore('keys', { keyPath: 'id', autoIncrement: true });
      }
    };
    req.onsuccess = () => resolve(req.result);
    req.onerror = () => reject(req.error);
  });
}

async function saveCredential(cred) {
  const db = await openVaultDB();
  return new Promise((resolve, reject) => {
    const tx = db.transaction('credentials', 'readwrite');
    tx.objectStore('credentials').put(cred);
    tx.oncomplete = () => resolve();
    tx.onerror = () => reject(tx.error);
  });
}

async function getCredential() {
  const db = await openVaultDB();
  return new Promise((resolve, reject) => {
    const tx = db.transaction('credentials', 'readonly');
    const req = tx.objectStore('credentials').getAll();
    req.onsuccess = () => resolve(req.result.length > 0 ? req.result[0] : null);
    req.onerror = () => reject(req.error);
  });
}

async function saveEncryptedKey(record) {
  const db = await openVaultDB();
  return new Promise((resolve, reject) => {
    const tx = db.transaction('keys', 'readwrite');
    tx.objectStore('keys').put(record);
    tx.oncomplete = () => resolve();
    tx.onerror = () => reject(tx.error);
  });
}

async function getEncryptedKeys() {
  const db = await openVaultDB();
  return new Promise((resolve, reject) => {
    const tx = db.transaction('keys', 'readonly');
    const req = tx.objectStore('keys').getAll();
    req.onsuccess = () => resolve(req.result);
    req.onerror = () => reject(req.error);
  });
}

async function deleteEncryptedKey(id) {
  const db = await openVaultDB();
  return new Promise((resolve, reject) => {
    const tx = db.transaction('keys', 'readwrite');
    tx.objectStore('keys').delete(id);
    tx.oncomplete = () => resolve();
    tx.onerror = () => reject(tx.error);
  });
}

// ── Crypto Helpers ─────────────────────────────────────
async function deriveAESKeyFromPRF(prfOutput) {
  const keyMaterial = await crypto.subtle.importKey(
    'raw', prfOutput, 'HKDF', false, ['deriveKey']
  );
  return crypto.subtle.deriveKey(
    { name: 'HKDF', salt: PRF_SALT, info: HKDF_INFO, hash: 'SHA-256' },
    keyMaterial,
    { name: 'AES-GCM', length: 256 },
    false,
    ['encrypt', 'decrypt']
  );
}

async function deriveAESKeyFromPassword(password, salt) {
  const enc = new TextEncoder();
  const keyMaterial = await crypto.subtle.importKey(
    'raw', enc.encode(password), 'PBKDF2', false, ['deriveKey']
  );
  return crypto.subtle.deriveKey(
    { name: 'PBKDF2', salt: salt, iterations: PBKDF2_ITERATIONS, hash: 'SHA-256' },
    keyMaterial,
    { name: 'AES-GCM', length: 256 },
    false,
    ['encrypt', 'decrypt']
  );
}

async function encryptPrivateKey(plaintext, key) {
  const iv = crypto.getRandomValues(new Uint8Array(12));
  const encoded = new TextEncoder().encode(plaintext);
  const encrypted = await crypto.subtle.encrypt(
    { name: 'AES-GCM', iv }, key, encoded
  );
  return { encrypted: new Uint8Array(encrypted), iv };
}

async function decryptPrivateKey(encrypted, iv, key) {
  const decrypted = await crypto.subtle.decrypt(
    { name: 'AES-GCM', iv }, key, encrypted
  );
  return new TextDecoder().decode(decrypted);
}

// ── Biometric Setup ────────────────────────────────────
async function setupBiometric() {
  const errEl = document.getElementById('setup-error');
  errEl.style.display = 'none';

  if (!window.PublicKeyCredential) {
    errEl.textContent = 'WebAuthn is not available in this browser. Use password setup instead.';
    errEl.style.display = 'block';
    return;
  }

  try {
    // 1. Create credential with PRF extension.
    const userId = crypto.getRandomValues(new Uint8Array(32));
    const credential = await navigator.credentials.create({
      publicKey: {
        rp: { name: 'Wallet', id: location.hostname },
        user: {
          id: userId,
          name: 'wallet-user',
          displayName: 'Wallet User'
        },
        challenge: crypto.getRandomValues(new Uint8Array(32)),
        pubKeyCredParams: [
          { type: 'public-key', alg: -7 },
          { type: 'public-key', alg: -257 }
        ],
        authenticatorSelection: {
          residentKey: 'preferred',
          userVerification: 'required'
        },
        extensions: { prf: {} }
      }
    });

    // 2. Try PRF eval regardless of prf.enabled — some browsers
    //    (Safari) report enabled:false but support PRF at assertion time.
    const transports = credential.response.getTransports ? credential.response.getTransports() : [];

    const assertion = await navigator.credentials.get({
      publicKey: {
        challenge: crypto.getRandomValues(new Uint8Array(32)),
        rpId: location.hostname,
        allowCredentials: [{
          type: 'public-key',
          id: credential.rawId,
          transports: transports
        }],
        userVerification: 'required',
        extensions: {
          prf: { eval: { first: PRF_SALT } }
        }
      }
    });

    const exts = assertion.getClientExtensionResults();
    if (!exts.prf || !exts.prf.results || !exts.prf.results.first) {
      errEl.textContent = 'Your authenticator does not support PRF encryption. Use password setup instead.';
      errEl.style.display = 'block';
      return;
    }

    // 3. PRF works — derive AES key and store credential.
    aesKey = await deriveAESKeyFromPRF(exts.prf.results.first);

    await saveCredential({
      id: 'primary',
      method: 'prf',
      credentialId: Array.from(new Uint8Array(credential.rawId)),
      rpId: location.hostname,
      transports: transports,
      createdAt: Date.now()
    });

    credMethod = 'prf';
    walletState = 'unlocked';
    decryptedKeys = [];
    storedKeyCount = 0;
    renderWalletBar();
    hideModal('setup-modal');
    showModal('addkey-modal');

  } catch (err) {
    if (err.name === 'NotAllowedError') {
      errEl.textContent = 'Biometric prompt was cancelled or timed out.';
    } else {
      errEl.textContent = err.message;
    }
    errEl.style.display = 'block';
  }
}

// ── Password Setup ─────────────────────────────────────
function showPasswordSetup() {
  hideModal('setup-modal');
  document.getElementById('setup-password').value = '';
  document.getElementById('setup-password-confirm').value = '';
  document.getElementById('password-setup-error').style.display = 'none';
  showModal('password-setup-modal');
}

async function setupWithPassword() {
  const pw = document.getElementById('setup-password').value;
  const confirm = document.getElementById('setup-password-confirm').value;
  const errEl = document.getElementById('password-setup-error');
  const btn = document.getElementById('btn-password-setup');
  errEl.style.display = 'none';

  if (!pw) {
    errEl.textContent = 'Please enter a password.';
    errEl.style.display = 'block';
    return;
  }
  if (pw.length < 8) {
    errEl.textContent = 'Password must be at least 8 characters.';
    errEl.style.display = 'block';
    return;
  }
  if (pw !== confirm) {
    errEl.textContent = 'Passwords do not match.';
    errEl.style.display = 'block';
    return;
  }

  btn.disabled = true;
  btn.textContent = 'Deriving key...';

  try {
    const salt = crypto.getRandomValues(new Uint8Array(32));
    aesKey = await deriveAESKeyFromPassword(pw, salt);

    await saveCredential({
      id: 'primary',
      method: 'password',
      pbkdf2Salt: Array.from(salt),
      createdAt: Date.now()
    });

    credMethod = 'password';
    walletState = 'unlocked';
    decryptedKeys = [];
    storedKeyCount = 0;
    renderWalletBar();
    hideModal('password-setup-modal');
    showModal('addkey-modal');

  } catch (err) {
    errEl.textContent = 'Setup failed: ' + err.message;
    errEl.style.display = 'block';
  } finally {
    btn.disabled = false;
    btn.textContent = 'Set Password';
  }
}

// ── Unlock ─────────────────────────────────────────────
async function unlockWallet() {
  if (credMethod === 'password') {
    document.getElementById('unlock-password').value = '';
    document.getElementById('password-unlock-error').style.display = 'none';
    showModal('password-unlock-modal');
    return;
  }
  // PRF unlock.
  const btn = document.querySelector('#wallet-actions .btn-primary');
  if (btn) { btn.disabled = true; btn.textContent = 'Unlocking...'; }

  try {
    const stored = await getCredential();
    if (!stored) throw new Error('No credential found.');

    const credentialId = new Uint8Array(stored.credentialId);
    const assertion = await navigator.credentials.get({
      publicKey: {
        challenge: crypto.getRandomValues(new Uint8Array(32)),
        rpId: stored.rpId,
        allowCredentials: [{
          type: 'public-key',
          id: credentialId.buffer,
          transports: stored.transports || []
        }],
        userVerification: 'required',
        extensions: {
          prf: { eval: { first: PRF_SALT } }
        }
      }
    });

    const exts = assertion.getClientExtensionResults();
    if (!exts.prf || !exts.prf.results || !exts.prf.results.first) {
      throw new Error('PRF evaluation failed.');
    }

    aesKey = await deriveAESKeyFromPRF(exts.prf.results.first);
    await decryptAllKeys();
    walletState = 'unlocked';
    renderWalletBar();
    refresh();

  } catch (err) {
    if (err.name !== 'NotAllowedError') {
      console.error('Unlock failed:', err);
    }
    renderWalletBar();
  }
}

async function unlockWithPassword() {
  const pw = document.getElementById('unlock-password').value;
  const errEl = document.getElementById('password-unlock-error');
  const btn = document.getElementById('btn-password-unlock');
  errEl.style.display = 'none';

  if (!pw) {
    errEl.textContent = 'Please enter your password.';
    errEl.style.display = 'block';
    return;
  }

  btn.disabled = true;
  btn.textContent = 'Unlocking...';

  try {
    const stored = await getCredential();
    if (!stored || !stored.pbkdf2Salt) throw new Error('No password credential found.');

    const salt = new Uint8Array(stored.pbkdf2Salt);
    aesKey = await deriveAESKeyFromPassword(pw, salt);

    // Try decrypting — if the password is wrong, decryption will fail.
    await decryptAllKeys();
    walletState = 'unlocked';
    renderWalletBar();
    hideModal('password-unlock-modal');
    refresh();

  } catch (err) {
    errEl.textContent = 'Wrong password or decryption failed.';
    errEl.style.display = 'block';
    aesKey = null;
  } finally {
    btn.disabled = false;
    btn.textContent = 'Unlock';
  }
}

async function decryptAllKeys() {
  const encryptedKeys = await getEncryptedKeys();
  decryptedKeys = [];
  for (const rec of encryptedKeys) {
    const plaintext = await decryptPrivateKey(
      new Uint8Array(rec.encrypted),
      new Uint8Array(rec.iv),
      aesKey
    );
    decryptedKeys.push({ id: rec.id, label: rec.label, address: rec.address, key: plaintext });
  }
  activeKeyIndex = 0;
  storedKeyCount = decryptedKeys.length;
}

// ── Lock ───────────────────────────────────────────────
function lockWallet() {
  for (let i = 0; i < decryptedKeys.length; i++) {
    decryptedKeys[i].key = '';
  }
  decryptedKeys = [];
  aesKey = null;
  activeKeyIndex = 0;
  walletState = 'locked';
  renderWalletBar();
  renderEndpoints();
}

// ── Import Key ─────────────────────────────────────────
async function doImportKey() {
  const labelInput = document.getElementById('import-label');
  const keyInput = document.getElementById('import-key');
  const errEl = document.getElementById('import-error');
  const btn = document.getElementById('btn-import-confirm');
  errEl.style.display = 'none';

  const label = labelInput.value.trim() || 'Key ' + (storedKeyCount + 1);
  let key = keyInput.value.trim();

  if (!key) {
    errEl.textContent = 'Please enter a private key.';
    errEl.style.display = 'block';
    return;
  }

  if (!key.startsWith('0x')) key = '0x' + key;

  if (!/^0x[0-9a-fA-F]{64}$/.test(key)) {
    errEl.textContent = 'Invalid key format. Expected 64 hex characters.';
    errEl.style.display = 'block';
    return;
  }

  if (!aesKey) {
    errEl.textContent = 'Wallet is not unlocked. Please unlock first.';
    errEl.style.display = 'block';
    return;
  }

  btn.disabled = true;
  btn.textContent = 'Encrypting...';

  try {
    await ensureEthers();
    const wallet = new ethers.Wallet(key);
    const address = wallet.address;

    const { encrypted, iv } = await encryptPrivateKey(key, aesKey);

    await saveEncryptedKey({
      label: label,
      address: address,
      encrypted: Array.from(encrypted),
      iv: Array.from(iv),
      createdAt: Date.now()
    });

    const allKeys = await getEncryptedKeys();
    const newest = allKeys[allKeys.length - 1];
    decryptedKeys.push({ id: newest.id, label: label, address: address, key: key });
    activeKeyIndex = decryptedKeys.length - 1;
    storedKeyCount = decryptedKeys.length;

    labelInput.value = '';
    keyInput.value = '';
    errEl.style.display = 'none';
    hideModal('import-modal');
    renderWalletBar();
    refresh();

  } catch (err) {
    errEl.textContent = 'Failed: ' + err.message;
    errEl.style.display = 'block';
  } finally {
    btn.disabled = false;
    btn.textContent = 'Import';
  }
}

// ── Generate Key ───────────────────────────────────────
async function generateKey() {
  const errEl = document.getElementById('addkey-error');
  errEl.style.display = 'none';

  if (!aesKey) {
    errEl.textContent = 'Wallet is not unlocked.';
    errEl.style.display = 'block';
    return;
  }

  try {
    await ensureEthers();
    const wallet = ethers.Wallet.createRandom();
    const key = wallet.privateKey;
    const address = wallet.address;
    const label = 'Key ' + (storedKeyCount + 1);

    const { encrypted, iv } = await encryptPrivateKey(key, aesKey);

    await saveEncryptedKey({
      label: label,
      address: address,
      encrypted: Array.from(encrypted),
      iv: Array.from(iv),
      createdAt: Date.now()
    });

    const allKeys = await getEncryptedKeys();
    const newest = allKeys[allKeys.length - 1];
    decryptedKeys.push({ id: newest.id, label: label, address: address, key: key });
    activeKeyIndex = decryptedKeys.length - 1;
    storedKeyCount = decryptedKeys.length;

    hideModal('addkey-modal');
    renderWalletBar();
    refresh();

  } catch (err) {
    errEl.textContent = 'Failed: ' + err.message;
    errEl.style.display = 'block';
  }
}

function showAddKeyModal() {
  document.getElementById('addkey-error').style.display = 'none';
  showModal('addkey-modal');
}

// ── Wallet Bar Rendering ───────────────────────────────
function renderWalletBar() {
  const statusEl = document.getElementById('wallet-status');
  const actionsEl = document.getElementById('wallet-actions');

  if (walletState === 'none') {
    statusEl.className = 'no-wallet';
    statusEl.textContent = 'No wallet configured';
    actionsEl.innerHTML = '<button class="btn btn-primary" onclick="showModal(\'setup-modal\')">Setup Wallet</button>';
  } else if (walletState === 'locked') {
    const methodLabel = credMethod === 'prf' ? 'biometric' : 'password';
    statusEl.className = 'label';
    statusEl.innerHTML = '<span class="lock-icon">&#128274;</span> Wallet locked' +
      (storedKeyCount > 0 ? ' <span class="key-badge">' + storedKeyCount + ' key' + (storedKeyCount !== 1 ? 's' : '') + '</span>' : '') +
      ' <span class="method-badge">' + methodLabel + '</span>';
    actionsEl.innerHTML = '<button class="btn btn-primary" onclick="unlockWallet()">Unlock</button>';
  } else if (walletState === 'unlocked') {
    let html = '';
    if (decryptedKeys.length > 0) {
      const active = decryptedKeys[activeKeyIndex];
      if (decryptedKeys.length > 1) {
        html += '<select class="key-selector" onchange="switchKey(this.value)">';
        for (let i = 0; i < decryptedKeys.length; i++) {
          const k = decryptedKeys[i];
          const sel = i === activeKeyIndex ? ' selected' : '';
          html += '<option value="' + i + '"' + sel + '>' + esc(k.label) + ' (' + k.address.slice(0, 6) + '...' + k.address.slice(-4) + ')</option>';
        }
        html += '</select>';
      }
      statusEl.className = 'address';
      statusEl.textContent = active.address;
    } else {
      statusEl.className = 'no-wallet';
      statusEl.textContent = 'No keys imported';
    }

    actionsEl.innerHTML = html +
      '<button class="btn btn-primary" onclick="showAddKeyModal()">Add Key</button>' +
      '<button class="btn" onclick="lockWallet()">Lock</button>';
  }
}

function switchKey(index) {
  activeKeyIndex = parseInt(index, 10);
  renderWalletBar();
  refresh();
}

function showImportModal() {
  document.getElementById('import-label').value = '';
  document.getElementById('import-key').value = '';
  document.getElementById('import-error').style.display = 'none';
  showModal('import-modal');
}

// ── Refresh ────────────────────────────────────────────
async function refresh() {
  try {
    const resp = await fetch('/api/status');
    const data = await resp.json();
    endpoints = data.endpoints || [];
    renderEndpoints();
  } catch (err) {
    console.error('status poll failed:', err);
  }
}

// ── Render ─────────────────────────────────────────────
function renderEndpoints() {
  const container = document.getElementById('endpoints-container');
  const walletAddress = getActiveAddress();

  if (endpoints.length === 0) {
    container.innerHTML =
      '<div class="empty-state">' +
        '<h2>No Endpoints Configured</h2>' +
        '<p>Add RPC endpoints to endpoints.json to get started.</p>' +
      '</div>';
    return;
  }

  let html = '<div class="endpoints">';
  for (const ep of endpoints) {
    const statusClass = ep.online ? 'status-online' : 'status-offline';
    const statusLabel = ep.online ? 'Online' : 'Offline';
    const chainId = ep.chain_id ? hexToDecimal(ep.chain_id) : '\u2014';
    const blockNum = ep.block_number ? hexToDecimal(ep.block_number) : '\u2014';
    const latencyClass = ep.latency_ms < 200 ? 'fast' : ep.latency_ms < 1000 ? 'medium' : 'slow';
    const urlAbbrev = abbreviateURL(ep.url);

    html += '<div class="ep-card">';
    html +=   '<div class="ep-card-header">';
    html +=     '<h3>' + esc(ep.name) + '</h3>';
    html +=     '<span class="' + statusClass + '">';
    html +=       '<span class="status-dot"></span>';
    html +=       '<span class="status-text">' + statusLabel + '</span>';
    html +=     '</span>';
    html +=   '</div>';
    html +=   '<div class="ep-card-body">';
    html +=     '<div class="ep-row">';
    html +=       '<span class="label">RPC</span>';
    html +=       '<span class="url-display" title="' + esc(ep.url) + '">' + esc(urlAbbrev) + '</span>';
    html +=     '</div>';
    html +=     '<div class="ep-row">';
    html +=       '<span class="label">Chain ID</span>';
    html +=       '<span class="value">' + chainId + '</span>';
    html +=     '</div>';
    html +=     '<div class="ep-row">';
    html +=       '<span class="label">Block</span>';
    html +=       '<span class="value">' + formatNumber(blockNum) + '</span>';
    html +=     '</div>';
    html +=     '<div class="ep-row">';
    html +=       '<span class="label">Latency</span>';
    html +=       '<span class="latency ' + latencyClass + '">' + ep.latency_ms + ' ms</span>';
    html +=     '</div>';

    if (walletAddress && ep.online) {
      html +=   '<div class="ep-row" id="balance-' + esc(ep.id) + '">';
      html +=     '<span class="label">Balance</span>';
      html +=     '<span class="value balance" data-ep="' + esc(ep.id) + '">...</span>';
      html +=   '</div>';
    }

    html +=   '</div>';
    html += '</div>';
  }
  html += '</div>';
  container.innerHTML = html;

  if (walletAddress) {
    fetchBalances(walletAddress);
  }
}

// ── Balances ───────────────────────────────────────────
async function fetchBalances(address) {
  for (const ep of endpoints) {
    if (!ep.online) continue;
    try {
      const resp = await fetch('/api/rpc/' + ep.id, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ method: 'eth_getBalance', params: [address, 'latest'] })
      });
      const data = await resp.json();
      if (data.result) {
        const el = document.querySelector('[data-ep="' + ep.id + '"]');
        if (el) {
          el.textContent = formatBalance(data.result) + ' ' + (ep.symbol || 'ETH');
        }
      }
    } catch (err) {
      console.error('balance fetch failed for ' + ep.name + ':', err);
    }
  }
}

// ── Active Address Helper ──────────────────────────────
function getActiveAddress() {
  if (walletState !== 'unlocked' || decryptedKeys.length === 0) return '';
  return decryptedKeys[activeKeyIndex].address;
}

// ── Ethers.js Lazy Load ────────────────────────────────
let ethersLoaded = false;
function ensureEthers() {
  if (ethersLoaded) return Promise.resolve();
  return new Promise((resolve, reject) => {
    const script = document.createElement('script');
    script.src = 'https://cdnjs.cloudflare.com/ajax/libs/ethers/6.13.4/ethers.umd.min.js';
    script.onload = () => { ethersLoaded = true; resolve(); };
    script.onerror = () => reject(new Error('Failed to load ethers.js'));
    document.head.appendChild(script);
  });
}

// ── Helpers ────────────────────────────────────────────
function hexToDecimal(hex) {
  if (!hex || hex === '0x') return '0';
  return parseInt(hex, 16).toString();
}

function formatNumber(n) {
  if (n === '\u2014') return n;
  return Number(n).toLocaleString();
}

function formatBalance(hexWei) {
  const wei = BigInt(hexWei);
  const ether = Number(wei) / 1e18;
  if (ether === 0) return '0';
  if (ether < 0.0001) return '< 0.0001';
  return ether.toFixed(4);
}

function abbreviateURL(url) {
  try {
    const u = new URL(url);
    let display = u.hostname;
    if (u.port) display += ':' + u.port;
    if (u.pathname !== '/') display += u.pathname;
    return display;
  } catch {
    return url;
  }
}

function esc(s) {
  const d = document.createElement('div');
  d.textContent = s || '';
  return d.innerHTML;
}

function showModal(id) {
  document.getElementById(id).classList.add('active');
}
function hideModal(id) {
  document.getElementById(id).classList.remove('active');
}

// Close modals on overlay click.
document.querySelectorAll('.modal-overlay').forEach(overlay => {
  overlay.addEventListener('click', (e) => {
    if (e.target === overlay) overlay.classList.remove('active');
  });
});

// Close modals on Escape key.
document.addEventListener('keydown', (e) => {
  if (e.key === 'Escape') {
    document.querySelectorAll('.modal-overlay.active').forEach(m => m.classList.remove('active'));
  }
});

// Submit password modals on Enter key.
document.addEventListener('keydown', (e) => {
  if (e.key !== 'Enter') return;
  if (document.getElementById('password-setup-modal').classList.contains('active')) {
    setupWithPassword();
  } else if (document.getElementById('password-unlock-modal').classList.contains('active')) {
    unlockWithPassword();
  } else if (document.getElementById('import-modal').classList.contains('active')) {
    doImportKey();
  }
});
</script>
</body>
</html>`
