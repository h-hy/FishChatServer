package controllers

/**
* @api {MessageReceived} /NotApi 融云MessageReceived
 * @apiDescription 用户收到新消息时推送
 * @apiName rongcloudMessageReceived
 * @apiGroup rongCloud
 * @apiSampleRequest off
 *
 * @apiSuccess {Number} messageId 消息id
 * @apiSuccess {String} fromUsername 发送用户名（设备则为IMEIxxxxxxxxxxxxx）
 * @apiSuccess {String} toUsername 接收用户名（设备则为IMEIxxxxxxxxxxxxx）
 * @apiSuccess {String} type 消息类型，目前只有voice
 * @apiSuccess {String} [mp3Url] Mp3文件地址，消息类型为voice时发送
 * @apiSuccess {String} created_at 消息创建时间，格式为Y-m-d H:i:s
 *
 * @apiSuccessExample Demo
 *     {
 *         "messageId": 123,
 *         "fromUser": "IMEI123456789101112",
 *         "toUser": "456",
 *         "type": "voice",
 *         "mp3Url": "http://baidu.com/a.mp3",
 *         "created_at": "2016-08-08 11:11:11",
 *     }
*/

/**
* @api {LocationUpdated} /NotApi 融云LocationUpdated
 * @apiDescription 设备定位发生改变时推送
 * @apiName rongcloudLocationUpdated
 * @apiGroup rongCloud
 * @apiSampleRequest off
 *
 * @apiSuccess {Number} messageId 消息id
 * @apiSuccess {String} IMEI 设备IMEI
 * @apiSuccess {String} nick 设备昵称
 * @apiSuccess {String} toUsername 接收用户名
 * @apiSuccess {String} locationType 定位类型,GPS/LBS
 * @apiSuccess {String} mapType 经纬度类型，高德地图为amap
 * @apiSuccess {String} lat 经度
 * @apiSuccess {String} lng 纬度
 * @apiSuccess {String} radius 定位精度半径，单位：米
 * @apiSuccess {String} created_at 消息创建时间，格式为Y-m-d H:i:s
 *
 * @apiSuccessExample Demo
 *     {
 *         "messageId": 123,
 *         "IMEI": "123",
 *         "nick": "456",
 *         "toUsername": "voice",
 *         "locationType": "GPS",
 *         "mapType": "amap",
 *         "lat": "22",
 *         "lng": "11",
 *         "radius": "11",
 *         "created_at": "2016-08-08 11:11:11"
 *     }
*/
