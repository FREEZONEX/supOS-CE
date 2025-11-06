import { useCallback, useState } from 'react';
import { Button, Form, App, Input } from 'antd';
import { addLabel } from '@/apis/inter-api/uns';
import { useTranslate } from '@/hooks';

import type { InitTreeDataFnType, UnsTreeNode } from '@/pages/uns/types';
import type { TreeStoreActions } from '../../store/types';
import ProModal from '@/components/pro-modal';

export interface LabelModalProps {
  successCallBack: InitTreeDataFnType;
  changeCurrentPath: (node: UnsTreeNode) => void;
  setTreeMap: TreeStoreActions['setTreeMap'];
  scrollTreeNode: (id: string) => void;
}

const Module = ({ successCallBack, changeCurrentPath, setTreeMap, scrollTreeNode }: LabelModalProps) => {
  const formatMessage = useTranslate();
  const [form] = Form.useForm();
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const { message } = App.useApp();

  const setModalOpen = useCallback(() => {
    setOpen(true);
  }, [setOpen]);

  const close = () => {
    setOpen(false);
    form.resetFields();
    setLoading(false);
  };

  const confirm = async () => {
    const formData = await form.validateFields();
    if (formData) {
      setLoading(true);
      console.log(formData);
      const { name } = formData;
      addLabel(name)
        .then((data: any) => {
          successCallBack({}, () => {
            changeCurrentPath({ key: data?.id, id: data?.id, pathType: 9 });
            setTreeMap(false);
            scrollTreeNode(data?.id);
          });
          message.success(formatMessage('uns.newSuccessfullyAdded'));
          setLoading(false);
          close();
        })
        .catch(() => {
          setLoading(false);
        });
    }
  };

  // 自定义校验空格
  const validateTrim = (_: any, value: string) => {
    if (value && value.trim() === '') {
      return Promise.reject(new Error(formatMessage('common.prohibitSpacesTip')));
    }
    return Promise.resolve();
  };

  const Dom = (
    <ProModal className="labelModalWrap" open={open} onCancel={close} title={formatMessage('uns.newLabel')} size="xxs">
      <Form
        name="addLabelForm"
        form={form}
        colon={false}
        initialValues={{
          withFlow: false,
          withDashboard: true,
        }}
        disabled={loading}
      >
        <Form.Item
          label={formatMessage('common.name')}
          name="name"
          rules={[
            { required: true },
            { max: 63 },
            { validator: validateTrim },
            { pattern: /^[\u4e00-\u9fa5a-zA-Z0-9_-]+$/, message: formatMessage('uns.nameFormat') },
          ]}
        >
          <Input />
        </Form.Item>
      </Form>
      <Button
        className="labelConfirm"
        color="primary"
        variant="solid"
        onClick={confirm}
        block
        style={{ marginTop: '10px' }}
        loading={loading}
        disabled={loading}
      >
        {formatMessage('common.save')}
      </Button>
    </ProModal>
  );
  return {
    LabelModal: Dom,
    setLabelOpen: setModalOpen,
  };
};
export default Module;
