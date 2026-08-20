import { getProjectAppDetail, replaceProjectApp, updateProjectApp, uploadAttachment } from '@/apis/core-api';
import ProModal from '@/components/pro-modal';
import { useTranslate } from '@/hooks';
import type { AppItem } from '@/pages/project/types';
import { readAppBundleContents, type AppBundleContents, type AppBundleMeta } from '@/utils/app-bundle-meta';
import { projectAppMediaUrl, uploadProjectAppMedia } from '@/utils/project-app-media';
import { Alert, App, Button, Flex, Input, Spin, Upload, type UploadFile, type UploadProps } from 'antd';
import type { RcFile } from 'antd/es/upload';
import { forwardRef, useCallback, useEffect, useImperativeHandle, useRef, useState, type MouseEvent } from 'react';
import { useParams } from 'react-router';
import '@/pages/uns/components/import-modal/index.scss';
import AppBundleUploadContent from './AppBundleUploadContent';
import bundleStyles from './AppBundleModal.module.scss';
import AppMediaFields, { type AppMediaKind } from './AppMediaFields';
import styles from './ReplaceAppModal.module.scss';

const { Dragger } = Upload;
const MAX_IMPORT_FILE_SIZE = 500 * 1024 * 1024;
const DEPLOY_POLL_INTERVAL_MS = 3000;
const MAX_DEPLOY_POLLS = 300;
const MAX_CONSECUTIVE_POLL_ERRORS = 5;

export interface ReplaceAppModalRef {
  onOpen: (app: AppItem) => void;
  onClose: (force?: boolean) => void;
}

export interface ReplaceAppModalProps {
  refreshRequest?: () => void;
}

