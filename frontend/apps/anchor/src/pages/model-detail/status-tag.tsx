import { Tag } from 'antd';
import { ModelStatus } from '../../api/models';
import { t } from '../../i18n';

// 对齐源端状态配色：Unknown 灰 / Parsing 黄 / Ready 蓝 / Error 红 / Running 绿
const STATUS_CONFIG: Record<number, { color: string; key: string }> = {
  [ModelStatus.Unknown]: { color: 'default', key: 'instance.status.unknown' },
  [ModelStatus.Parsing]: { color: 'gold', key: 'instance.status.parsing' },
  [ModelStatus.Ready]: { color: 'blue', key: 'instance.status.ready' },
  [ModelStatus.Error]: { color: 'red', key: 'instance.status.error' },
  4: { color: 'green', key: 'instance.status.running' },
};

export function StatusTag({ status, running }: { status: number; running?: boolean }) {
  const effective = running ? 4 : status;
  const config = STATUS_CONFIG[effective] ?? STATUS_CONFIG[ModelStatus.Unknown];
  return <Tag color={config.color}>{t(config.key)}</Tag>;
}
