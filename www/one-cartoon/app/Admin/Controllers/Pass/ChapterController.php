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
        $M = (new MangaChapter())->setTable($comic_id%256);
        $grid = new Grid($M);
        $grid->model()->orderBy('display_order','DESC');
        $grid->column('id', __('ID'))->sortable();
        $grid->column('title', '标题');
        $grid->column('display_order', '排序')->sortable();
        $grid->column('created_at', '创建时间')->sortable();
        $grid->column('updated_at', '更新时间')->sortable();
        /*配置*/
        $grid->disableCreateButton();
        $grid->disableExport();
        $grid->disableRowSelector();
        $grid->actions(function ($actions)use($comic_id){
            $actions->disableView();
            $actions->disableEdit();
            $actions->disableDelete();
            $url = CommonController::getCurrentUrl();
            $actions->add(Plane::rowAction('资源信息', $url."/{$actions->row->id}/edit?comic_id={$comic_id}", ['url' => $url."/{$actions->row->id}"]));
        });

        /*查询匹配*/
        $grid->filter(function ($filter) {
            $filter->equal('comic_id', '漫画id');
            $filter->like('title', '标题');
            $filter->equal('is_free', '付费状态')->select([0 => '免费', 1 => '收费']);
            $filter->between('create_at', '创建时间')->datetime();
        });
        return $grid;
    }

    public function edit($id, Content $content)
    {
        $comic_id = $_GET['comic_id'];
        $M = (new MangaChapter())->setTable($comic_id%256)->where('id',$id)->first();
        $W = new Wing();
        $W->fileInput('image')
            ->label('图片')
            ->settings([
                'maxFileCount' => 10,
                'maxFileSize' => 800 //单图限制800kb
            ])
            ->initialPreview(['files' => (array)json_decode($M->images), 'url' => env("IMG_DOMAIN")."/"]);
        return $W;
    }

    public function update($id)
    {
        $request = Request::capture();
        $data = $request->all();
        try {
            DB::beginTransaction();
            if (!$data['title']) throw new \Exception('标题参数必填');

        } catch (\Exception $e) {
            DB::rollBack();
            return Assistant::result(false, $e->getMessage());
        }
        DB::commit();
        return Assistant::result(true);
    }

    private function cascadeExampleData()
    {
        return [
            ["key" => 3, "val" => "基本", "nodes" => [
                ["key" => 6895, "val" => "可播放"],
                ["key" => 6896, "val" => "可下載"],
                ["key" => 6863, "val" => "含字幕"],
                ["key" => 6855, "val" => "含預覽圖"],
                ["key" => 6862, "val" => "含預覽視頻"]
            ]],
            ["key" => 10, "val" => "年份", "nodes" => [
                ["key" => 6868, "val" => "2021"],
                ["key" => 6867, "val" => "2020"],
                ["key" => 6866, "val" => "2019"],
                ["key" => 14, "val" => "2018"],
                ["key" => 15, "val" => "2017"],
                ["key" => 16, "val" => "2016"],
                ["key" => 6897, "val" => "2015"],
                ["key" => 18, "val" => "2014"],
                ["key" => 19, "val" => "2013"],
                ["key" => 20, "val" => "2012"],
                ["key" => 6898, "val" => "2011"],
                ["key" => 22, "val" => "2010"]
            ]],
            ["key" => 32, "val" => "主題", "nodes" => [
                ["key" => 443, "val" => "按摩油"],
                ["key" => 444, "val" => "成熟妈妈"],
                ["key" => 445, "val" => "綠帽男"]
            ]],
            ["key" => 86, "val" => "角色", "nodes" => [
                ["key" => 614, "val" => "空姐"],
                ["key" => 616, "val" => "繼女"],
                ["key" => 587, "val" => "少女"],
                ["key" => 617, "val" => "角色扮演"],
                ["key" => 588, "val" => "醫生\/護士"],
                ["key" => 618, "val" => "女友"],
                ["key" => 589, "val" => "性愛專家"],
                ["key" => 619, "val" => "女神"],
                ["key" => 377, "val" => "熟女"],
                ["key" => 590, "val" => "媽媽"],
                ["key" => 591, "val" => "女抖S"],
                ["key" => 592, "val" => "抖M"],
                ["key" => 593, "val" => "妻子"],
                ["key" => 601, "val" => "女傭"],
                ["key" => 603, "val" => "水管工"],
                ["key" => 604, "val" => "警察"],
                ["key" => 402, "val" => "女學生"],
                ["key" => 607, "val" => "女戰士"],
                ["key" => 608, "val" => "特務"],
                ["key" => 609, "val" => "老師"],
                ["key" => 610, "val" => "女服務員"],
                ["key" => 110, "val" => "秘書"]
            ]],
            ["key" => 133, "val" => "服裝", "nodes" => [
                ["key" => 550, "val" => "吊襪腰帶"],
                ["key" => 551, "val" => "背心"],
                ["key" => 552, "val" => "裙子"],
                ["key" => 553, "val" => "短裙"],
                ["key" => 554, "val" => "短褲"],
                ["key" => 555, "val" => "綁腿"],
                ["key" => 556, "val" => "太陽鏡"],
                ["key" => 557, "val" => "帽子"],
                ["key" => 558, "val" => "襯衫"],
                ["key" => 559, "val" => "睡衣"],
                ["key" => 560, "val" => "內褲"],
                ["key" => 561, "val" => "牛仔短褲"],
                ["key" => 542, "val" => "蕾絲"],
                ["key" => 543, "val" => "高跟鞋"],
                ["key" => 544, "val" => "絲襪"],
                ["key" => 545, "val" => "吊帶"],
                ["key" => 416, "val" => "比基尼"],
                ["key" => 691, "val" => "眼鏡"],
                ["key" => 546, "val" => "眼罩"],
                ["key" => 137, "val" => "制服"],
                ["key" => 547, "val" => "牛仔褲"],
                ["key" => 158, "val" => "緊身衣"]
            ]]
        ];
    }
}
