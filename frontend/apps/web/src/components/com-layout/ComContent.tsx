import type { CSSProperties, FC, ReactNode } from 'react';
import { Button, Tooltip } from 'antd';
import { useNavigate } from 'react-router';
import { DrillBack } from '@carbon/icons-react';
import classNames from 'classnames';
import { useTranslate } from '@/hooks';
import styles from './index.module.scss';
import { useThemeStore } from '@/stores/theme-store.ts';

export interface ComContentProps {
  border?: boolean;
  hasPadding?: boolean;
  children?: ReactNode;
  title?: ReactNode;
  extra?: ReactNode;
  style?: CSSProperties;
  hasBack?: boolean;
  backPath?: any;
  wrapperStyle?: CSSProperties;
  titleStyle?: CSSProperties;
  className?: string;
  mustShowTitle?: boolean;
  mustHasBack?: boolean;
}

const ComContent: FC<ComContentProps> = ({
  border = true,
  hasPadding,
  children,
  title,
  extra,
  style,
  hasBack = true,
  mustShowTitle = true,
  mustHasBack = true,
  backPath = -1,
  wrapperStyle,
  titleStyle,
  className,
}) => {
  const navigate = useNavigate();
  const isTop = useThemeStore((state) => state.isTop);
  const noTitle = mustShowTitle || !isTop;
  const noBack = mustHasBack || !isTop;
  const formatMessage = useTranslate();

  return (
    <div className={classNames(styles['com-content'], className)} style={wrapperStyle}>
      {noTitle && title && (
        <div
          style={{ paddingLeft: hasPadding && !isTop ? 300 : 25, ...(border ? {} : { border: 'none' }), ...titleStyle }}
          className="title"
        >
          <div
            style={{ flex: 1, overflow: 'hidden', whiteSpace: 'nowrap', textOverflow: 'ellipsis' }}
            title={typeof title === 'string' ? title : undefined}
          >
            {title}
          </div>
          {extra}
          {noBack && hasBack && (
            <Tooltip placement="bottom" title={formatMessage('common.back')}>
              <Button onClick={() => navigate(backPath)} type="primary" icon={<DrillBack style={{ marginTop: 2 }} />}>
                {formatMessage('common.back')}
              </Button>
            </Tooltip>
          )}
        </div>
      )}
      <div className="content" style={style}>
        {children}
      </div>
    </div>
  );
};

export default ComContent;
