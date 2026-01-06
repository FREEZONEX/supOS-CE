import { type FC, useState, useEffect } from 'react';
import { Form, Select, Flex, Button, Divider, App } from 'antd';
import {
  saveMount,
  getCollectorList,
  // getDeviceList
} from '@/apis/inter-api/uns';

export interface CollectorFromProps {
  formatMessage: any;
  close: (refreshTree?: boolean) => void;
}
const CollectorFrom: FC<CollectorFromProps> = ({ formatMessage, close }) => {
  const form = Form.useFormInstance();
  const { message } = App.useApp();
  const [loading, setLoading] = useState(false);
  const [collectorList, setCollectorList] = useState<any>([]);
  // const [deviceList, setDeviceList] = useState<any>([]);

  const source = Form.useWatch('source', form) || form.getFieldValue('source');

  useEffect(() => {
    getCollectorList({ collectorType: source }).then((res) => {
      setCollectorList(res || []);
      form.setFieldsValue({ collector: undefined, device: undefined });
    });
  }, [source]);

  // useEffect(() => {
  //   form.setFieldValue('device', undefined);
  // }, [deviceList]);

  // const collectorChange = (item: any) => {
  //   if (item?.value) {
  //     getDeviceList(item.value).then((res) => {
  //       setDeviceList(res || []);
  //     });
  //   }
  // };
  const save = () => {
    form.validateFields().then(async (values) => {
      const {
        targetFolder,
        collector,
        // device,
        dataType,
      } = values;
      const data = {
        targetType: 'folder',
        targetAlias: targetFolder,
        sourceType: source,
        dataType,
        extend: {
          sourceAlias: collector.value,
          sourceName: collector.label,
          // devices: device?.map((e: { label: string; value: string }) => ({ alias: e.value, name: e.label })) || [],
        },
      };
      setLoading(true);
      try {
        const res = await saveMount(data);
        if (res.code === 0) {
          message.success(formatMessage('appGui.saveSuccess'));
          close(true);
        }
        setLoading(false);
      } catch (error) {
        console.error(error);
        setLoading(false);
      }
    });
  };
  return (
    <>
      <Form.Item name="collector" label={formatMessage('uns.collectorName')} rules={[{ required: true }]}>
        <Select
          // onChange={collectorChange}
          options={collectorList}
          labelInValue
          fieldNames={{ label: 'name', value: 'alias' }}
        />
      </Form.Item>
      {/* <Form.Item name="device" label={formatMessage('uns.device')}>
        <Select
          mode="multiple"
          options={deviceList}
          labelInValue
          fieldNames={{ label: 'name', value: 'alias' }}
          allowClear
        />
      </Form.Item> */}
      <Divider style={{ borderColor: '#c6c6c6' }} />
      <Flex justify="flex-end" gap={10}>
        <Button color="primary" variant="solid" size="small" onClick={save} loading={loading}>
          {formatMessage('common.save')}
        </Button>
      </Flex>
    </>
  );
};

export default CollectorFrom;
