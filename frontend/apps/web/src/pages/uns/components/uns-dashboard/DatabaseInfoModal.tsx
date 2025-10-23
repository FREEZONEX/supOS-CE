import { forwardRef, useImperativeHandle, useState } from 'react';
import ProModal from '@/components/pro-modal';
import { useTranslate } from '@/hooks';
import { Col, Flex, Row } from 'antd';
import ComEllipsis from '@/components/com-ellipsis';
import ComCopy from '@/components/com-copy';

export interface ModalRef {
  onOpen: (props: any) => void;
  onClose: () => void;
}

export interface ModalProps {
  [key: string]: any;
}

const DatabaseInfoModal = forwardRef<ModalRef, ModalProps>((_props, ref) => {
  const [visible, setVisible] = useState(false);
  const formatMessage = useTranslate();
  const [info, setInfo] = useState<any>([]);

  const onOpen = (props: any) => {
    setVisible(true);
    if (!props || props?.dataType === 7) {
      return setInfo([
        {
          label: 'uns.dbEngine',
          value: '',
          span: 12,
          key: 1,
        },
        {
          label: 'uns.host',
          value: '',
          span: 12,
          key: 2,
        },
        {
          label: 'uns.port',
          value: '',
          span: 12,
          key: 3,
        },
        {
          label: 'uns.database',
          value: '',
          span: 12,
          key: 4,
        },
        {
          label: 'uns.schema',
          value: '',
          span: 12,
          key: 5,
        },
        {
          label: 'uns.table',
          value: '',
          span: 12,
          key: 6,
        },
        {
          label: 'uns.connectionString',
          value: '',
          span: 24,
          key: 7,
        },
      ]);
    } else {
      const dbEngine = [2, 6].includes(info?.dataType as number) ? 'PostgreSQL' : 'TimescaleDB';
      const port = dbEngine === 'PostgreSQL' ? 5432 : 2345;
      setInfo([
        {
          label: 'uns.dbEngine',
          value: dbEngine,
          span: 12,
          key: 1,
        },
        {
          label: 'uns.host',
          value: location.host,
          span: 12,
          key: 2,
        },
        {
          label: 'uns.port',
          value: port,
          span: 12,
          key: 3,
        },
        {
          label: 'uns.database',
          value: 'postgres',
          span: 12,
          key: 4,
        },
        {
          label: 'uns.schema',
          value: 'public',
          span: 12,
          key: 5,
        },
        {
          label: 'uns.table',
          value: props?.table,
          span: 12,
          key: 6,
        },
        {
          label: 'uns.connectionString',
          value: `jdbc:postgresql://${location.host}:${port}/postgres?currentSchema=public`,
          span: 24,
          key: 7,
        },
      ]);
    }
  };
  const onClose = () => {
    setVisible(false);
  };
  useImperativeHandle(ref, () => ({
    onOpen,
    onClose,
  }));

  return (
    <ProModal
      open={visible}
      onCancel={onClose}
      title={formatMessage('uns.databaseInfo')}
      width={650}
      styles={{
        body: {
          paddingBlockStart: 0,
        },
      }}
    >
      {() => {
        return (
          <Row gutter={[16, 16]} style={{ paddingTop: 16 }}>
            {info?.map((item: any) => {
              return (
                <Col key={item.key} span={item.span}>
                  <ComEllipsis style={{ fontWeight: 400, fontSize: 12, opacity: 0.7 }}>
                    {formatMessage(item.label)}
                  </ComEllipsis>
                  <Flex
                    title={item.value}
                    style={{
                      background: 'var(--supos-bg-color)',
                      padding: '4px 12px',
                      borderRadius: '3px',
                      border: '1px solid #E0E0E0',
                      marginTop: 8,
                    }}
                    align="center"
                    justify="space-between"
                  >
                    <pre
                      style={{
                        overflow: 'hidden',
                        whiteSpace: 'nowrap',
                        textOverflow: 'ellipsis',
                      }}
                    >
                      {item.value}
                    </pre>
                    <ComCopy textToCopy={item.value} />
                  </Flex>
                </Col>
              );
            })}
          </Row>
        );
      }}
    </ProModal>
  );
});

export default DatabaseInfoModal;
