import { useState, type FC } from 'react';
import ComContent from '@/components/com-layout/ComContent';
import ComLayout from '@/components/com-layout';
import ComSelect from '@/components/com-select';
import HelpTooltip from '@/components/help-tooltip';
import { useTranslate } from '@/hooks';
import styles from './index.module.scss';

type SessionDurationValue = '12h' | '1d' | '7d' | '30d' | 'permanent';

type SystemConfigurationProps = {
  title: string;
};

const SystemConfiguration: FC<SystemConfigurationProps> = ({ title }) => {
  const formatMessage = useTranslate();
  const [sessionDuration, setSessionDuration] = useState<SessionDurationValue>('12h');
  const sessionDurationOptions = [
    { value: '12h', label: formatMessage('settings.session12Hours') },
    { value: '1d', label: formatMessage('settings.session1Day') },
    { value: '7d', label: formatMessage('settings.session7Days') },
    { value: '30d', label: formatMessage('settings.session30Days') },
    { value: 'permanent', label: formatMessage('settings.sessionPermanent') },
  ];

  return (
    <ComLayout>
      <ComContent title={title} hasBack={false} className={styles['system-configuration']}>
        <section className={styles.section}>
          <h2 className={styles['section-title']}>{formatMessage('settings.login')}</h2>
          <div className={styles.row}>
            <div className={styles.label}>
              <span>{formatMessage('settings.sessionValidity')}</span>
              <HelpTooltip title={formatMessage('settings.sessionGlobalDesc')} />
            </div>
            <ComSelect
              className={styles.select}
              value={sessionDuration}
              options={sessionDurationOptions}
              aria-label={formatMessage('settings.sessionValidity')}
              onChange={(value) => setSessionDuration(value as SessionDurationValue)}
            />
          </div>
        </section>
      </ComContent>
    </ComLayout>
  );
};

export default SystemConfiguration;
