import type { MenuStoreProps, MenuStoreState } from './types.ts';
import { createStore } from 'zustand/vanilla';
import { immer } from 'zustand/middleware/immer';
import { createContext, type ReactNode, useContext, useState } from 'react';
import { useStoreWithEqualityFn } from 'zustand/traditional';
import { shallow } from 'zustand/shallow';
import { getRoutesResourceApi } from '@/apis/inter-api/resource.ts';
import { listToTree } from '@/pages/menu-configuration/store/utils.ts';

const initialState: MenuStoreState = {
  menuList: [],
  menuTree: [],
  contentType: null,
  selectNode: null,
  loading: false,
};

export const createMenuStore = (initProps?: Partial<MenuStoreProps>) => {
  return createStore<MenuStoreProps>()(
    immer((set) => ({
      ...initialState,
      ...initProps,
      requestMenu: async () => {
        set({
          loading: true,
        });
        return await getRoutesResourceApi()
          .then((data) => {
            set({
              menuList: data,
              menuTree: listToTree(data),
            });
            return data;
          })
          .finally(() =>
            set({
              loading: false,
            })
          );
      },
      setContentType: (type) =>
        set({
          contentType: type,
        }),
      setSelectNode: (node) => {
        return set({ selectNode: node });
      },
      setMenuInfo: (menuTree, menuList) => {
        return set({ menuTree, menuList });
      },
    }))
  );
};

const MenuStoreContext = createContext<ReturnType<typeof createMenuStore> | null>(null);

export function MenuStoreProvider({ children }: { children: ReactNode }) {
  const [TreeStoreProps] = useState(() => createMenuStore());

  return <MenuStoreContext.Provider value={TreeStoreProps}>{children}</MenuStoreContext.Provider>;
}

export function useMenuStore<U>(selector: (state: MenuStoreProps) => U) {
  const store = useContext(MenuStoreContext);

  if (store === null) {
    throw new Error('useTreeStore must be used within TreeStoreProvider');
  }

  return useStoreWithEqualityFn(store, selector, shallow);
}
