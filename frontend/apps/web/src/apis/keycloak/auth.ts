import { ApiWrapper, CustomAxiosConfigEnum } from '@/utils/request';

export const baseUrl = '';

const api = new ApiWrapper(baseUrl);

export const getKeycloakToken = async (data: any) =>
  api.post(
    `/keycloak/home/auth/realms/${import.meta.env.REACT_APP_HOMEPAGE_KEYCLOAK_REALMS || 'supos'}/protocol/openid-connect/token`,
    data,
    {
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      [CustomAxiosConfigEnum.BusinessResponse]: true,
      [CustomAxiosConfigEnum.NoCode]: true,
    }
  ); // 获取所以树数据
