import { Tag } from 'antd';
import { ArrowRight, DocumentAdd, FolderAdd, Renew, Draggable } from '@carbon/icons-react';
import { useTranslate } from '@/hooks';
import styles from './index.module.scss';

const EmptyDetail = () => {
  const formatMessage = useTranslate();

  return (
    <div className={styles['emptyDetail-wrap']}>
      <ul className={styles['detailInfo-list']}>
        <li className={styles['detailInfo-list-item']}>
          <Tag>
            {formatMessage('uns.guideClick')} &nbsp;
            <DocumentAdd size={12} />
            &nbsp;/&nbsp; <FolderAdd size={12} />
          </Tag>
          <ArrowRight className="icon-arrow" size={12} />
          {formatMessage('MenuConfiguration.guideCreateMenu')}
        </li>
        <li className={styles['detailInfo-list-item']}>
          <Tag>
            {formatMessage('uns.guideClick')} &nbsp; <Renew size={12} />
            {formatMessage('common.refresh')}
          </Tag>
          <ArrowRight className="icon-arrow" size={12} />
          {formatMessage('MenuConfiguration.guideRefreshMenu')}
        </li>
        <li className={styles['detailInfo-list-item']}>
          <Tag>
            {formatMessage('common.drag')} &nbsp; <Draggable size={12} />
            {formatMessage('MenuConfiguration.menu')}
          </Tag>
          <ArrowRight className="icon-arrow" size={12} />
          {formatMessage('MenuConfiguration.guideSortMenu')}
        </li>
      </ul>
    </div>
  );
};

export default EmptyDetail;
