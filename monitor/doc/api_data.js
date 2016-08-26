define({ "api": [
  {
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "optional": false,
            "field": "varname1",
            "description": "<p>No type.</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "varname2",
            "description": "<p>With type.</p>"
          }
        ]
      }
    },
    "type": "",
    "url": "",
    "version": "0.0.0",
    "filename": "./doc/main.js",
    "group": "D__git_client_RDAWatchServer_src_github_com_oikomi_FishChatServer_monitor_doc_main_js",
    "groupTitle": "D__git_client_RDAWatchServer_src_github_com_oikomi_FishChatServer_monitor_doc_main_js",
    "name": ""
  },
  {
    "type": "get",
    "url": "/device/:IMEI/chatRecord",
    "title": "拉取聊天记录",
    "name": "DeivceChatRecord",
    "group": "Device",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "username",
            "description": "<p>用户名</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "ticket",
            "description": "<p>用户接口调用凭据</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "allowedValues": [
              "\"timeDesc\""
            ],
            "optional": true,
            "field": "orderType",
            "defaultValue": "timeDesc",
            "description": "<p>顺序类型，默认为timeDesc</p>"
          },
          {
            "group": "Parameter",
            "type": "Number",
            "optional": true,
            "field": "startId",
            "defaultValue": "0",
            "description": "<p>开始拉取id</p>"
          },
          {
            "group": "Parameter",
            "type": "Number",
            "optional": true,
            "field": "length",
            "defaultValue": "10",
            "description": "<p>拉取长度</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "Number",
            "optional": false,
            "field": "id",
            "description": "<p>消息id</p>"
          },
          {
            "group": "Success 200",
            "type": "Number",
            "optional": false,
            "field": "direction",
            "description": "<p>语音方向，1为上行（设备-&gt;服务器），2为下行（服务器-&gt;设备）</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "type",
            "description": "<p>消息类型，目前为voice</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": true,
            "field": "voiceUrl",
            "description": "<p>语音url，消息类型为voice时返回</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "created_at",
            "description": "<p>消息产生时间，格式为Y-m-d H:i:s</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "status",
            "description": "<p>当前消息状态</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Success-Response:",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 0,\n    \"errmsg\": \"操作成功\",\n    \"data\": [{\n        \"id\": 2,\n        \"direction\": 2,\n        \"type\": \"voice\",\n        \"voiceUrl\": \"http://xxx.xxx.com/\",\n        \"created_at\": \"2016-01-01 00:00:00\",\n        \"status\": \"已发送到手表\"\n    },{\n        \"id\": 1,\n        \"direction\": 1,\n        \"type\": \"voice\",\n        \"voiceUrl\": \"http://xxx.xxx.com/\",\n        \"created_at\": \"2016-01-01 00:00:00\",\n        \"status\": \"已读\"\n    }]\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./controllers/device.go",
    "groupTitle": "Device",
    "sampleRequest": [
      {
        "url": "http://api.watch.h-hy.com:8080/v1/device/:IMEI/chatRecord"
      }
    ]
  },
  {
    "type": "post",
    "url": "/device/:IMEI/voice",
    "title": "发送聊天语音",
    "name": "DeivceSendVoice",
    "group": "Device",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "username",
            "description": "<p>用户名</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "ticket",
            "description": "<p>用户接口调用凭据</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "wechatMediaId",
            "description": "<p>微信提供的mediaId</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "Number",
            "optional": false,
            "field": "id",
            "description": "<p>消息id</p>"
          },
          {
            "group": "Success 200",
            "type": "Number",
            "optional": false,
            "field": "direction",
            "description": "<p>语音方向，1为上行（设备-&gt;服务器），2为下行（服务器-&gt;设备）</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "type",
            "description": "<p>消息类型，目前为voice</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "voiceUrl",
            "description": "<p>语音url</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "created_at",
            "description": "<p>消息产生时间，格式为Y-m-d H:i:s</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "status",
            "description": "<p>当前消息状态</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Success-Response:",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 0,\n    \"errmsg\": \"操作成功\",\n    \"data\": {\n        \"id\": 2,\n        \"direction\": 2,\n        \"type\": \"voice\",\n        \"voiceUrl\": \"http://xxx.xxx.com/\",\n        \"created_at\": \"2016-01-01 00:00:00\",\n        \"status\": \"发送中...\"\n    }\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./controllers/device.go",
    "groupTitle": "Device",
    "sampleRequest": [
      {
        "url": "http://api.watch.h-hy.com:8080/v1/device/:IMEI/voice"
      }
    ]
  },
  {
    "type": "post",
    "url": "/device/:IMEI/action/shutdown",
    "title": "设备关机",
    "name": "DeivceShutdown",
    "group": "Device",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "username",
            "description": "<p>用户名</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "ticket",
            "description": "<p>用户接口调用凭据</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "Number",
            "optional": false,
            "field": "messageId",
            "description": "<p>消息ID</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Success-Response:",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 0,\n    \"errmsg\": \"操作成功\",\n    \"data\": {\n        \"messageId\": 123,\n    }\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./controllers/device.go",
    "groupTitle": "Device",
    "sampleRequest": [
      {
        "url": "http://api.watch.h-hy.com:8080/v1/device/:IMEI/action/shutdown"
      }
    ]
  },
  {
    "type": "put",
    "url": "/device/:IMEI",
    "title": "更新设备信息",
    "description": "<p>本接口只需要传入需要更新的参数即可，无需更新的无需传入</p>",
    "name": "DeivceUpdate",
    "group": "Device",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "username",
            "description": "<p>用户名</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "ticket",
            "description": "<p>用户接口调用凭据</p>"
          },
          {
            "group": "Parameter",
            "type": "Number",
            "optional": true,
            "field": "work_model",
            "description": "<p>工作模式</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": true,
            "field": "emeregncyPhone",
            "description": "<p>设备紧急号码</p>"
          },
          {
            "group": "Parameter",
            "type": "Number",
            "optional": true,
            "field": "volume",
            "description": "<p>设备音量</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": true,
            "field": "nick",
            "description": "<p>设备昵称</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Request-Example:",
          "content": "传入需要更新的参数即可，无需更新的无需传入\nwork_model=1234567891011&nick=abc&volume=6&emeregncyPhone=13590210000",
          "type": "String"
        }
      ]
    },
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "Number",
            "optional": false,
            "field": "messageId",
            "description": "<p>消息ID</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Success-Response:",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 0,\n    \"errmsg\": \"操作成功\",\n    \"data\": {\n        \"messageId\": 123,\n    }\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./controllers/device.go",
    "groupTitle": "Device",
    "sampleRequest": [
      {
        "url": "http://api.watch.h-hy.com:8080/v1/device/:IMEI"
      }
    ]
  },
  {
    "type": "post",
    "url": "/device/:IMEI/action/location",
    "title": "设备实时定位",
    "name": "DeivceUpdateLocation",
    "group": "Device",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "username",
            "description": "<p>用户名</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "ticket",
            "description": "<p>用户接口调用凭据</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "Number",
            "optional": false,
            "field": "messageId",
            "description": "<p>消息ID</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Success-Response:",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 0,\n    \"errmsg\": \"操作成功\",\n    \"data\": {\n        \"messageId\": 123,\n    }\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./controllers/device.go",
    "groupTitle": "Device",
    "sampleRequest": [
      {
        "url": "http://api.watch.h-hy.com:8080/v1/device/:IMEI/action/location"
      }
    ]
  },
  {
    "type": "post",
    "url": "/device",
    "title": "用户绑定设备",
    "name": "deviceBinding",
    "group": "Device",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "username",
            "description": "<p>用户名</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "ticket",
            "description": "<p>用户接口调用凭据</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "IMEI",
            "description": "<p>设备IMEI</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "nick",
            "description": "<p>设备昵称</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Request-Example:",
          "content": "IMEI=1234567891011&nick=abc",
          "type": "String"
        }
      ]
    },
    "success": {
      "examples": [
        {
          "title": "Success-Response:",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 0,\n    \"errmsg\": \"操作成功\",\n    \"data\": {\n        \"IMEI\": \"123456789101112\",\n        \"nick\": \"123\",\n        \"status\": 1\n        \"work_model\": 1,\n        \"volume\": 6,\n        \"electricity\": 100,\n        \"emeregncyPhone\": \"13590210000\",\n    }\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./controllers/device.go",
    "groupTitle": "Device",
    "sampleRequest": [
      {
        "url": "http://api.watch.h-hy.com:8080/v1/device"
      }
    ]
  },
  {
    "type": "delete",
    "url": "/device/:IMEI?username=:user&ticket=:ticket",
    "title": "用户删除绑定设备",
    "description": "<p>特别说明：根据HTTP标准，DELETE方法的身份认证参数务必放在url中而不能放在body中。</p>",
    "name": "deviceDestory",
    "group": "Device",
    "success": {
      "examples": [
        {
          "title": "Success-Response:",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 0,\n    \"errmsg\": \"操作成功\",\n    \"data\": {\n    }\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./controllers/device.go",
    "groupTitle": "Device",
    "sampleRequest": [
      {
        "url": "http://api.watch.h-hy.com:8080/v1/device/:IMEI?username=:user&ticket=:ticket"
      }
    ]
  },
  {
    "type": "get",
    "url": "/device/:IMEI",
    "title": "查看用户设备详情",
    "name": "deviceDetail",
    "group": "Device",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "username",
            "description": "<p>用户名</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "ticket",
            "description": "<p>用户接口调用凭据</p>"
          }
        ]
      }
    },
    "success": {
      "examples": [
        {
          "title": "Success-Response:",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 0,\n    \"errmsg\": \"操作成功\",\n    \"data\": {\n        \"IMEI\": \"123456789101112\",\n        \"nick\": \"123\",\n        \"status\": 1\n        \"work_model\": 1,\n        \"volume\": 6,\n        \"electricity\": 100,\n        \"emeregncyPhone\": \"13590210000\",\n    }\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./controllers/device.go",
    "groupTitle": "Device",
    "sampleRequest": [
      {
        "url": "http://api.watch.h-hy.com:8080/v1/device/:IMEI"
      }
    ]
  },
  {
    "type": "get",
    "url": "/device",
    "title": "查看用户设备列表",
    "name": "deviceList",
    "group": "Device",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "username",
            "description": "<p>用户名</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "ticket",
            "description": "<p>用户接口调用凭据</p>"
          }
        ]
      }
    },
    "success": {
      "examples": [
        {
          "title": "Success-Response:",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 0,\n    \"errmsg\": \"操作成功\",\n    \"data\": [{\n        \"IMEI\": \"123456789101112\",\n        \"nick\": \"123\",\n        \"status\": 1,\n        \"work_model\": 1,\n        \"volume\": 6,\n        \"electricity\": 100,\n        \"emeregncyPhone\": \"13590210000\",\n    }]\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./controllers/device.go",
    "groupTitle": "Device",
    "sampleRequest": [
      {
        "url": "http://api.watch.h-hy.com:8080/v1/device"
      }
    ]
  },
  {
    "type": "get",
    "url": "/system",
    "title": "获取系统信息",
    "name": "systemInfo",
    "group": "System",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "pageUrl",
            "description": "<p>网页地址，用于JSSDK认证</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Request-Example:",
          "content": "pageUrl=http://www.xxx.com/index.html",
          "type": "String"
        }
      ]
    },
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "appId",
            "description": "<p>公众号的唯一标识</p>"
          },
          {
            "group": "Success 200",
            "type": "Number",
            "optional": false,
            "field": "timestamp",
            "description": "<p>生成签名的时间戳</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "nonceStr",
            "description": "<p>生成签名的随机串</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "signature",
            "description": "<p>签名</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Success-Response:",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 0,\n    \"errmsg\": \"操作成功\",\n    \"data\": {\n        \"appId\": \"1234567\",\n        \"timestamp\": 1234567\n        \"nonceStr\": \"nonceStr\",\n        \"signature\": \"signature\"\n    }\n}",
          "type": "json"
        }
      ]
    },
    "error": {
      "examples": [
        {
          "title": "Error-Response:",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 403,\n    \"errmsg\": \"Refer认证失败，请求被拒绝\",\n    \"data\": {\n    }\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./controllers/system.go",
    "groupTitle": "System",
    "sampleRequest": [
      {
        "url": "http://api.watch.h-hy.com:8080/v1/system"
      }
    ]
  },
  {
    "type": "get",
    "url": "/user/:username",
    "title": "查看用户信息",
    "description": "<p>接口中的融云token如果失效，需要调用“刷新融云密钥”接口来刷新</p>",
    "name": "userDetail",
    "group": "User",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "ticket",
            "description": "<p>用户接口调用凭据</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "username",
            "description": "<p>用户名</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "rongCloudAppKey",
            "description": "<p>融云AppKey</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "rongCloudToken",
            "description": "<p>融云token</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Success-Response:",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 0,\n    \"errmsg\": \"操作成功\",\n    \"data\": {\n        \"username\": \"13590210000\"\n        \"rongCloudAppKey\": \"rongCloudAppKey\"\n        \"rongCloudToken\": \"rongCloudToken\"\n    }\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./controllers/user.go",
    "groupTitle": "User",
    "sampleRequest": [
      {
        "url": "http://api.watch.h-hy.com:8080/v1/user/:username"
      }
    ]
  },
  {
    "type": "post",
    "url": "/user/:username/login",
    "title": "用户登录",
    "name": "userLogin",
    "group": "User",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "password",
            "description": "<p>用户密码</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": true,
            "field": "code",
            "description": "<p>微信用户授权码</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Request-Example:",
          "content": "password=111111&code=123",
          "type": "String"
        }
      ]
    },
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "username",
            "description": "<p>用户名</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "ticket",
            "description": "<p>用户接口调用凭据</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "正常回复",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 0,\n    \"errmsg\": \"操作成功\",\n    \"data\": {\n        \"username\": \"13590210000\",\n        \"ticket\": \"abcdefg\"\n    }\n}",
          "type": "json"
        }
      ]
    },
    "error": {
      "examples": [
        {
          "title": "用户名不存在回复",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 20003,\n    \"errmsg\": \"用户名不存在\",\n    \"data\": {\n    }\n}",
          "type": "json"
        },
        {
          "title": "用户密码错误回复",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 20004,\n    \"errmsg\": \"用户密码错误\",\n    \"data\": {\n    }\n}",
          "type": "json"
        },
        {
          "title": "微信用户授权码已失效回复",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 40001,\n    \"errmsg\": \"微信用户授权码已失效\",\n    \"data\": {\n    }\n}\nerrcode=40001的情况请重新发起微信页面授权",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./controllers/user.go",
    "groupTitle": "User",
    "sampleRequest": [
      {
        "url": "http://api.watch.h-hy.com:8080/v1/user/:username/login"
      }
    ]
  },
  {
    "type": "post",
    "url": "/user/:username/logout",
    "title": "用户退出登录",
    "name": "userLogout",
    "group": "User",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "ticket",
            "description": "<p>用户接口调用凭据</p>"
          }
        ]
      }
    },
    "success": {
      "examples": [
        {
          "title": "Success-Response:",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 0,\n    \"errmsg\": \"操作成功\",\n    \"data\": {\n    }\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./controllers/user.go",
    "groupTitle": "User",
    "sampleRequest": [
      {
        "url": "http://api.watch.h-hy.com:8080/v1/user/:username/logout"
      }
    ]
  },
  {
    "type": "post",
    "url": "/user/:username/resetPassword",
    "title": "用户找回密码",
    "name": "userResetPassword",
    "group": "User",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "password",
            "description": "<p>用户新密码</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "SMScode",
            "description": "<p>获取到的短信验证码</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Request-Example:",
          "content": "username=13590210000&password=123456&SMScode=123456",
          "type": "String"
        }
      ]
    },
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "username",
            "description": "<p>用户名</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "ticket",
            "description": "<p>用户接口调用凭据</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Success-Response:",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 0,\n    \"errmsg\": \"密码重置成功\",\n    \"data\": {\n        \"username\": \"13590210000\",\n        \"ticket\": \"abcdefg\"\n    }\n}",
          "type": "json"
        }
      ]
    },
    "error": {
      "examples": [
        {
          "title": "用户名不存在回复",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 20002,\n    \"errmsg\": \"用户名不存在\",\n    \"data\": {\n    }\n}",
          "type": "json"
        },
        {
          "title": "短信验证码错误回复",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 20013,\n    \"errmsg\": \"短信验证码错误\",\n    \"data\": {\n    }\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./controllers/user.go",
    "groupTitle": "User",
    "sampleRequest": [
      {
        "url": "http://api.watch.h-hy.com:8080/v1/user/:username/resetPassword"
      }
    ]
  },
  {
    "type": "get",
    "url": "/user/:username/SMSCode",
    "title": "获取短信验证码",
    "name": "userSMSCode",
    "group": "User",
    "success": {
      "examples": [
        {
          "title": "Success-Response:",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 0,\n    \"errmsg\": \"获取成功\",\n    \"data\": {\n    }\n}",
          "type": "json"
        }
      ]
    },
    "error": {
      "examples": [
        {
          "title": "秒级频率限制提示",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 20010,\n    \"errmsg\": \"短信验证码已发送，请60秒后再试\",\n    \"data\": {\n    }\n}",
          "type": "json"
        },
        {
          "title": "天级频率限制提示",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 20011,\n    \"errmsg\": \"每天最多发送10条验证码短信，请明天再试\",\n    \"data\": {\n    }\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./controllers/user.go",
    "groupTitle": "User",
    "sampleRequest": [
      {
        "url": "http://api.watch.h-hy.com:8080/v1/user/:username/SMSCode"
      }
    ]
  },
  {
    "type": "post",
    "url": "/user/",
    "title": "用户注册",
    "name": "userStore",
    "group": "User",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "username",
            "description": "<p>用户名</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "password",
            "description": "<p>用户密码</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": true,
            "field": "code",
            "description": "<p>微信用户授权码</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Request-Example:",
          "content": "username=13590210000&password=123456",
          "type": "String"
        }
      ]
    },
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "username",
            "description": "<p>用户名</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "ticket",
            "description": "<p>用户接口调用凭据</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Success-Response:",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 0,\n    \"errmsg\": \"注册成功\",\n    \"data\": {\n        \"username\": \"13590210000\",\n        \"ticket\": \"abcdefg\"\n    }\n}",
          "type": "json"
        }
      ]
    },
    "error": {
      "examples": [
        {
          "title": "用户名已经存在回复",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 20002,\n    \"errmsg\": \"用户名已经存在\",\n    \"data\": {\n    }\n}",
          "type": "json"
        },
        {
          "title": "微信用户授权码已失效回复",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 40001,\n    \"errmsg\": \"微信用户授权码已失效\",\n    \"data\": {\n    }\n}\nerrcode=40001的情况请重新发起微信页面授权",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./controllers/user.go",
    "groupTitle": "User",
    "sampleRequest": [
      {
        "url": "http://api.watch.h-hy.com:8080/v1/user/"
      }
    ]
  },
  {
    "type": "put",
    "url": "/user/:username",
    "title": "更新用户信息",
    "description": "<p>本接口只需要传入需要更新的参数即可</p>",
    "name": "userUpdate",
    "group": "User",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "ticket",
            "description": "<p>用户接口调用凭据</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": true,
            "field": "oldPassword",
            "description": "<p>用户旧密码</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": true,
            "field": "newPassword",
            "description": "<p>用户新密码</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Request-Example:",
          "content": "oldPassword=123456&newPassword=111111",
          "type": "String"
        }
      ]
    },
    "success": {
      "examples": [
        {
          "title": "Success-Response:",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 0,\n    \"errmsg\": \"操作成功\",\n    \"data\": {\n    }\n}",
          "type": "json"
        }
      ]
    },
    "error": {
      "examples": [
        {
          "title": "Error-Response:",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 20003,\n    \"errmsg\": \"原密码正确\",\n    \"data\": {\n    }\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./controllers/user.go",
    "groupTitle": "User",
    "sampleRequest": [
      {
        "url": "http://api.watch.h-hy.com:8080/v1/user/:username"
      }
    ]
  },
  {
    "type": "get",
    "url": "/user/:username/updateRongCloudToken",
    "title": "刷新融云密钥",
    "name": "userUpdateRongCloudToken",
    "group": "User",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "ticket",
            "description": "<p>用户接口调用凭据</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "rongCloudAppKey",
            "description": "<p>融云AppKey</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "rongCloudToken",
            "description": "<p>融云token</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Success-Response:",
          "content": "HTTP/1.1 200 OK\n{\n    \"errcode\": 0,\n    \"errmsg\": \"操作成功\",\n    \"data\": {\n        \"rongCloudAppKey\": \"rongCloudAppKey\"\n        \"rongCloudToken\": \"rongCloudToken\"\n    }\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./controllers/user.go",
    "groupTitle": "User",
    "sampleRequest": [
      {
        "url": "http://api.watch.h-hy.com:8080/v1/user/:username/updateRongCloudToken"
      }
    ]
  },
  {
    "type": "get",
    "url": "/wechat/authorize",
    "title": "微信请求授权",
    "name": "wechatAuthorize",
    "group": "Wechat",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "redirectUri",
            "description": "<p>授权后跳转地址（返回参数格式待商议）</p>"
          }
        ]
      }
    },
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "code",
            "description": "<p>用户授权码（注册和登陆时需带上）</p>"
          },
          {
            "group": "Success 200",
            "type": "bool",
            "optional": false,
            "field": "isLogined",
            "description": "<p>用户是否已经登陆（true/false）</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": true,
            "field": "username",
            "description": "<p>用户名（isLogined==true带上）</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": true,
            "field": "ticket",
            "description": "<p>用户接口调用凭据（isLogined==true带上）</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Success-Response:",
          "content": "HTTP/1.1 302 Found\nRedirect to {redirectUri}&code={code}&isLogined={isLogined}&username={username}&ticket={ticket}",
          "type": "json"
        }
      ]
    },
    "error": {
      "examples": [
        {
          "title": "Error-Response:",
          "content": "HTTP/1.1 403 Forbidden",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./controllers/wechat.go",
    "groupTitle": "Wechat"
  },
  {
    "type": "LocationUpdated",
    "url": "/NotApi",
    "title": "融云LocationUpdated",
    "description": "<p>设备定位发生改变时推送</p>",
    "name": "rongcloudLocationUpdated",
    "group": "rongCloud",
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "Number",
            "optional": false,
            "field": "messageId",
            "description": "<p>消息id</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "IMEI",
            "description": "<p>设备IMEI</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "nick",
            "description": "<p>设备昵称</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "toUsername",
            "description": "<p>接收用户名</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "locationType",
            "description": "<p>定位类型,GPS/LBS</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "mapType",
            "description": "<p>经纬度类型，高德地图为amap</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "lat",
            "description": "<p>经度</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "lng",
            "description": "<p>纬度</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "radius",
            "description": "<p>定位精度半径，单位：米</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "created_at",
            "description": "<p>消息创建时间，格式为Y-m-d H:i:s</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Demo",
          "content": "{\n    \"messageId\": 123,\n    \"IMEI\": \"123\",\n    \"nick\": \"456\",\n    \"toUsername\": \"voice\",\n    \"locationType\": \"GPS\",\n    \"mapType\": \"amap\",\n    \"lat\": \"22\",\n    \"lng\": \"11\",\n    \"radius\": \"11\",\n    \"created_at\": \"2016-08-08 11:11:11\"\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./controllers/rongCloud.go",
    "groupTitle": "rongCloud"
  },
  {
    "type": "MessageReceived",
    "url": "/NotApi",
    "title": "融云MessageReceived",
    "description": "<p>用户收到新消息时推送</p>",
    "name": "rongcloudMessageReceived",
    "group": "rongCloud",
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "Number",
            "optional": false,
            "field": "messageId",
            "description": "<p>消息id</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "fromUsername",
            "description": "<p>发送用户名（设备则为IMEIxxxxxxxxxxxxx）</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "toUsername",
            "description": "<p>接收用户名（设备则为IMEIxxxxxxxxxxxxx）</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "type",
            "description": "<p>消息类型，目前只有voice</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": true,
            "field": "mp3Url",
            "description": "<p>Mp3文件地址，消息类型为voice时发送</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "created_at",
            "description": "<p>消息创建时间，格式为Y-m-d H:i:s</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Demo",
          "content": "{\n    \"messageId\": 123,\n    \"fromUser\": \"IMEI123456789101112\",\n    \"toUser\": \"456\",\n    \"type\": \"voice\",\n    \"mp3Url\": \"http://baidu.com/a.mp3\",\n    \"created_at\": \"2016-08-08 11:11:11\",\n}",
          "type": "json"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "./controllers/rongCloud.go",
    "groupTitle": "rongCloud"
  }
] });
