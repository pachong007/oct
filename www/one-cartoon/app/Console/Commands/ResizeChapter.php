<?php

namespace App\Console\Commands;

use App\Models\Manga\MangaChapter;
use App\Models\Manga\MangaComic;
use App\Models\SourceChapter;
use App\Models\SourceImage;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;

class ResizeChapter extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'resize.chapter {--p=value}';

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
        $pick = $this->option('p');
        $path = public_path('').'/';
        $end = $pick*32;
        $start = $end - 32;
        for($start;$start<=$end;$start++) {
            $MC = (new MangaChapter())->setTable($start);
            $chapters = $MC->get();
            foreach ($chapters as $chapter){
                if($chapter->image_list){
                    continue;
                }
                $images = (array)json_decode($chapter->images,true);
                $new_images = [];
                foreach ($images as $image){
                    if(is_file($path.$image)) {
                        $imageInfo = getimagesize($path . $image);
                        $width = $imageInfo[0] ?? 0;
                        $height = $imageInfo[1] ?? 0;
                        $new_images[] = ['file'=>$image,'w'=>$width,'h'=>$height];
                    }
                }
                $MC->where('id',$chapter->id)->update(['image_list'=>json_encode($new_images)]);
            }
        }
    }

}
