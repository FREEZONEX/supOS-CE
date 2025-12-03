import { Code, Copy, DashboardReference, Delete, Home } from '@carbon/icons-react';
import { useClipboard } from '@/hooks';
import { type FC, useRef, useState } from 'react';
import { useTranslate } from '@/hooks';
import { Tooltip } from 'antd';
import { ButtonPermission } from '@/common-types/button-permission';
import styles from './AppUrlPreview.module.scss';
import { AuthWrapper } from '@/components/auth';
import { getBaseUrl, getBaseFileName } from '@/utils/url-util';

const AppUrlPreview: FC<any> = ({ editHandle, setHomepage, item, deleteHandle }) => {
  const formatMessage = useTranslate();
  const originUrl = getBaseUrl() + item.url;
  const buttonRef = useRef<any>(null);
  useClipboard(buttonRef, originUrl);
  const [isClicked, setIsClicked] = useState(false);
  const handleClick = () => {
    window.open(originUrl);
    setIsClicked(!isClicked);
  };
  return (
    <div className={styles['app-url-preview']} style={{ borderBottom: '1px  dashed rgba(47, 52, 59, 0.20)' }}>
      <span className="url">
        <Tooltip placement="top" title={formatMessage('appSpace.setHomepage')}>
          <Home
            size={20}
            style={{
              cursor: 'pointer',
              color: item.isHomePage ? 'var(--supos-theme-color)' : 'var(--supos-text-color)',
            }}
            onClick={setHomepage}
          />
        </Tooltip>
        <span className="single-url" title={getBaseFileName(item.url)}>
          {getBaseFileName(item.url)}
        </span>
      </span>
      <div className="box">
        <div title={originUrl} className="single-url">
          {originUrl}
        </div>
        <Tooltip placement="top" title={formatMessage('common.copy')}>
          <div ref={buttonRef} style={{ display: 'flex' }}>
            <Copy style={{ cursor: 'pointer' }} />
          </div>
        </Tooltip>
      </div>
      <AuthWrapper auth={ButtonPermission['appSpace.showPage']}>
        <Tooltip placement="top" title={formatMessage('appSpace.go')}>
          <div className="icon icon-blue" onClick={handleClick}>
            <DashboardReference className="icon-enter" />
            {formatMessage('appSpace.showPage')}
          </div>
        </Tooltip>
      </AuthWrapper>
      <AuthWrapper auth={ButtonPermission['appSpace.coding']}>
        <Tooltip placement="top" title={formatMessage('appSpace.coding')}>
          <div className="icon icon-blue" onClick={editHandle}>
            <Code />
            {formatMessage('appSpace.coding')}
          </div>
        </Tooltip>
      </AuthWrapper>
      <AuthWrapper auth={ButtonPermission['appSpace.deleteHTML']}>
        <Tooltip placement="top" title={formatMessage('appSpace.deleteHTML')}>
          <div className="icon icon-grey" onClick={deleteHandle}>
            <Delete
              style={{
                transform: 'rotate(180deg)',
              }}
            />
            {formatMessage('common.delete')}
          </div>
        </Tooltip>
      </AuthWrapper>
    </div>
  );
};

export default AppUrlPreview;
