import type { CSSProperties, FC, ReactNode } from 'react';
import { Flex } from 'antd';
import cx from 'classnames';
import { ChartLine, ClipboardList, Route, SendAlt } from '@/components/lucide-icon/carbon';
import { treeIconProps } from '@/components/lucide-icon/icon-props';
import { topicTypeNumber } from '@/apis/core-api/core-adapter';
import type { UnsTreeNode } from '@/pages/uns/types';
import './tree-icons.scss';

/** Figma Tier0 Namespace tree — node-id 13040:96411 */
export const UNS_TOPIC_ICON_COLORS = {
  state: 'var(--ui-uns-topic-state-icon)',
  action: 'var(--ui-uns-topic-action-icon)',
  metric: 'var(--ui-uns-topic-metric-icon)',
} as const;

export const UNS_TOPIC_FOLDER_BG = {
  state: 'var(--ui-uns-topic-state-bg)',
  action: 'var(--ui-uns-topic-action-bg)',
  metric: 'var(--ui-uns-topic-metric-bg)',
} as const;

export const UNS_TOPIC_FOLDER_TEXT = {
  state: 'var(--ui-uns-topic-state-text)',
  action: 'var(--ui-uns-topic-action-text)',
  metric: 'var(--ui-uns-topic-metric-text)',
} as const;

type UnsTopicTypeKey = keyof typeof UNS_TOPIC_FOLDER_BG;

const TOPIC_TYPE_KEY_MAP: Record<number, UnsTopicTypeKey | undefined> = {
  1: 'state',
  2: 'action',
  3: 'metric',
};

export const getUnsTopicTypeKey = (topicType: number): UnsTopicTypeKey | undefined => TOPIC_TYPE_KEY_MAP[topicType];

export const getUnsTopicFolderBackground = (topicType: number) => {
  const key = getUnsTopicTypeKey(topicType);
  return key ? UNS_TOPIC_FOLDER_BG[key] : undefined;
};

export const getUnsTopicFolderTextColor = (topicType: number) => {
  const key = getUnsTopicTypeKey(topicType);
  return key ? UNS_TOPIC_FOLDER_TEXT[key] : undefined;
};

export const getUnsTopicFolderTitleStyle = (
  topicType: number,
  options?: { height?: number | string; paddingRight?: number | string; paddingLeft?: number | string }
): CSSProperties => {
  const background = getUnsTopicFolderBackground(topicType);
  if (!background) return {};

  return {
    height: options?.height ?? '26px',
    backgroundColor: background,
    borderRadius: '3px',
    paddingLeft: options?.paddingLeft ?? '8px',
    paddingRight: options?.paddingRight ?? '8px',
    color: getUnsTopicFolderTextColor(topicType),
  };
};

export interface UnsTopicTypeIconWrapProps {
  topicType: number;
  size?: number;
  children: ReactNode;
}

export const UnsTopicTypeIconWrap: FC<UnsTopicTypeIconWrapProps> = ({ topicType, size = 36, children }) => {
  const topicKey = getUnsTopicTypeKey(topicType);

  return (
    <Flex
      align="center"
      justify="center"
      className={cx('uns-topic-type-icon-wrap', topicKey && `uns-topic-type-icon-wrap--${topicKey}`)}
      style={{ width: size, height: size }}
    >
      {children}
    </Flex>
  );
};

type TopicType = 1 | 2 | 3;

const TOPIC_FOLDER_NAMES: Record<number, string> = {
  1: 'State',
  2: 'Action',
  3: 'Metric',
};

const getTopicFolderLeafName = (node?: UnsTreeNode | null) => {
  const path = String(node?.path || node?.alias || '').trim();
  const pathLeaf = path.split('/').filter(Boolean).pop();
  return String(pathLeaf || node?.pathName || node?.name || node?.title || '').trim();
};

const getTopicTypeFromLeafName = (leafName: string) => {
  const normalized = leafName.trim().toLowerCase();
  for (const [topicType, folderName] of Object.entries(TOPIC_FOLDER_NAMES)) {
    if (normalized === folderName.toLowerCase()) {
      return Number(topicType);
    }
  }
  return 0;
};

