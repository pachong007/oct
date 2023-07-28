<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Redis;

class Run extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'run';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '';

    /**
     * Create a new command instance.
     *
     * @return void
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * Execute the console command.
     *
     * @return mixed
     */
    public function handle()
    {
        $this->dataInsert();
    }

    private function dataInsert()
    {
        DB::statement("DROP TABLE IF EXISTS `source_comic`;");
        DB::statement("
CREATE TABLE `source_comic` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `sid` int(11) NOT NULL DEFAULT '0',
  `source` tinyint(1) NOT NULL DEFAULT '1' COMMENT '采集源 1:快看 2:腾讯',
  `source_id` int(11) NOT NULL COMMENT '源漫画id',
  `source_url` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '源url',
  `cover` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '封面',
  `title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标题',
  `author` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '作者',
  `label` json NOT NULL COMMENT '标签',
  `category` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类',
  `region` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '地区',
  `chapter_count` int(11) NOT NULL DEFAULT '0' COMMENT '章节数量',
  `chapter_count_download` int(11) NOT NULL DEFAULT '0' COMMENT '章节数量(已下载)',
  `like` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '喜欢',
  `popularity` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '0' COMMENT '人气热度',
  `is_free` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0免费 1收费',
  `is_finish` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0连载 1完结',
  `description` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '描述',
  `source_data` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '源数据',
  `chapter_pick` int(11) NOT NULL DEFAULT '0' COMMENT '章节拨片',
  `retry` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0默认 1重抓',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0未审核 1通过',
  `last_chapter_update_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最新章节更新时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `source_id` (`source`,`source_id`) USING BTREE,
  UNIQUE KEY `source_uri` (`source_url`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=15936 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='采集-漫画';
");

        DB::statement("DROP TABLE IF EXISTS `source_chapter`;");
        DB::statement("
CREATE TABLE `source_chapter` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `sid` int(11) NOT NULL DEFAULT '0',
  `scid` int(11) NOT NULL DEFAULT '0',
  `comic_id` int(11) NOT NULL,
  `source` tinyint(1) NOT NULL DEFAULT '1' COMMENT '采集源 1:快看 2:腾讯',
  `source_chapter_id` int(11) NOT NULL COMMENT '源章节id',
  `source_url` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '源url',
  `cover` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '封面',
  `title` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标题',
  `sort` int(11) NOT NULL DEFAULT '0',
  `is_free` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0免费 1收费',
  `source_data` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0未审核 1通过',
  `retry` tinyint(1) NOT NULL DEFAULT '0',
  `view_type` TINYINT (1) NOT NULL DEFAULT '0' COMMENT '0条漫 1页漫'
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `source` (`source`,`comic_id`,`source_url`) USING BTREE,
  KEY `comic_id` (`comic_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=258616 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='采集-漫画章节';
        ");

        DB::statement("DROP TABLE IF EXISTS `source_image`;");
        DB::statement("
DROP TABLE IF EXISTS `source_image`;
CREATE TABLE `source_image` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `source` tinyint(1) DEFAULT '0' COMMENT '采集源 1:快看 2:腾讯',
  `comic_id` int(11) NOT NULL DEFAULT '0',
  `chapter_id` int(11) NOT NULL,
  `images` json NOT NULL,
  `source_data` json NOT NULL,
  `state` tinyint(1) NOT NULL DEFAULT '0' COMMENT '资源获取:0未开始 1已完成',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `chapter_id` (`chapter_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=363 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;
        ");

        DB::statement("DROP TABLE IF EXISTS `fail_info`;");
        DB::statement("
CREATE TABLE `fail_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `source` tinyint(1) NOT NULL DEFAULT '1' COMMENT '采集源 1:快看 2:腾讯',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0漫画列表 1漫画 2章节 3图片',
  `err` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '错误关键词',
  `url` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '地址',
  `info` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '失败信息记录',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `type` (`type`) USING BTREE,
  KEY `err` (`err`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;
        ");
    }
}
