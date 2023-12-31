<?php

namespace App\Console\Commands;

use App\Models\Chapter;
use App\Models\Comic;
use App\Models\ComicType;
use App\Models\Dbs;
use App\Models\Image;
use App\Models\McType;
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

    private $comicLimit = 5;
    private $chapterLimit = 20;

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
        $limit = 500;
        $sources = SourceComic::inRandomOrder()->where('chapter_count', '>', 0)->limit($limit)->get();

        $DBs = Dbs::select('name')->get();
        foreach ($DBs as $D) {
            $db = $D->name;
            config(["database.connections.mysql_${db}" => [
                'driver' => 'mysql',
                'host' => '107.148.191.71',
                'database' => $db,
                'username' => 'root',
                'password' => 'docker@6603',
                'charset' => 'utf8',
                'collation' => 'utf8_unicode_ci'
            ]]);

            $comic = new Comic();
            $comic->setConnection("mysql_${db}");
            $mct = new McType();
            $mct->setConnection("mysql_${db}");
            $ct = new ComicType();
            $ct->setConnection("mysql_${db}");

            $comicLimit = $this->comicLimit;

            foreach ($sources as $sourceComic) {
                if($comicLimit < 0)break;
                /*章节检查*/
                $chapterDone = SourceChapter::join('source_image', 'source_chapter.id', '=', 'source_image.chapter_id')
                    ->where('source_image.state', 1)
                    ->where('source_chapter.comic_id',$sourceComic->id)
                    ->count();
                if ($chapterDone === 0) {
                    continue;
                }
                $publish = Publish::where(['database' => $db, 'comic_id' => $sourceComic->id])->first();
                if (!$publish) {
                    $comicLimit--;
                    $mccComicId = $comic->insertGetId([
                        'name' => $sourceComic->title,
                        'yname' => '',
                        'pic' => rtrim(env("IMG_DOMAIN"), "/") . "/" . $sourceComic->cover,
                        'cid' => mt_rand(1, 4),
                        'serialize' => $sourceComic->is_finish == 1 ? '完结' : '连载',
                        'author' => $sourceComic->author,
                        'content' => $sourceComic->description,
                        'nums' => $sourceComic->chapter_count,
                        'score' => mt_rand(1, 9),
                        'did' => $sourceComic->id,
                        'ly' => 'kk',
                        'addtime' => time(),
                    ]);

                    $tag = $mct->where('fid', 1)->inRandomOrder()->first();
                    $tag2 = $mct->where('fid', 2)->inRandomOrder()->first();
                    $tag3 = $mct->where('fid', 3)->inRandomOrder()->first();
                    $tag4 = $mct->where('fid', 4)->inRandomOrder()->first();
                    $tag5 = $mct->where('fid', 5)->inRandomOrder()->first();
                    $typeInsert = [
                        ['mid' => $mccComicId, 'tid' => $tag->id],
                        ['mid' => $mccComicId, 'tid' => $tag2->id],
                        ['mid' => $mccComicId, 'tid' => $tag3->id],
                        ['mid' => $mccComicId, 'tid' => $tag4->id],
                        ['mid' => $mccComicId, 'tid' => $tag5->id],
                    ];
                    $ct->insert($typeInsert);

                    $pid = Publish::insertGetId([
                        'comic_id' => $sourceComic->id,
                        'chapter_id' => json_encode([]),
                        'source' => $sourceComic->source,
                        'database' => $db,
                        'publish_id' => $mccComicId,
                        'publish_chapter_id' => json_encode([])
                    ]);
                    $publish = Publish::where(['id' => $pid])->first();
                }

                if ($chapterDone <= count($publish->publish_chapter_id)) {
                    continue;
                }
                $res = $this->insertChapter($db, $sourceComic->id, $publish->publish_id, $publish->chapter_id);
                Publish::where('id',$publish->id)->update(['chapter_id'=>json_encode($res[0]),
                    'publish_chapter_id'=>json_encode($res[1])]);
                continue;
            }
        }
    }

    private function insertChapter($db, $comicId, $mccComicId, $chapterIds)
    {
        $chapterLimit = $this->chapterLimit;
        $chapters = SourceChapter::where('comic_id', $comicId)->get();
        $cha = new Chapter();
        $cha->setConnection("mysql_${db}");
        $img = new Image();
        $img->setConnection("mysql_${db}");

        $mccChapterIds = [];
        foreach ($chapters as $chapter) {
            if(in_array($chapter->id,$chapterIds))continue;
            if ($chapter->image && $chapter->image['state'] == 1) {
                $chapterLimit--;
                if($chapterLimit<0)break;
                $images = $chapter->image['images'];
                $cid = $cha->insertGetId([
                    'mid' => $mccComicId,
                    'xid' => $chapter->sort,
                    'name' => $chapter->title,
                    'jxurl' => $chapter->source_url,
                    'pnum' => count($images),
                    'addtime' => time(),
                ]);

                $insertImages = [];
                foreach ($images as $k => $image) {
                    $insertImages[] = [
                        'cid' => $cid,
                        'mid' => $mccComicId,
                        'img' => rtrim(env("IMG_DOMAIN"), "/") . "/" . $image,
                        'xid' => $k,
                        'md5' => ''
                    ];
                }
                if (!empty($insertImages)) {
                    $img->insert($insertImages);
                }
                $mccChapterIds[] = $cid;
                $chapterIds[] = $chapter->id;
            }
        }
        return [$chapterIds,$mccChapterIds];
    }
}
