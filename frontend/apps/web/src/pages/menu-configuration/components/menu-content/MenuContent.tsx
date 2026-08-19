import { App, Button, Form, Input, Tag } from 'antd';
import useTranslate from '@/hooks/useTranslate.ts';
import BasicInfo, { type FormItemType } from './BasicInfo.tsx';
import { useMenuStore } from '../../store/menuStore.tsx';
import { useEffect, useState } from 'react';
import { saveResourceApi } from '@/apis/core-api/resource.ts';
import { ButtonPermission } from '@/common-types/button-permission.ts';
import { AuthButton } from '@/components/auth';
import { WarningFilled } from '@/components/lucide-icon/carbon';
import styles from './MenuContent.module.scss';

const textItem = (name: string, label: string, hidden = false): FormItemType => ({
  formType: 'input',
  label,
  formProps: {
    name,
    hidden,
  },
});

const radioItem = (name: string, label: string, options: { label: string; value: number }[]): FormItemType => ({
  formType: 'radioGroup',
  label,
  formProps: {
    name,
  },
  childProps: {
    options,
  },
});

const MenuContent = () => {
  const formatMessage = useTranslate();
  const [form] = Form.useForm();
  const { message, modal } = App.useApp();
  const name = Form.useWatch('name', form);
  const { selectNode, menuList, requestMenu, setSelectNode, setContentType, contentType } = useMenuStore((state) => ({
    selectNode: state.selectNode,
    menuList: state.menuList,
    requestMenu: state.requestMenu,
    setSelectNode: state.setSelectNode,
    setContentType: state.setContentType,
    contentType: state.contentType,
  }));
  const [configs, setConfigs] = useState<FormItemType[]>([]);
  const [saving, setSaving] = useState(false);
  const [hasUserEdited, setHasUserEdited] = useState(false);
  const isEditable = !selectNode?.coreResourceId || selectNode?.editEnable !== false;
  const levelLabel =
    selectNode?.type === 1
      ? formatMessage('MenuConfiguration.level1')
      : selectNode?.type === 2
        ? formatMessage('MenuConfiguration.level2')
        : '';
  const displayTitle =
    name ||
    (selectNode?.showName
      ? formatMessage(selectNode.showName, undefined, selectNode.showName)
      : selectNode?.showName) ||
    formatMessage('common.new');
  const isAdding = contentType === 'addMenu' || contentType === 'addGroup';

  useEffect(() => {
    form.resetFields();
    const isPage = selectNode?.type === 2;
    const displayName = selectNode?.showName
      ? formatMessage(selectNode.showName, undefined, selectNode.showName)
      : selectNode?.showName;
    setConfigs([
      textItem('name', formatMessage('MenuConfiguration.menuName')),
      textItem('code', formatMessage('MenuConfiguration.menuCode')),
      {
        formType: 'menuIcon',
        label: formatMessage('MenuConfiguration.menuIcon'),
        formProps: {
          name: 'icon',
        },
      },
      ...(isPage
        ? [
            radioItem('urlType', formatMessage('MenuConfiguration.addressType'), [
              { label: formatMessage('MenuConfiguration.internalRoute'), value: 1 },
              { label: formatMessage('MenuConfiguration.externalUrl'), value: 2 },
            ]),
            textItem('url', formatMessage('MenuConfiguration.menuUrl')),
            radioItem('openType', formatMessage('MenuConfiguration.openMode'), [
              { label: formatMessage('MenuConfiguration.openCurrentTab'), value: 0 },
              { label: formatMessage('MenuConfiguration.openNewTab'), value: 1 },
            ]),
          ]
        : []),
      textItem('description', formatMessage('MenuConfiguration.menuDescription')),
    ]);
    if (selectNode) {
      form.setFieldsValue({
        ...selectNode,
        id: selectNode.id,
        name: displayName,
        description: selectNode.showDescription,
        urlType: selectNode.urlType ?? 1,
        openType: selectNode.openType ?? 0,
      });
    }
    setHasUserEdited(false);
  }, [form, formatMessage, isEditable, selectNode]);

  const onSave = async () => {
    if (!selectNode) return;
    const values = await form.validateFields();
    setSaving(true);
    try {
      const saved = await saveResourceApi({
        ...values,
        id: selectNode.coreResourceId || selectNode.id,
        parentId: menuList?.find((item) => item.id === selectNode.parentId)?.coreResourceId || 0,
        resourceKey: selectNode.resourceKey || values.code,
        routePath: values.url,
        enabled: selectNode.enable === false ? 0 : 1,
      });
      await requestMenu();
      setSelectNode({
        ...selectNode,
        ...values,
        coreResourceId: String(saved?.id || selectNode.coreResourceId || ''),
        url: values.url,
        showName: values.name,
        showDescription: values.description,
      });
      message.success(formatMessage('common.optsuccess'));
    } finally {
      setSaving(false);
    }
  };

  const closeAddDetail = () => {
    form.resetFields();
    setHasUserEdited(false);
    setSelectNode(null);
    setContentType(null);
  };

  const onCancel = () => {
    if (!isAdding) {
      form.resetFields();
      return;
    }

    if (!hasUserEdited) {
      closeAddDetail();
      return;
    }

    modal.confirm({
      title: null,
      icon: null,
      width: 420,
      className: 'menu-unsaved-confirm-modal',
      content: (
        <div className="menu-unsaved-confirm-content">
          <WarningFilled size={20} className="menu-unsaved-confirm-icon" aria-hidden />
          <span>{formatMessage('MenuConfiguration.confirmCloseUnsaved')}</span>
        </div>
      ),
      okText: formatMessage('common.confirm'),
      cancelText: formatMessage('common.cancel'),
      onOk: closeAddDetail,
      okButtonProps: {
        title: formatMessage('common.confirm'),
      },
      cancelButtonProps: {
        title: formatMessage('common.cancel'),
      },
    });
  };

  return (
    <div className={styles.menuContent}>
      <div className={styles.header}>
        <div className={styles.headerLeft}>
          <h2 className={styles.title} title={displayTitle}>
            {displayTitle}
          </h2>
          {levelLabel ? (
            <Tag bordered={false} className={styles.levelTag}>
              {levelLabel}
            </Tag>
          ) : null}
          {isEditable && !name ? (
            <Tag bordered={false} color="success" className={styles.newTag}>
              {formatMessage('common.new')}
            </Tag>
          ) : null}
        </div>
        <div className={styles.headerRight}>
          {!isEditable ? (
            <Tag bordered={false} className={styles.readonlyTag}>
              {formatMessage('route.systemReadonly')}
            </Tag>
          ) : (
            <div className={styles.actions}>
              <Button onClick={onCancel}>{formatMessage('common.cancel')}</Button>
              <AuthButton
                auth={ButtonPermission['MenuConfiguration.editMenu']}
                type="primary"
                loading={saving}
                onClick={onSave}
              >
                {formatMessage('common.save')}
              </AuthButton>
            </div>
          )}
        </div>
      </div>
      <div className={styles.body}>
        <Form
          className={styles.form}
          form={form}
          colon={false}
          disabled={!isEditable}
          onValuesChange={() => setHasUserEdited(true)}
        >
          <Form.Item name="coreResourceId" hidden>
            <Input />
          </Form.Item>
          <BasicInfo configs={configs} iconMenuItem={selectNode} iconDisabled={!isEditable} />
        </Form>
      </div>
    </div>
  );
};

export default MenuContent;
