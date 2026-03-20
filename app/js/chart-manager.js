/**
 * Chart Manager Service
 * Manages an array of Chart.js instances with dynamic add/remove/update
 */
import { DataMapper } from "./data-mapper.js";
import { CHART_TYPES } from "./constants.js";

export const ChartManager = {
  charts: new Map(), // chartId -> { instance, config, containerId }

  /**
   * Create a new chart instance
   * @param {string} containerId - DOM id of canvas container
   * @param {Object} config - Chart configuration
   * @param {Array} rawData - Sensor data array
   */
  create(containerId, config, rawData) {
    const canvas = document.createElement("canvas");
    canvas.id = `chart-${config.id}`;
    const container = document.getElementById(containerId);
    if (!container) {
      throw new Error(`Container #${containerId} not found`);
    }
    container.innerHTML = ""; // Clear existing
    container.appendChild(canvas);

    const chartData = DataMapper.transform(rawData, config.sensorFilter);

    const chartConfig = {
      type: config.type || CHART_TYPES.LINE,
      data: {
        labels: chartData.labels,
        datasets: chartData.datasets,
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        interaction: {
          mode: "index",
          intersect: false,
        },
        plugins: {
          zoom: {
            drag: {
              enabled: true, // ✅ Enable drag selection
              backgroundColor: "rgba(59, 130, 246, 0.3)",
              modifierKey: null, // No key required
            },
            pan: {
              enabled: true,
              mode: "x", // Pan horizontally only
              modifierKey: null, // No key required
            },
            zoom: {
              wheel: {
                enabled: true,
              },
              pinch: {
                enabled: true,
              },
              mode: "x", // Zoom horizontally only
              limitZoom: {
                min: 1000 * 60 * 5, // Min zoom: 5 minutes (ms)
                max: 1000 * 60 * 60 * 24 * 7, // Max zoom: 7 days (ms)
              },
            },
            limits: {
              x: { min: "original", max: "original" }, // Prevent panning beyond data
            },
          },
          legend: {
            position: "top",
            labels: { usePointStyle: true, pointStyle: "line" },
          },
          tooltip: {
            callbacks: {
              label: (context) => {
                const sensorId = context.dataset.label;
                const value = context.parsed.y;
                const timestamp = new Date(
                  context.parsed.x,
                ).toLocaleTimeString();
                return [
                  `Sensor: ${sensorId}`,
                  `Time: ${timestamp}`,
                  `Value: ${value}`,
                ];
              },
              title: () => "Reading Details",
            },
          },
        },
        scales: {
          y: {
            beginAtZero: false,
            title: { display: true, text: "Value" },
            grid: { color: "rgba(0,0,0,0.05)" },
          },
          x: {
            type: "time",
            time: {
              tooltipFormat: "MMM d, yyyy HH:mm:ss",
              displayFormats: {
                second: "HH:mm:ss",
                minute: "HH:mm",
                hour: "HH:mm",
                day: "MMM d",
              },
            },
            title: { display: true, text: "Time" },
            grid: { display: false },
          },
        },
      },
    };

    const instance = new Chart(canvas, chartConfig);

    this.charts.set(config.id, {
      instance,
      config,
      containerId,
      canvas,
    });

    return instance;
  },

  /**
   * Update an existing chart with new data
   */
  update(chartId, rawData) {
    const chart = this.charts.get(chartId);
    if (!chart) return false;

    const chartData = DataMapper.transform(rawData, chart.config.sensorFilter);
    chart.instance.data.labels = chartData.labels;
    chart.instance.data.datasets = chartData.datasets;
    chart.instance.update();
    return true;
  },

  /**
   * Remove a chart instance
   */
  remove(chartId) {
    const chart = this.charts.get(chartId);
    if (!chart) return false;

    chart.instance.destroy();
    const container = document.getElementById(chart.containerId);
    if (container) {
      container.innerHTML = "";
    }
    this.charts.delete(chartId);
    return true;
  },

  /**
   * Get chart instance by ID
   */
  get(chartId) {
    return this.charts.get(chartId)?.instance || null;
  },

  /**
   * Get all chart configs
   */
  getAllConfigs() {
    return Array.from(this.charts.values()).map((c) => c.config);
  },

  /**
   * Refresh all charts with new data
   */
  async refreshAll(fetchFn, baseUrl) {
    const rawData = await fetchFn(baseUrl);
    for (const [chartId] of this.charts) {
      this.update(chartId, rawData);
    }
  },
};
