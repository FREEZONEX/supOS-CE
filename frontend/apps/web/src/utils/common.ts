export const isJsonString = (value: any): boolean => {
  if (!value) return false;
  try {
    const jsonVal = JSON.parse(value);
    if (['[object Object]', '[object Array]'].includes(Object.prototype.toString.call(jsonVal))) {
      return true;
    } else {
      return false;
    }
    // eslint-disable-next-line
  } catch (err) {
    return false;
  }
};

type CopyCallback = (success: boolean) => void;

const fallbackCopyToClipboard = (text: string): boolean => {
  const textarea = document.createElement('textarea');
  const activeElement = document.activeElement as HTMLElement | null;
  const selection = document.getSelection();
  const originalRange = selection && selection.rangeCount > 0 ? selection.getRangeAt(0) : null;

  textarea.value = text;
  textarea.setAttribute('readonly', 'readonly');
  textarea.style.position = 'fixed';
  textarea.style.top = '0';
  textarea.style.left = '-9999px';
  textarea.style.opacity = '0';
  document.body.appendChild(textarea);

  try {
    textarea.focus();
    textarea.select();
    textarea.setSelectionRange(0, text.length);
    return document.execCommand('copy');
  } catch (err) {
    console.log(err);
    return false;
  } finally {
    document.body.removeChild(textarea);
    if (selection) {
      selection.removeAllRanges();
      if (originalRange) {
        selection.addRange(originalRange);
      }
    }
    activeElement?.focus?.();
  }
};

// 使用 const 定义一个函数，并同时指定参数和返回值的类型
export const copyToClipboard: (text: string, callback?: CopyCallback) => void = (
  text: string,
  callback?: CopyCallback
): void => {
  const handleResult = (success: boolean) => {
    if (callback) {
      callback(success);
    }
  };

  if (navigator.clipboard && typeof navigator.clipboard.writeText === 'function' && window.isSecureContext) {
    // 使用 Clipboard API
    const writeResult = navigator.clipboard.writeText(text);
    if (writeResult && typeof writeResult.then === 'function') {
      writeResult
        .then(() => {
          handleResult(true); // 成功时处理结果
        })
        .catch(() => {
          handleResult(fallbackCopyToClipboard(text)); // 失败时回退到旧方式
        });
    } else {
      handleResult(fallbackCopyToClipboard(text));
    }
  } else {
    handleResult(fallbackCopyToClipboard(text));
  }
};

export function canModifyParentHref() {
  try {
    if (window.parent === window || !window.parent) {
      return false;
    }

    const parentOrigin = window.parent.location.origin;
    const currentOrigin = window.location.origin;

    if (parentOrigin === currentOrigin) {
      // 同域
      return true;
    }

    return false;
  } catch (e) {
    console.log(e);
    // 跨域
    return -1;
  }
}
