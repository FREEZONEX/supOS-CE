import { createWithEqualityFn, type UseBoundStoreWithEqualityFn } from 'zustand/traditional';
import type { StoreApi } from 'zustand';
import { shallow } from 'zustand/vanilla/shallow';
import { storageOpt } from '@/utils';
import {
  SUPOS_PRIMARY_COLOR,
  SUPOS_REAL_THEME,
  SUPOS_STORAGE_MENU_TYPE,
  SUPOS_THEME,
} from '@/common-types/constans.ts';

export enum MenuTypeEnum {
  Fixed = 'fixed',
  Top = 'top',
}

export enum ThemeType {
  // 主题命名一定要以 light/dark-主题色为准
  Light = 'light',
  Dark = 'dark',
  System = 'system',
}

export enum PrimaryColorType {
  Blue = 'blue',
  Chartreuse = 'chartreuse',
}

export type MenuTypeProps = MenuTypeEnum.Fixed | MenuTypeEnum.Top;

export type TThemeStore = {
  // 主题 light dark
  theme: string;
  // 'blue' | 'chartreuse'
  primaryColor: string;
  // 真实的 light dark system
  _theme: string;
  menuType: MenuTypeProps;
  isTop: boolean;
};

// 设置跟节点类名
const setThemeRoot = (theme: string, primaryColor: string) => {
  const root = document.documentElement;
  switch (`${theme}-${primaryColor}`) {
    case 'dark-blue':
      {
        root.classList.remove('chartreuse', 'chartreuseDark');
        root.classList.add('dark');
      }
      break;
    case 'light-blue':
      {
        root.classList.remove('dark', 'chartreuse', 'chartreuseDark');
      }
      break;
    case 'dark-chartreuse':
      {
        root.classList.add('chartreuse', 'chartreuseDark', 'dark');
      }
      break;
    case 'light-chartreuse':
      {
        root.classList.remove('dark', 'chartreuseDark');
        root.classList.add('chartreuse');
      }
      break;
    default:
      {
        root.classList.remove('dark', 'chartreuse', 'chartreuseDark');
      }
      break;
  }
};

// chat2db主题
/**
 * primary-color: polar-blue,polar-green
 * theme light dark  darkDimmed
 * */
const setCha2dbTheme = (theme: string = ThemeType.Light, primaryColor: string = PrimaryColorType.Chartreuse) => {
  const _primaryColor = primaryColor === PrimaryColorType.Blue ? 'polar-blue' : 'polar-green';
  storageOpt.setOrigin('theme', theme);
  storageOpt.setOrigin('primary-color', _primaryColor);
};

export const useThemeStore: UseBoundStoreWithEqualityFn<StoreApi<TThemeStore>> = createWithEqualityFn(() => {
  const theme = storageOpt.getOrigin(SUPOS_THEME) || ThemeType.Light;
  const primaryColor = storageOpt.getOrigin(SUPOS_PRIMARY_COLOR) || PrimaryColorType.Chartreuse;
  const menuType = storageOpt.get(SUPOS_STORAGE_MENU_TYPE) || MenuTypeEnum.Top;
  setCha2dbTheme(theme);
  // 主题初始化
  setThemeRoot(theme, primaryColor);
  return {
    primaryColor,
    menuType,
    theme,
    _theme: storageOpt.getOrigin(SUPOS_REAL_THEME) || ThemeType.Light,
    isTop: menuType === MenuTypeEnum.Top,
  };
}, shallow);

// 设置菜单模式
export const setMenuType = (menuType: MenuTypeProps = MenuTypeEnum.Fixed) => {
  storageOpt.set(SUPOS_STORAGE_MENU_TYPE, menuType);
  useThemeStore.setState({
    menuType,
  });
};

// 设置主题
export const setTheme = (newTheme: ThemeType = ThemeType.Light) => {
  const oldPrimaryColor = useThemeStore.getState().primaryColor;
  storageOpt.setOrigin(SUPOS_REAL_THEME, newTheme);
  if (newTheme === ThemeType.System) {
    const theme = window.matchMedia('(prefers-color-scheme: dark)')?.matches ? ThemeType.Dark : ThemeType.Light;
    storageOpt.setOrigin(SUPOS_THEME, theme);
    storageOpt.setOrigin('dark-mode', theme === ThemeType.Dark ? 'on' : 'off');
    useThemeStore.setState({
      theme,
      _theme: newTheme,
    });
    setCha2dbTheme(theme, oldPrimaryColor);
    setThemeRoot(theme, oldPrimaryColor);
  } else {
    storageOpt.setOrigin(SUPOS_THEME, newTheme);
    storageOpt.setOrigin('dark-mode', newTheme === ThemeType.Dark ? 'on' : 'off');
    useThemeStore.setState({
      theme: newTheme,
      _theme: newTheme,
    });
    setCha2dbTheme(newTheme, oldPrimaryColor);
    setThemeRoot(newTheme, oldPrimaryColor);
  }
};

// 设置主题色
export const setPrimaryColor = (newPrimaryColor: PrimaryColorType = PrimaryColorType.Chartreuse) => {
  const oldTheme = useThemeStore.getState().theme;
  storageOpt.setOrigin(SUPOS_PRIMARY_COLOR, newPrimaryColor);
  useThemeStore.setState({
    primaryColor: newPrimaryColor,
  });
  setCha2dbTheme(oldTheme, newPrimaryColor);
  setThemeRoot(oldTheme, newPrimaryColor);
};

// 系统模式变化 设置 主题
export const setThemeBySystem = (isDark: boolean) => {
  const { _theme, primaryColor } = useThemeStore.getState();
  if (_theme === 'system') {
    storageOpt.setOrigin('dark-mode', isDark ? 'on' : 'off');
    const theme = isDark ? ThemeType.Dark : ThemeType.Light;
    storageOpt.setOrigin(SUPOS_THEME, theme);
    useThemeStore.setState({
      theme,
    });
    setCha2dbTheme(theme, primaryColor);
    setThemeRoot(theme, primaryColor);
  }
};
