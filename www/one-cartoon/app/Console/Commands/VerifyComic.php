<?php

namespace App\Console\Commands;

use App\Models\Chapter;
use App\Models\Comic;
use App\Models\Dbs;
use App\Models\Image;
use App\Models\Publish;
use App\Models\SourceChapter;
use App\Models\SourceComic;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;

class VerifyComic extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'verify.comic';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '自动审核漫画';

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
        $page = 0;
        $limit = 2;

        SourceComic::where('status', 2)
            ->where('updated_at', '<', date('Y-m-d', strtotime('7 days ago')))
            ->update(['status' => 0]);
        $sources = SourceComic::where('status', 0)->where('chapter_count', '>', 0)->offset($page * $limit)->limit($limit)->orderBy('created_at', 'ASC')->get();

        $DBs = Dbs::select('name')->get();
        foreach ($DBs as $D) {
            $db = $D->name;
            config(["database.connections.mysql_${db}" => [
                'driver'    => 'mysql',
                'host'      => '107.148.191.71',
                'database'  => $db,
                'username'  => 'root',
                'password'  => 'docker@6603',
                'charset'   => 'utf8',
                'collation' => 'utf8_unicode_ci'
            ]]);

            foreach ($sources as $sourceComic) {
                /*章节检查*/
                $chapterDone = SourceChapter::join('source_image', 'source_chapter.id', '=', 'source_image.chapter_id')->where('source_image.state', 1)->count();
                if ($chapterDone === 0) {
                    SourceComic::where('id', $sourceComic->id)->update(['status' => 2, 'updated_at' => date('Y-m-d H:i:s')]);
                    continue;
                }
                $publish = Publish::where(['database'=>$db,'publish_id'=>$sourceComic->id])->first();
                if(!$publish){
                    $comic = new Comic();
                    $comic->setConnection("mysql_${db}");
                    $comic->insertGetId([
                       'name'=>$sourceComic->title,
                        'yname'=>'',
                        'pic'=>rtrim(env("IMG_DOMAIN"),"/")."/".$sourceComic->cover,
                        'cid'=>1,
                        'serialize'=>$sourceComic->is_finish == 1 ? '完结':'连载',
                        'author'=>$sourceComic->author,
                        'content'=>$sourceComic->description,
                        'nums'=>$sourceComic->chapter_count,
                        'score'=>mt_rand(1,9),
                        'did'=>$sourceComic->id,
                        'ly'=>'kk',
                    ]);
                    Publish::insert([
                        'comic_id'=>$sourceComic->id,
                        'chapter_id'=>json_encode([]),
                        'source'=>$sourceComic->source,
                        'database'=>$db,
                        'publish_id'=>0,
                        'publish_chapter_id'=>json_encode([])
                    ]);
                    $publish = Publish::where(['database'=>$db,'publish_id'=>$sourceComic->id])->first();
                }
                if($publish){
                    if($chapterDone <= count($publish->publish_chapter_id)){
                        continue;
                    }
                    $this->insertChapter($db,$publish->publish_id,$publish->chapter_id);
                    continue;
                }
            }
        }
    }

    private function insertChapter($db,$mid,$chapterIds)
    {
        $chapters = SourceChapter::whereNotIn('id',$chapterIds)
            ->get()->toArray();
        foreach ($chapters as $chapter){
            if(!empty($chapter['image']) && $chapter['image']['state'] == 1){
                $images = $chapter['image']['images'];
                $cha = new Chapter();
                $cha->setConnection("mysql_${db}");
                $cid = $cha->insertGetId([
                    'mid'=>$mid,
                    'xid'=>$chapter['sort'],
                    'name'=>$chapter['title'],
                    'jxurl'=>$chapter['source_url'],
                    'pnum'=>count($images),
                    'addtime'=>time(),
                ]);

                $insertImages = [];
                foreach ($images as $k=>$image){
                    $insertImages[] = [
                      'cid'=>$cid,
                      'mid'=>$mid,
                      'img'=>rtrim(env("IMG_DOMAIN"),"/")."/".$image,
                      'xid'=>$k,
                      'md5'=>''
                    ];
                }
                if(!empty($insertImages)) {
                    $img = new Image();
                    $img->setConnection("mysql_${db}");
                    $img->insert($insertImages);
                }
            }
        }
    }
}
