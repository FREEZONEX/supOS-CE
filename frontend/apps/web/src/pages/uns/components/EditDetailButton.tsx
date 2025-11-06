import { useTranslate } from '@/hooks';
import { useEffect, useMemo, useState } from 'react';
import Icon from '@ant-design/icons';
import { App, Form, Flex, Button } from 'antd';
import ExpandedKeyFormList from '@/pages/uns/components/ExpandedKeyFormList.tsx';
import { modifyDetail, getInstanceInfo } from '@/apis/inter-api/uns.ts';
import { AuthButton } from '@/components/auth';
import OperationForm from '@/components/operation-form';
import ProModal from '@/components/pro-modal';
import FileEdit from '@/components/svg-components/FileEdit';
import { cloneDeep } from 'lodash-es';
import ExpressionForm from '@/pages/uns/components/use-create-modal/components/file/timeSeries/ExpressionForm';
import SearchSelect from '@/pages/uns/components/use-create-modal/components/SearchSelect.tsx';
import { getExpression } from '@/utils/uns';

type ReferType = {
  id: string;
  path: string;
  field: string;
  uts?: boolean;
  variableName: string;
};

type ReferItemType = {
  refer: {
    label: string;
    value: string;
  };
  field: string;
};

const extendToArr = (extend: { [key: string]: string }) => {
  if (!extend) return undefined;
  const arr: { key: string; value: string }[] = [];
  Object.keys(extend).forEach((item) => {
    arr.push({
      key: item,
      value: extend[item],
    });
  });
  return arr;
};

const extendToObj = (extend: { key: string; value: string }[]) => {
  if (!extend) return undefined;
  const obj: { [key: string]: string } = {};
  extend.forEach((item) => {
    obj[item.key] = item.value;
  });
  return obj;
};

