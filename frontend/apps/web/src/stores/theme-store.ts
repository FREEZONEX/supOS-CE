import { createWithEqualityFn, type UseBoundStoreWithEqualityFn } from 'zustand/traditional';
import type { StoreApi } from 'zustand';
import { shallow } from 'zustand/vanilla/shallow';
import { storageOpt } from '@/utils';
import { APP_PRIMARY_COLOR, APP_REAL_THEME, APP_STORAGE_MENU_TYPE, APP_THEME } from '@/common-types/constans.ts';

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

const DEFAULT_PRIMARY_COLOR = PrimaryColorType.Chartreuse;
const PRIMARY_COLOR_DEFAULT_VERSION_KEY = 'APP_PRIMARY_COLOR_DEFAULT_VERSION';
const PRIMARY_COLOR_DEFAULT_VERSION = 'chartreuse-default-v1';

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
const setCha2dbTheme = (theme: string = ThemeType.Light, primaryColor: string = DEFAULT_PRIMARY_COLOR) => {
  const _primaryColor = primaryColor === PrimaryColorType.Blue ? 'polar-blue' : 'polar-green';
  storageOpt.setOrigin('theme', theme);
  storageOpt.setOrigin('primary-color', _primaryColor);
};

const resolvePrimaryColor = () => {
  storageOpt.setOrigin(APP_PRIMARY_COLOR, DEFAULT_PRIMARY_COLOR);
  storageOpt.setOrigin(PRIMARY_COLOR_DEFAULT_VERSION_KEY, PRIMARY_COLOR_DEFAULT_VERSION);
  return DEFAULT_PRIMARY_COLOR;
};

export const useThemeStore: UseBoundStoreWithEqualityFn<StoreApi<TThemeStore>> = createWithEqualityFn(() => {
  const theme = storageOpt.getOrigin(APP_THEME) || ThemeType.Light;
  const primaryColor = resolvePrimaryColor();
  const menuType = storageOpt.get(APP_STORAGE_MENU_TYPE) || MenuTypeEnum.Top;
  setCha2dbTheme(theme, primaryColor);
  // 主题初始化
  setThemeRoot(theme, primaryColor);
  return {
    primaryColor,
    menuType,
    theme,
    _theme: storageOpt.getOrigin(APP_REAL_THEME) || ThemeType.Light,
    isTop: menuType === MenuTypeEnum.Top,
  };
}, shallow);

export const restoreThemeRootFromStore = () => {
  const { theme, primaryColor } = useThemeStore.getState();
  setThemeRoot(theme, primaryColor);
};

// 设置菜单模式
export const setMenuType = (menuType: MenuTypeProps = MenuTypeEnum.Fixed) => {
  storageOpt.set(APP_STORAGE_MENU_TYPE, menuType);
  useThemeStore.setState({
    menuType,
  });
};

// 设置主题
export const setTheme = (newTheme: ThemeType = ThemeType.Light) => {
  const oldPrimaryColor = useThemeStore.getState().primaryColor;
  storageOpt.setOrigin(APP_REAL_THEME, newTheme);
  if (newTheme === ThemeType.System) {
    const theme = window.matchMedia('(prefers-color-scheme: dark)')?.matches ? ThemeType.Dark : ThemeType.Light;
    storageOpt.setOrigin(APP_THEME, theme);
    storageOpt.setOrigin('dark-mode', theme === ThemeType.Dark ? 'on' : 'off');
    useThemeStore.setState({
      theme,
      _theme: newTheme,
    });
    setCha2dbTheme(theme, oldPrimaryColor);
    setThemeRoot(theme, oldPrimaryColor);
  } else {
    storageOpt.setOrigin(APP_THEME, newTheme);
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
export const setPrimaryColor = () => {
  const oldTheme = useThemeStore.getState().theme;
  storageOpt.setOrigin(APP_PRIMARY_COLOR, DEFAULT_PRIMARY_COLOR);
  storageOpt.setOrigin(PRIMARY_COLOR_DEFAULT_VERSION_KEY, PRIMARY_COLOR_DEFAULT_VERSION);
  useThemeStore.setState({
    primaryColor: DEFAULT_PRIMARY_COLOR,
  });
  setCha2dbTheme(oldTheme, DEFAULT_PRIMARY_COLOR);
  setThemeRoot(oldTheme, DEFAULT_PRIMARY_COLOR);
};

// 系统模式变化 设置 主题
export const setThemeBySystem = (isDark: boolean) => {
  const { _theme, primaryColor } = useThemeStore.getState();
  if (_theme === 'system') {
    storageOpt.setOrigin('dark-mode', isDark ? 'on' : 'off');
    const theme = isDark ? ThemeType.Dark : ThemeType.Light;
    storageOpt.setOrigin(APP_THEME, theme);
    useThemeStore.setState({
      theme,
    });
    setCha2dbTheme(theme, primaryColor);
    setThemeRoot(theme, primaryColor);
  }
};
