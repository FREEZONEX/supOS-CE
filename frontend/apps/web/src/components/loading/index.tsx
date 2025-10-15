import { Spin } from 'antd';
import styles from './index.module.scss';
import type { FC } from 'react';

const Loading: FC<any> = ({ children, spinning, ...restProps }) => {
  return (
    <Spin spinning={spinning} wrapperClassName={styles['login-wrapper-loading']} {...restProps}>
      {children}
    </Spin>
  );
};

export default Loading;
