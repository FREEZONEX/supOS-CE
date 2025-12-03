import useRemote from '@/hooks/useRemote';
import { useMatchRoute } from '@/hooks/useMatchRouter';
import type { FC } from 'react';
import type { PageProps } from '@/common-types';
import { Button, Result, Typography } from 'antd';
import { useTranslate } from '@/hooks';
import ErrorBoundary from '@/components/error-boundary';
const { Paragraph } = Typography;

const DynamicMFComponent: FC<PageProps> = (props) => {
  const matchRoute: any = useMatchRoute();
  const formatMessage = useTranslate();
  const {
    Module: Comp,
    reLoadRemote,
    errorMsg,
  } = useRemote({
    name: matchRoute?.moduleName ? matchRoute?.parentPath : matchRoute?.path,
    moduleName: matchRoute?.moduleName,
    location: props?.location,
  });
  if (errorMsg || Comp.error) {
    return (
      <Result
        status="error"
        title={formatMessage('plugin.errorTitle')}
        extra={[
          <Button
            type="primary"
            key="console"
            onClick={() => {
              reLoadRemote();
            }}
          >
            {formatMessage('plugin.retry')}
          </Button>,
        ]}
      >
        <Paragraph>{errorMsg || Comp?.error?.toString()}</Paragraph>
      </Result>
    );
  }

  return (
    <ErrorBoundary>
      <Comp {...props} />
    </ErrorBoundary>
  );
};

export default DynamicMFComponent;
