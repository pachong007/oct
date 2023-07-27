<?php

namespace App\Console\Commands;


use App\Models\SourceComic;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;


class CountChapterCount extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'count.chapter';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '统计章节数量';

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
        $limit = 1000;
        while (true) {
            $sources = SourceComic::offset($page * $limit)->where('last_chapter_update_at', '>' ,date('Y-m-d',strtotime('2 days ago')))->limit($limit)->orderBy('created_at', 'ASC')->pluck('id')->all();
            if(empty($sources))break;
            $page++;
            foreach ($sources as $id) {
                $result = DB::select("SELECT CountChapterDone({$id}) as result")[0];
                SourceComic::where('id',$id)->update(['chapter_count_download'=>$result->result]);
            }
        }
    }

}
