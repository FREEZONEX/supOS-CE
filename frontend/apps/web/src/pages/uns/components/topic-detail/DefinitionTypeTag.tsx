import type { FC } from 'react';
import { Tag } from 'antd';
import classNames from 'classnames';
import styles from './Definition.module.scss';

const TYPE_CLASS_MAP: Record<string, string> = {
  STRING: 'string',
  BOOLEAN: 'boolean',
};

const DefinitionTypeTag: FC<{ type?: string }> = ({ type }) => {
  const normalized = String(type || '')
    .trim()
    .toUpperCase();
  const kind = TYPE_CLASS_MAP[normalized] || 'default';

  return (
    <Tag bordered={false} className={classNames(styles.typePill, styles[`typePill--${kind}`])}>
      {normalized || '-'}
    </Tag>
  );
};

export default DefinitionTypeTag;
