import { useEffect, useState } from 'react';
import { Alert, Button, Empty, Spin, Tag } from 'antd';
import { testVisionTaskFrame, type VisionTask, type VisionTaskTestFrameResult } from '@/apis/core-api/task';
import ProModal from '@/components/pro-modal';
import { useTranslate } from '@/hooks';
import styles from './index.module.scss';

type TaskTestFrameModalProps = {
  task: VisionTask | null;
  onClose: () => void;
};

// 从一条检测记录里尽量解析出标签与置信度(worker 返回的字段名可能不同)。
const readLabel = (obj: Record<string, unknown>) => String(obj.label ?? obj.class ?? obj.name ?? obj.category ?? '-');
const readScore = (obj: Record<string, unknown>) => {
  const raw = obj.confidence ?? obj.score ?? obj.conf;
  const num = typeof raw === 'number' ? raw : Number(raw);
  return Number.isFinite(num) ? num : null;
};

const TaskTestFrameModal = ({ task, onClose }: TaskTestFrameModalProps) => {
  const formatMessage = useTranslate();
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<VisionTaskTestFrameResult | null>(null);
  const [error, setError] = useState<string>('');

  const run = async (id: number) => {
    setLoading(true);
    setError('');
    setResult(null);
    try {
      const res = await testVisionTaskFrame(id);
      if (res?.ok === false) {
        setError(res.error || formatMessage('Vision.task.testFailed'));
      } else {
        setResult(res);
      }
    } catch (e) {
      setError((e as Error)?.message || formatMessage('Vision.task.testFailed'));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    // 打开(task 变化)时自动跑一次单帧测试——由外部交互触发,非渲染派生状态。
    if (!task) return;
    // eslint-disable-next-line react-hooks/set-state-in-effect
    void run(task.id);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [task?.id]);

  const objects = result?.objects || [];

  return (
    <ProModal
      open={Boolean(task)}
      title={formatMessage('Vision.task.testFrame')}
      width={520}
      fullScreenable={false}
      onCancel={onClose}
      destroyOnHidden
    >
      <div className={styles.testFrameBody}>
        {loading && (
          <div className={styles.testFrameCenter}>
            <Spin />
            <span className={styles.testFrameHint}>{formatMessage('Vision.task.testRunning')}</span>
          </div>
        )}
        {!loading && error && <Alert type="error" showIcon message={error} />}
        {!loading && !error && result && (
          <>
            <div className={styles.testFrameSummary}>
              <span>
                {formatMessage('Vision.task.testObjects')}: <b>{objects.length}</b>
              </span>
              {typeof result.inferenceMs === 'number' && (
                <span>
                  {formatMessage('Vision.task.testLatency')}: <b>{result.inferenceMs} ms</b>
                </span>
              )}
            </div>
            {objects.length === 0 ? (
              <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={formatMessage('Vision.task.testNoObject')} />
            ) : (
              <div className={styles.testFrameTags}>
                {objects.map((obj, idx) => {
                  const score = readScore(obj);
                  return (
                    <Tag key={idx} bordered={false} color="processing">
                      {readLabel(obj)}
                      {score !== null && ` ${(score * 100).toFixed(0)}%`}
                    </Tag>
                  );
                })}
              </div>
            )}
          </>
        )}
      </div>
      <div className={styles.footer}>
        <span />
        <Button type="primary" loading={loading} disabled={!task} onClick={() => task && run(task.id)}>
          {formatMessage('Vision.task.testRerun')}
        </Button>
      </div>
    </ProModal>
  );
};

export default TaskTestFrameModal;
