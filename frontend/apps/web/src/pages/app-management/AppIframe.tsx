import type { PageProps } from '@/common-types';
import IframeWrapper from '@/components/iframe-wrapper';
import type { FC } from 'react';

const AppIframe: FC<PageProps> = ({ location }) => {
  const { state } = location || {};

  return <IframeWrapper title={state?.title} src={state?.src} />;
};

export default AppIframe;
