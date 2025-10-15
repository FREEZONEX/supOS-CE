import { Modal } from 'antd';
import { useState } from 'react';

// 图片渲染 + 点击放大预览，主要用于ChatBot中渲染图片。
const ImagePreview = (props: any) => {
  const [modalOpen, setModalOpen] = useState<boolean>(false);

  return (
    <>
      <img
        {...props}
        style={{ width: '100%', cursor: 'pointer' }}
        onClick={() => {
          setModalOpen(true);
        }}
      />
      <Modal
        width={1000}
        forceRender
        centered
        open={modalOpen}
        footer={null}
        onCancel={() => {
          setModalOpen(false);
        }}
        styles={{
          body: {
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            padding: '30px 24px',
          },
        }}
      >
        <img {...props} style={{ width: '100%' }} />
      </Modal>
    </>
  );
};

export default ImagePreview;
