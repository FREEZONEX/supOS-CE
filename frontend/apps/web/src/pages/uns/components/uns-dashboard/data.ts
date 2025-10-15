export const data = {
  nodes: [
    // {
    //   id: 'modbus',
    //   shape: 'modbus',
    //   x: 0,
    //   y: -8,
    //   data: {},
    // },
    // {
    //   id: 'api',
    //   shape: 'api',
    //   x: 0,
    //   y: 48,
    //   data: {},
    // },
    // {
    //   id: 'opcua',
    //   shape: 'opcua',
    //   x: 0,
    //   y: 104,
    //   data: {},
    // },
    // {
    //   id: 'opcda',
    //   shape: 'opcda',
    //   x: 0,
    //   y: 160,
    //   data: {},
    // },
    // {
    //   id: 'relation',
    //   shape: 'relation',
    //   x: 0,
    //   y: 216,
    //   data: {},
    // },
    // {
    //   id: 'icmp',
    //   shape: 'icmp',
    //   x: 0,
    //   y: 272,
    //   data: {},
    // },
    // {
    //   id: 'unknown',
    //   shape: 'unknown',
    //   x: 0,
    //   y: 328,
    //   data: {},
    // },
    {
      id: 'nodeRed',
      shape: 'nodeRed',
      x: 210,
      y: 160,
    },
    {
      id: 'mqtt',
      shape: 'mqtt',
      x: 420,
      y: 160,
    },
    {
      id: 'streamProcessing',
      shape: 'streamProcessing',
      x: 640,
      y: 40,
    },
    {
      id: 'tdEngine',
      shape: 'tdEngine',
      x: 640,
      y: 130,
    },
    {
      id: 'postgreSQL',
      shape: 'postgreSQL',
      x: 640,
      y: 190,
    },
    // {
    //   id: 'gui',
    //   shape: 'gui',
    //   x: 850,
    //   y: 100,
    // },
    // {
    //   id: 'apps',
    //   shape: 'apps',
    //   x: 850,
    //   y: 160,
    // },
    {
      id: 'grafana',
      shape: 'grafana',
      x: 850,
      y: 160,
    },
  ],
  edges: [
    // {
    //   shape: 'edge',
    //   source: 'modbus',
    //   target: 'nodeRed',
    //   vertices: [
    //     { x: 180, y: 14 },
    //     { x: 180, y: 180 },
    //   ],
    //   attrs: {
    //     // line 是选择器名称，选中的边的 path 元素
    //     line: {
    //       sourceMarker: { name: 'circle', r: 2 },
    //       targetMarker: { name: 'circle', r: 2 },
    //       stroke: 'var(--supos-theme-color)',
    //       strokeDasharray: 3,
    //       style: {},
    //     },
    //   },
    // },
    // {
    //   shape: 'edge',
    //   source: 'api',
    //   target: 'nodeRed',
    //   vertices: [
    //     { x: 180, y: 70 },
    //     { x: 180, y: 180 },
    //   ],
    //   attrs: {
    //     // line 是选择器名称，选中的边的 path 元素
    //     line: {
    //       sourceMarker: { name: 'circle', r: 2 },
    //       targetMarker: { name: 'circle', r: 2 },
    //       stroke: 'var(--supos-theme-color)',
    //       strokeDasharray: 3,
    //       style: {},
    //     },
    //   },
    // },
    // {
    //   shape: 'edge',
    //   source: 'opcua',
    //   target: 'nodeRed',
    //   vertices: [
    //     { x: 180, y: 126 },
    //     { x: 180, y: 180 },
    //   ],
    //   attrs: {
    //     // line 是选择器名称，选中的边的 path 元素
    //     line: {
    //       sourceMarker: { name: 'circle', r: 2 },
    //       targetMarker: { name: 'circle', r: 2 },
    //       stroke: 'var(--supos-theme-color)',
    //       strokeDasharray: 3,
    //       style: {},
    //     },
    //   },
    // },
    // {
    //   shape: 'edge',
    //   source: 'opcda',
    //   target: 'nodeRed',
    //   attrs: {
    //     // line 是选择器名称，选中的边的 path 元素
    //     line: {
    //       sourceMarker: { name: 'circle', r: 2 },
    //       targetMarker: { name: 'circle', r: 2 },
    //       stroke: 'var(--supos-theme-color)',
    //       strokeDasharray: 3,
    //       style: {},
    //     },
    //   },
    // },
    // {
    //   shape: 'edge',
    //   source: 'relation',
    //   target: 'nodeRed',
    //   vertices: [
    //     { x: 180, y: 238 },
    //     { x: 180, y: 180 },
    //   ],
    //   attrs: {
    //     // line 是选择器名称，选中的边的 path 元素
    //     line: {
    //       sourceMarker: { name: 'circle', r: 2 },
    //       targetMarker: { name: 'circle', r: 2 },
    //       stroke: 'var(--supos-theme-color)',
    //       strokeDasharray: 3,
    //       style: {},
    //     },
    //   },
    // },
    // {
    //   shape: 'edge',
    //   source: 'icmp',
    //   target: 'nodeRed',
    //   vertices: [
    //     { x: 180, y: 294 },
    //     { x: 180, y: 180 },
    //   ],
    //   attrs: {
    //     // line 是选择器名称，选中的边的 path 元素
    //     line: {
    //       sourceMarker: { name: 'circle', r: 2 },
    //       targetMarker: { name: 'circle', r: 2 },
    //       stroke: 'var(--supos-theme-color)',
    //       strokeDasharray: 3,
    //       style: {},
    //     },
    //   },
    // },
    // {
    //   shape: 'edge',
    //   source: 'unknown',
    //   target: 'nodeRed',
    //   vertices: [
    //     { x: 180, y: 350 },
    //     { x: 180, y: 180 },
    //   ],
    //   attrs: {
    //     // line 是选择器名称，选中的边的 path 元素
    //     line: {
    //       sourceMarker: { name: 'circle', r: 2 },
    //       targetMarker: { name: 'circle', r: 2 },
    //       stroke: 'var(--supos-theme-color)',
    //       strokeDasharray: 3,
    //       style: {},
    //     },
    //   },
    // },
    {
      shape: 'edge',
      source: 'nodeRed',
      target: 'mqtt',
      attrs: {
        // line 是选择器名称，选中的边的 path 元素
        line: {
          sourceMarker: { name: 'circle', r: 2 },
          targetMarker: { name: 'circle', r: 2 },

          stroke: 'var(--supos-theme-color)',
          strokeDasharray: 3,
          style: {
            animation: 'ant-line 60s infinite linear',
          },
        },
      },
    },
    {
      shape: 'edge',
      source: 'mqtt',
      target: 'postgreSQL',
      vertices: [
        { x: 610, y: 180 },
        { x: 610, y: 210 },
      ],
      attrs: {
        // line 是选择器名称，选中的边的 path 元素
        line: {
          sourceMarker: { name: 'circle', r: 2 },
          targetMarker: { name: 'circle', r: 2 },
          stroke: 'var(--supos-theme-color)',
          strokeDasharray: 3,
          style: {
            animation: 'ant-line 60s infinite linear',
          },
        },
      },
    },
    {
      shape: 'edge',
      source: 'mqtt',
      target: 'streamProcessing',
      vertices: [
        { x: 610, y: 180 },
        { x: 610, y: 60 },
      ],
      attrs: {
        // line 是选择器名称，选中的边的 path 元素
        line: {
          sourceMarker: { name: 'circle', r: 2 },
          targetMarker: { name: 'circle', r: 2 },
          stroke: 'var(--supos-theme-color)',
          strokeDasharray: 3,
          style: {
            animation: 'ant-line 60s infinite linear',
          },
        },
      },
    },
    {
      shape: 'edge',
      source: 'mqtt',
      target: 'tdEngine',
      vertices: [
        { x: 610, y: 180 },
        { x: 610, y: 150 },
      ],
      attrs: {
        // line 是选择器名称，选中的边的 path 元素
        line: {
          sourceMarker: { name: 'circle', r: 2 },
          targetMarker: { name: 'circle', r: 2 },
          stroke: 'var(--supos-theme-color)',
          strokeDasharray: 3,
          style: {
            animation: 'ant-line 60s infinite linear',
          },
        },
      },
    },
    {
      shape: 'edge',
      source: 'streamProcessing',
      target: 'tdEngine',
      attrs: {
        // line 是选择器名称，选中的边的 path 元素
        line: {
          sourceMarker: { name: 'circle', r: 2 },
          targetMarker: { name: 'circle', r: 2 },

          stroke: 'var(--supos-theme-color)',
          strokeDasharray: 3,
          style: {
            animation: 'ant-line 60s infinite linear',
          },
        },
      },
    },

    // {
    //   shape: 'edge',
    //   source: 'postgreSQL',
    //   target: 'gui',
    //   vertices: [
    //     { x: 820, y: 210 },
    //     { x: 820, y: 120 },
    //   ],
    //   attrs: {
    //     // line 是选择器名称，选中的边的 path 元素
    //     line: {
    //       sourceMarker: {
    //         name: 'circle',
    //         r: 2,
    //       },
    //       targetMarker: { name: 'circle', r: 2 },
    //
    //       stroke: 'var(--supos-theme-color)',
    //       strokeDasharray: 3,
    //       style: {
    //         animation: 'ant-line 60s infinite linear',
    //       },
    //     },
    //   },
    // },
    // {
    {
      shape: 'edge',
      source: 'tdEngine',
      target: 'grafana',
      vertices: [
        { x: 820, y: 150 },
        { x: 820, y: 180 },
      ],
      attrs: {
        // line 是选择器名称，选中的边的 path 元素
        line: {
          sourceMarker: {
            name: 'circle',
            r: 2,
          },
          targetMarker: { name: 'circle', r: 2 },

          stroke: 'var(--supos-theme-color)',
          strokeDasharray: 3,
          style: {
            animation: 'ant-line 60s infinite linear',
          },
        },
      },
    },
    {
      shape: 'edge',
      source: 'postgreSQL', //来源
      target: 'grafana', //目标
      vertices: [
        { x: 820, y: 210 },
        { x: 820, y: 180 },
      ],
      attrs: {
        // line 是选择器名称，选中的边的 path 元素
        line: {
          sourceMarker: {
            name: 'circle',
            r: 2,
          },
          targetMarker: { name: 'circle', r: 2 },

          stroke: 'var(--supos-theme-color)',
          strokeDasharray: 3,
          style: {
            animation: 'ant-line 60s infinite linear',
          },
        },
      },
    },
  ],
};