export const isTopicTypeFolder = (node?: UnsTreeNode | null) => {
  if (!node || node.pathType !== 0) {
    return false;
  }

  const leafName = getTopicFolderLeafName(node);
  const topicTypeFromLeaf = getTopicTypeFromLeafName(leafName);
  const dataType = Number(node.dataType || 0);
  if (dataType >= 1 && dataType <= 3) {
    return !topicTypeFromLeaf || topicTypeFromLeaf === dataType;
  }

  const topicTypeValue = topicTypeNumber(node.topicType);
  if (topicTypeValue && leafName.toLowerCase() === TOPIC_FOLDER_NAMES[topicTypeValue].toLowerCase()) {
    return true;
  }

  return topicTypeFromLeaf > 0;
};

export const getFolderTopicType = (node?: UnsTreeNode | null) => {
  if (!isTopicTypeFolder(node)) {
    return 0;
  }

  const leafName = getTopicFolderLeafName(node);
  const dataType = Number(node?.dataType || 0);
  if (dataType >= 1 && dataType <= 3) {
    return dataType;
  }

  const topicTypeValue = topicTypeNumber(node?.topicType);
  if (topicTypeValue) {
    return topicTypeValue;
  }

  return getTopicTypeFromLeafName(leafName);
};

export const getNodeTopicType = (node?: UnsTreeNode | null) =>
  node?.pathType === 0 ? getFolderTopicType(node) : Number(node?.parentDataType || 0);

const topicIconMap: Record<TopicType, { Icon: typeof ClipboardList; color: string }> = {
  1: { Icon: ClipboardList, color: UNS_TOPIC_ICON_COLORS.state },
  2: { Icon: SendAlt, color: UNS_TOPIC_ICON_COLORS.action },
  3: { Icon: ChartLine, color: UNS_TOPIC_ICON_COLORS.metric },
};

export const getTopicIconMeta = (topicType: number) => {
  if (topicType === 1 || topicType === 2 || topicType === 3) {
    return topicIconMap[topicType as TopicType];
  }
  return null;
};

const baseIconStyle: CSSProperties = { flexShrink: 0, marginRight: '5px' };

export interface UnsTreeNodeIconProps {
  dataNode: UnsTreeNode;
  topicType?: number;
  enableAutoCategorization?: boolean;
  isTopicTypeFolder?: (node: UnsTreeNode) => boolean;
  statusPrefix?: ReactNode;
}

export const UnsTreeNodeIcon: FC<UnsTreeNodeIconProps> = ({
  dataNode,
  topicType = 0,
  enableAutoCategorization = false,
  isTopicTypeFolder,
  statusPrefix,
}) => {
  const meta = getTopicIconMeta(topicType);
  const iconStyle = baseIconStyle;

  const renderTopicIcon = () => {
    const compactIconStyle = { ...iconStyle, marginRight: 0 };
    if (!meta) {
      return <Route {...treeIconProps} style={compactIconStyle} />;
    }
    const { Icon, color } = meta;
    return <Icon {...treeIconProps} style={{ ...compactIconStyle, color }} />;
  };

  if (dataNode.pathType === 0) {
    const topicFolderIcon = enableAutoCategorization && isTopicTypeFolder?.(dataNode);

    return (
      <Flex align="center">
        {statusPrefix ? <div style={{ width: 10, display: 'flex', alignItems: 'center' }}>{statusPrefix}</div> : null}
        {topicFolderIcon ? renderTopicIcon() : <Route {...treeIconProps} style={iconStyle} />}
      </Flex>
    );
  }

  if (dataNode.pathType === 2) {
    return (
      <Flex align="center">
        <div style={{ width: 10 }} />
        {meta ? (
          <meta.Icon {...treeIconProps} style={{ ...iconStyle, color: meta.color }} />
        ) : (
          <ClipboardList {...treeIconProps} style={iconStyle} />
        )}
      </Flex>
    );
  }

  return null;
};
