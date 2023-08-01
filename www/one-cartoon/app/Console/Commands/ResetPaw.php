<?php

namespace App\Console\Commands;


use App\Models\SourceChapter;
use App\Models\SourceImage;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Redis;

class ResetPaw extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'reset.paw';

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
        $this->reset('kk',1);
    }

    private function reset($source,$sourceId)
    {
        $redis = Redis::connection($source);

        $retry = $redis->lrange("source:comic:retry:task",0,-1);
        $redis->del("source:comic:retry:task");
        $list = $redis->lrange("source:comic:task",0,-1);
        $redis->del("source:comic:task");
        foreach ($retry as $r){
            $list[] = $r;
        }
        $newList = array_unique($list);
        foreach ($newList as $id){
            $redis->rpush("source:comic:task",$id);
        }


        $retry = $redis->lrange("source:comic:retry:chapter",0,-1);
        $redis->del("source:comic:retry:chapter");
        $list = $redis->lrange("source:comic:chapter",0,-1);
        $redis->del("source:comic:chapter");
        foreach ($retry as $r){
            $list[] = $r;
        }
exit;
        $chapters = SourceChapter::where('source',$sourceId)->where('status',0)->where('created_at','>',date('Y-m-d H:i:s',date('Y-m-d',strtotime('12 days ago'))))->get();
        foreach ($chapters as $chapter){
            if($chapter->is_free == 1) continue;
            if($chapter->retry > 9) continue;
            $img = SourceImage::where('chapter_id',$chapter->id)->first();
            if(!$img || $img->state == 0){
                $list[] = $chapter->id;
            }
        }

        $newList = array_unique($list);
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
}
