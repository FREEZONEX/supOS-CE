import { useTranslate } from '@/hooks';
import type { App } from '@/pages/launchpad/type';
import { ChevronDown, ChevronUp, Search } from '@/components/lucide-icon/carbon';
import { Input } from 'antd';
import { type FC, useEffect, useRef, useState } from 'react';
import styles from './AppSwitcher.module.scss';

interface AppSwitcherProps {
  apps: App[];
  currentAppName?: string;
  onAppChange: (appName: string, app?: App) => void;
}

const AppSwitcher: FC<AppSwitcherProps> = ({ apps, currentAppName, onAppChange }) => {
  const formatMessage = useTranslate();
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const [searchValue, setSearchValue] = useState('');
  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
        setDropdownOpen(false);
        setSearchValue('');
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const filteredApps = apps.filter((a) =>
    (a.displayName || a.appName).toLowerCase().includes(searchValue.toLowerCase())
  );

  const handleSelect = (appName: string, app?: App) => {
    setDropdownOpen(false);
    setSearchValue('');
    onAppChange(appName, app);
  };

  return (
    <div ref={dropdownRef} className={styles.appSwitcher}>
      <button className={styles.switcherTrigger} onClick={() => setDropdownOpen((v) => !v)}>
        <span>{formatMessage('Launchpad.switchApp', {}, 'Switch App')}</span>
        {dropdownOpen ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
      </button>
      {dropdownOpen && (
        <div className={styles.switcherDropdown}>
          <div className={styles.switcherSearch}>
            <Input
              prefix={<Search size={16} />}
              placeholder="Search..."
              value={searchValue}
              onChange={(e) => setSearchValue(e.target.value)}
              autoFocus
            />
          </div>
          <div className={styles.switcherList}>
            {filteredApps.map((a) => (
              <div
                key={a.appId || a.appName}
                className={`${styles.switcherItem} ${
                  a.appName === currentAppName || String(a.appId) === currentAppName ? styles.switcherItemActive : ''
                }`}
                onClick={() => handleSelect(String(a.appId || a.appName), a)}
              >
                <span className={styles.switcherItemName}>{a.displayName || a.appName}</span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default AppSwitcher;
