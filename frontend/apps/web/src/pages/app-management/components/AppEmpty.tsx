import type { FC } from 'react';
import { Subdirectory } from '@carbon/icons-react';
import { useNavigate } from 'react-router';
import styles from './AppEmpty.module.scss';
import { useTranslate } from '@/hooks';

const AppEmpty: FC<any> = ({ appName }) => {
  const navigate = useNavigate();
  const formatMessage = useTranslate();
  return (
    <div className={styles['app-empty']}>
      <span className="text">{formatMessage('appSpace.generateHere')}</span>
      <Subdirectory
        onClick={() => {
          navigate('/app-gui', {
            state: { appName },
          });
        }}
        className="icon"
        size={140}
      />
    </div>
  );
};

export default AppEmpty;
