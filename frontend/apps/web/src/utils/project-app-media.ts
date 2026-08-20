import { uploadAttachment } from '@/apis/core-api';

export const projectAppMediaUrl = (assetId?: number) => (assetId ? `/api/core/assets/${assetId}/download` : undefined);

export const uploadProjectAppMedia = async (file: File, ownerId: string | number, ownerType = 'projectAppMedia') => {
  const response = await uploadAttachment([{ value: file, name: 'files', fileName: file.name }], {
    ownerType,
    ownerId,
  });
  const list = response?.list ?? response?.data?.list;
  const uploaded = Array.isArray(list) ? list[0] : undefined;
  const assetId = Number(uploaded?.fileId || uploaded?.id || uploaded?.objectName || 0);
  if (!assetId) {
    throw new Error('upload-failed');
  }
  return assetId;
};
