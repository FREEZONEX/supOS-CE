import { ApiWrapper } from '@/utils/request';

const baseUrl = '/inter-api/supos/license';

const api = new ApiWrapper(baseUrl);

export const queryDeadline = async () => api.get('/deadline');
