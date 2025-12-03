import type { CopilotRefProps } from '@/components/copilotkit/CustomCopilotChat';
import { createContext, type RefObject, useContext } from 'react';

export const CopilotOperationContext = createContext<RefObject<CopilotRefProps | undefined> | undefined>(undefined);

export const useCopilotOperationContext = () => useContext(CopilotOperationContext);
