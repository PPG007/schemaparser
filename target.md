# ec.prepaidValueCardOrder （储值卡订单）

```js
{
    _id: ObjectId,
    accountId: ObjectId,
    createdAt: DateTime,
    updatedAt: DateTime,
    isDeleted: Boolean,
    memberId: ObjectId,
    number: String, // 订单编号
    totalAmount: Long, // 订单总金额
    payAmount: Long, // 实际支付金额
    tradeNo: String, // 交易流水号
    paidAt: DateTime, // 支付时间
    completedAt: DateTime, // 交易完成时间
    remarks: String, // 订单备注
    channel: {
        channelId: String,
        openId: String,
    },
    prepaidValueCards: [
        {
            id: ObjectId, // ec.storedValueCard._id
            total: Long, // 购买数量
            amount: Long, // 面值
            number: String, // 编号
            name: String, // 储值卡名称
            picture: String,
        },
    ]
}
```
