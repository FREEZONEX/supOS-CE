import { useState, useEffect, type FC, type CSSProperties, useRef } from 'react';
import { getLabelDetail, updateLabel, makeSingleLabel, updateLabelSubscribe } from '@/apis/inter-api/uns';
import { useTranslate } from '@/hooks';
import { Collapse, theme, Form, Flex, Button, message, Typography } from 'antd';
import { CaretRight, Copy, Tag } from '@carbon/icons-react';
import { ButtonPermission } from '@/common-types/button-permission.ts';
import SearchSelect from '@/pages/uns/components/use-create-modal/components/SearchSelect';
import Icon from '@ant-design/icons';
import FileList from './FileList';
import type { InitTreeDataFnType, UnsTreeNode } from '@/pages/uns/types';
import { AuthButton, AuthWrapper } from '@/components/auth';
import ProModal from '@/components/pro-modal';
import FileEdit from '@/components/svg-components/FileEdit';
import { hasPermission } from '@/utils/auth';
import Subscribe from '@/pages/uns/components/subscribe';
import ComDetailList from '@/components/com-detail-list';
import useClipboard from '@/hooks/useClipboard.ts';
import CustomParagraph from '@/components/custom-paragraph';
import { formatTimestamp } from '@/utils/format.ts';
import styles from './index.module.scss';

const { Title } = Typography;

const panelStyle: CSSProperties = {
  background: 'val(--supos-bg-color)',
  border: 'none',
};

export interface LabelDetailProps {
  currentNode: UnsTreeNode;
  initTreeData: InitTreeDataFnType;
  handleDelete: (node: UnsTreeNode) => void;
}

