import { ApiWrapper, CustomAxiosConfigEnum } from '@/utils/request';

const api = new ApiWrapper('/api/core/i18n');

/**
 * 获取语言包列表
 * 启停用语言包
 * 删除语言包
 */

export const getLangListApi = async (params?: { key?: string }) => {
  const keyword = params?.key?.trim().toLowerCase();
  return [
    { id: 1, languageCode: 'zh-CN', languageName: '中文（简体）', languageType: 1, hasUsed: true },
    { id: 2, languageCode: 'en-US', languageName: 'English', languageType: 1, hasUsed: true },
  ]
    .filter((item) => {
      if (!keyword) return true;
      return item.languageCode.toLowerCase().includes(keyword) || item.languageName.toLowerCase().includes(keyword);
    })
    .map((item) => ({
      ...item,
      label: item.languageName,
      value: item.languageCode,
    }));
};

export const getSystemI18Api = async (lang: string): Promise<{ messages: Record<string, string> }> =>
  api.get('/messages', { params: { lang } });

export const langEnableApi = async (data?: any) => api.put('/languages/enable', data);

export const deleteLangApi = async (code: string) => api.delete(`/languages/${code}`);

/**
 * 获取模块
 * @param moduleType  模块类型，0-所有；1-内置；2-自定义
 * @param keyword  关键词
 *
 * 删除模块
 */
export const getModulesListApi = async (params?: { moduleType?: number; keyword?: string }) =>
  api.get('/modules', { params }).then((res) => {
    return (res || [])?.map((item: any) => ({
      ...item,
      key: item.id,
      title: item.moduleName,
    }));
  });

export const deleteModuleApi = async (moduleCode: string) => api.delete(`/modules/${moduleCode}`);

/**
 * 获取词条/资源 列表
 * 新增词条
 * 编辑词条
 * 删除词条
 */
export const getResourcesListApi = async (data: any) =>
  api.post('/resources/search', data, {
    [CustomAxiosConfigEnum.BusinessResponse]: true,
  });

export const addResourcesApi = async (data: any) => api.post('/resources', data);

export const editResourcesApi = async (data: any) => api.put('/resources', data);

export const deleteResourcesApi = async (moduleCode?: string, key?: string) =>
  api.delete(`/resources/${moduleCode}/${key}`);

/**
 * 下载模板
 *
 * 导出语言包
 *
 * 下载语言包
 *
 * 获取导出记录
 *
 * 确认导出记录
 */

export const downloadTemplateApi = async (params?: Record<string, unknown>) =>
  api.get('/excel/template/download', { params, responseType: 'blob', [CustomAxiosConfigEnum.NoCode]: true });

export const exportLanguageApi = async (data?: Record<string, unknown>) => api.post('/excel/data/export', data);

export const downloadLanguageFileApi = async (params?: Record<string, unknown>) =>
  api.get('/excel/download', { params, responseType: 'blob', [CustomAxiosConfigEnum.NoCode]: true });

export const getLanguageRecordsApi = async (params?: Record<string, unknown>) =>
  api.get('/excel/data/getExportRecords', { params });

export const languageRecordConfirmApi = async (data?: Record<string, unknown>) =>
  api.post('/excel/data/exportRecordConfirm', data);

export const importLanguageFileApi = async (data: any) =>
  api.upload('/excel/template/import', data, {
    method: 'post',
  });
