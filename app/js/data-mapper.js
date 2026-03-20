/**
 * Data Mapper - Transforms raw API data for Chart.js
 */
import { ColorPalette } from "./color-palette.js";

export const DataMapper = {
  groupBySensor(rawData, sensorFilter = null) {
    let filtered = rawData;
    if (
      sensorFilter &&
      Array.isArray(sensorFilter) &&
      sensorFilter.length > 0
    ) {
      filtered = rawData.filter((item) =>
        sensorFilter.includes(item.sensor_id),
      );
    }

    return filtered.reduce((groups, item) => {
      const sensorId = item.sensor_id;
      if (!groups.has(sensorId)) {
        groups.set(sensorId, []);
      }
      groups.get(sensorId).push(item);
      return groups;
    }, new Map());
  },

  transform(rawData, sensorFilter = null) {
    if (!Array.isArray(rawData) || rawData.length === 0) {
      return { labels: [], datasets: [] };
    }

    const grouped = this.groupBySensor(rawData, sensorFilter);
    const datasets = [];

    for (const [sensorId, readings] of grouped) {
      readings.sort((a, b) => new Date(a.timestamp) - new Date(b.timestamp));

      const dataPoints = readings.map((r) => ({
        x: r.timestamp,
        y: r.value,
      }));

      const colors = ColorPalette.getColor(sensorId);

      datasets.push({
        label: sensorId,
        data: dataPoints,
        borderColor: colors.borderColor,
        backgroundColor: colors.backgroundColor,
        tension: 0.3,
        fill: false,
        pointRadius: 3,
        pointHoverRadius: 6,
        _meta: readings,
      });
    }

    return { labels: [], datasets };
  },
};
