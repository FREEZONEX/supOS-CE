import { type FC, useState } from 'react';
import classNames from 'classnames';
import { Subtract } from '@carbon/icons-react';
import { App } from 'antd';
import styles from './AppSpaceList.module.scss';
import { useTranslate } from '@/hooks';

import { ButtonPermission } from '@/common-types/button-permission';
import { AuthWrapper } from '@/components/auth';

const AppSpaceList: FC<any> = ({ list, getHtmlList, clickItemId, setClickItemIdItemId, deleteHandle }) => {
  const [hoveredItemId, setHoveredItemId] = useState(null);
  const { modal } = App.useApp();

  const formatMessage = useTranslate();

  return (
    <div className={styles['app-space-list']}>
      {list.map((item: any) => {
        return (
          <div
            key={item.name}
            className={classNames('common-card', item.name === clickItemId ? 'select-card' : 'can-select-card')}
            onClick={() => {
              setClickItemIdItemId(item.name);
              getHtmlList?.(item.name);
            }}
            onMouseEnter={() => setHoveredItemId(item.name)}
            onMouseLeave={() => setHoveredItemId(null)}
          >
            {item.name}
            {item.name === hoveredItemId ? (
              <AuthWrapper auth={ButtonPermission['appSpace.delete']}>
                <div
                  onClick={(e) => {
                    modal.confirm({
                      title: formatMessage('common.deleteConfirm'),
                      onOk: () => deleteHandle(item.name),
                      okButtonProps: {
                        title: formatMessage('common.confirm'),
                      },
                      cancelButtonProps: {
                        title: formatMessage('common.cancel'),
                      },
                    });
                    e.stopPropagation();
                  }}
                >
                  <Subtract />
                </div>
              </AuthWrapper>
            ) : null}
          </div>
        );
      })}
    </div>
  );
};

export default AppSpaceList;
