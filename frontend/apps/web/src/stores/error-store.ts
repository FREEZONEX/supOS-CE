import { createWithEqualityFn, type UseBoundStoreWithEqualityFn } from 'zustand/traditional';
import type { StoreApi } from 'zustand/index';
import { shallow } from 'zustand/vanilla/shallow';

export type TErrorStore = {
  errorInfo?: any;
};

export const useErrorStore: UseBoundStoreWithEqualityFn<StoreApi<TErrorStore>> = createWithEqualityFn(
  () => ({}),
  shallow
);

export const setErrorInfo = (result?: any) => {
  useErrorStore.setState({
    errorInfo: result,
  });
};
