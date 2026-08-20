import { createWithEqualityFn } from 'zustand/traditional';
import { shallow } from 'zustand/shallow';
import { getUserManageList } from '@/apis/core-api';
import { fetchAllPagesWithConcurrency } from '@/utils';

type UserManageItem = {
  id?: string;
  userId?: string;
  sub?: string;
  email?: string;
  userName?: string;
  firstName?: string;
  preferredUsername?: string;
};

type MemberOption = {
  label: string;
  value: string;
};

/**
 * Project 页面的公共状态。
 */
type ProjectCommonState = {
  /** 后端返回的成员列表，供成员弹窗直接使用 */
  memberOptions: MemberOption[];
  /** 成员列表接口请求状态 */
  memberOptionsLoading: boolean;
  /** 是否已经完成过一次成员列表拉取 */
  memberOptionsLoaded: boolean;
  /** 拉取成员列表：每次调用都以接口返回覆盖本地数据 */
  fetchMemberOptions: () => Promise<void>;
  /** 重置成员列表缓存，下次调用 fetchMemberOptions 时会重新请求 */
  resetMemberOptions: () => void;
};

export const PAGE_SIZE = 100;
export const MAX_CONCURRENT_PAGE_REQUESTS = 5;

/**
 * 这里使用 Zustand 单例 store，不依赖 Provider。
 * 只要在页面里调用 useProjectCommonStore，就会共享同一份状态。
 */
export const useProjectCommonStore = createWithEqualityFn<ProjectCommonState>()(
  (set, get) => ({
    memberOptions: [],
    memberOptionsLoading: false,
    memberOptionsLoaded: false,
    fetchMemberOptions: async () => {
      // 已加载或正在加载时不重复请求，避免每次打开弹窗都触发接口。
      if (get().memberOptionsLoaded || get().memberOptionsLoading) {
        return;
      }

      set({ memberOptionsLoading: true });

      try {
        const result = await fetchAllPagesWithConcurrency<UserManageItem>(
          async ({ pageNo, pageSize }) => {
            const response: any = await getUserManageList({
              pageNo,
              pageSize,
            });

            return {
              total: Number(response?.total ?? 0),
              data: (Array.isArray(response?.data) ? response.data : []) as UserManageItem[],
            };
          },
          {
            pageSize: PAGE_SIZE,
            maxConcurrency: MAX_CONCURRENT_PAGE_REQUESTS,
          }
        );

        const normalizedUsers = result.data;
        const options: MemberOption[] = [];
        const seenValues = new Set<string>();

        normalizedUsers.forEach((user) => {
          const value = String(user.id || user.userId || '');
          if (!value || seenValues.has(value)) {
            return;
          }
          seenValues.add(value);
          options.push({
            label: user.preferredUsername || user.userName || user.firstName || value,
            value,
          });
        });

        // 只存储接口原始返回数据；异常结构时退为空数组。
        set({
          memberOptions: options,
          memberOptionsLoaded: true,
        });
      } finally {
        // 无论成功失败都要结束 loading，避免 UI 一直转圈。
        set({ memberOptionsLoading: false });
      }
    },
    resetMemberOptions: () => {
      set({ memberOptionsLoaded: false, memberOptions: [] });
    },
  }),
  shallow
);
