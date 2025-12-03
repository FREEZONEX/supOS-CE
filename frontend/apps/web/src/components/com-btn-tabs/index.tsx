import { useEffect, useState } from 'react';
import classNames from 'classnames';
import './index.scss';

export interface OptionTypes {
  label: string;
  value: string;
}

interface PropsTypes {
  activeKey?: string;
  className?: string;
  options: OptionTypes[];
  style?: any;
  onSelect?: (item: OptionTypes) => void;
}

const ButtonTabs = (props: PropsTypes) => {
  const { activeKey: propsActiveKey, className, options, style, onSelect } = props;

  const [activeKey, setActiveKey] = useState<string>();

  useEffect(() => {
    setActiveKey(propsActiveKey);
  }, [propsActiveKey]);

  const handleSelect = (item: OptionTypes) => {
    if (onSelect) {
      return onSelect(item);
    }
    setActiveKey(item.value);
  };

  return (
    <div className={classNames('com-btn-tabs', className)} style={style}>
      {options.map((item) => {
        return (
          <div
            key={item.value}
            className={classNames('com-tab', { ['com-tab-active']: activeKey === item.value })}
            onClick={() => handleSelect(item)}
          >
            <div className="com-tab-label">{item.label}</div>
          </div>
        );
      })}
    </div>
  );
};

export default ButtonTabs;
