// 支持中文、英文部分特殊字符
import { sqlKeywords } from './sql-keywords.ts';

export const validInputPattern = /^[a-zA-Z0-9\u4e00-\u9fa5_\-/]*$/;

// 排除特殊字符
export const validSpecialCharacter = /^[^~`!@#$%^&*()_+={}[\]\\|;:'",<>./?]*$/;

// 中文和部分特殊字符
export const validNameRegex = /^[a-zA-Z0-9_\-.@&+]*$/;

export const validPicRegex = /^(jpg|jpeg|png|gif|bmp|webp)$/i;

//sql关键字
export const sqlKeywordsRegex = new RegExp(`^(?!.*\\b(${sqlKeywords.join('|')})\\b).*$`, 'i');

// passwordRegex
export const passwordRegex = /^[A-Za-z\d!@#$%^&*()_+\-=$$$${};':"\\|,.<>/?]+$/;

export const phoneRegex = /^\d{11}$/;
