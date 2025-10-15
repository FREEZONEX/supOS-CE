import { ApiWrapper, CustomAxiosConfigEnum } from '@/utils/request';

export const baseUrl = '/chat2db/api';

const api = new ApiWrapper(baseUrl);

/** 用户登录接口 */
export const userLogin = async (
  data: any = {
    password: 'chat2db',
    userName: 'chat2db',
  }
) =>
  api.post('/oauth/login_a', data, {
    [CustomAxiosConfigEnum.NoCode]: true,
  });

/** 获取用户信息 */
export const getUser = async () =>
  api.get('/oauth/user_a', {
    [CustomAxiosConfigEnum.NoCode]: true,
  });

/** 获取数据来源列表 */
export const getSourceList = async () =>
  api.get('/connection/datasource/list?pageNo=1&pageSize=1000&refresh=true', {
    [CustomAxiosConfigEnum.NoCode]: true,
  });

/** 获取数据库列表 */
export const getDatabaseList = async (params: any) =>
  api.get('/rdb/database/list', {
    params,
    [CustomAxiosConfigEnum.NoCode]: true,
  });

/** 获取schema列表 */
export const getSchemaList = async (params: any) =>
  api.get('/rdb/schema/list', {
    params,
    [CustomAxiosConfigEnum.NoCode]: true,
  });

/** 获取table列表 */
export const getTableList = async (params: any) =>
  api.get('/rdb/table/table_list', {
    params,
    [CustomAxiosConfigEnum.NoCode]: true,
  });

/** 获取字段列表 */
export const getColumnList = async (params: any) =>
  api.get('/rdb/ddl/column_list', {
    params,
    [CustomAxiosConfigEnum.NoCode]: true,
  });

/** 刷新数据库 */
export const getRefreshList = async (params: any) =>
  api.get('/rdb/table/list', {
    [CustomAxiosConfigEnum.NoCode]: true,
    params: {
      dataSourceName: '@postgresql',
      databaseType: 'POSTGRESQL',
      databaseName: 'postgres',
      schemaName: 'public',
      refresh: true,
      pageNo: 1,
      pageSize: 1000,
      ...params,
    },
  });
