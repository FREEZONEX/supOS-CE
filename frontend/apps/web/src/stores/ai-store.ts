import { createWithEqualityFn, type UseBoundStoreWithEqualityFn } from 'zustand/traditional';
import type { StoreApi } from 'zustand/index';
import { shallow } from 'zustand/vanilla/shallow';

interface AiResultProps {
  [key: string]: any;
}
export type TAiStore = {
  aiResult?: AiResultProps;
  aiOperationName?: string;
};

export const useAiStore: UseBoundStoreWithEqualityFn<StoreApi<TAiStore>> = createWithEqualityFn(() => ({}), shallow);

export const setAiResult = (key: string, result?: any) => {
  useAiStore.setState({
    aiResult: {
      ...useAiStore.getState().aiResult,
      [key]: result,
    },
  });
};

export const setAiOperationName = (aiOperationName: string) => {
  useAiStore.setState({
    aiOperationName,
  });
};