const Module: FC<LabelDetailProps> = (props) => {
  const {
    currentNode: { id },
    initTreeData,
    handleDelete,
  } = props;
  const formatMessage = useTranslate();
  const [activeList, setActiveList] = useState<string[]>(['fileList']);
  const { token } = theme.useToken();
  const [topicTitle, setTopicTitle] = useState('');
  const [modelInfo, setModelInfo] = useState<{ [key: string]: any }>({});
  const [isLabelVisible, setIsLabelVisible] = useState(false);
  const [loading, setLoading] = useState(false);

  const [form] = Form.useForm();
  const modelInfoRef = useRef(modelInfo);
  const fileListRef = useRef<any>(null);

  useEffect(() => {
    modelInfoRef.current = modelInfo;
  }, [modelInfo]);

  const getModel = (id: string) => {
    getLabelDetail(id)
      .then((data) => {
        setModelInfo(data);
        setTopicTitle(data?.labelName);
      })
      .catch((error) => {
        console.log(error);
      });
  };

  useEffect(() => {
    if (id) {
      getModel(id as string);
    }
  }, [id]);

  const onSave = async () => {
    const values = await form.validateFields();
    if (values) {
      setLoading(true);
      makeSingleLabel(values.unsId, id as string)
        .then(() => {
          message.success(formatMessage('common.optsuccess'));
          fileListRef.current?.getList?.(id);
          setIsLabelVisible(false);
          form.resetFields();
          setLoading(false);
        })
        .catch(() => {
          setLoading(false);
        });
    }
  };
  const handleAdd = () => {
    setIsLabelVisible(true);
  };

  const onDeleteHandle = () => {
    handleDelete({ key: '', id, pathType: 7 });
  };
  const { copy } = useClipboard();

  const items = [
    {
      key: 'detail',
      label: <span>{formatMessage('common.detail')}</span>,
      children: (
        <ComDetailList
          list={[
            {
              label: 'Topic',
              key: 'topic',
              hide: !modelInfo?.subscribeEnable,
              render: () => {
                return (
                  <div>
                    {modelInfo.topic}
                    <span
                      style={{ marginLeft: '5px', verticalAlign: 'sub', cursor: 'pointer' }}
                      onClick={() => copy(modelInfo.topic)}
                      title={formatMessage('common.copy')}
                    >
                      <Copy />
                    </span>
                  </div>
                );
              },
            },
            {
              label: formatMessage('uns.alias'),
              key: 'labelName',
            },
            {
              label: formatMessage('uns.subscriptionFrequency'),
              key: 'subscribeFrequency',
              hide: !modelInfo?.subscribeEnable,
              render: () => {
                return (
                  <CustomParagraph
                    className={styles.paragraph}
                    value={modelInfo.subscribeFrequency}
                    onChange={(value) => {
                      updateLabelSubscribe({ id, ...value }).then(() => {
                        message.success(formatMessage('uns.editSuccessful'));
                        getModel(id as string);
                      });
                    }}
                  />
                );
              },
            },
            {
              label: formatMessage('common.creationTime'),
              key: 'createTime',
              render: (item) => formatTimestamp(item),
            },
          ]}
          data={modelInfo}
        />
      ),
      style: panelStyle,
    },
    {
      key: 'fileList',
      label: formatMessage('common.fileList'),
      children: <FileList ref={fileListRef} labelId={id as string} />,
      style: panelStyle,
      extra: (
        <AuthButton
          auth={ButtonPermission['uns.labelDetail']}
          onClick={handleAdd}
          style={{
            border: '1px solid #C6C6C6',
            background: 'var(--supos-uns-button-color)',
            color: 'var(--supos-text-color)',
          }}
        >
          +
        </AuthButton>
      ),
    },
  ];
  const handleChangeSubscribe = async (enable: boolean, frequency?: string) => {
    await updateLabelSubscribe({ id, enable, frequency });
    getModel(id as string);
    message.success(enable ? formatMessage('uns.subscribeSuccessful') : formatMessage('uns.unsubscribeSuccessful'));
  };
  return (
    <div className="topicDetailWrap">
      <div className="topicDetailContent">
        <Flex className="detailTitle" gap={8} justify="flex-start" align="center">
          <Tag style={{ transform: 'rotate(90deg)' }} size={20} />
          <Title
            level={2}
            style={{ margin: 0, width: '100%', insetInlineStart: 0 }}
            editable={
              hasPermission(ButtonPermission['uns.labelDetail'])
                ? {
                    icon: (
                      <Icon
                        data-button-auth={ButtonPermission['uns.labelDetail']}
                        component={FileEdit}
                        style={{
                          fontSize: 25,
                          color: '#5A5A5A',
                          marginLeft: 5,
                        }}
                      />
                    ),
                    onChange: (val) => {
                      if (topicTitle === val || !val || val.trim() === '') return;
                      if (val.length > 63) {
                        return message.error(
                          formatMessage('uns.labelMaxLength', { label: formatMessage('common.name'), length: 63 })
                        );
                      }
                      const reg = /^[\u4e00-\u9fa5a-zA-Z0-9_-]+$/;
                      if (!reg.test(val)) {
                        return message.error(formatMessage('uns.nameFormat'));
                      }
                      updateLabel({
                        id,
                        labelName: val,
                      }).then(() => {
                        message.success(formatMessage('uns.editSuccessful'));
                        getModel(id as string);
                        initTreeData?.({});
                      });
                    },
                  }
                : false
            }
          >
            {topicTitle}
          </Title>
          <Subscribe
            hidden
            showModal
            value={modelInfo.subscribeEnable}
            topic={modelInfo.topic}
            subscribeFrequency={modelInfo.subscribeFrequency}
            onChange={handleChangeSubscribe}
          />
        </Flex>
        <div className="tableWrap">
          <Collapse
            bordered={false}
            collapsible="header"
            activeKey={activeList}
            onChange={(even) => setActiveList(even)}
            expandIcon={({ isActive }) => (
              <CaretRight
                size={20}
                style={{
                  rotate: isActive ? '90deg' : '0deg',
                  transition: '200ms',
                }}
              />
            )}
            items={items}
            style={{ background: token.colorBgContainer }}
          />
        </div>
        <AuthWrapper auth={ButtonPermission['uns.labelDetail']}>
          <div className="deleteBtnWrap" style={{ marginTop: 0 }}>
            <Button
              type="primary"
              style={{
                width: '100px',
                fontWeight: 'bold',
              }}
              onClick={onDeleteHandle}
            >
              {formatMessage('common.delete')}
            </Button>
          </div>
        </AuthWrapper>
      </div>
      <ProModal
        title={formatMessage('uns.addFile')}
        className="labelModalWrap"
        open={isLabelVisible}
        onCancel={() => {
          setIsLabelVisible(false);
          form.resetFields();
        }}
        size="xxs"
      >
        <Form colon={false} name="labelForm" disabled={loading} form={form}>
          <Form.Item
            label={formatMessage('uns.position')}
            name="unsId"
            style={{ marginBottom: 15 }}
            rules={[
              {
                required: true,
                message: formatMessage('rule.required'),
              },
            ]}
          >
            <SearchSelect
              popupMatchSelectWidth={400}
              apiParams={{ type: 2 }}
              disabledIds={modelInfo?.fileVoList?.map((item: any) => item.unsId) || []}
            />
          </Form.Item>
        </Form>
        <div style={{ marginTop: '20px' }}>
          <Button
            style={{ width: '48%', marginRight: '4%' }}
            size="large"
            color="default"
            variant="filled"
            disabled={loading}
            onClick={() => {
              setIsLabelVisible(false);
              form.resetFields();
              setLoading(false);
            }}
          >
            {formatMessage('common.cancel')}
          </Button>

          <Button
            className="labelConfirm"
            size="large"
            style={{ width: '48%' }}
            onClick={onSave}
            color="primary"
            loading={loading}
            variant="solid"
          >
            {formatMessage('common.save')}
          </Button>
        </div>
      </ProModal>
    </div>
  );
};
export default Module;