const ReplaceAppModal = forwardRef<ReplaceAppModalRef, ReplaceAppModalProps>(({ refreshRequest }, ref) => {
  const { projectId } = useParams<{ projectId: string }>();
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const [visible, setVisible] = useState(false);
  const [targetApp, setTargetApp] = useState<AppItem>();
  const [bundleMeta, setBundleMeta] = useState<AppBundleMeta>();
  // 名称与描述允许用户在替换前编辑（与导入弹窗一致），初值取包内 meta、回退目标应用信息
  const [replacementName, setReplacementName] = useState('');
  const [replacementDescription, setReplacementDescription] = useState('');
  const [fileList, setFileList] = useState<UploadFile[]>([]);
  const [uploadedFileInfo, setUploadedFileInfo] = useState<Record<string, unknown> | null>(null);
  const [uploading, setUploading] = useState(false);
  const [saveLoading, setSaveLoading] = useState(false);
  const [uploadError, setUploadError] = useState('');
  const [appIdMismatch, setAppIdMismatch] = useState(false);
  const [versionNotIncreased, setVersionNotIncreased] = useState(false);
  const [deployedAppId, setDeployedAppId] = useState<string>();
  const [deployStatus, setDeployStatus] = useState<string>();
  const [deployError, setDeployError] = useState<string>();
  const [iconAssetId, setIconAssetId] = useState<number>();
  const [coverAssetId, setCoverAssetId] = useState<number>();
  const [iconUrl, setIconUrl] = useState<string>();
  const [coverUrl, setCoverUrl] = useState<string>();
  const [uploadingMediaKind, setUploadingMediaKind] = useState<AppMediaKind>();
  const uploadRootRef = useRef<HTMLDivElement>(null);
  const pollTimerRef = useRef<ReturnType<typeof setInterval>>();
  const pollCountRef = useRef(0);

  const clearDeployPoll = useCallback(() => {
    if (pollTimerRef.current) {
      clearInterval(pollTimerRef.current);
      pollTimerRef.current = undefined;
    }
    pollCountRef.current = 0;
  }, []);

  const resetState = useCallback(() => {
    clearDeployPoll();
    setBundleMeta(undefined);
    setReplacementName('');
    setReplacementDescription('');
    setFileList([]);
    setUploadedFileInfo(null);
    setUploading(false);
    setSaveLoading(false);
    setUploadError('');
    setAppIdMismatch(false);
    setVersionNotIncreased(false);
    setDeployedAppId(undefined);
    setDeployStatus(undefined);
    setDeployError(undefined);
    setIconAssetId(undefined);
    setCoverAssetId(undefined);
    setIconUrl(undefined);
    setCoverUrl(undefined);
    setUploadingMediaKind(undefined);
  }, [clearDeployPoll]);

  const onOpen = (app: AppItem) => {
    resetState();
    setTargetApp(app);
    setVisible(true);
  };

  const onClose = (force?: boolean) => {
    if (!force && deployedAppId) {
      setVisible(false);
      return;
    }
    resetState();
    setTargetApp(undefined);
    setVisible(false);
  };

  const onCloseRef = useRef(onClose);
  onCloseRef.current = onClose;

  const getUploadedItem = (response: any) => {
    const list = response?.list ?? response?.data?.list;
    return Array.isArray(list) ? list[0] : null;
  };

  const getUploadErrorText = (error: any) => {
    const rawMessage =
      typeof error?.msg === 'string' ? error.msg : typeof error?.message === 'string' ? error.message : '';
    if (/^timeout of \d+ms exceeded$/i.test(rawMessage)) {
      return formatMessage('common.requestTimeout');
    }
    if (rawMessage) {
      return formatMessage(rawMessage, {}, rawMessage);
    }
    return formatMessage('common.serverBusy');
  };

  const updateUploadFileStatus = (file: RcFile, status: UploadFile['status']) => {
    setFileList([
      {
        uid: file.uid,
        name: file.name,
        size: file.size,
        type: file.type,
        status,
      },
    ]);
  };

  const validateBundleIdentity = useCallback(
    async (file: RcFile) => {
      if (!targetApp || !projectId) {
        throw new Error('invalid-target');
      }

      const contents = await readAppBundleContents(file);
      const { meta } = contents;
      // 与后端 validateReplaceIdentity 对齐：身份只按 appId 判定，包内 appId 对比目标应用记录的来源 App ID，
      // 不一致时硬阻断；workspaceId/projectId 仅记录不参与校验。
      // 目标应用未记录来源身份（手工创建/旧数据）时不预校验，交给后端裁决。
      setAppIdMismatch(Boolean(targetApp.sourceAppId) && meta.appId !== targetApp.sourceAppId);
      // 版本检查：参考 SaaS 规则（版本为递增整数），包内版本低于或等于当前版本时仅警告、不阻断。
      // 当前应用版本缺失或无法解析为数字时跳过比较。
      const currentVersion = Number(targetApp.version);
      const hasCurrentVersion = targetApp.version !== '' && !Number.isNaN(currentVersion);
      setVersionNotIncreased(meta.version !== undefined && hasCurrentVersion && meta.version <= currentVersion);
      setBundleMeta(meta);
      return contents;
    },
    [projectId, targetApp]
  );

  const uploadMediaFile = async (kind: AppMediaKind, file: File) => {
    if (!targetApp) {
      throw new Error('invalid-target');
    }
    setUploadingMediaKind(kind);
    try {
      const assetId = await uploadProjectAppMedia(file, targetApp.appId);
      if (kind === 'icon') {
        setIconAssetId(assetId);
        setIconUrl(projectAppMediaUrl(assetId));
      } else {
        setCoverAssetId(assetId);
        setCoverUrl(projectAppMediaUrl(assetId));
      }
    } catch (error) {
      setUploadError(getUploadErrorText(error));
    } finally {
      setUploadingMediaKind(undefined);
    }
  };

  const uploadBundleMedia = async (contents: AppBundleContents) => {
    if (contents.icon) {
      await uploadMediaFile('icon', contents.icon);
    }
    if (contents.cover) {
      await uploadMediaFile('cover', contents.cover);
    }
  };

  const doUpload = async (file: RcFile) => {
    setUploading(true);
    setUploadedFileInfo(null);
    setUploadError('');
    try {
      const response = await uploadAttachment([{ value: file, name: 'files', fileName: file.name }], {
        alias: '__templates__',
        ownerType: 'projectAppReplacement',
        ownerId: projectId,
        source: 'project',
      });
      const uploadedItem = getUploadedItem(response);
      if (!uploadedItem) {
        throw new Error('upload-failed');
      }
      setUploadedFileInfo(uploadedItem);
      updateUploadFileStatus(file, 'done');
      return true;
    } catch (error) {
      setUploadedFileInfo(null);
      setUploadError(getUploadErrorText(error));
      updateUploadFileStatus(file, 'error');
      return false;
    } finally {
      setUploading(false);
    }
  };

  const beforeUpload: UploadProps['beforeUpload'] = async (file) => {
    const fileType = file.name.split('.').pop()?.toLowerCase();
    if (fileType !== 'zip') {
      message.warning(formatMessage('common.fileFormatType', { fileType: '.zip' }));
      return Upload.LIST_IGNORE;
    }
    if (file.size > MAX_IMPORT_FILE_SIZE) {
      message.warning(formatMessage('common.fileSizeMax', { size: '500MB' }));
      return Upload.LIST_IGNORE;
    }

    setUploadError('');
    setBundleMeta(undefined);
    setUploadedFileInfo(null);
    setAppIdMismatch(false);
    setVersionNotIncreased(false);
    setIconAssetId(undefined);
    setCoverAssetId(undefined);
    setIconUrl(undefined);
    setCoverUrl(undefined);
    updateUploadFileStatus(file, 'uploading');
    let contents: AppBundleContents;
    try {
      contents = await validateBundleIdentity(file);
    } catch (error) {
      const reason = error instanceof Error ? error.message : '';
      setUploadError(
        reason === 'meta-missing'
          ? formatMessage(
              'project.replace.metaMissing',
              {},
              'This app package is missing required information and cannot be used for replacement.'
            )
          : formatMessage(
              'project.replace.invalidBundle',
              {},
              'This app package is invalid. Select a valid .zip app package.'
            )
      );
      updateUploadFileStatus(file, 'error');
      return false;
    }
    const uploaded = await doUpload(file);
    if (uploaded) {
      try {
        await uploadBundleMedia(contents);
      } catch (error) {
        setUploadError(getUploadErrorText(error));
      }
    }
    return false;
  };

  const openFileDialog = () => {
    if (!uploading && !saveLoading) {
      uploadRootRef.current?.querySelector<HTMLInputElement>('input[type="file"]')?.click();
    }
  };

  const handleUploadClick = (event: MouseEvent<HTMLDivElement>) => {
    const target = event.target as Element | null;
    if (!target?.closest('.ant-upload-list')) {
      openFileDialog();
    }
  };

  // 包内 meta 或目标应用变化时同步名称/描述初值（按后端上限截断，避免受控值超长标红）；
  // 用户编辑后仅在重新上传包/换目标时被覆盖
  useEffect(() => {
    setReplacementName((bundleMeta?.name || targetApp?.name || '').slice(0, 64));
    setReplacementDescription((bundleMeta?.description || targetApp?.description || '').slice(0, 200));
  }, [bundleMeta, targetApp]);

  const onReplace = async () => {
    if (!uploadedFileInfo || !bundleMeta) {
      message.warning(formatMessage('project.import.uploadRequired'));
      return;
    }
    if (!projectId || !targetApp) {
      message.warning(formatMessage('common.serverBusy'));
      return;
    }
    // appId 不一致硬阻断（与后端 validateReplaceIdentity 对齐），不允许确认放行
    if (appIdMismatch) {
      return;
    }

    setSaveLoading(true);
    setDeployError(undefined);
    try {
      const response = await replaceProjectApp(projectId, String(targetApp.appId), {
        ...uploadedFileInfo,
        description: replacementDescription.trim(),
      });
      setDeployedAppId(String(response?.appId || targetApp.appId));
    } catch (error) {
      setUploadError(getUploadErrorText(error));
      setSaveLoading(false);
    }
  };

  useEffect(() => {
    if (!deployedAppId || !projectId || !targetApp) {
      return;
    }

    pollCountRef.current = 0;
    let consecutiveErrors = 0;
    setDeployStatus('deploying');

    const checkDeploy = async () => {
      if (pollCountRef.current >= MAX_DEPLOY_POLLS) {
        // 达到轮询上限：部署结果未知，停止轮询并进入专用 timeout 终态，
        // 提示用户稍后从应用列表确认最终状态（不再轮询、不重复提交）。
        clearDeployPoll();
        setDeployedAppId(undefined);
        setSaveLoading(false);
        setDeployStatus('timeout');
        setDeployError(formatMessage('project.deployTimeout'));
        refreshRequest?.();
        return;
      }
      pollCountRef.current += 1;

      try {
        const detail = await getProjectAppDetail(projectId, String(targetApp.appId));
        consecutiveErrors = 0;
        const status = detail?.status || detail?.data?.status;
        const currentDeployStatus = detail?.deployStatus || detail?.data?.deployStatus;
        const currentDeployError = detail?.deployError || detail?.data?.deployError;
        setDeployStatus(currentDeployStatus || status);

        // deployStatus 优先：先判失败，避免 Replace 回退态（status=Active + deployStatus=failed）误判成功
        if (currentDeployStatus === 'failed' || status === 'Failed') {
          clearDeployPoll();
          setDeployedAppId(undefined);
          setSaveLoading(false);
          setDeployError(currentDeployError || formatMessage('project.deployFailed'));
          refreshRequest?.();
          return;
        }

        if (currentDeployStatus === 'success' || (!currentDeployStatus && status === 'Active')) {
          clearDeployPoll();
          await updateProjectApp(projectId, String(targetApp.appId), {
            name: replacementName.trim() || bundleMeta?.name || targetApp.name,
            description:
              replacementDescription.trim() || (bundleMeta?.description || targetApp.description || '').slice(0, 200),
            iconAssetId: iconAssetId || targetApp.iconAssetId || 0,
            coverAssetId: coverAssetId || targetApp.coverAssetId || 0,
          });
          message.success(formatMessage('project.replace.success', {}, 'App replaced successfully'));
          refreshRequest?.();
          onCloseRef.current(true);
          return;
        }
      } catch {
        consecutiveErrors += 1;
        if (consecutiveErrors >= MAX_CONSECUTIVE_POLL_ERRORS) {
          clearDeployPoll();
          setDeployedAppId(undefined);
          setSaveLoading(false);
          setDeployStatus('failed');
          setDeployError(formatMessage('project.deployFailed'));
          refreshRequest?.();
        }
      }
    };

    void checkDeploy();
    pollTimerRef.current = setInterval(() => void checkDeploy(), DEPLOY_POLL_INTERVAL_MS);
    return clearDeployPoll;
  }, [
    bundleMeta,
    clearDeployPoll,
    coverAssetId,
    deployedAppId,
    formatMessage,
    iconAssetId,
    message,
    projectId,
    refreshRequest,
    replacementDescription,
    replacementName,
    targetApp,
  ]);

  useImperativeHandle(ref, () => ({
    onOpen,
    onClose,
  }));

  const isDeploying = saveLoading && (deployStatus === 'deploying' || deployStatus === 'Deploying');

  return (
    <ProModal
      open={visible}
      onCancel={() => onClose()}
      title={formatMessage('project.replace.title', {}, 'Replace App')}
      className="importModalWrap attachment-upload-modal"
      width={544}
      fullScreenable={false}
      styles={{ body: { maxHeight: '68vh', overflowY: 'auto' } }}
    >
      <div ref={uploadRootRef} className={bundleStyles['upload-section']} onClick={handleUploadClick}>
        <Dragger
          className="uploadWrap"
          action=""
          accept=".zip"
          maxCount={1}
          fileList={fileList}
          disabled={uploading || saveLoading}
          beforeUpload={beforeUpload}
          openFileDialogOnClick={false}
          onRemove={() => {
            setFileList([]);
            setUploadedFileInfo(null);
            setBundleMeta(undefined);
            setUploadError('');
            setAppIdMismatch(false);
            setVersionNotIncreased(false);
            setIconAssetId(undefined);
            setCoverAssetId(undefined);
            setIconUrl(undefined);
            setCoverUrl(undefined);
          }}
        >
          <AppBundleUploadContent />
        </Dragger>
      </div>

      {appIdMismatch ? (
        <Alert
          className={styles['identity-warning']}
          type="error"
          showIcon
          message={formatMessage(
            'project.replace.appIdMismatchBlock',
            { bundleAppId: bundleMeta?.appId ?? '', currentAppId: String(targetApp?.appId ?? '') },
            'The App ID in the package does not match the current App and cannot be replaced.'
          )}
        />
      ) : null}

      {versionNotIncreased ? (
        <Alert
          className={styles['identity-warning']}
          type="warning"
          showIcon
          message={formatMessage(
            'project.replace.versionNotIncreased',
            { bundleVersion: String(bundleMeta?.version ?? ''), currentVersion: targetApp?.version ?? '' },
            'The package version is not higher than the current App version. The version will not increase after replacement.'
          )}
        />
      ) : null}

      {bundleMeta && uploadedFileInfo ? (
        <div className={bundleStyles['details-section']}>
          <div>
            <label className={bundleStyles['field-label']}>{formatMessage('common.name')}</label>
            <Input
              placeholder={formatMessage('apps.namePlaceholder')}
              value={replacementName}
              onChange={(e) => setReplacementName(e.target.value)}
              maxLength={64}
              disabled={uploading || saveLoading}
            />
          </div>
          <div className={bundleStyles['meta-grid']}>
            <div>
              <label className={bundleStyles['field-label']}>
                {formatMessage('project.replace.bundleAppId', {}, 'Bundle App ID')}
              </label>
              <Input value={bundleMeta.appId} disabled />
            </div>
            <div>
              <label className={bundleStyles['field-label']}>
                {formatMessage('project.replace.bundleVersion', {}, 'Bundle Version')}
              </label>
              <Input value={bundleMeta.version !== undefined ? String(bundleMeta.version) : ''} disabled />
            </div>
          </div>
          <div>
            <label className={bundleStyles['field-label']}>{formatMessage('common.description')}</label>
            <Input.TextArea
              placeholder={formatMessage('apps.descriptionPlaceholder')}
              value={replacementDescription}
              onChange={(e) => setReplacementDescription(e.target.value)}
              maxLength={200}
              rows={2}
              disabled={uploading || saveLoading}
            />
          </div>
          <AppMediaFields
            coverUrl={coverUrl}
            iconUrl={iconUrl}
            disabled={uploading || saveLoading || Boolean(uploadingMediaKind)}
            uploadingKind={uploadingMediaKind}
            onSelect={uploadMediaFile}
          />
        </div>
      ) : null}

      {uploadError ? (
        <Alert className={bundleStyles['error-alert']} type="error" showIcon message={uploadError} />
      ) : null}
      {deployError ? (
        <Alert className={bundleStyles['error-alert']} type="error" showIcon message={deployError} />
      ) : null}

      {isDeploying ? (
        <Flex className={bundleStyles['deploy-status']} align="center" gap={8}>
          <Spin size="small" />
          <span>{formatMessage('project.deploying')}</span>
        </Flex>
      ) : null}

      <Button
        className={bundleStyles['submit-button']}
        loading={saveLoading}
        color="primary"
        variant="solid"
        block
        disabled={uploading || Boolean(uploadingMediaKind) || isDeploying || !uploadedFileInfo || appIdMismatch}
        onClick={onReplace}
      >
        {formatMessage('project.replace.confirm', {}, 'Confirm Replace App')}
      </Button>
    </ProModal>
  );
});

ReplaceAppModal.displayName = 'ReplaceAppModal';

export default ReplaceAppModal;
