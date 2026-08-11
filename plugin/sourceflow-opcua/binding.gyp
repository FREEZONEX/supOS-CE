{
  "targets": [
    {
      "target_name": "opcua_open62541",
      "sources": [
        "src/addon.cc"
      ],
      "include_dirs": [
        "<!@(node -p \"require('node-addon-api').include\")",
        "/usr/local/include",
        "<!(node -e \"const fs=require('fs'),path=require('path'); const candidates=[path.join(process.cwd(),'../opcua/.build/open62541-install/include'),'/usr/local/include']; for(const c of candidates){ try{ fs.accessSync(path.join(c,'open62541','types.h')); console.log(c); process.exit(0); }catch(e){} } console.log('/usr/local/include');\")"
      ],
      "libraries": [
        "-L/usr/local/lib",
        "-L<!(node -e \"const fs=require('fs'),path=require('path'); const candidates=[path.join(process.cwd(),'../opcua/.build/open62541-install/lib'),'/usr/local/lib']; for(const c of candidates){ try{ fs.accessSync(path.join(c,'libopen62541.a')); console.log(c); process.exit(0); }catch(e){} } console.log('/usr/local/lib');\")",
        "-lopen62541"
      ],
      "dependencies": [
        "<!(node -p \"require('node-addon-api').gyp\")"
      ],
      "defines": [
        "NAPI_CPP_EXCEPTIONS"
      ],
      "cflags!": [
        "-fno-exceptions"
      ],
      "cflags_cc!": [
        "-fno-exceptions"
      ],
      "cflags_cc": [
        "-std=c++20"
      ],
      "conditions": [
        [
          "OS=='linux'",
          {
            "libraries": [
              "-lopen62541",
              "-lpthread",
              "-lm"
            ]
          }
        ]
      ]
    }
  ]
}
