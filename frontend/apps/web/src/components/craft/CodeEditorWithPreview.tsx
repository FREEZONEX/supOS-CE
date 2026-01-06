import { type FC, useEffect } from 'react';
import './index.scss';
import { usePropsValue } from '@/hooks';
import { Radio } from 'antd';

const svg = `<svg version="1.1" id="Layer_1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px"
	 width="92px" height="92px" viewBox="0 0 92 92" enable-background="new 0 0 92 92" xml:space="preserve">
  <path id="XMLID_32_" d="M46,53c-1.8,0-3.7-0.8-5-2.1c-1.3-1.3-2-3.1-2-4.9c0-1.8,0.8-3.6,2-5c1.3-1.3,3.1-2,5-2
	c1.8,0,3.6,0.8,4.9,2c1.3,1.3,2.1,3.1,2.1,5c0,1.8-0.8,3.6-2.1,4.9C49.6,52.2,47.8,53,46,53z" fill="rgb(141, 141, 141)"/>
</svg>`;

let tooltip: any = null;

type TypeProps = 'iframe' | 'shadow';
interface CodeEditorWithPreviewProps {
  code?: any;
  setComponentIds?: any;
  setShowTuneDialog?: any;
  setCode?: any;
  defaultType?: TypeProps;
  typeValue?: TypeProps;
  setTypeValue?: any;
}
const CodeEditorWithPreview: FC<CodeEditorWithPreviewProps> = ({
  code,
  setComponentIds,
  setShowTuneDialog,
  defaultType = 'shadow',
  setTypeValue,
  typeValue,
}) => {
  useEffect(() => {
    // eslint-disable-next-line react-hooks/immutability
    updatePreview();
  }, [code]);

  const [type, setType] = usePropsValue({
    defaultValue: defaultType,
    value: typeValue,
    onChange: setTypeValue,
  });

  useEffect(() => {
    return () => {
      if (tooltip) {
        document.body.removeChild(tooltip);
        tooltip = null; // 可选：清空引用
      }
    };
  }, []);

  const updatePreview = () => {
    const previewElement = document.getElementById('preview');
    if (previewElement) {
      if (!previewElement.shadowRoot) {
        const shadowRoot = previewElement.attachShadow({ mode: 'open' });
        const styleElement = document.createElement('style');
        const scriptElement = document.createElement('script');

        styleElement.textContent = `
        * {
          transition: outline 0.1s ease-in-out;
        }
        *[id]:hover, .hover-highlight {
          background-color: rgba(0, 0, 0, 0.5);
        }
      `;

        scriptElement.textContent = `
        document.body.addEventListener('mouseover', (e) => {
          if (e.target.id) {
            e.target.classList.add('hover-highlight');
          }
        });
        document.body.addEventListener('mouseout', (e) => {
          if (e.target.id) {
            e.target.classList.remove('hover-highlight');
          }
        });

      `;
        shadowRoot.appendChild(styleElement);

        shadowRoot.innerHTML = code;
        shadowRoot.appendChild(scriptElement);
      } else {
        previewElement.shadowRoot.innerHTML = code;
        const styleElement = document.createElement('style');
        const scriptElement = document.createElement('script');

        styleElement.textContent = `
       * {
          transition: outline 0.1s ease-in-out;
        }
        *[id]:hover, .hover-highlight {
          background-color: rgba(0, 0, 0, 0.5);
        }
      `;

        scriptElement.textContent = `
       document.body.addEventListener('mouseover', (e) => {
          if (e.target.id) {
            e.target.classList.add('hover-highlight');
          }
        });
        document.body.addEventListener('mouseout', (e) => {
          if (e.target.id) {
            e.target.classList.remove('hover-highlight');
          }
        });
      `;

        previewElement.shadowRoot.prepend(scriptElement);
        previewElement.shadowRoot.prepend(styleElement);
      }
      bindClicks();
      blindHovers();
    }
  };
  const blindHovers = () => {
    const previewElement = document.getElementById('preview');
    if (previewElement && previewElement.shadowRoot) {
      previewElement.shadowRoot.addEventListener('mouseover', handleHover);
      previewElement.shadowRoot.addEventListener('mouseover', updateTooltipPosition);
      previewElement.shadowRoot.addEventListener('mouseout', hideTooltip);
    }
  };
  const handleHover = (event: any) => {
    const path = event.composedPath();
    for (let i = 0; i < path.length; i++) {
      const element = path[i];
      if (element.id) {
        showTooltip(element.id, event);
        break;
      }
    }
  };

  const createTooltip = () => {
    tooltip = document.createElement('div');
    tooltip.id = 'id-tooltip';
    tooltip.style.position = 'absolute';
    tooltip.style.display = 'none';
    tooltip.style.backgroundColor = '#C7F564';
    tooltip.style.color = 'black';
    tooltip.style.fontSize = '12px';
    tooltip.style.padding = '5px';
    tooltip.style.borderRadius = '10px';
    document.body.appendChild(tooltip);
  };

  const showTooltip = (id: any, event: any) => {
    if (!tooltip) createTooltip();
    if (tooltip) {
      // eslint-disable-next-line react-hooks/immutability
      tooltip.textContent = id;
      // eslint-disable-next-line react-hooks/immutability
      tooltip.style.display = 'block';
      updateTooltipPosition(event);
    }
  };

  const hideTooltip = () => {
    if (tooltip) {
      tooltip.style.display = 'none';
    }
  };

  const updateTooltipPosition = (event: any) => {
    if (tooltip && tooltip.style.display === 'block') {
      tooltip.style.left = `${event.pageX + 10}px`;
      tooltip.style.top = `${event.pageY + 10}px`;
    }
  };
  const bindClicks = () => {
    const previewElement = document.getElementById('preview');
    if (previewElement && previewElement.shadowRoot) {
      previewElement.shadowRoot.addEventListener('click', handlePreviewClick);
    }
  };

  const handlePreviewClick = (event: any) => {
    const path = event.composedPath();
    for (let i = 0; i < path.length; i++) {
      const element = path[i];
      if (element.id) {
        console.log('id', element.id);
        setComponentIds?.([element.id]);
        setShowTuneDialog?.(true);
        break;
      }
    }
  };

  const base64SVG = btoa(svg);
  return (
    <div style={{ height: '100%' }}>
      <Radio.Group
        style={{ marginBottom: 2 }}
        defaultValue={defaultType}
        size="small"
        onChange={(e) => {
          setType(e?.target?.value);
        }}
        value={type}
      >
        <Radio.Button
          value="shadow"
          style={{ backgroundColor: 'var(--supos-bg-color)', color: 'var(--supos-text-color)' }}
        >
          shadow
        </Radio.Button>
        <Radio.Button
          value="iframe"
          style={{ backgroundColor: 'var(--supos-bg-color)', color: 'var(--supos-text-color)' }}
        >
          iframe
        </Radio.Button>
      </Radio.Group>
      {type === 'iframe' ? (
        <iframe
          srcDoc={code}
          style={{
            width: '100%',
            height: '100%',
          }}
        />
      ) : null}

      <div
        id="preview"
        className="dot-pattern"
        style={{
          backgroundImage: `url('data:image/svg+xml;base64,${base64SVG}')`,
          backgroundSize: '10px 10px',
          overflow: 'auto',
          display: type === 'shadow' ? 'block' : 'none',
          color: '#161616',
        }}
      ></div>
    </div>
  );
};

export default CodeEditorWithPreview;
