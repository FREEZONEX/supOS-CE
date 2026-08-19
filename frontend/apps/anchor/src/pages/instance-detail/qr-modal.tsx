import { Button, Checkbox, Modal, Radio, Spin, message as antdMessage } from 'antd';
import { Download, Launch } from '@carbon/icons-react';
import QRCode from 'qrcode';
import { useCallback, useEffect, useState } from 'react';
import { qrConfigUrl } from '../../api/instances';
import { getScene, listScenes, sceneQrConfigUrl } from '../../api/scenes';
import { t } from '../../i18n';

// 生成二维码弹窗（对齐云端 instance-qr-modal 两步式）：
// 第一步选分享模式（分享模型 / 分享场景——列出已包含该实例的场景单选）；
// 第二步展示二维码 + 下载 PNG + 复制链接。
// 二维码内容 = <origin>/anchor/viewer?configUrl=<免登录配置接口>，手机须与服务器同局域网。
type ShareMode = 'model' | 'scene';

interface SceneOption {
  id: number;
  name: string;
}

export default function InstanceQrModal({
  open,
  onClose,
  instance,
}: {
  open: boolean;
  onClose: () => void;
  instance: { id: number; name: string } | null;
}) {
  const [messageApi, contextHolder] = antdMessage.useMessage();
  const [step, setStep] = useState<'config' | 'preview'>('config');
  const [shareMode, setShareMode] = useState<ShareMode>('model');
  const [qrDataUrl, setQrDataUrl] = useState('');
  const [shareUrl, setShareUrl] = useState('');
  const [loading, setLoading] = useState(false);
  const [sceneLoading, setSceneLoading] = useState(false);
  const [sceneOptions, setSceneOptions] = useState<SceneOption[]>([]);
  const [selectedSceneId, setSelectedSceneId] = useState<number | null>(null);

  useEffect(() => {
    if (!open) return;
    setStep('config');
    setShareMode('model');
    setQrDataUrl('');
    setShareUrl('');
    setSelectedSceneId(null);
  }, [open, instance?.id]);

  // 加载包含该实例的场景（列表不含条目，需逐个取详情后本地过滤，场景数量级小可接受）。
  // 依赖按 instance.id：父组件高频重渲染（MQTT 实时消息）传新对象不应重启加载
  const instanceId = instance?.id;
  useEffect(() => {
    if (!open || !instanceId) {
      setSceneOptions([]);
      return;
    }
    let cancelled = false;
    setSceneLoading(true);
    listScenes({ page: 1, size: 200 })
      .then((result) => Promise.all((result.list || []).map((scene) => getScene(scene.id))))
      .then((scenes) => {
        if (cancelled) return;
        const matched = scenes
          .filter((scene) => (scene.items || []).some((item) => item.instanceId === instanceId))
          .map((scene) => ({ id: scene.id, name: scene.name }));
        setSceneOptions(matched);
        setSelectedSceneId(matched[0]?.id ?? null);
      })
      .catch(() => {
        if (!cancelled) setSceneOptions([]);
      })
      .finally(() => {
        if (!cancelled) setSceneLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [open, instanceId]);

  const generate = useCallback(async () => {
    if (!instance) return;
    const origin = window.location.origin;
    const configUrl =
      shareMode === 'scene' && selectedSceneId
        ? `${origin}${sceneQrConfigUrl(selectedSceneId)}`
        : `${origin}${qrConfigUrl(instance.id)}`;
    const focus = shareMode === 'scene' ? `&focusInstanceId=${instance.id}` : '';
    const url = `${origin}/anchor/viewer?configUrl=${encodeURIComponent(configUrl)}${focus}`;
    setLoading(true);
    try {
      const dataUrl = await QRCode.toDataURL(url, { width: 256, margin: 2, errorCorrectionLevel: 'M' });
      setQrDataUrl(dataUrl);
      setShareUrl(url);
      setStep('preview');
    } catch (e) {
      console.error('[anchor-qr] generate failed:', e);
      messageApi.error(t('instance.qr.genErr'));
    } finally {
      setLoading(false);
    }
  }, [instance, shareMode, selectedSceneId, messageApi]);

  const download = useCallback(() => {
    if (!qrDataUrl || !instance) return;
    const link = document.createElement('a');
    link.href = qrDataUrl;
    link.download = `${instance.name || 'instance'}-qr-code.png`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }, [instance, qrDataUrl]);

  const copyLink = useCallback(() => {
    if (!shareUrl) return;
    navigator.clipboard
      .writeText(shareUrl)
      .then(() => messageApi.success(t('instance.qr.linkCopied')))
      .catch(() => messageApi.error(t('instance.qr.copyFailed')));
  }, [shareUrl, messageApi]);

  const selectionCard = (mode: ShareMode, title: string, description: string, content?: React.ReactNode) => {
    const selected = shareMode === mode;
    return (
      <div
        role="button"
        tabIndex={0}
        className="flex min-h-[160px] w-full cursor-pointer rounded p-4 text-left transition-colors"
        style={{
          border: `1px solid ${selected ? 'var(--ui-theme-color)' : 'var(--ui-header-splitter-color, #e0e0e0)'}`,
          background: selected ? 'var(--ui-primary-bg)' : 'var(--ui-bg-color)',
        }}
        onClick={() => setShareMode(mode)}
        onKeyDown={(e) => {
          if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault();
            setShareMode(mode);
          }
        }}
      >
        <div className="flex w-full flex-col gap-4">
          <div className="flex items-start gap-2">
            <Radio checked={selected} className="mt-0.5 shrink-0" />
            <div>
              <div className="text-sm font-medium" style={{ color: 'var(--ui-text-color)' }}>
                {title}
              </div>
              <div className="mt-1 text-xs leading-5" style={{ color: 'var(--ui-description-card-color)' }}>
                {description}
              </div>
            </div>
          </div>
          {content}
        </div>
      </div>
    );
  };

  return (
    <Modal
      title={step === 'config' ? t('instance.qr.createTitle') : t('instance.qr.viewTitle')}
      open={open}
      onCancel={onClose}
      width={720}
      footer={
        step === 'preview' && qrDataUrl ? (
          <div className="flex items-center justify-end gap-2">
            <Button icon={<Download size={16} />} onClick={download} className="min-w-40">
              {t('instance.qr.download')}
            </Button>
            <Button type="primary" icon={<Launch size={16} />} onClick={copyLink} className="min-w-40">
              {t('instance.qr.copyLink')}
            </Button>
          </div>
        ) : null
      }
    >
      {contextHolder}
      {step === 'config' ? (
        <div className="space-y-5 py-2">
          <div className="grid grid-cols-2 gap-4">
            {selectionCard('model', t('instance.qr.shareModelTitle'), t('instance.qr.shareModelDesc'))}
            {selectionCard(
              'scene',
              t('instance.qr.shareSceneTitle'),
              t('instance.qr.shareSceneDesc'),
              <div className="space-y-2 pl-6">
                {sceneLoading ? (
                  <div className="text-xs" style={{ color: 'var(--ui-description-card-color)' }}>
                    {t('instance.qr.sceneLoading')}
                  </div>
                ) : sceneOptions.length > 0 ? (
                  sceneOptions.map((scene) => (
                    <label key={scene.id} className="flex items-center gap-2 text-sm">
                      <Checkbox
                        checked={selectedSceneId === scene.id}
                        onChange={() => {
                          setShareMode('scene');
                          setSelectedSceneId(scene.id);
                        }}
                      />
                      <span style={{ color: 'var(--ui-text-color)' }}>{scene.name}</span>
                    </label>
                  ))
                ) : (
                  <div className="text-xs" style={{ color: 'var(--ui-description-card-color)' }}>
                    {t('instance.qr.sceneNoModel')}
                  </div>
                )}
              </div>
            )}
          </div>
          <div className="flex justify-end">
            <Button
              type="primary"
              loading={loading}
              disabled={shareMode === 'scene' && !selectedSceneId}
              onClick={() => void generate()}
            >
              {t('instance.qr.generate')}
            </Button>
          </div>
        </div>
      ) : (
        <div className="flex flex-col items-center gap-4 py-4">
          {loading ? (
            <div className="flex h-64 w-64 items-center justify-center">
              <Spin tip={t('instance.qr.genLoading')} />
            </div>
          ) : (
            <img
              width={256}
              height={256}
              src={qrDataUrl}
              alt={instance?.name || 'QR Code'}
              className="rounded-lg shadow-sm"
              style={{ border: '1px solid var(--ui-header-splitter-color, #e0e0e0)' }}
            />
          )}
        </div>
      )}
    </Modal>
  );
}
