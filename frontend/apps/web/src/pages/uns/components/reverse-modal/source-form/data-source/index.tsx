import { type FC, useState, useEffect } from 'react';
import { Form, Select, Flex, Button, Divider, App } from 'antd';
import ComCheckbox from '@/components/com-checkbox';
// import HelpTooltip from '@/components/help-tooltip';
import { saveMount, getSourceList } from '@/apis/inter-api/uns';

export interface DataSourceFromProps {
  formatMessage: any;
  close: (refreshTree?: boolean) => void;
}

const DataSourceForm: FC<DataSourceFromProps> = ({ formatMessage, close }) => {
  const form = Form.useFormInstance();
  const { message } = App.useApp();
  const [loading, setLoading] = useState(false);
  const [mqttList, setMqttList] = useState<{ name: string; alias: string; sourceType: string }[]>([]);

  const source = Form.useWatch('source', form) || form.getFieldValue('source');

  useEffect(() => {
    getSourceList({ sourceType: source }).then((res) => {
      setMqttList(res || []);
    });
  }, [source]);

  const save = () => {
    form.validateFields().then(async (values) => {
      const { targetFolder, dataType, dataSource, persistence, dashboard, syncMeta } = values;
      const data = {
        targetType: 'folder',
        targetAlias: targetFolder,
        sourceType: mqttList.find((item) => item?.alias === dataSource.value)?.sourceType,
        dataType,
        extend: {
          sourceAlias: dataSource.value,
          sourceName: dataSource.label,
        },
        persistence,
        dashboard,
        syncMeta,
      };
      setLoading(true);
      try {
        const res = await saveMount(data);
        if (res.code === 0 || res.code === 200) {
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

  const renderBatchChecks = () => {
    return (
      <Flex gap={8}>
        <Form.Item name="persistence" valuePropName="checked" initialValue={false} noStyle>
          <ComCheckbox>{formatMessage('uns.batchPersistence')}</ComCheckbox>
        </Form.Item>
        <Form.Item name="dashboard" valuePropName="checked" initialValue={false} noStyle>
          <ComCheckbox>{formatMessage('uns.batchAutoDashboard')}</ComCheckbox>
        </Form.Item>
        {/* <Form.Item name="syncMeta" valuePropName="checked" initialValue={false} noStyle>
          <ComCheckbox>
            <Flex gap={4} align="center">
              <span>{formatMessage('uns.metadataSynchronization')}</span>
              <HelpTooltip title={formatMessage('uns.metadataSynchronizationTooltip')} />
            </Flex>
          </ComCheckbox>
        </Form.Item> */}
      </Flex>
    );
  };

  return (
    <>
      <Form.Item name="dataSource" label={formatMessage('streams.dataSource')} rules={[{ required: true }]}>
        <Select
          options={mqttList}
          placeholder={formatMessage('streams.dataSource')}
          fieldNames={{ label: 'name', value: 'alias' }}
          labelInValue
        />
      </Form.Item>
      <Divider style={{ borderColor: '#c6c6c6' }} />
      <Flex justify="space-between">
        {renderBatchChecks()}
        <Button color="primary" variant="solid" size="small" onClick={save} loading={loading}>
          {formatMessage('common.save')}
        </Button>
      </Flex>
    </>
  );
};

export default DataSourceForm;
