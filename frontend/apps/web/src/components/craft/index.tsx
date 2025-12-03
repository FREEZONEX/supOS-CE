import { type FC, useEffect, useState } from 'react';
import { Button, Flex, Input, message, Table, Typography } from 'antd';
import CodeEditorWithPreview from '../craft/CodeEditorWithPreview';
import { CopilotTask, useCopilotContext } from '@copilotkit/react-core';
import { Subtract } from '@carbon/icons-react';
// import { searchGraphql } from '@/apis/hasura/graphql.ts';
import { getGeneratedPrompt } from '../craft/util.ts';
import ComSelect from '../com-select';
import { useCopilotOperationContext } from '@/layout/context';
import { useTranslate } from '@/hooks';
import { AuthButton } from '../auth';
import { ButtonPermission } from '@/common-types/button-permission';
// import { getBaseUrl } from '@/utils/url-util.ts';
import ProModal from '../pro-modal/index.tsx';
import { setAiResult, useAiStore } from '@/stores/ai-store.ts';
const { Title } = Typography;

export const defaultUI = [
  `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>MES Dashboard</title>
    <style>
        h2 {
            text-align: center;
            font-size: 24px;
            margin-bottom: 20px;
        }
        body {
            font-family: Arial, sans-serif;
            background-color: #F4F4F4; /* 背景色 */
            margin: 0;
            padding: 0;
            color: #2C3E50; /* 文本色 */
        }

        .container {
            background-color: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            width: 300px;
        }

        label {
            font-size: 14px;
        }

        input {
            width: 100%;
            padding: 8px;
            margin: 8px 0 0;
            border: 1px solid #ccc;
            border-radius: 4px;
            box-sizing: border-box;
        }

        button {
            width: 100%;
            padding: 10px;
            background-color: #4CAF50;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
            margin-bottom: 10px;
        }

        .rest {
            background-color: rgb(172, 174, 177);
        }

        .text {
            text-align: center;
        }

        .sign {
            color: #4CAF50;
        }

    </style>
</head>
<body>
    <div class="container">
        <h2>sign In</h2>
        <form action="/submit" method="post">
          <label for="username">Username</label>
          <input type="text" id="username" name="username" required><br><br>
          <label for="password">Password</label>
          <input type="password" id="password" name="password" required><br><br>
          <button type="submit">Submit</button>
          <button class="rest" type="reset" class="">Back to Login</button>
          <div class="text">
            <span>Dont have an account? </span>
            <span class="sign">Sign up</span>
          </div>
        </form>
    </div>
</body>
</html>`,
];

