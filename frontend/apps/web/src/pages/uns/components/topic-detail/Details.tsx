import type { FC, ReactNode } from 'react';
import { History, Aggregation } from './protocol-table';
import { useTranslate, useClipboard } from '@/hooks';
import { Copy } from '@/components/lucide-icon/carbon';
import { toolbarIconProps } from '@/components/lucide-icon/icon-props';
import { formatTimestamp } from '@/utils/format';

interface DetailsProps {
  instanceInfo: { [key: string]: any };
  updateTime?: number;
  websocketData?: any;
}

const Details: FC<DetailsProps> = ({ instanceInfo, updateTime }) => {
  const formatMessage = useTranslate();
  const { copy } = useClipboard();
  const emptyText = '-';

  const displayValue = (value: any) => {
    if (value === undefined || value === null || value === '') return emptyText;
    if (typeof value === 'boolean') return formatMessage(value ? 'uns.true' : 'uns.false');
    return value;
  };

  const displayTime = (value: any) => {
    const formatted = formatTimestamp(value);
    return formatted || emptyText;
  };

  const renderValue = (label: string, value: ReactNode) => (
    <div className="detailItem">
      <div className="detailKey">{label}</div>
      <div className="detailValue">{value}</div>
    </div>
  );

  const renderCopyValue = (label: string, value: any) => {
    const text = displayValue(value);
    const canCopy = text !== emptyText;
    return renderValue(
      label,
      <span className="detailCopyValue">
        <span className="detailCopyText">{text}</span>
        {canCopy ? (
          <button
            type="button"
            className="detailCopyButton"
            onClick={() => copy(String(text))}
            title={formatMessage('common.copy')}
          >
            <Copy {...toolbarIconProps} />
          </button>
        ) : null}
      </span>
    );
  };

  const topicValue = instanceInfo.topic || instanceInfo.path || instanceInfo.namespace || instanceInfo.alias;

  const dataTypeMap: { [key: number]: string } = {
    1: formatMessage('uns.timeSeries'),
    2: formatMessage('uns.relational'),
    3: formatMessage('uns.realtimeCalculation'),
    4: formatMessage('uns.historicalCalculation'),
    6: formatMessage('uns.aggregation'),
    7: formatMessage('uns.reference'),
    8: formatMessage('uns.jsonb'),
  };
  const renderProtocolTable = (protocol: { [key: string]: any }) => {
    if (instanceInfo.dataType === 4) return <History protocol={protocol} dataPath={instanceInfo.dataPath} />;
    if (instanceInfo.dataType === 6) return <Aggregation protocol={protocol} refers={instanceInfo.refers || []} />;
    return null;
  };

  const mountTypeMap: { [key: number]: string } = {
    16: formatMessage('uns.grpcGateway'),
    50: formatMessage('streams.dataSource'),
    51: formatMessage('streams.dataSource'),
    52: formatMessage('streams.dataSource'),
    100: formatMessage('streams.dataSource'),
  };

  const fileTypeMap: { [key: number]: string } = {
    1: formatMessage('uns.state'),
    2: formatMessage('uns.action'),
    3: formatMessage('uns.metric'),
  };
  const mountLabel = instanceInfo.mount
    ? `${mountTypeMap[instanceInfo.mount?.mountType || 100]}（${
        instanceInfo.mount?.displayName || instanceInfo.mount?.mountSource || emptyText
      }）`
    : '';

  return (
    <>
      {renderCopyValue(formatMessage('uns.topic'), topicValue)}
      {renderCopyValue(formatMessage('uns.alias'), instanceInfo.alias)}
      {renderValue(formatMessage('uns.displayName'), displayValue(instanceInfo.displayName))}
      {renderValue(formatMessage('uns.description'), displayValue(instanceInfo.description))}
      {instanceInfo.mount && renderValue(formatMessage('uns.mountDataSource'), mountLabel)}
      {renderValue(formatMessage('uns.databaseType'), displayValue(dataTypeMap[instanceInfo.dataType]))}
      {instanceInfo.dataType !== 7 &&
        renderValue(
          formatMessage('uns.persistence'),
          formatMessage(instanceInfo.persistence ? 'uns.true' : 'uns.false')
        )}
      {instanceInfo.protocol && renderProtocolTable(instanceInfo.protocol)}
      {instanceInfo.showExpression && (
        <>
          {renderValue(
            formatMessage('common.expression'),
            displayValue(instanceInfo.showExpression.replace(/\$(.*?)#/g, '$1'))
          )}
          {renderValue(
            formatMessage('uns.reference'),
            displayValue(instanceInfo?.refers?.find((e: { uts: boolean }) => e.uts)?.path)
          )}
        </>
      )}
      {renderValue(
        formatMessage('common.creationTime'),
        displayTime(instanceInfo.createTime || instanceInfo.createdTime)
      )}
      {instanceInfo.dataType === 7 &&
        renderValue(formatMessage('uns.referenceTarget'), displayValue(instanceInfo?.refers?.[0]?.path))}
      {![3, 4].includes(instanceInfo.dataType) &&
        renderValue(
          formatMessage('common.latestUpdate'),
          displayTime(updateTime || instanceInfo.updateTime || instanceInfo.updatedTime)
        )}
      {renderValue(formatMessage('uns.namespace'), displayValue(instanceInfo.path || instanceInfo.namespace))}
      {renderValue(formatMessage('uns.originalName'), displayValue(instanceInfo.name))}
      {instanceInfo.extend &&
        Object.keys(instanceInfo.extend).map((item: string, index: number) => (
          <div className="detailItem" key={index}>
            <div className="detailKey">{item}</div>
            <div className="detailValue">{displayValue(instanceInfo.extend[item])}</div>
          </div>
        ))}
      {instanceInfo.parentDataType &&
        renderValue(formatMessage('uns.filesType'), displayValue(fileTypeMap[instanceInfo.parentDataType]))}
    </>
  );
};
export default Details;
