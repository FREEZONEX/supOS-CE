import { ApiWrapper, CustomAxiosConfigEnum } from '@/utils';

const api = new ApiWrapper('/api/core/cli-auth');

export type CliAuthStatus = 'pending' | 'completed' | 'expired' | 'denied';

export interface CliAuthStatusReq {
  setupCode: string;
}

export interface CliAuthStatusResp {
  status: CliAuthStatus;
  apiKey: string;
  workspaceName: string;
  expiresAt: string;
}

export interface CliAuthBindReq {
  setupCode: string;
  name?: string;
  permission?: string;
}

export interface CliAuthBindResp {
  apiKey: string;
  name: string;
  permission: string;
  workspaceName?: string;
}

export const getCliAuthStatus = async (data: CliAuthStatusReq): Promise<CliAuthStatusResp> =>
  api.post('/status', data, { [CustomAxiosConfigEnum.NoMessage]: true });

export const cliAuthBind = async (data: CliAuthBindReq): Promise<CliAuthBindResp> =>
  api.post('/bind', data, { [CustomAxiosConfigEnum.NoMessage]: true });
