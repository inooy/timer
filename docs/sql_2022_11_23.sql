CREATE TABLE `timer_delay_task`
(
    `id`          bigint UNSIGNED                                               NOT NULL AUTO_INCREMENT COMMENT '消息id',
    `topic`       varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci  NULL DEFAULT '' COMMENT '主题，可以用作业务名称',
    `body`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT '' COMMENT '消息体',
    `retry`       int                                                           NULL DEFAULT 0 COMMENT '失败重试次数，-1表示一直重试，0表示不重试',
    `retry_count` int                                                           NULL DEFAULT 0 COMMENT '已经重试次数',
    `delay`       bigint                                                        NULL DEFAULT 0 COMMENT '任务延迟时间，单位：毫秒',
    `status`      tinyint                                                       NULL DEFAULT 1 COMMENT '当前消息状态：1：延迟队列中 2：就绪中 3：消费中 4：消费成功 5重试队列 6死亡队列',
    `create_time` bigint                                                        NULL DEFAULT 0 COMMENT '消息创建时间',
    `update_time` bigint                                                        NULL DEFAULT 0 COMMENT '更新时间',
    `version`     int                                                           NULL DEFAULT 0 COMMENT '数据版本',
    PRIMARY KEY (`id`)
);