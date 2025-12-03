import { type FC, useState, useEffect } from 'react';
import { Form, Flex, Button, Input } from 'antd';
import { ChevronRight, ChevronLeft } from '@carbon/icons-react';
import { useTranslate } from '@/hooks';
import FieldsFormList from '@/pages/uns/components/use-create-modal/components/FieldsFormList';
import ComRadio from '@/components/com-radio';

export interface TopicToUnsFieldsListProps {
  types?: string[];
}

const TopicToUnsFieldsList: FC<TopicToUnsFieldsListProps> = ({ types }) => {
  const form = Form.useFormInstance();
  const formatMessage = useTranslate();

  const next = Form.useWatch('next', form) || form.getFieldValue('next');
  const jsonList = Form.useWatch('jsonList', form) || form.getFieldValue('jsonList');
  const jsonDataPath = Form.useWatch('jsonDataPath', form) || form.getFieldValue('jsonDataPath');
  const [hasMoreJson, setHasMoreJson] = useState(false);

  useEffect(() => {
    setHasMoreJson(jsonList?.length > 1);
  }, [jsonList]);

  const nextStep = () => {
    form.setFieldsValue({ next: true });
  };

  const back = () => {
    const fields = jsonList.find((item: any) => item.dataPath === jsonDataPath)?.fields;
    form.setFieldsValue({ fields, next: false });
  };

  return (
    <div className="dashedWrap">
      <Form.Item name="next" hidden initialValue={false}>
        <Input />
      </Form.Item>
      <Form.Item name="jsonList" hidden initialValue={[]}>
        <Input />
      </Form.Item>

      {hasMoreJson && !next && (
        <Form.Item name="jsonDataPath" label={formatMessage('uns.schemaGenerated')}>
          <ComRadio
            style={{ flexWrap: 'wrap' }}
            options={jsonList.map((item: any) => ({ label: item.dataPath, value: item.dataPath }))}
            onChange={(e) => {
              const fields = jsonList.find((item: any) => item.dataPath === e.target.value)?.fields;
              form.setFieldsValue({ fields });
            }}
          />
        </Form.Item>
      )}
      <FieldsFormList types={types} disabled={!next && hasMoreJson} showMainKey={next} showWrap={false} />
      {hasMoreJson && (
        <Flex justify="flex-end" style={{ marginTop: next ? '20px' : '' }} gap={10}>
          {next ? (
            <Button color="primary" variant="filled" size="small" icon={<ChevronLeft />} onClick={back}>
              {formatMessage('common.back')}
            </Button>
          ) : (
            <Button
              color="primary"
              variant="filled"
              size="small"
              icon={<ChevronRight />}
              iconPosition="end"
              onClick={nextStep}
            >
              {formatMessage('common.next')}
            </Button>
          )}
        </Flex>
      )}
    </div>
  );
};

export default TopicToUnsFieldsList;
