import { ApiWrapper, CustomAxiosConfigEnum } from '@/utils/request';

const baseUrl = '/inter-api/supos/uns';

const api = new ApiWrapper(baseUrl);

export const searchTreeData = async (params?: Record<string, unknown>) => api.get('/search', { params }); // 获取所以树数据
export const getTreeData = async (params?: Record<string, unknown>) => api.get('/tree', { params }); // 获取所以树数据
export const getTypes = async (params?: Record<string, unknown>) => api.get('/types', { params }); // 获取数据类型
export const getLastMsg = async (params?: Record<string, unknown>) => api.get('/getLastMsg', { params }); // 获取最新msg
export const addModel = async (data: any) => api.post('/model', data); // 新增model
export const detectModel = async (data: any) => api.post('/model/detect', data); // 校验model
export const editModel = async (data: any) => api.put('/model', data); // 修改model
export const getModelInfo = async (params?: Record<string, unknown>) => api.get('/model', { params, _noMessage: true }); // 查询模型字段声明
export const getInstanceInfo = async (params?: Record<string, unknown>) =>
  api.get('/instance', { params, _noMessage: true }); // 查询模型字段声明
export const deleteTreeNode = async (params?: Record<string, unknown>) => api.delete('', { params }); // 删除树节点
export const getDashboardList = async (params?: Record<string, unknown>) =>
  api.get('/dashboard', {
    params,
    [CustomAxiosConfigEnum.BusinessResponse]: true,
  }); // 获取dashboard
export const addDashboard = async (data: any) => api.post('/dashboard', data); // 新增dashboard
export const editDashboard = async (data: any) => api.put('/dashboard', data); // 编辑dashboard
export const getDashboardByUns = async (unsAlias: string, config?: any) =>
  api.get(`/dashboard/getByUns?unsAlias=${unsAlias}`, config); // 获取dashboard信息
export const bindDashboardForUns = async (data: any) => api.post(`/dashboard/bindUns`, data); // 获取dashboard信息
// 置顶
export const markDashboard = async (id: string) => api.post('/dashboard/mark', { id });
export const unmarkDashboard = async (id: string) => api.delete(`/dashboard/unmark?id=${id}`);
export const deleteDashboard = async (uid: string) => api.delete(`/dashboard/${uid}`); // 删除dashboard
/**
 * 首次传入 checkSmallFile: true ,
 * 如果是小文件，直接返回下载内容
 * 大文件 返回 {code: 200, Msg: "ok"} 再次传入参数就为下载内容（checkSmallFile不传）
 * 完全没数据返回 code: 204
 * 错误信息返回 code: 300
 */
export const exportExcel = async (data: any) =>
  api.post('/importExport/export', data, {
    [CustomAxiosConfigEnum.NoCode]: true,
  }); //导出excel
export const searchRestField = async (data: any) => api.post('/searchRestField', data); // 从RestApi搜系模型字段
export const getDashboardDetail = async (id: any) => api.get(`/dashboard/${id}`); // 获取dashboard详情

// 报警列表页
export const getAlertList = async (params?: Record<string, unknown>) =>
  api
    .get('/search', {
      params,
      [CustomAxiosConfigEnum.BusinessResponse]: true,
    })
    .then((data: any) => {
      return {
        ...data,
        pageNo: data?.page?.pageNo,
        pageSize: data?.page?.pageSize,
        total: data?.page?.total,
      };
    });

// 报警列表 options使用
export const getAlertForSelect = async (params?: Record<string, unknown>) =>
  api
    .get('/search', {
      params,
    })
    .then((data: any) => {
      return data?.map?.((item: any) => ({
        label: item.name,
        value: item.id,
      }));
    });

// 轮询获取拓扑图状态
export const getTopologyStatus = async (params?: Record<string, unknown>) =>
  api.get('/topology', { params, [CustomAxiosConfigEnum.NoMessage]: true });
export const getAllLabel = async (params?: Record<string, unknown>) => api.get('/allLabel', { params }); // 获取所有标签
export const addLabel = async (name?: Record<string, unknown>) =>
  api.post(`/label?name=${name ? encodeURIComponent(String(name)) : ''}`); // 获取所有标签
