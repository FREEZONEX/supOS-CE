import { Tag } from 'antd';
import { ArrowRight, ChartLine, ClipboardList, Route, SendAlt, SquarePlus } from '@/components/lucide-icon/carbon';
import { useTranslate } from '@/hooks';
import './index.scss';

const EmptyDetail = () => {
  const formatMessage = useTranslate();

  return (
    <div className="emptyDetail-wrap">
      <ul className="detailInfo-list">
        <li className="detailInfo-list-item">
          <Tag>
            {formatMessage('uns.guideClick')} &nbsp;
            <SquarePlus size={12} strokeWidth={1.75} aria-hidden />
          </Tag>
          <ArrowRight className="icon-arrow" size={12} />
          {formatMessage('uns.guideBuildUnsWay')}
        </li>
        <li className="detailInfo-list-item">
          <Tag>{formatMessage('uns.guideRightClick')}</Tag>
          <ArrowRight className="icon-arrow" size={12} />
          {formatMessage('uns.guideQuickExpand')}
        </li>
        <li className="detailInfo-list-item">
          <Tag>
            {formatMessage('uns.guideClick')} &nbsp; <Route size={12} strokeWidth={1.75} aria-hidden />
            {formatMessage('uns.guidePath')}
          </Tag>
          <ArrowRight className="icon-arrow" size={12} />
          {formatMessage('uns.guideBrowseUns')}
        </li>
        <li className="detailInfo-list-item">
          <Tag>
            {formatMessage('uns.guideClick')} &nbsp; <ChartLine size={12} strokeWidth={1.75} aria-hidden /> /{' '}
            <SendAlt size={12} strokeWidth={1.75} aria-hidden /> /{' '}
            <ClipboardList size={12} strokeWidth={1.75} aria-hidden />
            {formatMessage('uns.guideTopic')}
          </Tag>
          <ArrowRight className="icon-arrow" size={12} />
          {formatMessage('uns.guideManageUns')}
        </li>
      </ul>
    </div>
  );
};

export default EmptyDetail;
