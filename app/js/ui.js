/**
 * UI Service - DOM helpers and modal management
 */
import { DEFAULT_CHART_CONFIG, CHART_TYPES } from "./constants.js";

export const UI = {
  statusEl: document.getElementById("status-msg"),
  chartsContainer: document.getElementById("charts-container"),
  addChartBtn: document.getElementById("btn-add-chart"),
  resetConfigBtn: document.getElementById("btn-reset-config"),

  setStatus(message, type = "success") {
    this.statusEl.textContent = message;
    this.statusEl.className = `status-msg ${type}`;
  },

  /**
   * Create DOM structure for a new chart card
   */
  createChartCard(config) {
    const card = document.createElement("article");
    card.className = "chart-card";
    card.dataset.chartId = config.id;

    card.innerHTML = `
            <header class="chart-card-header">
                <h3>${config.title}</h3>
                <div class="chart-actions">
                    <button class="btn-text btn-small btn-remove" data-action="remove">Remove</button>
                    <button class="btn-text btn-small btn-refresh" data-action="refresh">Refresh</button>
                </div>
            </header>
            <div class="chart-card-body">
                <div class="chart-container" id="container-${config.id}"></div>
            </div>
            <footer class="chart-card-footer">
                <span class="chart-meta">Created: ${new Date(config.createdAt).toLocaleTimeString()}</span>
            </footer>
        `;

    return card;
  },

  /**
   * Show modal for adding a new chart
   */
  showAddChartModal(onSubmit) {
    // ✅ Close any existing modal first
    const existingModal = document.querySelector(".modal-overlay");
    if (existingModal) {
      existingModal.remove();
    }
    const modal = document.createElement("div");
    modal.className = "modal-overlay";
    modal.innerHTML = `
            <div class="modal">
                <h3>Add New Chart</h3>
                <form class="modal-form" id="add-chart-form">
                    <label>
                        Title
                        <input type="text" name="title" value="New Chart" required>
                    </label>
                    <label>
                        Chart Type
                        <select name="type">
                            <option value="${CHART_TYPES.LINE}">Line</option>
                            <option value="${CHART_TYPES.BAR}">Bar</option>
                            <option value="${CHART_TYPES.SCATTER}">Scatter</option>
                        </select>
                    </label>
                    <label>
                        Filter Sensors (comma-separated IDs, empty for all)
                        <input type="text" name="sensorFilter" placeholder="temp1,temp2">
                    </label>
                    <div class="modal-actions">
                        <button type="button" class="btn-text" id="modal-cancel">Cancel</button>
                        <button type="submit" class="btn-primary">Create Chart</button>
                    </div>
                </form>
            </div>
        `;

    document.body.appendChild(modal);

    const form = modal.querySelector("#add-chart-form");
    const cancelBtn = modal.querySelector("#modal-cancel");

    const handleSubmit = (e) => {
      e.preventDefault();
      const formData = new FormData(form);
      const sensorFilter = formData.get("sensorFilter")
        ? formData
            .get("sensorFilter")
            .split(",")
            .map((s) => s.trim())
            .filter(Boolean)
        : null;

      onSubmit({
        title: formData.get("title"),
        type: formData.get("type"),
        sensorFilter,
      });
      modal.remove();
    };

    form.addEventListener("submit", handleSubmit);
    cancelBtn.addEventListener("click", () => modal.remove());
  },

  /**
   * Bind global event listeners
   */
  // In js/ui.js - bindEvents method

  bindEvents({ onAddChart, onResetConfig, onChartAction }) {
    // Add Chart Button
    this.addChartBtn.addEventListener("click", () => {
      this.showAddChartModal(onAddChart);
    });

    // ✅ Reset Button - Direct listener, always active
    this.resetConfigBtn.addEventListener("click", (e) => {
      e.preventDefault();
      e.stopPropagation(); // Prevent any parent handlers from blocking

      // Confirm immediately
      if (confirm("Clear stored API URL and reload?")) {
        // Clear storage directly here as backup
        sessionStorage.removeItem("APP_API_BASE_URL");
        // Force reload
        window.location.href = window.location.href;
      }
    });

    // Chart Card Actions (Remove/Refresh) - Delegated
    this.chartsContainer.addEventListener("click", (e) => {
      const btn = e.target.closest("button[data-action]");
      if (!btn) return;

      const card = btn.closest(".chart-card");
      if (!card) return;

      const chartId = card.dataset.chartId;
      const action = btn.dataset.action;

      onChartAction({ chartId, action });
    });
  },
};
