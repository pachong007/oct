<?php

namespace App\Console\Commands;


use App\Models\Manga\MangaComic;
use App\Models\Pri\Category;
use App\Models\Pri\Tag;
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
        $limit = 200;

        SourceComic::where('status', 2)
            ->where('updated_at', '<', date('Y-m-d',strtotime('14 days ago')))
            ->update(['status' => 0]);

        $time = time();
        $categories = Category::pluck('id', 'title')->all();
        $letters = Category::pluck('letter', 'id')->all();
        $tags = Tag::pluck('id', 'title')->all();
        $sources = SourceComic::where('status', 0)->where('chapter_count', '>', 0)->offset($page * $limit)->limit($limit)->orderBy('created_at', 'ASC')->get();
        foreach ($sources as $sourceComic) {
            if (MangaComic::where('old_name', trim($sourceComic->title))->exists()) {
                SourceComic::where('id', $sourceComic->id)->update(['status' => 3, 'updated_at' => date('Y-m-d H:i:s')]);
                continue;
            }
            /*章节检查*/
            $result = DB::select("SELECT CountChapterDone({$sourceComic->id}) as result")[0];
            if ($result->result == 0) {
                SourceComic::where('id', $sourceComic->id)->update(['status' => 2, 'updated_at' => date('Y-m-d H:i:s')]);
                continue;
            }

            if (isset($categories[$sourceComic->category])) {
                $cateId = $categories[$sourceComic->category];
            } else if ($sourceComic->category == '') {
                $cateId = $categories['其它'];
            } else {
                Category::insert(['title' => $sourceComic->category]);
                $categories = Category::pluck('title', 'id')->all();
                $cateId = $categories[$sourceComic->category];
                $letters = Category::pluck('letter', 'id')->all();
            }

            $tag_ids = [];
            foreach ($sourceComic->label as $label) {
                if (isset($tags[$label])) {
                    $tag_ids[] = $tags[$label];
                } else if ($label == '') {
                    continue;
                } else {
                    Tag::insert(['title' => $label]);
                    $tags = Tag::pluck('id', 'title')->all();
                    $tag_ids[] = $tags[$label];
                }
            }
            $horizontal_cover = '';
            if($sourceComic->source == 1){
                $horizontal_cover = $sourceComic->cover_h;
            }else if($sourceComic->source == 2){
                $horizontal_cover = $sourceComic->cover;
            }
            $mangaId = MangaComic::insertGetId([
                'source_comic_id' => $sourceComic->id,
                'sort_id' => $cateId,
                'sorts' => $cateId,
                'sorts_letter' => isset($letters[$cateId]) ? $letters[$cateId] : '',
                'type' => 2,
                'name' => trim($sourceComic->title),
                'old_name' => trim($sourceComic->title),
                'source' => $sourceComic->source,
                'author' => $sourceComic->author,
                'description' => $sourceComic->description,
                'keywords' => $sourceComic->title,
                'tags' => join(',', $tag_ids),
                'vertical_cover' => $sourceComic->cover,
                'horizontal_cover' => $horizontal_cover,
                'status' => 1,
                'issue_time' => strtotime($sourceComic->created_at),
                'is_finish' => $sourceComic->is_finish,
                'created_at' => $time,
                'updated_at' => $time
            ]);
            MangaComic::where('comic_id', $mangaId)->update(['tab' => $mangaId % 256]);
            SourceComic::where('id', $sourceComic->id)->update(['status' => 1]);
        }

    }

}
