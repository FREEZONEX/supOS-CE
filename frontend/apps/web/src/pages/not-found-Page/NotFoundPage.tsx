import { Button, Result } from 'antd';
import { useTranslate } from '@/hooks';
import useNavigateForIframe from '@/hooks/useNavigateForIframe';

const NotFoundPage = () => {
  const formatMessage = useTranslate();
  const { security, onClick } = useNavigateForIframe({ path: '/uns' });

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
