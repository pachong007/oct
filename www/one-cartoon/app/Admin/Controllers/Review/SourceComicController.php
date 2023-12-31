<?php
namespace App\Admin\Controllers\Review;


use App\Admin\Actions\Post\BatchComicRetry;
use App\Admin\Controllers\CommonController;
use App\Models\Comic;
use App\Models\SourceChapter;
use App\Models\SourceComic;
use App\Models\SourceImage;
use DLP\Layer\Dialog;
use DLP\Assembly\Wing;
use DLP\Tool\Assistant;
use DLP\Widget\Plane;
use Encore\Admin\Admin;
use Encore\Admin\Controllers\AdminController;
use Encore\Admin\Form;
use Encore\Admin\Grid;
use Encore\Admin\Layout\Content;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Redis;

class SourceComicController extends AdminController
{
    protected $title = '漫画审核';

    protected function grid()
    {
        $grid = new Grid(new SourceComic());
        $grid->model()->orderBy('id', 'DESC');

        $grid->column('id', __('ID'))->sortable();
        $grid->column('title', '标题')->width(220);
        Admin::script("_component.imgDelay('.cover',{zoom:true});");
        $grid->column('cover', '封面')->display(function ($v) {
            if (preg_match("/^http/", $v)) {
                $url = $v;
            } else {
                $url = env('IMG_DOMAIN') . '/' . $v;
            }
            return "<div style='width:100px;height:60px'><img data-src='{$url}' class='cover img img-thumbnail' style='max-width:100px;height: 100%;' /></div>";
        });
        Admin::script("_component.imgDelay('.cover2',{zoom:true});");
        $grid->column('cover_h', '横板封面')->display(function ($v) {
            if (preg_match("/^http/", $v)) {
                $url = $v;
            } else {
                $url = env('IMG_DOMAIN') . '/' . $v;
            }
            return "<div style='width:100px;height:60px'><img data-src='{$url}' class='cover2 img img-thumbnail' style='max-width:100px;height: 100%;' /></div>";
        });
        $grid->column('status', '审核')->using([0 => '待审核', 1 => '通过',2=>'待处理',3=>'同名暂存'])->dot([0 => 'info', 1 => 'success',2=>'warning',3=>'danger']);
        $grid->column('category', '分类');
        $grid->column('region', '地区');
        $grid->column('is_free', '付费状态')->using([0 => '免费', 1 => '收费'])->dot([0 => 'success',1=>'danger']);
        $grid->column('is_finish', '连载状态')->using([0 => '连载中', 1 => '完结']);
        $grid->column('chasm', '章节断层')->using([0 => '连续', 1 => '断层']);
        $grid->column('source', '采集源')->display(function ($v) {
            $source_url = $this->source_url;
            if ($v == 1) {
                return "<a href='$source_url' target='_blank'>快看</a>";
            } else {
                return "<a href='$source_url' target='_blank'>腾讯</a>";
            }
        });
        $grid->column('chapter_count', '章节数量');
        $grid->column('created_at', '创建时间')->sortable();
        $grid->column('updated_at', '更新时间')->sortable();
        $grid->column('章节列表')->display(function () {
            return "<a href='/admin/source_chapter?comic_id={$this->id}' target='_blank'>章节列表</a>";
        });
        /*配置*/
        $grid->disableCreateButton();
        $grid->disableExport();
        $grid->actions(function ($actions) {
            $actions->disableView();
            $actions->disableEdit();
            $actions->disableDelete();
            $url = CommonController::getCurrentUrl();
            $dialog = (new Dialog())->info('已经存在相同漫画 是否进行覆盖操作?')->button('测试');
            $actions->add(Plane::rowAction('发布漫画',
                $url . "/{$actions->row->id}/edit",
                ['url' => $url . "/{$actions->row->id}", 'callback'=> <<<EOF
            function(response){
                let tip = document.createElement('div');
                tip.style = "display: flex;align-items: center;justify-content: center;height: 30px;";
                tip.innerText = "已经存在相同漫画 是否进行覆盖操作?";
                {$dialog}
            }
EOF
            ]));
            $actions->add(Plane::rowAction('数据源',
                rtrim(config('app.url'), '/') . '/admin/source.comic/'.$actions->row->id)
                ->withoutBind());
        });

        $grid->tools(function ($tools) {
            $tools->batch(function ($batch) {
                $batch->disableDelete();
                $batch->add(new BatchComicRetry());
            });
        });

        /*查询匹配*/
        $grid->filter(function ($filter) {
            $filter->like('title', '标题');
            $filter->equal('region', '地区查询');
            $filter->equal('source', '采集源')->select([1 => '快看', 2 => '腾讯']);
            $filter->equal('is_free', '付费状态')->select([0 => '免费', 1 => '收费']);
            $filter->equal('is_finish', '连载状态')->select([0 => '连载中', 1 => '完结']);
            $filter->equal('status', '审核状态')->select([0 => '待审核', 1 => '通过',2=>'待处理']);
            $filter->equal('chasm', '章节断层')->select([0 => '连续', 1 => '断层']);
            $filter->group('chapter_count','章节数量', function ($group) {
                $group->gt('大于');
                $group->lt('小于');
                $group->nlt('不小于');
                $group->ngt('不大于');
                $group->equal('等于');
            });
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
            throw new \Exception('漫画已经发布');
            if (!$data['title']) throw new \Exception('标题参数必填');
            $sourceComic = SourceComic::where('id', $id)->first();
            if($sourceComic->status == 1) throw new \Exception('漫画已经发布');
            Comic::insert([
                "source_comic_id" => $sourceComic->id,
                "cover" => $sourceComic->cover,
                "title" => $data['title'],
                "author" => $sourceComic->author,
                "label" => json_encode($sourceComic->label),
                "category" => $sourceComic->category,
                "region" => $sourceComic->region,
                "chapter_count" => 0,
                "like" => $sourceComic->like,
                "popularity" => $sourceComic->popularity,
                "is_finish" => $sourceComic->is_finish,
                "description" => $sourceComic->description
            ]);
            SourceComic::where('id',$sourceComic->id)->update(['status'=>1]);
        } catch (\Exception $e) {
            return Assistant::result(false, $e->getMessage());
        }
        return Assistant::result(true);
    }

    protected function form($id = '')
    {
        $form = new Form(new SourceComic());
        /*配置*/
        CommonController::disableDetailConf($form);
        $form->builder()->setTitle('审核漫画');
        $form->display('id', 'ID');
        $form->text('title', '标题')->required();
        $form->textarea('title', '标题');
        $form->image('cover', '封面')->options(['maxFileSize' => 1024]);
        return $form;
    }

    public function sourceComic($id)
    {
        $W = new Wing();
        $comic = SourceComic::where('id', $id)->first();
        $W->display('id')->label('ID')->value($comic->id);
        $W->text('title')->label('标题')->value($comic->title);
        $W->text('source_url')->label('源地址')->value($comic->source_url);
        $W->datepicker('create_at')->label('创建时间')->value($comic->created_at);
        $source_data = $comic->source_data;
        $W->textarea('source')->label('源信息')->rows(13)->value($source_data);
        $W->section(function (Wing $W)use($id){
            $W->button('删除所有章节')->setStyle(['width'=>'200px'])->bindDialog(function (Dialog $dialog)use($id){
                $dialog->info('确认删除资源？')->button('确认');
            },['url'=>rtrim(config('app.url'), '/') . '/admin/source.comic/del/'.$id])->color('red');

            $W->button('重抓所有章节')->setStyle(['width'=>'200px'])->bindDialog(function (Dialog $dialog)use($id){
                $dialog->info('重抓所有章节？')->button('确认');
            },['url'=>rtrim(config('app.url'), '/') . '/admin/source.comic/retry/'.$id])->color('blue');
        });
        return $W->form();
    }

    public function del($id)
    {
        $chapters = SourceChapter::where('comic_id',$id)->get();
        foreach ($chapters as $chapter){
            SourceImage::where('comic_id',$id)->where('chapter_id',$chapter->id)->delete();
        }
        SourceChapter::where('comic_id',$id)->delete();
        SourceComic::where('id',$id)->update(['chapter_count'=>0,'chapter_pick'=>0]);
        return Assistant::result(true);
    }

    public function retry($id)
    {
        $comic = SourceComic::where('id',$id)->first();
        if($comic->source == 1){
            $redis = Redis::connection('kk');
        }else{
            $redis = Redis::connection('tx');
        }
        SourceImage::where('comic_id',$id)->delete();
        SourceChapter::where('comic_id',$id)->delete();
        SourceComic::where('id',$id)->update(['chapter_count'=>0,'chapter_pick'=>0]);
        $redis->lpush("source:comic:task",$id);
        return Assistant::result(true);
    }
}
