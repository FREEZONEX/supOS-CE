import { Tag } from 'antd';
import { ArrowRight, AddLarge, ChartLine, WatsonHealth3DCurveAutoColon, SendAlt, Document } from '@carbon/icons-react';
import { useTranslate } from '@/hooks';
import './index.scss';

const EmptyDetail = () => {
  const formatMessage = useTranslate();

  return (
    <div className="emptyDetail-wrap">
      <ul className="detailInfo-list">
        <li className="detailInfo-list-item">
          <Tag>{formatMessage('uns.guideSwitchTabs')}</Tag>
          <ArrowRight className="icon-arrow" size={12} />
          {formatMessage('uns.guideBuild')} <span className="tag-info">{formatMessage('uns.treeList')}</span>/
          <span className="tag-info">{formatMessage('common.template')}</span>/
          <span className="tag-info">{formatMessage('common.label')}</span>
        </li>
        <li className="detailInfo-list-item">
          <Tag>
            {formatMessage('uns.guideClick')} &nbsp;
            <AddLarge size={12} />
          </Tag>
          <ArrowRight className="icon-arrow" size={12} />
          {formatMessage('uns.guideBuildUnsWay', {
            namespace: <span className="tag-info">{formatMessage('uns.treeList')}</span>,
          })}
        </li>
        <li className="detailInfo-list-item">
          <Tag>
            {formatMessage('uns.guideClick')} &nbsp; <WatsonHealth3DCurveAutoColon size={12} />
            {formatMessage('uns.model')}
          </Tag>
          <ArrowRight className="icon-arrow" size={12} />
          {formatMessage('uns.guideBrowseUns', {
            namespace: <span className="tag-info">{formatMessage('uns.treeList')}</span>,
          })}
        </li>
        <li className="detailInfo-list-item">
          <Tag>
            {formatMessage('uns.guideClick')} &nbsp; <ChartLine size={12} /> / <SendAlt size={12} /> /
            <Document size={12} />
            {formatMessage('uns.instance')}
          </Tag>
          <ArrowRight className="icon-arrow" size={12} />
          {formatMessage('uns.guideManageUns', {
            namespace: <span className="tag-info">{formatMessage('uns.treeList')}</span>,
          })}
        </li>
      </ul>
    </div>
  );
};

export default EmptyDetail;
