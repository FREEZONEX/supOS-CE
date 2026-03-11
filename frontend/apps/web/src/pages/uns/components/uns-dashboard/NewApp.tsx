import { Flex } from 'antd';
import { AddLarge } from '@carbon/icons-react';
import styles from './index.module.scss';
import useTranslate from '@/hooks/useTranslate.ts';
import AddAppModal from './AddAppModal.tsx';
import { useRef } from 'react';

const NewApp = () => {
  const formatMessage = useTranslate();
  const addAppRef = useRef<any>(null);

  return (
    <>
      <Flex
        className={styles['new-app']}
        justify="center"
        vertical
        align="center"
        style={{ height: '100%' }}
        gap={8}
        onClick={() => {
          addAppRef.current?.onOpen();
        }}
      >
        <AddLarge size={32} />
        {formatMessage('uns.deployNewApp')}
      </Flex>
      <AddAppModal ref={addAppRef} />
    </>
  );
};

export default NewApp;
