import { type CSSProperties, type FC, useMemo, useRef, useState } from 'react';
import classNames from 'classnames';
import { Col, Row, Tabs, type TabsProps } from 'antd';
import StickyBox from 'react-sticky-box';
import HtmlField from './HtmlField';
import CodeSnippetField from './CodeSnippetField';
import styles from './index.module.scss';
import CollapseField from './CollapseField';
import ComCopyContent from '../com-copy/ComCopyContent';

interface ItemChildTypes {
  key: string; // 唯一标识
  label?: string; // 标题
  subTitle?: string; // 子标题
  value?: string; // 内容
  type?: string; // 组件类型
  minCollapsedNumberOfRows?: number; // 最小行数，仅对type为codeSnippet生效
  maxCollapsedNumberOfRows?: number; // 最大行数，仅对type为codeSnippet生效
  isJSON?: boolean; // 是否是json，仅对type为codeSnippet生效
  collapseItems?: any[]; // 折叠项，仅对type为collapse生效
  style?: CSSProperties;
}

export interface TabItems {
  key: string;
  label: string;
  leftFormItems: ItemChildTypes[]; // 左侧表单项
  rightFormItems?: ItemChildTypes[]; // 右侧表单项
}

interface PropsTypes {
  className?: string;
  labelClassName?: string;
  title?: string;
  tabItems?: TabItems[];
}

const components: { [x: string]: FC<any> } = {
  copyContent: ComCopyContent,
  html: HtmlField,
  codeSnippet: CodeSnippetField,
  collapse: CollapseField,
};

const ServerDemo = (props: PropsTypes) => {
  const { className, labelClassName, title, tabItems } = props;

  const rightFormItemsRef = useRef<any>(null);

  const [tab, setTab] = useState<string>('');
  const [heights, setHeights] = useState<number[]>([]);

  const activeKey = useMemo(() => {
    return tab || tabItems?.[0]?.key;
  }, [tab, tabItems]);

  const activeItem = useMemo(() => {
    return tabItems?.find((item) => item.key === activeKey);
  }, [tabItems, activeKey]);

  const height = useMemo(() => {
    if (!heights.length) return 'auto';

    const h = heights?.reduce((prev, curr) => {
      return prev + curr;
    }, 0);

    if (!h) return 'auto';

    if (rightFormItemsRef.current?.offsetHeight > h) return rightFormItemsRef.current?.offsetHeight;

    return h + 2;
  }, [heights]);

  const renderComp = (item: any, index: number, direction: 'left' | 'right') => {
    const {
      key,
      type,
      label,
      subTitle,
      value,
      collapseItems,
      minCollapsedNumberOfRows,
      maxCollapsedNumberOfRows,
      isJSON,
      style,
    } = item;
    const Comp = components[type] ?? ComCopyContent;
    const params: { [x: string]: any } = { label, subTitle, labelClassName, style };

    switch (type) {
      case 'copyContent':
        params.textToCopy = value;
        break;
      case 'html':
        params.value = value;
        break;
      case 'collapse':
        params.value = collapseItems;
        break;
      case 'codeSnippet':
        params.onSizeChange = (params: any) => {
          if (direction === 'right' && params?.height) {
            setHeights((v) => {
              v[index] = params?.height;
              return [...v];
            });
          }
        };
        params.value = value;
        params.minCollapsedNumberOfRows = minCollapsedNumberOfRows;
        params.maxCollapsedNumberOfRows = maxCollapsedNumberOfRows;
        params.isJSON = isJSON;
        break;
      default:
        params.textToCopy = value;
        break;
    }
    return <Comp key={key} className={styles.fieldItem} {...params} />;
  };

  const items = useMemo(() => {
    return (tabItems || []).map((item: any) => {
      return {
        label: item.label,
        key: item.key,
        children: item.leftFormItems?.map((item: ItemChildTypes, index: number) => renderComp(item, index, 'left')),
      };
    });
  }, [tabItems]);

  const renderTabBar: TabsProps['renderTabBar'] = (props, DefaultTabBar) => (
    <StickyBox offsetTop={-12} offsetBottom={0} style={{ zIndex: 1 }}>
      <DefaultTabBar {...props} />
    </StickyBox>
  );

  return (
    <div className={classNames(styles['upload-data'], className)}>
      <Row className={styles['upload-data-info']} gutter={12}>
        <Col span={activeItem?.rightFormItems?.length ? 12 : 24} style={{ height }}>
          <div className={styles.leftFormItems}>
            {title && (
              <div className={styles.title} style={{ paddingBottom: items?.length > 1 ? 0 : 15 }}>
                {title}
              </div>
            )}
            {items?.length > 1 ? (
              <Tabs renderTabBar={renderTabBar} activeKey={activeKey} onChange={setTab} items={items}></Tabs>
            ) : (
              <div>{items[0]?.children}</div>
            )}
          </div>
        </Col>
        {!!activeItem?.rightFormItems?.length && (
          <Col span={12}>
            <div className={styles.rightFormItems}>
              <div ref={rightFormItemsRef}>
                {activeItem.rightFormItems?.map((item: ItemChildTypes, index: number) =>
                  renderComp(item, index, 'right')
                )}
              </div>
            </div>
          </Col>
        )}
      </Row>
    </div>
  );
};

export default ServerDemo;
