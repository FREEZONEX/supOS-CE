import Shepherd, { type TourOptions } from 'shepherd.js';
import './index.scss';
import { getIntl } from '@/stores/i18n-store.ts';

// 统一参数
export const defaultConfig: TourOptions = {
  // 是否显示黑色遮罩层
  useModalOverlay: true,
  // 键盘按钮控制步骤
  keyboardNavigation: false,
  // 这里是创建了一个默认的 导航组件
  defaultStepOptions: {
    classes: 'shepherd-custom-classes', // 可以自定义类名，方便调整一些样式什么的不会影响到其他的
    // 显示关闭按钮
    cancelIcon: {
      enabled: true,
    },
    // "center" | "end" | "start" | "nearest"
    scrollTo: { inline: 'center', block: 'center' },
    // 高亮元素四周要填充的空白像素
    modalOverlayOpeningPadding: 6,
    // 空白像素的圆角
    modalOverlayOpeningRadius: 3,
    // modalOverlayOpeningxOffset: -10,
    // modalOverlayOpeningYOffset: -20,
    buttons: [
      {
        // 定义的按钮
        action() {
          return this.back();
        },
        text: getIntl('common.prev'),
        classes: 'prev-class',
      },
      {
        action() {
          return this.next();
        },
        text: getIntl('common.next'),
      },
    ],
  },
};

const shepherd = (props = {}) => {
  const newProps: TourOptions = {
    ...defaultConfig,
    ...props,
  };
  return new Shepherd.Tour(newProps);
};

export { shepherd };
