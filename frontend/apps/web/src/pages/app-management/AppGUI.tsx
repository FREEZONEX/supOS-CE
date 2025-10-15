import { Flex } from 'antd';
import { type FC, useRef, useState } from 'react';
import DeployForm from '@/pages/app-management/components/DeployForm';
import type { PageProps } from '@/common-types';
import { useTranslate } from '@/hooks';
import { ButtonPermission } from '@/common-types/button-permission';
import { AuthButton } from '@/components/auth';
import ComDrawer from '@/components/com-drawer';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import Board from '@/components/craft';
const Module: FC<PageProps> = ({ location }) => {
  const formatMessage = useTranslate();

  const { state } = location || {};
  const [show, setShow] = useState(false);
  const boardCodeRef = useRef(null);
  const getHtmlContent = () => {
    return boardCodeRef.current;
  };
  return (
    <ComLayout>
      <ComContent
        title={
          <Flex justify="flex-end" align="center" style={{ height: '100%' }}>
            <AuthButton
              auth={ButtonPermission['appGui.deploy']}
              style={{ width: 102 }}
              type="primary"
              onClick={() => setShow((pre) => !pre)}
            >
              {formatMessage('appGui.deploy')}
            </AuthButton>
          </Flex>
        }
      >
        <Board boardCodeRef={boardCodeRef} />
      </ComContent>
      <ComDrawer title=" " open={show} onClose={() => setShow(false)}>
        <DeployForm show={show} setShow={setShow} getHtmlContent={getHtmlContent} appName={state?.appName} />
      </ComDrawer>
    </ComLayout>
  );
};

export default Module;
