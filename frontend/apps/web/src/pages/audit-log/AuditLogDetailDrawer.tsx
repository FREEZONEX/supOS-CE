import { Button, Drawer, Tag, Tooltip, Typography } from 'antd';
import type { AuditLogItem } from '@/apis/core-api/audit-log';
import { Close } from '@/components/lucide-icon/carbon';
import { useTranslate } from '@/hooks';
import { formatTimestamp } from '@/utils/format';
import {
  getAuditActionLabel,
  getAuditModuleLabel,
  getAuditOperatorLabel,
  getAuditResourceLabel,
  getAuditResultLabel,
  getAuditScopeLabel,
} from './presentation';
import styles from './index.module.scss';

const { Paragraph } = Typography;

const formatDetailJson = (value?: string) => {
  const raw = value?.trim();
  if (!raw) return '{}';
  try {
    return JSON.stringify(JSON.parse(raw), null, 2);
  } catch {
    return raw;
  }
};

type AuditLogDetailDrawerProps = {
  open: boolean;
  loading: boolean;
  detail?: AuditLogItem | null;
  onClose: () => void;
};

export const AuditLogDetailDrawer = ({ open, loading, detail, onClose }: AuditLogDetailDrawerProps) => {
  const formatMessage = useTranslate();
  const translate = (id: string, opt?: any, defaultMessage?: string) => formatMessage(id, opt, defaultMessage);
  const t = (id: string, fallback: string, opt?: Record<string, unknown>) => translate(id, opt, fallback);

  return (
    <Drawer
      rootClassName={styles.auditDrawer}
      open={open}
      width={720}
      title={t('auditLog.detailTitle', 'Audit Details')}
      closable={false}
      destroyOnClose
      onClose={onClose}
      extra={
        <Tooltip title={t('common.close', 'Close')}>
          <Button color="default" variant="text" onClick={onClose} icon={<Close size={20} />} />
        </Tooltip>
      }
      classNames={{
        body: styles.auditDrawerBody,
      }}
    >
      {loading ? (
        <Paragraph>{t('common.loading', 'Loading...')}</Paragraph>
      ) : detail ? (
        <div className={styles.detailPanel}>
          <section className={styles.detailSection}>
            <div className={styles.sectionHeader}>{t('auditLog.overview', 'Overview')}</div>
            <div className={styles.detailGrid}>
              <div className={styles.detailItem}>
                <span className={styles.detailLabel}>{t('auditLog.time', 'Time')}</span>
                <span className={styles.detailValue}>{formatTimestamp(detail.createdAt) || '-'}</span>
              </div>
              <div className={styles.detailItem}>
                <span className={styles.detailLabel}>{t('auditLog.operator', 'Operator')}</span>
                <span className={styles.detailValue}>{getAuditOperatorLabel(detail)}</span>
              </div>
              <div className={styles.detailItem}>
                <span className={styles.detailLabel}>{t('auditLog.scope', 'Scope')}</span>
                <span className={styles.detailValue}>{getAuditScopeLabel(detail, translate)}</span>
              </div>
              <div className={styles.detailItem}>
                <span className={styles.detailLabel}>{t('auditLog.module', 'Module')}</span>
                <span className={styles.detailValue}>
                  {getAuditModuleLabel(detail.scopeType, detail.scopeName, detail.resType, translate)}
                </span>
              </div>
              <div className={styles.detailItem}>
                <span className={styles.detailLabel}>{t('auditLog.action', 'Action')}</span>
                <span className={styles.detailValue}>{getAuditActionLabel(detail.businessType, translate)}</span>
              </div>
              <div className={styles.detailItem}>
                <span className={styles.detailLabel}>{t('auditLog.result', 'Result')}</span>
                <Tag className={detail.code >= 200 && detail.code < 300 ? styles.successTag : styles.failedTag}>
                  {getAuditResultLabel(detail, translate)}
                </Tag>
              </div>
              <div className={`${styles.detailItem} ${styles.fullSpan}`}>
                <span className={styles.detailLabel}>{t('auditLog.resource', 'Resource')}</span>
                <span className={styles.detailValue}>{getAuditResourceLabel(detail, translate)}</span>
              </div>
            </div>
          </section>

          <section className={styles.detailSection}>
            <div className={styles.sectionHeader}>{t('auditLog.request', 'Request')}</div>
            <div className={styles.detailGrid}>
              <div className={styles.detailItem}>
                <span className={styles.detailLabel}>{t('auditLog.code', 'Code')}</span>
                <span className={styles.detailValue}>{detail.code || '-'}</span>
              </div>
              <div className={styles.detailItem}>
                <span className={styles.detailLabel}>{t('auditLog.requestId', 'Request ID')}</span>
                <span className={styles.detailValue}>{detail.requestId || '-'}</span>
              </div>
              <div className={styles.detailItem}>
                <span className={styles.detailLabel}>{t('auditLog.ip', 'IP')}</span>
                <span className={styles.detailValue}>{detail.ip || '-'}</span>
              </div>
              <div className={`${styles.detailItem} ${styles.fullSpan}`}>
                <span className={styles.detailLabel}>{t('auditLog.userAgent', 'User Agent')}</span>
                <span className={styles.detailValue}>{detail.userAgent || '-'}</span>
              </div>
            </div>
          </section>

          <section className={styles.detailSection}>
            <div className={styles.sectionHeader}>{t('auditLog.detailJson', 'Details JSON')}</div>
            <pre className={styles.detailJson}>{formatDetailJson(detail.detailJson)}</pre>
          </section>
        </div>
      ) : (
        <div className={styles.emptyDetail}>{t('auditLog.noDetail', 'No details available')}</div>
      )}
    </Drawer>
  );
};
