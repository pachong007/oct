<?php


namespace App\Admin\Controllers\Pass;


use App\Admin\Controllers\CommonController;
use App\Models\Comic;
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

    /**
     * CREATE USER 'root'@'172.18.0.6' IDENTIFIED BY 'docker@6603';
     * GRANT ALL PRIVILEGES ON *.* TO 'root'@'172.18.0.6' WITH GRANT OPTION;
     * FLUSH PRIVILEGES;
     * @return Grid
     */
    protected function grid()
    {
        $db = 'fxkexie_cn';
        config(["database.connections.mysql_${db}" => [
            'driver'    => 'mysql',
            'host'      => '107.148.191.71',
            'database'  => $db,
            'username'  => 'root',
            'password'  => 'docker@6603',
            'charset'   => 'utf8',
            'collation' => 'utf8_unicode_ci'
        ]]);
        $comic = new Comic();
        $comic->setConnection("mysql_${db}");
        $grid = new Grid($comic);

        $grid->column('id', __('ID'))->sortable();
        $grid->column('name', '标题');
        Admin::script("_component.imgDelay('.cover',{zoom:true});");
        $grid->column('pic', '封面')->display(function ($v){
            if(preg_match("/^http/",$v)){
                $url = $v;
            }else {
                $url = env('IMG_DOMAIN') . '/' . $v;
            }
            return "<div style='width:100px;height:60px'><img data-src='{$url}' class='cover img img-thumbnail' style='max-width:100px;height: 100%;' /></div>";
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
        });
        return $grid;
    }

    protected function form($id='')
    {

    }
}
