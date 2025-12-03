import card from '@/assets/app-display/card.svg';
import cardBg from '@/assets/app-display/card-bg.svg';
import cardDark from '@/assets/app-display/card-dark.svg';
import cardChartBg from '@/assets/app-display/card-chartreuse-bg.svg';
import useAddModal from '@/pages/app-management/components/useAddModal';
import { useNavigate } from 'react-router';
import { type FC, useState } from 'react';
import { AddLarge, Close } from '@carbon/icons-react';
import { destroyApp } from '@/apis/inter-api/apps';
import styles from './AppList.module.scss';
import { useThemeStore } from '@/stores/theme-store.ts';

const AppList: FC<any> = ({ list, successCallBack }) => {
  const navigate = useNavigate();
  const { theme, primaryColor } = useThemeStore((state) => ({
    theme: state.theme,
    primaryColor: state.primaryColor,
  }));
  const { ModalDom, setModalOpen } = useAddModal({ successCallBack });
  const [hoveredItemId, setHoveredItemId] = useState(null);
  const onAddHandle = () => {
    setModalOpen(true);
  };
  const goHandle = (item: any) => {
    if (item.homepage) {
      navigate('/app-iframe', {
        state: { title: item?.name, src: item?.homepage },
      });
    } else {
      navigate('/app-space', {
        state: { name: item?.name },
      });
    }
  };

  const deleteHandle = (appName: string) => {
    destroyApp(appName).then(() => {
      successCallBack();
    });
  };
  const imgHandle = (name: string, hoveredItemId: string) => {
    if (theme.includes('dark')) {
      return cardDark;
    }
    if (hoveredItemId === name) {
      if (primaryColor.includes('chartreuse')) {
        return cardChartBg;
      }
      return cardBg;
    }
    return card;
  };

  return (
    <div className={styles['app-list']}>
      {ModalDom}
      {
        // <AuthWrapper auth={ButtonPermission['appDisplay.add']}>
        <div className="add-card" onClick={onAddHandle}>
          <AddLarge size={106} color="var(--supos-text-color)" />
        </div>
        // </AuthWrapper>
      }
      {list.map((item: any) => (
        <div
          onClick={(e) => {
            e.stopPropagation();
            goHandle(item);
          }}
          onMouseEnter={() => setHoveredItemId(item.name)}
          onMouseLeave={() => setHoveredItemId(null)}
          className="card"
          key={item.name}
        >
          {item.name === hoveredItemId && (
            // <AuthWrapper auth={ButtonPermission['appDisplay.delete']}>
            <div className="icon">
              <Close
                size={20}
                onClick={(e) => {
                  e.stopPropagation();
                  deleteHandle(item.name);
                }}
              />
            </div>
            // </AuthWrapper>
          )}

          <div className="name" title={item.name}>
            <div className="name-text">{item.name}</div>
          </div>
          <img src={imgHandle(item.name, hoveredItemId || '')} />
        </div>
      ))}
    </div>
  );
};

export default AppList;
