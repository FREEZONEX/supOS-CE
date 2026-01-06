import type { FC } from 'react';
import { Form, Divider } from 'antd';
import { useTranslate, useFormValue } from '@/hooks';
import { useI18nStore } from '@/stores/i18n-store';
import ComRadio from '@/components/com-radio';
import FieldsFormList from '@/pages/uns/components/use-create-modal/components/FieldsFormList';
import ModelFieldsForm from '@/pages/uns/components/use-create-modal/components/file/ModelFieldsForm';
import ReverseGeneration from '@/pages/uns/components/use-create-modal/components/file/ReverseGeneration';

interface AttributeTypeFormProps {
  types?: string[];
  addNamespaceForAi?: { [key: string]: any };
  setAddNamespaceForAi?: (e: any) => void;
  dataType?: number;
  templateList?: any[];
}
const AttributeTypeForm: FC<AttributeTypeFormProps> = ({
  types,
  addNamespaceForAi,
  setAddNamespaceForAi,
  dataType,
  templateList,
}) => {
  const formatMessage = useTranslate();
  const lang = useI18nStore((state) => state.lang);
  const form = Form.useFormInstance();
  const attributeType = useFormValue('attributeType', form);

  const renderContent = () => {
    switch (attributeType) {
      case 1:
        return (
          <FieldsFormList
            showWrap={false}
            types={types}
            addNamespaceForAi={addNamespaceForAi}
            setAddNamespaceForAi={setAddNamespaceForAi}
            showMoreBtn={dataType === 1}
          />
        );
      case 2:
        return <ModelFieldsForm options={(templateList || []).slice(1)} types={types} />;
      case 3:
        return <ReverseGeneration require={dataType !== 8} onlyJson types={types} />;
      default:
        return null;
    }
  };
  return (
    <>
      {/*<Divider style={{ borderColor: '#c6c6c6' }} />*/}
      <div className="dashedWrap">
        <Form.Item
          name="attributeType"
          label={formatMessage('uns.attributeGenerationMethod')}
          initialValue={1}
          tooltip={{
            title: (
              <div>
                <span>• {formatMessage('common.custom')}</span> —&nbsp;
                {formatMessage('uns.attributeGenerationMethodTooltip-Custom')}
                <br />
                <span>• {formatMessage('common.template')}</span> —&nbsp;
                {formatMessage('uns.attributeGenerationMethodTooltip-Template')}
                <br />
                <span>• {formatMessage('uns.reverseGeneration')}</span> —&nbsp;
                {formatMessage('uns.attributeGenerationMethodTooltip-ReverseGeneration')}
              </div>
            ),
          }}
          className={lang === 'en-US' ? 'customLabelStyle' : ''}
        >
          <ComRadio
            options={[
              { label: formatMessage('common.custom'), value: 1 },
              { label: formatMessage('common.template'), value: 2 },
              { label: formatMessage('uns.reverseGeneration'), value: 3 },
            ]}
            onChange={() => {
              form.setFieldsValue({
                fields: [{}],

                modelId: undefined,
                jsonData: undefined,
                jsonList: [],
                jsonDataPath: undefined,
                source: undefined,
                dataSource: undefined,
                table: undefined,
                next: false,
                mainKey: undefined,
              });
            }}
          />
        </Form.Item>
        <Divider style={{ borderColor: '#c6c6c6' }} dashed />
        {renderContent()}
      </div>
    </>
  );
};

export default AttributeTypeForm;
