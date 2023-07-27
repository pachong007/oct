<?php

namespace App\Console\Commands;


use App\Models\Manga\MangaComic;
use App\Models\Pri\Category;
use App\Models\Pri\Tag;
use App\Models\SourceComic;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Redis;

class RefreshKK extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'refresh.kk';

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
        $redis = Redis::connection("kk");
        $retryQueue = $redis->lrange("source:comic:task",0,-1);
        $redis->del("source:comic:task");
        $comics = SourceComic::where('source', 1)
            ->where('created_at', '>', date('2023-07-13 22:00:00'))
            ->where('chapter_count',0)->inRandomOrder()->take(100)->pluck('id')->all();
        SourceComic::whereIn('id',$comics)->update(['retry'=>7]);
        $retryQueue =  array_unique(array_merge($retryQueue,$comics));
        shuffle($retryQueue);
        foreach ($retryQueue as $id){
            $redis->rpush("source:comic:task",$id);
        }
    }

}
