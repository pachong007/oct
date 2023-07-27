<?php

namespace App\Console\Commands;


use App\Models\Manga\MangaComic;
use App\Models\SourceChapter;
use App\Models\SourceComic;
use App\Models\SourceImage;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Redis;

class Republish extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'republish';

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
        $redis = Redis::connection('kk');
        for($i=0;$i<=1000;$i++) {
            $id = $redis->rpop("source:republish:comic");
            if ($id == false || $id == "") {
                return;
            }
            $id = (int)$id;
            $sourceComic = SourceComic::where('id', $id)->first();
            if (!$sourceComic) continue;

            $comic = MangaComic::where('source_comic_id', $id)->where('source',1)->first();
            if (!$comic) continue;
            MangaComic::where('source_comic_id', $id)->where('source',1)->update([
                'vertical_cover'=>$sourceComic->cover,
                'horizontal_cover'=>$sourceComic->cover_h
            ]);
        }
    }

}
