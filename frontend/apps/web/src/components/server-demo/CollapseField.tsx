import type { CSSProperties } from 'react';
import { CaretRightOutlined } from '@ant-design/icons';
import { Collapse, theme } from 'antd';
import type { CollapseProps } from 'antd';
import ServerDemo from '.';

const CollapseField = (props: any) => {
  const { className, value } = props;
  const { token } = theme.useToken();

  const panelStyle: React.CSSProperties = {
    marginBottom: 24,
    background: token.colorFillAlter,
    borderRadius: token.borderRadiusLG,
    border: 'none',
  };

  const getItems: (panelStyle: CSSProperties) => CollapseProps['items'] = (panelStyle) => {
    return value?.map((item: any) => {
      const { key, label, children } = item;
      return {
        key,
        label,
        children: <ServerDemo tabItems={[{ key, label, leftFormItems: children }]} />,
        style: panelStyle,
      };
    });
  };

  return (
    <Collapse
      className={className}
      bordered={false}
      defaultActiveKey={['1', '2', '3', '4', '5']}
      expandIcon={({ isActive }) => <CaretRightOutlined rotate={isActive ? 90 : 0} />}
      style={{ background: 'var(--supos-switchwrap-active-bg-color)' }}
      items={getItems(panelStyle)}
    />
  );
};

export default CollapseField;
