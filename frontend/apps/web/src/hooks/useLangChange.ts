import { useEffect, useRef } from 'react';
import { useI18nStore } from '@/stores/i18n-store.ts';
import { useTabsContext } from '@/contexts/tabs-context.ts';

const useLangChange = ({ route }: { route?: string }) => {
  const lang = useI18nStore((state) => state.lang);
  const isFirstRender = useRef(true);
  const { TabsContext } = useTabsContext();

  useEffect(() => {
    if (isFirstRender.current) {
      isFirstRender.current = false;
    } else {
      if (!route) return;
      TabsContext?.current?.onRefreshTab?.(route);
    }
  }, [lang]);
};

export default useLangChange;
