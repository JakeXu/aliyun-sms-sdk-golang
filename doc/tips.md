# 使用阿里云短信服务注意事项

### 名词解释

以下几个变量名称在调用 `Send` 函数时使用）

`phones`: 用户接收的手机号, 多个手机号时以逗号 `,` 分隔

`signName`: 签名名称, 在管理后台获取, 可自行创建

`templateCode`: 模板CODE, 在管理后台获取, 可自行创建

`paramStr`: 参数字符串, 序列化的JSON格式

### 准备工作

[直达帮助页面](https://help.aliyun.com/document_detail/44346.html?spm=5176.doc44348.6.103.z0JAmF)

1. 新建短信签名

    用户需要先`新建短信签名`, 阿里云审核通过后会得到一个`签名名称`, 此`签名名称`即为`signName`

2. 新建模板

    * 用户需要先`新建模板`, 阿里云审核通过后会得到一个`模板CODE`, 此`模板CODE`即为`templateCode`

    * 用户在创建模板的时候，会在模板中添加变量（一个或多个），如下为一个例子：

        ```
          尊敬的用户，您的${device_id}(${device_name})设备已离线

        ```

        上面的`device_id`和`device_name`既是模板变量，用户在使用此SDK时，需要把如下所示转换为JSON字符串格式：

        ```
          {"device_id":"T0000001","device_name":"一号设备"}

        ```

        以上字符串即为 `paramStr`.

        下面提供一个golang方式转换方法：

        ```
          A. 为每一个模板CODE建立一个模板参数结构体，实现一个String()方法

            type AlarmOfflineDevice struct {
            	DeviceId   string `json:"device_id"`
            	DeviceName string `json:"device_name"`
            }
            
            func (a *AlarmOfflineDevice) String() string {
            	body, err := json.Marshal(a)
            	if err != nil {
            		return ""
            	}
            	return string(body)
            }

          B. 在设置paramStr参数时，直接使用String()方法产生：

          	d := new(AlarmOfflineDevice)
          	d.DeviceId = "T0000001"
          	d.DeviceName = "测试设备"
          	paramStr ：= d.String()

          	// 这里paramstring将会是符合要求的结果：{"device_id":"T0000001","device_name":"一号设备"}
          	fmt.Println(paramStr)
        ```



