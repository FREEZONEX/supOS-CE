import { Upload as UploadIcon } from '@/components/lucide-icon/carbon';
import { useTranslate } from '@/hooks';

const AppBundleUploadContent = () => {
  const formatMessage = useTranslate();

  return (
    <div className="upload-drag-content">
      <UploadIcon size={32} strokeWidth={1.5} className="upload-drag-icon" />
      <p className="upload-hint-primary">{formatMessage('project.package.upload', {}, 'Upload App Package')}</p>
      <p className="upload-hint-secondary">
        {formatMessage('project.package.supportHint', {}, 'Supports .zip app packages up to 500MB.')}
      </p>
    </div>
  );
};

export default AppBundleUploadContent;
