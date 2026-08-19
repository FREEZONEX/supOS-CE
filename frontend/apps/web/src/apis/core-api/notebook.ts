import { CustomAxiosConfigEnum } from '@/utils/request';
import { coreApi } from './core-adapter';

export const getNotebookList = async (data: Record<string, unknown>): Promise<any> =>
  coreApi.get('/notebooks', {
    params: {
      listType: data?.listType,
      folderId: data?.folderId,
      search: data?.search,
      favoriteOnly: data?.favoriteOnly,
    },
  });

export const getNotebook = async (id: string | number): Promise<any> => coreApi.get(`/notebooks/${id}`);

export const createNotebook = async (data: Record<string, unknown>): Promise<any> => coreApi.post('/notebooks', data);

export const updateNotebook = async (id: string | number, data: Record<string, unknown>): Promise<any> =>
  coreApi.put(`/notebooks/${id}`, data);

export const deleteNotebook = async (id: string | number): Promise<any> => coreApi.delete(`/notebooks/${id}`);

export const cloneNotebook = async (id: string | number): Promise<any> => coreApi.post(`/notebooks/${id}/clone`, {});

export const shutdownNotebook = async (id: string | number): Promise<any> =>
  coreApi.post(`/notebooks/${id}/shutdown`, {});

export const favoriteNotebook = async (id: string | number): Promise<any> => {
  void id;
  return {};
};
export const unfavoriteNotebook = async (id: string | number): Promise<any> => {
  void id;
  return {};
};

export const getNotebookContent = async (id: string | number): Promise<any> => coreApi.get(`/notebooks/${id}/content`);

export const saveNotebookContent = async (id: string | number, data: Record<string, unknown>) =>
  coreApi.put(`/notebooks/${id}/content`, data);

export const getFolderTree = async (): Promise<any> => coreApi.get('/notebook-folders');

export const createFolder = async (data: Record<string, unknown>): Promise<any> =>
  coreApi.post('/notebook-folders', data);

export const updateFolder = async (id: string | number, data: Record<string, unknown>): Promise<any> =>
  coreApi.put(`/notebook-folders/${id}`, data);

export const deleteFolder = async (id: string | number): Promise<any> => coreApi.delete(`/notebook-folders/${id}`);

export const importNotebook = async (formData: FormData): Promise<any> =>
  coreApi.post('/notebook-import', formData, { [CustomAxiosConfigEnum.NoMessage]: true });

export const exportNotebook = async (id: string | number, data: { exportFormat: string; includeCode?: boolean }) =>
  coreApi.post(`/notebooks/${id}/export`, data, {
    [CustomAxiosConfigEnum.NoCode]: true,
    responseType: 'blob',
  });

export const createNotebookSnapshot = async (id: string | number, data: Record<string, unknown>): Promise<any> =>
  coreApi.post(`/notebooks/${id}/snapshots`, data);

export const getNotebookSnapshotList = async (id: string | number, data: Record<string, unknown>): Promise<any> =>
  coreApi.get(`/notebooks/${id}/snapshots`, {
    params: {
      pageNo: data?.pageNo,
      pageSize: data?.pageSize,
    },
  });

export const revertNotebookSnapshot = async (id: string | number, snapshotId: string | number): Promise<any> =>
  coreApi.post(`/notebooks/${id}/snapshots/${snapshotId}/revert`, {});

export const saveAsNotebookSnapshot = async (
  id: string | number,
  snapshotId: string | number,
  data: Record<string, unknown>
): Promise<any> => coreApi.post(`/notebooks/${id}/snapshots/${snapshotId}/save-as`, data);
