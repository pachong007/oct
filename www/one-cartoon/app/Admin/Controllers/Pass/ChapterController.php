<?php


namespace App\Admin\Controllers\Pass;


use App\Admin\Controllers\CommonController;
use App\Models\Chapter;
use App\Models\Comic;
use App\Models\Manga\MangaChapter;
use DLP\Assembly\Wing;
use DLP\Tool\Assistant;
use DLP\Widget\Plane;
use Encore\Admin\Controllers\AdminController;
use Encore\Admin\Form;
use Encore\Admin\Grid;
use Encore\Admin\Layout\Content;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;

class ChapterController extends AdminController
{
    protected $title = '章节';

    protected function grid()
    {
        $comic_id = $_GET['comic_id'];

    }

}
