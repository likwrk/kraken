/**
 * Application Constants
 */
export const DEVICE_ID = "b2ee2e5c-dd93-41b6-be86-cfb1fda6c480";

export const CHART_TYPES = {
  LINE: "line",
  BAR: "bar",
  SCATTER: "scatter",
};

export const DEFAULT_CHART_CONFIG = {
  id: null, // Auto-generated
  title: "New Chart",
  type: CHART_TYPES.LINE,
  sensorFilter: null, // null = all sensors, or array of sensor_ids
  refreshInterval: null, // ms, null = no auto-refresh
  createdAt: null,
};
