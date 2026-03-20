/**
 * Color Palette Utility
 */
export const ColorPalette = {
  cache: new Map(),

  getColor(key) {
    if (this.cache.has(key)) {
      return this.cache.get(key);
    }
    let hash = 0;
    for (let i = 0; i < key.length; i++) {
      hash = key.charCodeAt(i) + ((hash << 5) - hash);
    }
    const hue = Math.abs(hash) % 360;
    const borderColor = `hsl(${hue}, 70%, 50%)`;
    const backgroundColor = `hsla(${hue}, 70%, 50%, 0.1)`;

    const result = { borderColor, backgroundColor };
    this.cache.set(key, result);
    return result;
  },
};
