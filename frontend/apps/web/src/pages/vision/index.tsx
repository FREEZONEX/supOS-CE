import { type FC, useCallback } from 'react';
import { useSearchParams } from 'react-router';
import { Tabs } from 'antd';
import { Cctv } from 'lucide-react';
import { useTranslate } from '@/hooks';
import CameraInputsPanel from './CameraInputsPanel';
import AlgorithmsPanel from './AlgorithmsPanel';
import AnalysisTasksPanel from './AnalysisTasksPanel';
import ResultsPanel from './ResultsPanel';
import styles from './index.module.scss';

type VisionTab = 'camera' | 'algorithm' | 'task' | 'result';

const normalizeTab = (value: string | null): VisionTab => {
  if (value === 'camera' || value === 'algorithm' || value === 'task' || value === 'result') {
    return value;
  }
  return 'camera';
};

const VisionPage: FC = () => {
  const formatMessage = useTranslate();
  const [searchParams, setSearchParams] = useSearchParams();
  const tab = normalizeTab(searchParams.get('tab'));

  const setTab = useCallback(
    (nextTab: string) => {
      const next = new URLSearchParams(searchParams);
      next.set('tab', nextTab);
      setSearchParams(next, { replace: true });
    },
    [searchParams, setSearchParams]
  );

  return (
    <div className={styles.visionPage}>
      <div className={styles.visionHeader}>
        <div className={styles.pageHeader}>
          <Cctv size={20} strokeWidth={1.75} />
          <span className={styles.pageTitle}>{formatMessage('Vision.title')}</span>
        </div>
        <Tabs
          className={styles.visionTabs}
          activeKey={tab}
          onChange={setTab}
          items={[
            { key: 'camera', label: formatMessage('Vision.camera.title') },
            { key: 'algorithm', label: formatMessage('Vision.algorithm.title') },
            { key: 'task', label: formatMessage('Vision.task.title') },
            { key: 'result', label: formatMessage('Vision.results.title') },
          ]}
        />
      </div>
      <div className={styles.panelHost}>
        {tab === 'camera' && <CameraInputsPanel />}
        {tab === 'algorithm' && <AlgorithmsPanel />}
        {tab === 'task' && <AnalysisTasksPanel />}
        {tab === 'result' && <ResultsPanel />}
      </div>
    </div>
  );
};

export default VisionPage;
