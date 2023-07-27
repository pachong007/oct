<?php

namespace App\Console\Commands;

use App\Models\Manga\MangaChapter;
use App\Models\Manga\MangaComic;
use App\Models\SourceChapter;
use App\Models\SourceComic;
use App\Models\SourceImage;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;

class VerifyChapter extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'verify.chapter';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '自动审核章节';

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
        $result = DB::select('SELECT ComicPublishQueue() as result')[0];
        $path = public_path('') . '/';
        if ($result->result == null) return;
        $comics = [];
        foreach (explode(',', $result->result) as $item) {
            $tmp = explode('-', $item);
            $comics[] = ['source_comic_id' => $tmp[0], 'comic_id' => $tmp[1], 'tab' => $tmp[2]];
        }

        if (empty($comics)) return;
        $time = time();
        foreach ($comics as $comic) {
            $chapters = SourceChapter::where('comic_id', $comic['source_comic_id'])->where('is_free', 0)->where('status', 0)->get()->toArray();
            if (empty($chapters)) continue;
            $MC = (new MangaChapter())->setTable($comic['tab']);
            foreach ($chapters as $chapter) {
                $image = SourceImage::where('chapter_id', $chapter['id'])->where('state', 1)->select('id','images')->first();
                if (!$image) {
                    continue;
                }
                $new_images = [];
                foreach ($image->images as $img) {
                    if (is_file($path . $img)) {
                        $imageInfo = getimagesize($path . $img);
                        $width = $imageInfo[0] ?? 0;
                        $height = $imageInfo[1] ?? 0;
                        $new_images[] = ['file' => $img, 'w' => $width, 'h' => $height];
                    }
                }

                $MC->insert([
                    'comic_id' => $comic['comic_id'],
                    'source_comic_id' => $comic['source_comic_id'],
                    'source_chapter_id' => $chapter['id'],
                    'source_image_id' => $image->id,
                    'title' => $chapter['title'],
                    'display_order' => $chapter['sort'],
                    'total_photos' => count($image->images),
                    'images' => json_encode($image->images),
                    'image_list' => json_encode($new_images),
                    'created_at' => strtotime($chapter['created_at']),
                    'updated_at' => $time,
                    'deleted_at' => 0
                ]);

                SourceChapter::where('id', $chapter['id'])->update(['status' => 1]);
            }
            $chapter_total = $MC->where('comic_id', $comic['comic_id'])->count();
            $last = $MC->where('comic_id', $comic['comic_id'])->orderBy('display_order', 'DESC')->first();

            $data = [
                'total_chapters' => $chapter_total,
                'free_chapters' => $chapter_total,
                'last_chapter_id' => $last ? $last->id : 0,
                'last_update_time' => $time,
                'last_chapter_time' => $last ? $last->created_at : $time,
            ];
            MangaComic::where('comic_id', $comic['comic_id'])->update($data);
        }
    }

}
