/**
 * QR Code Generator — Frontend Application
 * Handles tab switching, form collection, API calls, preview, and download.
 */

(function () {
  'use strict';

  // ═══════════════════════════════════════════
  // DOM References
  // ═══════════════════════════════════════════
  const $ = (sel) => document.querySelector(sel);
  const $$ = (sel) => document.querySelectorAll(sel);

  const tabs = $$('.tab');
  const tabContents = $$('.tab-content');
  const btnGenerate = $('#btn-generate');
  const btnDownload = $('#btn-download');
  const previewArea = $('#preview-area');
  const previewImage = $('#preview-image');
  const previewPlaceholder = $('#preview-placeholder');
  const infoBadges = $('#info-badges');
  const badgeType = $('#badge-type');
  const badgeSize = $('#badge-size');
  const badgeEC = $('#badge-ec');
  const sizeSlider = $('#qr-size');
  const sizeValue = $('#size-value');
  const fgColor = $('#fg-color');
  const bgColor = $('#bg-color');
  const fgHex = $('#fg-hex');
  const bgHex = $('#bg-hex');
  const toastContainer = $('#toast-container');

  // ═══════════════════════════════════════════
  // State
  // ═══════════════════════════════════════════
  let activeTab = 'text';
  let currentBlob = null;
  let debounceTimer = null;

  // ═══════════════════════════════════════════
  // Tab Switching
  // ═══════════════════════════════════════════
  tabs.forEach((tab) => {
    tab.addEventListener('click', () => {
      const target = tab.dataset.tab;
      if (target === activeTab) return;

      // Update tab buttons
      tabs.forEach((t) => {
        t.classList.remove('tab--active');
        t.setAttribute('aria-selected', 'false');
      });
      tab.classList.add('tab--active');
      tab.setAttribute('aria-selected', 'true');

      // Update tab content
      tabContents.forEach((c) => c.classList.remove('tab-content--active'));
      $(`#content-${target}`).classList.add('tab-content--active');

      activeTab = target;
    });
  });

  // ═══════════════════════════════════════════
  // Color Picker Sync
  // ═══════════════════════════════════════════
  function syncColor(colorInput, hexInput) {
    colorInput.addEventListener('input', () => {
      hexInput.value = colorInput.value;
    });
    hexInput.addEventListener('input', () => {
      const val = hexInput.value;
      if (/^#[0-9A-Fa-f]{6}$/.test(val)) {
        colorInput.value = val;
      }
    });
  }

  syncColor(fgColor, fgHex);
  syncColor(bgColor, bgHex);

  // ═══════════════════════════════════════════
  // Size Slider
  // ═══════════════════════════════════════════
  sizeSlider.addEventListener('input', () => {
    sizeValue.textContent = `${sizeSlider.value}px`;
  });

  // ═══════════════════════════════════════════
  // Content Collection
  // ═══════════════════════════════════════════
  function getContent() {
    switch (activeTab) {
      case 'text':
        return $('#input-text').value.trim();
      case 'url':
        return $('#input-url').value.trim();
      case 'email':
        return $('#input-email').value.trim();
      case 'phone':
        return $('#input-phone').value.trim();
      case 'wifi':
        return JSON.stringify({
          ssid: $('#wifi-ssid').value.trim(),
          password: $('#wifi-password').value.trim(),
          encryption: $('#wifi-encryption').value,
          hidden: $('#wifi-hidden').checked,
        });
      default:
        return '';
    }
  }

  function validateContent(content) {
    if (!content || content.length === 0) {
      return 'Please enter some content to generate a QR code.';
    }

    if (activeTab === 'email' && !content.includes('@')) {
      return 'Please enter a valid email address.';
    }

    if (activeTab === 'wifi') {
      try {
        const wifi = JSON.parse(content);
        if (!wifi.ssid) return 'WiFi network name (SSID) is required.';
      } catch {
        return 'Invalid WiFi data.';
      }
    }

    return null;
  }

  // ═══════════════════════════════════════════
  // Error Correction Label
  // ═══════════════════════════════════════════
  function getECLabel(value) {
    const labels = { L: 'Low', M: 'Medium', Q: 'Quartile', H: 'High' };
    return labels[value] || 'Medium';
  }

  // ═══════════════════════════════════════════
  // Content Type Label
  // ═══════════════════════════════════════════
  function getTypeLabel(type) {
    const icons = { text: '📝', url: '🔗', email: '📧', phone: '📞', wifi: '📶' };
    const names = { text: 'Text', url: 'URL', email: 'Email', phone: 'Phone', wifi: 'WiFi' };
    return `${icons[type] || '📝'} ${names[type] || 'Text'}`;
  }

  // ═══════════════════════════════════════════
  // Toast Notifications
  // ═══════════════════════════════════════════
  function showToast(message, type = 'error') {
    const toast = document.createElement('div');
    toast.className = `toast${type === 'success' ? ' toast--success' : ''}`;
    toast.innerHTML = `
      <span>${type === 'error' ? '❌' : '✅'} ${message}</span>
      <button class="toast__close" aria-label="Close">&times;</button>
    `;

    // Remove existing toasts
    toastContainer.querySelectorAll('.toast').forEach((t) => t.remove());
    toastContainer.appendChild(toast);

    // Close button
    toast.querySelector('.toast__close').addEventListener('click', () => {
      toast.style.animation = 'toastSlideIn 0.3s ease reverse';
      setTimeout(() => toast.remove(), 300);
    });

    // Auto-dismiss after 5s
    setTimeout(() => {
      if (toast.parentNode) {
        toast.style.animation = 'toastSlideIn 0.3s ease reverse';
        setTimeout(() => toast.remove(), 300);
      }
    }, 5000);
  }

  // ═══════════════════════════════════════════
  // Generate QR Code
  // ═══════════════════════════════════════════
  async function generateQR() {
    const content = getContent();
    const error = validateContent(content);

    if (error) {
      showToast(error);
      return;
    }

    // Start loading state
    btnGenerate.classList.add('btn-generate--loading');
    btnGenerate.disabled = true;

    const payload = {
      content: content,
      content_type: activeTab,
      size: parseInt(sizeSlider.value, 10),
      error_correction: $('#error-correction').value,
      foreground_color: fgHex.value,
      background_color: bgHex.value,
    };

    try {
      const response = await fetch('/api/generate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        const err = await response.json().catch(() => ({ message: 'Unknown error' }));
        throw new Error(err.message || `Server error (${response.status})`);
      }

      // Get the image blob
      const blob = await response.blob();
      currentBlob = blob;

      // Display the image
      const url = URL.createObjectURL(blob);
      previewImage.src = url;
      previewImage.style.display = 'block';
      previewPlaceholder.style.display = 'none';
      previewArea.classList.add('preview-area--has-image');

      // Update badges
      badgeType.textContent = getTypeLabel(activeTab);
      badgeSize.textContent = `📐 ${sizeSlider.value}px`;
      badgeEC.textContent = `🛡️ ${getECLabel($('#error-correction').value)}`;
      infoBadges.style.display = 'flex';

      // Enable download
      btnDownload.disabled = false;

      showToast('QR code generated successfully!', 'success');
    } catch (err) {
      showToast(err.message || 'Failed to generate QR code. Please try again.');
    } finally {
      btnGenerate.classList.remove('btn-generate--loading');
      btnGenerate.disabled = false;
    }
  }

  // ═══════════════════════════════════════════
  // Download QR Code
  // ═══════════════════════════════════════════
  function downloadQR() {
    if (!currentBlob) return;

    // Re-create blob with explicit PNG MIME type to ensure correct download
    const pngBlob = new Blob([currentBlob], { type: 'image/png' });
    const url = URL.createObjectURL(pngBlob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `qrcode-${activeTab}-${Date.now()}.png`;
    a.style.display = 'none';
    document.body.appendChild(a);

    // Use setTimeout to ensure the link is in the DOM before clicking
    setTimeout(() => {
      a.click();
      document.body.removeChild(a);
      // Revoke after a longer delay to allow download to start
      setTimeout(() => URL.revokeObjectURL(url), 1000);
      showToast('QR code downloaded!', 'success');
    }, 0);
  }

  // ═══════════════════════════════════════════
  // Event Listeners
  // ═══════════════════════════════════════════
  btnGenerate.addEventListener('click', generateQR);
  btnDownload.addEventListener('click', downloadQR);

  // Keyboard: Enter to generate
  document.addEventListener('keydown', (e) => {
    if (e.key === 'Enter' && !e.shiftKey && e.target.tagName !== 'TEXTAREA') {
      e.preventDefault();
      generateQR();
    }
  });

  // ═══════════════════════════════════════════
  // Auto-fill demo on first load
  // ═══════════════════════════════════════════
  (function setDefaults() {
    $('#input-text').value = 'Hello, World! 👋\nScan this QR code.';
    $('#input-url').value = 'https://github.com';
    $('#input-email').value = 'hello@example.com';
    $('#input-phone').value = '+1 234 567 8900';
    $('#wifi-ssid').value = 'MyNetwork';
    $('#wifi-password').value = '';
  })();
})();
