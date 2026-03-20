/**
 * Configuration Service
 */
export const Config = {
  STORAGE_KEY: "APP_API_BASE_URL",
  DEFAULT_BASE_URL: "http://localhost:12345",

  getUrl() {
    return sessionStorage.getItem(this.STORAGE_KEY);
  },

  setUrl(url) {
    if (!url) return;
    sessionStorage.setItem(this.STORAGE_KEY, url.replace(/\/$/, ""));
  },

  clear() {
    sessionStorage.removeItem(this.STORAGE_KEY);
  },

  async ensureUrl() {
    let url = this.getUrl();
    if (!url) {
      url = prompt(
        "Please enter the Base API URL (e.g., https://api.example.com):",
        this.DEFAULT_BASE_URL,
      );
      if (url) {
        this.setUrl(url);
      } else {
        throw new Error("Configuration cancelled by user.");
      }
    }
    return this.getUrl();
  },
};
