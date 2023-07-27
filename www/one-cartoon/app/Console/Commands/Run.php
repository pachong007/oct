<?php

namespace App\Console\Commands;


use App\Models\Manga\MangaChapter;
use App\Models\Manga\MangaComic;
use App\Models\Pri\Category;
use App\Models\SourceChapter;
use App\Models\SourceComic;
use App\Models\SourceImage;
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
        $comics = SourceComic::where('source',1)
            ->where('created_at','>',date('2023-07-14'))
            ->pluck('id')->all();
        $comics = array_chunk($comics,100);
        foreach ($comics as $comicIds){
            try{
                DB::beginTransaction();
                $mgs = MangaComic::where('source',1)->whereIn('source_comic_id',$comicIds)->select('comic_id','tab')->get();
                $mgIds = [];
                foreach ($mgs as $mg){
                    $MC = (new MangaChapter())->setTable($mg->tab);
                    $MC->where('comic_id',$mg->comic_id)->delete();
                    $mgIds[] = $mg->comic_id;
                }
                MangaComic::whereIn('comic_id',$mgIds)->delete();
                SourceComic::whereIn('id',$comicIds)->update(['status'=>0]);
                SourceChapter::whereIn('comic_id',$comicIds)->update(['status'=>0]);
                DB::commit();
            }catch (\Exception $e){
                DB::rollBack();
            }
        }
    }

    private function cate()
    {
        $baijiaxing = array(
            '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
            'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
            'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
            'u', 'v', 'w', 'x', 'y', 'z',
            '赵', '钱', '孙', '李', '周', '吴', '郑', '王', '冯', '陈',
            '褚', '卫', '蒋', '沈', '韩', '杨', '朱', '秦', '尤', '许',
            '何', '吕', '施', '张', '孔', '曹', '严', '华', '金', '魏',
            '陶', '姜', '戚', '谢', '邹', '喻', '柏', '水', '窦', '章',
            '云', '苏', '潘', '葛', '奚', '范', '彭', '郎', '鲁', '韦',
            '昌', '马', '苗', '凤', '花', '方', '俞', '任', '袁', '柳',
            '酆', '鲍', '史', '唐', '费', '廉', '岑', '薛', '雷', '贺',
            '倪', '汤', '滕', '殷', '罗', '毕', '郝', '邬', '安', '常',
            '乐', '于', '时', '傅', '皮', '卞', '齐', '康', '伍', '余'
        );
        $categories = Category::where('id','>',0)->get();
        foreach ($categories as $k=>$category){
            Category::where('id',$category->id)->update(['letter'=>$baijiaxing[$k]]);
        }
        $categories = Category::pluck('letter','id')->all();
        $comics = MangaComic::where('comic_id','>',0)->get();
        foreach ($comics as $comic){
            $sorts = explode(',',$comic->sorts);
            $letters = [];
            foreach ($sorts as $sort){
                if(isset($categories[$sort])){
                    $letters[] = $categories[$sort];
                }
            }
            $letters = join(',',$letters);
            MangaComic::where('comic_id',$comic->comic_id)->update(['sorts_letter'=>$letters]);
        }
    }

    private function tx()
    {
        $redis = Redis::connection('tx');
        $list = $redis->lrange("source:comic:chapter",0,-1);

        $chapters = SourceChapter::where('source',2)->where('status',0)->get();
        foreach ($chapters as $chapter){
            if($chapter->is_free == 1) continue;
            $img = SourceImage::where('chapter_id',$chapter->id)->first();
            if(!$img || $img->state == 0){
                $list[] = $chapter->id;
            }
        }

        $newList = array_unique($list);
        $redis->del("source:comic:chapter");
        foreach ($newList as $id){
            $chapter = SourceChapter::where('id',$id)->first();
            if($chapter) {
                if($chapter->is_free == 1) continue;
                if(SourceImage::where('chapter_id',$id)->where('state',1)->exists()){
                    continue;
                }
                $redis->rpush("source:comic:chapter", $id);
            }
        }
    }

    private function kk()
    {
        $redis = Redis::connection('kk');
        $list = $redis->lrange("source:comic:chapter",0,-1);

        $chapters = SourceChapter::where('source',1)->where('status',0)->get();
        foreach ($chapters as $chapter){
            if($chapter->is_free == 1) continue;
            $img = SourceImage::where('chapter_id',$chapter->id)->first();
            if(!$img || $img->state == 0){
                $list[] = $chapter->id;
            }
        }

        $newList = array_unique($list);
        $redis->del("source:comic:chapter");
        foreach ($newList as $id){
            $chapter = SourceChapter::where('id',$id)->first();
            if($chapter) {
                if($chapter->is_free == 1) continue;
                if(SourceImage::where('chapter_id',$id)->where('state',1)->exists()){
                    continue;
                }
                $redis->rpush("source:comic:chapter", $id);
            }
        }
    }

    private function ctt()
    {
        for ($i=0;$i<=256;$i++){
            DB::statement("ALTER TABLE manhua.`bw_comic_chapter_$i` ADD COLUMN image_list JSON;");
        }
    }

    private function ct()
    {
        SourceComic::where('id','>',0)->update(['status'=>0]);
        SourceChapter::where('id','>',0)->update(['status'=>0]);
        for ($i=0;$i<=256;$i++){
            DB::statement("
CREATE TABLE manhua.`bw_comic_chapter_$i` (
`id` INT ( 11 ) UNSIGNED NOT NULL AUTO_INCREMENT,
`title` VARCHAR ( 500 ) NOT NULL DEFAULT '' COMMENT '章节名',
`cover` VARCHAR ( 500 ) NOT NULL DEFAULT '' COMMENT '封面',
`c_cover` VARCHAR ( 500 ) DEFAULT NULL COMMENT '加密封面',
`comic_id` INT ( 11 ) UNSIGNED NOT NULL DEFAULT '0' COMMENT '漫画id',
`source_comic_id` INT ( 11 ) NOT NULL DEFAULT '0',
`source_chapter_id` INT ( 11 ) NOT NULL DEFAULT '0',
`source_image_id` INT ( 11 ) NOT NULL DEFAULT '0',
`images` json DEFAULT NULL,
`display_order` INT ( 11 ) UNSIGNED NOT NULL DEFAULT '0' COMMENT '章节序号',
`is_vip` TINYINT ( 1 ) UNSIGNED NOT NULL DEFAULT '0' COMMENT '是否vip章节 0不是 1是',
`img_type` TINYINT ( 1 ) UNSIGNED NOT NULL DEFAULT '1',
`total_photos` INT ( 11 ) UNSIGNED NOT NULL DEFAULT '0' COMMENT '章节图片数',
`ip` VARCHAR ( 128 ) NOT NULL DEFAULT '' COMMENT '上传ip',
`status` TINYINT ( 3 ) UNSIGNED NOT NULL DEFAULT '1',
`created_at` int(11) unsigned NOT NULL COMMENT '创建时间',
`updated_at` int(11) unsigned NOT NULL COMMENT '更新时间',
`deleted_at` int(11) unsigned NOT NULL COMMENT '删除时间',
`is_book_coupon_pay` TINYINT ( 1 ) DEFAULT '0' COMMENT '漫画章节书券解锁  默认为0 免费观看  1为需要书券解锁观看',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_comic` (`comic_id`) USING BTREE,
  KEY `idx_display` (`display_order`) USING BTREE,
  KEY `idx_create` (`created_at`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='漫画章节表';
");

        }
    }
}
