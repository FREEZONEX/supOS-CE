/**
 * @description 模块联邦配置
 * */
export const mfConfig = {
  name: 'supos-ce/host',
  manifest: true,
  exposes: {
    './components': './src/components/index.ts',
    './utils': './src/utils/index.ts',
    './hooks': './src/hooks/index.ts',
    './apis': './src/apis/inter-api/index.ts',
    './button-permission': './src/common-types/button-permission.ts',
    './constans': './src/common-types/constans.ts',
    './i18nStore': './src/stores/i18n-store.ts',
    './baseStore': './src/stores/base/index.ts',
    './tabs-lifecycle-context': './src/contexts/tabs-lifecycle-context.ts',
    './useTabsContext': './src/contexts/tabs-context.ts',
  },
  shared: {
    react: {
      singleton: true,
      requiredVersion: '18.3.1',
    },
    'react-dom': {
      singleton: true,
      requiredVersion: '18.3.1',
    },
    'react-router': {
      singleton: true,
      requiredVersion: '7.9.4',
    },
    antd: {
      singleton: true,
      requiredVersion: '5.27.4',
    },
    '@ant-design/icons': {
      singleton: true,
      requiredVersion: '6.1.0',
    },
    ahooks: {
      singleton: true,
      requiredVersion: '3.9.5',
    },
    '@carbon/icons-react': {
      singleton: true,
      requiredVersion: '11.69.0',
    },
    'lodash-es': {
      singleton: true,
      requiredVersion: '4.17.21',
    },
    sass: {
      singleton: true,
      requiredVersion: '1.93.2',
    },
  },
};