export const deleteLabel = async (id: string | number) =>
  api.delete(`/label?id=${id}`, {
    [CustomAxiosConfigEnum.BusinessResponse]: true,
  }); // 删除标签
export const getAllTemplate = async (data: any, config?: any) => api.post('/template/pageList', data, config); // 获取所有模板
export const getTemplateDetail = async (params: any) => api.get('/template', { params }); // 获取模板详情
export const addTemplate = async (data: any) => api.post('/template', data); // 新增模板
export const makeLabel = async (unsId: string, data: any) => api.post(`/makeLabel?unsId=${unsId}`, data); // 新增model
export const deleteTemplate = async (id: string | number) => api.delete(`/template?id=${id}`); // 删除模板
export const editTemplateName = async (params: any) =>
  api.put(`/template`, null, {
    params,
  }); // 删除模板

// 标签
export const getLabelDetail = async (id: string) => api.get(`/label/detail?id=${id}`); // 获取标签详情
export const getLabelPath = async () => api.get(`/search?type=2`); // 获取path
export const updateLabel = async (data: any) =>
  api.put('/label', data, {
    [CustomAxiosConfigEnum.BusinessResponse]: true,
  }); // 修改标签
export const getLabelUnsId = async (path: any) => api.get(`/instance?id=${path}`); // 获取标签详情
export const verifyFileName = async (params: any) => api.get(`/name/duplication`, { params }); // 检验文件夹文件重名
export const triggerRestApi = async (params: any) => api.post('/triggerRestApi', {}, { params }); // 主档触发restapi
export const ds2fs = async (data: any) => api.post('/ds2fs', data); // 外部数据源转uns接口
export const json2fs = async (data: any) => api.post('/json2fs', data); // JSON 转 UNS 接口
export const json2fsTree = async (data: any) => api.post('/json2fs/tree', data); // JSON 转 UNSTree 接口
export const batchReverser = async (data: any) =>
  api.post('/batch', data, {
    [CustomAxiosConfigEnum.BusinessResponse]: true,
    [CustomAxiosConfigEnum.NoMessage]: true,
  }); // 批量提交topic

export const modifyModel = async (data: any) => api.put('/name', data); // 修改文件夹或文件
export const modifyDetail = async (data: any) => api.put('/detail', data); // 修改文件夹或文件详情

export const getUnsLazyTree = async (
  data: { parentId?: string; keyword?: string; pageNo: number; pageSize: number; searchType?: number },
  config?: any
) => api.post('/condition/tree', data, config);

export const pageListUnsByTemplate = async (params: any) =>
  api.get(`/label/pageListUnsByTemplate`, { params, [CustomAxiosConfigEnum.BusinessResponse]: true }); // 模版关联文件夹文件
export const pageListUnsByLabel = async (params: any) =>
  api.get(`/label/pageListUnsByLabel`, { params, [CustomAxiosConfigEnum.BusinessResponse]: true }); // 标签关联文件夹文件

export const cancelLabel = async (id: string, data: any) => api.delete(`/cancelLabel?unsId=${id}`, {}, data); // 删除标签关联的文件

export const makeSingleLabel = async (unsId: string, labelId: string) =>
  api.post(`/makeSingleLabel?unsId=${unsId}&labelId=${labelId}`); // 增加标签关联的文件

// ====================== 国际化相关 - start ==============
// 设置国际化语言
export const updatePersonConfigApi = async (data: { userId: string; mainLanguage: string }) =>
  api.post('/person/config', data);
// 获取国际化语言
export const getPersonConfigApi = async (userId: string) =>
  api.get('/person/config', {
    params: {
      userId,
    },
  });
// 获取系统语言
export const getSystemI18Api = async (lang: string) =>
  api.get('/i18n/messages', {
    params: {
      lang,
    },
  });

// 获取插件语言
export const getPlugI18Api = async (lang: string, pluginId: string[]) =>
  api.get('/i18n/messages/plugin', {
    params: {
      lang,
      pluginId,
    },
  });

// ====================== 国际化相关 - end ================

export const createDashboard = async (alias: string) => await api.post(`/dashboard/createGrafanaByUns/${alias}`);

