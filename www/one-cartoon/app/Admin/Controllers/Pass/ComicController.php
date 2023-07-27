<?php


namespace App\Admin\Controllers\Pass;


use App\Admin\Controllers\CommonController;
use App\Models\Manga\MangaComic;
use DLP\Tool\Assistant;
use DLP\Widget\Plane;
use Encore\Admin\Admin;
use Encore\Admin\Controllers\AdminController;
use Encore\Admin\Form;
use Encore\Admin\Grid;
use Encore\Admin\Layout\Content;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;

class ComicController extends AdminController
{
    protected $title = '漫画';

    protected function grid()
    {
        $grid = new Grid(new MangaComic());
        $grid->model()->orderBy('created_at','DESC');

        $grid->column('comic_id', __('ID'))->sortable();
        $grid->column('name', '标题');
        Admin::script("_component.imgDelay('.cover',{zoom:true});");
        $grid->column('vertical_cover', '封面')->display(function ($v){
            if(preg_match("/^http/",$v)){
                $url = $v;
            }else {
                $url = env('IMG_DOMAIN') . '/' . $v;
            }
            return "<div style='width:100px;height:60px'><img data-src='{$url}' class='cover img img-thumbnail' style='max-width:100px;height: 100%;' /></div>";
        });
        $grid->column('sort_id', '分类');
        $grid->column('tags', '标签');
        $grid->column('is_finish', '连载状态')->using([0 => '连载中', 1 => '完结']);
        $grid->column('source', '采集源')->display(function ($v) {
            if ($v == 1) {
                return "快看";
            } else {
                return "腾讯";
            }
        });
        $grid->column('total_chapters', '章节数量');
        $grid->column('created_at', '创建时间')->sortable();
        $grid->column('updated_at', '更新时间')->sortable();
        $grid->column('章节列表')->display(function () {
            return "<a href='/admin/chapter?comic_id={$this->comic_id}' target='_blank'>章节列表</a>";
        });
        /*配置*/
        $grid->disableCreateButton();
        $grid->disableExport();
        $grid->disableRowSelector();
        $grid->actions(function ($actions){
            $actions->disableView();
            $actions->disableEdit();
            $actions->disableDelete();
            $url = CommonController::getCurrentUrl();
            $actions->add(Plane::rowAction('编辑', $url."/{$actions->row->comic_id}/edit", ['url' => $url."/{$actions->row->comic_id}"]));
        });

        /*查询匹配*/
        $grid->filter(function ($filter) {
            $filter->like('title', '标题');
            $filter->equal('is_finish', '连载状态')->select([0 => '连载中', 1 => '完结']);
            $filter->between('created_at', '创建时间')->datetime();
        });
        return $grid;
    }

    public function edit($id, Content $content)
    {
        $content = $content
            ->body($this->form($id)->edit($id));
        return Plane::form($content);
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

    protected function form($id='')
    {
        $form = new Form(new MangaComic());
        /*配置*/
        CommonController::disableDetailConf($form);
        $form->builder()->setTitle('编辑漫画');
        $form->display('comic_id', 'ID');
        $form->text('name', '标题')->required();
        $form->radio('status', '上下架')->options([1 => '上架中', 2 => '下架'])->default(1);
        $form->image('vertical_cover', '封面')->options(['maxFileSize' => 1024]);
        $form->textarea('description','内容简介');
        return $form;
    }
}
