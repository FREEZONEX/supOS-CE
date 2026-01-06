import { type RefObject, useEffect, useState, useRef } from 'react';
import Clipboard from 'clipboard';
import { useTranslate } from '@/hooks';
import { App } from 'antd';

// 复制的hooks
const useClipboard = (buttonRef: RefObject<HTMLElement> | null = null, textToCopy?: string, msg?: string) => {
  const [isCopied, setIsCopied] = useState(false); // 状态：是否成功复制
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  // 如果没有传入ref，则创建一个临时的DOM元素
  const internalRef = useRef<HTMLDivElement | null>(null);
  // 用于存储clipboard实例的ref
  const clipboardRef = useRef<Clipboard | null>(null);

  // 确定使用哪个ref
  const targetRef = buttonRef || internalRef;

  // 创建临时DOM元素的辅助函数
  const createTempElement = () => {
    if (!buttonRef && !internalRef.current) {
      const tempElement = document.createElement('div');
      tempElement.style.position = 'absolute';
      tempElement.style.top = '-9999px';
      tempElement.style.left = '-9999px';
      document.body.appendChild(tempElement);
      internalRef.current = tempElement;
    }
  };

  // 清理临时DOM元素的辅助函数
  const cleanupTempElement = () => {
    if (!buttonRef && internalRef.current && document.body.contains(internalRef.current)) {
      document.body.removeChild(internalRef.current);
      internalRef.current = null;
    }
  };

  useEffect(() => {
    // 使用辅助函数创建临时DOM元素
    createTempElement();

    if (!targetRef.current) {
      return () => {
        cleanupTempElement(); // 清理临时DOM元素
      };
    }

    if (!textToCopy) {
      return () => {
        cleanupTempElement(); // 清理临时DOM元素
      };
    }
    const clipboard = new Clipboard(targetRef.current, {
      text: () => textToCopy, // 复制的内容
    });

    clipboardRef.current = clipboard;

    clipboard.on('success', () => {
      setIsCopied(true); // 更新状态为已复制
      message.success(msg ?? `${formatMessage('common.copySuccess')}`).then(() => {
        setIsCopied(false);
      });
    });

    clipboard.on('error', () => {
      setIsCopied(false); // 复制失败时状态重置
      console.error('Failed to copy text');
    });

    return () => {
      clipboard.destroy(); // 清除事件监听和实例
      cleanupTempElement(); // 清理临时DOM元素
    };
  }, [buttonRef, targetRef, textToCopy]);

  // 手动触发复制的方法
  const copy = (text?: string) => {
    // 如果传入了新的文本，则使用新文本；否则使用初始设置的文本
    const textToUse = text !== undefined ? text : textToCopy;

    if (!textToUse) return;
    // 确保临时DOM元素存在
    createTempElement();

    if (!targetRef.current) return false;

    // 如果clipboard实例不存在，则创建一个新的
    if (!clipboardRef.current) {
      clipboardRef.current = new Clipboard(targetRef.current, {
        text: () => textToUse,
      });

      // 设置一次性事件监听
      clipboardRef.current.on('success', () => {
        setIsCopied(true);
        message.success(msg ?? `${formatMessage('common.copySuccess')}`).then(() => {
          setIsCopied(false);
        });
        if (clipboardRef.current) {
          clipboardRef.current.destroy();
          clipboardRef.current = null;
        }
        // 确保在复制成功后也清理临时DOM元素
        cleanupTempElement();
      });

      clipboardRef.current.on('error', () => {
        setIsCopied(false);
        console.error('Failed to copy text');
        if (clipboardRef.current) {
          clipboardRef.current.destroy();
          clipboardRef.current = null;
        }
        // 确保在复制失败时也清理临时DOM元素
        cleanupTempElement();
      });
    }

    try {
      targetRef.current.click();
      return true;
    } catch (error) {
      console.error('Failed to trigger copy:', error);
      // 发生错误时清理资源
      if (clipboardRef.current) {
        clipboardRef.current.destroy();
        clipboardRef.current = null;
      }
      // 确保在错误情况下也清理临时DOM元素
      cleanupTempElement();
      return false;
    }
  };

  return { isCopied, copy };
};

export default useClipboard;
