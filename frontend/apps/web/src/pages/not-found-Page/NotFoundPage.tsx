import { Button, Result } from 'antd';
import { useTranslate } from '@/hooks';
import useNavigateForIframe from '@/hooks/useNavigateForIframe';
import { useBaseStore } from '@/stores/base';

const NotFoundPage = () => {
  const formatMessage = useTranslate();
  const homePage = useBaseStore((state) => state.currentUserInfo?.homePage) || '/uns';
  const { security, onClick } = useNavigateForIframe({ path: homePage });

  return (
    <Result
      status="404"
      title={<span style={{ color: 'var(--supos-text-color)' }}>{formatMessage('common.notFound')}</span>}
      subTitle={<span style={{ color: 'var(--supos-text-color)' }}>{formatMessage('common.pageNotFound')}</span>}
      style={{ backgroundColor: 'var(--supos-bg-color)' }}
      extra={
        security && (
          <Button type="primary" onClick={onClick}>
            {formatMessage('common.goHome')}
          </Button>
        )
      }
    />
  );
};

export default NotFoundPage;
