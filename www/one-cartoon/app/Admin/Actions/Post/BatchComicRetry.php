<?php

namespace App\Admin\Actions\Post;

use App\Models\SourceChapter;
use App\Models\SourceComic;
use App\Models\SourceImage;
use Encore\Admin\Actions\BatchAction;
use Illuminate\Database\Eloquent\Collection;
use Illuminate\Support\Facades\Redis;

class BatchComicRetry extends BatchAction
{
    public $name = '漫画批量重抓';

    public function handle(Collection $collection)
    {
        foreach ($collection as $model) {
            $id = $model->id;
            if($model->source == 1){
                $redis = Redis::connection('kk');
            }else{
                $redis = Redis::connection('tx');
            }
            SourceImage::where('comic_id',$id)->delete();
            SourceChapter::where('comic_id',$id)->delete();
            SourceComic::where('id',$id)->update(['chapter_count'=>0,'chapter_pick'=>0,'retry'=>7]);
            $redis->lpush("source:comic:task",$id);
        }

        return $this->response()->success('进入重抓队列...')->refresh();
    }

}
