import { useEffect, useRef } from 'react';
import { useLocation } from 'react-router';
import { filter, find, isEmpty, map } from 'lodash-es';
import { SUPOS_USER_GUIDE_ROUTES } from '@/common-types/constans';
import { useTranslate } from '@/hooks';
import { storageOpt } from '@/utils/storage';
import { shepherd } from '@/components/shepherd';
import { MenuTypeEnum, setMenuType, useThemeStore } from '@/stores/theme-store.ts';

/**
 * 使用 新人引导 步骤
 * 注意：对某个页面添加steps时，请务必在 stores -> base -> index.tsx -> GuidePagePaths 中添加该页面路由
 * @param steps 初始化步骤数据
 * @param startStepId 指定起始的步骤id
 */
export const useGuideSteps = (steps: any[] = [], startStepId?: string) => {
  const _menuType = useThemeStore((state) => state.menuType);
  const pathname = useLocation().pathname;
  const tour = useRef(shepherd()).current;
  const formatMessage = useTranslate();

  useEffect(() => {
    const startTour = () => {
      // 过滤出存在的步骤数据（根据元素是否存在进行判断）
      const availableSteps = filter(
        steps,
        (step) =>
          isEmpty(step.attachTo?.element) || (step.attachTo?.element && document.querySelector(step.attachTo.element))
      );
      // 如果有步骤数据
      if (availableSteps && availableSteps.length > 0) {
        if (availableSteps.length === 1 && !availableSteps[0].buttons) {
          availableSteps[0].buttons = [
            {
              action() {
                return this.complete();
              },
              text: formatMessage('global.tipDone'),
            },
          ];
        } else {
          if (!availableSteps[0].buttons) {
            availableSteps[0].buttons = [
              {
                action() {
                  return this.complete();
                },
                text: formatMessage('global.tipExit'),
                classes: 'prev-class',
              },
              {
                action() {
                  return tour.next();
                },
                text: formatMessage('common.next'),
              },
            ];
          }
          if (!availableSteps[availableSteps.length - 1].buttons) {
            availableSteps[availableSteps.length - 1].buttons = [
              {
                action() {
                  return this.back();
                },
                text: formatMessage('common.prev'),
                classes: 'prev-class',
              },
              {
                action() {
                  return this.complete();
                },
                text: formatMessage('global.tipDone'),
              },
            ];
          }
        }
        tour.addSteps(availableSteps);
        // 如果存在指定起始的stepId，则从指定的步骤开始
        // 否则从第一个步骤开始
        if (startStepId) {
          tour.show(startStepId);
        } else {
          tour.start();
        }
        tour.on('cancel', () => {
          tour.complete();
        });
        tour.on('complete', () => {
          // 监听完成时把当前路由isVisited设置为已浏览
          const currentUserGuideRoute = storageOpt.get(SUPOS_USER_GUIDE_ROUTES);
          storageOpt.set(
            SUPOS_USER_GUIDE_ROUTES,
            map(currentUserGuideRoute, (route) => (route?.url === pathname ? { ...route, isVisited: true } : route))
          );
        });
      }
    };

    const userGuideRoute = storageOpt.get(SUPOS_USER_GUIDE_ROUTES);
    const currentRoute = find(userGuideRoute, (route) => route?.url === pathname);
    // 当前路由没有被访问过，则初始化当前路由的步骤数据
    if (currentRoute && currentRoute?.isVisited === false) {
      const menuType = _menuType;
      if (menuType !== MenuTypeEnum.Top) {
        setMenuType(MenuTypeEnum.Top);
        setTimeout(() => {
          // 需要等菜单渲染后才能初始化数据，此时需导航指引的id才存在
          startTour();
        }, 200);
      } else {
        startTour();
      }
    }
  }, []);

  return {
    tour,
  };
};
