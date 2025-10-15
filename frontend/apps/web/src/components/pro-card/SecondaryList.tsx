import { Col, Row } from 'antd';
import './secondaryList.scss';
import type { CSSProperties } from 'react';

const SecondaryList = ({
  options,
  colon = true,
}: {
  options: {
    label: string;
    content: string;
    span: number;
    labelStyle?: CSSProperties;
    contentStyle?: CSSProperties;
    wrapperStyle?: CSSProperties;
    key: string | number;
  }[];
  colon?: boolean;
}) => {
  return (
    <Row style={{ overflow: 'hidden', margin: '4px 0' }} className="secondaryList">
      {options?.map(({ label, labelStyle, contentStyle, span, content, wrapperStyle, key }) => {
        return (
          <Col span={span} style={{ display: 'flex', ...wrapperStyle }} key={key}>
            <div style={labelStyle ? labelStyle : { maxWidth: 100 }} className="span-ellipsis" title={label}>
              {label}
            </div>
            {colon && <span style={{ paddingRight: 4 }}>:</span>}
            <div
              style={contentStyle ? contentStyle : { flex: 1, minWidth: 0 }}
              className="span-ellipsis"
              title={content}
            >
              {content}
            </div>
          </Col>
        );
      })}
    </Row>
  );
};

export default SecondaryList;
