# calc24-muti
一个基于UDP协议的多人算24点游戏
----------------

#Feature

 - 随机产生4个小于14的数字作为题目
 - 自动判断是否有解，若否则重新出题
 - 登录无状态，用户名作为唯一标识符
 - 先算出答案的为胜利者

#Usage

    go get github.com/xuzhenglun/calc24-muti
    go build service.go
    go build client.go


# Todo：

 - 存在各种BUG，服务端对数据输入做的校验不多
 - 客户端界面不友好
 


----------


happy hacking
