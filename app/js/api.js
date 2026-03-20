/**
 * API Service
 */
import { DEVICE_ID } from "./constants.js";

export const API = {
  async fetchSensors(baseUrl, deviceId = DEVICE_ID) {
    const response = await fetch(`${baseUrl}/sensors?device=${deviceId}`);
    if (!response.ok) {
      throw new Error(`API Error: ${response.status}`);
    }
    return await response.json();
  },

  async fetchSensorsSet(baseUrl, deviceId = DEVICE_ID) {
    const response = await fetch(`${baseUrl}/sensors_set?device=${deviceId}`);
    if (!response.ok) {
      throw new Error(`API Error: ${response.status}`);
    }
    return await response.json();
  },
};
