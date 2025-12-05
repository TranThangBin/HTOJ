// Client-side i18n manager with localStorage cache
class I18n {
  constructor() {
    this.translations = {};
    this.currentLang = this.getLangFromCookie() || 'en';
    this.cacheVersion = '1.0'; // Increment this to invalidate cache
    this.init();
  }

  getLangFromCookie() {
    const match = document.cookie.match(/lang=([^;]+)/);
    return match ? match[1] : null;
  }

  setLangCookie(lang) {
    const days = 30;
    const date = new Date();
    date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
    document.cookie = `lang=${lang};expires=${date.toUTCString()};path=/`;
  }

  getCacheKey(lang) {
    return `i18n_${lang}_v${this.cacheVersion}`;
  }

  async init() {
    await this.loadTranslations(this.currentLang);
  }

  async loadTranslations(lang) {
    // Try to get from localStorage cache first
    const cacheKey = this.getCacheKey(lang);
    const cached = localStorage.getItem(cacheKey);
    const etagKey = `${cacheKey}_etag`;
    const cachedETag = localStorage.getItem(etagKey);
    
    if (cached && cachedETag) {
      try {
        // Validate cache with server using ETag
        const response = await fetch(`/api/translations/${lang}`, {
          headers: {
            'If-None-Match': cachedETag
          }
        });
        
        if (response.status === 304) {
          // Cache is still valid
          this.translations = JSON.parse(cached);
          this.currentLang = lang;
          this.updatePageContent();
          console.log(`Cache validated, loaded ${lang} translations from localStorage`);
          return;
        }
        
        // Cache is outdated, fetch new data
        if (response.ok) {
          this.translations = await response.json();
          const newETag = response.headers.get('ETag');
          
          // Update cache
          localStorage.setItem(cacheKey, JSON.stringify(this.translations));
          localStorage.setItem(etagKey, newETag);
          this.currentLang = lang;
          this.updatePageContent();
          console.log(`Cache updated, loaded ${lang} translations from server`);
          return;
        }
      } catch (error) {
        console.warn('Failed to validate cache, using cached version:', error);
        // Use cached version if validation fails
        this.translations = JSON.parse(cached);
        this.currentLang = lang;
        this.updatePageContent();
        return;
      }
    }

    try {
      const response = await fetch(`/api/translations/${lang}`);
      if (response.ok) {
        this.translations = await response.json();
        this.currentLang = lang;
        const etag = response.headers.get('ETag');
        
        // Save to localStorage cache
        try {
          localStorage.setItem(cacheKey, JSON.stringify(this.translations));
          if (etag) {
            localStorage.setItem(etagKey, etag);
          }
          console.log(`Cached ${lang} translations to localStorage`);
        } catch (error) {
          console.warn('Failed to cache translations:', error);
          this.clearOldCache();
        }
        
        this.updatePageContent();
      }
    } catch (error) {
      console.error('Failed to load translations:', error);
    }
  }

  clearOldCache() {
    // Remove old translation cache entries
    const keys = Object.keys(localStorage);
    keys.forEach(key => {
      if (key.startsWith('i18n_') && !key.includes(`_v${this.cacheVersion}`)) {
        localStorage.removeItem(key);
        console.log(`Removed old cache: ${key}`);
      }
    });
  }

  t(key) {
    return this.translations[key] || key;
  }

  async switchLanguage(lang) {
    if (lang === this.currentLang) return;
    
    await this.loadTranslations(lang);
    this.setLangCookie(lang);
  }

  updatePageContent() {
    // Update all elements with data-i18n attribute
    document.querySelectorAll('[data-i18n]').forEach(el => {
      const key = el.getAttribute('data-i18n');
      el.textContent = this.t(key);
    });
  }
}

// Initialize i18n
const i18n = new I18n();

// Setup language switcher event listeners
document.addEventListener('DOMContentLoaded', () => {
  document.querySelectorAll('[data-lang]').forEach(link => {
    link.addEventListener('click', async (e) => {
      e.preventDefault();
      const lang = e.currentTarget.getAttribute('data-lang');
      await i18n.switchLanguage(lang);
    });
  });
});
