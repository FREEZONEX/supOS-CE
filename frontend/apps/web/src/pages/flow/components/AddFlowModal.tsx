import { forwardRef, useImperativeHandle, useMemo, useState } from 'react';
import { App, Flex, Form } from 'antd';
import ProModal from '@/components/pro-modal';
import { useTranslate } from '@/hooks';
import OperationForm from '@/components/operation-form';
import ComButton from '@/components/com-button';
import ComRadio from '@/components/com-radio';
import { addFlow as addSourceFlow } from '@/apis/core-api/flow';
import { addFlow as addEventFlow } from '@/apis/core-api/event-flow';
import { ButtonPermission } from '@/common-types/button-permission';
import { hasPermission } from '@/utils/auth';
import { validInputPattern } from '@/utils/pattern';
import type { FlowKind } from '../types';
import styles from './AddFlowModal.module.scss';

export interface AddFlowModalOpenOptions {
  flowKind?: FlowKind;
  groupId?: number;
  lockFlowKind?: boolean;
}

export interface AddFlowModalRef {
  onOpen: (options?: AddFlowModalOpenOptions) => void;
  onClose: () => void;
}

export interface AddFlowModalProps {
  onSuccess?: () => void;
}

const canCreateSource = () => hasPermission(ButtonPermission['SourceFlow.add']);
const canCreateEvent = () => hasPermission(ButtonPermission['EventFlow.add']);

const AddFlowModal = forwardRef<AddFlowModalRef, AddFlowModalProps>(({ onSuccess }, ref) => {
  const [visible, setVisible] = useState(false);
  const [lockFlowKind, setLockFlowKind] = useState(false);
  const formatMessage = useTranslate();
  const [form] = Form.useForm();
  const { message } = App.useApp();

  const flowKindOptions = useMemo(
    () =>
      [
        canCreateSource() ? { label: formatMessage('home.sourceFlow'), value: 'source' as const } : null,
        canCreateEvent() ? { label: formatMessage('home.eventFlow'), value: 'event' as const } : null,
      ].filter(Boolean) as { label: string; value: FlowKind }[],
    [formatMessage]
  );

  const onClose = () => {
    form.resetFields();
    setVisible(false);
    setLockFlowKind(false);
  };

  const onOpen = (options?: AddFlowModalOpenOptions) => {
    const preferredKind =
      options?.flowKind && flowKindOptions.some((item) => item.value === options.flowKind)
        ? options.flowKind
        : flowKindOptions[0]?.value || 'source';

    form.setFieldsValue({
      flowKind: preferredKind,
      groupId: options?.groupId,
      flowName: undefined,
      description: undefined,
    });
    setLockFlowKind(Boolean(options?.lockFlowKind) || flowKindOptions.length <= 1);
    setVisible(true);
  };

  const onSave = async () => {
    const values = await form.validateFields();
    const api = values.flowKind === 'event' ? addEventFlow : addSourceFlow;
    return api({
      flowName: values.flowName,
      description: values.description,
      groupId: values.groupId,
    }).then(() => {
      onClose();
      onSuccess?.();
      message.success(formatMessage('common.optsuccess'));
    });
  };

  useImperativeHandle(ref, () => ({
    onOpen,
    onClose,
  }));

  return (
    <ProModal
      open={visible}
      onCancel={onClose}
      title={formatMessage('common.newFlow')}
      width={500}
      styles={{
        header: {
          marginBottom: 0,
        },
        body: {
          paddingBlockStart: 16,
        },
      }}
    >
      {() => (
        <OperationForm
          form={form}
          onCancel={onClose}
          onSave={onSave}
          formConfig={{
            layout: 'vertical',
            labelCol: { span: 24 },
            wrapperCol: { span: 24 },
            className: `operation-form ${styles.form}`,
          }}
          formItemOptions={[
            {
              name: 'groupId',
              hidden: true,
            },
            {
              name: 'flowName',
              label: formatMessage('common.name'),
              rules: [
                { required: true, message: formatMessage('rule.required') },
                {
                  max: 128,
                  message: formatMessage('uns.labelMaxLength', {
                    label: formatMessage('common.name'),
                    length: 128,
                  }),
                },
                { pattern: validInputPattern, message: formatMessage('rule.flowNameIllegal') },
              ],
              properties: {
                placeholder: formatMessage('common.commonPlaceholder'),
                allowClear: true,
              },
            },
            {
              name: 'description',
              label: formatMessage('uns.description'),
              properties: {
                placeholder: formatMessage('common.commonPlaceholder'),
                allowClear: true,
              },
            },
            ...(lockFlowKind
              ? [
                  {
                    name: 'flowKind',
                    hidden: true,
                  },
                ]
              : [
                  {
                    name: 'flowKind',
                    label: formatMessage('flow.flowType'),
                    rules: [{ required: true, message: formatMessage('rule.required') }],
                    component: (
                      <div className={styles.flowTypeRadio}>
                        <ComRadio options={flowKindOptions} />
                      </div>
                    ),
                  },
                ]),
          ]}
          style={{ padding: 0, overflow: 'visible' }}
          footer={
            <Flex gap="10px" justify="end">
              <ComButton
                color="default"
                variant="filled"
                onClick={onClose}
                title={formatMessage('common.cancel')}
              >
                {formatMessage('common.cancel')}
              </ComButton>
              <ComButton type="primary" variant="solid" onClick={onSave} title={formatMessage('common.save')}>
                {formatMessage('common.save')}
              </ComButton>
            </Flex>
          }
        />
      )}
    </ProModal>
  );
});

AddFlowModal.displayName = 'AddFlowModal';

export default AddFlowModal;
