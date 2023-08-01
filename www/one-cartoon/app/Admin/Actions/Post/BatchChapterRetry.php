<?php

namespace App\Admin\Actions\Post;

use App\Models\SourceChapter;
use App\Models\SourceImage;
use Encore\Admin\Actions\BatchAction;
use Illuminate\Database\Eloquent\Collection;
use Illuminate\Support\Facades\Redis;

class BatchChapterRetry extends BatchAction
{
    public $name = '章节批量重抓';

    public function handle(Collection $collection)
    {
        foreach ($collection as $model) {
            $id = $model->id;
            SourceImage::where('chapter_id',$id)->update(['state'=>0,'images'=>json_encode([])]);
            if($model->source == 1){
                $redis = Redis::connection('kk');
            }else{
                $redis = Redis::connection('tx');
            }
            $redis->lpush("source:comic:chapter",$id);
            SourceChapter::where('id',$id)->update(['retry'=>7]);
        }

        return $this->response()->success('进入重抓队列...')->refresh();
    }

}
