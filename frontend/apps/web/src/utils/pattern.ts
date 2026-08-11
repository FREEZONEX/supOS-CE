// 支持中文、英文部分特殊字符
import { sqlKeywords } from './sql-keywords.ts';

export const validInputPattern = /^[a-zA-Z0-9\u4e00-\u9fa5_\-/]*$/;

// 排除特殊字符
export const validSpecialCharacter = /^[^~`!@#$%^&*()_+={}[\]\\|;:'",<>./?]*$/;

// 中文、其他语言字母和部分特殊字符
export const validNameRegex = /^[\p{L}\p{N}_\-.@&+]*$/u;

export const validPicRegex = /^(jpg|jpeg|png|gif|bmp|webp)$/i;

//sql关键字
export const sqlKeywordsRegex = new RegExp(`^(?!.*\\b(${sqlKeywords.join('|')})\\b).*$`, 'i');

// passwordRegex
export const passwordRegex = /^[A-Za-z\d!@#$%^&*()_+\-=$$$${};':"\\|,.<>/?]+$/;
export const passwordStrengthRegex = /^(?=.*[A-Za-z])(?=.*\d).{8,20}$/;

export const phoneRegex = /^\d{11}$/;
