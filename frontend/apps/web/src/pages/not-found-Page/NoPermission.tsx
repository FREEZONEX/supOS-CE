import { Button, Result } from 'antd';
import { useTranslate } from '@/hooks';
import useNavigateForIframe from '@/hooks/useNavigateForIframe';
import { useBaseStore } from '@/stores/base';

const NotFoundPage = () => {
  const formatMessage = useTranslate();
  const homePage = useBaseStore((state) => state.currentUserInfo?.homePage) || '/uns';
  const { security, onClick } = useNavigateForIframe({ path: homePage, replaceCurrent: true });

  return (
    <Result
      status="403"
      title={403}
      subTitle={<span style={{ color: 'var(--supos-text-color)' }}>{formatMessage('common.pageNoPermission')}</span>}
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
