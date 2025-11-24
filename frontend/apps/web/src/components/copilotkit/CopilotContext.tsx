import {
  useCopilotAction,
  useCopilotReadable,
  useCoAgent,
  type CatchAllActionRenderProps,
} from '@copilotkit/react-core';
import { type FC, type ReactNode, useEffect, useState } from 'react';
import { useLocalStorage } from '@/hooks';
import { useCopilotChatSuggestions } from '@copilotkit/react-ui';
import { useLocation } from 'react-router';
import MermaidCom from './sub-com/MermaidCom.tsx';
import { attempt, isError } from 'lodash-es';
import { useBaseStore } from '@/stores/base';
import MCPToolCall from './sub-com/McpToolCall.tsx';

const JSONParse = (str: string | null) => {
  let json = attempt(JSON.parse, str);

  if (isError(json)) {
    json = {};
  }
  return json;
};

export enum Page {
  Uns = 'uns',
}

export enum UnsPageOperations {
  UnsOperation = 'unsOperation',
}

export const AVAILABLE_OPERATIONS_PER_PAGE = {
  [Page.Uns]: Object.values(UnsPageOperations),
};

// Local storage key for saving agent state
const STORAGE_KEY = 'mcp-agent-state';
const STORAGE_MODEL_KEY = 'mcp-agent-model-state';

const SuggestionsContext = ({ children }: { children: ReactNode }) => {
  // 提示语
  useCopilotChatSuggestions({
    instructions: `根据useMermaid等action生成提示语，方便用户操作`,
    minSuggestions: 1,
    maxSuggestions: 1,
  });

  return children;
};

// 通用的readable和action放这里
const CopilotContext: FC<{ children: ReactNode; copilotCatRef: any }> = ({ children }) => {
  const { menuGroupNoSub, systemInfo } = useBaseStore((state) => ({
    systemInfo: state.systemInfo,
    menuGroupNoSub: state.menuGroup?.filter((f) => !f.subMenu),
  }));
  const routeMap = menuGroupNoSub?.map((item) => item.showName ?? item.code);
  const [currentPage, setPage] = useState<string>('');
  const location = useLocation();
  useEffect(() => {
    setPage(location.pathname);
  }, [location.pathname]);
  // 获取mcpclient存储的本地配置
  const mcpConfigStr: any = useLocalStorage(STORAGE_KEY);
  const modelConfigStr: any = useLocalStorage(STORAGE_MODEL_KEY);

  const modelConfig = JSONParse(modelConfigStr);

  // Initialize agent state with the data from localStorage
  const { state: agentState, setState: setAgentState } = useCoAgent<any>({
    name: 'sample_agent',
    initialState: {
      modelSdk: modelConfig?.modelSdk || 'openai',
      model: modelConfig?.model,
      apiKey: modelConfig?.apiKey,
      mcp_config: JSONParse(mcpConfigStr),
    },
  });

  useEffect(() => {
    setAgentState({ ...agentState, mcp_config: JSONParse(mcpConfigStr) });
  }, [mcpConfigStr]);

  useEffect(() => {
    const newModelConfig = JSONParse(modelConfigStr);
    setAgentState({
      ...agentState,
      modelSdk: newModelConfig?.modelSdk || 'openai',
      model: newModelConfig?.model,
      apiKey: newModelConfig?.apiKey,
    });
  }, [modelConfigStr]);

  // 数据可读
  useCopilotReadable({
    description: 'pages是页面集合;mcpServiceList是目前已有的mcp服务器列表，如果没有不要调用;',
    value: {
      pages: routeMap,
      operations: AVAILABLE_OPERATIONS_PER_PAGE,
      currentPage,
      mcpServiceList: [],
    },
  });

  useCopilotAction({
    name: 'useMermaid',
    description: `
    # Mermaid图表绘制职责与脚本说明
    ## 职责说明
    你是一个图表绘制专家，分析用户的问题和数据等，使用Mermaid语法输出专业的图表。
    支持的图表类型：流程图、时序图、类图、状态图、柱状图、实体关系图、用户旅程图、甘特图、饼图、象限图、需求图、Gitgraph 图、C4 图、思维导图、时间线图、ZenUML、桑基图、XY 图、框图、数据包图、看板图、架构图
    **举例**
    XY图示例:
      xychart-beta
      title "Sales Revenue"
      x-axis [jan, feb, mar, apr, may, jun, jul]
      y-axis "Revenue (in $)" 4000 --> 11000
      bar [5000, 6000, 7500, 8200, 9500, 10500, 11000]
      line [5000, 6000, 7500, 8200, 9500, 10500, 11000]

    ###注意
    解析用户需求，选择合适的图表输出，只输出图表即可，不要有额外的回复`,
    parameters: [
      {
        name: 'code',
        type: 'string',
        description: '输出的Mermaid语法的code',
        required: true,
        enum: routeMap,
      },
    ],
    followUp: false,
    render: (props: any) => {
      const { code } = props.args;
      const { status } = props;
      // if (status === 'inProgress') {
      //   return <InlineLoading status="active" description="Loading data..." />;
      // }
      if (status === 'complete') {
        return <MermaidCom code={`${code}`} />;
      }
      return <></>;
    },
  });

  // add a custom action renderer for all actions
  // useCopilotAction({
  //   name: '*',
  //   render: ({ name, args, status, result }: any) => {
  //     return <ToolCallRenderer name={name} args={args} status={status || 'unknown'} result={result} />;
  //   },
  // });

  useCopilotAction({
    /**
     * The asterisk (*) matches all tool calls
     */
    name: '*',
    render: ({ name, status, args, result }: CatchAllActionRenderProps<[]>) => (
      <MCPToolCall status={status} name={name} args={args} result={result} />
    ),
  });

  return systemInfo?.llmType === 'openai' ? <SuggestionsContext>{children}</SuggestionsContext> : children;
};

export default CopilotContext;
