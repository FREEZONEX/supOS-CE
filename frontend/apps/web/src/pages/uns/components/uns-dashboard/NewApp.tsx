import { Flex } from 'antd';
import { AddLarge } from '@carbon/icons-react';
import styles from './index.module.scss';
import useTranslate from '../../../../hooks/useTranslate.ts';

const NewApp = () => {
  const formatMessage = useTranslate();

  return (
    <Flex className={styles['new-app']} justify="center" vertical align="center" style={{ height: '100%' }} gap={8}>
      <AddLarge size={32} />
      {formatMessage('uns.deployNewApp')}
    </Flex>
  );
};

export default NewApp;
