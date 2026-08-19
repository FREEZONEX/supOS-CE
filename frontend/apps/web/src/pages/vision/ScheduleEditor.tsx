import { TimePicker } from 'antd';
import dayjs from 'dayjs';
import type { VisionScheduleRule } from '@/apis/core-api/task';
import { Add, TrashCan } from '@/components/lucide-icon/carbon';
import { useTranslate } from '@/hooks';
import styles from './index.module.scss';

type ScheduleEditorProps = {
  value: VisionScheduleRule[];
  onChange: (rules: VisionScheduleRule[]) => void;
};

// 周一=0 … 周日=6(对齐 worker),展示按 Mon..Sun。
const WEEKDAYS = [0, 1, 2, 3, 4, 5, 6];
const WEEKDAY_KEYS = ['mon', 'tue', 'wed', 'thu', 'fri', 'sat', 'sun'];

// 设计稿:整块浅灰底,每条规则一张白卡(Weekdays / Start Time / End Time 三组带小标题),
// 底部一张只放加号的白卡用于新增。选中的星期是黑底白字,不用主题绿。
const ScheduleEditor = ({ value, onChange }: ScheduleEditorProps) => {
  const formatMessage = useTranslate();
  const rules = value.length > 0 ? value : [{ weekdays: [0, 1, 2, 3, 4], start: '08:00', end: '18:00' }];

  const update = (idx: number, patch: Partial<VisionScheduleRule>) =>
    onChange(rules.map((r, i) => (i === idx ? { ...r, ...patch } : r)));

  const toggleWeekday = (idx: number, day: number) => {
    const set = new Set(rules[idx].weekdays);
    if (set.has(day)) set.delete(day);
    else set.add(day);
    update(idx, { weekdays: [...set].sort((a, b) => a - b) });
  };

  const addRule = () => onChange([...rules, { weekdays: [0, 1, 2, 3, 4], start: '08:00', end: '18:00' }]);
  const removeRule = (idx: number) => onChange(rules.filter((_, i) => i !== idx));

  return (
    <div className={styles.scheduleEditor}>
      {rules.map((rule, idx) => (
        <div key={idx} className={styles.scheduleCard}>
          <div className={styles.scheduleField}>
            <span className={styles.scheduleFieldLabel}>{formatMessage('Vision.task.weekdaysLabel')}</span>
            <div className={styles.weekdayGroup}>
              {WEEKDAYS.map((day) => (
                <button
                  key={day}
                  type="button"
                  className={`${styles.weekdayBtn} ${rule.weekdays.includes(day) ? styles.weekdayOn : ''}`}
                  onClick={() => toggleWeekday(idx, day)}
                >
                  {formatMessage(`Vision.task.weekday.${WEEKDAY_KEYS[day]}`)}
                </button>
              ))}
            </div>
          </div>
          <div className={styles.scheduleField}>
            <span className={styles.scheduleFieldLabel}>{formatMessage('Vision.task.startTime')}</span>
            <TimePicker
              size="small"
              format="HH:mm"
              allowClear={false}
              className={styles.scheduleTime}
              value={dayjs(rule.start, 'HH:mm')}
              onChange={(v) => update(idx, { start: v ? v.format('HH:mm') : '00:00' })}
            />
          </div>
          <div className={styles.scheduleField}>
            <span className={styles.scheduleFieldLabel}>{formatMessage('Vision.task.endTime')}</span>
            <TimePicker
              size="small"
              format="HH:mm"
              allowClear={false}
              className={styles.scheduleTime}
              value={dayjs(rule.end, 'HH:mm')}
              onChange={(v) => update(idx, { end: v ? v.format('HH:mm') : '23:59' })}
            />
          </div>
          {rules.length > 1 && (
            <button
              type="button"
              className={styles.scheduleDel}
              aria-label={formatMessage('common.delete')}
              onClick={() => removeRule(idx)}
            >
              <TrashCan size={16} />
            </button>
          )}
        </div>
      ))}
      <button
        type="button"
        className={styles.scheduleAdd}
        aria-label={formatMessage('Vision.task.addSchedule')}
        onClick={addRule}
      >
        <Add size={16} />
      </button>
    </div>
  );
};

export default ScheduleEditor;
