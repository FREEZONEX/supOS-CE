import type { FC } from 'react';
import { useTranslate } from '@/hooks';

interface HistoryProps {
  protocol: { [key: string]: any };
  dataPath: string;
}
interface windowContentType {
  windowType: string;
  options: { [key: string]: any };
}

const History: FC<HistoryProps> = ({ protocol, dataPath }) => {
  const formatMessage = useTranslate();

  const windowContent = (window: windowContentType) => {
    switch (window?.windowType) {
      case 'INTERVAL':
        return (
          <>
            <div className="detailItem">
              <div className="detailKey">{formatMessage('streams.intervalValue')}</div>
              <div>{window?.options?.intervalValue}</div>
            </div>
            <div className="detailItem">
              <div className="detailKey">{formatMessage('streams.intervalOffset')}</div>
              <div>{window?.options?.intervalOffset}</div>
            </div>
          </>
        );
      case 'SESSION':
        return (
          <div className="detailItem">
            <div className="detailKey">{formatMessage('streams.tolValue')}</div>
            <div>{window?.options?.tolValue}</div>
          </div>
        );
      case 'STATE_WINDOW':
        return (
          <div className="detailItem">
            <div className="detailKey">{formatMessage('streams.referenceField')}</div>
            <div>{window?.options?.field}</div>
          </div>
        );
      case 'EVENT_WINDOW':
        return (
          <>
            <div className="detailItem">
              <div className="detailKey">{formatMessage('streams.startWith')}</div>
              <div>{window?.options?.startWith}</div>
            </div>
            <div className="detailItem">
              <div className="detailKey">{formatMessage('streams.endWith')}</div>
              <div>{window?.options?.endWith}</div>
            </div>
          </>
        );
      case 'COUNT_WINDOW':
        return (
          <>
            <div className="detailItem">
              <div className="detailKey">{formatMessage('streams.countValue')}</div>
              <div>{window?.options?.countValue}</div>
            </div>
            <div className="detailItem">
              <div className="detailKey">{formatMessage('streams.slidingValue')}</div>
              <div>{window?.options?.slidingValue}</div>
            </div>
          </>
        );
      default:
        return null;
    }
  };

  return (
    <>
      <div className="detailItem">
        <div className="detailKey">{formatMessage('streams.dataSource')}</div>
        <div>{dataPath}</div>
      </div>
      {protocol?.whereCondition && (
        <div className="detailItem">
          <div className="detailKey">{formatMessage('streams.whereCondition')}</div>
          <div>{protocol.whereCondition.replace(/\$(.*?)#/g, '$1')}</div>
        </div>
      )}
      {protocol?.havingCondition && (
        <div className="detailItem">
          <div className="detailKey">{formatMessage('streams.havingCondition')}</div>
          <div>{protocol.havingCondition.replace(/\$(.*?)#/g, '$1')}</div>
        </div>
      )}
      {protocol.window && (
        <>
          <div className="detailItem">
            <div className="detailKey">{formatMessage('streams.windowType')}</div>
            <div>{protocol.window.windowType}</div>
          </div>
          {windowContent(protocol.window)}
        </>
      )}
      {protocol?.trigger && (
        <div className="detailItem">
          <div className="detailKey">{formatMessage('streams.trigger')}</div>
          <div>{protocol.trigger}</div>
        </div>
      )}
      {protocol?.waterMark && (
        <div className="detailItem">
          <div className="detailKey">{formatMessage('streams.watermark')}</div>
          <div>{protocol.waterMark}</div>
        </div>
      )}
      {protocol?.deleteMark && (
        <div className="detailItem">
          <div className="detailKey">{formatMessage('streams.deleteMark')}</div>
          <div>{protocol.deleteMark}</div>
        </div>
      )}
      {[true, false].includes(protocol?.fillHistory) && (
        <div className="detailItem">
          <div className="detailKey">{formatMessage('streams.fillHistory')}</div>
          <div>{formatMessage(protocol?.fillHistory ? 'uns.true' : 'uns.false')}</div>
        </div>
      )}
      {[true, false].includes(protocol?.ignoreUpdate) && (
        <div className="detailItem">
          <div className="detailKey">{formatMessage('streams.ignoreUpdate')}</div>
          <div>{formatMessage(protocol?.ignoreUpdate ? 'uns.true' : 'uns.false')}</div>
        </div>
      )}
      {[true, false].includes(protocol?.ignoreExpired) && (
        <div className="detailItem">
          <div className="detailKey">{formatMessage('streams.ignoreExpired')}</div>
          <div>{formatMessage(protocol?.ignoreExpired ? 'uns.true' : 'uns.false')}</div>
        </div>
      )}
      {protocol?.startTime && (
        <div className="detailItem">
          <div className="detailKey">{formatMessage('streams.startTime')}</div>
          <div>{protocol.startTime}</div>
        </div>
      )}
      {protocol?.endTime && (
        <div className="detailItem">
          <div className="detailKey">{formatMessage('streams.endTime')}</div>
          <div>{protocol.endTime}</div>
        </div>
      )}
    </>
  );
};
export default History;
