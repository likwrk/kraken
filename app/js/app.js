/**
 * Application Bootstrap
 */
import { Config } from "./config.js";
import { API } from "./api.js";
import { ChartManager } from "./chart-manager.js";
import { UI } from "./ui.js";
import { DEFAULT_CHART_CONFIG } from "./constants.js";

/**
 * Generate unique ID for charts
 */
const generateId = () =>
  `chart_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;

/**
 * Application Controller
 */
const App = {
  rawData: null, // Cache fetched data

  async init() {
    try {
      const baseUrl = await Config.ensureUrl();
      UI.setStatus(`Connected to ${baseUrl}`);

      this.rawData = await API.fetchSensors(baseUrl);

      if (ChartManager.charts.size === 0) {
        this.addChart({
          ...DEFAULT_CHART_CONFIG,
          id: generateId(),
          title: "All Sensors",
          createdAt: new Date().toISOString(),
        });
      }

      UI.bindEvents({
        onAddChart: (config) =>
          this.addChart({
            ...DEFAULT_CHART_CONFIG,
            ...config,
            id: generateId(),
            createdAt: new Date().toISOString(),
          }),
        onResetConfig: () => {
          Config.clear();
          window.location.reload();
        },
        onChartAction: ({ chartId, action }) => {
          if (action === "remove") {
            ChartManager.remove(chartId);
            const card = document.querySelector(`[data-chart-id="${chartId}"]`);
            if (card) card.remove();
          } else if (action === "refresh") {
            this.refreshChart(chartId);
          }
        },
      });

      UI.setStatus(`Loaded ${this.rawData.length} readings.`, "success");
    } catch (error) {
      console.error("App initialization error:", error);

      // ✅ Critical: Show clear error with recovery option
      UI.setStatus(`Error: ${error.message}. Try resetting API URL.`, "error");

      // ✅ Ensure reset button is still visible and active
      const resetBtn = document.getElementById("btn-reset-config");
      if (resetBtn) {
        resetBtn.style.display = "inline-block";
        resetBtn.disabled = false;
        resetBtn.style.opacity = "1";
      }

      // ✅ Re-bind reset listener as fallback (in case previous was lost)
      resetBtn?.addEventListener(
        "click",
        () => {
          if (confirm("Clear stored API URL and reload?")) {
            sessionStorage.removeItem("APP_API_BASE_URL");
            window.location.reload();
          }
        },
        { once: true },
      );
    }
  },

  /**
   * Add a new chart to the dashboard
   */
  addChart(config) {
    // Create DOM card
    const card = UI.createChartCard(config);
    UI.chartsContainer.appendChild(card);

    // Create Chart.js instance
    ChartManager.create(`container-${config.id}`, config, this.rawData);
  },

  /**
   * Refresh a specific chart with latest data
   */
  async refreshChart(chartId) {
    try {
      const baseUrl = Config.getUrl();
      if (!baseUrl) throw new Error("API URL not configured");

      const freshData = await API.fetchSensors(baseUrl);
      this.rawData = freshData; // Update cache
      ChartManager.update(chartId, freshData);
      UI.setStatus("Chart refreshed.", "success");
    } catch (error) {
      UI.setStatus(`Refresh failed: ${error.message}`, "error");
    }
  },

  /**
   * Refresh all charts (for future auto-refresh feature)
   */
  async refreshAll() {
    try {
      const baseUrl = Config.getUrl();
      if (!baseUrl) return;

      this.rawData = await API.fetchSensors(baseUrl);
      ChartManager.refreshAll(() => Promise.resolve(this.rawData), baseUrl);
    } catch (error) {
      console.error("Bulk refresh error:", error);
    }
  },
};

// Bootstrap
App.init();

// Optional: Expose for debugging
window.SensorApp = App;
