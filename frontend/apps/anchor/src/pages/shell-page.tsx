import { Empty } from 'antd';
import { getLocale, t } from '../i18n';

const MODULE_TITLES = {
  model: { zh: '模型库', en: 'Model Library' },
  scene: { zh: '场景', en: 'Scenes' },
} as const;

export function ShellPage({ module }: { module: keyof typeof MODULE_TITLES }) {
  const locale = getLocale();
  const title = MODULE_TITLES[module][locale === 'zh-CN' ? 'zh' : 'en'];
  return (
    <div className="flex h-full flex-col">
      <div className="page-title-bar">
        <div className="flex-1 truncate">{title}</div>
      </div>
      <div className="flex min-h-0 flex-1 items-center justify-center">
        <Empty description={t('shell.migrating')} />
      </div>
    </div>
  );
}