// 获取导出结果
export const getUnsExportRecordsApi = async (params?: Record<string, unknown>) =>
  api.get('/excel/data/getExportRecords', params);

// 确认导出记录
export const unsExportRecordConfirmApi = async (params?: Record<string, unknown>) =>
  api.post('/excel/data/exportRecordConfirm', params);

// 文件下载
export const downloadUnsFile = async (params?: Record<string, unknown>) =>
  api.get('/importExport/file/download', { params, responseType: 'blob', [CustomAxiosConfigEnum.NoCode]: true });

export const detectIfRemoveApi = (params: { id: any }) => api.get('/detectIfRemove', { params });

// 订阅
export const updateModelSubscribe = async (params: any) => {
  const { id, enable, frequency } = params;
  let query = `id=${id}&enable=${enable}`;
  if (frequency) {
    query += `&frequency=${frequency}`;
  }
  return api.put(`/model/subscribe?${query}`);
};

// 订阅
export const updateTemplateSubscribe = async (params: any) => {
  const { id, enable, frequency } = params;
  let query = `id=${id}&enable=${enable}`;
  if (frequency) {
    query += `&frequency=${frequency}`;
  }
  return api.put(`/template/subscribe?${query}`);
};

// 订阅
export const updateLabelSubscribe = async (params: any) => {
  const { id, enable, frequency } = params;
  let query = `id=${id}&enable=${enable}`;
  if (frequency) {
    query += `&frequency=${frequency}`;
  }
  return api.put(`/label/subscribe?${query}`);
};

// 查询文件夹订阅
export const subscribeFolderPage = async (params?: Record<string, unknown>) =>
  api.post('/subscribe/folder', params, {
    [CustomAxiosConfigEnum.BusinessResponse]: true,
  });

// 查询文件订阅
export const subscribeFilePage = async (params?: Record<string, unknown>) =>
  api.post('/subscribe/file', params, {
    [CustomAxiosConfigEnum.BusinessResponse]: true,
  });

// 查询模板订阅
export const subscribeTemplatePage = async (params?: Record<string, unknown>) =>
  api.post('/subscribe/template', params, {
    [CustomAxiosConfigEnum.BusinessResponse]: true,
  });

// 查询标签订阅
export const subscribeLabelPage = async (params?: Record<string, unknown>) =>
  api.post('/subscribe/label', params, {
    [CustomAxiosConfigEnum.BusinessResponse]: true,
  });

export const batchWriteFileValue = async (data: any) =>
  api.post('/file/current/batchUpdate', data, {
    [CustomAxiosConfigEnum.BusinessResponse]: true,
    [CustomAxiosConfigEnum.NoCode]: true,
  }); // 文件数据写值

// schema 获取接口
export const getFileSchema = async () => api.get('/file/schema', { [CustomAxiosConfigEnum.NoCode]: true });
export const getFolderSchema = async () => api.get('/folder/schema', { [CustomAxiosConfigEnum.NoCode]: true });
export const getTemplateSchema = async () => api.get('/template/schema', { [CustomAxiosConfigEnum.NoCode]: true });
export const getLabelSchema = async () => api.get('/label/schema', { [CustomAxiosConfigEnum.NoCode]: true });
export const checkDashboardIsExist = async (params?: Record<string, unknown>) =>
  api.get('/dashboard/isExist', { params, [CustomAxiosConfigEnum.NoCode]: true });
export const getEmptyFolder = async () => api.get(`/folder/empty`); // 获取所有空文件夹
export const saveMount = async (data: any) =>
  api.post('/mount', data, {
    [CustomAxiosConfigEnum.BusinessResponse]: true,
  }); // 挂载采集器提交
export const getCollectorList = async (params?: Record<string, unknown>) =>
  api.get('/mount/source/collector', { params }); // 获取采集器列表
export const getSourceList = async (params?: Record<string, unknown>) => api.get('/mount/source', { params }); // 获取数据源列表
export const pasteUns = async (data?: { sourceId?: any; targetId?: any; newF?: Record<string, unknown> }) =>
  api.post('/paste', data, {
    [CustomAxiosConfigEnum.BusinessResponse]: true,
    [CustomAxiosConfigEnum.NoCode]: true,
  }); // 黏贴uns文件文件夹