const Board: FC<any> = ({ boardCodeRef }) => {
  const formatMessage = useTranslate();
  const aiResult = useAiStore((state) => state.aiResult);
  const copilotOperation = useCopilotOperationContext();
  // const apiUrl = getBaseUrl();
  const initialCode = () => {
    if (typeof window !== 'undefined') {
      const savedCode = localStorage.getItem('code');
      return savedCode ? JSON.parse(savedCode) : defaultUI;
    }
    return defaultUI;
  };
  const [code, setCode] = useState(initialCode);
  const [codeToDisplay, setCodeToDisplay] = useState(code[code.length - 1] || '');
  const [codeCommand, setCodeCommand] = useState({
    prompt: '',
    databaseName: null,
  });
  const [tuneCommand, setTuneCommand] = useState('');
  const [componentIds, setComponentIds] = useState([]);
  const [selectedIndex, setSelectedIndex] = useState(code?.length - 1);
  const [showTuneDialog, setShowTuneDialog] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const setCodeToDisplayFn = (code: any) => {
    if (boardCodeRef) {
      boardCodeRef.current = code;
    }
    setCodeToDisplay(code);
  };
  useEffect(() => {
    const savedCode = localStorage.getItem('code');
    if (savedCode) {
      setCode(JSON.parse(savedCode));
    }
    const codeCommandStr = localStorage.getItem('codeCommand');
    if (codeCommandStr) {
      const codeCommandItems = JSON.parse(codeCommandStr)[code.length - 1];
      setCodeCommand({ ...codeCommand, prompt: codeCommandItems?.prompt || '' });
    }
  }, []);
  useEffect(() => {
    localStorage.setItem('code', JSON.stringify(code));
    if (code.length > 0) {
      setSelectedIndex(code.length - 1);
      setCodeToDisplayFn(code[code.length - 1]);
    }
  }, [code]);

  useEffect(() => {
    if (aiResult?.['app-gui']) {
      if (isLoading) {
        message.error(formatMessage('appGui.aiError'));
        return;
      }
      console.log(copilotOperation);
      copilotOperation?.current?.setOpen?.(false);
      const prompt = aiResult?.['app-gui'] || '';
      setAiResult('app-gui', null);
      setIsLoading(true);
      setCodeCommand({ ...codeCommand, prompt: prompt });
      getBaseDataFields();
    }
  }, [aiResult?.['app-gui']]);

  const tuneComponents = new CopilotTask({
    instructions: `当前页面的代码为${codeToDisplay}.` + tuneCommand,
    actions: [
      {
        name: 'generateCode',
        description: `修改且仅修改id为${componentIds.join(
          ','
        )}的元素的样式。返回修改后的html代码, 不修改的区域保持不变。注意，返回值是修改后完整的html页面代码！`,
        parameters: [
          {
            name: 'html_code',
            type: 'string',
            description: 'Code to be generated',
            required: true,
          },
        ],
        handler: async ({ html_code }) => {
          try {
            setCode((prev: any) => [...prev, html_code]);
            setCodeToDisplayFn(html_code);
          } finally {
            setIsLoading(false);
          }
        },
      },
    ],
  });
  const context = useCopilotContext();
  const actions: any = [
    {
      name: 'generateCode',
      description:
        // 'generate a single page react app, only react code is allowed, and do not include the import and export lines',
        '生成一个完整的HTML页面代码，只能生成html! 做到页面美观， 生成的每一个element都有一个唯一的id',
      parameters: [
        {
          name: 'html_code',
          type: 'string',
          description: 'Code to be generated',
          required: true,
        },
      ],
      handler: async ({ html_code }: any) => {
        try {
          setCode((prev: any) => {
            const oldCodeCommand: any = localStorage.getItem('codeCommand');
            if (oldCodeCommand) {
              const oldList: any[] = JSON.parse(oldCodeCommand) || [];
              localStorage.setItem('codeCommand', JSON.stringify([...oldList, codeCommand]));
            } else {
              localStorage.setItem(
                'codeCommand',
                JSON.stringify([
                  ...(prev?.map(() => ({
                    prompt: '',
                    databaseName: null,
                  })) || []),
                  codeCommand,
                ])
              );
            }
            return [...prev, html_code];
          });
          setCodeToDisplayFn(html_code);
        } finally {
          setIsLoading(false);
        }
      },
    },
  ];
  const getBaseDataFields = () => {
    if (!codeCommand.databaseName) {
      const generateCode = new CopilotTask({
        instructions: getGeneratedPrompt({ codeCommand }),
        actions: actions,
      });
      generateCode.run(context);
      return;
    }
    // searchGraphql({
    //   query: String.raw`{
    //           __type(name: "${codeCommand.databaseName}") {
    //             fields {
    //               name
    //               type {
    //                 name
    //                 ofType {
    //                   name
    //                 }
    //               }
    //             }
    //           }
    //         }`,
    // })
    //   .then((data) => {
    //     const fileds = data?.data?.__type?.fields?.map?.((field: any) => ({
    //       name: field.name,
    //       type: field.type?.name || field.type?.ofType?.name,
    //     }));
    //     const fieldsName = fileds?.map((item: any) => item.name).join(' ');
    //     const generateCode = new CopilotTask({
    //       instructions: getGeneratedPrompt({ apiUrl, codeCommand, fieldsName, fileds }),
    //       actions: actions,
    //     });
    //     generateCode.run(context);
    //   })
    //   .catch(() => {
    //     const generateCode = new CopilotTask({
    //       instructions: getGeneratedPrompt({ codeCommand }),
    //       actions: actions,
    //       includeCopilotActions: false,
    //     });
    //     generateCode.run(context);
    //   });
  };
  // const [databaseList, setDatabaseList] = useState<any>([]);
  useEffect(() => {
    //     searchGraphql({
    //       query: String.raw`query {
    //   __type(name: "query_root") {
    //     fields {
    //       name
    //     }
    //   }
    // }`,
    //     }).then((data: any) => {
    //       setDatabaseList(
    //         data?.data?.__type?.fields
    //           ?.filter((item: any) => !(item.name.includes('_by_pk') || item.name.includes('_aggregate')))
    //           ?.map((item: any) => ({
    //             label: item.name,
    //             value: item.name,
    //           }))
    //       );
    //     });
  }, []);
  return (
    <Flex style={{ height: '100%', width: '100%', padding: '40px 20px' }} gap="20px">
      <Flex
        style={{
          width: 350,
          height: '100%',
        }}
        vertical
      >
        <Title level={5} style={{ marginBottom: 0 }}>
          {formatMessage('appGui.databaseName')}
        </Title>
        <div>
          <ComSelect
            value={codeCommand.databaseName}
            onChange={(value) => {
              setCodeCommand({
                ...codeCommand,
                databaseName: value,
              });
            }}
            showSearch
            options={[]}
            style={{ margin: '8px 0', width: '100%' }}
            placeholder={formatMessage('appGui.yourCodeDatabaseName')}
            allowClear
          />
        </div>
        <Title level={5} style={{ marginBottom: 0 }}>
          {formatMessage('uns.description')}
        </Title>
        <div>
          <Input.TextArea
            rows={12}
            style={{ margin: '8px 0' }}
            placeholder={formatMessage('appGui.enterCodeCommand')}
            value={codeCommand.prompt}
            onChange={(e) => {
              setCodeCommand({
                ...codeCommand,
                prompt: e.target.value,
              });
            }}
          />
        </div>

        <AuthButton
          auth={ButtonPermission['appGui.generate']}
          disabled={isLoading}
          onClick={() => {
            setIsLoading(true);
            getBaseDataFields();
          }}
          style={{ marginTop: 8, borderColor: 'var(--supos-text-color)' }}
          block
          color="default"
          variant="solid"
        >
          {formatMessage('appSpace.newgenerate')}
        </AuthButton>
        <Title level={5} style={{ marginBottom: 0, marginTop: 5 }}>
          {formatMessage('common.history')}
        </Title>
        <div
          style={{
            flex: 1,
            marginTop: 8,
            overflow: 'auto',
            padding: '10px 8px 8px',
            backgroundColor: 'var(--supos-craft-bg-color)',
          }}
        >
          <Table
            className="carft-table"
            size="small"
            columns={[
              {
                dataIndex: 'version',
                render: (_, __, index) => (
                  <Flex justify="space-between" align="center">
                    <span>v1.0.0-beta.{index}</span>
                    {code.length > 1 && (
                      <Subtract
                        onClick={(e) => {
                          e.stopPropagation();
                          if (code.length === 1) return;
                          setCode((pre: any) => {
                            const newCode = pre?.filter((_p: any, i: number) => i !== index) || [];
                            const oldCodeCommand: any = localStorage.getItem('codeCommand');
                            if (oldCodeCommand) {
                              const oldList: any[] = JSON.parse(oldCodeCommand) || [];
                              localStorage.setItem(
                                'codeCommand',
                                JSON.stringify(oldList?.filter((_p: any, i: number) => i !== index))
                              );
                            }
                            return newCode;
                          });
                        }}
                        style={{ cursor: 'pointer' }}
                      />
                    )}
                  </Flex>
                ),
              },
            ]}
            dataSource={code.map((c: any, index: number) => ({
              key: index,
              index,
              code: c,
            }))}
            rowClassName={(record) => (record.index === selectedIndex ? 'selected-row' : '')}
            onRow={(record: any) => ({
              onClick: () => {
                const codeCommandStor = localStorage.getItem('codeCommand');
                if (codeCommandStor) {
                  const codeCommandItems = JSON.parse(codeCommandStor)[record.index];
                  setCodeCommand({ ...codeCommand, prompt: codeCommandItems?.prompt || '' });
                }
                setSelectedIndex(record.index);
                setCodeToDisplayFn(record.code);
              },
            })}
            loading={isLoading}
            pagination={false}
            rowKey="index"
          />
        </div>
      </Flex>
      <div
        style={{
          flex: 1,
          height: '100%',
        }}
      >
        <CodeEditorWithPreview
          code={codeToDisplay}
          setCode={setCodeToDisplay}
          setComponentIds={setComponentIds}
          setShowTuneDialog={setShowTuneDialog}
        />
      </div>
      <ProModal size="sm" maskClosable={false} open={showTuneDialog} onCancel={() => setShowTuneDialog(false)}>
        <Title style={{ fontWeight: 400, marginBottom: 0 }} type="secondary" level={4}>
          Changing element: {componentIds}
        </Title>
        <Title style={{ margin: 0 }} type="secondary" level={5}>
          Now you are changing element username only.
        </Title>
        <Input.TextArea
          rows={6}
          value={tuneCommand}
          onChange={(e) => setTuneCommand(e.target.value)}
          // id="code"
          placeholder="Type your description here"
          style={{ backgroundColor: 'white', marginTop: 20 }}
        />
        <Button
          onClick={() => {
            setIsLoading(true);
            tuneComponents.run(context);
            setShowTuneDialog(false);
          }}
          style={{ marginTop: 8 }}
          block
          color="default"
          variant="solid"
        >
          Start Tuning!
        </Button>
      </ProModal>
    </Flex>
  );
};

export default Board;