const EditDetailButton = ({ auth, type = 'file', modelInfo, getModel }: any) => {
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const [loading, setLoading] = useState(false);
  const [show, setShow] = useState(false);
  const [step, setStep] = useState(1);
  const [form] = Form.useForm();

  const scrollToTop = () => {
    const editModalBody = document.querySelector('.editModalBody');
    if (editModalBody) {
      editModalBody.scrollTop = 0;
    }
  };

  const onClose = () => {
    setShow(false);
    setStep(1);
    form.resetFields();
  };

  useEffect(() => {
    if (show)
      setTimeout(() => {
        scrollToTop();
      });
  }, [show, step]);

  const handleBackfill = async () => {
    const {
      alias,
      pathName,
      displayName,
      description,
      accessLevel,
      withSave2db: save2db,
      extend,
      labelList,
      refers,
      expression,
      dataType,
    } = modelInfo;

    const backfillForm = {
      alias,
      pathName,
      displayName,
      description,
      save2db,
      accessLevel,
      extend: extendToArr(extend || []),
      labelNames: (labelList || [])?.map((i: any) => ({ label: i.labelName, value: i.id })),
    };
    if (type === 'file') {
      if (dataType === 3) {
        //实时计算
        const refersRes = await Promise.all(refers.map((e: ReferType) => getInstanceInfo({ id: e.id })));
        const _refers = refers.map((refer: ReferType, i: number) => ({
          ...refer,
          refer: {
            label: refer.path,
            value: refer.id,
          },
          fields: refersRes[i]?.fields?.filter?.(
            (t: any) => !(t.systemField || ['BLOB', 'LBLOB'].includes(t.type))
          ) || [{ name: refer.field }],
        }));

        Object.assign(backfillForm, {
          refers: _refers,
          expression: getExpression(refers, expression),
        });
      }
      if (dataType === 7) {
        const referId = modelInfo?.refers?.length
          ? {
              label: modelInfo?.refers?.[0]?.path,
              value: modelInfo?.refers?.[0]?.id,
            }
          : undefined;

        Object.assign(backfillForm, {
          referId,
        });
      }
    }

    form.setFieldsValue(backfillForm);
  };

  useEffect(() => {
    if (show) handleBackfill();
  }, [show]);

  const onSave = async () => {
    await form.validateFields();
    const info = cloneDeep(form.getFieldsValue(true));
    const { dataType } = modelInfo;
    const {
      save2db,
      accessLevel,
      labelNames,

      refers,
      expression,

      referId,

      ...restInfo
    } = info;
    if (type === 'file') {
      if (dataType === 3) {
        //实时计算-函数计算
        restInfo.refers = refers.map((item: ReferItemType, index: number) => {
          return {
            id: item?.refer?.value,
            field: item.field,
            variableName: `a${index + 1}`,
            variableGroup: 0,
          };
        });
        restInfo.expression = expression ? expression.replace(/\$(.*?)#/g, '$1') : '';
      }

      if (dataType === 7) {
        //模型
        restInfo.refers = referId?.value ? [{ id: referId.value }] : [];
      }
      restInfo.labelNames = labelNames?.map((e: any) => e.label || e.value) || [];
    }
    setLoading(true);
    modifyDetail({
      ...restInfo,
      extend: extendToObj(info?.extend),
      save2db: type === 'file' && ![7].includes(dataType) ? save2db : undefined,
      accessLevel: type === 'file' && [1, 2].includes(dataType) ? accessLevel : undefined,
    })
      .then(() => {
        onClose();
        message.success(formatMessage('uns.editSuccessful'));
        getModel?.(info);
      })
      .finally(() => {
        setLoading(false);
      });
  };

  const referIdChange = (option: any) => {
    if (option.key) {
      getInstanceInfo({ id: option.key }).then((res) => {
        form.setFieldsValue({
          accessLevel: res.accessLevel || 'READ_ONLY',
        });
      });
    } else {
      form.setFieldsValue({
        accessLevel: 'READ_ONLY',
      });
    }
  };

  const formItemOptions = useMemo(() => {
    switch (step) {
      case 1:
        return [
          {
            label: formatMessage('common.name'),
            name: 'pathName',
            properties: {
              disabled: true,
            },
          },
          {
            label: formatMessage('uns.alias'),
            name: 'alias',
            properties: {
              disabled: true,
            },
          },
          {
            label: formatMessage('uns.displayName'),
            name: 'displayName',
            rules: [{ max: 128 }],
          },
          {
            type: 'TextArea',
            label: type === 'file' ? formatMessage('uns.fileDescription') : formatMessage('uns.folderDescription'),
            name: 'description',
            rules: [{ max: 255 }],
          },
          {
            component: (
              <SearchSelect
                apiParams={{ type: 2, normal: true }}
                labelInValue
                onChange={referIdChange}
                onClear={() => form.setFieldsValue({ accessLevel: undefined })}
              />
            ),
            label: formatMessage('uns.referenceTarget'),
            name: 'referId',
            noShowKey: type === 'file' && modelInfo.dataType === 7 ? undefined : 'hidden',
          },
          {
            type: 'Select',
            label: formatMessage('uns.writDownData'),
            name: 'accessLevel',
            initialValue: 'READ_ONLY',
            properties: {
              options: [
                { label: formatMessage('uns.true'), value: 'READ_WRITE' },
                { label: formatMessage('uns.false'), value: 'READ_ONLY' },
              ],
              disabled: modelInfo.mount || (type === 'file' && modelInfo.dataType === 7),
            },
            noShowKey: ![1, 2, 7].includes(modelInfo.dataType) && type === 'file' ? 'hidden' : 'folder',
          },
          {
            type: 'Checkbox',
            name: 'save2db',
            properties: {
              label: formatMessage('uns.persistence'),
              style: { marginLeft: 5 },
              disabled: modelInfo.mount,
            },
            noShowKey: [7].includes(modelInfo.dataType) && type === 'file' ? 'hidden' : 'folder',
            valuePropName: 'checked',
          },
          {
            type: 'TagSelect',
            label: formatMessage('common.label'),
            name: 'labelNames',
            noShowKey: 'folder',
            properties: {
              tagMaxLen: 63,
            },
          },
          {
            type: 'divider',
          },
          {
            render: () => <ExpandedKeyFormList />,
          },
        ]
          .filter((f: any) => (!f.noShowKey || f.noShowKey !== type) && f.noShowKey !== 'hidden')
          .map((e: any) => {
            delete e.noShowKey;
            return e;
          });
      case 2:
        if (type === 'file' && modelInfo.dataType === 3) {
          return [
            {
              render: () => <ExpressionForm apiParams={{ calculationType: 1 }} />,
            },
          ];
        }
        return [];
      default:
        return [];
    }
  }, [type, modelInfo?.dataType, modelInfo?.mount, modelInfo.calculationType, step]);

  const footer = useMemo(() => {
    return (
      <Flex gap="10px" justify="end">
        {step === 1 ? (
          <Button
            style={{
              height: '40px',
              backgroundColor: 'var(--supos-uns-button-color)',
              color: 'var(--supos-text-color)',
            }}
            color="default"
            variant="filled"
            onClick={onClose}
            block
          >
            {formatMessage('common.cancel')}
          </Button>
        ) : (
          <Button
            style={{
              height: '40px',
              backgroundColor: 'var(--supos-uns-button-color)',
              color: 'var(--supos-text-color)',
            }}
            color="default"
            variant="filled"
            onClick={() => setStep?.(step - 1)}
            disabled={loading}
            block
          >
            {formatMessage('common.prev')}
          </Button>
        )}
        {type === 'file' && modelInfo.dataType === 3 && step === 1 ? (
          <Button style={{ height: '40px' }} type="primary" variant="solid" onClick={() => setStep?.(step + 1)} block>
            {formatMessage('common.next')}
          </Button>
        ) : (
          <Button style={{ height: '40px' }} type="primary" variant="solid" onClick={onSave} loading={loading} block>
            {formatMessage('common.save')}
          </Button>
        )}
      </Flex>
    );
  }, [step, modelInfo?.dataType, loading, getModel, type]);

  const renderFrom = useMemo(() => {
    if (!show) return null;
    return (
      <OperationForm
        formConfig={{
          layout: 'vertical',
          labelCol: { span: undefined },
          wrapperCol: { span: undefined },
        }}
        style={{ padding: 0 }}
        form={form}
        formItemOptions={formItemOptions}
        buttonConfig={{ block: true }}
        footer={<span />}
      />
    );
  }, [formItemOptions, show]);

  return (
    <>
      <AuthButton
        auth={auth}
        onClick={() => setShow(true)}
        style={{ border: '1px solid #C6C6C6', background: 'var(--supos-uns-button-color)' }}
        icon={
          <Icon
            data-button-auth={auth}
            component={FileEdit}
            style={{
              fontSize: 16,
              color: 'var(--supos-text-color)',
            }}
          />
        }
      />
      <ProModal
        title={
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <span>{formatMessage('uns.editDetails')}</span>
          </div>
        }
        onCancel={onClose}
        open={show}
        size="xs"
        afterClose={() => {
          form.resetFields();
        }}
        styles={{
          content: { padding: 0 },
          header: { padding: '20px 24px 10px', margin: 0 },
          body: { padding: '0 24px 0', margin: 0, maxHeight: 'calc(80vh - 122px)', overflowY: 'auto' },
          footer: { padding: '0 24px 20px' },
        }}
        footer={footer}
        classNames={{ body: 'editModalBody' }}
        destroyOnHidden
      >
        {renderFrom}
      </ProModal>
    </>
  );
};

export default EditDetailButton;
