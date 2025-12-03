import type { CSSProperties, FC, ReactNode } from 'react';
import ComCopy from '../com-copy/index';
import { Flex } from 'antd';
import './ComCopyContent.scss';
import classNames from 'classnames';

const ComCopyContent: FC<{
  className?: string;
  textToCopy?: string;
  label?: ReactNode;
  labelClassName?: string;
  style?: CSSProperties;
}> = ({ textToCopy = ' ', className, label, labelClassName, style }) => {
  return (
    <Flex align="center" className={classNames('com-copy-content', className)}>
      {label && <div className={classNames('label', labelClassName)}>{label}</div>}
      <Flex className={'content'} justify="space-between" style={style}>
        <div className={'text'} title={textToCopy}>
          {textToCopy}
        </div>
        <ComCopy textToCopy={textToCopy} />
      </Flex>
    </Flex>
  );
};

export default ComCopyContent;
