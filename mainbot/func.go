package mainbot

import "github.com/0ojixueseno0/go-Pichubot-base/pichumod"

//* Message Event
//? 消息事件
var OnPrivateMsg func(eventinfo pichumod.MessagePrivate) = func(eventinfo pichumod.MessagePrivate) {} // 私聊消息事件
var OnGroupMsg func(eventinfo pichumod.MessageGroup) = func(eventinfo pichumod.MessageGroup) {}       // 群聊消息事件

//* Notice Event
//? 提醒事件
var OnGroupUpload func(eventinfo pichumod.GroupUpload) = func(eventinfo pichumod.GroupUpload) {}       // 群文件上传
var OnGroupAdmin func(eventinfo pichumod.GroupAdmin) = func(eventinfo pichumod.GroupAdmin) {}          // 群管理员变动
var OnGroupDecrease func(eventinfo pichumod.GroupDecrease) = func(eventinfo pichumod.GroupDecrease) {} // 群成员减少
var OnGroupIncrease func(eventinfo pichumod.GroupIncrease) = func(eventinfo pichumod.GroupIncrease) {} // 群成员增加
var OnGroupBan func(eventinfo pichumod.GroupBan) = func(eventinfo pichumod.GroupBan) {}                // 群聊禁言
var OnFriendAdd func(eventinfo pichumod.FriendAdd) = func(eventinfo pichumod.FriendAdd) {}             // 已经添加好友后的事件
var OnGroupRecall func(eventinfo pichumod.GroupRecall) = func(eventinfo pichumod.GroupRecall) {}       // 群消息撤回(群聊)
var OnFriendRecall func(eventinfo pichumod.FriendRecall) = func(eventinfo pichumod.FriendRecall) {}    // 好友消息撤回(私聊)
var OnNotify func(eventinfo pichumod.Notify) = func(eventinfo pichumod.Notify) {}                      // 群内戳一戳 群红包运气王 群成员荣誉变更

//* Request Event
//? 请求事件
var OnFriendRequest func(eventinfo pichumod.FriendRequest) = func(eventinfo pichumod.FriendRequest) {} // 加好友请求
var OnGroupRequest func(eventinfo pichumod.GroupRequest) = func(eventinfo pichumod.GroupRequest) {}    // 加群请求/邀请

//* Meta Event
//? 元事件
var OnMetaLifecycle func(eventinfo pichumod.MetaLifecycle) = func(eventinfo pichumod.MetaLifecycle) {} // 生命周期
var OnMetaHeartbeat func(eventinfo pichumod.MetaHeartbeat) = func(eventinfo pichumod.MetaHeartbeat) {} // 心跳包
