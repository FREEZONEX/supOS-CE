import ComEmpty from '@/components/com-empty';
import { useTranslate } from '@/hooks';
import styles from './index.module.scss';

const EmptyDetail = () => {
  const formatMessage = useTranslate();

  return (
    <div className={styles['emptyDetail-wrap']}>
      <ComEmpty description={formatMessage('MenuConfiguration.selectResource')} />
    </div>
  );
};

export default EmptyDetail;
