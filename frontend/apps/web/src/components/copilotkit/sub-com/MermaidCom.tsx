import { useEffect, useRef, useState } from 'react';
import mermaid from 'mermaid';
import { Modal } from 'antd';

const Mermaid = ({ code }: { code: string }) => {
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
  const [isInvalid, setIsInvalid] = useState<boolean>(false);

  const graphDivRef = useRef<HTMLDivElement | null>(null);
  const graphPreviewRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    if (code) {
      mermaid
        .parse(code)
        .then(() => {
          setIsInvalid(false);
          mermaid.render(`mermaid-${new Date().valueOf()}`, code).then(({ svg }) => {
            // console.log('mermaid', code, svg);
            if (graphDivRef.current) {
              graphDivRef.current.innerHTML = svg;
            }
          });
        })
        .catch((error) => {
          console.log('error', error);
          setIsInvalid(true);
        });
    }
  }, [code]);

  // useEffect(() => {
  //   // mermaid.contentLoaded();
  // }, []);

  return (
    <div>
      {isInvalid ? (
        <>Chart loading error!</>
      ) : (
        <>
          <Modal
            width={1200}
            // height={1000}
            forceRender
            centered
            open={isModalOpen}
            onOk={() => {
              setIsModalOpen(false);
            }}
            onCancel={() => {
              setIsModalOpen(false);
            }}
          >
            <div
              ref={graphPreviewRef}
              style={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
              }}
            ></div>
          </Modal>
          <div
            ref={graphDivRef}
            onClick={() => {
              if (graphPreviewRef.current) {
                setIsModalOpen(true);
                graphPreviewRef.current.innerHTML = graphDivRef.current?.innerHTML || '';
              }
            }}
          ></div>
        </>
      )}
    </div>
  );
};

export default Mermaid;
