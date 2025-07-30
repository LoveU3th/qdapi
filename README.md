## 使用方法
起点app随便访问一下福利中心  
抓包`https://h5.if.qidian.com`任意链接中的以下四个值:  
+ header中的  
  `SDKSign`
+ header->cookie中的  
  `QDInfo`
  `ywkey` (每月需要更新一次)
  `ywguid`
放入`config.json`中

### `config.json`解释
+ 示例文件是多账号(同时处理每个账号任务),请自行增减
+ `SDKSign,QDInfo,ywkey,ywguid`为抓包数据
+ `TaskType`是要运行的任务
  + 1:一小时一个的宝箱
  + 2:每天的8个任务
    + 其实可以瞬间执行完,但是不知道他有没有风险控制,暂定15s间隔(`task.go L24`),需要2min(8*15s)
  + 3:看3个得10点的任务
  + 4:更多任务  
    + 目前仅支持**前往游戏中心玩游戏10分钟奖励10点币**(每30s给一次心跳累计20次,需要10min)
    + 由于没有破解imei相关的加密,多账号时只能使用多台手机的cookie,否则后台会不计算时长
  + 所有任务并行同时处理,如2需要8*15s=2min,4需要10min.同时运行2,4需要10min

### 本地运行  
+ [github](https://github.com/skeeto/w64devkit/releases)
+ [lanzouq](https://wwtw.lanzouq.com/icI8r322p9na)
+ [cdn](https://app.parap.dpdns.org/qdapi/qdapi.zip)
+ 修改`config.json`中数据并运行`qdapi.exe`/`run.bat`

### github action 每天自动运行(推荐)
+ fork[本项目](https://github.com/pzx521521/qdapi)
+ **推荐方式**：使用 GitHub Actions 的 Secrets 功能保护配置信息
  + 详细配置步骤请查看 [SECRETS_SETUP.md](./SECRETS_SETUP.md)
  + 这样可以避免配置信息泄露，更加安全
+ **传统方式**：直接修改 `config.json` 文件
  + 如果使用此方式，建议将仓库设为私有

### 其他平台的运行
+ 修改`config.json`中数据并运行cmd/main.go
+ vercel-function
+ nas/路由器/手机,可以编译对应平台执行文件
+ 云函数对时间有要求的请删掉任务4(前往游戏中心玩游戏10分钟奖励10点币)

## 其他注意事项
+ 如果提示`领取失败，请升级至最新版本领取`是版本号低了,自行代码修改版本号或者抓包最新app更新`QDInfo`和`SDKSign`
+ 可能会导致任何的损失(如封号),概不负责

## 破解流程如下:
https://www.jianshu.com/p/58ec69e04983
QDInfo 算法如下:

```python
import base64
from Crypto.Cipher import DES3
key = b'0821CAAD409B84020821CAAD'
data = base64.b64decode('SO+aPyWTJ02k4C9FkkB29fACDXIsJx4pAGbhVI07D8hjHPOEsCFgpJ99gS3kYIjunO+UrcWbhPgIlUSo3XxdoisFnouWF80qfP+9nYAPZWuWE/x7ukJhxq8DEJW+n90UAoC6t3e9KFYaJ/yFFUfggDVS6xpzIkTxCCDps2WxRBcdvOXoA5I5/i3jrw8wJqw0DmbxzkSOoKB1T5VHx/VjWCoYTuW8fA5DlGMQL+4lQldYUANNM1Aarp6oD16p7Rqc9JpGyHOOnKF3tDxv8vGv0ElZszGBKKqK70o3d0OzvfmgFhyXErR92g==')
cryptor = DES3.new(key, DES3.MODE_CBC, b'00000000')
print(cryptor.decrypt(data))
```
SDKSign 算法如下:

```python
import base64

from Crypto.Cipher import DES3

key = b'8YV#U2Butm,VutR2B_W[*}6t'
data = base64.b64decode('fwU0VSlfsV/NtCFBjpJarbYpi9mlbLU/EDzhOVoz2RdtheX+SLpjTy8L2+gA InschgJSs1O5vbtFpSZ6+GPI8iEd6QhtwlTz8ODKLNM1r+aH0A8sY5+lP6la DPt/GpDgPvW5ZvKHiqnIqFEJHRoPYEshR2+cAq03JfcYLPvSfE7DpuHLVA2F mRtLGCdVWmTujc/5Lb+/Cmk=')
cryptor = DES3.new(key, DES3.MODE_CBC, b'01234567')
print(cryptor.decrypt(data))
```
# Todo
+ [ ] ywkey 每月需要更新,目前没有自动每月更新
+ MoreRewardTab
+ [x] 104=前往游戏中心玩游戏10分钟奖励10点币
  + 目前没有破解imei相关的加密,多账号时只能使用多台手机的cookie,否则后台会不计算时长

[![pVYaMmd.jpg](https://s21.ax1x.com/2025/07/28/pVYaMmd.jpg)](https://imgse.com/i/pVYaMmd)

已经稳定用了一年了,放出来不知道会不会被修