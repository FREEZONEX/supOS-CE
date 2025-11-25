import localZhCN from '@/locale/zh-CN.json';
import localEnUS from '@/locale/en-US.json';
import { getSystemI18Api } from '@/apis/inter-api/uns.ts';

type I18n = 'zh-CN' | 'en-US';
export type I18nData = { [x: string]: string };

const localSources: { [x: string]: any } = {
  'zh-CN': localZhCN,
  'en-US': localEnUS,
};

// antd 所有语言包
const antSources: Record<string, () => Promise<any>> = {
  ar_EG: () => import('antd/es/locale/ar_EG'),
  az_AZ: () => import('antd/es/locale/az_AZ'),
  bg_BG: () => import('antd/es/locale/bg_BG'),
  bn_BD: () => import('antd/es/locale/bn_BD'),
  by_BY: () => import('antd/es/locale/by_BY'),
  ca_ES: () => import('antd/es/locale/ca_ES'),
  cs_CZ: () => import('antd/es/locale/cs_CZ'),
  da_DK: () => import('antd/es/locale/da_DK'),
  de_DE: () => import('antd/es/locale/de_DE'),
  el_GR: () => import('antd/es/locale/el_GR'),
  en_GB: () => import('antd/es/locale/en_GB'),
  en_US: () => import('antd/es/locale/en_US'),
  es_ES: () => import('antd/es/locale/es_ES'),
  et_EE: () => import('antd/es/locale/et_EE'),
  fa_IR: () => import('antd/es/locale/fa_IR'),
  fi_FI: () => import('antd/es/locale/fi_FI'),
  fr_BE: () => import('antd/es/locale/fr_BE'),
  fr_CA: () => import('antd/es/locale/fr_CA'),
  fr_FR: () => import('antd/es/locale/fr_FR'),
  ga_IE: () => import('antd/es/locale/ga_IE'),
  gl_ES: () => import('antd/es/locale/gl_ES'),
  he_IL: () => import('antd/es/locale/he_IL'),
  hi_IN: () => import('antd/es/locale/hi_IN'),
  hr_HR: () => import('antd/es/locale/hr_HR'),
  hu_HU: () => import('antd/es/locale/hu_HU'),
  hy_AM: () => import('antd/es/locale/hy_AM'),
  id_ID: () => import('antd/es/locale/id_ID'),
  is_IS: () => import('antd/es/locale/is_IS'),
  it_IT: () => import('antd/es/locale/it_IT'),
  ja_JP: () => import('antd/es/locale/ja_JP'),
  ka_GE: () => import('antd/es/locale/ka_GE'),
  kk_KZ: () => import('antd/es/locale/kk_KZ'),
  km_KH: () => import('antd/es/locale/km_KH'),
  kmr_IQ: () => import('antd/es/locale/kmr_IQ'),
  kn_IN: () => import('antd/es/locale/kn_IN'),
  ko_KR: () => import('antd/es/locale/ko_KR'),
  ku_IQ: () => import('antd/es/locale/ku_IQ'),
  lt_LT: () => import('antd/es/locale/lt_LT'),
  lv_LV: () => import('antd/es/locale/lv_LV'),
  mk_MK: () => import('antd/es/locale/mk_MK'),
  ml_IN: () => import('antd/es/locale/ml_IN'),
  mn_MN: () => import('antd/es/locale/mn_MN'),
  ms_MY: () => import('antd/es/locale/ms_MY'),
  nb_NO: () => import('antd/es/locale/nb_NO'),
  ne_NP: () => import('antd/es/locale/ne_NP'),
  nl_BE: () => import('antd/es/locale/nl_BE'),
  nl_NL: () => import('antd/es/locale/nl_NL'),
  pl_PL: () => import('antd/es/locale/pl_PL'),
  pt_BR: () => import('antd/es/locale/pt_BR'),
  pt_PT: () => import('antd/es/locale/pt_PT'),
  ro_RO: () => import('antd/es/locale/ro_RO'),
  ru_RU: () => import('antd/es/locale/ru_RU'),
  si_LK: () => import('antd/es/locale/si_LK'),
  sk_SK: () => import('antd/es/locale/sk_SK'),
  sl_SI: () => import('antd/es/locale/sl_SI'),
  sr_RS: () => import('antd/es/locale/sr_RS'),
  sv_SE: () => import('antd/es/locale/sv_SE'),
  ta_IN: () => import('antd/es/locale/ta_IN'),
  th_TH: () => import('antd/es/locale/th_TH'),
  tk_TK: () => import('antd/es/locale/tk_TK'),
  tr_TR: () => import('antd/es/locale/tr_TR'),
  uk_UA: () => import('antd/es/locale/uk_UA'),
  ur_PK: () => import('antd/es/locale/ur_PK'),
  ur_UZ: () => import('antd/es/locale/uz_UZ.js'),
  vi_VN: () => import('antd/es/locale/vi_VN'),
  zh_CN: () => import('antd/es/locale/zh_CN'),
  zh_HK: () => import('antd/es/locale/zh_HK'),
  zh_TW: () => import('antd/es/locale/zh_TW'),
};

export const normalizedLangFn = (locale: string) => locale.replace('-', '_');
// 动态加载antd语言包
export const loadAntdLocale = async (locale: string) => {
  const loader = antSources[normalizedLangFn(locale)] || antSources['en_US']; // 默认英语
  const module = await loader();
  return module.default as I18nData; // 返回默认导出
};

// 动态加载语言包并转换为react-intl期望的格式
export const loadMessages = async (lang: I18n) => {
  try {
    // 加载本地语言包
    const localMessages = localSources[lang];

    // 尝试从 服务器加载语言包
    let backEndMessages = {};
    try {
      const content = await getSystemI18Api(lang);
      backEndMessages = content?.messages || {};
      // backEndMessages = content;
    } catch (e) {
      console.log(e);
    }
    // 合并语言包，以后端服务器存储的为准
    const messages = { ...localMessages, ...backEndMessages };
    // 如果两个来源都没有加载到语言包，则抛出错误
    if (Object.keys(messages).length === 0) {
      throw new Error(`Failed to load any language file for ${lang}`);
    }

    const data: I18nData = messages; // 初始值可以是默认语言包

    return data;
  } catch (error) {
    console.error(`Error loading language file for ${lang}:`, error);
    return {};
  }
};
