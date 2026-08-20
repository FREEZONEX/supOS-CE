const deviceTokenPattern = /^[a-f0-9]{64}$/;

export const normalizeDeviceToken = (token?: string) => (token || '').trim().toLowerCase();

export const isValidDeviceToken = (token?: string) => deviceTokenPattern.test(normalizeDeviceToken(token));

export const computeDeviceTokenCheckCode = (token?: string) => {
  const normalized = normalizeDeviceToken(token);
  if (!isValidDeviceToken(normalized)) {
    return '';
  }

  let checksum = 0;
  for (let i = 0; i < normalized.length; i += 1) {
    const nibble = Number.parseInt(normalized[i], 16);
    checksum = (checksum + nibble * (i + 1)) & 0xffff;
    checksum = ((checksum << 1) | (checksum >> 15)) & 0xffff;
  }

  return checksum.toString(16).padStart(4, '0').toUpperCase();
};
