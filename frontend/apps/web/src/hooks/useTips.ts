import { useEffect, useRef, useState } from 'react';
import { storageOpt } from '@/utils/storage';
import { SUPOS_USER_GUIDE_ROUTES, SUPOS_USER_LAST_LOGIN_ENABLE } from '@/common-types/constans';
import { shepherd } from '@/components/shepherd';
import { updateTipsEnable } from '@/apis/inter-api/user-manage';
import { find, shuffle } from 'lodash-es';
import { useLocation } from 'react-router';
import { useTranslate } from '@/hooks';
import { setUserTipsEnable, useBaseStore } from '@/stores/base';

const getCheckboxText = (checked: boolean, formatMessage: any) =>
  `<input type="checkbox" ${checked ? 'checked' : ''} /> <span class="checkbox-label"> ${formatMessage('global.tipCheckbox', 'Don’t show again at next login')}</span>`;

/**
 * 使用 tips （只在每次重新登录或者免登录时，并且用户允许提示时触发）
 * @param tips 初始化步骤数据
 */
export const useTips = (originTips: any[] = []) => {
  const pathname = useLocation().pathname;
  const formatMessage = useTranslate();

  const userTipsEnable = useBaseStore((state) => state.userTipsEnable);
  const tour = useRef(shepherd()).current;
  const [checked, setChecked] = useState<boolean>(true);
  const loginEnable = storageOpt.getOrigin(SUPOS_USER_LAST_LOGIN_ENABLE); // 是否为免登录

  useEffect(() => {
    // 每次登录出现的tips都打乱步骤顺序，随机出现tips功能点
    const tips = shuffle(originTips);
    // 初始化步骤数据
    if (tips.length === 1) {
      if (!tips[0].buttons) {
        tips[0].buttons = [
          {
            action() {
              setChecked(!checked);
            },
            text: getCheckboxText(checked, formatMessage),
            classes: 'prev-class checkbox-class',
          },
          {
            action() {
              return this.complete();
            },
            text: formatMessage('global.tipDone', 'Done'),
          },
        ];
        if (loginEnable !== 'true') {
          // 免登录时不支持勾选禁用，该状态需要绑定具体用户
          tips[0].buttons.splice(0, 1);
        }
      }
    } else {
      tips.forEach((tip, i) => {
        if (!tip.buttons) {
          tip.buttons = [
            {
              action() {
                setChecked(!checked);
              },
              text: getCheckboxText(checked, formatMessage),
              classes: 'prev-class checkbox-class',
            },
            {
              action() {
                return this.back();
              },
              text: formatMessage('global.tipLastOne', 'Back'),
              classes: 'prev-class',
            },
            {
              action() {
                return this.next();
              },
              text: formatMessage('global.tipNextOne', 'Next'),
            },
          ];
          if (i === 0) {
            tip.buttons[1] = {
              action() {
                return this.complete();
              },
              text: formatMessage('global.tipExit', 'Exit'),
              classes: 'prev-class',
            };
          }
          if (i === tips.length - 1) {
            tip.buttons[tip.buttons.length - 1] = {
              action() {
                return this.complete();
              },
              text: formatMessage('global.tipDone', 'Done'),
            };
          }
          if (loginEnable !== 'true') {
            // 免登录时不支持勾选禁用，该状态需要绑定具体用户
            tip.buttons.splice(0, 1);
          }
        }
      });
    }
    tour.addSteps(tips);
    tour.on('cancel', () => {
      tour.complete();
    });
    tour.on('complete', () => {
      // 监听完成时设置为不开启
      setUserTipsEnable('0');
    });
  }, []);

  useEffect(() => {
    const userGuideRoute = storageOpt.get(SUPOS_USER_GUIDE_ROUTES);
    const currentRoute = find(userGuideRoute, (route) => route?.url === pathname);
    // 新手导航存在则先不触发，新手导航优先级更高
    if (currentRoute && currentRoute?.isVisited === false) {
      setUserTipsEnable('0');
      return;
    }
    // 监听到userTipsEnable改变时判断是否触发启用
    if (userTipsEnable === '1' && !tour.isActive()) {
      tour.start();
    }
  }, [userTipsEnable]);

  useEffect(() => {
    if (loginEnable !== 'true') {
      // 免登录时不支持勾选禁用，该状态需要绑定具体用户
      return;
    }
    tour.steps.forEach((step) => {
      step.updateStepOptions({
        ...step.options,
        buttons: step.options.buttons?.map((btn, i) =>
          i === 0
            ? {
                action() {
                  setChecked(!checked);
                },
                text: getCheckboxText(checked, formatMessage),
                classes: 'prev-class checkbox-class',
              }
            : btn
        ),
      });
    });
    const fn = () => {
      // if (checked) {
      // 调后端接口设置不再触发tips
      updateTipsEnable(checked ? 0 : 1);
      // }
      // 监听完成时设置为不开启
      setUserTipsEnable('0');
    };
    tour.off('complete', fn); // 先删除之前的事件绑定函数
    tour.on('complete', fn);
  }, [checked]);

  return {
    tour,
  };
};
